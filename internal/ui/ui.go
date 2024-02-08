package ui

import (
	"fmt"
	"strings"
)

var frame = map[string]string{
	"corner": "+",
	"border": "-",
	"pipe":   "|",
}

func createBorder(maxLen int) string {
	const padding = 2

	border := strings.Builder{}
	border.WriteString(frame["corner"])

	for i := 0; i < maxLen+padding; i++ {
		border.WriteString(frame["border"])
	}

	border.WriteString(frame["corner"])

	return border.String()
}

func SingleFrameDisplay(text string, maxLen int) {
	border := createBorder(maxLen)

	centerRow := strings.Builder{}
	centerRow.WriteString(frame["pipe"] + " ")
	centerRow.WriteString(text + " ")
	centerRow.WriteString(frame["pipe"])

	formatted := border + "\n" +
		centerRow.String() + "\n" +
		border

	fmt.Println(formatted)
}

func MultilineFrameDisplay(rows []string, maxLen int) {
    border := createBorder(maxLen)
    fmt.Println(border);

    for _, text := range rows {
        strLen := len(text)
        padding := maxLen - strLen + 1

        line := strings.Builder{}
        line.WriteString(frame["pipe"] + " ")
        line.WriteString(addTrailingWhitespace(text, padding))
        line.WriteString(frame["pipe"])
        fmt.Println(line.String())
    }

    fmt.Println(border);
}

func addTrailingWhitespace(input string, count int) string {
    sb := strings.Builder{}
    sb.WriteString(input)

    for i := 0; i < count; i++ {
        sb.WriteByte(' ')
    }

    return sb.String()
}
