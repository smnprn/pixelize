package ui

import (
	"fmt"
	"strings"
)

func CreateResultScreen(m Model) string {
	var builder strings.Builder

	var success string
	if m.errStatus != nil {
		success = m.style.Failure.Render("failure")
	} else {
		success = m.style.Success.Render("success")
	}

	exit := m.style.Info.Render("Use 'ctrl+c' to exit or 'ctrl+p' to start again")
	fmt.Fprintf(&builder, "Operation status: %s\n", success)
	if m.errStatus != nil {
		fmt.Fprintf(&builder, "%s\n", m.errStatus)
	}

	fmt.Fprint(&builder, exit)
	return builder.String()
}
