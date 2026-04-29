package ws

import (
	"sync"

	"github.com/orbit/orbit/internal/metrics"
)

// Gateway manages the active local WebSocket connections.
type Gateway struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	mu         sync.RWMutex

	maxConnsPerUser int
	userConns       map[string]int // userID -> active connection count
}

func NewGateway(maxConnsPerUser int) *Gateway {
	return &Gateway{
		Register:        make(chan *Client),
		Unregister:      make(chan *Client),
		Clients:         make(map[*Client]bool),
		maxConnsPerUser: maxConnsPerUser,
		userConns:       make(map[string]int),
	}
}

// ConnCount returns the current number of open connections for a userID.
// Returns (count, overLimit) — overLimit is true if adding one more would exceed the cap.
func (g *Gateway) AllowConnection(userID string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.userConns[userID] < g.maxConnsPerUser
}

func (g *Gateway) Run() {
	for {
		select {
		case client := <-g.Register:
			g.mu.Lock()
			g.Clients[client] = true
			g.userConns[client.UserID]++
			metrics.ActiveConnections.Inc()
			g.mu.Unlock()

		case client := <-g.Unregister:
			g.mu.Lock()
			if _, ok := g.Clients[client]; ok {
				delete(g.Clients, client)
				close(client.Send)
				metrics.ActiveConnections.Dec()
				if g.userConns[client.UserID] > 0 {
					g.userConns[client.UserID]--
				}
				if g.userConns[client.UserID] == 0 {
					delete(g.userConns, client.UserID)
				}
			}
			g.mu.Unlock()
		}
	}
}

func (g *Gateway) BroadcastLocal(channel string, env interface{}) {
	// Routing is handled by the router; left as extension point.
}

