package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const AllNonManagementText = "all non management accounts"

var (
	gray = lipgloss.Color("240")
	pink = lipgloss.Color("205")

	helpStyle              = lipgloss.NewStyle().Foreground(gray)
	filterPlaceholderStyle = lipgloss.NewStyle().Foreground(gray).Faint(true)
	cursorStyle            = lipgloss.NewStyle().Foreground(pink)

	helpTextMultipleChoice = "↑/↓/←/→: Navigate • Space: Select • Enter: Confirm"
	helpTextSingleChoice   = "↑/↓/←/→: Navigate • Enter: Select"
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
	single        bool // when true, enter selects the current cursor item and quits
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quit = true
			return m, tea.Quit
		case "enter":
			if m.single {
				if len(m.filtered) > 0 {
					chosen := m.filtered[m.currentPage*m.pageSize+m.cursor]
					m.selected = map[string]struct{}{chosen: {}}
					return m, tea.Quit
				}
				return m, nil
			}
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
		case "space":
			if m.single {
				return m, nil
			}
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

func (m model) View() tea.View {
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

			if m.single {
				s += fmt.Sprintf("%s %s\n", cursor, choice)
				continue
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

	if !m.single && len(m.selected) == 0 {
		s += "\n" + helpStyle.Render("Please select at least one item")
	}

	helpText := helpTextMultipleChoice
	if m.single {
		helpText = helpTextSingleChoice
	}
	s += "\n" + helpStyle.Render(helpText)

	return tea.NewView(s)
}

// Prompt shows a multi-select TUI and returns the original indexes of the
// chosen items.
func Prompt(question string, choices []string) ([]int, error) {
	return runPrompt(question, choices, false)
}

// PromptSingle shows a single-select TUI and returns the original index of the
// chosen item, or -1 if the user quit without selecting.
func PromptSingle(question string, choices []string) (int, error) {
	indexes, err := runPrompt(question, choices, true)
	if err != nil {
		return -1, err
	}
	if len(indexes) == 0 {
		return -1, nil
	}
	return indexes[0], nil
}

func runPrompt(question string, choices []string, single bool) ([]int, error) {
	allChoicesMap := make(map[int]string, len(choices))
	for i, choice := range choices {
		allChoicesMap[i] = choice
	}

	m := model{
		question:      question,
		choices:       choices,
		filtered:      choices,
		selected:      make(map[string]struct{}),
		allChoicesMap: allChoicesMap,
		pageSize:      10,
		single:        single,
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
