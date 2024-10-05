## Crogo 🐊

Same as [Croissant 🥐](https://github.com/super7ramp/croissant), but in Go.

Work in progress.

### Build

```shell
go build -o crogo cmd/crogo/main.go
```

### Usage

See `crogo --help`:

```
🐊 Welcome to Crogo, a crossword solver that bites

Examples:

$ crogo "...,...,..." # The grid is a comma-separated list of rows.
[[B A A] [A B B] [B A A]]

$ crogo "A..,B..,C.." # '.' means an empty cell
[[A B A] [B A B] [C H A]]

$ crogo "ALL,...,..." --count 3 # --count allows to get more than one solution
[[A L L] [B A A] [A B B]]
[[A L L] [K A A] [A B B]]
[[A L L] [K A A] [E B B]]

Usage:
  crogo <GRID> [flags]

Flags:
  -c, --count int   The desired number of solutions (default 1)
  -h, --help        help for crogo
```
