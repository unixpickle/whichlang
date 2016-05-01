package idtree

import (
	"strconv"
	"strings"

	"github.com/unixpickle/weakai/idtrees"
	"github.com/unixpickle/whichlang/tokens"
)

func Train(freqs map[string][]tokens.Freqs) *Classifier {
	allWords := map[string]bool{}
	entryCount := 0
	for _, samples := range freqs {
		for _, sample := range samples {
			for word := range sample {
				allWords[word] = true
			}
			entryCount++
		}
	}

	dataSet := &idtrees.DataSet{
		Entries: make([]idtrees.Entry, 0, entryCount),
		Fields:  make([]idtrees.Field, 0),
	}

	for lang, list := range freqs {
		for _, wordMap := range list {
			entry := &treeEntry{
				language:    lang,
				freqs:       wordMap,
				fieldValues: []idtrees.Value{},
			}
			dataSet.Entries = append(dataSet.Entries, entry)
		}
	}

	for word := range allWords {
		idtrees.CreateBisectingFloatFields(dataSet, func(e idtrees.Entry) float64 {
			return e.(*treeEntry).freqs[word]
		}, func(e idtrees.Entry, v idtrees.Value) {
			te := e.(*treeEntry)
			te.fieldValues = append(te.fieldValues, v)
		}, strings.Replace(word, "%", "%%", -1)+" > %f")
	}

	tree := idtrees.GenerateTree(dataSet)
	res := convertTree(tree)
	centerThresholdsRoot(res, freqs)

	return res
}

func convertTree(t *idtrees.TreeNode) *Classifier {
	if t.BranchField == nil {
		if t.LeafValue == nil {
			lang := "Unknown"
			return &Classifier{LeafClassification: &lang}
		} else {
			lang := t.LeafValue.String()
			return &Classifier{LeafClassification: &lang}
		}
	}

	comps := strings.Split(t.BranchField.String(), " ")
	if len(comps) != 3 {
		panic("unknown branch field: " + t.BranchField.String())
	}
	val, _ := strconv.ParseFloat(comps[2], 64)
	res := &Classifier{
		Keyword:   comps[0],
		Threshold: val,
	}
	res.FalseBranch = convertTree(t.Branches[idtrees.BoolValue(false)])
	res.TrueBranch = convertTree(t.Branches[idtrees.BoolValue(true)])
	return res
}

type treeEntry struct {
	language    string
	freqs       tokens.Freqs
	fieldValues []idtrees.Value
}

func (t *treeEntry) FieldValues() []idtrees.Value {
	return t.fieldValues
}

func (t *treeEntry) Class() idtrees.Value {
	return idtrees.StringValue(t.language)
}
