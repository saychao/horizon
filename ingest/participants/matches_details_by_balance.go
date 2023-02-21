package participants

import (
	"fmt"

	"gitlab.com/tokend/go/xdr"
)

type MatchesDetailsByBalance struct {
	Details map[string]*MatchesDetails
}

func (m *MatchesDetailsByBalance) ToParticipants(participants []Participant) []Participant {
	// temporary map to find existing participant for details merging
	byBalance := map[string]Participant{}
	for _, participant := range participants {
		byBalance[participant.BalanceID.AsString()] = participant
	}

	for _, detail := range m.Details {
		pending := detail.ToParticipant()
		if existing, ok := byBalance[pending.BalanceID.AsString()]; ok {
			// merge details in case same balance participated twice
			pending = m.mergeParticipantDetails(existing, pending)
		}
		byBalance[pending.BalanceID.AsString()] = pending
	}

	// converting back to slice
	participants = make([]Participant, 0, len(byBalance))
	for _, participant := range byBalance {
		participants = append(participants, participant)
	}

	return participants
}

func (m *MatchesDetailsByBalance) mergeParticipantDetails(a Participant, b Participant) Participant {
	aDetails, ok := a.Details.([]MatchIngestDetails)
	if !ok {
		panic(fmt.Sprintf("expected slice of details got %T", a.Details))
	}
	bDetails, ok := b.Details.([]MatchIngestDetails)
	if !ok {
		panic(fmt.Sprintf("expected slice of details got %T", b.Details))
	}
	a.Details = append(aDetails, bDetails...)
	return Participant{
		a.AccountID,
		a.BalanceID,
		append(aDetails, bDetails...),
	}
}

func NewMatchesDetailsByBalance() *MatchesDetailsByBalance {
	return &MatchesDetailsByBalance{
		Details: make(map[string]*MatchesDetails),
	}
}

func (m *MatchesDetailsByBalance) add(source xdr.AccountId, balance xdr.BalanceId, baseAsset, quoteAsset xdr.AssetCode, isBuy bool, match *Match) {
	matches, ok := m.Details[balance.AsString()]
	if !ok {
		matches = NewMatchesDetails(source, balance, string(baseAsset), string(quoteAsset), isBuy)
		m.Details[balance.AsString()] = matches
	}

	matches.Add(match)
}

func (m *MatchesDetailsByBalance) Add(source xdr.AccountId, balance []xdr.BalanceId, baseAsset, quoteAsset xdr.AssetCode, isBuy bool, match *Match) {
	for i := range balance {
		m.add(source, balance[i], baseAsset, quoteAsset, isBuy, match)
	}
}
