package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {

	ctx := c.Request.Context()

	fact, err := fetchCatFact(ctx)
	if err != nil {
		log.Printf("Error fetching cat fact: %v", err)
		fact = "Unable to fetch cat fact at this time. Did you know cats are amazing creatures?"
	}

	c.JSON(200, gin.H{
		"status": "success",
		"user": gin.H{
			"email": "isaacchimdi@gmail.com",
			"name":  "Isaac Ubani",
			"stack": "Go/Gin",
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"fact":      fact,
	})

}

func fetchCatFact(ctx context.Context) (string, error) {

	catFactAPIURL := "https://catfact.ninja/fact"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, catFactAPIURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch cat fact: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cat fact API returned status: %d", resp.StatusCode)
	}

	type CatFactResponse struct {
		Fact   string `json:"fact"`
		Length int    `json:"length"`
	}
	var catFact CatFactResponse
	if err := json.NewDecoder(resp.Body).Decode(&catFact); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return catFact.Fact, nil
}
