package generator

import (
	"math"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

//ID - used to generate sequential ID in following format 64-bits[32-bits of significant; 32-bits of sequence]
type ID struct {
	Significant int32
	seq         uint32
}

//NewID - creates new instance of ID. If significant is > MaxInt32 - panics
func NewID(significant uint32) *ID {
	if significant > math.MaxInt32 {
		panic(errors.From(errors.New("trying to create new generator with significant > MaxInt32"), logan.F{
			"significant": significant,
		}))
	}
	return &ID{
		Significant: int32(significant),
		seq:         0,
	}
}

//NewIDI32 - creates new instance of ID. If significant is < 0 - panics
func NewIDI32(significant int32) *ID {
	if significant < 0 {
		panic(errors.New("negative significant is not supported"))
	}

	return NewID(uint32(significant))
}

//Next - generates new sequential ID
func (g *ID) Next() int64 {
	if g.seq == math.MaxUint32 {
		panic(errors.New("failed to generate next ID - overflow of sequence"))
	}

	g.seq++
	return (int64(g.Significant) << 32) | int64(g.seq)
}

func MakeIDUint32(significant int32, sequence uint32) int64 {
	if significant < 0 {
		panic(errors.New("negative significant is not supported"))
	}

	return (int64(significant) << 32) | int64(sequence)
}

func MakeIDUint16(significant int32, seq1, seq2 uint16) int64 {
	if significant < 0 {
		panic(errors.New("negative significant is not supported"))
	}

	return (int64(significant) << 32) | int64(seq1)<<16 | int64(seq2)
}

// Function gets sequence id from 64 bit integer
// Passed format should be following:
// 32 most significant bits are sequence id
// 32 least significant bits are some other staff
//
// this scheme is used in poll/vote IDs
func GetSeqFromInt64(id int64) int32 {
	return int32(id >> 32)
}
