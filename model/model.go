package model

import (
	"bufio"
	"log"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zt64/logcat-tui/logcat"
)

var (
	verboseStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	debugStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffcc00"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	fatalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Bold(true)
)

type model struct {
	channel        chan lineMsg
	lines          []string
	terminalHeight int
	scrollOffset   int
	lineLimit      int
	autoscroll     bool
}

func InitializeModel() model {
	c := make(chan lineMsg)
	go startLogcat(c)
	return model{
		channel:    c,
		lineLimit:  10000,
		autoscroll: true,
	}
}

func waitForActivity(sub chan lineMsg) tea.Cmd {
	return func() tea.Msg {
		return lineMsg(<-sub)
	}
}

func (m model) Init() tea.Cmd {
	tea.SetWindowTitle("Bubble Tea Example")
	return tea.Batch(
		waitForActivity(m.channel),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case lineMsg:
		m.lines = append(m.lines, colorize(msg.Message, msg.Priority))
		if len(m.lines) > m.lineLimit {
			drop := len(m.lines) - m.lineLimit
			m.lines = m.lines[drop:] // Keep only the last 'lineLimit' lines
			m.scrollOffset -= drop
		}
		if m.autoscroll && m.scrollOffset < len(m.lines)-m.terminalHeight {
			m.scrollOffset = len(m.lines) - m.terminalHeight
		}
		if m.scrollOffset < 0 {
			m.scrollOffset = 0
		}
		return m, waitForActivity(m.channel)
	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEscape:
			return m, tea.Quit
		case tea.KeyCtrlL:
			m.scrollOffset = 0
			m.autoscroll = true
		case tea.KeyUp:
			m.autoscroll = false
			if m.scrollOffset > 0 {
				m.scrollOffset--
			}
		case tea.KeyDown:
			if m.scrollOffset < len(m.lines)-m.terminalHeight {
				m.scrollOffset++
			}
			if m.scrollOffset == len(m.lines)-m.terminalHeight {
				m.autoscroll = true
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	start := m.scrollOffset
	end := start + m.terminalHeight
	if end >= len(m.lines) {
		end = len(m.lines) - 1
	}

	visibleLines := m.lines[start : end+1]
	return strings.Join(visibleLines, "\n")
}

func colorize(s string, p logcat.Priority) string {
	var style lipgloss.Style
	switch p {
	case logcat.PriorityVerbose:
		style = verboseStyle
	case logcat.PriorityDebug:
		style = debugStyle
	case logcat.PriorityInfo:
		style = infoStyle
	case logcat.PriorityWarn:
		style = warnStyle
	case logcat.PriorityError:
		style = errorStyle
	case logcat.PriorityFatal:
		style = fatalStyle
	default:
		return s
	}

	return style.Render(s)
}

func startLogcat(lines chan<- lineMsg) {
	cmd := exec.Command("adb", "logcat", "-v", "epoch")
	reader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	buf := bufio.NewReader(reader)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			close(lines)
			return
		}
		lineStr := string(line)
		if lineStr == "--------- beginning of system" {
			continue
		}
		lines <- parseLine(lineStr)
	}
}
