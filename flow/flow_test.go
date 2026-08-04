package flow

import (
	"strings"
	"testing"

	"lem-in/graph"
)

func makeGraph(rooms []string, links []string, start, end string, n int) *graph.Graph {
	g := graph.NewGraph()
	g.NumAnts = n
	g.StartRoom = start
	g.EndRoom = end
	for _, r := range rooms {
		parts := strings.Fields(r)
		name := parts[0]
		_ = g.AddRoom(&graph.Room{Name: name})
	}
	for _, l := range links {
		p := strings.Split(l, "-")
		_ = g.AddLink(p[0], p[1])
	}
	return g
}

func TestSinglePath(t *testing.T) {
	g := makeGraph(
		[]string{"s 0 0", "a 1 0", "e 2 0"},
		[]string{"s-a", "a-e"},
		"s", "e", 3,
	)
	paths, err := FindPaths(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 1 {
		t.Fatalf("expected 1 path, got %d", len(paths))
	}
	if len(paths[0]) != 3 {
		t.Fatalf("expected path length 3, got %d: %v", len(paths[0]), paths[0])
	}
}

func TestTwoPaths(t *testing.T) {
	g := makeGraph(
		[]string{"s 0 0", "a 1 0", "b 1 1", "e 2 0"},
		[]string{"s-a", "s-b", "a-e", "b-e"},
		"s", "e", 4,
	)
	paths, err := FindPaths(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(paths))
	}
}

func TestNoPath(t *testing.T) {
	g := makeGraph(
		[]string{"s 0 0", "e 1 0"},
		[]string{},
		"s", "e", 1,
	)
	_, err := FindPaths(g)
	if err == nil {
		t.Fatal("expected error for no path")
	}
}
