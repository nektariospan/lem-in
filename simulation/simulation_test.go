package simulation

import (
	"strings"
	"testing"

	"lem-in/flow"
	"lem-in/solver"
)

func TestSimpleLinearPath(t *testing.T) {
	assignments := []solver.Assignment{
		{Path: flow.Path{"s", "a", "e"}, NumAnts: 3},
	}
	turns := Run(assignments, "e")
	if len(turns) != 4 {
		t.Fatalf("expected 4 turns, got %d: %v", len(turns), turns)
	}
}

func TestTwoPathsNoConflict(t *testing.T) {
	assignments := []solver.Assignment{
		{Path: flow.Path{"s", "a", "e"}, NumAnts: 2},
		{Path: flow.Path{"s", "b", "e"}, NumAnts: 2},
	}
	turns := Run(assignments, "e")
	if len(turns) != 3 {
		t.Fatalf("expected 3 turns, got %d: %v", len(turns), turns)
	}
}

func TestNoRoomConflict(t *testing.T) {
	assignments := []solver.Assignment{
		{Path: flow.Path{"s", "a", "b", "e"}, NumAnts: 3},
	}
	turns := Run(assignments, "e")
	for _, turn := range turns {
		roomCount := make(map[string]int)
		for _, move := range strings.Fields(turn) {
			parts := strings.SplitN(move, "-", 2)
			if len(parts) == 2 && parts[1] != "e" {
				roomCount[parts[1]]++
				if roomCount[parts[1]] > 1 {
					t.Fatalf("room conflict in turn %q: room %s occupied twice", turn, parts[1])
				}
			}
		}
	}
}
