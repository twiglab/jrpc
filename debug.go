package jrpc

import (
	"encoding/json"
	"net/http"
	"reflect"
)

// A MethodReference is a reference of JSON-RPC method.
type MethodReference struct {
	Name    string  `json:"name"`
	Handler string  `json:"handler"`
	Params  *Schema `json:"params,omitempty"`
	Result  *Schema `json:"result,omitempty"`
}

// ServeDebug views registered method list.
func (mr *MethodRepository) ServeDebug(w http.ResponseWriter, r *http.Request) { // nolint: unparam
	ms := mr.Methods()
	if len(ms) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	l := make([]*MethodReference, 0, len(ms))
	for k, md := range ms {
		l = append(l, makeMethodReference(k, md))
	}
	w.Header().Set(contentTypeKey, contentTypeValue)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	if err := enc.Encode(l); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func makeMethodReference(k string, md Metadata) *MethodReference {
	mr := &MethodReference{
		Name: k,
	}
	tv := reflect.TypeOf(md.Handler)
	if tv.Kind() == reflect.Ptr {
		tv = tv.Elem()
	}
	mr.Handler = tv.Name()
	if md.Params != nil {
		mr.Params = Reflect(md.Params)
	}
	if md.Result != nil {
		mr.Result = Reflect(md.Result)
	}
	return mr
}

type DebugHandler struct {
	*MethodRepository
}

func Debug(mr *MethodRepository) *DebugHandler {
	return &DebugHandler{
		MethodRepository: mr,
	}
}

func (d *DebugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.ServeDebug(w, r)
}
