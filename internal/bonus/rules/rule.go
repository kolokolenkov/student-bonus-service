package rules

type Rule interface {
	Calculate(hours float64, rate float64) float64
}
