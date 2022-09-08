package repl

import (
	"bufio"
	"fmt"
	"interrupter/evaluator"
	"interrupter/lexer"
	"interrupter/parser"
	"io"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	prompt := ">> "
	fmt.Fprint(out, "Enter in Fork Language!\n")
	for {
		fmt.Fprint(out, prompt)
		if !scanner.Scan() {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		// p.PrintAllToken(out)
		prog := p.ParseProgram()
		if p.Errors() != nil {
			printParserErrors(out, p.Errors())
			continue
		}
		obj := evaluator.Eval(prog)
		if obj != nil {
			_, _ = io.WriteString(out, obj.Inspect())
			_, _ = io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "parse error: \n")
	for _, err := range errors {
		io.WriteString(out, "\t")
		io.WriteString(out, err)
		io.WriteString(out, "\n")
	}
}
