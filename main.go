package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type URL struct {
	ID             string `bson:"_id" json:"_id"`
	ShortURL       string `bson:"short_url" json:"short_url"`
	LongURL        string `bson:"long_url" json:"long_url"`
	CreatedAt      string `bson:"created_at" json:"created_at"`
	NumberOfVisits int    `bson:"number_of_visits" json:"number_of_visits"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var (
	client     *mongo.Client
	collection *mongo.Collection
)

func main() {
	// Retrieve MongoDB connection string from .env file
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	connectionString := os.Getenv("MONGODB_CONNECTION_STRING")
	if connectionString == "" {
		log.Fatal("MONGODB_CONNECTION_STRING environment variable is not set")
	}

	// Connect to MongoDB
	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverApi)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database("url-shortener")
	collection = database.Collection("urls")

	// Create Gin router
	router := gin.Default()

	router.LoadHTMLFiles("static/index.html")
	router.Static("/static", "./static")

	router.GET("/", index)
	router.GET("/status", status)
	router.GET("/:shortURL", redirect)
	router.GET("/shorten", shorten)
	router.GET("/retrieve", retrieve)

	// Start the server
	router.Run(":8080")
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func status(c *gin.Context) {
	response := gin.H{
		"status":  200,
		"message": "API is live. Read the documentation at https://github.com/Jeusto/1pt.one",
	}

	c.JSON(http.StatusOK, response)
}

func redirect(c *gin.Context) {
	shortURL := c.Param("shortURL")

	url, err := getURL(shortURL)
	if err != nil {
		c.Redirect(http.StatusFound, "/?url_not_found="+shortURL)
		return
	}

	incrementUrlVisits(shortURL)
	c.Redirect(http.StatusFound, url.LongURL)
}

func shorten(c *gin.Context) {
	shortURL := c.Query("short")
	longURL := c.Query("long")

	response, err := addURL(shortURL, longURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func retrieve(c *gin.Context) {
	shortURL := c.Query("short")

	response, err := getURL(shortURL)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
