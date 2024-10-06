package cmd

import (
	"crogo/pkg/crogo"
	"crogo/pkg/dictionaries"
	"crogo/pkg/solver"
	"errors"
	"fmt"
	"iter"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// count is the desired number of solutions.
var count int

// solverName is the name of the desired solver.
var solverName string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "crogo <GRID>",
	Short: "Solve a crossword grid",
	Long: `🐊 Welcome to Crogo, a crossword solver that bites

Examples:

$ crogo "...,...,..." # The grid is a comma-separated list of rows.
[[B A A] [A B B] [B A A]]

$ crogo "A..,B..,C.." # '.' means an empty cell
[[A B A] [B A B] [C H A]]

$ crogo "ALL,...,..." --count 3 # --count allows to get more than one solution
[[A L L] [B A A] [A B B]]
[[A L L] [K A A] [A B B]]
[[A L L] [K A A] [E B B]]

`,
	Args: cobra.ExactArgs(1),
	RunE: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&count, "count", "c", 1, "the desired number of solutions")
	rootCmd.Flags().StringVarP(&solverName, "solver", "s", "logicng", "the desired solver backend. Possible values are: logicng, gini")
}

func run(_ *cobra.Command, args []string) error {
	crossword, errCrossword := crosswordFrom(args[0])
	s, errSolver := solverFrom(solverName)
	if errCrossword != nil || errSolver != nil {
		return errors.Join(errCrossword, errSolver)
	}
	solutions := crossword.SolveWith(s)
	iterateAndPrint(solutions)
	return nil
}

func crosswordFrom(crosswordArg string) (*crogo.Crossword, error) {
	lines := strings.Split(crosswordArg, ",")
	runes := make([][]rune, len(lines))
	for i, line := range lines {
		runes[i] = []rune(line)
	}
	return crogo.NewCrossword(runes, dictionaries.Ukacd())
}

func solverFrom(solverName string) (solver.ConfigurableSolver, error) {
	switch solverName {
	case "logicng":
		return solver.NewLogicNgSolver(), nil
	case "gini":
		return solver.NewGiniSolver(), nil
	default:
		return nil, fmt.Errorf("unknown solver: %s", solverName)
	}
}

func iterateAndPrint(solutions crogo.Solutions) {
	getNextSolution, stop := iter.Pull(solutions)
	defer stop()
	for i := 1; i <= count; i++ {
		nextSolution, found := getNextSolution()
		if !found {
			if i == 1 {
				fmt.Println("No solution found.")
			} else {
				fmt.Println("No more solution.")
			}
			break
		}
		fmt.Printf("%c\n", nextSolution)
	}
}
