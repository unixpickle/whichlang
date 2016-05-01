package tokens

// Freqs maps words to their frequencies.
// The frequency for a token X equal to
// the number of occurrences of X, divided
// by the total number of tokens in the
// document.
type Freqs map[string]float64

// Freqs converts word counts into a
// frequency map.
func (c Counts) Freqs() Freqs {
	var totalCount int
	for _, count := range c {
		totalCount += count
	}
	res := Freqs{}
	for word, count := range c {
		res[word] = float64(count) / float64(totalCount)
	}
	return res
}
