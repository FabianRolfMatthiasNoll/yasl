package main

import (
	"log"

	"github.com/FabianRolfMatthiasNoll/yasl/internal/repository"
)

func main() {
	// Initialize the repository
	repo, err := repository.NewRepository()
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}
	defer repo.Close()


	listID, err := repo.CreateList("Kaufland");
	if err != nil {
		log.Fatalf("failed to create list: %v", err)
	}
	log.Printf("created list: %v", listID)
	itemID, err := repo.CreateItem(listID, "Milch", "Milchprodukte");
	if err != nil {
		log.Fatalf("failed to create item: %v", err)
	}
	log.Printf("created item: %v", itemID)
}
