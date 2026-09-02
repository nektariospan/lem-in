package simulation

import (
	"fmt"
	"strings"

	"lem-in/solver"
)

type antState struct {
	id   int
	path []string
	step int
	done bool
}

func Run(assignments []solver.Assignment, endRoom string) []string {
	ants := make([]*antState, 0)
	antID := 1
	for _, a := range assignments {
		for i := 0; i < a.NumAnts; i++ {
			ants = append(ants, &antState{
				id:   antID,
				path: a.Path,
				step: 0,
			})
			antID++
		}
	}

	var turns []string
	type pathGroup struct {
		path      []string
		launched  int
		totalAnts int
	}

	groups := make([]*pathGroup, len(assignments))
	for i, a := range assignments {
		groups[i] = &pathGroup{
			path:      a.Path,
			totalAnts: a.NumAnts,
		}
	}

	antLaunchTurn := make([]int, len(ants))
	{
		gi := 0
		launchCounter := make([]int, len(groups))
		for ai, ant := range ants {
			_ = ant
			antLaunchTurn[ai] = launchCounter[gi] + 1
			launchCounter[gi]++
			if launchCounter[gi] >= groups[gi].totalAnts {
				gi++
			}
		}
	}

	turn := 1
	for {
		occupied := make(map[string]bool)
		usedEdges := make(map[string]bool)
		var moveParts []string

		sortByStepDesc(ants)

		for _, ant := range ants {
			if ant.done || antLaunchTurn[ant.id-1] > turn {
				continue
			}
			nextStep := ant.step + 1
			if nextStep >= len(ant.path) {
				ant.done = true
				continue
			}
			fromRoom := ant.path[ant.step]
			nextRoom := ant.path[nextStep]
			edgeKey := fromRoom + "->" + nextRoom

			if usedEdges[edgeKey] {
				continue
			}
			usedEdges[edgeKey] = true

			if nextRoom != endRoom && occupied[nextRoom] {
				continue
			}

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
		if turn > 100000 {
			break
		}
	}

	return turns
}

func sortByStepDesc(ants []*antState) {
	for i := 1; i < len(ants); i++ {
		for j := i; j > 0 && ants[j].step > ants[j-1].step; j-- {
			ants[j], ants[j-1] = ants[j-1], ants[j]
		}
	}
}
