package main

import (
	"bufio"
	"fmt"
	"io"
	"monkey-interpreter/evaluator"
	"monkey-interpreter/lexer"
	"monkey-interpreter/object"
	"monkey-interpreter/parser"
	"os"
)

func main() {
	Start(os.Stdin, os.Stdout)
}

func Start(in io.Reader, out io.Writer) {
	stdinReader := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Print(">> ")

		if !stdinReader.Scan() {
			return
		}

		l := lexer.New(stdinReader.Text())
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			for _, msg := range p.Errors() {
				io.WriteString(out, "\t"+msg.Error()+"\n")
			}
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
