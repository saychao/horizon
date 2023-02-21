package storage

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// ParticipantEffect - helper struct to store `operation participants`
type ParticipantEffect struct {
	repo *pgdb.DB
}

// NewOpParticipants - creates new instance of `ParticipantEffect`
func NewOpParticipants(repo *pgdb.DB) *ParticipantEffect {
	return &ParticipantEffect{
		repo: repo,
	}
}

func convertParticipantEffectToParams(participant history2.ParticipantEffect) []interface{} {
	return []interface{}{
		participant.ID, participant.AccountID, participant.BalanceID, participant.AssetCode,
		participant.Effect, participant.OperationID,
	}
}

// Insert - stores slice of the participant effects in one batch
func (p *ParticipantEffect) Insert(participants []history2.ParticipantEffect) error {
	columns := []string{"id, account_id, balance_id, asset_code, effect, operation_id"}
	err := history2ParticipantEffectBatchInsert(p.repo, participants, "participant_effects", columns, convertParticipantEffectToParams)
	if err != nil {
		return errors.Wrap(err, "failed to insert participant effects")
	}
	return nil
}
