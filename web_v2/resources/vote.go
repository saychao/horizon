package resources

import (
	"github.com/saychao/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
	"strconv"
)

func NewVoteKey(voteID int64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatInt(voteID, 10),
		Type: regources.VOTES,
	}
}

func NewVote(record history2.Vote) regources.Vote {
	return regources.Vote{
		Key: NewVoteKey(record.ID),
		Attributes: regources.VoteAttributes{
			VoteData: regources.VoteData(record.VoteData),
		},
		Relationships: regources.VoteRelationships{
			Voter: NewAccountKey(record.VoterID).AsRelation(),
			Poll:  NewPollKey(record.PollID).AsRelation(),
		},
	}
}
