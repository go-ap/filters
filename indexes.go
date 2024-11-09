package filters

import (
	"github.com/RoaringBitmap/roaring"
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/filters/index"
)

func (ff Checks) IndexMatch(indexes map[index.Type]index.Indexable) *roaring.Bitmap {
	if len(ff) == 0 {
		return nil
	}

	ands := make([]*roaring.Bitmap, 0)
	for _, f := range ff {
		switch ff := f.(type) {
		case naturalLanguageValCheck:
			switch ff.typ {
			case byName:
				// NOTE(marius): the naturalLanguageValChecks have this idiosyncrasy of doing name searches for
				// both Name and PreferredUsername fields, so until we split them, we should use the same logic here.
				ors := []*roaring.Bitmap{
					index.GetBitmaps[string](indexes[index.ByName], ff.checkValue),
					index.GetBitmaps[string](indexes[index.ByPreferredUsername], ff.checkValue),
				}
				ands = append(ands, roaring.FastOr(ors...))
			case bySummary:
				ands = append(ands, index.GetBitmaps[string](indexes[index.BySummary], ff.checkValue))
			case byContent:
				ands = append(ands, index.GetBitmaps[string](indexes[index.ByContent], ff.checkValue))
			default:
			}
		case withTypes:
			ors := make([]*roaring.Bitmap, 0)
			for _, tf := range ff {
				ors = append(ors, index.GetBitmaps[string](indexes[index.ByType], string(tf)))
			}
			if len(ors) > 0 {
				ands = append(ands, roaring.FastOr(ors...))
			}
		case actorChecks:
			ors := make([]*roaring.Bitmap, 0)
			if values := objectCheckValues(ff); len(values) > 0 {
				for _, val := range values {
					ors = append(ors, index.GetBitmaps[vocab.IRI](indexes[index.ByActor], val))
				}
			}
			if len(ors) > 0 {
				ands = append(ands, roaring.FastOr(ors...))
			}
		case objectChecks:
			ors := make([]*roaring.Bitmap, 0)
			if values := objectCheckValues(ff); len(values) > 0 {
				for _, val := range values {
					ors = append(ors, index.GetBitmaps[vocab.IRI](indexes[index.ByObject], val))
				}
			}
			if len(ors) > 0 {
				ands = append(ands, roaring.FastOr(ors...))
			}
		case attributedToEquals:
			ands = append(ands, index.GetBitmaps[vocab.IRI](indexes[index.ByAttributedTo], string(ff)))
		case authorized:
			ands = append(ands, index.GetBitmaps[vocab.IRI](indexes[index.ByRecipients], string(ff)))
		case recipients:
			ands = append(ands, index.GetBitmaps[vocab.IRI](indexes[index.ByRecipients], string(ff)))
		}
	}
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
	}
	return values
}

// SearchIndex does a fast index search for the received filters.
func SearchIndex(i *index.Index, ff ...Check) ([]vocab.IRI, error) {
	bmp := Checks(ff).IndexMatch(i.Indexes)
	result := make([]vocab.IRI, 0, bmp.GetCardinality())
	bmp.Iterate(func(x uint32) bool {
		if iri, ok := i.Ref[x]; ok {
			result = append(result, iri)
		}
		return true
	})
	return result, nil
}
