package state_transition

import (
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"github.com/prysmaticlabs/go-ssz"
)

func (st *StateTransition) ApplyBlock(oldState *core.State, body *core.BlockBody) (newState *core.State, err error) {
	newState = core.CopyState(oldState)

	// process
	if err := st.ProcessNewPoolRequests(newState, body.NewPoolReq); err != nil {
		return nil,err
	}

	// bump epoch
	newState.CurrentEpoch = body.Epoch
	// apply seed
	newSeed, err := shared.MixSeed(
		shared.SliceToByte32(oldState.Seeds[len(oldState.Seeds) - 1].Bytes), // previous seed
		shared.SliceToByte32(body.Randao[:32]))
	if err != nil {
		return nil, err
	}
	newState.Seeds = append(newState.Seeds, &core.EpochAndBytes{
		Epoch:                newState.CurrentEpoch,
		Bytes:                newSeed[:],
	})
	// add block root
	root, err := ssz.HashTreeRoot(body)
	if err != nil {
		return nil, err
	}
	newState.BlockRoots = append(newState.BlockRoots, &core.EpochAndBytes{
		Epoch:                newState.CurrentEpoch,
		Bytes:               root[:],
	})
	// state root
	root, err = ssz.HashTreeRoot(newState)
	if err != nil {
		return nil, err
	}
	newState.StateRoots = append(newState.StateRoots, &core.EpochAndBytes{
		Epoch:                newState.CurrentEpoch,
		Bytes:                root[:],
	})
	// ETH1 and ETH2 block/ epoch
	newState.ETH1Blocks = append(newState.ETH1Blocks, body.ETH1Block)
	newState.ETH2Epochs = append(newState.ETH2Epochs, body.ETH2Epoch)

	return newState, nil
}