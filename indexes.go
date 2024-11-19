package filters

import (
	"github.com/RoaringBitmap/roaring"
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/filters/index"
)

func extractBitmaps(checks Checks, indexes map[index.Type]index.Indexable) []*roaring.Bitmap {
	result := make([]*roaring.Bitmap, 0)
	for _, check := range checks {
		switch fil := check.(type) {
		case idEquals:
			result = append(result, index.GetBitmaps[vocab.IRI](indexes[index.ByID], vocab.IRI(fil)))
		case iriEquals:
			result = append(result, index.GetBitmaps[vocab.IRI](indexes[index.ByID], vocab.IRI(fil)))
		case checkAny:
			anys := extractBitmaps(Checks(fil), indexes)
			result = append(result, roaring.FastOr(anys...))
		case checkAll:
			alls := extractBitmaps(Checks(fil), indexes)
			result = append(result, roaring.FastAnd(alls...))
		case naturalLanguageValCheck:
			switch fil.typ {
			case byName:
				// NOTE(marius): the naturalLanguageValChecks have this idiosyncrasy of doing name searches for
				// both Name and PreferredUsername fields, so until we split them, we should use the same logic here.
				ors := []*roaring.Bitmap{
					index.GetBitmaps[string](indexes[index.ByName], fil.checkValue),
					index.GetBitmaps[string](indexes[index.ByPreferredUsername], fil.checkValue),
				}
				result = append(result, roaring.FastOr(ors...))
			case bySummary:
				result = append(result, index.GetBitmaps[string](indexes[index.BySummary], fil.checkValue))
			case byContent:
				result = append(result, index.GetBitmaps[string](indexes[index.ByContent], fil.checkValue))
			default:
			}
		case withTypes:
			ors := make([]*roaring.Bitmap, 0)
			for _, tf := range fil {
				ors = append(ors, index.GetBitmaps[string](indexes[index.ByType], string(tf)))
			}
			if len(ors) > 0 {
				result = append(result, roaring.FastOr(ors...))
			}
		case actorChecks:
			ors := make([]*roaring.Bitmap, 0)
			if values := objectCheckValues(fil); len(values) > 0 {
				for _, val := range values {
					ors = append(ors, index.GetBitmaps[vocab.IRI](indexes[index.ByActor], vocab.IRI(val)))
				}
			}
			if len(ors) > 0 {
				result = append(result, roaring.FastOr(ors...))
			}
		case objectChecks:
			ors := make([]*roaring.Bitmap, 0)
			if values := objectCheckValues(fil); len(values) > 0 {
				for _, val := range values {
					ors = append(ors, index.GetBitmaps[vocab.IRI](indexes[index.ByObject], vocab.IRI(val)))
				}
			}
			if len(ors) > 0 {
				result = append(result, roaring.FastOr(ors...))
			}
		case attributedToEquals:
			result = append(result, index.GetBitmaps[vocab.IRI](indexes[index.ByAttributedTo], vocab.IRI(fil)))
		case authorized:
			result = append(result, index.GetBitmaps[vocab.IRI](indexes[index.ByRecipients], vocab.IRI(fil)))
		case recipients:
			result = append(result, index.GetBitmaps[vocab.IRI](indexes[index.ByRecipients], vocab.IRI(fil)))
		}
	}
	return result
}

func (ff Checks) IndexMatch(indexes map[index.Type]index.Indexable) *roaring.Bitmap {
	if len(ff) == 0 {
		return nil
	}

	// NOTE(marius): A normal list of Check functions in this package corresponds
	// to a filter equivalent of All(Checks...).
	// We can therefore use an AND operator for the bitmaps.
	ands := extractBitmaps(ff, indexes)
	return roaring.FastAnd(ands...)
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
		if ie, ok := af.(withTypes); ok {
			for _, typ := range ie {
				values = append(values, string(typ))
			}
		}
	}
	return values
}

// SearchIndex does a fast index search for the received filters.
func SearchIndex(i *index.Index, ff ...Check) ([]vocab.IRI, error) {
	bmp := Checks(ff).IndexMatch(i.Indexes)
	result := make([]vocab.IRI, 0, bmp.GetCardinality())
	it := bmp.Iterator()
	for x := it.PeekNext(); it.HasNext(); x = it.Next() {
		if iri, ok := i.Ref[x]; ok {
			result = append(result, iri)
		}
	}
	return result, nil
}
