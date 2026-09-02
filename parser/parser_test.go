package parser

import (
	"strings"
	"testing"

	"lem-in/graph"
)

func parse(input string) error {
	_, _, err := ParseInput(strings.NewReader(input))
	return err
}

func TestValidSimple(t *testing.T) {
	input := "3\n##start\nstart 0 0\n##end\nend 5 0\na 1 0\nstart-a\na-end"
	_, _, err := ParseInput(strings.NewReader(input))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestZeroAnts(t *testing.T) {
	if err := parse("0\n##start\ns 0 0\n##end\ne 1 0\ns-e"); err == nil {
		t.Fatal("expected error for 0 ants")
	}
}

func TestNegativeAnts(t *testing.T) {
	if err := parse("-1\n##start\ns 0 0\n##end\ne 1 0\ns-e"); err == nil {
		t.Fatal("expected error for negative ants")
	}
}

func TestMissingStart(t *testing.T) {
	input := "2\nroom 0 0\n##end\nend 1 0\nroom-end"
	if err := parse(input); err == nil {
		t.Fatal("expected error for missing ##start")
	}
}

func TestMissingEnd(t *testing.T) {
	input := "2\n##start\nstart 0 0\nroom 1 0\nstart-room"
	if err := parse(input); err == nil {
		t.Fatal("expected error for missing ##end")
	}
}

func TestDuplicateRoom(t *testing.T) {
	input := "2\n##start\ns 0 0\n##end\ne 1 0\ns 2 0\ns-e"
	if err := parse(input); err == nil {
		t.Fatal("expected error for duplicate room")
	}
}

func TestRoomNameStartsWithL(t *testing.T) {
	input := "2\n##start\nLstart 0 0\n##end\ne 1 0\nLstart-e"
	if err := parse(input); err == nil {
		t.Fatal("expected error for room starting with L")
	}
}

func TestNoPath(t *testing.T) {
	input := "2\n##start\ns 0 0\n##end\ne 1 0\na 2 0\nb 3 0\na-b"
	if err := parse(input); err == nil {
		t.Fatal("expected error for no path between start and end")
	}
}

func TestUnknownRoomInLink(t *testing.T) {
	input := "2\n##start\ns 0 0\n##end\ne 1 0\ns-ghost"
	if err := parse(input); err == nil {
		t.Fatal("expected error for link to unknown room")
	}
}

func TestStartEqualsEnd(t *testing.T) {
	g := graph.NewGraph()
	g.NumAnts = 1
	_ = g.AddRoom(&graph.Room{Name: "r", X: 0, Y: 0})
	g.StartRoom = "r"
	g.EndRoom = "r"
	if err := g.Validate(); err == nil {
		t.Fatal("expected error for start == end")
	}
}

func TestCommentLinesPreserved(t *testing.T) {
	input := "2\n# comment\n##start\ns 0 0\n##end\ne 1 0\ns-e"
	_, raw, err := ParseInput(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, l := range raw {
		if l == "# comment" {
			found = true
		}
	}
	if !found {
		t.Fatal("comment line not preserved in raw output")
	}
}

func TestOnlyStartAndEndAreSpecialDirectives(t *testing.T) {
	input := "2\n#rooms\n##start\ns 0 0\n##end\ne 1 0\ns-e"
	if err := parse(input); err != nil {
		t.Fatalf("expected #rooms and other comments to be ignored, got: %v", err)
	}
}

func TestUnknownDoubleHashDirectiveIsIgnoredAsComment(t *testing.T) {
	input := "2\n##rooms\n##start\ns 0 0\n##end\ne 1 0\ns-e"
	if err := parse(input); err != nil {
		t.Fatalf("expected ##rooms to be ignored as a comment, got: %v", err)
	}
}
