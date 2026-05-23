package transport

import (
	"log"
	"net/http"
	"time"

	"github.com/doter2410/avito-project/internal/courier"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Logger     *log.Logger
	HttpServer *http.Server
}

func NewServer(port *string, logI *log.Logger, courierHandler *courier.Handler) *Server {

	r := chi.NewRouter()

	r.Post("/couriers", courierHandler.CreateCourier)
	r.Get("/couriers/{id}", courierHandler.GetCourier)
	r.Get("/couriers", courierHandler.GetAllCouriers)
	r.Put("/couriers/{id}", courierHandler.PutUpdCourier)

	var server = &Server{
		Logger: logI,
		HttpServer: &http.Server{
			Addr:         ":" + *port,
			Handler:      r,
			ErrorLog:     logI,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}
	return server
}
