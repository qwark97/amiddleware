package amiddleware

import (
	"net/http"

	"github.com/qwark97/alog"
)

type Layer func(http.ResponseWriter, *http.Request) (interrupt bool, err *LayerError)

type LayerError struct {
	HTTPStatus int
	Message    string
}

func (le *LayerError) Error() string {
	return le.Message
}

func New(log alog.Logger) Middleware {
	return Middleware{log: log}
}

type Middleware struct {
	log    alog.Logger
	layers []Layer
}

func (mw Middleware) Use(f Layer) Middleware {
	mw.layers = append(mw.layers, f)
	return mw
}

func (mw Middleware) With(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		for _, layer := range mw.layers {
			if interrupt, err := layer(w, req); err != nil {
				mw.log.Error(err.Error())
				http.Error(w, err.Message, err.HTTPStatus)
				return
			} else if interrupt {
				return
			}
		}
		handler(w, req)
	})
}
