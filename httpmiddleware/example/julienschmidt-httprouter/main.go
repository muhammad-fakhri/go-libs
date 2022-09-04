package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	middleware "github.com/muhammad-fakhri/go-libs/httpmiddleware"
	logger "github.com/muhammad-fakhri/go-libs/log"
	"github.com/rs/cors"
)

type AStruct struct {
	Name string
}

func main() {
	logger := logger.NewSLogger("test middleware ingress")
	logIngressMid := middleware.NewIngressLogMiddleware(logger)

	router := httprouter.New()
	router.PanicHandler = RecoverWrap // this panic handler will override the panic handler in ingress middleware

	router.GET("/hello", Hello)
	router.GET("/hello-panic", HelloPanic)
	handler := logIngressMid.Enforce(router)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "Authorization", "X-Tenant", "X-User-ID"},
	})
	srv := &http.Server{
		Addr:    ":8181",
		Handler: c.Handler(handler),
	}

	logger.Info(context.TODO(), "web starting in 8181")
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalln("failed starting web on", srv.Addr, err)
	}
}

func Hello(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("some", "value")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello"))
}

func HelloPanic(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	time.Sleep(123 * time.Millisecond)
	testPanic(nil)

	w.Header().Set("some", "value")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello"))
}

func testPanic(as *AStruct) {
	log.Println(as.Name)
}

// RecoverWrap custom panic handler
func RecoverWrap(w http.ResponseWriter, r *http.Request, err interface{}) {
	log.Println("[wrap] recovered wrap ", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("service currently unavailable"))
}
