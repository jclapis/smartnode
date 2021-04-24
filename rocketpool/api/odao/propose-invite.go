package odao

import (
    "fmt"

    "github.com/ethereum/go-ethereum/common"
    "github.com/rocket-pool/rocketpool-go/dao/trustednode"
    "github.com/urfave/cli"
    "golang.org/x/sync/errgroup"

    "github.com/rocket-pool/smartnode/shared/services"
    "github.com/rocket-pool/smartnode/shared/types/api"
)


func canProposeInvite(c *cli.Context, memberAddress common.Address) (*api.CanProposeTNDAOInviteResponse, error) {

    // Get services
    if err := services.RequireNodeTrusted(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.CanProposeTNDAOInviteResponse{}

    // Sync
    var wg errgroup.Group

    // Check if proposal cooldown is active
    wg.Go(func() error {
        nodeAccount, err := w.GetNodeAccount()
        if err != nil {
            return err
        }
        proposalCooldownActive, err := getProposalCooldownActive(rp, nodeAccount.Address)
        if err == nil {
            response.ProposalCooldownActive = proposalCooldownActive
        }
        return err
    })

    // Check if member exists
    wg.Go(func() error {
        memberExists, err := trustednode.GetMemberExists(rp, memberAddress, nil)
        if err == nil {
            response.MemberAlreadyExists = memberExists
        }
        return err
    })

    // Wait for data
    if err := wg.Wait(); err != nil {
        return nil, err
    }

    // Update & return response
    response.CanPropose = !(response.ProposalCooldownActive || response.MemberAlreadyExists)
    return &response, nil

}


func proposeInvite(c *cli.Context, memberAddress common.Address, memberId, memberEmail string) (*api.ProposeTNDAOInviteResponse, error) {

    // Get services
    if err := services.RequireNodeTrusted(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.ProposeTNDAOInviteResponse{}

    // Get transactor
    opts, err := w.GetNodeAccountTransactor()
    if err != nil {
        return nil, err
    }

    // Submit proposal
    message := fmt.Sprintf("invite %s (%s)", memberId, memberEmail)
    proposalId, hash, err := trustednode.ProposeInviteMember(rp, message, memberAddress, memberId, memberEmail, opts)
    if err != nil {
        return nil, err
    }
    response.ProposalId = proposalId
    response.TxHash = hash

    // Return response
    return &response, nil

}

