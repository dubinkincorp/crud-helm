package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"
	"regexp"
)

type UnexpectedError struct {
	code    int
	message string
}

func main() {
	pool, err := openDbConnection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	//delayed function
	defer pool.Close(context.Background())

	fmt.Println("Start server")

	var handler RegExpHandler
	route := Route{
		pool:      pool,
		errorBody: errorBody(),
	}

	compiledPattern, err := regexp.Compile(`.*\\user\\[0-9]+`)
	if err != nil {
		log.Fatal(err)
	}

	handler.HandleFunc(compiledPattern, route.Handle)
	http.ListenAndServe(":80", nil)
}

func openDbConnection() (*pgxpool.Conn, error) {
	const envDbUrl = "DATABASE_URI"

	conn, err := pgxpool.Connect(context.Background(), os.Getenv(envDbUrl))

	return conn, err
}

func errorBody() *[]byte {
	body := UnexpectedError{0, "Something went wrong"}

	js, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	return &js
}
