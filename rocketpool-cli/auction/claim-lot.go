package auction

import (
    "fmt"
    "strconv"

    "github.com/rocket-pool/rocketpool-go/utils/eth"
    "github.com/urfave/cli"

    "github.com/rocket-pool/smartnode/shared/services/rocketpool"
    "github.com/rocket-pool/smartnode/shared/types/api"
    cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
    "github.com/rocket-pool/smartnode/shared/utils/math"
)


func claimFromLot(c *cli.Context) error {

    // Get RP client
    rp, err := rocketpool.NewClientFromCtx(c)
    if err != nil { return err }
    defer rp.Close()

    // Get lot details
    lots, err := rp.AuctionLots()
    if err != nil {
        return err
    }

    // Get claimable lots
    claimableLots := []api.LotDetails{}
    for _, lot := range lots.Lots {
        if lot.ClaimAvailable {
            claimableLots = append(claimableLots, lot)
        }
    }

    // Check for claimable lots
    if len(claimableLots) == 0 {
        fmt.Println("No lots are available for RPL claims.")
        return nil
    }

    // Get selected lots
    var selectedLots []api.LotDetails
    if c.String("lot") == "all" {

        // Select all claimable lots
        selectedLots = claimableLots

    } else if c.String("lot") != "" {

        // Get selected lot index
        selectedIndex, err := strconv.ParseUint(c.String("lot"), 10, 64)
        if err != nil {
            return fmt.Errorf("Invalid lot ID '%s': %w", c.String("lot"), err)
        }

        // Get matching lot
        found := false
        for _, lot := range claimableLots {
            if lot.Details.Index == selectedIndex {
                selectedLots = []api.LotDetails{lot}
                found = true
                break
            }
        }
        if !found {
            return fmt.Errorf("Lot %d is not available for RPL claims.", selectedIndex)
        }

    } else {

        // Prompt for lot selection
        options := make([]string, len(claimableLots) + 1)
        options[0] = "All available lots"
        for li, lot := range claimableLots {
            options[li + 1] = fmt.Sprintf("lot %d (%.6f ETH bid @ %.6f ETH per RPL)", lot.Details.Index, math.RoundDown(eth.WeiToEth(lot.Details.AddressBidAmount), 6), math.RoundDown(eth.WeiToEth(lot.Details.CurrentPrice), 6))
        }
        selected, _ := cliutils.Select("Please select a lot to claim RPL from:", options)

        // Get lots
        if selected == 0 {
            selectedLots = claimableLots
        } else {
            selectedLots = []api.LotDetails{claimableLots[selected - 1]}
        }

    }

    // Claim RPL from lots
    for _, lot := range selectedLots {
        if _, err := rp.ClaimFromLot(lot.Details.Index); err != nil {
            fmt.Printf("Could not claim RPL from lot %d: %s.\n", lot.Details.Index, err)
        } else {
            fmt.Printf("Successfully claimed RPL from lot %d.\n", lot.Details.Index)
        }
    }

    // Return
    return nil

}

