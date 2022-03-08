package main

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
	"gqlgen-ent/config"
	"gqlgen-ent/ent"
	"gqlgen-ent/ent/migrate"
	"gqlgen-ent/graph"
	"gqlgen-ent/graph/generated"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultPort = "8080"

func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
	// "Tesla"というモデルの車を新しく作成します
	tesla, err := client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", tesla)

	// "Ford"というモデルの車を新しく作成します
	ford, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", ford)

	// 新しいユーザーを作成し、2台の車を所有させます
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		AddCars(tesla, ford).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", a8m)
	return a8m, nil
}

func QueryCars(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}
	log.Println("returned cars:", cars)

	for _, car := range cars {
		owner, err := car.QueryOwner().Only(ctx)
		if err != nil {
			return fmt.Errorf("filed querying car %q owner: %w", car.Model, err)
		}
		log.Printf("car %q owner: %q\n", car.Model, owner.Name)
	}
	return nil
}

func main() {
	url := config.CreateDatabaseConnectionURL()
	log.Println(url)
	client, err := ent.Open(config.RetrieveDialect(), url)
	if err != nil {
		log.Fatalf("failed connecting to postgres: %v", err)
	}
	defer client.Close()
	ctx := context.Background()
	if err := client.Schema.Create(
		ctx,
		migrate.WithGlobalUniqueID(true)); err != nil {
		log.Fatalf("failed creating schema resource: %v", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
