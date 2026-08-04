package main

import (
	"fmt"
	"os"
	"strings"

	"lem-in/flow"
	"lem-in/graph"
	"lem-in/parser"
	"lem-in/simulation"
	"lem-in/solver"
)

func main() {
	var input *os.File
	switch len(os.Args) {
	case 1:
		input = os.Stdin
	case 2:
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, graph.ErrInvalidData.Error())
			os.Exit(1)
		}
		defer f.Close()
		input = f
	default:
		fmt.Fprintln(os.Stderr, graph.ErrInvalidData.Error())
		os.Exit(1)
	}

	g, rawLines, err := parser.ParseInput(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, graph.ErrInvalidData.Error())
		os.Exit(1)
	}

	paths, err := flow.FindPaths(g)
	if err != nil {
		fmt.Fprintln(os.Stderr, graph.ErrInvalidData.Error())
		os.Exit(1)
	}

	assignments := solver.DistributeAnts(paths, g.NumAnts)
	moves := simulation.Run(assignments, g.EndRoom)

	fmt.Println(strings.Join(rawLines, "\n"))
	fmt.Println()
	for _, line := range moves {
		fmt.Println(line)
	}
}
