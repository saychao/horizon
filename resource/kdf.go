package resource

import (
	"math"
)

type KdfParams struct {
	Algorithm string  `json:"algorithm"`
	Bits      uint    `json:"bits"`
	N         float64 `json:"n"`
	R         uint    `json:"r"`
	P         uint    `json:"p"`
}

func (p *KdfParams) Populate() {
	p.Algorithm = "scrypt"
	p.Bits = 256
	p.N = math.Pow(2, 12)
	p.R = 8
	p.P = 1
}
