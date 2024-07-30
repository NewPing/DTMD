package main

import (
	_ "DTMD_API/docs" // replace with your actual project path
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// album represents data about a record album.
type member struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type lobby struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Members []member `json:"list"`
}

// albums slice to seed record album data.
var lobbys = []lobby{
	{ID: "1", Name: "AromaticA", Members: []member{{ID: "1", Name: "AromaticA"}}},
	{ID: "2", Name: "HeckLeggedJoe", Members: []member{{ID: "1", Name: "AromaticA"}}},
	{ID: "3", Name: "OnlyNew", Members: []member{{ID: "1", Name: "AromaticA"}}},
	{ID: "4", Name: "Britney", Members: []member{{ID: "1", Name: "AromaticA"}}},
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	// Swagger setup
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Routes
	router.GET("/members", getMembers)

	router.Run("localhost:8080")
}

// ListAccounts godoc
// @Summary      List of members
// @Description  get a list of members
// @Tags         members
// @Accept       json
// @Produce      json
// @Success      200  {array}   main.member
// @Failure      400
// @Failure      404
// @Failure      500
// @Router /members [get]
func getMembers(c *gin.Context) {
	c.JSON(http.StatusOK, lobbys[0].Members)
}
