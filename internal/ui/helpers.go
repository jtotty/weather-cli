package ui

import (
	"fmt"
	"strings"
)

func CreateBorder(maxLen int) string {
	border := strings.Builder{}
	for range maxLen {
		border.WriteString("-")
	}

	return border.String()
}

func Spacer() {
	fmt.Print("\n\n")
}
