// Package simulation runs the turn-by-turn ant movement and produces the
// output move lines.
package simulation

import (
	"fmt"
	"strings"

	"lem-in/solver"
)

// antState tracks the position of a single ant during the simulation.
type antState struct {
	id   int
	path []string // full path from start to end
	step int      // index of the room the ant currently occupies
	done bool
}

// Run executes the simulation given the path assignments and returns each
// turn's move string (e.g. "L1-roomA L3-roomB").
// endRoom is the name of ##end so we can allow unlimited occupancy there.
func Run(assignments []solver.Assignment, endRoom string) []string {
	// Build the flat list of ants, grouped by path (shortest path first).
	// Within a path, assign ant IDs in order: they enter one per turn.
	ants := make([]*antState, 0)
	antID := 1
	for _, a := range assignments {
		for i := 0; i < a.NumAnts; i++ {
			ants = append(ants, &antState{
				id:   antID,
				path: a.Path,
				step: 0, // index 0 = start room (already there, not printed)
			})
			antID++
		}
	}

	var turns []string

	// Each ant enters the path at a staggered turn so they don't collide at
	// the start room. Ant i on a path starts moving on turn i+1 (1-based).
	// We track per-path launch offset.
	type pathGroup struct {
		path      []string
		launched  int // how many ants have been launched so far
		totalAnts int
	}

	// Map each ant to its group index and launch turn.
	// Ants on the same path are launched one per turn.
	groups := make([]*pathGroup, len(assignments))
	for i, a := range assignments {
		groups[i] = &pathGroup{
			path:      a.Path,
			totalAnts: a.NumAnts,
		}
	}

	// Assign each ant a launch turn (turn on which it makes its first move).
	antGroupIdx := make([]int, len(ants))
	antLaunchTurn := make([]int, len(ants))
	{
		gi := 0
		launchCounter := make([]int, len(groups))
		for ai, ant := range ants {
			// Find the group for this ant (sequential assignment)
			_ = ant
			antGroupIdx[ai] = gi
			antLaunchTurn[ai] = launchCounter[gi] + 1 // 1-based
			launchCounter[gi]++
			// Advance group index when group is full
			if launchCounter[gi] >= groups[gi].totalAnts {
				gi++
			}
		}
	}

	turn := 1
	for {
		// occupied tracks which rooms are claimed this turn (excluding end).
		occupied := make(map[string]bool)

		var moveParts []string

		// Priority: ants further along their path move first (avoid blocking).
		// Sort ants by descending step (in-place stable, k is small).
		sortByStepDesc(ants)

		for _, ant := range ants {
			if ant.done {
				continue
			}
			// Not yet launched
			if antLaunchTurn[ant.id-1] > turn {
				continue
			}
			nextStep := ant.step + 1
			if nextStep >= len(ant.path) {
				ant.done = true
				continue
			}
			nextRoom := ant.path[nextStep]

			// Can we move there?
			if nextRoom != endRoom && occupied[nextRoom] {
				continue // blocked this turn
			}

			// Move!
			ant.step = nextStep
			if nextRoom == endRoom {
				ant.done = true
			} else {
				occupied[nextRoom] = true
			}
			moveParts = append(moveParts, fmt.Sprintf("L%d-%s", ant.id, nextRoom))
		}

		if len(moveParts) > 0 {
			turns = append(turns, strings.Join(moveParts, " "))
		}

		// Check if all ants are done.
		allDone := true
		for _, ant := range ants {
			if !ant.done {
				allDone = false
				break
			}
		}
		if allDone {
			break
		}

		turn++

		// Safety: if we've exceeded a reasonable bound, break (shouldn't happen).
		if turn > 100000 {
			break
		}
	}

	return turns
}

// sortByStepDesc sorts ants by descending step value (furthest along first).
// Uses insertion sort since slices are small in practice.
func sortByStepDesc(ants []*antState) {
	for i := 1; i < len(ants); i++ {
		for j := i; j > 0 && ants[j].step > ants[j-1].step; j-- {
			ants[j], ants[j-1] = ants[j-1], ants[j]
		}
	}
}
