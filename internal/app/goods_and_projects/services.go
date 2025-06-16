package app

import "github.com/dsaime/goods-and-projects/internal/service"

type services struct {
	goods *service.Goods
}

func (s *services) Goods() *service.Goods {
	return s.goods
}

func initServices(repos *repositories, adapterss *adapters) *services {
	return &services{
		goods: &service.Goods{
			Repo:             repos.goods,
			GoodsEventLogger: adapterss.goodsEventLogger,
		},
	}
}
