package state

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/prysmaticlabs/go-ssz"
)

var helperFunc src.NonSpecFunctions

type State struct {
	pools          []*Pool
	currentEpoch   uint64
	blockProducers []*BlockProducer
	seed           [32]byte
}

func (state *State) Root() ([32]byte,error) {
	return ssz.HashTreeRoot(state)
}

func (state *State) GetPools() []core.IPool {

}

func (state *State) GetPool(id uint64) core.IPool {
	for _, p := range state.pools {
		if p.id == id {
			return p
		}
	}
	return nil
}

func (state *State) AddNewPool(pool core.IPool) error {
	if found := state.GetPool(pool.GetId()); found != nil {
		return fmt.Errorf("pool already exists")
	}

	state.pools = append(state.pools, pool)
	return nil
}

func (state *State) GetBlockProducers() []core.IBlockProducer {
	return state.blockProducers
}

func (state *State) GetBlockProducer(id uint64) core.IBlockProducer {
	for _, bp := range state.GetBlockProducers() {
		if bp.GetId() == id {
			return bp
		}
	}
	return nil
}

func (state *State) GetCurrentEpoch() uint64 {

}

func (state *State) GetSeed() [32]byte {
	return state.seed
}

func (state *State) SetSeed(seed [32]byte) {
	state.seed = seed
}

func (state *State) GetPastSeed(epoch uint64) [32]byte {

}


func (state *State) Copy() (core.IState, error) {
	copiedPools := make([]*Pool, len(state.pools))
	for i, p := range state.pools {
		newP, err := p.Copy()
		if err != nil {
			return nil, err
		}
		copiedPools[i] = newP
	}

	copiedBps := make([]*BlockProducer, len(state.blockProducers))
	for i, bp := range state.blockProducers {
		newBP, err := bp.Copy()
		if err != nil {
			return nil, err
		}
		copiedBps[i] = newBP
	}

	return &State{
		pools:          copiedPools,
		blockProducers: copiedBps,
		seed:           state.seed,
	}, nil
}
