package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"api-golang/database"
)

func init() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		content, err := ioutil.ReadFile(os.Getenv("DATABASE_URL_FILE"))
		if err != nil {
			log.Fatal(err)
		}
		databaseUrl = string(content)
	}

	errDB := database.InitDB(databaseUrl)
	if errDB != nil {
		log.Fatalf("⛔ Unable to connect to database: %v\n", errDB)
	} else {
		log.Println("DATABASE CONNECTED 🥇")
	}

}

func main() {

	r := gin.Default()
	var tm time.Time
	var reqCount int

	r.GET("/", func(c *gin.Context) {
		database.InsertView(c)
		tm, reqCount = database.GetTimeAndRequestCount(c)
		c.JSON(200, gin.H{
			"api":          "go",
			"currentTime":  tm,
			"requestCount": reqCount,
		})
	})

	r.GET("/health", func(c *gin.Context) {
		dbStatus := "UP"

		// Assuming your database package exposes the raw sql.DB or a Ping method
		// If the DB is down, we change the status and can optionally return a 503 Service Unavailable
		if err := database.Ping(c); err != nil {
			dbStatus = "DOWN"
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "DOWN",
				"database": dbStatus,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "UP",
			"database": dbStatus,
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		_, _ = database.GetTimeAndRequestCount(c)
		c.JSON(200, "pong")
	})

	port := os.Getenv("PORT")
	if port == "" {
		// Defaulting to 8000 to deconflict with unprivileged nginx container
		port = "8000"
	}

	r.Run(":" + port) // listen and serve on 0.0.0.0:8000 (or "PORT" env var if set)}
}
