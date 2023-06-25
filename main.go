package main

import (
	"log"
	"net/http"
	"os"
	"sensibull/stocks-api/routes"

	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) { c.JSON(http.StatusOK, "OK") }

func main() {
	router := gin.Default()
	router.GET("/health", health)

	routes.InitRoutes(router)

	err := router.Run(":" + os.Getenv("PORT"))

	if err != nil {
		log.Fatal(err)
	}
}
