package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

// ResponseConstraint is custom constraint type adapter response.
type ResponseConstraint interface {
	any
}

// Adapter is wrapper func type for error handler.
type Adapter[Response ResponseConstraint] func(w http.ResponseWriter, r *http.Request) (Response, error)

// HandlerAdapter is middleware handler to process error.
func HandlerAdapter[ResponseType ResponseConstraint](a Adapter[ResponseType]) http.HandlerFunc {
	null := make(map[string]interface{})
	response := &Response{
		Version: Version{
			Label:  "v1",
			Number: "0.1.0",
		},
		Data:       null,
		Pagination: Pagination{},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if ver, ok := r.Context().Value(CtxVersion).(Version); ok {
			response.Version = ver
		}
		payload, err := a(w, r)
		if err != nil {
			code, ok := r.Context().Value(CtxStatusCode).(int)
			if !ok || code < 1 {
				code = http.StatusInternalServerError
			}
			w.Header().Set(ContentType.String(), ApplicationJSON.String())
			w.Header().Set(XContentTypOptions.String(), "nosniff")
			response.Meta = Meta{
				Code:    strconv.Itoa(code),
				Message: err.Error(),
			}
			bytes, err := json.Marshal(response)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, err = w.Write(bytes)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
		if pagination, ok := r.Context().Value(CtxPagination).(Pagination); ok {
			response.Pagination = pagination
		}
		response.Data = payload
		*r = *r.WithContext(context.WithValue(r.Context(), CtxPayload, response))
	}
}

// SemanticVersion is middleware handler to semantic versioning.
func SemanticVersion(label string, version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*r = *r.WithContext(context.WithValue(r.Context(), CtxVersion, Version{
				Label:  label,
				Number: version,
			}))
			next.ServeHTTP(w, r)
		})
	}
}
