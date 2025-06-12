package register_handler

import (
	"github.com/dsaime/goods-and-projects/internal/controller/http2"
	"github.com/dsaime/goods-and-projects/internal/controller/http2/middleware"
	"github.com/dsaime/goods-and-projects/internal/service"
)

func GoodsList(router http2.Router) {
	router.HandleFunc(
		"GET /goods/list",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {
			input := service.GoodsIn{
				Limit:  http2.FormInt(context, "limit"),
				Offset: http2.FormInt(context, "offset"),
			}

			return context.Services().Goods().Goods(input)
		})
}

func GoodsCreate(router http2.Router) {
	type requestBody struct {
		Name string `json:"name"`
	}

	router.HandleFunc(
		"POST /goods/create",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			input := service.CreateGoodIn{
				ID:        http2.FormInt(context, "id"),
				ProjectID: http2.FormInt(context, "projectId"),
				Name:      rb.Name,
			}

			return context.Services().Goods().CreateGood(input)
		})
}

func GoodsUpdate(router http2.Router) {
	type requestBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	router.HandleFunc(
		"PATH /goods/update",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			input := service.UpdateGoodIn{
				ID:          http2.FormInt(context, "id"),
				ProjectID:   http2.FormInt(context, "projectId"),
				Name:        rb.Name,
				Description: rb.Description,
			}

			return context.Services().Goods().UpdateGood(input)
		})
}

func GoodsDelete(router http2.Router) {
	router.HandleFunc(
		"DELETE /goods/remove",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {

			input := service.DeleteGoodIn{
				ID:        http2.FormInt(context, "id"),
				ProjectID: http2.FormInt(context, "projectId"),
			}

			return context.Services().Goods().DeleteGood(input)
		})
}

func GoodsReprioritiize(router http2.Router) {
	type requestBody struct {
		NewPriority string `json:"newPriority"`
	}
	router.HandleFunc(
		"PATH /goods/reprioritiize",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			input := service.ReprioritiizeGoodIn{
				ID:          http2.FormInt(context, "id"),
				ProjectID:   http2.FormInt(context, "projectId"),
				NewPriority: rb.NewPriority,
			}

			return context.Services().Goods().ReprioritiizeGood(input)
		})
}
