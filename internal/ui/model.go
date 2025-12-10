package ui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jaytaylor/html2text"
	"github.com/pkg/browser"
	"github.com/slatkin/goflux/pkg/config"
	"github.com/slatkin/goflux/pkg/miniflux"
)

type State int

const (
	StateLoading State = iota
	StateList
	StateReading
	StateError
)

type Model struct {
	Client *miniflux.Client
	Config config.Config

	State    State
	Entries  []miniflux.FeedEntry
	Cursor   int
	Offset   int
	Selected *miniflux.FeedEntry

	Viewport viewport.Model
	Help     help.Model

	Err error
}

func NewModel(cfg config.Config) Model {
	client := miniflux.NewClient(cfg.ServerUrl, cfg.ApiKey, cfg.AllowInvalidCerts)
	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().Padding(1, 2)

	return Model{
		Client:   client,
		Config:   cfg,
		State:    StateLoading,
		Viewport: vp,
		Help:     help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.fetchUnreadEntries,
	)
}

func (m Model) fetchUnreadEntries() tea.Msg {
	entries, err := m.Client.GetUnreadEntries(50, 0) // TODO: Pagination
	if err != nil {
		return ErrorMsg(err)
	}
	return EntriesMsg(entries)
}

func (m Model) fetchContent(entryID int) tea.Cmd {
	return func() tea.Msg {
		content, err := m.Client.FetchOriginalContent(entryID)
		if err != nil {
			return ErrorMsg(err)
		}
		return EntryContentMsg{EntryID: entryID, Content: content}
	}
}

func (m Model) toggleReadStatus(entryID int, currentStatus miniflux.ReadStatus) tea.Cmd {
	return func() tea.Msg {
		newStatus := currentStatus.Toggle()
		err := m.Client.ChangeEntryReadStatus([]int{entryID}, newStatus)
		if err != nil {
			return ErrorMsg(err)
		}
		return ActionDoneMsg{Action: "read_toggle", Success: true}
	}
}

func (m Model) markAsRead(entryID int) tea.Cmd {
	return func() tea.Msg {
		err := m.Client.ChangeEntryReadStatus([]int{entryID}, miniflux.ReadStatusRead)
		if err != nil {
			return ErrorMsg(err)
		}
		return ActionDoneMsg{Action: "mark_read", Success: true}
	}
}

func (m Model) saveEntry(entryID int) tea.Cmd {
	return func() tea.Msg {
		err := m.Client.SaveEntry(entryID)
		if err != nil {
			return ErrorMsg(err)
		}
		return ActionDoneMsg{Action: "save_entry", Success: true}
	}
}

