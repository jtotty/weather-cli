package ui

import (
	"fmt"
	"strings"
)

func CreateBorder(maxLen int) string {
	border := strings.Builder{}
	for i := 0; i < maxLen; i++ {
		border.WriteString("-")
	}

	return border.String()
}

func Spacer() {
    fmt.Println("")
}
