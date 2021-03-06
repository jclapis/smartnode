package odao

import (
    "bytes"

    "github.com/rocket-pool/rocketpool-go/dao"
    "github.com/rocket-pool/rocketpool-go/dao/trustednode"
    rptypes "github.com/rocket-pool/rocketpool-go/types"
    "github.com/urfave/cli"
    "golang.org/x/sync/errgroup"

    "github.com/rocket-pool/smartnode/shared/services"
    "github.com/rocket-pool/smartnode/shared/types/api"
)


func canCancelProposal(c *cli.Context, proposalId uint64) (*api.CanCancelTNDAOProposalResponse, error) {

    // Get services
    if err := services.RequireNodeTrusted(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.CanCancelTNDAOProposalResponse{}

    // Sync
    var wg errgroup.Group

    // Check proposal exists
    wg.Go(func() error {
        proposalCount, err := dao.GetProposalCount(rp, nil)
        if err == nil {
            response.DoesNotExist = (proposalId > proposalCount)
        }
        return err
    })

    // Check proposal state
    wg.Go(func() error {
        proposalState, err := dao.GetProposalState(rp, proposalId, nil)
        if err == nil {
            response.InvalidState = !(proposalState == rptypes.Pending || proposalState == rptypes.Active)
        }
        return err
    })

    // Check proposer address
    wg.Go(func() error {
        nodeAccount, err := w.GetNodeAccount()
        if err != nil {
            return err
        }
        proposerAddress, err := dao.GetProposalProposerAddress(rp, proposalId, nil)
        if err == nil {
            response.InvalidProposer = !bytes.Equal(proposerAddress.Bytes(), nodeAccount.Address.Bytes())
        }
        return err
    })

    // Wait for data
    if err := wg.Wait(); err != nil {
        return nil, err
    }

    // Update & return response
    response.CanCancel = !(response.DoesNotExist || response.InvalidState || response.InvalidProposer)
    return &response, nil

}


func cancelProposal(c *cli.Context, proposalId uint64) (*api.CancelTNDAOProposalResponse, error) {

    // Get services
    if err := services.RequireNodeTrusted(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.CancelTNDAOProposalResponse{}

    // Get transactor
    opts, err := w.GetNodeAccountTransactor()
    if err != nil {
        return nil, err
    }

    // Cancel proposal
    txReceipt, err := trustednode.CancelProposal(rp, proposalId, opts)
    if err != nil {
        return nil, err
    }
    response.TxHash = txReceipt.TxHash

    // Return response
    return &response, nil

}

