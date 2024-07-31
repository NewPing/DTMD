package main

import (
	_ "DTMD_API/docs" // replace with your actual project path
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"crypto/rand"
	"math/big"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

var lobbys map[string]lobby

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

// createLobbyRequest represents the request body for creating a new lobby.
type createLobbyRequest struct {
	Name string `json:"name" binding:"required"`
}

func main() {
	lobbys = make(map[string]lobby)
	router := gin.Default()
	router.Use(cors.Default())
	// Swagger setup
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Routes
	router.GET("/members", getMembers)
	router.POST("/lobbies", createLobby)

	router.Run("localhost:8080")
}

// generateRandomID generates a random string of a given length using the specified character set
func generateRandomPin(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var pin string
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return ""
		}
		pin += string(charset[randomIndex.Int64()])
	}
	return pin
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
	id := generateRandomPin(6)
	for {
		if _, exists := lobbys[id]; exists {
			id = generateRandomPin(6)
			break
		}
		else {
			break
		}
	}
	// Create a new lobby
	var newLobby = lobby{
		ID:      id,
		Name:    req.Name,
		Members: []member{},
	}
	lobbys[id] = newLobby

	// Return the ID of the new lobby
	c.JSON(http.StatusOK, gin.H{"pin": newLobby.PIN})
}
