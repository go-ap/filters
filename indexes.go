package filters

import (
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/filters/index"
)

// AggregateFilters converts the received [Checks] into a list of [index.BasicFilter] objects
// that can be used to query an [index.Index].
func AggregateFilters(fil ...Check) []index.BasicFilter {
	types := make([]index.BasicFilter, 0, len(fil))
	for _, f := range fil {
		switch ff := f.(type) {
		case naturalLanguageValCheck:
			switch ff.typ {
			case byName:
				types = append(types, index.BasicFilter{Type: index.ByName, Values: []string{ff.checkValue}})
				types = append(types, index.BasicFilter{Type: index.ByPreferredUsername, Values: []string{ff.checkValue}})
			case bySummary:
				types = append(types, index.BasicFilter{Type: index.BySummary, Values: []string{ff.checkValue}})
			case byContent:
				types = append(types, index.BasicFilter{Type: index.ByContent, Values: []string{ff.checkValue}})
			default:
			}
		case actorChecks:
			values := make([]string, 0)
			for _, af := range ff {
				if ie, ok := af.(iriEquals); ok {
					values = append(values, vocab.IRI(ie).String())
				}
			}
			if len(values) > 0 {
				types = append(types, index.BasicFilter{Type: index.ByActor, Values: values})
			}
		case objectChecks:
			values := make([]string, 0)
			for _, af := range ff {
				if ie, ok := af.(iriEquals); ok {
					values = append(values, vocab.IRI(ie).String())
				}
			}
			if len(values) > 0 {
				types = append(types, index.BasicFilter{Type: index.ByObject, Values: values})
			}
		case attributedToEquals:
			ie := vocab.IRI(ff)
			types = append(types, index.BasicFilter{Type: index.ByAttributedTo, Values: []string{ie.String()}})
		case withTypes:
			values := make([]string, 0)
			for _, tf := range ff {
				values = append(values, string(tf))
			}
			if len(values) > 0 {
				types = append(types, index.BasicFilter{Type: index.ByType, Values: values})
			}
		case authorized:
			ie := vocab.IRI(ff)
			types = append(types, index.BasicFilter{Type: index.ByRecipients, Values: []string{ie.String()}})
		}
	}
	return types
}
