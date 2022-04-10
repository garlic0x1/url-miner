package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func request(u string) string {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Println("Error creating request", u)
	}

	var client http.Client

	client = http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	// perform request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error performing request", u)
		// this shows when response is nil
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bodyBytes)
}
