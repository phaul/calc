package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/phaul/calc/evaluator"
	"github.com/phaul/calc/parser"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	for {
		input, _ := r.ReadString('\n')
		t, err := parser.Parse(input)
		if len(t) > 0 {
			t[0].PrettyPrint()
		}
    fmt.Println("> ", evaluator.Evaluate(t[0]))
		if err != nil {
			fmt.Println(err)
		}
	}
}
