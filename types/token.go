package types

import (
	"fmt"

	"github.com/phaul/calc/combinator"
)

type TokenType int

const (
	// InvalidToken should not be produced by the lexer, however the parser uses it in compound AST nodes
	InvalidToken = TokenType(iota)
	EOL          // end of line
	IntLit       // integer literal
	FloatLit     // float literal
	VarName      // variable name
	Op           // one of +, -, *, /, =
	Paren        // one of (, )
)

// Token as produced by the lexer
type Token struct {
	Value string    // the slice into the input containing the token value
	Type  TokenType // the type of the token
}

func (t Token) String() string {
	if t.Type == EOL {
		return fmt.Sprintf("<%v>", t.Type)
	}
	return fmt.Sprintf("<\"%v\" %v>", t.Value, t.Type)
}

// fulfills the combinator.Token interface
func (t Token) Node() combinator.Node {
	return Node{Token: t}
}
