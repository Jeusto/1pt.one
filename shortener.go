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

type DbEntry struct {
	ShortURL       string `bson:"short_url" json:"short_url"`
	LongURL        string `bson:"long_url" json:"long_url"`
	CreatedAt      string `bson:"created_at" json:"created_at"`
	NumberOfVisits int    `bson:"number_of_visits" json:"number_of_visits"`
	FileContent    string `bson:"file_content" json:"file_content"`
}

func getURL(shortURL string) (*DbEntry, error) {
	result := collection.FindOne(context.Background(), bson.M{"short_url": shortURL})

	fmt.Println(result)
	if result.Err() != nil {
		return nil, fmt.Errorf("The provided short URL does not exist")
	}

	var item DbEntry
	err := result.Decode(&item)
	if err != nil {
		return nil, err
	}

	return &item, nil
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

func addURL(shortURL, longURL string) (*DbEntry, error) {
	// Check if URL is valid
	if !urlIsValid(longURL) {
		return nil, fmt.Errorf("The provided URL is not valid")
	}

	// Check if short URL doesnt already exist
	if shortURL != "" && urlExists(shortURL) {
		return nil, fmt.Errorf("The provided short URL already exists")
	}

	// Generate a new short URL if none is provided
	if shortURL == "" {
		shortURL = generateNanoid()
	}

	// Prepend "http://" if URL doesn't have a scheme
	if !strings.HasPrefix(longURL, "http://") && !strings.HasPrefix(longURL, "https://") {
		longURL = "http://" + longURL
	}

	// Insert new URL into database
	url := DbEntry{
		ShortURL:       shortURL,
		LongURL:        longURL,
		CreatedAt:      time.Now().Format("2006-01-02 15:04:05"),
		NumberOfVisits: 0,
	}

	_, err := collection.InsertOne(context.Background(), url)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func addFile(fileContent string) (*DbEntry, error) {
	// Generate a new short URL
	shortURL := generateNanoid()

	// Insert new URL into database
	url := DbEntry{
		ShortURL:       shortURL,
		FileContent:    fileContent,
		CreatedAt:      time.Now().Format("2006-01-02 15:04:05"),
		NumberOfVisits: 0,
	}

	_, err := collection.InsertOne(context.Background(), url)
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

func generateNanoid() string {
	var id string
	var err error

	for {
		id, err = gonanoid.Nanoid(4)
		if err != nil {
			panic(err)
		}

		if !urlExists(id) {
			break
		}
	}

	return id
}

func shortUrlIsValid(shortIdentifier string) bool {
	pattern := `^[a-zA-Z0-9-_]*$`
	match, err := regexp.MatchString(pattern, shortIdentifier)

	if err != nil {
		return false
	}

	return match
}
