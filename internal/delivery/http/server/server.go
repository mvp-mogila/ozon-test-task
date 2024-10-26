package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/mvp-mogila/ozon-test-task/gen/graph"
	"github.com/mvp-mogila/ozon-test-task/internal/delivery/graphql"
	"github.com/mvp-mogila/ozon-test-task/pkg/logger"
)

const defaultRequestTimeout = 5

var requestID uint64

func NewHTTPServer(pu graphql.PostsUsecase, cu graphql.CommentsUsecase, requestTimeout int, logger *log.Logger) *handler.Server {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graphql.Resolver{
			PostUsecase:    pu,
			CommentUsecase: cu,
		},
	}))

	if requestTimeout == 0 {
		requestTimeout = defaultRequestTimeout
	}

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", LoggingMiddleware(logger)(RequestTimeoutMiddleware(requestTimeout)(srv)))

	return srv
}

func RequestTimeoutMiddleware(requestTimeout int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header["Connection"] != nil && r.Header["Connection"][0] != "Upgrade" || r.Header["Upgrade"] != nil && r.Header["Upgrade"][0] != "websocket" {
				ctx, cancel := context.WithTimeout(r.Context(), time.Duration(requestTimeout)*time.Second)
				defer cancel()
				req := r.WithContext(ctx)

				doneCh := make(chan struct{})

				go func() {
					next.ServeHTTP(w, req)
					close(doneCh)
				}()

				select {
				case <-doneCh:
				case <-ctx.Done():
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusRequestTimeout)
					jsonResponse := map[string]string{
						"error": "Request was too complicated and timed out",
					}

					json.NewEncoder(w).Encode(jsonResponse)
				}
			} else {
				log.Println("Serving new websocket connection...")
				next.ServeHTTP(w, r)
			}
		})
	}
}

func LoggingMiddleware(l *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.SetPrefix(fmt.Sprintf("REQUEST %d  ", atomic.AddUint64(&requestID, 1)))
			l.Println("New request")
			newCtx := logger.AddToContext(r.Context(), l)
			req := r.WithContext(newCtx)

			next.ServeHTTP(w, req)
		})
	}
}
