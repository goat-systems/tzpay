package tzkt

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

/*
Delegators -
See: https://api.tzkt.io/#operation/Rewards_GetRewardSplit
ExtraFields:
	- NetRewards     int     `json:"net_rewards"`
	- GrossRewards   int     `json:"gross_rewards"`
	- Share          float64 `json:"share"`
	- Fee            int     `json:"fee"`
	- DexterContract string  `json:"dexter_contra
*/
type Delegators []Delegator

/*
Delegator -
See: https://api.tzkt.io/#operation/Rewards_GetRewardSplit
ExtraFields:
	- NetRewards     int     `json:"net_rewards"`
	- GrossRewards   int     `json:"gross_rewards"`
	- Share          float64 `json:"share"`
	- Fee            int     `json:"fee"`
	- DexterContract string  `json:"dexter_contra
*/
type Delegator struct {
	Address            string              `json:"address"`
	Balance            int                 `json:"balance"`
	CurrentBalance     int                 `json:"currentBalance"`
	Emptied            bool                `json:"emptied"`
	NetRewards         int                 `json:"net_rewards"`
	GrossRewards       int                 `json:"gross_rewards"`
	Share              float64             `json:"share"`
	Fee                int                 `json:"fee"`
	LiquidityProviders []LiquidityProvider `json:"liquidity_providers,omitempty"`
	BlackListed        bool                `json:"blacklisted,omitempty"`
}

/*
LiquidityProvider -
This is an extra structure to be embedded in Delegators for the purpose of paying
out Dexter Exchange Contracts.
*/
type LiquidityProvider struct {
	Address      string  `json:"address"`
	Balance      int     `json:"balance"`
	NetRewards   int     `json:"net_rewards"`
	GrossRewards int     `json:"gross_rewards"`
	Share        float64 `json:"share"`
	Fee          int     `json:"fee"`
	BlackListed  bool    `json:"blacklisted"`
}

/*
RewardsSplit -
See: https://api.tzkt.io/#operation/Rewards_GetRewardSplit
*/
type RewardsSplit struct {
	Cycle                       int        `json:"cycle"`
	StakingBalance              int        `json:"stakingBalance"`
	DelegatedBalance            int        `json:"delegatedBalance"`
	NumDelegators               int        `json:"numDelegators"`
	ExpectedBlocks              float64    `json:"expectedBlocks"`
	ExpectedEndorsements        float64    `json:"expectedEndorsements"`
	FutureBlocks                int        `json:"futureBlocks"`
	FutureBlockRewards          int        `json:"futureBlockRewards"`
	FutureBlockDeposits         int        `json:"futureBlockDeposits"`
	OwnBlocks                   int        `json:"ownBlocks"`
	OwnBlockRewards             int        `json:"ownBlockRewards"`
	ExtraBlocks                 int        `json:"extraBlocks"`
	ExtraBlockRewards           int        `json:"extraBlockRewards"`
	MissedOwnBlocks             int        `json:"missedOwnBlocks"`
	MissedOwnBlockRewards       int        `json:"missedOwnBlockRewards"`
	MissedExtraBlocks           int        `json:"missedExtraBlocks"`
	MissedExtraBlockRewards     int        `json:"missedExtraBlockRewards"`
	UncoveredOwnBlocks          int        `json:"uncoveredOwnBlocks"`
	UncoveredOwnBlockRewards    int        `json:"uncoveredOwnBlockRewards"`
	UncoveredExtraBlocks        int        `json:"uncoveredExtraBlocks"`
	UncoveredExtraBlockRewards  int        `json:"uncoveredExtraBlockRewards"`
	BlockDeposits               int        `json:"blockDeposits"`
	FutureEndorsements          int        `json:"futureEndorsements"`
	FutureEndorsementRewards    int        `json:"futureEndorsementRewards"`
	FutureEndorsementDeposits   int        `json:"futureEndorsementDeposits"`
	Endorsements                int        `json:"endorsements"`
	EndorsementRewards          int        `json:"endorsementRewards"`
	MissedEndorsements          int        `json:"missedEndorsements"`
	MissedEndorsementRewards    int        `json:"missedEndorsementRewards"`
	UncoveredEndorsements       int        `json:"uncoveredEndorsements"`
	UncoveredEndorsementRewards int        `json:"uncoveredEndorsementRewards"`
	EndorsementDeposits         int        `json:"endorsementDeposits"`
	OwnBlockFees                int        `json:"ownBlockFees"`
	ExtraBlockFees              int        `json:"extraBlockFees"`
	MissedOwnBlockFees          int        `json:"missedOwnBlockFees"`
	MissedExtraBlockFees        int        `json:"missedExtraBlockFees"`
	UncoveredOwnBlockFees       int        `json:"uncoveredOwnBlockFees"`
	UncoveredExtraBlockFees     int        `json:"uncoveredExtraBlockFees"`
	DoubleBakingRewards         int        `json:"doubleBakingRewards"`
	DoubleBakingLostDeposits    int        `json:"doubleBakingLostDeposits"`
	DoubleBakingLostRewards     int        `json:"doubleBakingLostRewards"`
	DoubleBakingLostFees        int        `json:"doubleBakingLostFees"`
	DoubleEndorsingRewards      int        `json:"doubleEndorsingRewards"`
	DoubleEndorsingLostDeposits int        `json:"doubleEndorsingLostDeposits"`
	DoubleEndorsingLostRewards  int        `json:"doubleEndorsingLostRewards"`
	DoubleEndorsingLostFees     int        `json:"doubleEndorsingLostFees"`
	RevelationRewards           int        `json:"revelationRewards"`
	RevelationLostRewards       int        `json:"revelationLostRewards"`
	RevelationLostFees          int        `json:"revelationLostFees"`
	Delegators                  Delegators `json:"delegators"`
	OperationLink               []string   `json:"operation_links,omitempty"`
	BakerRewards                int        `json:"baker_rewards,omitempty"`
	BakerShare                  float64    `json:"baker_share,omitempty"`
	BakerCollectedFees          int        `json:"collected_fees,omitempty"`
}

/*
GetRewardsSplit -
See: https://api.tzkt.io/#operation/Rewards_GetRewardSplit
*/
func (t *Tzkt) GetRewardsSplit(delegate string, cycle int, options ...URLParameters) (RewardsSplit, error) {
	resp, err := t.get(fmt.Sprintf("/v1/rewards/split/%s/%d", delegate, cycle), options...)
	if err != nil {
		return RewardsSplit{}, errors.Wrapf(err, "failed to get reward split")
	}

	var rewardsSplit RewardsSplit
	if err := json.Unmarshal(resp, &rewardsSplit); err != nil {
		return RewardsSplit{}, errors.Wrap(err, "failed to get reward split")
	}

	return rewardsSplit, nil
}
