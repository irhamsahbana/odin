package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

type ctxStatusCodeKey struct {
	Name string
}

// String - print context value.
func (r *ctxStatusCodeKey) String() string {
	return "context value " + r.Name
}

type ctxPayloadKey struct {
	Name string
}

func (r *ctxPayloadKey) String() string {
	return "context value " + r.Name
}

type ctxVersionKey struct {
	Name string
}

func (r *ctxVersionKey) String() string {
	return "context value " + r.Name
}

var (
	// CtxVersion context key value for http status code.
	CtxVersion = ctxVersionKey{Name: "context version"}
	// CtxPayload context key value for http status code.
	CtxPayload = ctxPayloadKey{Name: "context payload"}
	// CtxStatusCode context key value for http status code.
	CtxStatusCode = ctxStatusCodeKey{Name: "context status code"}
)

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
