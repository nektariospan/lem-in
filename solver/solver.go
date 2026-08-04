// Package solver decides how many ants travel on each path to minimise turns.
package solver

import "lem-in/flow"

// Assignment pairs a path with the number of ants that will walk it.
type Assignment struct {
	Path    flow.Path
	NumAnts int
}

// DistributeAnts computes the optimal ant-to-path assignment that minimises
// the total number of turns (makespan).
//
// Formula:
//
//	For a set of paths with lengths L[0] ≤ L[1] ≤ ... ≤ L[k-1], the minimum
//	turns T is the smallest integer such that:
//	    Σ max(0, T - L[i] + 1) >= N
//	Each path i then gets  n[i] = T - L[i] + 1  ants (≥ 1).
//
// We try using 1 path, then 2, etc., and pick the combination that yields the
// fewest turns. This handles cases where a longer additional path still helps.
func DistributeAnts(paths []flow.Path, numAnts int) []Assignment {
	if len(paths) == 0 || numAnts <= 0 {
		return nil
	}

	bestTurns := 1<<31 - 1
	bestK := 1

	// Try using the first k paths (already sorted by length ascending).
	for k := 1; k <= len(paths); k++ {
		t := turnsForK(paths[:k], numAnts)
		if t < bestTurns {
			bestTurns = t
			bestK = k
		}
	}

	chosen := paths[:bestK]
	assignments := make([]Assignment, bestK)
	for i, p := range chosen {
		n := bestTurns - (len(p) - 1) + 1 - 1
		// n[i] = T - L[i] + 1  where L[i] = len(p)-1 (number of steps/tunnels)
		n = bestTurns - (len(p) - 1) + 1
		if n < 1 {
			n = 1
		}
		assignments[i] = Assignment{Path: p, NumAnts: n}
	}

	// Adjust rounding: total might exceed numAnts due to ceiling effects.
	// Trim excess from the longest paths first.
	total := 0
	for _, a := range assignments {
		total += a.NumAnts
	}
	excess := total - numAnts
	for i := len(assignments) - 1; i >= 0 && excess > 0; i-- {
		reduce := assignments[i].NumAnts - 1
		if reduce > excess {
			reduce = excess
		}
		assignments[i].NumAnts -= reduce
		excess -= reduce
	}

	// Remove paths with 0 ants.
	result := assignments[:0]
	for _, a := range assignments {
		if a.NumAnts > 0 {
			result = append(result, a)
		}
	}
	return result
}

// turnsForK computes the minimum turns when using exactly the first k paths.
// paths must be sorted by length ascending.
func turnsForK(paths []flow.Path, numAnts int) int {
	// L[i] = number of tunnels (steps) on path i = len(path)-1
	// Start T at the shortest path length and increment until capacity ≥ numAnts.
	T := len(paths[0]) - 1 // minimum possible turns
	for {
		total := 0
		for _, p := range paths {
			l := len(p) - 1
			slots := T - l + 1
			if slots > 0 {
				total += slots
			}
		}
		if total >= numAnts {
			return T
		}
		T++
	}
}
