package main

import (
	"log"
	"time"
)

// Stillinger - en liste med stillingsannonser.
type Stllinger struct {
	Content []struct {
		UUID          string    `json:"uuid"`
		Published     time.Time `json:"published"`
		Expires       time.Time `json:"expires"`
		WorkLocations []struct {
			Country    string `json:"country"`
			City       string `json:"city"`
			PostalCode string `json:"postalCode"`
			County     string `json:"county"`
			Municipal  string `json:"municipal"`
		} `json:"workLocations"`
		Title          string   `json:"title"`
		Description    string   `json:"description"`
		Source         string   `json:"source"`
		ApplicationDue string   `json:"applicationDue"`
		Occupations    []string `json:"occupations"`
		Link           string   `json:"link"`
		Employer       struct {
			Name        string `json:"name"`
			Orgnr       string `json:"orgnr"`
			Description string `json:"description"`
		} `json:"employer"`
	} `json:"content"`
	TotalElements int    `json:"totalElements"`
	PageNumber    int    `json:"pageNumber"`
	PageSize      int    `json:"pageSize"`
	TotalPages    int    `json:"totalPages"`
	First         bool   `json:"first"`
	Last          bool   `json:"last"`
	Sort          string `json:"sort"`
}

func main() {
	log.Println("hello")
}
