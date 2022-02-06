package app

import (
	"fmt"
	"io"
	"tsk/db"
	"tsk/list"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mergestat/timediff"
	"github.com/muesli/reflow/truncate"
)

var (
	normalTitle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
			Padding(0, 0, 0, 1)
	normalDesc = normalTitle.Copy().
			Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})
	selectedTitle = lipgloss.NewStyle().
			BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
			Padding(0, 0, 0, 1)
	selectedDesc = selectedTitle.Copy().
			Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"})
	dimmedTitle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
			Padding(0, 0, 0, 1)
	dimmedDesc = dimmedTitle.Copy().
			Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	timeAddedStyle    = lipgloss.NewStyle().PaddingLeft(5).Faint(true).Render
)

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.SetSpacing(1)

	d.RenderFunc = func(w io.Writer, m list.Model, index int, listItem list.Item) {
		task, ok := listItem.(db.Task)
		if !ok {
			return
		}

		var (
			title, desc, checkbox string
			matchedRunes          []int
			s                     = &d.Styles
		)

		if i, ok := listItem.(db.Task); ok {
			title = i.Value
			desc = timeAddedStyle(timediff.TimeDiff(task.TimeAdded))
			if i.Completed {
				checkbox = fmt.Sprint("[x]")
			} else {
				checkbox = fmt.Sprint("[ ]")
			}
		} else {
			return
		}

		width := m.Width()
		// Prevent text from exceeding list width
		if width > 0 {
			textwidth := uint(width - normalTitle.GetPaddingLeft() - normalTitle.GetPaddingRight())
			title = truncate.StringWithTail(title, textwidth, "…")
			desc = truncate.StringWithTail(desc, textwidth, "…")
		}

		// Conditions
		var (
			isSelected  = index == m.Index()
			emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
			isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
		)

		if isFiltered && index < len(m.VisibleItems()) {
			// Get indices of matched characters
			matchedRunes = m.MatchesForItem(index)
		}

		if emptyFilter {
			title = dimmedTitle.Render(title)
			desc = dimmedDesc.Render(desc)
			checkbox = dimmedDesc.Render(checkbox)
		} else if isSelected && m.FilterState() != list.Filtering {
			if isFiltered {
				// Highlight matches
				unmatched := selectedTitle.Inline(true)
				matched := unmatched.Copy().Inherit(s.FilterMatch)
				title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
			}
			title = selectedTitle.Render(title)
			desc = selectedTitle.Render(desc)
			checkbox = selectedTitle.Render(checkbox)
		} else {
			if isFiltered {
				// Highlight matches
				unmatched := normalTitle.Inline(true)
				matched := unmatched.Copy().Inherit(s.FilterMatch)
				title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
			}
			title = normalTitle.Render(title)
			desc = normalDesc.Render(desc)
			checkbox = normalDesc.Render(checkbox)
		}

		str := fmt.Sprintf("%s %s\n%s", checkbox, title, desc)
		fmt.Fprintf(w, str)
	}
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		task, ok := m.SelectedItem().(db.Task)
		if !ok {
			return nil
		}
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.complete):
				var err error
				var status string
				if task.Completed {
					task, err = db.MarkActive(task.Key)
					status = ""
				} else {
					task, err = db.MarkDone(task.Key)
					status = statusMessageStyle("Completed: " + fmt.Sprintf("\"%s\"", task.Value) + " 🎉")
				}
				if err != nil {
					m.NewStatusMessage("Error marking task as complete")
				}
				m.SetItem(m.Index(), task)
				return m.NewStatusMessage(status)

			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				return m.NewStatusMessage(statusMessageStyle("Deleted " + fmt.Sprintf("\"%s\"", task.Value)))
			}
		}

		return nil
	}

	help := []key.Binding{keys.complete, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type delegateKeyMap struct {
	complete key.Binding
	remove   key.Binding
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.complete,
		d.remove,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.complete,
			d.remove,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		complete: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "toggle complete"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}
