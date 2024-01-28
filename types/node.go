package types

import (
	"fmt"
	"strings"

	"github.com/phaul/calc/types/token"
)

type Node struct {
	Token    token.Type
	Children []Node
}

func (n Node) PrettyPrint() { recurse(0, n) }

func recurse(depth int, n Node) {
	fmt.Printf("%s %v\n", strings.Repeat(" ", 3*depth), n.Token)
	for _, c := range n.Children {
		recurse(depth+1, c)
	}
}
