package gomq

// Matcher is the interface that wraps MatchString function.
type Matcher interface {

	// MatchString returns true if the input pattern matches.
	MatchString(string) bool
}

// ExactMatcher matches the pattern only if they are equal.
type ExactMatcher string

// MatchString is the implementation of ExactMatcher for Matcher interface.
// It returns true if pattern equals Matcher.
func (em ExactMatcher) MatchString(pattern string) bool {
	return string(em) == pattern
}
