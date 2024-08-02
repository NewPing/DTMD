package models

import (
	"slices"
	"sync"
	"time"
)

type Member struct {
	ID                 string
	Name               string
	lastHeartBeat      time.Time
	updateInstructions []int
	newChatMessages    []ChatMessage
	mu                 sync.Mutex
}

// NewMember creates a new Member
func NewMember(id, name string) *Member {
	return &Member{
		ID:                 id,
		Name:               name,
		lastHeartBeat:      time.Now(),
		updateInstructions: make([]int, 0),
		newChatMessages:    make([]ChatMessage, 0),
	}
}

// GetID returns the ID of the member
func (m *Member) GetID() string {
	return m.ID
}

// SetID sets the ID of the member
func (m *Member) SetID(id string) {
	m.ID = id
}

// GetName returns the Name of the member
func (m *Member) GetName() string {
	return m.Name
}

// SetName sets the Name of the member
func (m *Member) SetName(name string) {
	m.Name = name
}

// GetLastHeartBeat returns the LastHeartBeat of the member
func (m *Member) GetLastHeartBeat() time.Time {
	return m.lastHeartBeat
}

// SetLastHeartBeat sets the LastHeartBeat of the member
func (m *Member) SetLastHeartBeat(lastHeartBeat time.Time) {
	m.lastHeartBeat = lastHeartBeat
}

// GetUpdateInstruction returns the UpdateInstruction slice
func (m *Member) GetUpdateInstructions() []int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.updateInstructions
}

// AddUpdateInstruction adds an update instruction if it is not already present
func (m *Member) AddUpdateInstruction(instruction int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !slices.Contains(m.updateInstructions, instruction) {
		m.updateInstructions = append(m.updateInstructions, instruction)
	}
}

// ClearUpdateInstructions clears the update instructions
func (m *Member) ClearUpdateInstructions() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updateInstructions = []int{}
}

// GetNewChatMessages returns the NewChatMessages slice
func (m *Member) GetNewChatMessages() []ChatMessage {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.newChatMessages
}

// AddNewChatMessage adds a new chat message to the NewChatMessages slice
func (m *Member) AddNewChatMessage(message ChatMessage) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.newChatMessages = append(m.newChatMessages, message)
}

// ClearNewChatMessages clears the new chat messages
func (m *Member) ClearNewChatMessages() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.newChatMessages = []ChatMessage{}
}
