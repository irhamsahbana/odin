package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

// ResponseType - Custom type to hold value for find and replace on context value response type.
type ResponseType int

// Declare related constants for each ResponseType starting with index 1.
const (
	CtxPagination ResponseType = iota
	CtxVersion
	CtxPayload
	CtxStatusCode
)

func (r ResponseType) String() string {
	return [...]string{
		"pagination-key",
		"version-key",
		"payload-key",
		"status-code-key",
	}[r]
}

// Index - Return index of the Constant.
func (r ResponseType) Index() int {
	return int(r)
}

// Meta holds the response definition for the Meta entity.
type Meta struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"error_message,omitempty"`
}

// Version holds the response definition for the Version entity.
type Version struct {
	Label  string `json:"label,omitempty"`
	Number string `json:"number,omitempty"`
}

// Pagination holds the response definition for the Pagination entity.
type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"per_page"`
	Size  int `json:"page_count,omitempty"`
	Total int `json:"total_count,omitempty"`
}

// Response holds the response definition for the Response entity.
type Response struct {
	Meta       `json:"meta"`
	Version    `json:"version"`
	Pagination `json:"pagination,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

// JSON sends a JSON response with status code.
func JSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
		code, ok := r.Context().Value(CtxStatusCode).(int)
		if !ok || code < 1 {
			code = http.StatusOK
		}
		if payload, ok := r.Context().Value(CtxPayload).(*Response); ok {
			payload.Meta = Meta{
				Code: strconv.Itoa(code),
			}
			buf := &bytes.Buffer{}
			enc := json.NewEncoder(buf)
			enc.SetEscapeHTML(true)
			if err := enc.Encode(payload); err != nil {
				log.Error().Err(ErrInternalServerError(w, r, err)).Msg("JSON")
				return
			}
			w.WriteHeader(code)
			w.Header().Set(ContentType.String(), ApplicationJSON.String())
			_, err := w.Write(buf.Bytes())
			if err != nil {
				log.Error().Err(ErrInternalServerError(w, r, err)).Msg("JSON")
				return
			}
			return
		}
		http.Error(w, fmt.Errorf("payload cannot cast to response").Error(), http.StatusNotImplemented)
	}
}

// Paging send a Pagination data.
func Paging(r *http.Request, p Pagination) {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxPagination, p))
}
