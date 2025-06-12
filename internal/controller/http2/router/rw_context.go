package router

import (
	"net/http"

	"github.com/dsaime/goods-and-projects/internal/controller/http2"
)

// rwContext представляет контекст HTTP-запроса
type rwContext struct {
	requestID string
	request   *http.Request
	writer    http.ResponseWriter
	services  http2.Services
}

func (r *rwContext) Writer() http.ResponseWriter {
	return r.writer
}

func (r *rwContext) RequestID() string {
	return r.requestID
}

func (r *rwContext) Request() *http.Request {
	return r.request
}

func (r *rwContext) Services() http2.Services {
	return r.services
}

func (r *rwContext) SetRequestID(id string) {
	r.requestID = id
}

func (r *rwContext) SetRequest(req *http.Request) {
	r.request = req
}
