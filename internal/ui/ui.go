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

func createBorder(maxLength int) strings.Builder {
	const padding = 2

	border := strings.Builder{}
	border.WriteString(frame["corner"])

	for i := 0; i < maxLength+padding; i++ {
		border.WriteString(frame["border"])
	}

	border.WriteString(frame["corner"])

	return border
}

func SingleFrameDisplay(text string, maxLength int) {
	border := createBorder(maxLength)

	centerRow := strings.Builder{}
	centerRow.WriteString(frame["pipe"] + " ")
	centerRow.WriteString(text + " ")
	centerRow.WriteString(frame["pipe"])

	formatted := border.String() + "\n" +
		centerRow.String() + "\n" +
		border.String()

	fmt.Println(formatted)
}

func MultilineFrameDisplay(text []string, maxLength int) {
    // border := createBorder(maxLength)
}
