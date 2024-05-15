package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/smnprn/pixelize/fileops"
	"github.com/smnprn/pixelize/ui"
	"github.com/sunshineplan/imgconv"
)

type model struct {
	homeForm       *huh.Form
	conversionForm *huh.Form
	resizeForm     *huh.Form
	currentPage    tea.Model
	width          int
	height         int
	style          *styles
	errStatus      error
}

type styles struct {
	BorderColor lipgloss.AdaptiveColor
	Framed      lipgloss.Style
	Success     lipgloss.Style
	Failure     lipgloss.Style
	Info        lipgloss.Style
}

func DefaultStyle() *styles {
	s := new(styles)

	s.BorderColor = lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}

	s.Framed = lipgloss.
		NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		Padding(1).
		Width(55)

	s.Success = lipgloss.
		NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}).
		Bold(true)

	s.Failure = lipgloss.
		NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"}).
		Bold(true)

	s.Info = lipgloss.
		NewStyle().
		Foreground(lipgloss.Color("#3C3C3C")).Align(lipgloss.Center)

	return s
}

func NewModel() model {
	style := DefaultStyle()
	m := model{
		homeForm:       ui.HomePage(),
		conversionForm: ui.ConversionPage(),
		resizeForm:     ui.ResizePage(),
		style:          style,
	}

	m.currentPage = m.homeForm

	return m
}

func (m model) Init() tea.Cmd {
	return m.homeForm.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd
	currentPage, cmd := m.currentPage.Update(msg)
	if f, ok := currentPage.(*huh.Form); ok {
		m.currentPage = f
		cmds = append(cmds, cmd)
	}

	if m.homeForm.State == huh.StateCompleted {
		operation := m.homeForm.GetInt("operation")

		switch operation {
		case 0:
			m.currentPage = m.conversionForm
			m.currentPage.Init()
		case 1:
			m.currentPage = m.resizeForm
			m.currentPage.Init()
		}
	}

	if m.conversionForm.State == huh.StateCompleted {
		oldFileName := m.conversionForm.GetString("oldFileName")
		newFileName := m.conversionForm.GetString("newFileName")
		format := m.conversionForm.Get("format")
		confirm := m.conversionForm.GetBool("confirm")

		if confirm {
			m.errStatus = fileops.Convert(oldFileName, newFileName, format.(imgconv.Format))
			if m.errStatus != nil {
				m.currentPage = nil
			}
		}

		m.currentPage = nil
	}

	if m.resizeForm.State == huh.StateCompleted {
		fileName := m.resizeForm.GetString("fileName")
		resizeMode := m.resizeForm.GetInt("resizeMode")
		newSizeStr := m.resizeForm.GetString("newSizeStr")
		confirm := m.resizeForm.GetBool("confirm")

		if confirm {
			newSize, err := strconv.ParseFloat(newSizeStr, 64)
			if err != nil {
				m.errStatus = errors.New("could not parse new size")
			}

			m.errStatus = fileops.Resize(fileName, resizeMode, newSize)
			if m.errStatus != nil {
				m.currentPage = nil
			}
		}

		m.currentPage = nil
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.width == 0 {
		return "loading..."
	}

	var styledForm string
	if m.currentPage != nil {
		styledForm = m.style.Framed.Render(m.currentPage.View())
	} else {
		var success string
		if m.errStatus != nil {
			success = m.style.Failure.Render("failure (" + m.errStatus.Error() + ")")
		} else {
			success = m.style.Success.Render("success")
		}

		exit := m.style.Info.Render("You can exit the program using 'esc' or 'ctrl+c'")
		var builder strings.Builder
		fmt.Fprintf(&builder, "Operation status: %s\n", success)
		fmt.Fprintf(&builder, exit)
		styledForm = m.style.Framed.Render(builder.String())
	}

	centeredForm := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		styledForm,
	)

	return centeredForm
}

func main() {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
