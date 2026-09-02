package flow

import (
	"lem-in/graph"
)

// Path is an ordered slice of room names from ##start to ##end (inclusive).
type Path []string

// FindPaths returns the valid paths from start to end, sorted by length.
func FindPaths(g *graph.Graph) ([]Path, error) {
	net := buildNetwork(g)

	src := net.outNode[g.StartRoom]
	dst := net.inNode[g.EndRoom]

	// Snapshot original capacity so we can detect which edges carry flow.
	before := snapshotCap(net.cap)

	flowVal := net.edmondsKarp(src, dst)
	if flowVal == 0 {
		return nil, graph.ErrInvalidData
	}

	// Extract individual paths by tracing flow-carrying edges.
	// An edge u→v carries flow if before[u][v] > net.cap[u][v].
	paths := make([]Path, 0, flowVal)
	for i := 0; i < flowVal; i++ {
		p := tracePath(net, before, src, dst, g)
		if p == nil {
			break
		}
		paths = append(paths, p)
	}

	// Sort by length ascending (simple insertion sort – k is small).
	for i := 1; i < len(paths); i++ {
		for j := i; j > 0 && len(paths[j]) < len(paths[j-1]); j-- {
			paths[j], paths[j-1] = paths[j-1], paths[j]
		}
	}

	return paths, nil
}

// tracePath follows one unit of flow through the residual graph, from src to
// dst. It "consumes" flow along the way by restoring capacity on used edges,
// so subsequent calls extract different paths.
func tracePath(net *network, before snapshot, src, dst int, g *graph.Graph) Path {
	// Build reverse lookup: node ID → room name
	nodeToRoom := make(map[int]string, len(g.Rooms)*2)
	for name := range g.Rooms {
		nodeToRoom[net.inNode[name]] = name
		nodeToRoom[net.outNode[name]] = name
	}

	path := []int{src}
	visited := make(map[int]bool)
	visited[src] = true

	cur := src
	for cur != dst {
		moved := false
		for v := 0; v < net.size; v++ {
			if visited[v] {
				continue
			}
			// Edge carries flow: original capacity was used
			if before[cur][v] > net.cap[cur][v] {
				// Consume this flow unit
				before[cur][v]--
				net.cap[cur][v]++ // restore residual so it won't be re-used
				path = append(path, v)
				visited[v] = true
				cur = v
				moved = true
				break
			}
		}
		if !moved {
			return nil
		}
	}

	// Convert node IDs → room names, deduplicate split-node pairs.
	var rooms Path
	seen := make(map[string]bool)
	for _, node := range path {
		name, ok := nodeToRoom[node]
		if !ok {
			continue
		}
		if !seen[name] {
			seen[name] = true
			rooms = append(rooms, name)
		}
	}
	// Make sure end room is included
	endName := nodeToRoom[dst]
	if len(rooms) == 0 || rooms[len(rooms)-1] != endName {
		rooms = append(rooms, endName)
	}

	return rooms
}
