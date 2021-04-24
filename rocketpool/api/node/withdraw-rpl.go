package node

import (
    "context"
    "math/big"

    "github.com/rocket-pool/rocketpool-go/node"
    "github.com/rocket-pool/rocketpool-go/settings/protocol"
    "github.com/urfave/cli"
    "golang.org/x/sync/errgroup"

    "github.com/rocket-pool/smartnode/shared/services"
    "github.com/rocket-pool/smartnode/shared/types/api"
)


func canNodeWithdrawRpl(c *cli.Context, amountWei *big.Int) (*api.CanNodeWithdrawRplResponse, error) {

    // Get services
    if err := services.RequireNodeRegistered(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    ec, err := services.GetEthClient(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.CanNodeWithdrawRplResponse{}

    // Get node account
    nodeAccount, err := w.GetNodeAccount()
    if err != nil {
        return nil, err
    }

    // Data
    var wg errgroup.Group
    var rplStake *big.Int
    var minimumRplStake *big.Int
    var currentBlock uint64
    var rplStakedBlock uint64
    var withdrawalDelay uint64

    // Get RPL stake
    wg.Go(func() error {
        var err error
        rplStake, err = node.GetNodeRPLStake(rp, nodeAccount.Address, nil)
        return err
    })

    // Get minimum RPL stake
    wg.Go(func() error {
        var err error
        minimumRplStake, err = node.GetNodeMinimumRPLStake(rp, nodeAccount.Address, nil)
        return err
    })

    // Get current block
    wg.Go(func() error {
        header, err := ec.HeaderByNumber(context.Background(), nil)
        if err == nil {
            currentBlock = header.Number.Uint64()
        }
        return err
    })

    // Get RPL staked block
    wg.Go(func() error {
        var err error
        rplStakedBlock, err = node.GetNodeRPLStakedBlock(rp, nodeAccount.Address, nil)
        return err
    })

    // Get withdrawal delay
    wg.Go(func() error {
        var err error
        withdrawalDelay, err = protocol.GetRewardsClaimIntervalBlocks(rp, nil)
        return err
    })

    // Wait for data
    if err := wg.Wait(); err != nil {
        return nil, err
    }

    // Check data
    var remainingRplStake big.Int
    remainingRplStake.Sub(rplStake, amountWei)
    response.InsufficientBalance = (amountWei.Cmp(rplStake) > 0)
    response.MinipoolsUndercollateralized = (remainingRplStake.Cmp(minimumRplStake) < 0)
    response.WithdrawalDelayActive = ((currentBlock - rplStakedBlock) < withdrawalDelay)

    // Update & return response
    response.CanWithdraw = !(response.InsufficientBalance || response.MinipoolsUndercollateralized || response.WithdrawalDelayActive)
    return &response, nil

}


func nodeWithdrawRpl(c *cli.Context, amountWei *big.Int) (*api.NodeWithdrawRplResponse, error) {

    // Get services
    if err := services.RequireNodeRegistered(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.NodeWithdrawRplResponse{}

    // Get transactor
    opts, err := w.GetNodeAccountTransactor()
    if err != nil {
        return nil, err
    }

    // Withdraw RPL
    hash, err := node.WithdrawRPL(rp, amountWei, opts)
    if err != nil {
        return nil, err
    }
    response.TxHash = hash

    // Return response
    return &response, nil

}

