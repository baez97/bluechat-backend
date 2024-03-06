package main

import (
	"bluechat-server/database"
	"bluechat-server/graph"
	"bluechat-server/graph/model"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
		router := chi.NewRouter()
		router.Use(cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:8080"},
			AllowCredentials: true,
			Debug:            true,
		}).Handler)

    port := os.Getenv("PORT")
    if port == "" {
        port = defaultPort
    }

    // Open a connection to the PostgreSQL database
    db, err := sql.Open("postgres", "user=josemanuelbaezsoriano dbname=chat password=4323415 sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create the GraphQL server with the Resolver instance
    srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
        Database: database.Database{SQL: db},
        ChatObservers: map[string]chan []*model.Message{},
    }}))

		srv.AddTransport(transport.POST{})
		srv.AddTransport(transport.Websocket{
			KeepAlivePingInterval: 10 * time.Second,
			Upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		})

    // Set up routes for the GraphQL playground and GraphQL endpoint
    router.Handle("/", playground.Handler("GraphQL playground", "/query"))
    router.Handle("/query", srv)
		err = http.ListenAndServe(":8080", router)

    log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
		if err != nil {
			panic(err)
		}
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func corsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Set CORS headers
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
					w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

					// Handle OPTIONS requests (preflight requests)
					if r.Method == "OPTIONS" {
							w.WriteHeader(http.StatusOK)
							return
					}

					// Call the next handler
					next.ServeHTTP(w, r)
			})
	}
}