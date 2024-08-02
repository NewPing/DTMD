package models

import (
	"sync"
)

type Lobby struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	members []*Member // Use a slice of pointers to Member
	mu      sync.Mutex
}

// NewLobby creates a new Lobby
func NewLobby(id, name string) *Lobby {
	return &Lobby{
		ID:      id,
		Name:    name,
		members: make([]*Member, 0),
	}
}

// GetID returns the ID of the lobby
func (l *Lobby) GetID() string {
	return l.ID
}

// SetID sets the ID of the lobby
func (l *Lobby) SetID(id string) {
	l.ID = id
}

// GetName returns the Name of the lobby
func (l *Lobby) GetName() string {
	return l.Name
}

// SetName sets the Name of the lobby
func (l *Lobby) SetName(name string) {
	l.Name = name
}

// GetMembers returns the Members slice
func (l *Lobby) GetMembers() []*Member {
	l.mu.Lock()
	defer l.mu.Unlock()
	// Return a copy of the slice to prevent external modifications
	membersCopy := make([]*Member, len(l.members))
	copy(membersCopy, l.members)
	return membersCopy
}

// AddMember adds a new member to the Members slice
func (l *Lobby) AddMember(mem *Member) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.members = append(l.members, mem)
}

// RemoveMember removes a member from the Members slice by ID
func (l *Lobby) RemoveMember(id string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for i, mem := range l.members {
		if mem.GetID() == id {
			l.members = append(l.members[:i], l.members[i+1:]...)
			break
		}
	}
}
