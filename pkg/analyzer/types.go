package analyzer

type AnomalySeverity string

const (
	SeverityCritical AnomalySeverity = "critical"
	SeverityHigh     AnomalySeverity = "high"
	SeverityMedium   AnomalySeverity = "medium"
	SeverityLow      AnomalySeverity = "low"
	SeverityInfo     AnomalySeverity = "info"
)

type Anomaly struct {
	Severity    AnomalySeverity `json:"severity"`
	Title       string          `json:"title"`
	Message     string          `json:"message"`
	Remediation string          `json:"remediation,omitempty"`
	RuleID      string          `json:"ruleId"`
	DocURL      string          `json:"docUrl,omitempty"`
}

type ResourceAnalysis struct {
	Anomalies []Anomaly `json:"anomalies"`
	Summary   string    `json:"summary,omitempty"`
	Score     int       `json:"score,omitempty"`
}