func openUrl(url string) tea.Cmd {
	return func() tea.Msg {
		_ = browser.OpenURL(url)
		return nil
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case keyMatches(msg, Keys.Quit):
			return m, tea.Quit
		case keyMatches(msg, Keys.Back):
			if m.State == StateReading {
				m.State = StateList
				m.Viewport.SetContent("") // Clear content to save memory? Or keep it.
				// m.Viewport.GotoTop() // Reset position?
			}
			return m, nil
		}

		switch m.State {
		case StateList:
			switch {
			case keyMatches(msg, Keys.Up):
				if m.Cursor > 0 {
					m.Cursor--
				}
			case keyMatches(msg, Keys.Down):
				if m.Cursor < len(m.Entries)-1 {
					m.Cursor++
				}
			case keyMatches(msg, Keys.Enter):
				if len(m.Entries) > 0 {
					m.Selected = &m.Entries[m.Cursor]
					m.State = StateReading

					// Auto-mark as read if unread
					if m.Selected.Status == miniflux.ReadStatusUnread {
						m.Selected.Status = miniflux.ReadStatusRead
						m.Entries[m.Cursor].Status = miniflux.ReadStatusRead
						m.Viewport.SetContent(renderEntryContent(m.Selected, m.Viewport.Width))
						m.Viewport.GotoTop()
						return m, m.markAsRead(m.Selected.ID)
					}

					// Load content if not present or just show what we have
					// Miniflux usually sends content in list, but we might want original
					m.Viewport.SetContent(renderEntryContent(m.Selected, m.Viewport.Width))
					m.Viewport.GotoTop()
				}
			case keyMatches(msg, Keys.Refresh):
				m.State = StateLoading
				return m, m.fetchUnreadEntries
			case keyMatches(msg, Keys.ToggleReadList):
				if len(m.Entries) > 0 {
					entry := m.Entries[m.Cursor]
					// Update locally immediately for UI responsiveness
					if entry.Status == miniflux.ReadStatusRead {
						m.Entries[m.Cursor].Status = miniflux.ReadStatusUnread
					} else {
						m.Entries[m.Cursor].Status = miniflux.ReadStatusRead
					}
					// Send API request
					return m, m.toggleReadStatus(entry.ID, entry.Status)
				}
			case keyMatches(msg, Keys.Save):
				if len(m.Entries) > 0 {
					entry := m.Entries[m.Cursor]
					return m, m.saveEntry(entry.ID)
				}
			case keyMatches(msg, Keys.OpenBrowser):
				if len(m.Entries) > 0 {
					entry := m.Entries[m.Cursor]
					return m, openUrl(entry.URL)
				}
			}
		case StateReading:
			// Handle component specific keys first if needed, or global keys above
			switch {
			case keyMatches(msg, Keys.ToggleRead):
				if m.Selected != nil {
					entry := m.Selected
					// We need to find it in the list to update it too?
					// Simplification: just update selected and send request.
					if entry.Status == miniflux.ReadStatusRead {
						entry.Status = miniflux.ReadStatusUnread
					} else {
						entry.Status = miniflux.ReadStatusRead
					}
					// Also update in list if possible
					if m.Cursor >= 0 && m.Cursor < len(m.Entries) && m.Entries[m.Cursor].ID == entry.ID {
						m.Entries[m.Cursor].Status = entry.Status
					}

					return m, m.toggleReadStatus(entry.ID, entry.Status)
				}
			case keyMatches(msg, Keys.Save):
				if m.Selected != nil {
					return m, m.saveEntry(m.Selected.ID)
				}
			case keyMatches(msg, Keys.OpenBrowser):
				if m.Selected != nil {
					return m, openUrl(m.Selected.URL)
				}
			}

			// Forward other keys to viewport (scrolling)
			m.Viewport, cmd = m.Viewport.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		m.Viewport.Width = msg.Width
		m.Viewport.Height = msg.Height - 1 // Leave room for status bar

		if m.State == StateReading && m.Selected != nil {
			// Re-render content with new width
			m.Viewport.SetContent(renderEntryContent(m.Selected, msg.Width))
		}

	case EntriesMsg:
		m.Entries = msg
		m.State = StateList
		m.Cursor = 0

	case ErrorMsg:
		m.Err = msg
		m.State = StateError

	case EntryContentMsg:
		// Handle original content fetching if we implement that feature fully
		// For now we rely on the feed content
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.State {
	case StateLoading:
		return "Loading..."
	case StateError:
		return fmt.Sprintf("Error: %v", m.Err)
	case StateReading:
		return m.Viewport.View()
	case StateList:
		return m.viewList()
	}
	return ""
}

func (m Model) viewList() string {
	var s strings.Builder
	s.WriteString(StyleTitle.Render("Miniflux Feeds") + "\n\n")

	// Simple windowing for long lists
	start := 0
	end := len(m.Entries)
	height := 20 // Arbitrary default, usually bound to WindowSize
	if m.Viewport.Height > 0 {
		height = m.Viewport.Height - 4 // Title + padding
	}

	if m.Cursor >= height {
		start = m.Cursor - height + 1
	}
	if end > start+height {
		end = start + height
	}

	for i := start; i < end; i++ {
		entry := m.Entries[i]
		cursor := " "
		style := StyleBase

		if m.Cursor == i {
			cursor = ">"
			style = StyleSelected
		}

		if entry.Status == miniflux.ReadStatusUnread {
			// style = style.Foreground(lipgloss.Color(m.Config.Theme.UnreadColor))
			// Simpler handling of theme for now
			if m.Cursor != i {
				style = StyleStatusUnread
			}
		} else {
			if m.Cursor != i {
				style = StyleStatusRead
			}
		}

		line := fmt.Sprintf("%s %s", cursor, truncate(entry.Title, 80))
		s.WriteString(style.Render(line) + "\n")
	}

	return s.String()
}

func renderEntryContent(entry *miniflux.FeedEntry, width int) string {
	content := entry.Content
	text, err := html2text.FromString(content, html2text.Options{PrettyTables: true})
	if err == nil {
		content = text
	}

	// Extract links to footnotes
	// html2text output format: "text ( url )"
	// Regex to match ( url )
	re := regexp.MustCompile(`\( (https?://[^)]+) \)`)

	var links []string
	content = re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) > 1 {
			links = append(links, parts[1])
			return fmt.Sprintf("[%d]", len(links))
		}
		return match
	})

	// Append footnotes
	if len(links) > 0 {
		content += "\n\n"
		for i, link := range links {
			content += fmt.Sprintf("[%d] %s\n", i+1, link)
		}
	}

	markdown := fmt.Sprintf("# %s\n\n**%s**\n\n%s",
		entry.Title,
		entry.Feed.Title,
		content,
	)

	// Create a style request wrapping to the effective viewport width
	// Padding is (1, 2) which means 2 on left + 2 on right = 4 total horizontal padding
	wrapWidth := width - 4
	if wrapWidth < 20 {
		wrapWidth = 20
	}

	return lipgloss.NewStyle().Width(wrapWidth).Render(markdown)
}

func keyMatches(k tea.KeyMsg, b key.Binding) bool {
	return key.Matches(k, b)
}

func truncate(s string, l int) string {
	if len(s) > l {
		return s[:l] + "..."
	}
	return s
}
