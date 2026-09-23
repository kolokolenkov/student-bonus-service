package rules

type ScholarshipRule struct{}

func NewScholarshipRule() *ScholarshipRule {
	return &ScholarshipRule{}
}

func (r *ScholarshipRule) Calculate(
	hours float64,
	rate float64,
) float64 {
	return hours * rate
}
