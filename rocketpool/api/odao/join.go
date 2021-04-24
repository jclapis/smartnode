package odao

import (
    "math/big"

    "github.com/ethereum/go-ethereum/common"
    "github.com/rocket-pool/rocketpool-go/dao/trustednode"
    tndao "github.com/rocket-pool/rocketpool-go/dao/trustednode"
    tnsettings "github.com/rocket-pool/rocketpool-go/settings/trustednode"
    "github.com/rocket-pool/rocketpool-go/tokens"
    "github.com/urfave/cli"
    "golang.org/x/sync/errgroup"

    "github.com/rocket-pool/smartnode/shared/services"
    "github.com/rocket-pool/smartnode/shared/types/api"
)


func canJoin(c *cli.Context) (*api.CanJoinTNDAOResponse, error) {

    // Get services
    if err := services.RequireNodeRegistered(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.CanJoinTNDAOResponse{}

    // Get node account
    nodeAccount, err := w.GetNodeAccount()
    if err != nil {
        return nil, err
    }

    // Data
    var wg errgroup.Group
    var nodeRplBalance *big.Int
    var rplBondAmount *big.Int

    // Check proposal actionable status
    wg.Go(func() error {
        proposalActionable, err := getProposalIsActionable(rp, nodeAccount.Address, "invited")
        if err == nil {
            response.ProposalExpired = !proposalActionable
        }
        return err
    })

    // Check if already a member
    wg.Go(func() error {
        isMember, err := tndao.GetMemberExists(rp, nodeAccount.Address, nil)
        if err == nil {
            response.AlreadyMember = isMember
        }
        return err
    })

    // Get node RPL balance
    wg.Go(func() error {
        var err error
        nodeRplBalance, err = tokens.GetRPLBalance(rp, nodeAccount.Address, nil)
        return err
    })

    // Get RPL bond amount
    wg.Go(func() error {
        var err error
        rplBondAmount, err = tnsettings.GetRPLBond(rp, nil)
        return err
    })

    // Wait for data
    if err := wg.Wait(); err != nil {
        return nil, err
    }

    // Check data
    response.InsufficientRplBalance = (nodeRplBalance.Cmp(rplBondAmount) < 0)

    // Update & return response
    response.CanJoin = !(response.ProposalExpired || response.AlreadyMember || response.InsufficientRplBalance)
    return &response, nil

}


func approveRpl(c *cli.Context) (*api.JoinTNDAOApproveResponse, error) {

    // Get services
    if err := services.RequireNodeRegistered(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.JoinTNDAOApproveResponse{}

    // Data
    var wg errgroup.Group
    var rocketDAONodeTrustedActionsAddress *common.Address
    var rplBondAmount *big.Int

    // Get oracle node actions contract address
    wg.Go(func() error {
        var err error
        rocketDAONodeTrustedActionsAddress, err = rp.GetAddress("rocketDAONodeTrustedActions")
        return err
    })

    // Get RPL bond amount
    wg.Go(func() error {
        var err error
        rplBondAmount, err = tnsettings.GetRPLBond(rp, nil)
        return err
    })

    // Wait for data
    if err := wg.Wait(); err != nil {
        return nil, err
    }

    // Approve RPL allowance
    if opts, err := w.GetNodeAccountTransactor(); err != nil {
        return nil, err
    } else if hash, err := tokens.ApproveRPL(rp, *rocketDAONodeTrustedActionsAddress, rplBondAmount, opts); err != nil {
        return nil, err
    } else {
        response.ApproveTxHash = hash
    }

    // Return response
    return &response, nil

}


func waitForApprovalAndJoin(c *cli.Context, hash common.Hash) (*api.JoinTNDAOJoinResponse, error) {

    // Get services
    if err := services.RequireNodeRegistered(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Wait for the RPL approval TX to successfully get mined
    _, err = trustednode.WaitForTransaction(rp, hash)
    if err != nil {
        return nil, err
    }

    // Response
    response := api.JoinTNDAOJoinResponse{}

    // Join
    if opts, err := w.GetNodeAccountTransactor(); err != nil {
        return nil, err
    } else if hash, err := tndao.Join(rp, opts); err != nil {
        return nil, err
    } else {
        response.JoinTxHash = hash
    }

    // Return response
    return &response, nil

}

