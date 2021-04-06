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
	ID        int    `json: "structure_id"`
	Name      string `json: "name"`
	Lon       string `json: "longitude"`
	Lat       string `json: "latitude"`
	Year      int    `json: "year_constructed"`
	Type      string `json: "type_description"`
	Length    int    `json: "structure_length"`
	Community int    `json: "community"`
	Count     int    `json: "count"`
}

//Creating a struct for the vessel table
type Vessel struct {
	ID          string `json: "vessel_id"`
	Description string `json: "vessel_description"`
}

//Creating a struct for the results table ** this is not its final form, don't be afraid

type Result struct {
	Sum 						int `json: "sum"`
	ID 							int `json: "s_id"`
	Day							string `json: "day"`
	Distance				float64 `json: "avg_dist"`
	Counts					int 		`json: "counts"`
	UnqCount 				int `json: "num_unique"`
	Shiptype1 			int `json: "shiptype1"`
	Shiptype1Count  int `json: "shiptype1_num"`
	Shiptype1Unq    int `json: "shiptype1_unqnum"`
	Shiptype2 			int `json: "shiptype2"`
	Shiptype2Count  int `json: "shiptype2_num"`
	Shiptype2Unq    int `json: "shiptype2_unqnum"`
	Shiptype3 			int `json: "shiptype3"`
	Shiptype3Count  int `json: "shiptype3_num"`
	Shiptype3Unq    int `json: "shiptype3_unqnum"`
	Shiptype4 			int `json: "shiptype4"`
	Shiptype4Count  int `json: "shiptype4_num"`
	Shiptype4Unq    int `json: "shiptype4_unqnum"`
	Shiptype5 			int `json: "shiptype5"`
	Shiptype5Count  int `json: "shiptype5_num"`
	Shiptype5Unq    int `json: "shiptype5_unqnum"`
	Shiptype6 			int `json: "shiptype6"`
	Shiptype6Count  int `json: "shiptype6_num"`
	Shiptype6Unq    int `json: "shiptype6_unqnum"`
	Shiptype7 			int `json: "shiptype7"`
	Shiptype7Count  int `json: "shiptype7_num"`
	Shiptype7Unq    int `json: "shiptype7_unqnum"`
	Shiptype8 			int `json: "shiptype8"`
	Shiptype8Count  int `json: "shiptype8_num"`
	Shiptype8Unq    int `json: "shiptype8_unqnum"`
	Shiptype9 			int `json: "shiptype9"`
	Shiptype9Count  int `json: "shiptype9_num"`
	Shiptype9Unq    int `json: "shiptype9_unqnum"`
	Shiptype10 			int `json: "shiptype10"`
	Shiptype10Count int `json: "shiptype10_num"`
	Shiptype10Unq   int `json: "shiptype10_unqnum"`
	Shiptype11 			int `json: "shiptype11"`
	Shiptype11Count int `json: "shiptype11_num"`
	Shiptype11Unq   int `json: "shiptype11_unqnum"`
	Shiptype12 			int `json: "shiptype12"`
	Shiptype12Count int `json: "shiptype12_num"`
	Shiptype12Unq   int `json: "shiptype12_unqnum"`
	Shiptype13 			int `json: "shiptype13"`
	Shiptype13Count int `json: "shiptype13_num"`
	Shiptype13Unq   int `json: "shiptype13_unqnum"`
	Shiptype14 			int `json: "shiptype14"`
	Shiptype14Count int `json: "shiptype14_num"`
	Shiptype14Unq   int `json: "shiptype14_unqnum"`
	Shiptype15 			int `json: "shiptype15"`
	Shiptype15Count int `json: "shiptype15_num"`
	Shiptype15Unq   int `json: "shiptype15_unqnum"`
	Shiptype16 			int `json: "shiptype16"`
	Shiptype16Count int `json: "shiptype16_num"`
	Shiptype16Unq   int `json: "shiptype16_unqnum"`
	Shiptype17 			int `json: "shiptype17"`
	Shiptype17Count int `json: "shiptype17_num"`
	Shiptype17Unq   int `json: "shiptype17_unqnum"`
	Shiptype18 			int `json: "shiptype18"`
	Shiptype18Count int `json: "shiptype18_num"`
	Shiptype18Unq   int `json: "shiptype18_unqnum"`
	Shiptype19 			int `json: "shiptype19"`
	Shiptype19Count int `json: "shiptype19_num"`
	Shiptype19Unq   int `json: "shiptype19_unqnum"`
	Shiptype20 			int `json: "shiptype20"`
	Shiptype20Count int `json: "shiptype20_num"`
	Shiptype20Unq   int `json: "shiptype20_unqnum"`
	Shiptype21 			int `json: "shiptype21"`
	Shiptype21Count int `json: "shiptype21_num"`
	Shiptype21Unq   int `json: "shiptype21_unqnum"`
	Shiptype22 			int `json: "shiptype22"`
	Shiptype22Count int `json: "shiptype22_num"`
	Shiptype22Unq   int `json: "shiptype22_unqnum"`




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

	c.Header("Access-Control-Allow-Origin", "http://localhost:3000") //<== USE THIS LINE FOR DEVELOPMENT ON LOCAL MACHINE
	//c.Header("Access-Control-Allow-Origin", "https://strange-tome-305601.ue.r.appspot.com/") //<== USE THIS LINE FOR PRODUCTION
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	//Return data from query
	c.JSON(http.StatusOK, &vessels)
}

