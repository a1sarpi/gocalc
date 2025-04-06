package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/a1sarpi/gocalc/src/evaluation"
	"github.com/a1sarpi/gocalc/src/tokenizer"
	"github.com/spf13/cobra"
)

var (
	useRadians bool
)

var rootCmd = &cobra.Command{
	Use:   "calc",
	Short: "Advanced calculator with RPN",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: expression is required")
			os.Exit(1)
		}
		input := strings.Join(args, " ")
		processExpression(input)
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&useRadians, "radians", "r", false, "Use radians for trigonometric functions")
}

func main() {
	rootCmd.SetOut(os.Stdout)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func processExpression(input string) {
	tokens, err := tokenizer.Tokenize(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	rpn, err := evaluation.ToRPN(tokens)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	result, err := evaluation.Calculate(rpn, useRadians)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%.6g\n", result)
}
