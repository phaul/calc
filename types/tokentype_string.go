// Code generated by "stringer -type TokenType token.go"; DO NOT EDIT.

package types

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[InvalidToken-0]
	_ = x[EOL-1]
	_ = x[IntLit-2]
	_ = x[FloatLit-3]
	_ = x[Name-4]
	_ = x[Sticky-5]
	_ = x[NotSticky-6]
}

const _TokenType_name = "InvalidTokenEOLIntLitFloatLitNameStickyNotSticky"

var _TokenType_index = [...]uint8{0, 12, 15, 21, 29, 33, 39, 48}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
