package ui

import (
	"errors"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/smnprn/pixelize/fileops"
	"github.com/sunshineplan/imgconv"
)

/*
* Page indexes:
* 0 - Home
* 1 - Conversion
* 2 - Resize
 */

type model struct {
	pages            []tea.Model
	currentPageIndex int
	completed        bool
	width            int
	height           int
	style            *Styles
	errStatus        error
}

func NewModel() model {
	style := DefaultStyle()
	m := model{
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

func (m model) Init() tea.Cmd {
	return m.pages[m.currentPageIndex].Init()
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

func (m model) View() string {
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

func (m *model) ChangePageOperation(pageIndex int) error {
	switch pageIndex {
	case 0:
		homeForm := m.pages[m.currentPageIndex].(*huh.Form)
		operation := homeForm.GetInt("operation")

		switch operation {
		case 0:
			m.currentPageIndex = 1
		case 1:
			m.currentPageIndex = 2
		}
	case 1:
		conversionForm := m.pages[m.currentPageIndex].(*huh.Form)

		oldFileName := conversionForm.GetString("oldFileName")
		newFileName := conversionForm.GetString("newFileName")
		format := conversionForm.Get("format")
		confirm := conversionForm.GetBool("confirm")

		if confirm {
			m.errStatus = fileops.Convert(oldFileName, newFileName, format.(imgconv.Format))
			if m.errStatus != nil {
				m.completed = true
			}
		}

		m.completed = true
	case 2:
		resizeForm := m.pages[m.currentPageIndex].(*huh.Form)

		fileName := resizeForm.GetString("fileName")
		resizeMode := resizeForm.GetInt("resizeMode")
		newSizeStr := resizeForm.GetString("newSizeStr")
		confirm := resizeForm.GetBool("confirm")

		if confirm {
			newSize, err := strconv.ParseFloat(newSizeStr, 64)
			if err != nil {
				m.errStatus = errors.New("could not parse new size")
			}

			m.errStatus = fileops.Resize(fileName, resizeMode, newSize)
			if m.errStatus != nil {
				m.completed = true
			}
		}

		m.completed = true
	}

	return nil
}
