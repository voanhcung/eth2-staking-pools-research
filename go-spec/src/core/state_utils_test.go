package core

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCommitteeShuffling(t *testing.T) {
	// test state
	pools := 5
	bpInPool := 128
	bps := make([]*BlockProducer, pools * bpInPool)
	for i := 0 ; i < len(bps) ; i++ {
		bps[i] = &BlockProducer{
			Id:      uint64(i),
			Balance: 1000,
			Stake:   0,
			Slashed: false,
			Active:  true,
			PubKey:  []byte(fmt.Sprintf("pubkey %d", i)),
		}
	}

	state := &State{
		GenesisTime:          0,
		CurrentEpoch:         0,
		BlockRoots:           nil,
		StateRoots:           nil,
		Seeds:                []*EpochAndBytes{
			&EpochAndBytes{
				Epoch:               0,
				Bytes:                 []byte("seedseedseedseedseedseedseedseed"),
			},
			&EpochAndBytes{
				Epoch:               1,
				Bytes:                 []byte("sdddseedseedseedseedseedseedseed"),
			},
		},
		BlockProducers:       bps,
		Pools:                nil,
		Slashings:            nil,
	}

	tests := []struct{
		name string
		epoch uint64
		poolId uint64
		dkgReqId uint64
		expectedPoolCommittee []uint64
		expectedBlockVotingCommittee []uint64
		expectedBlockProposer uint64
	}{
		{
			name:"epoch 1: pool id:1, dkg req:1",
			epoch: 1,
			poolId: 1,
			dkgReqId:1,
			expectedPoolCommittee: []uint64{493,19,340,403,27,250,419,446,604,209,84,6,343,64,16,137,626,451,282,322,296,369,576,28,427,452,165,287,552,600,330,308,61,333,305,632,156,388,405,131,365,187,524,132,180,107,161,18,444,97,585,29,326,538,342,3,364,579,197,411,40,30,431,533,543,357,561,303,465,138,386,90,389,401,261,391,352,116,285,501,233,521,110,316,253,371,328,46,569,73,186,93,321,309,41,98,613,624,195,550,519,539,360,67,141,382,44,106,473,367,377,109,152,276,616,266,467,598,103,537,361,164,297,481,573,515,460,314},
			expectedBlockVotingCommittee: []uint64{447,592,293,502,129,115,279,311,480,40,249,1,119,555,603,78,281,629,174,201,375,256,237,611,324,471,381,152,236,65,29,517,285,465,58,341,467,32,329,260,80,335,312,252,489,560,251,38,186,390,10,319,537,43,294,527,214,188,216,127,615,107,27,617,30,83,155,463,230,63,232,549,28,444,448,185,94,562,34,624,382,55,153,265,240,97,581,261,172,533,374,543,4,524,409,454,308,72,60,183,371,231,432,267,618,495,564,349,469,566,378,376,635,101,276,587,318,623,242,21,104,112,637,321,208,103,500,275},
			expectedBlockProposer: 17,
		},
		{
			name:"epoch 1: pool id:2, dkg req:2",
			epoch: 1,
			poolId: 2,
			dkgReqId:2,
			expectedPoolCommittee: []uint64{518,583,464,375,196,145,323,582,527,164,141,627,449,429,579,500,30,387,103,29,617,434,538,74,140,474,497,501,370,229,204,535,365,174,421,213,180,599,114,453,600,547,42,408,102,18,555,20,526,393,433,302,151,498,432,461,610,513,413,591,631,85,269,633,377,459,263,562,548,84,296,355,248,289,268,603,557,313,208,485,334,622,625,173,221,493,153,348,266,374,616,537,250,572,139,237,382,604,86,505,439,492,340,327,575,396,195,95,70,343,33,632,9,184,88,240,570,371,187,127,255,511,69,611,419,178,568,32},
			expectedBlockVotingCommittee: []uint64{447,592,293,502,129,115,279,311,480,40,249,1,119,555,603,78,281,629,174,201,375,256,237,611,324,471,381,152,236,65,29,517,285,465,58,341,467,32,329,260,80,335,312,252,489,560,251,38,186,390,10,319,537,43,294,527,214,188,216,127,615,107,27,617,30,83,155,463,230,63,232,549,28,444,448,185,94,562,34,624,382,55,153,265,240,97,581,261,172,533,374,543,4,524,409,454,308,72,60,183,371,231,432,267,618,495,564,349,469,566,378,376,635,101,276,587,318,623,242,21,104,112,637,321,208,103,500,275},
			expectedBlockProposer: 17,
		},
		{
			name:"epoch 2: pool id:1, dkg req:1",
			epoch: 2,
			poolId: 2,
			dkgReqId:2,
			expectedPoolCommittee: []uint64{116,301,182,224,376,506,435,421,548,540,447,226,487,458,232,496,463,328,373,375,178,210,56,205,344,541,425,525,117,177,231,609,189,619,615,204,300,108,235,513,95,174,213,225,345,192,436,563,483,223,532,206,516,271,330,493,591,23,195,618,130,341,3,608,181,634,456,146,381,546,96,267,184,290,494,251,287,512,269,566,292,560,12,326,323,418,211,134,26,617,20,244,613,61,538,197,630,159,281,148,460,1,123,557,509,365,259,398,520,171,389,104,558,209,58,426,246,109,607,550,80,378,359,317,585,399,445,289},
			expectedBlockVotingCommittee: []uint64{28,345,32,278,304,498,322,452,423,147,399,46,21,355,38,404,403,118,270,273,248,141,354,466,467,154,618,596,33,96,378,318,90,132,420,497,14,128,64,617,266,442,15,372,319,334,0,216,461,477,357,212,516,68,183,347,284,185,388,472,305,531,592,542,384,259,622,173,144,570,39,462,65,138,448,441,362,81,367,499,136,434,134,582,507,454,301,568,161,36,481,558,320,478,310,534,342,235,336,131,427,620,291,364,623,111,513,12,369,18,522,389,208,366,376,533,398,521,249,75,79,267,346,151,120,240,566,588},
			expectedBlockProposer: 20,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pc,err := PoolCommittee(state, test.poolId, test.epoch)
			require.NoError(t, err)
			require.EqualValues(t, test.expectedPoolCommittee, pc)

			voting,err := BlockVotingCommittee(state, test.epoch)
			require.NoError(t, err)
			require.EqualValues(t, test.expectedBlockVotingCommittee, voting)

			proposer,err := GetBlockProposer(state, test.epoch)
			require.NoError(t, err)
			require.EqualValues(t, test.expectedBlockProposer, proposer)
		})
	}
}
