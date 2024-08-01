package main

import (
	_ "DTMD_API/docs" // replace with your actual project path
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"crypto/rand"
	"math/big"
	rand2 "math/rand/v2"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

var lobbys map[string]lobby

// album represents data about a record album.
type member struct {
	ID                 string        `json:"id"`
	Name               string        `json:"name"`
	UpdateInstructions []int         `json:"updateInstructions"`
	NewChatMessages    []ChatMessage `json:"newChatMessages"`
}

type ChatMessage struct {
	Sender  string `json:"id"`
	Message string `json:"name"`
}

type lobby struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Members []member `json:"members"`
}

// createLobbyRequest represents the request body for creating a new lobby.
type createLobbyRequest struct {
	Name string `json:"name" binding:"required"`
}

// createLobbyRequest represents the request body for creating a new lobby.
type joinLobbyRequest struct {
	Nickname string `json:"nickname" binding:"required"`
}

// createLobbyRequest represents the request body for creating a new lobby.
type rollDiceRequest struct {
	MemberID      string `json:"MemberID" binding:"required"`
	IsPrivateRoll *int   `json:"IsPrivateRoll" binding:"required"`
	NumberOfRolls *int   `json:"NumberOfRolls" binding:"required"`
	DiceType      *int   `json:"DiceType" binding:"required"`
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
	router.POST("/lobbies/:id/rolldice", rollDice)
	router.GET("/lobbies/:id/name", getLobbyName)
	router.GET("/lobbies/:id/members", getLobbyMembers)
	router.GET("/lobbies/:id/members/:id2/updates", getUpdateInstructions)
	router.GET("/lobbies/:id/members/:id2/messages", getNewChatMessages)

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
// @Produce      plain
// @Param        id   path      string  true  "Lobby ID"
// @Param        lobby body joinLobbyRequest true "Join Lobby"
// @Success      200  {string} string
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
		ID:                 generateUniqueMemberID(),
		Name:               req.Nickname,
		UpdateInstructions: []int{},
	}

	if lobby, exists := lobbys[id]; exists {
		notifyLobbyMembers(id, InstructionUpdateLobbyMembers)
		notifyLobbyMembers(id, InstructionUpdateChat)
		lobby.Members = append(lobby.Members, newMember)
		lobbys[id] = lobby

		c.String(http.StatusOK, newMember.ID)
	} else {
		c.String(http.StatusBadRequest, newMember.ID)
	}
}

// RollDice godoc
// @Summary      Roll dice post request
// @Description  roll dice and send back integer result in string form
// @Tags         lobbies
// @Accept       json
// @Produce      plain
// @Param        id    path     string  true  "Lobby ID"
// @Param        lobby body     rollDiceRequest true "Roll Dice Request"
// @Success      200   {string}    string
// @Failure      400
// @Failure      500
// @Router /lobbies/{id}/rolldice [post]
func rollDice(c *gin.Context) {
	id := c.Param("id")
	var req rollDiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var msg = "Has rolled " + strconv.Itoa(*req.NumberOfRolls) + " w" + strconv.Itoa(*req.DiceType)
	var result = 0
	for i := 0; i < *req.NumberOfRolls; i++ {
		var number = GenerateRandomNumber(*req.DiceType)
		result = result + number
	}

	msg += ", result: " + strconv.Itoa(result)

	if lobby, exists := lobbys[id]; exists {
		for i := range lobby.Members {
			if *req.IsPrivateRoll == 0 || lobby.Members[i].ID == req.MemberID {
				lobby.Members[i].NewChatMessages = append(lobby.Members[i].NewChatMessages, ChatMessage{Sender: GetUserNameByID(id, req.MemberID), Message: msg})
				if *req.IsPrivateRoll == 0 {
					notifyLobbyMembers(id, InstructionUpdateChat)
				} else {
					notifyLobbyMember(id, req.MemberID, InstructionUpdateChat)
				}

				lobbys[id] = lobby
				c.JSON(http.StatusOK, result)
			}
		}
	} else {
		c.String(http.StatusBadRequest, strconv.Itoa(-1))
	}
}

