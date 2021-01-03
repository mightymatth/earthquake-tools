package main

import "fmt"

func SourceLinkHTML(sourceType, ID string) string {
	switch SourceType(sourceType) {
	case emsc:
		return fmt.Sprintf(
			`<a href="https://www.emsc-csem.org/Earthquake/earthquake.php?id=%s">%s/%s</a>`,
			ID, sourceType, ID)
	default:
		return fmt.Sprintf(`<code>%s/%s</code>`, sourceType, ID)
	}
}

type SourceType string

const (
	emsc SourceType = "EMSC-RTS"
)
