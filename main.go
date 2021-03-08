package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	_ "github.com/lib/pq"
)

//Setting a constant for the database information
const (
	instanceConnection = "35.196.23.166" //"strange-tome-305601:us-east1:vessel-data"
	databaseName       = "postgres"
	user               = "postgres"
	password           = "capstone"
)

//Creating a struct for the structure table
type Structure struct {
	ID       int    `json: "structure_id"`
	Location string `json: "location"`
	Lon      string `json: "geog"`
	Lat      string `json: "geog"`
	Year     int    `json: year_constructed"`
	Type     int    `json: structure_type"`
}

//Creating a struct for the vessel table
type Vessel struct {
	ID          string `json: "vessel_id"`
	Description string `json: "vessel_description"`
}

//Open database connection and return a reference to the database
func OpenConnection() *sql.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", instanceConnection, user, password, databaseName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func getVessels(c *gin.Context) {
	//Open DB connection
	db := OpenConnection()

	//Grab everything from vessels table
	rows, err := db.Query("SELECT * vessels")
	if err != nil {
		log.Fatal(err)
	}

	//Create instance of Vessel struct
	var vessels []Vessel

	//For the results of the query add information to the instance of Vessel created above
	for rows.Next() {
		var vessel Vessel
		rows.Scan(&vessel.ID, &vessel.Description)
		vessels = append(vessels, vessel)
	}

	//Set response type to JSON
	c.Header("Content-Type", "application/json")

	//Next Two Lines are needed for cors and cross origin requests

	//c.Header("Access-Control-Allow-Origin", "http://localhost:3000") //<== USE THIS LINE FOR DEVELOPMENT ON LOCAL MACHINE
	c.Header("Access-Control-Allow-Origin", "https://strange-tome-305601.ue.r.appspot.com/") //<== USE THIS LINE FOR PRODUCTION
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	//Return data from query
	c.JSON(http.StatusOK, &vessels)
}

func getStructures(c *gin.Context) {
	//Open DB connection
	db := OpenConnection()

	rows, err := db.Query("SELECT  structure_id, location, ST_X(geog::geometry), ST_Y(geog::geometry), year_constructed, structure_type FROM structures")
	if err != nil {
		panic(err)
	}

	//Create instance of Structure struct
	var structures []Structure

	//For the results of the query add information to the instance of Structure created above
	for rows.Next() {
		var structure Structure
		rows.Scan(&structure.ID, &structure.Location, &structure.Lon, &structure.Lat, &structure.Year, &structure.Type)
		structures = append(structures, structure)
	}

	//Set reponse type to JSON
	c.Header("Content-Type", "application/json")

	//Next Two Lines are needed for cors and cross origin requests
	//c.Header("Access-Control-Allow-Origin", "http://localhost:3000") //<== USE THIS LINE FOR DEVELOPMENT ON LOCAL MACHINE
	c.Header("Access-Control-Allow-Origin", "https://strange-tome-305601.ue.r.appspot.com/") //<== USE THIS LINE FOR PRODUCTION
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	//Return data from query
	c.JSON(http.StatusOK, &structures)
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
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

		//Setting Api routes
		api.GET("/vessel", getVessels)
		api.GET("/structure", getStructures)

	}

	// Start and run the server
	log.Printf("Listening on port %s", port)
	router.Run(":" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
