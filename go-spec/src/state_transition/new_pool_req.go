package state_transition

import (
	"fmt"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/core"
	"github.com/bloxapp/eth2-staking-pools-research/go-spec/src/shared"
	"sort"
)

func (st *StateTransition) ProcessNewPoolRequests(state *core.State, requests []*core.CreateNewPoolRequest) error {
	for _, req := range requests {
		leader := core.GetBlockProducer(state, req.StartEpoch)
		if leader == nil {
			return fmt.Errorf("could not find new pool req leader")
		}

		// verify leader is correct
		if req.LeaderBlockProducer != leader.Id {
			return fmt.Errorf("new pool req leader incorrect")
		}
		if core.GetPool(state, req.Id) != nil {
			return fmt.Errorf("new pool id == req id, this is already exists")
		}
		// TODO - check that network has enough capitalization
		// TODO - check leader is not part of DKG Committee

		// get DKG participants
		committee, err := core.DKGCommittee(state, req.Id, req.StartEpoch)
		if err != nil {
			return err
		}
		sort.Slice(committee, func(i int, j int) bool {
			return committee[i] < committee[j]
		})

		switch req.GetStatus() {
		case 0:
			// TODO if i'm the DKDG leader act uppon it
		case 1: // successful
			// get committee
			committee, err := core.DKGCommittee(state, req.Id, req.StartEpoch)
			sort.Slice(committee, func(i int, j int) bool {
				return committee[i] < committee[j]
			})

			state.Pools = append(state.Pools, &core.Pool{
				Id:              req.Id,
				PubKey:          req.GetCreatePubKey(),
				SortedCommittee: committee,
			})
			if err != nil {
				return err
			}

			// reward/ penalty
			for i := 0 ; i < len(committee) ; i ++ {
				bp := core.GetBlockProducer(state, committee[i])
				if bp == nil {
					return fmt.Errorf("could not find BP %d", committee[i])
				}
				partic := req.GetParticipation()
				if shared.IsBitSet(partic[:], uint64(i)) {
					err := core.IncreaseBPBalance(bp, core.TestConfig().DKGReward)
					if err != nil {
						return err
					}
				} else {
					err := core.DecreaseBPBalance(bp, core.TestConfig().DKGReward)
					if err != nil {
						return err
					}
				}
			}

			// special reward for leader
			err = core.IncreaseBPBalance(leader, 3* core.TestConfig().DKGReward)
			if err != nil {
				return err
			}
		case 2: // un-successful
			for i := 0 ; i < len(committee) ; i ++ {
				bp := core.GetBlockProducer(state, committee[i])
				if bp == nil {
					return fmt.Errorf("could not find BP %d", committee[i])
				}
				err := core.DecreaseBPBalance(bp, core.TestConfig().DKGReward)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}