package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aygp-dr/zero-trust-orchestrator/internal/policy"
)

// Styles
var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("252")).Underline(true)
	cursorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))

	statusEnforced   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	statusMonitoring = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	statusDisabled   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	complianceHigh = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	complianceMed  = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	complianceLow  = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	detailLabel = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("252"))
	detailValue = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	detailBox   = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(1, 2)
)

type view int

const (
	dashboardView view = iota
	detailView
)

type model struct {
	policies []policy.Policy
	cursor   int
	view     view
}

func initialModel() model {
	return model{
		policies: policy.MockPolicies(),
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.view == detailView {
				m.view = dashboardView
				return m, nil
			}
			return m, tea.Quit
		case "j", "down":
			if m.view == dashboardView && m.cursor < len(m.policies)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.view == dashboardView && m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			if m.view == dashboardView {
				m.view = detailView
			}
		case "esc", "backspace":
			m.view = dashboardView
		case "?":
			// Could toggle help; for now help is always visible
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.view {
	case detailView:
		return m.renderDetail()
	default:
		return m.renderDashboard()
	}
}

func (m model) renderDashboard() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Zero Trust Policy Dashboard"))
	b.WriteString("\n\n")

	// Table header
	b.WriteString(fmt.Sprintf("   %-26s %-14s %-12s %-22s %s\n",
		headerStyle.Render("Policy"),
		headerStyle.Render("Status"),
		headerStyle.Render("Compliance"),
		headerStyle.Render("Last Audit"),
		headerStyle.Render("Violations/24h"),
	))

	// Table rows
	for i, p := range m.policies {
		prefix := "  "
		if i == m.cursor {
			prefix = cursorStyle.Render("> ")
		}

		name := fmt.Sprintf("%-24s", p.Name)
		status := renderStatus(p.Status)
		compliance := renderCompliance(p.CompliancePct)
		audit := fmt.Sprintf("%-20s", formatAuditTime(p.LastAudit))
		violations := fmt.Sprintf("%d", p.Violations24h)

		b.WriteString(fmt.Sprintf("%s %-24s %-23s %-21s %-20s %s\n",
			prefix, name, status, compliance, audit, violations))
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("j/k: navigate  enter: details  q: quit"))
	return b.String()
}

func (m model) renderDetail() string {
	p := m.policies[m.cursor]

	var b strings.Builder
	b.WriteString(titleStyle.Render(fmt.Sprintf("Policy: %s", p.Name)))
	b.WriteString("\n\n")

	rows := []struct{ label, value string }{
		{"Status", string(p.Status)},
		{"Compliance", fmt.Sprintf("%.1f%%", p.CompliancePct)},
		{"Compliance Level", p.ComplianceLevel()},
		{"Last Audit", p.LastAudit.Format("2006-01-02 15:04:05")},
		{"Violations (24h)", fmt.Sprintf("%d", p.Violations24h)},
	}

	var content strings.Builder
	for _, r := range rows {
		label := detailLabel.Render(fmt.Sprintf("%-20s", r.label))
		value := detailValue.Render(r.value)
		// Apply color to status and compliance values
		switch r.label {
		case "Status":
			value = renderStatus(policy.Status(r.value))
		case "Compliance":
			value = renderCompliance(p.CompliancePct)
		case "Compliance Level":
			value = renderComplianceLevel(p.ComplianceLevel())
		}
		content.WriteString(fmt.Sprintf("%s  %s\n", label, value))
	}

	b.WriteString(detailBox.Render(content.String()))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("esc/backspace: back  q: back"))
	return b.String()
}

func renderStatus(s policy.Status) string {
	switch s {
	case policy.Enforced:
		return statusEnforced.Render(fmt.Sprintf("%-12s", s))
	case policy.Monitoring:
		return statusMonitoring.Render(fmt.Sprintf("%-12s", s))
	case policy.Disabled:
		return statusDisabled.Render(fmt.Sprintf("%-12s", s))
	default:
		return fmt.Sprintf("%-12s", s)
	}
}

func renderCompliance(pct float64) string {
	text := fmt.Sprintf("%-10s", fmt.Sprintf("%.1f%%", pct))
	switch {
	case pct >= 90:
		return complianceHigh.Render(text)
	case pct >= 70:
		return complianceMed.Render(text)
	default:
		return complianceLow.Render(text)
	}
}

func renderComplianceLevel(level string) string {
	switch level {
	case "high":
		return complianceHigh.Render(level)
	case "medium":
		return complianceMed.Render(level)
	default:
		return complianceLow.Render(level)
	}
}

func formatAuditTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

func main() {
	jsonFlag := flag.Bool("json", false, "Output policies as JSON and exit")
	flag.Parse()

	if *jsonFlag {
		policies := policy.MockPolicies()
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(policies); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
