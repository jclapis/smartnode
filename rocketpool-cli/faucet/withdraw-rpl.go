package faucet

import (
    "fmt"

    "github.com/rocket-pool/rocketpool-go/utils/eth"
    "github.com/urfave/cli"

    "github.com/rocket-pool/smartnode/shared/services/rocketpool"
    "github.com/rocket-pool/smartnode/shared/utils/math"
)


func withdrawRpl(c *cli.Context) error {

    // Get RP client
    rp, err := rocketpool.NewClientFromCtx(c)
    if err != nil { return err }
    defer rp.Close()

    // Check RPL can be withdrawn
    canWithdraw, err := rp.CanFaucetWithdrawRpl()
    if err != nil {
        return err
    }
    if !canWithdraw.CanWithdraw {
        fmt.Println("Cannot withdraw RPL from the faucet:")
        if canWithdraw.InsufficientFaucetBalance {
            fmt.Println("The faucet does not have any RPL for withdrawal")
        }
        if canWithdraw.InsufficientAllowance {
            fmt.Println("You don't have any allowance remaining for the withdrawal period")
        }
        if canWithdraw.InsufficientNodeBalance {
            fmt.Println("You don't have enough GoETH to pay the faucet withdrawal fee")
        }
        return nil
    }

    // Withdraw RPL
    response, err := rp.FaucetWithdrawRpl()
    if err != nil {
        return err
    }

    // Log & return
    fmt.Printf("Successfully withdrew %.6f RPL from the faucet.\n", math.RoundDown(eth.WeiToEth(response.Amount), 6))
    return nil

}

