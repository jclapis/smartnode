package node

import (
    "math/big"

    "github.com/rocket-pool/rocketpool-go/tokens"
    "github.com/urfave/cli"

    "github.com/rocket-pool/smartnode/shared/services"
    "github.com/rocket-pool/smartnode/shared/types/api"
)


func canNodeSwapRpl(c *cli.Context, amountWei *big.Int) (*api.CanNodeSwapRplResponse, error) {

    // Get services
    if err := services.RequireNodeWallet(c); err != nil { return nil, err }
    if err := services.RequireRocketStorage(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.CanNodeSwapRplResponse{}

    // Get node account
    nodeAccount, err := w.GetNodeAccount()
    if err != nil {
        return nil, err
    }

    // Check node fixed-supply RPL balance
    fixedSupplyRplBalance, err := tokens.GetFixedSupplyRPLBalance(rp, nodeAccount.Address, nil)
    if err != nil {
        return nil, err
    }
    response.InsufficientBalance = (amountWei.Cmp(fixedSupplyRplBalance) > 0)

    // Update & return response
    response.CanSwap = !response.InsufficientBalance
    return &response, nil

}


func nodeSwapRpl(c *cli.Context, amountWei *big.Int) (*api.NodeSwapRplResponse, error) {

    // Get services
    if err := services.RequireNodeWallet(c); err != nil { return nil, err }
    if err := services.RequireRocketStorage(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.NodeSwapRplResponse{}

    // Get RPL contract address
    rocketTokenRPLAddress, err := rp.GetAddress("rocketTokenRPL")
    if err != nil {
        return nil, err
    }

    // Approve fixed-supply RPL allowance
    if opts, err := w.GetNodeAccountTransactor(); err != nil {
        return nil, err
    } else if txReceipt, err := tokens.ApproveFixedSupplyRPL(rp, *rocketTokenRPLAddress, amountWei, opts); err != nil {
        return nil, err
    } else {
        response.ApproveTxHash = txReceipt.TxHash
    }

    // Swap fixed-supply RPL for RPL
    if opts, err := w.GetNodeAccountTransactor(); err != nil {
        return nil, err
    } else if txReceipt, err := tokens.SwapFixedSupplyRPLForRPL(rp, amountWei, opts); err != nil {
        return nil, err
    } else {
        response.SwapTxHash = txReceipt.TxHash
    }

    // Return response
    return &response, nil

}

