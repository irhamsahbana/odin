// Package rest is port adapter via http/s protocol
// # This manifest was generated by ymir. DO NOT EDIT.
package rest

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"

	"gitlab.playcourt.id/nanang_suryadi/odin/pkg/shared/security"
)

const (
	defaultMemory = 32 << 20 // 32 MB
)

// decoderBody detects the correct decoderBody for use on an HTTP request and
// marshals into a given interface.
func decoderBody(r *http.Request, i any) (err error) {
	if r.ContentLength == 0 {
		return nil
	}
	switch GetRequestContentType(r) {
	case MIMEApplicationJSON:
		err = DecodeJSON(r.Body, i)
	case MIMETextXML:
		err = DecodeXML(r.Body, i)
	case MIMEApplicationForm, MIMEMultipartForm:
		err = DecodeForm(r, i)
	default:
		err = errors.New("unable to decode the request content type")
	}
	return err
}

// DecodeJSON decodes a given reader into an interface using the json decoderBody.
func DecodeJSON(r io.Reader, v any) error {
	defer func(dst io.Writer, src io.Reader) {
		_, err := io.Copy(dst, src)
		if err != nil {

		}
	}(io.Discard, r)
	return json.NewDecoder(r).Decode(v)
}

// DecodeXML decodes a given reader into an interface using the xml decoderBody.
func DecodeXML(r io.Reader, v any) error {
	defer func(dst io.Writer, src io.Reader) {
		_, err := io.Copy(dst, src)
		if err != nil {

		}
	}(io.Discard, r)
	return xml.NewDecoder(r).Decode(v)
}

// DecodeForm decodes a given reader into an interface using the form decoderBody.
func DecodeForm(r *http.Request, v any) error {
	if strings.HasPrefix(r.Header.Get(HeaderContentType.String()), MIMEMultipartForm.String()) {
		if err := r.ParseMultipartForm(defaultMemory); err != nil {
			return err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return err
		}
	}
	decoding := schema.NewDecoder()
	return decoding.Decode(&v, r.Form)
}

type Binder[T any] struct {
	str *T
}

// Bind implements all bind request raw data.
// Binding is done in following order: 1) path params; 2) query params; 3) request body. Each step COULD override previous
// step bind values.
func Bind[T any](r *http.Request, v *T) (*Binder[T], error) {
	binder := &Binder[T]{}
	if err := bindURLParams(r.Context(), v); err != nil {
		return nil, err
	}
	method := r.Method
	if method == http.MethodGet ||
		method == http.MethodDelete ||
		method == http.MethodHead {
		if err := bindQueryParams(r, v); err != nil {
			return nil, err
		}
	}
	if err := decoderBody(r, v); err != nil {
		return nil, err
	}
	binder.str = v
	return binder, nil
}

// Validate implements value validations for structs and individual fields based on tags.
func (b *Binder[T]) Validate() error {
	var errs []error
	validators := security.Validate(b.str)
	for n := range validators {
		v := validators[n]
		errs = append(errs, errors.New(v.Message))
	}
	if len(errs) < 1 {
		return nil
	}
	return errors.Join(errs...)
}

func bindURLParams(ctx context.Context, v any) error {
	urls := chi.RouteContext(ctx).URLParams
	names := urls.Keys
	values := urls.Values
	params := map[string][]string{}
	for i, name := range names {
		params[name] = []string{values[i]}
	}
	decoding := schema.NewDecoder()
	return decoding.Decode(v, params)
}

func bindQueryParams(r *http.Request, v any) error {
	queries := r.URL.Query()
	decoding := schema.NewDecoder()
	return decoding.Decode(v, queries)
}
