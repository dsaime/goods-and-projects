package register_handler

import (
	"github.com/dsaime/goods-and-projects/internal/controller/http2"
	"github.com/dsaime/goods-and-projects/internal/controller/http2/middleware"
	"github.com/dsaime/goods-and-projects/internal/service"
)

// GoodsList регистрирует обработчик для получения списка существующих товаров.
//
// Метод: GET /goods/list
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

// GoodsCreate регистрирует обработчик для создания товаров.
//
// Метод: POST /goods/create
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
				ProjectID: http2.FormInt(context, "projectId"),
				Name:      rb.Name,
			}

			return context.Services().Goods().CreateGood(input)
		})
}

// GoodsUpdate регистрирует обработчик для обновления некоторых полей товаров.
//
// Метод: PATCH /goods/update
func GoodsUpdate(router http2.Router) {
	type requestBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	router.HandleFunc(
		"PATCH /goods/update",
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

// GoodsDelete регистрирует обработчик для удаления товаров.
//
// Метод: DELETE /goods/remove
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

// GoodsReprioritiize регистрирует обработчик для обновления приоритета товаров.
//
// Метод: PATCH /goods/reprioritiize
func GoodsReprioritiize(router http2.Router) {
	type requestBody struct {
		NewPriority int `json:"newPriority"`
	}
	router.HandleFunc(
		"PATCH /goods/reprioritiize",
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
