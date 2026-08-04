// Package graph defines the core data structures for the lem-in colony map.
package graph

import "errors"

// ErrInvalidData is returned for any malformed input.
var ErrInvalidData = errors.New("ERROR: invalid data format")

// Room represents a node (chamber) in the ant colony graph.
type Room struct {
	Name string
	X, Y int
}

// Graph holds the full parsed colony map.
type Graph struct {
	Rooms     map[string]*Room    // room name → Room
	Links     map[string][]string // adjacency list (undirected)
	StartRoom string
	EndRoom   string
	NumAnts   int
}

// NewGraph allocates an empty Graph.
func NewGraph() *Graph {
	return &Graph{
		Rooms: make(map[string]*Room),
		Links: make(map[string][]string),
	}
}

// AddRoom inserts a room; returns ErrInvalidData on duplicate.
func (g *Graph) AddRoom(r *Room) error {
	if _, exists := g.Rooms[r.Name]; exists {
		return ErrInvalidData
	}
	g.Rooms[r.Name] = r
	return nil
}

// AddLink creates an undirected edge; returns ErrInvalidData if either room
// is unknown or if the link already exists.
func (g *Graph) AddLink(a, b string) error {
	if _, ok := g.Rooms[a]; !ok {
		return ErrInvalidData
	}
	if _, ok := g.Rooms[b]; !ok {
		return ErrInvalidData
	}
	if a == b {
		return ErrInvalidData
	}
	// Deduplicate
	for _, nb := range g.Links[a] {
		if nb == b {
			return ErrInvalidData
		}
	}
	g.Links[a] = append(g.Links[a], b)
	g.Links[b] = append(g.Links[b], a)
	return nil
}

// Validate performs post-parse sanity checks.
func (g *Graph) Validate() error {
	if g.StartRoom == "" || g.EndRoom == "" {
		return ErrInvalidData
	}
	if g.StartRoom == g.EndRoom {
		return ErrInvalidData
	}
	if g.NumAnts <= 0 {
		return ErrInvalidData
	}
	// BFS reachability: start must be able to reach end.
	if !g.reachable(g.StartRoom, g.EndRoom) {
		return ErrInvalidData
	}
	return nil
}

// reachable returns true if dst is reachable from src via BFS.
func (g *Graph) reachable(src, dst string) bool {
	visited := make(map[string]bool)
	queue := []string{src}
	visited[src] = true
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == dst {
			return true
		}
		for _, nb := range g.Links[cur] {
			if !visited[nb] {
				visited[nb] = true
				queue = append(queue, nb)
			}
		}
	}
	return false
}
