package main

import (
	_ "DTMD_API/docs" // replace with your actual project path
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// album represents data about a record album.
type member struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// albums slice to seed record album data.
var members = []member{
	{ID: "1", Name: "AromaticA"},
	{ID: "2", Name: "HeckLeggedJoe"},
	{ID: "3", Name: "OnlyNew"},
	{ID: "4", Name: "Britney"},
}

func main() {
	router := gin.Default()
	// Swagger setup
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Routes
	router.GET("/members", getMembers)

	router.Run("localhost:8080")
}

// @Router /members [get]
func getMembers(c *gin.Context) {
	c.JSON(http.StatusOK, members)
}
