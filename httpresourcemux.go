package httpresourcemux

import "net/http"

type resourceMux struct {
	ID string
	CreateMethod string
	ReadMethod string
	UpdateMethod string
	DeleteMethod string
	ListMethod string
	CreateHandler http.Handler
	ReadHandler http.Handler
	ListHandler http.Handler
	UpdateHandler http.Handler
	DeleteHandler http.Handler
}

type ResourceMuxOptions struct {
	CreateMethod string
	ReadMethod string
	UpdateMethod string
	DeleteMethod string
	ListMethod string
}

type ResourceHandlers struct {
	CreateHandler http.Handler
	ReadHandler http.Handler
	ListHandler http.Handler
	UpdateHandler http.Handler
	DeleteHandler http.Handler
}

func NewMux(ID string, opts *ResourceMuxOptions, handlers ResourceHandlers) *resourceMux {
	mux := &resourceMux{
		ID: ID,
		CreateMethod: http.MethodPost,
		ReadMethod: http.MethodGet,
		UpdateMethod: http.MethodPut,
		DeleteMethod: http.MethodDelete,
		ListMethod: http.MethodGet,
		CreateHandler: handlers.CreateHandler,
		ReadHandler: handlers.ReadHandler,
		ListHandler: handlers.ListHandler,
		UpdateHandler: handlers.UpdateHandler,
		DeleteHandler: handlers.DeleteHandler,
	}

	if opts != nil {
		if opts.CreateMethod != "" {
			mux.CreateMethod = opts.CreateMethod
		}

		if opts.ReadMethod != "" {
			mux.CreateMethod = opts.ReadMethod
		}

		if opts.UpdateMethod != "" {
			mux.CreateMethod = opts.UpdateMethod
		}

		if opts.DeleteMethod != "" {
			mux.CreateMethod = opts.UpdateMethod
		}

		if opts.ListMethod != "" {
			mux.CreateMethod = opts.ListMethod
		}
	}

	return mux
}

func (mux *resourceMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if mux.ID != "" {
		if r.Method == mux.UpdateMethod {
			serve(mux.UpdateHandler,w,r)
			return
		}

		if r.Method == mux.ReadMethod {
			serve(mux.ReadHandler,w,r)
			return
		}

		if r.Method == mux.DeleteMethod {
			serve(mux.DeleteHandler,w,r)
			return
		}
	}

	if r.Method == mux.CreateMethod {
		serve(mux.CreateHandler,w,r)
		return
	}

	if r.Method == mux.ListMethod {
		serve(mux.ListHandler,w,r)
		return
	}

	http.NotFound(w,r)
}

func serve(handler http.Handler, w http.ResponseWriter, r *http.Request) {
	if handler != nil {
		handler.ServeHTTP(w,r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
