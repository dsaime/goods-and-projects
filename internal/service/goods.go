package service

type Goods struct{}

type GoodsIn struct{}
type GoodsOut struct{}

func (g *Goods) Goods(in GoodsIn) (GoodsOut, error) {}

type CreateGoodIn struct{}
type CreateGoodOut struct{}

func (g *Goods) CreateGood(in CreateGoodIn) (CreateGoodOut, error) {}

type UpdateGoodIn struct {
	ID        int
	ProjectID int
}

type UpdateGoodOut struct{}

func (g *Goods) UpdateGood(in UpdateGoodIn) (UpdateGoodOut, error) {}

type DeleteGoodIn struct{}
type DeleteGoodOut struct{}

func (g *Goods) DeleteGood(in DeleteGoodIn) (DeleteGoodOut, error) {}

type ReprioritiizeGoodIn struct{}
type ReprioritiizeGoodOut struct{}

func (g *Goods) ReprioritiizeGood(in ReprioritiizeGoodIn) (ReprioritiizeGoodOut, error) {}
