package core2

type SaleIDParticipantsCount struct {
	SaleID           uint64 `db:"sale_id"`
	ParticipantCount int64  `db:"p_count"`
}

type SaleParticipantsMap map[uint64]int64

func SaleParticipantsToMap(sp []SaleIDParticipantsCount) SaleParticipantsMap {
	saleParticipants := make(SaleParticipantsMap)

	for idx := range sp {
		sale := sp[idx]
		saleParticipants[sale.SaleID] = sale.ParticipantCount
	}

	return saleParticipants
}
