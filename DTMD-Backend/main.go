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
	ID                 string `json:"id"`
	Name               string `json:"name"`
	UpdateInstructions []int  `json:"int"`
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

// createLobbyRequest represents the request body for creating a new lobby.
type joinLobbyRequest struct {
	Nickname string `json:"nickname" binding:"required"`
}

const (
	InstructionUpdateLobbyMembers = 0
	InstructionUpdateChat         = 1
)

func main() {
	lobbys = make(map[string]lobby)
	router := gin.Default()
	router.Use(cors.Default())
	// Swagger setup
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Routes

	router.POST("/lobbies", createLobby)
	router.POST("/lobbies/:id/members", joinLobby)
	router.GET("/lobbies/:id/members", getLobbyMembers)
	router.GET("/lobbies/:id/members/:id2/updates", getUpdateInstructions)

	router.Run("localhost:8080")
}

// CreateLobby godoc
// @Summary      Create a new lobby
// @Description  create a new lobby with the given name
// @Tags         lobbies
// @Accept       json
// @Produce      plain
// @Param        lobby body createLobbyRequest true "Create Lobby"
// @Success      200 {string} string
// @Failure      400
// @Failure      500
// @Router /lobbies [post]
func createLobby(c *gin.Context) {
	var req createLobbyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := generateUniqueLobbyID()

	// Create a new lobby
	var newLobby = lobby{
		ID:      id,
		Name:    req.Name,
		Members: []member{},
	}
	lobbys[id] = newLobby

	// Return the ID of the new lobby
	c.String(http.StatusOK, newLobby.ID)
}

// JoinLobby godoc
// @Summary      Join an existing Lobby
// @Description  lets a user join a lobby
// @Tags         lobbies
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Lobby ID"
// @Param        lobby body joinLobbyRequest true "Join Lobby"
// @Success      200  {object} string
// @Failure      400
// @Failure      500
// @Router /lobbies/{id}/members [post]
func joinLobby(c *gin.Context) {
	id := c.Param("id")
	var req joinLobbyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newMember = member{
		ID:   generateUniqueMemberID(),
		Name: req.Nickname,
	}

	if lobby, exists := lobbys[id]; exists {
		notifyLobbyMembers(id, InstructionUpdateLobbyMembers)
		lobby.Members = append(lobby.Members, newMember)
		lobbys[id] = lobby

		c.JSON(http.StatusOK, gin.H{"id": newMember.ID})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"id": ""})
	}
}

// GetLobbyMembers godoc
// @Summary      Get members of a lobby
// @Description  get members of a specific lobby by ID
// @Tags         lobbies
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Lobby ID"
// @Success      200  {array}   []string
// @Failure      400
// @Failure      404
// @Router /lobbies/{id}/members [get]
func getLobbyMembers(c *gin.Context) {
	id := c.Param("id")

	if lobby, exists := lobbys[id]; exists {
		var membersNames []string

		for _, m := range lobby.Members {
			membersNames = append(membersNames, m.Name)
		}

		c.JSON(http.StatusOK, gin.H{"id": membersNames})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"id": ""})
	}
}

// GetMemberUpdateInstructions godoc
// @Summary      get member update instructions
// @Description  Get update instructions of a specific member
// @Tags         lobbies
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Lobby ID"
// @Param        id2   path      string  true  "Member ID"
// @Success      200  {array}   []int
// @Failure      400
// @Failure      404
// @Router /lobbies/{id}/members/{id2}/updates [get]
func getUpdateInstructions(c *gin.Context) {
	lobbyID := c.Param("id")
	memberID := c.Param("id2")

	var lobbys2 = lobbys
	print(lobbys2)

	if lobby, exists := lobbys[lobbyID]; exists {
		for i := range lobby.Members {
			if lobby.Members[i].ID == memberID {
				c.JSON(http.StatusOK, gin.H{"id": lobby.Members[i].UpdateInstructions})
				lobby.Members[i].UpdateInstructions = []int{}
				lobbys[lobbyID] = lobby
				return
			}
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"id": ""})
}

func generateUniqueLobbyID() string {
	id := generateRandomPin(6)
	for {
		if _, exists := lobbys[id]; exists {
			id = generateRandomPin(6)
		} else {
			break
		}
	}
	return id
}

func generateUniqueMemberID() string {
	id := generateRandomPin(12)
	for {
		if _, exists := lobbys[id]; exists {
			id = generateRandomPin(6)
		} else {
			break
		}
	}
	return id
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

func notifyLobbyMembers(lobbyID string, updateInstructionType int) {
	if lobby, exists := lobbys[lobbyID]; exists {

		for i := range lobby.Members {
			lobby.Members[i].UpdateInstructions = append(lobby.Members[i].UpdateInstructions, updateInstructionType)
		}
	}
}
