package main

import "sort"

type Rating struct {
	Correct int
	Total   int
}

func (r *Rating) Frac() float64 {
	return float64(r.Correct) / float64(r.Total)
}

type OverallRating struct {
	Rating
	LangRatings []*LangRating
}

func NewOverallRating(correct, total int, l []*LangRating) *OverallRating {
	sorter := ratingSorter(make([]*LangRating, len(l)))
	copy(sorter, l)
	sort.Sort(sorter)
	return &OverallRating{
		Rating: Rating{
			Correct: correct,
			Total:   total,
		},
		LangRatings: []*LangRating(sorter),
	}
}

func (o *OverallRating) LongestLangName() string {
	var longest string
	for _, lang := range o.LangRatings {
		if len(lang.Language) > len(longest) {
			longest = lang.Language
		}
	}
	return longest
}

type LangRating struct {
	Rating
	Language string
}

type ratingSorter []*LangRating

func (r ratingSorter) Len() int {
	return len(r)
}

func (r ratingSorter) Less(i, j int) bool {
	return r[i].Frac() > r[j].Frac()
}

func (r ratingSorter) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
