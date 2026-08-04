// Package parser reads and validates the lem-in input format.
package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"lem-in/graph"
)

// ParseInput reads from r and returns:
//   - the populated Graph
//   - raw input lines (for printing later)
//   - an error (graph.ErrInvalidData) on any malformed input
func ParseInput(r io.Reader) (*graph.Graph, []string, error) {
	scanner := bufio.NewScanner(r)
	var rawLines []string

	// ---- helper: read next non-empty raw line ----
	readLine := func() (string, bool) {
		for scanner.Scan() {
			line := scanner.Text()
			rawLines = append(rawLines, line)
			return line, true
		}
		return "", false
	}

	// 1. Ant count
	antLine, ok := readLine()
	if !ok {
		return nil, nil, graph.ErrInvalidData
	}
	// Skip leading comment lines before ant count
	for strings.HasPrefix(antLine, "#") {
		antLine, ok = readLine()
		if !ok {
			return nil, nil, graph.ErrInvalidData
		}
	}
	numAnts, err := strconv.Atoi(strings.TrimSpace(antLine))
	if err != nil || numAnts <= 0 {
		return nil, rawLines, graph.ErrInvalidData
	}

	g := graph.NewGraph()
	g.NumAnts = numAnts

	// 2. Parse rooms and links
	nextIsStart := false
	nextIsEnd := false
	parsingLinks := false // once we see first link line, rooms are done

	for scanner.Scan() {
		line := scanner.Text()
		rawLines = append(rawLines, line)

		// Empty lines — skip
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Directives
		if line == "##start" {
			nextIsStart = true
			nextIsEnd = false
			continue
		}
		if line == "##end" {
			nextIsEnd = true
			nextIsStart = false
			continue
		}

		// Other comments
		if strings.HasPrefix(line, "#") {
			nextIsStart = false
			nextIsEnd = false
			continue
		}

		// Link line: contains '-' but no spaces (rooms have spaces)
		if strings.Contains(line, "-") && !strings.Contains(line, " ") {
			parsingLinks = true
			nextIsStart = false
			nextIsEnd = false
			if err := parseLink(g, line); err != nil {
				return nil, rawLines, err
			}
			continue
		}

		// Room line
		if parsingLinks {
			// A room after links have started is invalid
			return nil, rawLines, graph.ErrInvalidData
		}

		room, err := parseRoom(line)
		if err != nil {
			return nil, rawLines, err
		}

		if nextIsStart {
			if g.StartRoom != "" {
				return nil, rawLines, graph.ErrInvalidData
			}
			g.StartRoom = room.Name
			nextIsStart = false
		} else if nextIsEnd {
			if g.EndRoom != "" {
				return nil, rawLines, graph.ErrInvalidData
			}
			g.EndRoom = room.Name
			nextIsEnd = false
		}

		if err := g.AddRoom(room); err != nil {
			return nil, rawLines, err
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return nil, rawLines, fmt.Errorf("ERROR: invalid data format")
	}

	if err := g.Validate(); err != nil {
		return nil, rawLines, err
	}

	return g, rawLines, nil
}

// parseRoom parses a line of the form "name x y".
func parseRoom(line string) (*graph.Room, error) {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return nil, graph.ErrInvalidData
	}
	name := fields[0]
	// Room names must not start with 'L' or '#'
	if len(name) == 0 || name[0] == 'L' || name[0] == '#' {
		return nil, graph.ErrInvalidData
	}
	x, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil, graph.ErrInvalidData
	}
	y, err := strconv.Atoi(fields[2])
	if err != nil {
		return nil, graph.ErrInvalidData
	}
	return &graph.Room{Name: name, X: x, Y: y}, nil
}

// parseLink parses a line of the form "roomA-roomB".
func parseLink(g *graph.Graph, line string) error {
	parts := strings.Split(line, "-")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return graph.ErrInvalidData
	}
	return g.AddLink(parts[0], parts[1])
}
