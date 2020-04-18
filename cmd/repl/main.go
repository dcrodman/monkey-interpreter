package main

import (
	"bufio"
	"fmt"
	"io"
	"monkey-interpreter/lexer"
	"monkey-interpreter/parser"
	"os"
)

func main() {
	Start(os.Stdin, os.Stdout)
}

func Start(in io.Reader, out io.Writer) {
	stdinReader := bufio.NewScanner(in)

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

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}
