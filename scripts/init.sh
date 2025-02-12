#!/bin/bash
set -e

echo "Waiting Neo4j r1..."
until cypher-shell -a neo4j-r1:7687 -u neo4j -p password "RETURN 1" &>/dev/null; do
  sleep 2
done

echo "Waiting Neo4j r2..."
until cypher-shell -a neo4j-r2:7687 -u neo4j -p password "RETURN 1" &>/dev/null; do
  sleep 2
done

echo "Creating data in neo4j-r1..."
cypher-shell -a neo4j-r1:7687 -u neo4j -p password "CREATE (:Server { id: 'r1', name: 'neo4j-r1'});"

echo "Creating data in neo4j-r2..."
cypher-shell -a neo4j-r2:7687 -u neo4j -p password "CREATE (:Server { id: 'r2', name: 'neo4j-r2'});"

echo "Init complete!"
tail -f /dev/null