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

//Creating a struct for the timeseries queries

type VesselTripSeries struct {
	Sum int    `json: "sum"`
	Day string `json: "day"`
}

//Creating a struct for vessel distribution queries

type VesselTypeDistribution struct {
	Count  float64 `json: "total"`
	Vessel string  `json: "vessel_description"`
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

func getTimeSeries(c *gin.Context) { //times series of vessel trip counts
	//Open DB connection
	db := OpenConnection()

	rows, err := db.Query("SELECT day, sum(counts) FROM results GROUP BY(day, s_id) HAVING s_id = " + c.Param("id") + "") // S_ID NEEDS TO BE SET BASED ON WHAT STRUCTURE IS CHOSEN
	if err != nil {
		panic(err)
	}

	var timeSeries []VesselTripSeries

	for rows.Next() {
		var result VesselTripSeries
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

func getUnqVessels(c *gin.Context) { //distribution of vessel type by unique vessels broken out by structure
	//Open DB connection
	db := OpenConnection()

	rows, err := db.Query("WITH structure AS ( " +
		"SELECT * FROM results WHERE s_id = " + c.Param("id") + "), " + //THE STRUCTURE ID HAS TO BE CHANGED TO A VARIABLE PASSED IN FROM THE FRONT END BASED ON WHICH STRUCTURE IS BEING QUERIED
		"t AS ( " +
		"SELECT shiptype1 AS vessel_type, sum(shiptype1_unqnum) AS t_sum FROM structure " +
		"GROUP BY (shiptype1) " +
		"UNION ALL " +
		"SELECT shiptype2, sum(shiptype2_unqnum) FROM structure " +
		"GROUP BY (shiptype2) " +
		"UNION ALL " +
		"SELECT shiptype3, sum(shiptype3_unqnum) FROM structure " +
		"GROUP BY (shiptype3) " +
		"UNION ALL " +
		"SELECT shiptype4, sum(shiptype4_unqnum) FROM structure " +
		"GROUP BY (shiptype4) " +
		"UNION ALL " +
		"SELECT shiptype5, sum(shiptype5_unqnum) FROM structure " +
		"GROUP BY (shiptype5) " +
		"UNION ALL " +
		"SELECT shiptype6, sum(shiptype6_unqnum) FROM structure " +
		"GROUP BY (shiptype6) " +
		"UNION ALL " +
		"SELECT shiptype7, sum(shiptype7_unqnum) FROM structure " +
		"GROUP BY (shiptype7) " +
		"UNION ALL " +
		"SELECT shiptype8, sum(shiptype8_unqnum) FROM structure " +
		"GROUP BY (shiptype8) " +
		"UNION ALL " +
		"SELECT shiptype9, sum(shiptype9_unqnum) FROM structure " +
		"GROUP BY (shiptype9) " +
		"UNION ALL " +
		"SELECT shiptype10, sum(shiptype10_unqnum) FROM structure " +
		"GROUP BY (shiptype10) " +
		"UNION ALL " +
		"SELECT shiptype11, sum(shiptype11_unqnum) FROM structure " +
		"GROUP BY (shiptype11) " +
		"UNION ALL " +
		"SELECT shiptype12, sum(shiptype12_unqnum) FROM structure " +
		"GROUP BY (shiptype12) " +
		"UNION ALL " +
		"SELECT shiptype13, sum(shiptype13_unqnum) FROM structure " +
		"GROUP BY (shiptype13) " +
		"UNION ALL " +
		"SELECT shiptype14, sum(shiptype14_unqnum) FROM structure " +
		"GROUP BY (shiptype14) " +
		"UNION ALL " +
		"SELECT shiptype15, sum(shiptype15_unqnum) FROM structure " +
		"GROUP BY (shiptype15) " +
		"UNION ALL " +
		"SELECT shiptype16, sum(shiptype16_unqnum) FROM structure " +
		"GROUP BY (shiptype16) " +
		"UNION ALL " +
		"SELECT shiptype17, sum(shiptype17_unqnum) FROM structure " +
		"GROUP BY (shiptype17) " +
		"UNION ALL " +
		"SELECT shiptype18, sum(shiptype18_unqnum) FROM structure " +
		"GROUP BY (shiptype18) " +
		"UNION ALL " +
		"SELECT shiptype19, sum(shiptype19_unqnum) FROM structure " +
		"GROUP BY (shiptype19) " +
		"UNION ALL " +
		"SELECT shiptype20, sum(shiptype20_unqnum) FROM structure " +
		"GROUP BY (shiptype20) " +
		"UNION ALL " +
		"SELECT shiptype21, sum(shiptype21_unqnum) FROM structure " +
		"GROUP BY (shiptype21) " +
		"UNION ALL " +
		"SELECT shiptype22, sum(shiptype22_unqnum) FROM structure " +
		"GROUP BY (shiptype22) ), " +
		"t2 AS ( " +
		"SELECT vessel_type, sum(t_sum) AS total FROM t GROUP BY(vessel_type)), " +
		"vessel_sum AS ( " +
		"SELECT vessel_type, total FROM t2 WHERE vessel_type is not null AND total is not null) " +
		"SELECT total, vessel_description FROM vessel_sum, vessels WHERE vessel_sum.vessel_type = vessels.vessel_id")

	if err != nil {
		panic(err)
	}

	var uniqueVessels []VesselTypeDistribution

	for rows.Next() {
		var result VesselTypeDistribution
		rows.Scan(&result.Count, &result.Vessel)
		uniqueVessels = append(uniqueVessels, result)
	}

	c.Header("Content-Type", "application/json")

	//Next Two Lines are needed for cors and cross origin requests
	c.Header("Access-Control-Allow-Origin", "http://localhost:3000") //<== USE THIS LINE FOR DEVELOPMENT ON LOCAL MACHINE
	//c.Header("Access-Control-Allow-Origin", "https://strange-tome-305601.ue.r.appspot.com/") //<== USE THIS LINE FOR PRODUCTION
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	//Set reponse type to JSON
	c.JSON(http.StatusOK, &uniqueVessels)

}

func getVesselTripCount(c *gin.Context) { //gets the trip count of each vessel type

	//Open DB connection
	db := OpenConnection()

	rows, err := db.Query("WITH structure AS ( " +
		"SELECT shiptype1, shiptype1_num, shiptype2, shiptype2_num, shiptype3, shiptype3_num, shiptype4, shiptype4_num, shiptype5, shiptype5_num, " +
		"shiptype6, shiptype6_num, shiptype7, shiptype7_num, shiptype8, shiptype8_num, shiptype9, shiptype9_num, shiptype10, shiptype10_num, shiptype11, " +
		"shiptype11_num, shiptype12, shiptype12_num, shiptype13, shiptype13_num, shiptype14, shiptype14_num, shiptype15, shiptype15_num, shiptype16, shiptype16_num, " +
		"shiptype17, shiptype17_num, shiptype18, shiptype18_num, shiptype19, shiptype19_num, shiptype20, shiptype20_num, shiptype21, shiptype21_num, shiptype22, " +
		"shiptype22_num, s_id FROM results WHERE s_id = " + c.Param("id") + "), " +
		"t AS ( " +
		"SELECT shiptype1 AS vessel_type, sum(shiptype1_num) AS t_sum FROM structure " +
		"GROUP BY (shiptype1) " +
		"UNION ALL " +
		"SELECT shiptype2, sum(shiptype2_num) FROM structure " +
		"GROUP BY (shiptype2) " +
		"UNION ALL " +
		"SELECT shiptype3, sum(shiptype3_num) FROM structure " +
		"GROUP BY (shiptype3) " +
		"UNION ALL " +
		"SELECT shiptype4, sum(shiptype4_num) FROM structure " +
		"GROUP BY (shiptype4)" +
		"UNION ALL " +
		"SELECT shiptype5, sum(shiptype5_num) FROM structure " +
		"GROUP BY (shiptype5) " +
		"UNION ALL " +
		"SELECT shiptype6, sum(shiptype6_num) FROM structure " +
		"GROUP BY (shiptype6) " +
		"UNION ALL " +
		"SELECT shiptype7, sum(shiptype7_num) FROM structure " +
		"GROUP BY (shiptype7) " +
		"UNION ALL " +
		"SELECT shiptype8, sum(shiptype8_num) FROM structure " +
		"GROUP BY (shiptype8) " +
		"UNION ALL " +
		"SELECT shiptype9, sum(shiptype9_num) FROM structure " +
		"GROUP BY (shiptype9) " +
		"UNION ALL " +
		"SELECT shiptype10, sum(shiptype10_num) FROM structure " +
		"GROUP BY (shiptype10) " +
		"UNION ALL " +
		"SELECT shiptype11, sum(shiptype11_num) FROM structure " +
		"GROUP BY (shiptype11) " +
		"UNION ALL " +
		"SELECT shiptype12, sum(shiptype12_num) FROM structure " +
		"GROUP BY (shiptype12) " +
		"UNION ALL " +
		"SELECT shiptype13, sum(shiptype13_num) FROM structure " +
		"GROUP BY (shiptype13) " +
		"UNION ALL " +
		"SELECT shiptype14, sum(shiptype14_num) FROM structure " +
		"GROUP BY (shiptype14) " +
		"UNION ALL " +
		"SELECT shiptype15, sum(shiptype15_num) FROM structure " +
		"GROUP BY (shiptype15) " +
		"UNION ALL " +
		"SELECT shiptype16, sum(shiptype16_num) FROM structure " +
		"GROUP BY (shiptype16) " +
		"UNION ALL " +
		"SELECT shiptype17, sum(shiptype17_num) FROM structure " +
		"GROUP BY (shiptype17) " +
		"UNION ALL " +
		"SELECT shiptype18, sum(shiptype18_num) FROM structure " +
		"GROUP BY (shiptype18) " +
		"UNION ALL " +
		"SELECT shiptype19, sum(shiptype19_num) FROM structure " +
		"GROUP BY (shiptype19) " +
		"UNION ALL " +
		"SELECT shiptype20, sum(shiptype20_num) FROM structure " +
		"GROUP BY (shiptype20) " +
		"UNION ALL " +
		"SELECT shiptype21, sum(shiptype21_num) FROM structure " +
		"GROUP BY (shiptype21) " +
		"UNION ALL " +
		"SELECT shiptype22, sum(shiptype22_num) FROM structure " +
		"GROUP BY (shiptype22) ), " +
		"t2 AS ( " +
		"SELECT vessel_type, sum(t_sum) AS total FROM t GROUP BY(vessel_type)), " +
		"vessel_sum AS ( " +
		"SELECT vessel_type, total FROM t2 WHERE vessel_type is not null AND total is not null) " +
		"SELECT total, vessel_description FROM vessel_sum, vessels WHERE vessel_sum.vessel_type = vessels.vessel_id")
	if err != nil {
		panic(err)
	}

	var vesselCounts []VesselTypeDistribution

	for rows.Next() {
		var result VesselTypeDistribution
		rows.Scan(&result.Count, &result.Vessel)
		vesselCounts = append(vesselCounts, result)
	}

	//Set reponse type to JSON
	c.Header("Content-Type", "application/json")

	//Next Two Lines are needed for cors and cross origin requests
	c.Header("Access-Control-Allow-Origin", "http://localhost:3000") //<== USE THIS LINE FOR DEVELOPMENT ON LOCAL MACHINE
	//c.Header("Access-Control-Allow-Origin", "https://strange-tome-305601.ue.r.appspot.com/") //<== USE THIS LINE FOR PRODUCTION
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	c.JSON(http.StatusOK, &vesselCounts)
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
		api.GET("/timeseries/:id", getTimeSeries)
		api.GET("/uniquevessels/:id", getUnqVessels)
		api.GET("/vesseltripcounts/:id", getVesselTripCount)

	}

	// Start and run the server
	log.Printf("Listening on port %s", port)
	router.Run(":" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
