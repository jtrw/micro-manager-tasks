package server

import (
	"context"
	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jtrw/go-rest"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	taskHandler "micro-manager-tasks/m/v2/app/handler"
	"net/http"
	"time"
	//"fmt"
)

type Server struct {
	Listen         string
	PinSize        int
	MaxPinAttempts int
	MaxExpire      time.Duration
	WebRoot        string
	Secret         string
	Version        string
	Client         *mongo.Client
}

func (s Server) Run(ctx context.Context) error {
	log.Printf("[INFO] activate rest server")
	log.Printf("[INFO] Listen: %s", s.Listen)

	httpServer := &http.Server{
		Addr:              s.Listen,
		Handler:           s.routes(),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	go func() {
		<-ctx.Done()
		if httpServer != nil {
			if clsErr := httpServer.Close(); clsErr != nil {
				log.Printf("[ERROR] failed to close proxy http server, %v", clsErr)
			}
		}
	}()

	err := httpServer.ListenAndServe()
	log.Printf("[WARN] http server terminated, %s", err)

	if err != http.ErrServerClosed {
		return errors.Wrap(err, "server failed")
	}
	return err
}

func (s Server) routes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.RealIP)
	router.Use(middleware.Throttle(1000), middleware.Timeout(60*time.Second))
	//router.Use(rest.AppInfo("events", "Jrtw", s.Version), rest.Ping)
	router.Use(rest.Ping)
	router.Use(tollbooth_chi.LimitHandler(tollbooth.NewLimiter(10, nil)))
	router.Use(middleware.Logger)

	handler := taskHandler.NewHandler(s.Client)
	router.Route(
		"/api/v1", func(r chi.Router) {
			r.Post("/tasks", handler.CreateTask)
			r.Get("/tasks/{uuid}", handler.onShowTaskInfo)
			r.Post("/subtask", handler.AddSubTask)
			r.Get("/status", handler.CheckStatus)
		},
	)

	router.Get(
		"/robots.txt", func(w http.ResponseWriter, r *http.Request) {
			render.PlainText(w, r, "User-agent: *\nDisallow: /\n")
		},
	)

	return router
}
