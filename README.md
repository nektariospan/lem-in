# lem-in

A Go implementation of the **lem-in** ant colony pathfinding problem from Zone01.

The program reads an ant colony map, finds the optimal set of paths from `##start` to `##end`, and simulates all ants moving through the colony in the minimum possible number of turns.

---

## Usage

```bash
# Compile
go build -o lem-in .

# Run with a file argument
./lem-in example00.txt

# Run via stdin
./lem-in < example00.txt

# Run without compiling
go run . example00.txt
```

---

## Input Format

```
<number_of_ants>
##start
<start_room_name> <x> <y>
<room_name> <x> <y>
...
##end
<end_room_name> <x> <y>
<roomA>-<roomB>
<roomC>-<roomD>
...
```

- Lines starting with `#` (but not `##start` / `##end`) are treated as comments and are printed but ignored.
- Room names must not start with `L` or `#`.
- Coordinates must be integers.
- Links are undirected; each pair of rooms can have at most one tunnel.

### Example input

```
4
##start
0 0 3
2 2 5
3 4 0
##end
1 8 3
0-2
2-3
3-1
```

---

## Output Format

The program prints the original input, a blank line, then one line per turn showing every ant that moved:

```
4
##start
0 0 3
...

L1-2
L1-3 L2-2
L1-1 L2-3 L3-2
L2-1 L3-3 L4-2
L3-1 L4-3
L4-1
```

- `Lx-y` — ant number `x` moved to room `y`
- Multiple moves on one line = same turn
- Only ants that moved are shown per turn

---

## Error Handling

Any invalid input prints to **stderr** and exits with code `1`:

```
ERROR: invalid data format
```

Invalid cases include:

| Condition | Detected by |
|---|---|
| Ant count ≤ 0 or not an integer | parser |
| Missing `##start` or `##end` | parser / graph.Validate |
| `##start` == `##end` | graph.Validate |
| Duplicate room name | graph.AddRoom |
| Room name starts with `L` or `#` | parser.parseRoom |
| Invalid room coordinates | parser.parseRoom |
| Link to unknown room | graph.AddLink |
| Duplicate or self-loop link | graph.AddLink |
| No path between start and end | graph.Validate (BFS) |

---

## How It Works

### 1. Parsing (`parser/`)

Reads the input line by line using `bufio.Scanner`. Builds a `Graph` struct (rooms + adjacency list), stores raw lines for reprinting, and validates all format rules.

### 2. Node-Split Max Flow (`flow/`)

To respect the **room capacity = 1** rule, every room is split into two nodes:

```
room_in  →  room_out   (capacity 1)
##start / ##end        (capacity ∞)
tunnel A-B  →  A_out → B_in  and  B_out → A_in  (capacity 1)
```

**Edmonds-Karp** (BFS-based max flow) is run on this network to find the maximum number of room-disjoint paths. Each unit of flow corresponds to one usable path. Paths are extracted from the residual graph and sorted shortest-first.

### 3. Ant Distribution (`solver/`)

Given paths of lengths `L[0] ≤ L[1] ≤ ... ≤ L[k-1]` and `N` total ants, the minimum number of turns `T` satisfies:

```
sum( max(0, T - L[i] + 1) for all i ) >= N
```

Path `i` is assigned `n[i] = T - L[i] + 1` ants. This minimises the turn when the **last ant** arrives (the makespan).

### 4. Simulation (`simulation/`)

Ants on the same path are staggered one turn apart to avoid collisions at the start. Each turn:

1. Ants furthest along their path move first (prevents forward blocking)
2. Before moving, check the destination room is not already occupied
3. Mark destination as occupied (unless it is `##end`, which has unlimited capacity)
4. Collect all `Lx-room` tokens for this turn and emit one output line

### 5. Output (`main.go`)

Prints all raw input lines, a blank line, then each turn's move string.

---

## Project Structure

```
lem-in/
├── main.go                  Entry point — arg/stdin handling, pipeline wiring
├── graph/
│   └── graph.go             Room and Graph types; AddRoom, AddLink, Validate
├── parser/
│   ├── parser.go            Full input parser with all validation rules
│   └── parser_test.go       Unit tests for parsing edge cases
├── flow/
│   ├── flow.go              Node-split capacity network + Edmonds-Karp
│   ├── paths.go             Path extraction from residual graph
│   └── flow_test.go         Unit tests for single path, two paths, no path
├── solver/
│   ├── solver.go            Ant distribution formula; optimal turn calculation
│   └── solver_test.go       Unit tests for distribution correctness
├── simulation/
│   ├── simulation.go        Turn-by-turn simulation with occupancy tracking
│   └── simulation_test.go   Unit tests for movement, conflicts, turn count
└── tests/
    └── map01-06.txt         Sample input maps for manual testing
```

---

## Running Tests

```bash
go test ./...
```

All packages have unit tests covering:

- Valid inputs and expected turn counts
- Parser: zero ants, missing start/end, duplicate rooms, bad room names, unknown links
- Flow: single path, two parallel paths, no path
- Solver: equal-length paths, different-length paths, distribution formula
- Simulation: no room conflicts, correct turn count, all ants arrive

---

## Constraints

- Written entirely in **Go**
- Only **standard library** packages are used (`bufio`, `fmt`, `os`, `strings`, `strconv`, `errors`)
- No global state; all data flows through explicit function parameters
