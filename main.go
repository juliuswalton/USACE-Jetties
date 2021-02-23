package main

import (
	"log"
	"net/http"
	"os"
  
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)
/*func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
	}
	fmt.Fprint(w, "Hello, World!")
}*/

func main() {
  port := os.Getenv("PORT")
	if port == "" {
			port = "8080"
			log.Printf("Defaulting to port %s", port)
	}


  // Set the router as the default one shipped with Gin
  router := gin.Default()
  
  // Serve frontend static files
  router.Use(static.Serve("/", static.LocalFile("./views", true)))
  
  // Setup route group for the API
  api := router.Group("/api")
  {
    api.GET("/", func(c *gin.Context) {
      c.JSON(http.StatusOK, gin.H {
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