package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	uri      = "bolt://traefik:7687"
	username = "neo4j"
	password = "password"
)

func welcome(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Welcome to hollow test")
}

func queryNeo4j() {
	ctx := context.Background()
	time.Sleep(2 * time.Second)

	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)
		fmt.Printf("Performing query %d...\n", i+1)

		driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
		if err != nil {
			log.Printf("Error creating driver: %v", err)
			return
		}
		defer driver.Close(ctx)

		if err = driver.VerifyConnectivity(ctx); err != nil {
			log.Printf("Cannot connect to Neo4j: %v", err)
			return
		}

		session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
		defer session.Close(ctx)

		result, err := session.Run(ctx, "MATCH (s:Server) RETURN s", nil)
		if err != nil {
			log.Printf("Error executing query: %v", err)
			return
		}

		for result.Next(ctx) {
			record := result.Record()
			if sValue, found := record.Get("s"); found {
				fmt.Printf("Result: %v\n", sValue)
			} else {
				fmt.Println("Record does not contain key 's'")
			}
		}

		if err = result.Err(); err != nil {
			log.Printf("Error iterating results: %v", err)
		}
	}

	fmt.Println("Test tear down!")
}

func queryHandler(w http.ResponseWriter, _ *http.Request) {
	go queryNeo4j()
	fmt.Fprintln(w, "Query execution started, check logs for results..")
}

func main() {
	http.HandleFunc("/", welcome)
	http.HandleFunc("/test", queryHandler)
	fmt.Println("Server started on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
