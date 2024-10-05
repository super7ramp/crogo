package cmd

import (
	"crogo/pkg/crogo"
	"crogo/pkg/dictionaries"
	"fmt"
	"iter"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// count is the desired number of solutions.
var count int

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "crogo <GRID>",
	Short: "Solve a crossword grid",
	Long: `🐊 Welcome to Crogo, a crossword solver that bites hard.

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
	rootCmd.Flags().IntVarP(&count, "count", "c", 1, "The desired number of solutions")
}

func run(_ *cobra.Command, args []string) error {
	crossword, err := crosswordFrom(args[0])
	if err != nil {
		return err
	}
	solveAndPrintSolutionsOf(crossword)
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

func solveAndPrintSolutionsOf(crossword *crogo.Crossword) {
	solutions := crossword.Solve()
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
