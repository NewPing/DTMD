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

// createLobbyRequest represents the request body for creating a new lobby.
type createLobbyRequest struct {
	Name string `json:"name" binding:"required"`
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	// Swagger setup
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Routes
	router.GET("/members", getMembers)
	router.POST("/lobbies", createLobby)

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

// CreateLobby godoc
// @Summary      Create a new lobby
// @Description  create a new lobby with the given name
// @Tags         lobbies
// @Accept       json
// @Produce      json
// @Param        lobby body createLobbyRequest true "Create Lobby"
// @Success      200  {object} int
// @Failure      400
// @Failure      500
// @Router /lobbies [post]
func createLobby(c *gin.Context) {
	var req createLobbyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new lobby
	var newLobby = lobby{
		ID:      generateID(),
		Name:    req.Name,
		Members: []member{},
	}
	lobbys = append(lobbys, newLobby)

	// Return the ID of the new lobby
	c.JSON(http.StatusOK, gin.H{"id": newLobby.ID})
}

// generateID generates a new unique ID for a lobby
func generateID() string {
	return "12345"
}
