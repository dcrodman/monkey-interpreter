# Monkey

A simple interpreter for the Monkey language built by following along with the book 
[Writing an Interpreter in Go by Thorsten Ball](https://www.amazon.com/Writing-Interpreter-Go-Thorsten-Ball-ebook/dp/B01N2T1VD2/ref=sr_1_1?crid=2FLUXA8PN01BD&keywords=writing+an+interpreter+in+go&qid=1584926025&sprefix=writing+an+inter%2Caps%2C215&sr=8-1).

This project is my own implementation based on the concepts from this book and in many (if not all)
cases will vary significantly from the book's. 

## Usage

The REPL can be found in `cmd/repl` and provides an interactive interpreter like you may be
used to from other dynamic languages. 

    # Assumes GOBIN=$(project)/.bin
    go install cmd/repl
    .bin/repl
    Monkey Version 0.1
    >> let x = 10;
    >> 