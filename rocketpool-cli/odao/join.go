package odao

import (
    "fmt"
    "math/big"

    "github.com/rocket-pool/rocketpool-go/utils/eth"
    "github.com/urfave/cli"

    "github.com/rocket-pool/smartnode/shared/services/rocketpool"
    cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
    "github.com/rocket-pool/smartnode/shared/utils/math"
)


func join(c *cli.Context) error {

    // Get RP client
    rp, err := rocketpool.NewClientFromCtx(c)
    if err != nil { return err }
    defer rp.Close()

    // Get node status
    status, err := rp.NodeStatus()
    if err != nil {
        return err
    }

    // Check for fixed-supply RPL balance
    if status.AccountBalances.FixedSupplyRPL.Cmp(big.NewInt(0)) > 0 {

        // Confirm swapping RPL
        if (c.Bool("swap") || cliutils.Confirm(fmt.Sprintf("The node has a balance of %.6f old RPL. Would you like to swap it for new RPL before transferring your bond?", math.RoundDown(eth.WeiToEth(status.AccountBalances.FixedSupplyRPL), 6)))) {

            // Swap RPL
            if _, err := rp.NodeSwapRpl(status.AccountBalances.FixedSupplyRPL); err != nil {
                return err
            }

            // log
            fmt.Printf("Successfully swapped %.6f old RPL for new RPL.\n", math.RoundDown(eth.WeiToEth(status.AccountBalances.FixedSupplyRPL), 6))
            fmt.Println("")

        }

    }

    // Check if node can join the oracle DAO
    canJoin, err := rp.CanJoinTNDAO()
    if err != nil {
        return err
    }
    if !canJoin.CanJoin {
        fmt.Println("Cannot join the oracle DAO:")
        if canJoin.ProposalExpired {
            fmt.Println("The proposal for you to join the oracle DAO does not exist or has expired.")
        }
        if canJoin.AlreadyMember {
            fmt.Println("The node is already a member of the oracle DAO.")
        }
        if canJoin.InsufficientRplBalance {
            fmt.Println("The node does not have enough RPL to pay the RPL bond.")
        }
        return nil
    }

    // Prompt for confirmation
    if !(c.Bool("yes") || cliutils.Confirm("Are you sure you want to join the oracle DAO? Your RPL bond will be locked until you leave.")) {
        fmt.Println("Cancelled.")
        return nil
    }

    // Join the oracle DAO
    if _, err := rp.JoinTNDAO(); err != nil {
        return err
    }

    // Log & return
    fmt.Println("Successfully joined the oracle DAO.")
    return nil

}

