package policy

import "testing"

func TestComplianceLevel(t *testing.T) {
	tests := []struct {
		name string
		pct  float64
		want string
	}{
		{"high at 100", 100.0, "high"},
		{"high at 90", 90.0, "high"},
		{"high at 95.5", 95.5, "high"},
		{"medium at 89.9", 89.9, "medium"},
		{"medium at 70", 70.0, "medium"},
		{"medium at 75", 75.0, "medium"},
		{"low at 69.9", 69.9, "low"},
		{"low at 0", 0.0, "low"},
		{"low at 45.8", 45.8, "low"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Policy{CompliancePct: tt.pct}
			if got := p.ComplianceLevel(); got != tt.want {
				t.Errorf("ComplianceLevel() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMockPolicies(t *testing.T) {
	policies := MockPolicies()
	if len(policies) != 6 {
		t.Fatalf("expected 6 policies, got %d", len(policies))
	}

	names := map[string]bool{}
	for _, p := range policies {
		names[p.Name] = true
	}

	expected := []string{
		"network-segmentation",
		"identity-verification",
		"device-trust",
		"data-classification",
		"access-logging",
		"encryption-enforcement",
	}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("missing expected policy: %s", name)
		}
	}
}

func TestMockPoliciesStatuses(t *testing.T) {
	policies := MockPolicies()
	statusCount := map[Status]int{}
	for _, p := range policies {
		statusCount[p.Status]++
	}
	if statusCount[Enforced] != 3 {
		t.Errorf("expected 3 enforced policies, got %d", statusCount[Enforced])
	}
	if statusCount[Monitoring] != 2 {
		t.Errorf("expected 2 monitoring policies, got %d", statusCount[Monitoring])
	}
	if statusCount[Disabled] != 1 {
		t.Errorf("expected 1 disabled policy, got %d", statusCount[Disabled])
	}
}

func TestMockPoliciesComplianceRange(t *testing.T) {
	for _, p := range MockPolicies() {
		if p.CompliancePct < 0 || p.CompliancePct > 100 {
			t.Errorf("policy %q has compliance %f out of [0,100] range", p.Name, p.CompliancePct)
		}
	}
}

func TestMockPoliciesViolationsNonNegative(t *testing.T) {
	for _, p := range MockPolicies() {
		if p.Violations24h < 0 {
			t.Errorf("policy %q has negative violations: %d", p.Name, p.Violations24h)
		}
	}
}
