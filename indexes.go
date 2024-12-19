package filters

import (
	"github.com/RoaringBitmap/roaring/roaring64"
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/filters/index"
)

var hFn = index.HashFn

func extractBitmapsForSubprop(checks Checks, indexes map[index.Type]index.Indexable, typ index.Type) []*roaring64.Bitmap {
	found := roaring64.FastAnd(extractBitmaps(checks, indexes)...)
	if found.GetCardinality() == 0 {
		return nil
	}

	iter := found.Iterator()
	if iter == nil {
		return nil
	}

	refs := make([]uint64, 0, found.GetCardinality())
	for x := iter.PeekNext(); iter.HasNext(); x = iter.Next() {
		refs = append(refs, x)
	}

	return index.GetBitmaps[uint64](indexes[typ], refs...)
}

const (
	ByID                = index.ByID
	ByType              = index.ByType
	ByName              = index.ByName
	ByPreferredUsername = index.ByPreferredUsername
	BySummary           = index.BySummary
	ByContent           = index.ByContent
	ByActor             = index.ByActor
	ByObject            = index.ByObject
	ByRecipients        = index.ByRecipients
	ByAttributedTo      = index.ByAttributedTo
	ByInReplyTo         = index.ByInReplyTo
)

func extractBitmaps(checks Checks, indexes map[index.Type]index.Indexable) []*roaring64.Bitmap {
	result := make([]*roaring64.Bitmap, 0)
	for _, check := range checks {
		switch fil := check.(type) {
		case notCrit:
			toExclude := roaring64.FastOr(extractBitmaps(Checks(fil), indexes)...)
			if toExclude.GetCardinality() == 0 {
				continue
			}
			all := roaring64.FastOr(index.GetBitmaps[uint64](indexes[ByID])...)
			if all.GetCardinality() > 0 {
				all.AndNot(toExclude)
				result = append(result, all)
			}
		case idEquals:
			result = append(result, index.GetBitmaps[uint64](indexes[ByID], hFn(vocab.IRI(fil)))...)
		case iriEquals:
			result = append(result, index.GetBitmaps[uint64](indexes[ByID], hFn(vocab.IRI(fil)))...)
		case checkAny:
			anys := extractBitmaps(Checks(fil), indexes)
			result = append(result, roaring64.FastOr(anys...))
		case checkAll:
			alls := extractBitmaps(Checks(fil), indexes)
			result = append(result, roaring64.FastAnd(alls...))
		case naturalLanguageValCheck:
			switch fil.typ {
			case byName:
				// NOTE(marius): the naturalLanguageValChecks have this idiosyncrasy of doing name searches for
				// both Name and PreferredUsername fields, so until we split them, we should use the same logic here.
				ors := make([]*roaring64.Bitmap, 0)
				ors = append(ors, index.GetBitmaps[string](indexes[ByName], fil.checkValue)...)
				ors = append(ors, index.GetBitmaps[string](indexes[ByPreferredUsername], fil.checkValue)...)
				if len(ors) > 0 {
					result = append(result, roaring64.FastOr(ors...))
				}
			case bySummary:
				result = append(result, index.GetBitmaps[string](indexes[BySummary], fil.checkValue)...)
			case byContent:
				result = append(result, index.GetBitmaps[string](indexes[ByContent], fil.checkValue)...)
			default:
			}
		case withTypes:
			ors := make([]*roaring64.Bitmap, 0)
			for _, tf := range fil {
				ors = append(ors, index.GetBitmaps[string](indexes[ByType], string(tf))...)
			}
			if len(ors) > 0 {
				result = append(result, roaring64.FastOr(ors...))
			}
		case actorChecks:
			actorRefs := extractBitmapsForSubprop(Checks(fil), indexes, ByActor)
			result = append(result, roaring64.FastOr(actorRefs...))
		case objectChecks:
			objectRefs := extractBitmapsForSubprop(Checks(fil), indexes, ByObject)
			result = append(result, roaring64.FastOr(objectRefs...))
		case attributedToEquals:
			result = append(result, index.GetBitmaps[uint64](indexes[ByAttributedTo], hFn(vocab.IRI(fil)))...)
		case inReplyToEquals:
			result = append(result, index.GetBitmaps[uint64](indexes[ByInReplyTo], hFn(vocab.IRI(fil)))...)
		case authorized:
			if iri := vocab.IRI(fil); iri.Equals(vocab.PublicNS, true) {
				result = append(result, index.GetBitmaps[uint64](indexes[ByRecipients], hFn(iri))...)
			} else {
				result = append(result,
					roaring64.FastOr(
						index.GetBitmaps[uint64](indexes[ByRecipients], hFn(vocab.PublicNS), hFn(iri))...,
					),
				)
			}
		case recipients:
			result = append(result, index.GetBitmaps[uint64](indexes[ByRecipients], hFn(vocab.IRI(fil)))...)
		}
	}
	return result
}

func (ff Checks) IndexMatch(indexes map[index.Type]index.Indexable) *roaring64.Bitmap {
	if len(ff) == 0 {
		return nil
	}

	// NOTE(marius): A normal list of Check functions in this package corresponds
	// to a filter equivalent of All(Checks...).
	// We can therefore use an AND operator for the bitmaps.
	ands := extractBitmaps(ff, indexes)
	return roaring64.FastAnd(ands...)
}

// SearchIndex does a fast index search for the received filters.
func SearchIndex(i *index.Index, ff ...Check) ([]vocab.IRI, error) {
	bmp := Checks(ff).IndexMatch(i.Indexes)

	if bmp.GetCardinality() == 0 {
		return nil, nil
	}

	it := bmp.Iterator()
	if it == nil {
		return nil, nil
	}

	result := make([]vocab.IRI, 0, bmp.GetCardinality())
	for x := it.PeekNext(); it.HasNext(); x = it.Next() {
		if iri, ok := i.Ref[x]; ok {
			result = append(result, iri)
		}
	}
	return result, nil
}
