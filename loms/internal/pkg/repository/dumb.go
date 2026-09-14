package repository

import "math/rand"

type DumbRepo struct {
}

func NewDumbRepo() *DumbRepo {
	return &DumbRepo{}
}

func (DumbRepo) GetStocks(n uint32) uint64 {
	return rand.Uint64() % 1001
}
