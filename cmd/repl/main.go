package main

import (
	"bufio"
	"fmt"
	"monkey-interpreter/lexer"
	"os"
)

func main() {
	stdinReader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">> ")

		if !stdinReader.Scan() {
			return
		}

		tokenizer := lexer.New(stdinReader.Text())
		t := tokenizer.NextToken()

		for ; t.Type != lexer.EOF; t = tokenizer.NextToken() {
			fmt.Printf("%#v\n", t)
		}
		// Strip the EOF.
		tokenizer.NextToken()
	}
}
