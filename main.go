package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	//_ "github.com/lib/pq"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
)

const (
	instanceConnection = "strange-tome-305601:us-east1:vessel-data"
	databaseName       = "postgres"
	user               = "postgres"
	password           = "capstone"
)

func main() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", instanceConnection, user, password, databaseName)
	db, err := sql.Open("cloudsqlpostgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	rows, err := db.Query("SELECT * FROM vessels")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var (
		vessel_id          string
		vessel_description string
	)

	for rows.Next() {
		err := rows.Scan(&vessel_id, &vessel_description)
		if err != nil {
			panic(err)
		}
		log.Printf("%s: %s", vessel_id, vessel_description)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./build", true)))

	// Setup route group for the API
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	// Start and run the server

	log.Printf("Listening on port %s", port)
	router.Run(":" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
