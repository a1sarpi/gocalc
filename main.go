package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/a1sarpi/gocalc/src/evaluation"
	"github.com/a1sarpi/gocalc/src/tokenizer"
	"github.com/spf13/cobra"
)

var useRadians bool

var rootCmd = &cobra.Command{
	Use:   "calc",
	Short: "Advanced calculator with RPN",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("expression is required")
		}
		input := strings.Join(args, " ")
		return processExpression(input)
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&useRadians, "radians", "r", false, "Use radians for trigonometric functions")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Panic recovered: %v\n", r)
			os.Exit(1)
		}
	}()

	rootCmd.SetOut(os.Stdout)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func processExpression(input string) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Panic recovered in processExpression: %v\n", r)
		}
	}()

	input = strings.TrimPrefix(input, "calc")
	input = strings.TrimSpace(input)

	tokens, err := tokenizer.Tokenize(input)
	if err != nil {
		return fmt.Errorf("tokenization error: %v", err)
	}

	rpn, err := evaluation.ToRPN(tokens)
	if err != nil {
		return fmt.Errorf("RPN conversion error: %v", err)
	}

	result, err := evaluation.Calculate(rpn, useRadians)
	if err != nil {
		return fmt.Errorf("calculation error: %v", err)
	}

	fmt.Fprintf(os.Stdout, "%.15f\n", result)
	return nil
}
