package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/a1sarpi/gocalc/src/evaluation"
	"github.com/a1sarpi/gocalc/src/stack"
	"github.com/a1sarpi/gocalc/src/tokenizer"
	"github.com/spf13/cobra"
)

var (
	useRadians bool
	debug      bool
)

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
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug output")
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

	evaluation.Debug = debug
	stack.Debug = debug

	if debug {
		fmt.Printf("Input expression: %q\n", input)
	}

	// Удаляем имя команды из входной строки, если оно присутствует
	input = strings.TrimPrefix(input, "calc")
	input = strings.TrimSpace(input)

	tokens, err := tokenizer.Tokenize(input)
	if err != nil {
		return fmt.Errorf("tokenization error: %v", err)
	}
	if debug {
		fmt.Printf("Tokens: %v\n", tokens)
	}

	rpn, err := evaluation.ToRPN(tokens)
	if err != nil {
		return fmt.Errorf("RPN conversion error: %v", err)
	}
	if debug {
		fmt.Printf("RPN: %v\n", rpn)
	}

	result, err := evaluation.Calculate(rpn, useRadians)
	if err != nil {
		return fmt.Errorf("calculation error: %v", err)
	}

	if debug {
		fmt.Printf("Result before formatting: %g\n", result)
	}
	fmt.Fprintf(os.Stdout, "%.15f\n", result)
	if debug {
		fmt.Printf("Result after formatting: %.15f\n", result)
	}
	return nil
}
