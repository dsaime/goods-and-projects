package router

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/dsaime/goods-and-projects/internal/controller/http2"
)

var (
	ErrJsonMarshalResponseData = errors.New("json marshal response data")
	ErrWriteResponseBytes      = errors.New("write response bytes")
	ErrParseRequestURL         = errors.New("parse request url")
)

func (c *Router) modulation(handle http2.HandlerFuncRW) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			respData   any
			b          []byte
			err        error
			httpStatus int
		)

		// Получить значения из URL
		if err = r.ParseForm(); err != nil {
			log(slog.Warn, r, errors.Join(ErrParseRequestURL, err))
		}

		// Выполнить обработку запроса
		respData, err = handle(&rwContext{
			request:  r,
			writer:   w,
			services: c.Services,
		})
		if err != nil {
			// Если есть ошибка
			respData, httpStatus = errHttpResponse(err)
			w.WriteHeader(httpStatus)
		}

		// Если ответ это редирект, выполнить редирект
		if redirect, ok := respData.(http2.Redirect); ok {
			http.Redirect(w, r, redirect.URL, redirect.Code)
			return
		}

		// Если ответ это строка, перезаписать структурой
		if s, ok := respData.(string); ok {
			// Если ответ это строка
			respData = ResponseMsg{Message: s}
		}

		// Сериализация ответа
		if b, err = json.Marshal(respData); err != nil {
			err = errors.Join(ErrJsonMarshalResponseData, err)
			respData, httpStatus = errHttpResponse(err)
			w.WriteHeader(httpStatus)
			b, _ = json.Marshal(respData)
			log(slog.Error, r, err)
		}

		// Отправить ответ
		w.Header().Set("Content-Type", "application/json")
		if _, err = w.Write(b); err != nil {
			err = errors.Join(ErrWriteResponseBytes, err)
			_, httpStatus = errHttpResponse(err)
			w.WriteHeader(httpStatus)
			log(slog.Error, r, err)
		}
	}
}

func log(lvlFn func(msg string, args ...any), r *http.Request, err error) {
	lvlFn("modulation: "+err.Error(),
		slog.String("url", r.RequestURI),
		slog.String("host", r.Host),
		slog.String("referer", r.Referer()),
	)
}

type ResponseError struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}

type ResponseMsg struct {
	Message string `json:"message"`
}
