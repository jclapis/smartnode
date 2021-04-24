package odao

import (
    "math/big"

    "github.com/rocket-pool/rocketpool-go/settings/trustednode"
    "github.com/urfave/cli"

    "github.com/rocket-pool/smartnode/shared/services"
    "github.com/rocket-pool/smartnode/shared/types/api"
)


func canProposeSetting(c *cli.Context) (*api.CanProposeTNDAOSettingResponse, error) {

    // Get services
    if err := services.RequireNodeTrusted(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.CanProposeTNDAOSettingResponse{}

    // Get node account
    nodeAccount, err := w.GetNodeAccount()
    if err != nil {
        return nil, err
    }

    // Check if proposal cooldown is active
    proposalCooldownActive, err := getProposalCooldownActive(rp, nodeAccount.Address)
    if err != nil {
        return nil, err
    }
    response.ProposalCooldownActive = proposalCooldownActive

    // Update & return response
    response.CanPropose = !response.ProposalCooldownActive
    return &response, nil

}


func proposeSettingMembersQuorum(c *cli.Context, quorum float64) (*api.ProposeTNDAOSettingMembersQuorumResponse, error) {

    // Get services
    if err := services.RequireNodeTrusted(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.ProposeTNDAOSettingMembersQuorumResponse{}

    // Get transactor
    opts, err := w.GetNodeAccountTransactor()
    if err != nil {
        return nil, err
    }

    // Submit proposal
    proposalId, hash, err := trustednode.ProposeQuorum(rp, quorum, opts)
    if err != nil {
        return nil, err
    }
    response.ProposalId = proposalId
    response.TxHash = hash

    // Return response
    return &response, nil

}


func proposeSettingMembersRplBond(c *cli.Context, bondAmountWei *big.Int) (*api.ProposeTNDAOSettingMembersRplBondResponse, error) {

    // Get services
    if err := services.RequireNodeTrusted(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.ProposeTNDAOSettingMembersRplBondResponse{}

    // Get transactor
    opts, err := w.GetNodeAccountTransactor()
    if err != nil {
        return nil, err
    }

    // Submit proposal
    proposalId, hash, err := trustednode.ProposeRPLBond(rp, bondAmountWei, opts)
    if err != nil {
        return nil, err
    }
    response.ProposalId = proposalId
    response.TxHash = hash

    // Return response
    return &response, nil

}


func proposeSettingMinipoolUnbondedMax(c *cli.Context, unbondedMinipoolMax uint64) (*api.ProposeTNDAOSettingMinipoolUnbondedMaxResponse, error) {

    // Get services
    if err := services.RequireNodeTrusted(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.ProposeTNDAOSettingMinipoolUnbondedMaxResponse{}

    // Get transactor
    opts, err := w.GetNodeAccountTransactor()
    if err != nil {
        return nil, err
    }

    // Submit proposal
    proposalId, hash, err := trustednode.ProposeMinipoolUnbondedMax(rp, unbondedMinipoolMax, opts)
    if err != nil {
        return nil, err
    }
    response.ProposalId = proposalId
    response.TxHash = hash

    // Return response
    return &response, nil

}


func proposeSettingProposalCooldown(c *cli.Context, proposalCooldownBlocks uint64) (*api.ProposeTNDAOSettingProposalCooldownResponse, error) {

    // Get services
    if err := services.RequireNodeTrusted(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.ProposeTNDAOSettingProposalCooldownResponse{}

    // Get transactor
    opts, err := w.GetNodeAccountTransactor()
    if err != nil {
        return nil, err
    }

    // Submit proposal
    proposalId, hash, err := trustednode.ProposeProposalCooldown(rp, proposalCooldownBlocks, opts)
    if err != nil {
        return nil, err
    }
    response.ProposalId = proposalId
    response.TxHash = hash

    // Return response
    return &response, nil

}


func proposeSettingProposalVoteBlocks(c *cli.Context, proposalVoteBlocks uint64) (*api.ProposeTNDAOSettingProposalVoteBlocksResponse, error) {

    // Get services
    if err := services.RequireNodeTrusted(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.ProposeTNDAOSettingProposalVoteBlocksResponse{}

    // Get transactor
    opts, err := w.GetNodeAccountTransactor()
    if err != nil {
        return nil, err
    }

    // Submit proposal
    proposalId, hash, err := trustednode.ProposeProposalVoteBlocks(rp, proposalVoteBlocks, opts)
    if err != nil {
        return nil, err
    }
    response.ProposalId = proposalId
    response.TxHash = hash

    // Return response
    return &response, nil

}


func proposeSettingProposalVoteDelayBlocks(c *cli.Context, proposalDelayBlocks uint64) (*api.ProposeTNDAOSettingProposalVoteDelayBlocksResponse, error) {

    // Get services
    if err := services.RequireNodeTrusted(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.ProposeTNDAOSettingProposalVoteDelayBlocksResponse{}

    // Get transactor
    opts, err := w.GetNodeAccountTransactor()
    if err != nil {
        return nil, err
    }

    // Submit proposal
    proposalId, hash, err := trustednode.ProposeProposalVoteDelayBlocks(rp, proposalDelayBlocks, opts)
    if err != nil {
        return nil, err
    }
    response.ProposalId = proposalId
    response.TxHash = hash

    // Return response
    return &response, nil

}


func proposeSettingProposalExecuteBlocks(c *cli.Context, proposalExecuteBlocks uint64) (*api.ProposeTNDAOSettingProposalExecuteBlocksResponse, error) {

    // Get services
    if err := services.RequireNodeTrusted(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.ProposeTNDAOSettingProposalExecuteBlocksResponse{}

    // Get transactor
    opts, err := w.GetNodeAccountTransactor()
    if err != nil {
        return nil, err
    }

    // Submit proposal
    proposalId, hash, err := trustednode.ProposeProposalExecuteBlocks(rp, proposalExecuteBlocks, opts)
    if err != nil {
        return nil, err
    }
    response.ProposalId = proposalId
    response.TxHash = hash

    // Return response
    return &response, nil

}


func proposeSettingProposalActionBlocks(c *cli.Context, proposalActionBlocks uint64) (*api.ProposeTNDAOSettingProposalActionBlocksResponse, error) {

    // Get services
    if err := services.RequireNodeTrusted(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }

    // Response
    response := api.ProposeTNDAOSettingProposalActionBlocksResponse{}

    // Get transactor
    opts, err := w.GetNodeAccountTransactor()
    if err != nil {
        return nil, err
    }

    // Submit proposal
    proposalId, hash, err := trustednode.ProposeProposalActionBlocks(rp, proposalActionBlocks, opts)
    if err != nil {
        return nil, err
    }
    response.ProposalId = proposalId
    response.TxHash = hash

    // Return response
    return &response, nil

}

