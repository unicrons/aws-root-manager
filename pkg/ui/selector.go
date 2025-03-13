package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const AllNonManagementText = "all non management accounts"

type model struct {
	question      string
	choices       []string
	filtered      []string
	filter        string
	cursor        int
	selected      map[string]struct{}
	allChoicesMap map[int]string // Maps original index to choice text
	quit          bool
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
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		case " ":
			if len(m.filtered) > 0 {
				currentChoice := m.filtered[m.cursor]
				if _, ok := m.selected[currentChoice]; ok {
					delete(m.selected, currentChoice)
				} else {
					m.selected[currentChoice] = struct{}{}
				}
			}
		case "enter":
			return m, tea.Quit
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
	if len(m.filtered) == 0 {
		m.cursor = 0
	} else if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
}

func (m model) View() string {
	s := fmt.Sprintf("%s: [Use arrows to move, space to select, enter to confirm, type to filter]\n\n", m.question)
	s += fmt.Sprintf("Filter: %s\n\n", m.filter)

	if len(m.filtered) == 0 {
		s += "No matches found.\n"
	} else {
		for i, choice := range m.filtered {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			checked := " "
			if _, ok := m.selected[choice]; ok {
				checked = "x"
			}

			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}
	}

	s += "\n"

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
		quit: false,
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
