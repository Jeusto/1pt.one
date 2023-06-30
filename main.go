package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	router.POST("/", uploadFile)
	router.GET("/status", getStatus)
	router.GET("/:shortURL", getShortUrl)
	router.GET("/shorten", shortenUrl)
	router.GET("/retrieve", getStats)

	// Start the server
	router.Run(":8080")
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func getStatus(c *gin.Context) {
	response := gin.H{
		"status":  200,
		"message": "API is live. Read the documentation at https://github.com/Jeusto/1pt.one",
	}

	c.JSON(http.StatusOK, response)
}

func getShortUrl(c *gin.Context) {
	shortURL := c.Param("shortURL")

	item, err := getURL(shortURL)
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusFound, "/?url_not_found="+shortURL)
		return
	}

	incrementUrlVisits(shortURL)
	if item.LongURL != "" {
		c.Redirect(http.StatusFound, item.LongURL)
	} else {
		c.String(http.StatusOK, item.FileContent)
	}
}

func shortenUrl(c *gin.Context) {
	shortURL := c.Query("short")
	longURL := c.Query("long")

	response, err := addURL(shortURL, longURL)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func getStats(c *gin.Context) {
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

func uploadFile(c *gin.Context) {
	codeBytes, _ := ioutil.ReadAll(c.Request.Body)
	file := string(codeBytes)

	shortURL, err := addFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.String(http.StatusOK, "https://1pt.one/"+shortURL.ShortURL+"\n")
}
