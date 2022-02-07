package app

import (
	"fmt"
	"os"
	"tsk/db"
	"tsk/list"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	addInputStyle = lipgloss.NewStyle().Padding(0, 0, 1, 2)

	addCursor = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"})

	addPrompt = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"})
)

type listKeyMap struct {
	togglePagination  key.Binding
	toggleHelpMenu    key.Binding
	addItem           key.Binding
	cancelWhileAdding key.Binding
	acceptWhileAdding key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		addItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
		cancelWhileAdding: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		acceptWhileAdding: key.NewBinding(
			key.WithKeys("enter"),
		),
	}
}

type model struct {
	list          list.Model
	keys          *listKeyMap
	delegate      list.DefaultDelegate
	delegateKeys  *delegateKeyMap
	addInput      textinput.Model
	showAddInput  bool
	width, height int
}

func (m *model) setShowAddInput(b bool) bool {
	m.showAddInput = b
	m.list.SetShowTitle(!b)
	m.list.SetShowFilter(!b)
	m.list.SetFilteringEnabled(!b)
	if b {
		m.addInput.CursorEnd()
		m.addInput.Focus()
	}
	m.updateKeybindings()
	return m.showAddInput
}

func (m *model) setSize(width int, height int) (int, int) {
	m.width = width
	m.height = height
	return m.width, m.height
}

func (m *model) updateKeybindings() {
	switch m.showAddInput {
	case true:
		m.keys.cancelWhileAdding.SetEnabled(true)
		m.keys.addItem.SetEnabled(false)
		m.delegateKeys.remove.SetEnabled(false)
		m.delegateKeys.complete.SetEnabled(false)
		m.list.KeyMap.ShowFullHelp.SetEnabled(false)
		m.list.KeyMap.CloseFullHelp.SetEnabled(false)
		m.list.KeyMap.Quit.SetEnabled(false)
		m.list.KeyMap.CursorUp.SetEnabled(false)
		m.list.KeyMap.CursorDown.SetEnabled(false)
		m.list.KeyMap.PrevPage.SetEnabled(false)
		m.list.KeyMap.NextPage.SetEnabled(false)
		m.list.KeyMap.GoToStart.SetEnabled(false)
		m.list.KeyMap.GoToEnd.SetEnabled(false)
	default:
		m.keys.cancelWhileAdding.SetEnabled(false)
		m.keys.addItem.SetEnabled(true)
		m.delegateKeys.remove.SetEnabled(true)
		m.delegateKeys.complete.SetEnabled(true)
		m.list.KeyMap.ShowFullHelp.SetEnabled(true)
		m.list.KeyMap.CloseFullHelp.SetEnabled(true)
		m.list.KeyMap.Quit.SetEnabled(true)
		m.list.KeyMap.CursorUp.SetEnabled(true)
		m.list.KeyMap.CursorDown.SetEnabled(true)
		m.list.KeyMap.PrevPage.SetEnabled(true)
		m.list.KeyMap.NextPage.SetEnabled(true)
		m.list.KeyMap.GoToStart.SetEnabled(true)
		m.list.KeyMap.GoToEnd.SetEnabled(true)
	}
}

func (m *model) handleAdding(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(msg, m.keys.cancelWhileAdding):
			m.setShowAddInput(false)
			m.keys.addItem.SetEnabled(true)
			m.keys.cancelWhileAdding.SetEnabled(false)
		case key.Matches(msg, m.keys.acceptWhileAdding):
			newItem, err := db.CreateTask(m.addInput.Value())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			m.list.InsertItem(0, newItem)
			m.setShowAddInput(false)
			return m.list.NewStatusMessage(statusMessageStyle("Added " + "\"" + newItem.Value + "\""))
		}
	}

	newAddInput, inputCmd := m.addInput.Update(msg)
	m.addInput = newAddInput
	cmds = append(cmds, inputCmd)

	return tea.Batch(cmds...)
}

func NewModel(tasks []db.Task) model {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	items := make([]list.Item, 0, len(tasks))
	for _, task := range tasks {
		items = append(items, task)
	}

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	taskList := list.New(items, delegate, 0, 0)
	taskList.Title = "Todos"
	taskList.Styles.Title = titleStyle
	taskList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.addItem,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	addInput := textinput.New()
	addInput.Prompt = "Add: "
	addInput.PromptStyle = addPrompt
	addInput.CursorStyle = addCursor
	addInput.Focus()

	return model{
		list:         taskList,
		addInput:     addInput,
		keys:         listKeys,
		delegateKeys: delegateKeys,
		delegate:     delegate,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := appStyle.GetPadding()
		w, h := m.setSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)
		m.list.SetSize(w, h)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		// TODO: update this to actually add a new item to the list
		case key.Matches(msg, m.keys.addItem):
			m.setShowAddInput(true)
			return m, textinput.Blink
		}
	}

	// This will also call our delegate's update function.
	if m.showAddInput {
		cmds = append(cmds, m.handleAdding(msg))
	} else {
		newListModel, cmd := m.list.Update(msg)
		m.list = newListModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var (
		sections []string
	)

	if m.showAddInput {
		v := addInputStyle.Render(m.addInput.View())
		sections = append(sections, v)
		m.list.SetHeight(m.height - lipgloss.Height(v))
	}

	sections = append(sections, m.list.View())

	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}
