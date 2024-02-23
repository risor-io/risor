package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
)

func ListenAndServe(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if err := arg.RequireRange("http.listen_and_serve", 1, 2, args); err != nil {
		return err
	}
	addr, errObj := object.AsString(args[0])
	if errObj != nil {
		return errObj
	}
	callFn, ok := object.GetCloneCallFunc(ctx)
	if !ok {
		return object.Errorf("http.listen_and_serve: no clone-call function found in context")
	}
	var handler http.Handler
	if numArgs == 2 {
		switch fn := args[1].(type) {
		case http.Handler:
			handler = fn
		case *object.Function:
			handler = HandlerFunc(fn, callFn)
		default:
			return object.Errorf("type error: unsupported http handler type: %s", fn.Type())
		}
	} else {
		handler = http.DefaultServeMux
	}

	var wg sync.WaitGroup
	var listenErr error
	server := &http.Server{Addr: addr, Handler: handler}
	wg.Add(1)
	go func() {
		defer wg.Done()
		listenErr = server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
	case <-stop:
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return object.NewError(err)
	}
	wg.Wait()
	if listenErr != nil {
		return object.NewError(listenErr)
	}
	return object.Nil
}

func ListenAndServeTLS(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if err := arg.RequireRange("http.listen_and_serve_tls", 3, 4, args); err != nil {
		return err
	}
	addr, errObj := object.AsString(args[0])
	if errObj != nil {
		return errObj
	}
	certFile, errObj := object.AsString(args[1])
	if errObj != nil {
		return errObj
	}
	keyFile, errObj := object.AsString(args[2])
	if errObj != nil {
		return errObj
	}
	callFn, ok := object.GetCloneCallFunc(ctx)
	if !ok {
		return object.Errorf("http.listen_and_serve_tls: no clone-call function found in context")
	}
	var handler http.Handler
	if numArgs == 4 {
		switch fn := args[3].(type) {
		case http.Handler:
			handler = fn
		case *object.Function:
			handler = HandlerFunc(fn, callFn)
		default:
			return object.Errorf("type error: unsupported http handler type: %s", fn.Type())
		}
	} else {
		handler = http.DefaultServeMux
	}

	var wg sync.WaitGroup
	var listenErr error
	server := &http.Server{Addr: addr, Handler: handler}
	wg.Add(1)
	go func() {
		defer wg.Done()
		listenErr = server.ListenAndServeTLS(certFile, keyFile)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return object.NewError(err)
	}
	wg.Wait()
	if listenErr != nil {
		return object.NewError(listenErr)
	}
	return object.Nil
}

func HandlerFunc(fn *object.Function, callFunc object.CallFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := NewResponseWriter(w)
		req := NewRequest(r)
		result, err := callFunc(r.Context(), fn, []object.Object{res, req})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		switch result := result.(type) {
		case *object.Error:
			http.Error(w, result.Value().Error(), http.StatusInternalServerError)
		case *object.String,
			*object.ByteSlice,
			*object.Map,
			*object.List:
			// Map and list objects will be converted to JSON, while strings and
			// byte slices will be written as-is.
			res.Write(result)
		case *object.NilType, *object.Int:
			// Nothing more to do when the result is nil or an int. An int is
			// treated as a special case because it will be the return value of
			// a handler that ends with a w.write() call, which returns an int.
		default:
			http.Error(w, "type error: unsupported http handler return type",
				http.StatusInternalServerError)
		}
	})
}
