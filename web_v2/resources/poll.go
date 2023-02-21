package resources

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

func NewPollKey(ID int64) regources.Key {
	return regources.NewKeyInt64(ID, regources.POLLS)
}

func NewParticipationKey(ID int64) regources.Key {
	return regources.NewKeyInt64(ID, regources.POLL_OUTCOME)
}

func NewPoll(record history2.Poll) regources.Poll {
	return regources.Poll{
		Key: NewPollKey(record.ID),
		Attributes: regources.PollAttributes{
			VoteConfirmationRequired: record.VoteConfirmationRequired,
			PermissionType:           record.PermissionType,
			NumberOfChoices:          record.NumberOfChoices,
			StartTime:                record.StartTime,
			EndTime:                  record.EndTime,
			CreatorDetails:           record.Details,
			PollData:                 regources.PollData(record.Data),
			PollState:                record.State,
		},
		Relationships: regources.PollRelationships{
			ResultProvider: NewAccountKey(record.ResultProviderID).AsRelation(),
			Owner:          NewAccountKey(record.OwnerID).AsRelation(),
			Participation:  NewParticipationKey(record.ID).AsRelation(),
		},
	}
}
func NewParticipation(id int64, historyVotes []history2.Vote) regources.PollParticipation {
	outcome := regources.PollParticipation{
		Key: NewParticipationKey(id),
		Relationships: regources.PollParticipationRelationships{
			Votes: &regources.RelationCollection{},
		},
	}
	for _, v := range historyVotes {
		vote := NewVoteKey(v.ID)
		outcome.Relationships.Votes.Data = append(outcome.Relationships.Votes.Data, vote)
	}

	return outcome
}
