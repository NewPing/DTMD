package main

import (
	_ "DTMD_API/docs" // replace with your actual project path
	"DTMD_API/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"crypto/rand"
	"math/big"
	rand2 "math/rand/v2"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

var lobbyManager *models.LobbyManager

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
	lobbyManager = models.NewLobbyManager()
	router := gin.Default()
	router.Use(cors.Default())
	// Swagger setup
	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Routes

	router.POST("/api/lobbies", createLobby)
	router.POST("/api/lobbies/:id/members", joinLobby)
	router.POST("/api/lobbies/:id/rolldice", rollDice)
	router.GET("/api/lobbies/:id/name", getLobbyName)
	router.GET("/api/lobbies/:id/chathistory", getLobbyChatHistory)
	router.GET("/api/lobbies/:id/members", getLobbyMembers)
	router.GET("/api/lobbies/:id/members/:id2/updates", getUpdateInstructions)
	router.GET("/api/lobbies/:id/members/:id2/messages", getNewChatMessages)

	go startBackgroundWorker() //removes inactive clients and closes lobbys if no members are present

	router.Run("0.0.0.0:8080")
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
// @Router /api/lobbies [post]
func createLobby(c *gin.Context) {
	var req createLobbyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	id := generateUniqueLobbyID()

	newLobby := models.NewLobby(id, req.Name)
	lobbyManager.AddLobby(newLobby)

	c.String(http.StatusOK, newLobby.GetID())
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
// @router /api/lobbies/{id}/members [post]
func joinLobby(c *gin.Context) {
	id := c.Param("id")
	var req joinLobbyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	newMember := models.NewMember(generateUniqueMemberID(id), req.Nickname)

	lobby, exists := lobbyManager.GetLobby(id)

	if !exists {
		c.String(http.StatusNotFound, "")
		return
	}

	notifyLobbyMembers(id, InstructionUpdateLobbyMembers)
	lobby.AddMember(newMember)

	c.String(http.StatusOK, newMember.GetID())

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
// @router /api/lobbies/{id}/rolldice [post]
func rollDice(c *gin.Context) {
	id := c.Param("id")
	var req rollDiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	msg := "Has rolled " + strconv.Itoa(*req.NumberOfRolls) + " w" + strconv.Itoa(*req.DiceType)
	result := 0
	for i := 0; i < *req.NumberOfRolls; i++ {
		number := GenerateRandomNumber(*req.DiceType)
		result += number
	}

	msg += ", result: " + strconv.Itoa(result)

	lobby, exists := lobbyManager.GetLobby(id)
	if !exists {
		c.String(http.StatusNotFound, "")
		return
	}

	for _, member := range lobby.GetMembers() {
		if *req.IsPrivateRoll == 0 || member.GetID() == req.MemberID {
			chatMessage := models.ChatMessage{Sender: GetUserNameByID(id, req.MemberID), Message: msg, Timestamp: time.Now()}
			member.AddNewChatMessage(chatMessage)
			lobby.AddMessageToChatHistory(chatMessage)
			notifyLobbyMember(id, req.MemberID, InstructionUpdateChat)
		}
	}

	c.String(http.StatusOK, strconv.Itoa(result))
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
// @router /api/lobbies/{id}/name [get]
func getLobbyName(c *gin.Context) {
	id := c.Param("id")

	lobby, exists := lobbyManager.GetLobby(id)
	if !exists {
		c.String(http.StatusNotFound, "")
		return
	}
	c.String(http.StatusOK, lobby.GetName())
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
// @router /api/lobbies/{id}/members [get]
func getLobbyMembers(c *gin.Context) {
	id := c.Param("id")

	var membersNames []string

	lobby, exists := lobbyManager.GetLobby(id)
	if !exists {
		c.JSON(http.StatusNotFound, membersNames)
		return
	}

	for _, member := range lobby.GetMembers() {
		membersNames = append(membersNames, member.GetName())
	}
	c.JSON(http.StatusOK, membersNames)
}

// GetLobbyChatHistory godoc
// @Summary      Get the chat history of a lobby
// @Description  Get the history of all chat messages of a specific lobby
// @Tags         lobby
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Lobby ID"
// @Success      200  {array}   ChatMessage
// @Failure      400
// @Failure      404
// @router /api/lobbies/{id}/chathistory [get]
func getLobbyChatHistory(c *gin.Context) {
	id := c.Param("id")

	lobby, exists := lobbyManager.GetLobby(id)
	if !exists {
		c.JSON(http.StatusNotFound, []models.ChatMessage{})
		return
	}

	c.JSON(http.StatusOK, lobby.GetChatHistory())
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
// @router /api/lobbies/{id}/members/{id2}/messages [get]
func getNewChatMessages(c *gin.Context) {
	lobbyID := c.Param("id")
	memberID := c.Param("id2")

	lobby, exists := lobbyManager.GetLobby(lobbyID)
	if !exists {
		c.JSON(http.StatusNotFound, []models.ChatMessage{})
		return
	}

	for _, member := range lobby.GetMembers() {
		if member.GetID() == memberID {
			messages := member.GetNewChatMessages()
			member.ClearNewChatMessages()
			c.JSON(http.StatusOK, messages)
			return
		}
	}

	c.JSON(http.StatusBadRequest, []models.ChatMessage{})
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
// @router /api/lobbies/{id}/members/{id2}/updates [get]
func getUpdateInstructions(c *gin.Context) {
	lobbyID := c.Param("id")
	memberID := c.Param("id2")

	lobby, exists := lobbyManager.GetLobby(lobbyID)
	if !exists {
		c.JSON(http.StatusNotFound, []int{})
		return
	}

	for _, member := range lobby.GetMembers() {
		if member.GetID() == memberID {
			instructions := member.GetUpdateInstructions()
			member.ClearUpdateInstructions()
			member.SetLastHeartBeat(time.Now())
			c.JSON(http.StatusOK, instructions)
			return
		}
	}

	c.JSON(http.StatusBadRequest, []int{})
}

func GetUserNameByID(lobbyID, userID string) string {
	lobby, exists := lobbyManager.GetLobby(lobbyID)
	if !exists {
		return "undefined"
	}

	for _, member := range lobby.GetMembers() {
		if member.GetID() == userID {
			return member.GetName()
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
		if _, exists := lobbyManager.GetLobby(id); exists {
			id = generateRandomPin(6)
		} else {
			break
		}
	}
	return id
}

func generateUniqueMemberID(lobbyID string) string {
	id := generateRandomPin(12)

	lobby, exists := lobbyManager.GetLobby(lobbyID)
	if !exists {
		return "-1"
	}

	for {
		isDuplicate := false
		for _, member := range lobby.GetMembers() {
			if member.GetID() == id {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			// ID is unique, return it
			break
		}
		// Generate a new ID if the current one is not unique
		id = generateRandomPin(12)
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
	lobby, exists := lobbyManager.GetLobby(lobbyID)
	if !exists {
		return
	}

	for _, member := range lobby.GetMembers() {
		member.AddUpdateInstruction(updateInstructionType)
	}
}

func notifyLobbyMember(lobbyID string, memberID string, updateInstructionType int) {
	lobby, exists := lobbyManager.GetLobby(lobbyID)
	if !exists {
		return
	}

	for _, member := range lobby.GetMembers() {
		if member.GetID() == memberID {
			member.AddUpdateInstruction(updateInstructionType)
		}
	}
}

// Function that runs the background worker
func startBackgroundWorker() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Do background work
			doBackgroundWork()
		}
	}
}

// background work function - removes disconnected members and closes empty lobbys
func doBackgroundWork() {
	for _, lobby := range lobbyManager.GetAllLobbies() {
		for _, member := range lobby.GetMembers() {
			if time.Since(member.GetLastHeartBeat()) > time.Minute {
				lobby.RemoveMember(member.GetID())
				notifyLobbyMembers(lobby.GetID(), InstructionUpdateLobbyMembers)
			}
		}

		if len(lobby.GetMembers()) == 0 {
			lobbyManager.RemoveLobby(lobby.GetID())
		}
	}
}
