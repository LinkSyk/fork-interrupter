package main

import (
	"interrupter/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
