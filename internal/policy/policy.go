package policy

import "time"

// Status represents the enforcement state of a policy.
type Status string

const (
	Enforced   Status = "enforced"
	Monitoring Status = "monitoring"
	Disabled   Status = "disabled"
)

// Policy represents a zero-trust security policy.
type Policy struct {
	Name          string
	Status        Status
	CompliancePct float64
	LastAudit     time.Time
	Violations24h int
}

// ComplianceLevel returns "high" (>=90), "medium" (>=70), or "low" (<70).
func (p Policy) ComplianceLevel() string {
	switch {
	case p.CompliancePct >= 90:
		return "high"
	case p.CompliancePct >= 70:
		return "medium"
	default:
		return "low"
	}
}

// MockPolicies returns the six predefined security policies.
func MockPolicies() []Policy {
	now := time.Now()
	return []Policy{
		{
			Name:          "network-segmentation",
			Status:        Enforced,
			CompliancePct: 94.2,
			LastAudit:     now.Add(-2 * time.Hour),
			Violations24h: 12,
		},
		{
			Name:          "identity-verification",
			Status:        Enforced,
			CompliancePct: 98.7,
			LastAudit:     now.Add(-30 * time.Minute),
			Violations24h: 3,
		},
		{
			Name:          "device-trust",
			Status:        Monitoring,
			CompliancePct: 76.5,
			LastAudit:     now.Add(-6 * time.Hour),
			Violations24h: 28,
		},
		{
			Name:          "data-classification",
			Status:        Enforced,
			CompliancePct: 89.1,
			LastAudit:     now.Add(-1 * time.Hour),
			Violations24h: 7,
		},
		{
			Name:          "access-logging",
			Status:        Monitoring,
			CompliancePct: 82.3,
			LastAudit:     now.Add(-4 * time.Hour),
			Violations24h: 15,
		},
		{
			Name:          "encryption-enforcement",
			Status:        Disabled,
			CompliancePct: 45.8,
			LastAudit:     now.Add(-24 * time.Hour),
			Violations24h: 52,
		},
	}
}
