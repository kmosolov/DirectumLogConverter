package DirectumLogConverter

import (
	"bytes"
	"unicode/utf8"
)

func FitWidth(input string, width int) string {
	runeCount := utf8.RuneCountInString(input)
	spaces := width - runeCount
	if spaces <= 0 {
		return input[runeCount-width:]
	}
	buf := bytes.NewBuffer(make([]byte, 0, spaces+len(input)))
	for i := 0; i < spaces; i++ {
		buf.WriteRune(' ')
	}
	buf.WriteString(input)
	return buf.String()
}
