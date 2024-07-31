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

	"github.com/google/uuid"
)

// album represents data about a record album.
type member struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type lobby struct {
	ID      uuid.UUID `json:"id"`
	PIN     string    `json:"pin"`
	Name    string    `json:"name"`
	Members []member  `json:"list"`
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
	router.POST("/lobbys/:lobbyName", createLobby)

	router.Run("localhost:8080")
}

// generateRandomID generates a random string of a given length using the specified character set
func generateRandomID(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var id string
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		id += string(charset[randomIndex.Int64()])
	}
	return id, nil
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

// creates a new lobby
// @Summary      Creates a new Lobby
// @Description  creates a new lobby
// @Tags         lobbys
// @Accept       json
// @Produce      json
// @Param lobbyName    header     string  true  "lobby name"
// @Success      200  {object} int
// @Failure      400
// @Failure      404
// @Failure      500
// @Router /lobbys [post]
func createLobby(c *gin.Context) {
	lobbyName := c.Param("lobbyName")
	unique_id := uuid.New()
	var newLobby = lobby{ID: "1", Name: lobbyName, Members: []member{{ID: "1", Name: "AromaticA"}}}
	lobbys = append(lobbys, newLobby)
	c.JSON(http.StatusOK, newLobby.Name)
}
