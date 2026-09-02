# lem-in

Go implementation of the lem-in ant colony problem.

The program reads a map, validates it, finds valid paths from `##start` to `##end`, distributes the ants, and simulates their movement while respecting room and tunnel constraints.

## Usage

```bash
# run a sample file from the project root
go run . example00.txt

# also works when you pass the explicit path
go run . tests/example00.txt

# stdin works too
cat tests/example00.txt | go run .
```

If a bare filename is passed, the program falls back to the `tests/` directory automatically.

## Input format

```text
<number_of_ants>
##start
<room_name> <x> <y>
...
##end
<room_name> <x> <y>
<roomA>-<roomB>
...
```

Rules:
- `##start` and `##end` are the only special directives.
- Any other `#...` line is treated as a comment and ignored.
- Room names must not start with `L` or `#`.
- Coordinates must be integers.
- Links are undirected and each tunnel may only be used once per turn.

## Output format

The program prints the original map, followed by the turn count and the generated movement lines.

```text
Number of turns: 6
L1-2
L1-3 L2-2
...
```

Each move uses the format:

```text
L<number>-<room_name>
```

Multiple moves on the same line happen in the same turn.

## Validation rules

The parser rejects invalid input such as:
- ant count not positive
- missing `##start` or `##end`
- duplicate room names
- invalid room coordinates
- self-links or links to unknown rooms
- no valid path between start and end

## Example files

All examples and validation maps are stored in the `tests/` directory, including:
- `example00.txt` to `example07.txt`
- `badexample00.txt` and `badexample01.txt`
- small generated valid checks like `random_valid01.txt`

## Testing

```bash
go test ./...
```

This project includes unit tests for parsing, flow calculation, solver logic, and simulation behavior.

## Project structure

```text
lem-in/
├── main.go
├── graph/
├── parser/
├── flow/
├── solver/
├── simulation/
├── tests/
└── README.md
```
