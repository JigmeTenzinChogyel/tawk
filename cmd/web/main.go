package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/JigmeTenzinChogyel/tawk/ent"
	"github.com/JigmeTenzinChogyel/tawk/ent/user"
	"github.com/JigmeTenzinChogyel/tawk/graph"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	_ "github.com/lib/pq"
)

func main() {


	// set up data base
	client := openDB()


	// server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/query", graphqlHandler())
	e.GET("/", playgroundHandler())
	
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })
	
	// e.GET("/user/create", func(ctx echo.Context) error {
	// 	c := context.Background()

	// 	u, err := CreateUser(c, client)
	// 	if err!=nil {
	// 		log.Fatal(err)
	// 	}

	// 	return ctx.String(http.StatusAccepted, u.Name)
	// })

	e.GET("/users", func(ctx echo.Context) error {
		c := context.Background()

		u, err := QueryUser(c, client)
		if err!=nil {
			log.Fatal(err)
		}

		return ctx.String(http.StatusAccepted, u.Name)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

// Defining the Graphql handler
func graphqlHandler() echo.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	// Server setup:
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})

	// h.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	// h.SetQueryCache(lru.New[*ast.QueryDocument, ](1000))

	h.Use(extension.Introspection{})
	// h.Use(extension.AutomaticPersistedQuery{
	// 	Cache: lru.New[string](100),
	// })

	return func(ctx echo.Context) error {
		h.ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// Defining the Playground handler
func playgroundHandler() echo.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")
	
	return func(ctx echo.Context) error {
		h.ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}


func openDB() (*ent.Client) {
	db, err := ent.Open("postgres","host=localhost port=5432 user=mybhutan dbname=tawk password=admin sslmode=disable")
    if err != nil {
        log.Fatalf("failed opening connection to postgres: %v", err)
    }

    // Run the auto migration tool.
    if err := db.Schema.Create(context.Background()); err != nil {
        log.Fatalf("failed creating schema resources: %v", err)
    }

	return db
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
    u, err := client.User.
        Create().
        SetAge(30).
        SetName("a8m").
        Save(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed creating user: %w", err)
    }
    log.Println("user was created: ", u)
    return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
    u, err := client.User.
        Query().
        Where(user.Name("a8m")).
        // `Only` fails if no user found,
        // or more than 1 user returned.
        Only(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed querying user: %w", err)
    }
    log.Println("user returned: ", u)
    return u, nil
}