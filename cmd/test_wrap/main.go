package main

import (
	"fmt"
	"strings"

	"github.com/jaytaylor/html2text"
)

func main() {
	// Create a long paragraph
	longText := strings.Repeat("This is a long sentence that should hopefully not be wrapped by default but we want to check if it inserts newlines after 80 characters or so. ", 10)
	input := fmt.Sprintf("<p>%s</p>", longText)

	text, _ := html2text.FromString(input, html2text.Options{})

	fmt.Printf("Original length: %d\n", len(text))
	fmt.Printf("Contains newlines: %v\n", strings.Contains(text, "\n"))

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if i < 3 {
			fmt.Printf("Line %d length: %d\n", i, len(line))
		}
	}
}
