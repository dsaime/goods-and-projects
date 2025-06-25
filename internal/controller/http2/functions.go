package http2

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
)

// DecodeBody декодирует тело запроса в структуру dst.
func DecodeBody(context Context, dst any) error {
	if err := json.NewDecoder(context.Request().Body).Decode(dst); err != nil {
		return fmt.Errorf("json Decoder Decode: %w", err)
	}

	return nil
}

// FormStr возвращает значение query-параметра или form-data запроса
func FormStr(context Context, name string) string {
	return context.Request().FormValue(name)
}

// FormInt возвращает значение query-параметра или form-data запроса
func FormInt(context Context, name string) int {
	valStr := FormStr(context, name)
	val, _ := strconv.Atoi(valStr)
	return val
}

// WrapHandlerWithMiddlewares оборачивает обработчик h всеми переданными middleware.
// Middlewares применяются в обратном порядке — от последнего к первому,
// таким образом формируя цепочку, где первый middleware выполняется первым.
func WrapHandlerWithMiddlewares(h HandlerFunc, middlewares ...Middleware) HandlerFuncRW {
	hrw := AdaptToRW(h)
	for _, mw := range slices.Backward(middlewares) {
		hrw = mw(hrw)
	}

	return hrw
}

// AdaptToRW оборачивает HandlerFunc, преобразуя его в HandlerFuncRW.
// Позволяет использовать обработчики, принимающие простой Context,
// в системе, где требуется RWContext.
func AdaptToRW(handlerFunc HandlerFunc) HandlerFuncRW {
	return func(context RWContext) (any, error) {
		return handlerFunc(context)
	}
}
