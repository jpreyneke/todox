package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewDB,
		NewRepository,
		NewService,
		NewHandler,
		NewRouter,
	),
	fx.Invoke(StartServer),
)

func NewDB() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		GetEnv("DB_USER", "user"),
		GetEnv("DB_PASSWORD", "password"),
		GetEnv("DB_HOST", "localhost"),
		GetEnv("DB_PORT", "3306"),
		GetEnv("DB_NAME", "todox"),
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		slog.Error("Failed to connect to database",
			"error", err,
			"host", GetEnv("DB_HOST", "localhost"),
			"port", GetEnv("DB_PORT", "3306"),
		)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	maxOpen := GetEnvInt("DB_MAX_OPEN_CONNS", 25)
	maxIdle := GetEnvInt("DB_MAX_IDLE_CONNS", 5)
	maxLifetime := GetEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute)

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(maxLifetime)

	slog.Info("Connected to database",
		"host", GetEnv("DB_HOST", "localhost"),
		"max_open_conns", maxOpen,
		"max_idle_conns", maxIdle,
		"conn_max_lifetime", maxLifetime,
	)
	return db, nil
}

func NewRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())
	r.Use(MetricsMiddleware())
	r.Use(AuthMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	handler.RegisterRoutes(r)

	return r
}

func StartServer(lc fx.Lifecycle, router *gin.Engine) {
	port := GetEnv("PORT", "8080")

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("Starting HTTP server",
				"port", port,
				"read_timeout", srv.ReadTimeout,
				"write_timeout", srv.WriteTimeout,
			)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					slog.Error("Server error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("Shutting down server gracefully...")
			return srv.Shutdown(ctx)
		},
	})
}
