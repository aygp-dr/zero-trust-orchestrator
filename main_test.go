package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInitialModel(t *testing.T) {
	m := initialModel()
	if len(m.policies) != 6 {
		t.Fatalf("expected 6 policies, got %d", len(m.policies))
	}
	if m.cursor != 0 {
		t.Errorf("expected cursor at 0, got %d", m.cursor)
	}
	if m.view != dashboardView {
		t.Errorf("expected dashboard view, got %v", m.view)
	}
}

func TestNavigateDown(t *testing.T) {
	m := initialModel()
	for i := 0; i < 5; i++ {
		updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		m = updated.(model)
	}
	if m.cursor != 5 {
		t.Errorf("expected cursor at 5, got %d", m.cursor)
	}
	// Should not go past last item
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(model)
	if m.cursor != 5 {
		t.Errorf("cursor should stay at 5, got %d", m.cursor)
	}
}

func TestNavigateUp(t *testing.T) {
	m := initialModel()
	// Move down first
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(model)
	if m.cursor != 2 {
		t.Fatalf("expected cursor at 2, got %d", m.cursor)
	}
	// Move up
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(model)
	if m.cursor != 1 {
		t.Errorf("expected cursor at 1, got %d", m.cursor)
	}
}

func TestNavigateUpAtTop(t *testing.T) {
	m := initialModel()
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(model)
	if m.cursor != 0 {
		t.Errorf("cursor should stay at 0, got %d", m.cursor)
	}
}

func TestEnterDetailView(t *testing.T) {
	m := initialModel()
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	if m.view != detailView {
		t.Errorf("expected detail view after enter, got %v", m.view)
	}
}

func TestEscReturnsToList(t *testing.T) {
	m := initialModel()
	// Enter detail view
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	// Press esc
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEscape})
	m = updated.(model)
	if m.view != dashboardView {
		t.Errorf("expected dashboard view after esc, got %v", m.view)
	}
}

func TestQuitFromDetailReturnsToList(t *testing.T) {
	m := initialModel()
	// Enter detail view
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	// Press q - should go back to dashboard, not quit
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m = updated.(model)
	if m.view != dashboardView {
		t.Errorf("expected dashboard view after q in detail, got %v", m.view)
	}
	if cmd != nil {
		t.Error("should not quit from detail view with q")
	}
}

func TestQuitFromDashboard(t *testing.T) {
	m := initialModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Error("expected quit command from dashboard")
	}
}

func TestDashboardViewRender(t *testing.T) {
	m := initialModel()
	out := m.View()
	if out == "" {
		t.Error("dashboard view should not be empty")
	}
	// Check that policy names appear in the output
	for _, p := range m.policies {
		if !containsStr(out, p.Name) {
			t.Errorf("dashboard should contain policy name %q", p.Name)
		}
	}
}

func TestDetailViewRender(t *testing.T) {
	m := initialModel()
	m.view = detailView
	m.cursor = 0
	out := m.View()
	if !containsStr(out, m.policies[0].Name) {
		t.Errorf("detail view should contain selected policy name %q", m.policies[0].Name)
	}
	if !containsStr(out, "Compliance") {
		t.Error("detail view should contain Compliance label")
	}
}

func TestCursorDoesNotMoveInDetailView(t *testing.T) {
	m := initialModel()
	m.view = detailView
	m.cursor = 2
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(model)
	if m.cursor != 2 {
		t.Errorf("cursor should not move in detail view, got %d", m.cursor)
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && searchStr(s, substr)
}

func searchStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
