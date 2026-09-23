package rules

type PracticeRule struct{}

func NewPracticeRule() *PracticeRule {
	return &PracticeRule{}
}

func (r *PracticeRule) Calculate(
	hours float64,
	rate float64,
) float64 {
	return hours * rate
}
