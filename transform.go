package DirectumLogConverter

import (
	"bytes"
	"unicode/utf8"
)

func FitWidth(input string, width int) string {
	spaces := width - utf8.RuneCountInString(input)
	if spaces <= 0 {
		return input[:width]
	}
	buf := bytes.NewBuffer(make([]byte, 0, spaces+len(input)))
	for i := 0; i < spaces; i++ {
		buf.WriteRune(' ')
	}
	buf.WriteString(input)
	return buf.String()
}
