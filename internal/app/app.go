package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mvp-mogila/ozon-test-task/internal/delivery/http/server"
	"github.com/mvp-mogila/ozon-test-task/internal/pkg/config"
	inmemrepo "github.com/mvp-mogila/ozon-test-task/internal/repository/inmemory"
	postgresrepo "github.com/mvp-mogila/ozon-test-task/internal/repository/postgres"
	"github.com/mvp-mogila/ozon-test-task/internal/usecase"
	loggerpkg "github.com/mvp-mogila/ozon-test-task/pkg/logger"
	"github.com/mvp-mogila/ozon-test-task/pkg/postgres"
)

var postgresConnTimeout = 5 * time.Second

type App struct {
	cfg     *config.Config
	server  *handler.Server
	pgxConn *pgxpool.Pool
	logger  *log.Logger
}

func NewApp() *App {
	logger := loggerpkg.NewLogger()

	cfg := config.LoadConfig()

	var (
		postUsecase    *usecase.PostUsecase
		commentUsecase *usecase.CommentUsecase
		pgxConn        *pgxpool.Pool
	)

	commentNitifier := usecase.NewCommentNotificationService()

	if cfg.UseInMemoryStorage {
		pgxConn = nil

		postRepo := inmemrepo.NewPostInMemoryRepository()
		commentRepo := inmemrepo.NewCommentInMemoryRepository()

		postUsecase = usecase.NewPostUsecase(postRepo)
		commentUsecase = usecase.NewCommentUsecase(commentRepo, postUsecase, commentNitifier)
	} else {
		pgxCfg := config.InitPostgresConfig(cfg)

		ctx, cancel := context.WithTimeout(context.Background(), postgresConnTimeout)
		defer cancel()
		pgxDatabase := postgres.NewPgxDatabase(ctx, pgxCfg)
		pgxConn = pgxDatabase

		postRepo := postgresrepo.NewPostPostgresRepository(pgxDatabase)
		commentRepo := postgresrepo.NewCommentPostgresRepository(pgxDatabase)

		postUsecase = usecase.NewPostUsecase(postRepo)
		commentUsecase = usecase.NewCommentUsecase(commentRepo, postUsecase, commentNitifier)
	}

	srv := server.NewHTTPServer(postUsecase, commentUsecase, cfg.RequestTimeout, logger)

	return &App{
		cfg:     cfg,
		server:  srv,
		pgxConn: pgxConn,
		logger:  logger,
	}
}

func (a *App) Run() error {
	log.Printf("Server is listening on %s...\n", a.cfg.Host+":"+a.cfg.Port)
	if err := http.ListenAndServe(a.cfg.Host+":"+a.cfg.Port, nil); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("Server closed")
			return nil
		}
		return err
	}
	return nil
}

func (a *App) Stop() {
	if a.pgxConn != nil {
		a.pgxConn.Close()
	}
}
