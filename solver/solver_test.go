package solver

import (
	"testing"

	"lem-in/flow"
)

func TestDistributeSinglePath(t *testing.T) {
	paths := []flow.Path{{"s", "a", "e"}}
	assignments := DistributeAnts(paths, 3)
	if len(assignments) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(assignments))
	}
	if assignments[0].NumAnts != 3 {
		t.Fatalf("expected 3 ants on path, got %d", assignments[0].NumAnts)
	}
}

func TestDistributeTwoEqualPaths(t *testing.T) {
	paths := []flow.Path{
		{"s", "a", "e"},
		{"s", "b", "e"},
	}
	assignments := DistributeAnts(paths, 4)
	total := 0
	for _, a := range assignments {
		total += a.NumAnts
	}
	if total != 4 {
		t.Fatalf("expected total 4 ants, got %d", total)
	}
}

func TestDistributeTwoDiffLengthPaths(t *testing.T) {
	paths := []flow.Path{
		{"s", "r1", "e"},
		{"s", "r2", "r3", "e"},
	}
	assignments := DistributeAnts(paths, 4)
	total := 0
	for _, a := range assignments {
		total += a.NumAnts
	}
	if total != 4 {
		t.Fatalf("expected total 4 ants, got %d", total)
	}
}

func TestTurnsForKFormula(t *testing.T) {
	paths := []flow.Path{{"s", "e"}}
	got := turnsForK(paths, 5)
	if got != 5 {
		t.Fatalf("expected 5 turns, got %d", got)
	}
}
