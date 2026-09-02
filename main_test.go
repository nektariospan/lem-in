package main

import "testing"

func TestResolveInputPathUsesTestsDirFallback(t *testing.T) {
	path, err := resolveInputPath("example00.txt")
	if err != nil {
		t.Fatalf("expected resolveInputPath to find example00.txt, got error: %v", err)
	}
	if path != "tests/example00.txt" {
		t.Fatalf("expected tests/example00.txt, got %q", path)
	}
}

func TestResolveInputPathKeepsExplicitPath(t *testing.T) {
	path, err := resolveInputPath("tests/example00.txt")
	if err != nil {
		t.Fatalf("expected resolveInputPath to accept explicit path, got error: %v", err)
	}
	if path != "tests/example00.txt" {
		t.Fatalf("expected tests/example00.txt, got %q", path)
	}
}

func TestTurnSummaryUsesMoveCount(t *testing.T) {
	moves := []string{"L1-a", "L2-b", "L3-c"}
	if got := turnSummary(moves); got != "Number of turns: 3" {
		t.Fatalf("expected turn summary to report 3 turns, got %q", got)
	}
}
