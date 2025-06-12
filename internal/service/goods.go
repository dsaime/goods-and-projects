package service

type Goods struct{}

type GoodsIn struct {
	Limit  int
	Offset int
}
type GoodsOut struct{}

func (g *Goods) Goods(in GoodsIn) (GoodsOut, error) {
	panic("implement me")
}

type CreateGoodIn struct {
	ID        int
	ProjectID int
	Name      string
}
type CreateGoodOut struct{}

func (g *Goods) CreateGood(in CreateGoodIn) (CreateGoodOut, error) {
	panic("implement me")
}

type UpdateGoodIn struct {
	ID          int
	ProjectID   int
	Name        string
	Description string
}

type UpdateGoodOut struct{}

func (g *Goods) UpdateGood(in UpdateGoodIn) (UpdateGoodOut, error) {
	panic("implement me")
}

type DeleteGoodIn struct {
	ID        int
	ProjectID int
}
type DeleteGoodOut struct{}

func (g *Goods) DeleteGood(in DeleteGoodIn) (DeleteGoodOut, error) {
	panic("implement me")
}

type ReprioritiizeGoodIn struct {
	ID          int
	ProjectID   int
	NewPriority string
}
type ReprioritiizeGoodOut struct{}

func (g *Goods) ReprioritiizeGood(in ReprioritiizeGoodIn) (ReprioritiizeGoodOut, error) {
	panic("implement me")
}
