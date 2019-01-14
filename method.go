package jrpc

import (
	"errors"
	"sync"
)

type (
	MethodRepository struct {
		m sync.RWMutex
		r map[string]Metadata
	}
	Metadata struct {
		Handler Handler
		Params  interface{}
		Result  interface{}
	}
)

// NewMethodRepository returns new MethodRepository.
func NewMethodRepository() *MethodRepository {
	return &MethodRepository{
		m: sync.RWMutex{},
		r: make(map[string]Metadata),
	}
}

// TakeHandler takes jsonrpc.Func in MethodRepository.
func (mr *MethodRepository) TakeHandler(r *Request) (Handler, *Error) {
	if r.Method == "" || r.Version != Version {
		return nil, ErrInvalidParams()
	}

	mr.m.RLock()
	md, ok := mr.r[r.Method]
	mr.m.RUnlock()
	if !ok {
		return nil, ErrMethodNotFound()
	}

	return md.Handler, nil
}

func (mr *MethodRepository) RegisterHandler(method string, h Handler, params, result interface{}) {
	if method == "" || h == nil {
		panic(errors.New("jrpc: handler name and function should not be empty"))
	}
	mr.m.Lock()
	mr.r[method] = Metadata{
		Handler: h,
		Params:  params,
		Result:  result,
	}
	mr.m.Unlock()
}

func (mr *MethodRepository) RegisterExHandler(h ExHandler) {
	mr.RegisterHandler(h.Name(), h, h.Params(), h.Result())
}

// Methods returns registered methods.
func (mr *MethodRepository) Methods() map[string]Metadata {
	mr.m.RLock()
	ml := make(map[string]Metadata, len(mr.r))
	for k, md := range mr.r {
		ml[k] = md
	}
	mr.m.RUnlock()
	return ml
}
