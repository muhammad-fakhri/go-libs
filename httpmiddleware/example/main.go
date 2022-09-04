package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	middleware "github.com/muhammad-fakhri/go-libs/httpmiddleware"
	logger "github.com/muhammad-fakhri/go-libs/log"
)

type AStruct struct {
	Name string
}

func main() {
	logger := logger.NewSLogger("test middleware ingress")
	logIngressMid := middleware.NewIngressLogMiddleware(logger)

	handler := logIngressMid.Enforce(http.HandlerFunc(Hello))
	handlerReqID := RequestIDMiddleware(logIngressMid.Enforce(http.HandlerFunc(Hello)))
	handlerPanic := logIngressMid.Enforce(http.HandlerFunc(HelloPanic))
	// to override the ingress middleware's panic handler, place your custom handler first like this.
	handlerPanic2 := logIngressMid.Enforce(RecoverWrap(http.HandlerFunc(HelloPanic)))

	http.Handle("/hello", handler)
	http.Handle("/hello-panic", handlerPanic)
	http.Handle("/hello-panic-wrap", handlerPanic2)
	http.Handle("/hello-request-id", handlerReqID)

	logger.Info(context.TODO(), "web starting in 8181")
	if err := http.ListenAndServe(":8181", nil); err != nil {
		panic(err)
	}
}

func Hello(w http.ResponseWriter, r *http.Request) {
	responseBodyBytes, _ := ioutil.ReadAll(r.Body)
	logger := logger.NewSLogger("test middleware ingress handler")
	logger.Info(r.Context(), "check context id and country")

	// responses
	w.Header().Set("some", "value")
	w.WriteHeader(http.StatusOK)

	w.Write(responseBodyBytes)
	w.Write([]byte("check if all response being logged"))
}

func HelloPanic(w http.ResponseWriter, r *http.Request) {
	responseBodyBytes, _ := ioutil.ReadAll(r.Body)
	logger := logger.NewSLogger("test middleware ingress handler")
	logger.Info(r.Context(), "check context id and country")

	time.Sleep(123 * time.Millisecond)
	testPanic(nil)

	// responses
	w.Header().Set("some", "value")
	w.WriteHeader(http.StatusOK)

	w.Write(responseBodyBytes)
	w.Write([]byte("check if all response being logged"))
}

func testPanic(as *AStruct) {
	log.Println(as.Name)
}

// RecoverWrap custom panic handler
func RecoverWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				log.Println("[wrap] recovered wrap ", r)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("service currently unavailable"))
			}
		}()

		h.ServeHTTP(w, r)
	})
}

// RequestIDMiddleware assign request id middleware
func RequestIDMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]string, 0)
		data[logger.ContextCountryKey] = "id"
		data[logger.ContextIdKey] = "abcdefghijklmnopq"

		ctx := context.WithValue(r.Context(), logger.ContextDataMapKey, data)
		r2 := r.Clone(ctx)

		h.ServeHTTP(w, r2)
	})
}
