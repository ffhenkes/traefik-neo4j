package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	ctx := context.Background()

	uri := "bolt://traefik:7687"
	username := "neo4j"
	password := "password"

	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatalf("Error creating driver: %v", err)
	}
	defer func() {
		if err = driver.Close(ctx); err != nil {
			log.Printf("Erro closing driver: %v", err)
		}
	}()

	if err = driver.VerifyConnectivity(ctx); err != nil {
		log.Fatalf("Unable to connect to Neo4j: %v", err)
	}
	fmt.Println("Connection succesful with Neo4j via Traefik.")

	// Execute queries in different sessions
	// Each iteration opens new session thus a new connection
	for i := 0; i < 5; i++ {
		// Pause for logs inspect
		time.Sleep(2 * time.Second)
		fmt.Printf("Performing query %d...\n", i+1)

		session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
		defer session.Close(ctx)

		result, err := session.Run(ctx, "MATCH (s:Server) RETURN s", nil)
		if err != nil {
			log.Printf("Error performing query: %v", err)
			session.Close(ctx)
			continue
		}

		// Print results
		for result.Next(ctx) {
			record := result.Record()
			sValue, found := record.Get("s")
			if !found {
				fmt.Println("Bad record.. ")
				continue
			}
			fmt.Printf("Results: %v\n", sValue)
		}

		if err = result.Err(); err != nil {
			log.Printf("Error inside loop: %v", err)
		}

		session.Close(ctx)
	}
	fmt.Println("Test tear down!")
}
