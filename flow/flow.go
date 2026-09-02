// Package flow implements node-split max-flow (Edmonds-Karp) to find the
// maximum set of room-disjoint paths from ##start to ##end.
package flow

import (
	"lem-in/graph"
)

// nodeID returns the integer IDs for the in/out split of a room.
// start and end are not split (they have unlimited capacity).
// For room r: in-node = 2*index, out-node = 2*index+1.
// start uses only its out-node; end uses only its in-node.

// network is the internal max-flow representation.
type network struct {
	// cap[u][v] = remaining capacity on edge u→v
	cap  [][]int
	size int // total number of nodes

	// mapping: room name → (inNode, outNode)
	inNode  map[string]int
	outNode map[string]int
}

// buildNetwork creates the flow graph with one in/out node per room.
func buildNetwork(g *graph.Graph) *network {
	// Assign indices: for each room, allocate 2 nodes (in, out).
	// start and end get the same treatment but their split-edge has ∞ capacity.
	idx := 0
	inNode := make(map[string]int, len(g.Rooms))
	outNode := make(map[string]int, len(g.Rooms))

	for name := range g.Rooms {
		inNode[name] = idx
		outNode[name] = idx + 1
		idx += 2
	}

	size := idx
	cap := make([][]int, size)
	for i := range cap {
		cap[i] = make([]int, size)
	}

	// Internal node edges (in → out)
	for name := range g.Rooms {
		u := inNode[name]
		v := outNode[name]
		if name == g.StartRoom || name == g.EndRoom {
			cap[u][v] = 1<<31 - 1 // ∞
		} else {
			cap[u][v] = 1 // room capacity = 1
		}
	}

	// Tunnel edges: for each undirected link A-B add:
	//   A_out → B_in  (capacity 1)
	//   B_out → A_in  (capacity 1)
	seen := make(map[string]bool)
	for a, neighbors := range g.Links {
		for _, b := range neighbors {
			key := a + "~" + b
			rev := b + "~" + a
			if seen[key] || seen[rev] {
				continue
			}
			seen[key] = true
			cap[outNode[a]][inNode[b]] = 1
			cap[outNode[b]][inNode[a]] = 1
		}
	}

	return &network{cap: cap, size: size, inNode: inNode, outNode: outNode}
}

// bfs finds an augmenting path from src to dst in the residual graph.
// Returns the parent array (parent[v] = u means u→v on the path), or nil.
func (n *network) bfs(src, dst int) []int {
	parent := make([]int, n.size)
	for i := range parent {
		parent[i] = -1
	}
	parent[src] = src
	queue := []int{src}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		if u == dst {
			return parent
		}
		for v := 0; v < n.size; v++ {
			if parent[v] == -1 && n.cap[u][v] > 0 {
				parent[v] = u
				queue = append(queue, v)
			}
		}
	}
	return nil
}

// edmondsKarp runs BFS-based max flow from src to dst, mutating cap in place.
// Returns the flow value.
func (n *network) edmondsKarp(src, dst int) int {
	total := 0
	for {
		parent := n.bfs(src, dst)
		if parent == nil {
			break
		}
		// Find bottleneck
		flow := 1<<31 - 1
		for v := dst; v != src; v = parent[v] {
			u := parent[v]
			if n.cap[u][v] < flow {
				flow = n.cap[u][v]
			}
		}
		// Update residual capacities
		for v := dst; v != src; v = parent[v] {
			u := parent[v]
			n.cap[u][v] -= flow
			n.cap[v][u] += flow
		}
		total += flow
	}
	return total
}

// originalCap stores the capacities before flow so we can detect used edges.
// We snapshot it before running edmondsKarp.
type snapshot [][]int

func snapshotCap(cap [][]int) snapshot {
	s := make(snapshot, len(cap))
	for i, row := range cap {
		s[i] = make([]int, len(row))
		copy(s[i], row)
	}
	return s
}
