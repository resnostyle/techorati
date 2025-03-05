package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	_ "github.com/resnostyle/techorati/internal/server/inventory"
)

type UcibiServer struct {
	Server *http.Server
	Router *echo.Echo
}

func NewServer(port string) UcibiServer {
	router := echo.New()

	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	// Prometheus
	router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%v", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	return UcibiServer{
		Server: srv,
		Router: router,
	}
}

func (me *UcibiServer) Start() {
	me.Server.Handler = me.Router
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Info().Msg(fmt.Sprintf("listening on :%v", me.Server.Addr))
		if err := me.Server.ListenAndServe(); err != nil {
			log.Err(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(15))
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	me.Server.Shutdown(ctx)

	log.Info().Msg("shutting down")
	os.Exit(0)
}