// GetLobbyName godoc
// @Summary      get lobby name
// @Description  return the name of the specified lobby
// @Tags         lobbies
// @Accept       json
// @Produce      plain
// @Param        id   path      string  true  "Lobby ID"
// @Success      200  {string} string
// @Failure      400
// @Failure      500
// @Router /lobbies/{id}/name [get]
func getLobbyName(c *gin.Context) {
	id := c.Param("id")

	if lobby, exists := lobbys[id]; exists {
		c.String(http.StatusOK, lobby.Name)
	} else {
		c.String(http.StatusBadRequest, "")
	}
}

// GetLobbyMembers godoc
// @Summary      Get members of a lobby
// @Description  get members of a specific lobby by ID
// @Tags         member
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Lobby ID"
// @Success      200  {array}   string
// @Failure      400
// @Failure      404
// @Router /lobbies/{id}/members [get]
func getLobbyMembers(c *gin.Context) {
	id := c.Param("id")

	var membersNames []string
	if lobby, exists := lobbys[id]; exists {
		for _, m := range lobby.Members {
			membersNames = append(membersNames, m.Name)
		}

		c.JSON(http.StatusOK, membersNames)
	} else {
		c.JSON(http.StatusBadRequest, membersNames)
	}
}

// GetNewChatMessages godoc
// @Summary      get new messages
// @Description  get all new chat messages for this specific member
// @Tags         member
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Lobby ID"
// @Param        id2  path      string  true  "Member ID"
// @Success      200  {array}   ChatMessage
// @Failure      400
// @Failure      404
// @Router /lobbies/{id}/members/{id2}/messages [get]
func getNewChatMessages(c *gin.Context) {
	lobbyID := c.Param("id")
	memberID := c.Param("id2")

	if lobby, exists := lobbys[lobbyID]; exists {
		for i := range lobby.Members {
			if lobby.Members[i].ID == memberID {
				c.JSON(http.StatusOK, lobby.Members[i].NewChatMessages)
				lobby.Members[i].NewChatMessages = []ChatMessage{}
				lobbys[lobbyID] = lobby
				return
			}
		}
	}

	c.JSON(http.StatusBadRequest, []ChatMessage{})
}

// GetMemberUpdateInstructions godoc
// @Summary      get member update instructions
// @Description  Get update instructions of a specific member
// @Tags         member
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Lobby ID"
// @Param        id2   path      string  true  "Member ID"
// @Success      200  {array}   int
// @Failure      400
// @Failure      404
// @Router /lobbies/{id}/members/{id2}/updates [get]
func getUpdateInstructions(c *gin.Context) {
	lobbyID := c.Param("id")
	memberID := c.Param("id2")

	if lobby, exists := lobbys[lobbyID]; exists {
		for i := range lobby.Members {
			if lobby.Members[i].ID == memberID {
				c.JSON(http.StatusOK, lobby.Members[i].UpdateInstructions)
				lobby.Members[i].UpdateInstructions = []int{}
				lobbys[lobbyID] = lobby
				return
			}
		}
	}

	c.JSON(http.StatusBadRequest, []int{})
}

func GetUserNameByID(lobbyID, userID string) string {
	if lobby, exists := lobbys[lobbyID]; exists {
		for i := range lobby.Members {
			if lobby.Members[i].ID == userID {
				return lobby.Members[i].Name
			}
		}
	}
	return "undefined"
}

func GenerateRandomNumber(xmax int) int {
	return rand2.IntN(xmax) + 1

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
			containsInstructionType := slices.Contains(lobby.Members[i].UpdateInstructions, updateInstructionType)

			if !containsInstructionType {
				lobby.Members[i].UpdateInstructions = append(lobby.Members[i].UpdateInstructions, updateInstructionType)
			}
		}
	}
}

func notifyLobbyMember(lobbyID string, memberID string, updateInstructionType int) {
	if lobby, exists := lobbys[lobbyID]; exists {

		for i := range lobby.Members {
			if lobby.Members[i].ID == memberID {
				containsInstructionType := slices.Contains(lobby.Members[i].UpdateInstructions, updateInstructionType)

				if !containsInstructionType {
					lobby.Members[i].UpdateInstructions = append(lobby.Members[i].UpdateInstructions, updateInstructionType)
				}
			}
		}
	}
}
