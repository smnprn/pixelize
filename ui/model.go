package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

/*
 * Page indexes:
 * 0 - Home
 * 1 - Conversion
 * 2 - Resize
 */

type Model struct {
	pages            []tea.Model
	currentPageIndex int
	completed        bool
	width            int
	height           int
	style            *Styles
	errStatus        error
}

func NewModel() Model {
	style := DefaultStyle()
	m := Model{
		pages: []tea.Model{
			HomePage(),
			ConversionPage(),
			ResizePage(),
		},
		currentPageIndex: 0,
		style:            style,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return m.pages[m.currentPageIndex].Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "ctrl+p":
			m.ChangePageOperation(-1)
		}
	}

	var cmds []tea.Cmd
	_, cmd := m.pages[m.currentPageIndex].Update(msg)
	cmds = append(cmds, cmd)

	if _, ok := m.pages[m.currentPageIndex].(*huh.Form); ok {
		if m.pages[m.currentPageIndex].(*huh.Form).State == huh.StateCompleted {
			m.ChangePageOperation(m.currentPageIndex)
			m.pages[m.currentPageIndex].Init()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var styledForm string
	if m.completed {
		styledForm = m.style.Framed.Render(CreateResultScreen(m))
	} else {
		styledForm = m.style.Framed.Render(m.pages[m.currentPageIndex].View())
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

func (m *Model) ChangePageOperation(pageIndex int) {
	switch pageIndex {
	case 0:
		terminateHomePage(m)
	case 1:
		terminateConversionPage(m)
	case 2:
		terminateResizePage(m)
	case -1:
		resetPages(m)
	}
}
