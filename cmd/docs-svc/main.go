package main

import (
	"log"
	"net/http"
	"time"

	"docs.svc.plus/internal/config"
	httpapi "docs.svc.plus/internal/http"
	"docs.svc.plus/internal/observability"
)

func main() {
	stopTracing, ddEnabled := observability.ConfigureDDTrace("web-saas-content")
	defer stopTracing()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	app, err := httpapi.NewApp(cfg)
	if err != nil {
		log.Fatalf("create app: %v", err)
	}

	if cfg.ReloadInterval > 0 {
		go func() {
			ticker := time.NewTicker(cfg.ReloadInterval)
			defer ticker.Stop()
			for range ticker.C {
				result := app.Reload(true)
				if !result.Reloaded {
					log.Printf("content reload failed: %s", result.Message)
				}
			}
		}()
	}

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           observability.WrapHTTPHandler(app.Routes(), "web-saas-content", "content.request", ddEnabled),
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("docs.svc.plus listening on :%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %v", err)
	}
}
