package classes

import (
	"sync"
)

type LobbyManager struct {
	lobbys map[string]*Lobby
	mu     sync.RWMutex
}

// NewLobbyManager creates a new LobbyManager
func NewLobbyManager() *LobbyManager {
	return &LobbyManager{
		lobbys: make(map[string]*Lobby),
	}
}

// GetLobby retrieves a lobby by ID
func (lm *LobbyManager) GetLobby(id string) (*Lobby, bool) {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	lobby, exists := lm.lobbys[id]
	return lobby, exists
}

// AddLobby adds a new lobby to the manager
func (lm *LobbyManager) AddLobby(l *Lobby) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	lm.lobbys[l.GetID()] = l
}

// RemoveLobby removes a lobby from the manager by ID
func (lm *LobbyManager) RemoveLobby(id string) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	delete(lm.lobbys, id)
}

// GetAllLobbies returns a copy of all lobbies
func (lm *LobbyManager) GetAllLobbies() []*Lobby {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	lobbies := make([]*Lobby, 0, len(lm.lobbys))
	for _, l := range lm.lobbys {
		lobbies = append(lobbies, l)
	}
	return lobbies
}
