package router

import (
	"log/slog"
	"net/http"

	"github.com/dsaime/goods-and-projects/internal/controller/http2"
)

// Router обрабатывает HTTP-запросы
type Router struct {
	Services http2.Services
	http.ServeMux
}

// HandleFunc регистрирует обработчик для переданного паттерна и middleware над ним.
func (c *Router) HandleFunc(pattern string, chain []http2.Middleware, handlerFunc http2.HandlerFunc) {
	handlerFuncRW := http2.WrapHandlerWithMiddlewares(handlerFunc, chain...)
	httpHandlerFunc := c.modulation(handlerFuncRW)
	c.ServeMux.HandleFunc(pattern, httpHandlerFunc)
	slog.Info("Router.HandleFunc: Зарегистрирован новый обработчик",
		slog.String("pattern", pattern),
	)
}
