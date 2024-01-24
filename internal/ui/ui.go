package ui

import "fmt"

var frame = map[string]string{
    "corner": "+",
    "border": "-",
}

func FrameDisplay(text string) {
    fmt.Println(text)
}
