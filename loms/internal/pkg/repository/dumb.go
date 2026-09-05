package repository

type DumbRepo struct {
}

func NewDumbRepo() *DumbRepo {
	return &DumbRepo{}
}

func (DumbRepo) GetStocks(n uint32) uint64 {
	return 1000
}
