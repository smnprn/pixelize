package main

import (
	"log"
	"strconv"

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
	currentPage    *huh.Form
	width          int
	height         int
	style          *Styles
}

type Styles struct {
	BorderColor lipgloss.AdaptiveColor
	Home        lipgloss.Style
}

func DefaultStyle() *Styles {
	s := new(Styles)

	s.BorderColor = lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}
	s.Home = lipgloss.
		NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		Padding(1).
		Width(80)

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
			fileops.Convert(oldFileName, newFileName, format.(imgconv.Format))
		}
	}

	if m.resizeForm.State == huh.StateCompleted {
		fileName := m.resizeForm.GetString("fileName")
		resizeMode := m.resizeForm.GetInt("resizeMode")
		newSizeStr := m.resizeForm.GetString("newSizeStr")
		confirm := m.resizeForm.GetBool("confirm")

		if confirm {
			newSize, err := strconv.ParseFloat(newSizeStr, 64)
			if err != nil {
				log.Fatal(err)
				return m, tea.Quit
			}

			fileops.Resize(fileName, resizeMode, newSize)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.width == 0 {
		return "loading..."
	}

	styledForm := m.style.Home.Render(m.currentPage.View())
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
