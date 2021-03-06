package node

import (
    "fmt"
    "math/big"
    "strconv"

    "github.com/rocket-pool/rocketpool-go/utils/eth"
    "github.com/urfave/cli"

    "github.com/rocket-pool/smartnode/shared/services/rocketpool"
    cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
    "github.com/rocket-pool/smartnode/shared/utils/math"
)


func nodeWithdrawRpl(c *cli.Context) error {

    // Get RP client
    rp, err := rocketpool.NewClientFromCtx(c)
    if err != nil { return err }
    defer rp.Close()

    // Get withdrawal mount
    var amountWei *big.Int
    if c.String("amount") == "max" {

        // Get node status
        status, err := rp.NodeStatus()
        if err != nil {
            return err
        }

        // Set amount to maximum withdrawable amount
        var maxAmount big.Int
        if status.RplStake.Cmp(status.MinimumRplStake) > 0 {
            maxAmount.Sub(status.RplStake, status.MinimumRplStake)
        }
        amountWei = &maxAmount

    } else if c.String("amount") != "" {

        // Parse amount
        withdrawalAmount, err := strconv.ParseFloat(c.String("amount"), 64)
        if err != nil {
            return fmt.Errorf("Invalid withdrawal amount '%s': %w", c.String("amount"), err)
        }
        amountWei = eth.EthToWei(withdrawalAmount)

    } else {

        // Get node status
        status, err := rp.NodeStatus()
        if err != nil {
            return err
        }

        // Get maximum withdrawable amount
        var maxAmount big.Int
        if status.RplStake.Cmp(status.MinimumRplStake) > 0 {
            maxAmount.Sub(status.RplStake, status.MinimumRplStake)
        }

        // Prompt for maximum amount
        if cliutils.Confirm(fmt.Sprintf("Would you like to withdraw the maximum amount of staked RPL (%.6f RPL)?", math.RoundDown(eth.WeiToEth(&maxAmount), 6))) {
            amountWei = &maxAmount
        } else {

            // Prompt for custom amount
            inputAmount := cliutils.Prompt("Please enter an amount of staked RPL to withdraw:", "^\\d+(\\.\\d+)?$", "Invalid amount")
            withdrawalAmount, err := strconv.ParseFloat(inputAmount, 64)
            if err != nil {
                return fmt.Errorf("Invalid withdrawal amount '%s': %w", inputAmount, err)
            }
            amountWei = eth.EthToWei(withdrawalAmount)

        }

    }

    // Check RPL can be withdrawn
    canWithdraw, err := rp.CanNodeWithdrawRpl(amountWei)
    if err != nil {
        return err
    }
    if !canWithdraw.CanWithdraw {
        fmt.Println("Cannot withdraw staked RPL:")
        if canWithdraw.InsufficientBalance {
            fmt.Println("The node's staked RPL balance is insufficient.")
        }
        if canWithdraw.MinipoolsUndercollateralized {
            fmt.Println("Remaining staked RPL is not enough to collateralize the node's minipools.")
        }
        if canWithdraw.WithdrawalDelayActive {
            fmt.Println("The withdrawal delay period has not passed.")
        }
        return nil
    }

    // Prompt for confirmation
    if !(c.Bool("yes") || cliutils.Confirm(fmt.Sprintf("Are you sure you want to withdraw %.6f staked RPL? This may decrease your node's RPL rewards.", math.RoundDown(eth.WeiToEth(amountWei), 6)))) {
        fmt.Println("Cancelled.")
        return nil
    }

    // Withdraw RPL
    if _, err := rp.NodeWithdrawRpl(amountWei); err != nil {
        return err
    }

    // Log & return
    fmt.Printf("Successfully withdrew %.6f staked RPL.\n", math.RoundDown(eth.WeiToEth(amountWei), 6))
    return nil

}

