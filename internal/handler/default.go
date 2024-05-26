package handler

import "net/http"

type DefaultHandler struct{}

func (h *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome"))
}
