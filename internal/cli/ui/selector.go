package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const AllNonManagementText = "all non management accounts"

var (
	gray = lipgloss.Color("240")
	pink = lipgloss.Color("205")

	helpStyle              = lipgloss.NewStyle().Foreground(gray)
	filterPlaceholderStyle = lipgloss.NewStyle().Foreground(gray).Faint(true)
	cursorStyle            = lipgloss.NewStyle().Foreground(pink)

	helpTextMultipleChoice = "↑/↓/←/→: Navigate • Space: Select • Enter: Confirm"
)

type model struct {
	question      string
	choices       []string
	filtered      []string
	filter        string
	cursor        int
	selected      map[string]struct{}
	allChoicesMap map[int]string // Maps original index to choice text
	quit          bool
	pageSize      int
	currentPage   int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quit = true
			return m, tea.Quit
		case "enter":
			if len(m.selected) > 0 {
				return m, tea.Quit
			}
		case "up":
			if m.cursor > 0 {
				m.cursor--
			} else if m.currentPage > 0 {
				m.currentPage--
				m.cursor = m.pageSize - 1
			}
		case "down":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			} else if (m.currentPage+1)*m.pageSize < len(m.filtered) {
				m.currentPage++
				m.cursor = 0
			}
		case "left":
			if m.currentPage > 0 {
				m.currentPage--
				m.cursor = 0
			}
		case "right":
			if (m.currentPage+1)*m.pageSize < len(m.filtered) {
				m.currentPage++
				m.cursor = 0
			}
		case " ":
			if len(m.filtered) > 0 {
				currentChoice := m.filtered[m.currentPage*m.pageSize+m.cursor]
				if _, ok := m.selected[currentChoice]; ok {
					delete(m.selected, currentChoice)
				} else {
					m.selected[currentChoice] = struct{}{}
				}
			}
		case "backspace":
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.updateFilteredChoices()
			}
		default:
			if len(msg.String()) == 1 {
				m.filter += msg.String()
				m.updateFilteredChoices()
			}
		}
	}

	return m, nil
}

func (m *model) updateFilteredChoices() {
	if m.filter == "" {
		m.filtered = m.choices
	} else {
		m.filtered = []string{}
		for _, choice := range m.choices {
			if strings.Contains(strings.ToLower(choice), strings.ToLower(m.filter)) {
				m.filtered = append(m.filtered, choice)
			}
		}
	}

	// Reset pagination when filter changes
	m.currentPage = 0
	m.cursor = 0
}

func (m model) View() string {
	s := fmt.Sprintf("%s:\n\n", m.question)

	filterText := "Filter: "
	if m.filter == "" {
		filterText += filterPlaceholderStyle.Render("type to filter")
	} else {
		filterText += m.filter
	}
	s += fmt.Sprintf("%s\n\n", helpStyle.Render(filterText))

	if len(m.filtered) == 0 {
		s += "No matches found.\n"
	} else {
		start := m.currentPage * m.pageSize
		end := min(start+m.pageSize, len(m.filtered))

		for i, choice := range m.filtered[start:end] {
			cursor := " "
			if m.cursor == i {
				cursor = cursorStyle.Render(">")
			}

			checked := " "
			if _, ok := m.selected[choice]; ok {
				checked = cursorStyle.Render("x")
			}

			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}

		if len(m.filtered) > m.pageSize {
			s += fmt.Sprintf("\n%s\n",
				helpStyle.Render(fmt.Sprintf("Page %d/%d (Showing %d-%d of %d items)",
					m.currentPage+1,
					(len(m.filtered)-1)/m.pageSize+1,
					start+1,
					end,
					len(m.filtered))))
		}
	}

	if len(m.selected) == 0 {
		s += "\n" + helpStyle.Render("Please select at least one item")
	}

	s += "\n" + helpStyle.Render(helpTextMultipleChoice)

	return s
}

func Prompt(question string, choices []string) ([]int, error) {
	m := model{
		question: question,
		choices:  choices,
		filtered: choices,
		filter:   "",
		cursor:   0,
		selected: make(map[string]struct{}),
		allChoicesMap: func() map[int]string {
			m := make(map[int]string)
			for i, choice := range choices {
				m[i] = choice
			}
			return m
		}(),
		quit:        false,
		pageSize:    10,
		currentPage: 0,
	}

	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return nil, err
	}

	finalModel, ok := result.(model)
	if !ok {
		return nil, fmt.Errorf("failed to parse the final model")
	}

	if finalModel.quit {
		return nil, fmt.Errorf("selector interrupted")
	}

	selectedIndexes := []int{}
	for choice := range finalModel.selected {
		for i, c := range finalModel.allChoicesMap {
			if c == choice {
				selectedIndexes = append(selectedIndexes, i)
				break
			}
		}
	}

	return selectedIndexes, nil
}
