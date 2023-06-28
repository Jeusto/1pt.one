package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/matoous/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
)

func getURL(shortURL string) (*URL, error) {
	result := collection.FindOne(context.Background(), bson.M{"short_url": shortURL})
	if result.Err() != nil {
		return nil, fmt.Errorf("The provided short URL does not exist")
	}

	var url URL
	err := result.Decode(&url)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func incrementUrlVisits(shortURL string) error {
	url, err := getURL(shortURL)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(context.Background(),
		bson.M{"short_url": shortURL},
		bson.M{"$set": bson.M{"number_of_visits": url.NumberOfVisits + 1}},
	)

	if err != nil {
		return err
	}

	return nil
}

func addURL(shortURL, longURL string) (*URL, error) {
	// Check if short URL doesnt already exist
	if urlExists(shortURL) {
		return nil, fmt.Errorf("The provided short URL already exists")
	}

	// Prepend "http://" if URL doesn't have a scheme
	if !strings.HasPrefix(longURL, "http://") && !strings.HasPrefix(longURL, "https://") {
		longURL = "http://" + longURL
	}

	// Check if URL is valid
	if !urlIsValid(longURL) {
		return nil, fmt.Errorf("The provided URL is not valid")
	}

	// Generate a new short URL if none is provided
	if shortURL == "" {
		for {
			shortURL = generateShortURL()
			if !urlExists(shortURL) {
				break
			}
		}

	}

	// Add the new URL to database
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	uniqueId, err := gonanoid.Nanoid()
	if err != nil {
		return nil, err
	}

	url := URL{
		ID:             uniqueId,
		ShortURL:       shortURL,
		LongURL:        longURL,
		CreatedAt:      currentTime,
		NumberOfVisits: 0,
	}

	_, err = collection.InsertOne(context.Background(), url)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func urlExists(shortURL string) bool {
	result := collection.FindOne(context.Background(), bson.M{"short_url": shortURL})
	return result.Err() == nil
}

func urlIsValid(url string) bool {
	pattern := `(https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z0-9]{2,}(\.[a-zA-Z0-9]{2,})(\.[a-zA-Z0-9]{2,})?`

	match, err := regexp.MatchString(pattern, url)

	fmt.Println(match)
	if err != nil {
		return false
	}

	return match
}

func generateShortURL() string {
	id, err := gonanoid.Nanoid(4)

	if err != nil {
		panic(err)
	}

	return id
}
