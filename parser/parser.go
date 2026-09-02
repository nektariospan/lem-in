package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"lem-in/graph"
)

// ParseInput validates and loads the map format used by lem-in.
func ParseInput(r io.Reader) (*graph.Graph, []string, error) {
	scanner := bufio.NewScanner(r)
	var rawLines []string

	readLine := func() (string, bool) {
		for scanner.Scan() {
			line := scanner.Text()
			rawLines = append(rawLines, line)
			return line, true
		}
		return "", false
	}

	antLine, ok := readLine()
	if !ok {
		return nil, nil, graph.ErrInvalidData
	}
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

	nextIsStart := false
	nextIsEnd := false
	parsingLinks := false

	for scanner.Scan() {
		line := scanner.Text()
		rawLines = append(rawLines, line)

		if strings.TrimSpace(line) == "" {
			continue
		}

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

		if strings.HasPrefix(line, "#") {
			nextIsStart = false
			nextIsEnd = false
			continue
		}

		if strings.Contains(line, "-") && !strings.Contains(line, " ") {
			parsingLinks = true
			nextIsStart = false
			nextIsEnd = false
			if err := parseLink(g, line); err != nil {
				return nil, rawLines, err
			}
			continue
		}

		if parsingLinks {
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

// parseRoom turns a room line into a graph room with integer coordinates.
func parseRoom(line string) (*graph.Room, error) {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return nil, graph.ErrInvalidData
	}
	name := fields[0]
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

func parseLink(g *graph.Graph, line string) error {
	parts := strings.Split(line, "-")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return graph.ErrInvalidData
	}
	return g.AddLink(parts[0], parts[1])
}
