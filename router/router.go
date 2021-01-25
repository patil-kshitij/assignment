package router

import (
	"assignment/router/muxrouter"
	"net/http"
)

func NewRouter() http.Handler {
	return muxrouter.GetMuxRouter()
}
