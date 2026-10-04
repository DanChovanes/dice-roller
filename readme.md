# whim

`whim` is a tiny Go command-line dice roller for tabletop and story-driven rolls. It parses dice expressions, rolls them with a configurable randomness table, and prints a readable result along with optional source details.

It is intentionally small, fast, and easy to script in terminals, shells, and other tooling.

## Features

- Roll standard dice: `d2`, `d4`, `d6`, `d8`, `d10`, `d12`, `d20`, `d100`
- Accept multi-roll expressions: `4d6`, `2d10+3`, `d20adv`, `d20dis`
- Combine multiple groups in one command: `2d6+1 d20adv 3d8`
- Accept whitespace- or comma-separated inputs
- Optional `-i` flag to reveal the randomness source used
- ANSI color output with `--color=auto|always|never`
- Respects `NO_COLOR` when color output is automatic
- Built-in validation for invalid expressions, oversized roll batches, and impossible die sizes

## Installation

### From source

```bash
go install ./cmd/whim
```

### Run directly

```bash
go run ./cmd/whim --help
```

## Usage

```bash
whim [OPTIONS] DICE...
```

For the built-in quick start and notation guide, run `whim --help` (or `whim -h`).
The `-i` option shows which randomness source produced the roll, and
`--color=auto|always|never` controls colored output.

Examples:

```bash
whim d20
whim 4d6+2
whim 2d20adv+5
whim -i 3d8 d20dis+1
whim --color=always 1d100
```

## Supported expression syntax

- `d20` — one twenty-sided die
- `4d6` — four six-sided dice
- `2d10+3` — two ten-sided dice with a +3 modifier
- `d20adv` — roll 1d20 with advantage
- `d20dis` — roll 1d20 with disadvantage
- `d20adv+8` or `d20+8adv` — advantage plus modifier

The parser supports whitespace-separated expressions and comma-separated expressions, such as:

```bash
whim "4d6, d8, d20adv+2"
whim "4d6 d8 d20adv+2"
```

## Example output

```text
$ whim -i 2d6+3 d8
The dice have spoken!

  1. 2d6+3
     rolls: [4] [2]
     modifier: +3
     ==> TOTAL: 9

  2. 1d8
     rolls: [7]
     ==> TOTAL: 7

A whisper from the source:
  - Go math/rand pseudo-random generator (non-cryptographic)
```

## Project structure

- `cmd/whim` — CLI entrypoint
- `internal/dice` — expression parsing and roll logic
- `internal/randomness` — randomness sources and table definitions
- `project-ideas.md` — notes describing possible narrative randomness sources

## Notes

The default source is Go's `math/rand` pseudo-random generator. It is useful for simple local rolls, but isn't really **fun**, ya know? 

## License

This project does not currently declare a license. If you plan to distribute or publish it, add a license before doing so.
