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
	// Traefik address; using service name since we're on the same docker-compose network
	uri := "bolt://traefik:7687"
	username := "neo4j"
	password := "password"

	// For each query, create a new driver instance to force a new connection
	for i := 0; i < 5; i++ {
		// Sleep for 2 seconds to allow log visualization and avoid overwhelming connections
		time.Sleep(2 * time.Second)
		fmt.Printf("Performing query %d...\n", i+1)

		// Create a new driver instance for each iteration, which forces a new TCP connection
		driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
		if err != nil {
			log.Fatalf("Error creating driver: %v", err)
		}

		// Verify connectivity with Neo4j
		if err = driver.VerifyConnectivity(ctx); err != nil {
			log.Fatalf("Cannot connect to Neo4j: %v", err)
		}

		// Create a new session (from the new driver) with read access mode
		session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})

		// Execute the query: "MATCH (s:Server) RETURN s"
		result, err := session.Run(ctx, "MATCH (s:Server) RETURN s", nil)
		if err != nil {
			log.Printf("Error executing query: %v", err)
			session.Close(ctx)
			driver.Close(ctx)
			continue
		}

		// Iterate over returned records and print the results
		for result.Next(ctx) {
			record := result.Record()
			sValue, found := record.Get("s")
			if !found {
				fmt.Println("Record does not contain key 's'")
				continue
			}
			fmt.Printf("Result: %v\n", sValue)
		}

		// Check for errors during result iteration
		if err = result.Err(); err != nil {
			log.Printf("Error iterating results: %v", err)
		}

		// Close the session and driver to ensure the connection is terminated
		if err = session.Close(ctx); err != nil {
			log.Printf("Error closing session: %v", err)
		}
		if err = driver.Close(ctx); err != nil {
			log.Printf("Error closing driver: %v", err)
		}
	}

	fmt.Println("Test tear down!")
}
