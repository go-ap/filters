package filters

import (
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/filters/index"
)

// AggregateFilters converts the received [Checks] into a list of [index.BasicFilter] objects
// that can be used to query an [index.Index].
func AggregateFilters(fil ...Check) []index.BasicFilter {
	if len(fil) == 0 {
		return nil
	}

	types := make([]index.BasicFilter, 0, len(fil))
	for _, f := range fil {
		switch ff := f.(type) {
		case naturalLanguageValCheck:
			switch ff.typ {
			case byName:
				types = append(types, index.BasicFilter{Type: index.ByName, Values: []string{ff.checkValue}})
				// NOTE(marius): the naturalLanguageValChecks have this idiosyncrasy of doing name searches for
				// both Name and PreferredUsername fields.
				types = append(types, index.BasicFilter{Type: index.ByPreferredUsername, Values: []string{ff.checkValue}})
			case bySummary:
				types = append(types, index.BasicFilter{Type: index.BySummary, Values: []string{ff.checkValue}})
			case byContent:
				types = append(types, index.BasicFilter{Type: index.ByContent, Values: []string{ff.checkValue}})
			default:
			}
		case actorChecks:
			if values := objectCheckValues(ff); len(values) > 0 {
				types = append(types, index.BasicFilter{Type: index.ByActor, Op: index.OPEq, Values: values})
			}
		case objectChecks:
			if values := objectCheckValues(ff); len(values) > 0 {
				types = append(types, index.BasicFilter{Type: index.ByObject, Op: index.OPEq, Values: values})
			}
		case attributedToEquals:
			ie := vocab.IRI(ff)
			types = append(types, index.BasicFilter{Type: index.ByAttributedTo, Op: index.OPEq, Values: []string{ie.String()}})
		case withTypes:
			values := make([]string, 0)
			for _, tf := range ff {
				values = append(values, string(tf))
			}
			if len(values) > 0 {
				types = append(types, index.BasicFilter{Type: index.ByType, Op: index.OPEq, Values: values})
			}
		case authorized:
			ie := vocab.IRI(ff)
			types = append(types, index.BasicFilter{Type: index.ByRecipients, Op: index.OPEq, Values: []string{ie.String()}})
		case recipients:
			ie := vocab.IRI(ff)
			types = append(types, index.BasicFilter{Type: index.ByRecipients, Op: index.OPEq, Values: []string{ie.String()}})
		}
	}
	return types
}

func objectCheckValues(ff []Check) []string {
	values := make([]string, 0, len(ff))
	for _, af := range ff {
		if ie, ok := af.(iriEquals); ok {
			values = append(values, vocab.IRI(ie).String())
		}
		if ie, ok := af.(idEquals); ok {
			values = append(values, vocab.IRI(ie).String())
		}
	}
	return values
}
