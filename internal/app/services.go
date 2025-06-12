package app

import "github.com/dsaime/goods-and-projects/internal/service"

type services struct {
	goods *service.Goods
}

func (s *services) Goods() *service.Goods {
	return s.goods
}

func initServices(repos *repositories, adaps *adapters) *services {
	return &services{
		goods: &service.Goods{},
	}
}