func getStructures(c *gin.Context) {
	//Open DB connection
	db := OpenConnection()

	rows, err := db.Query("SELECT s.structure_id, s.name, s.longitude, s.latitude, " +
		"s.year_constructed, st.type_description, s.structure_length, s.community, SUM (r.counts) count " +
		"FROM structures s " +
		"JOIN structure_types st ON s.structure_type = st.type_id " +
		"JOIN results r ON s.structure_id = r.s_id " +
		"GROUP BY s.structure_id, st.type_description")
	if err != nil {
		panic(err)
	}

	//Create instance of Structure struct
	var structures []Structure

	//For the results of the query add information to the instance of Structure created above
	for rows.Next() {
		var structure Structure
		rows.Scan(&structure.ID, &structure.Name, &structure.Lon, &structure.Lat, &structure.Year, &structure.Type, &structure.Length, &structure.Community, &structure.Count)
		structures = append(structures, structure)
	}

	//Set reponse type to JSON
	c.Header("Content-Type", "application/json")

	//Next Two Lines are needed for cors and cross origin requests
	c.Header("Access-Control-Allow-Origin", "http://localhost:3000") //<== USE THIS LINE FOR DEVELOPMENT ON LOCAL MACHINE
	//c.Header("Access-Control-Allow-Origin", "https://strange-tome-305601.ue.r.appspot.com/") //<== USE THIS LINE FOR PRODUCTION
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	//Return data from query
	c.JSON(http.StatusOK, &structures)
}

func getTimeSeries(c *gin.Context){ //times series of vessel trip counts
	//Open DB connection
	db := OpenConnection()

	rows, err := db.Query("SELECT day, sum(counts) FROM results GROUP BY day")
	if err != nil {
		panic(err)
	}

	var timeSeries []Result

	for rows.Next() {
		var result Result
		rows.Scan(&result.Day, &result.Sum)
		timeSeries = append(timeSeries, result)
	}

	//Set reponse type to JSON
	c.Header("Content-Type", "application/json")

	//Next Two Lines are needed for cors and cross origin requests
	c.Header("Access-Control-Allow-Origin", "http://localhost:3000") //<== USE THIS LINE FOR DEVELOPMENT ON LOCAL MACHINE
	//c.Header("Access-Control-Allow-Origin", "https://strange-tome-305601.ue.r.appspot.com/") //<== USE THIS LINE FOR PRODUCTION
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	c.JSON(http.StatusOK, &timeSeries)

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
		api.GET("/timeseries", getTimeSeries)

	}

	// Start and run the server
	log.Printf("Listening on port %s", port)
	router.Run(":" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
