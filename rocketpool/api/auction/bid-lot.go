package auction

import (
    "math/big"

    "github.com/rocket-pool/rocketpool-go/auction"
    "github.com/rocket-pool/rocketpool-go/settings/protocol"
    "github.com/urfave/cli"
    "golang.org/x/sync/errgroup"

    "github.com/rocket-pool/smartnode/shared/services"
    "github.com/rocket-pool/smartnode/shared/types/api"
)


func canBidOnLot(c *cli.Context, lotIndex uint64) (*api.CanBidOnLotResponse, error) {

    // Get services
    if err := services.RequireNodeWallet(c); err != nil { return nil, err }
    if err := services.RequireRocketStorage(c); err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.CanBidOnLotResponse{}

    // Sync
    var wg errgroup.Group

    // Check if lot exists
    wg.Go(func() error {
        lotExists, err := auction.GetLotExists(rp, lotIndex, nil)
        if err == nil {
            response.DoesNotExist = !lotExists
        }
        return err
    })

    // Check if lot bidding has ended
    wg.Go(func() error {
        biddingEnded, err := getLotBiddingEnded(rp, lotIndex)
        if err == nil {
            response.BiddingEnded = biddingEnded
        }
        return err
    })

    // Check lot remaining RPL amount
    wg.Go(func() error {
        remainingRpl, err := auction.GetLotRemainingRPLAmount(rp, lotIndex, nil)
        if err == nil {
            response.RPLExhausted = (remainingRpl.Cmp(big.NewInt(0)) == 0)
        }
        return err
    })

    // Check if lot bidding is enabled
    wg.Go(func() error {
        bidOnLotEnabled, err := protocol.GetBidOnLotEnabled(rp, nil)
        if err == nil {
            response.BidOnLotDisabled = !bidOnLotEnabled
        }
        return err
    })

    // Wait for data
    if err := wg.Wait(); err != nil {
        return nil, err
    }

    // Update & return response
    response.CanBid = !(response.DoesNotExist || response.BiddingEnded || response.RPLExhausted || response.BidOnLotDisabled)
    return &response, nil

}


func bidOnLot(c *cli.Context, lotIndex uint64, amountWei *big.Int) (*api.BidOnLotResponse, error) {

    // Get services
    if err := services.RequireNodeWallet(c); err != nil { return nil, err }
    if err := services.RequireRocketStorage(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.BidOnLotResponse{}

    // Get transactor
    opts, err := w.GetNodeAccountTransactor()
    if err != nil {
        return nil, err
    }
    opts.Value = amountWei

    // Bid on lot
    hash, err := auction.PlaceBid(rp, lotIndex, opts)
    if err != nil {
        return nil, err
    }
    response.TxHash = hash

    // Return response
    return &response, nil

}

