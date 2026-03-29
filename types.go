package filters

import (
	"slices"
	"strings"

	vocab "github.com/go-ap/activitypub"
)

// VocabularyTypesFilter converts the received list of strings to a list of ActivityVocabularyTypes
// that can be used with the HasType filter function.
// The individual strings are not validated against the known vocabulary types.
func VocabularyTypesFilter(types ...string) vocab.ActivityVocabularyTypes {
	r := make(vocab.ActivityVocabularyTypes, 0, len(types))
	for _, t := range types {
		typ := vocab.ActivityVocabularyType(t)
		if slices.Contains(vocab.Types, typ) {
			r = append(r, typ)
		}
	}
	return r
}

// HasType checks an activitypub.Object's Type property against a series of values.
// If any of the ty values matches, the function returns true.
func HasType(ty ...vocab.ActivityVocabularyType) Check {
	return withTypes(ty)
}

type withTypes vocab.ActivityVocabularyTypes

func (tt withTypes) Match(it vocab.Item) bool {
	if vocab.IsNil(it) || it.GetType() == nil {
		return len(tt) == 0
	}
	matchFn := func(ob vocab.Item) bool {
		withType := vocab.ActivityVocabularyTypes{}
		if typ := ob.GetType(); typ != nil {
			withType = typ.AsTypes()
		}
		return vocab.AnyTypes(tt...).Match(withType...)
	}

	if !vocab.IsItemCollection(it) {
		return matchFn(it)
	}

	itemsHaveType := false
	_ = vocab.OnItem(it, func(it vocab.Item) error {
		itemsHaveType = itemsHaveType || matchFn(it)
		return nil
	})
	return itemsHaveType
}

func (tt withTypes) GoString() string {
	ss := strings.Builder{}
	ss.WriteString("type=[")
	for i, typ := range tt {
		ss.WriteString(string(typ))
		if i < len(tt)-1 {
			ss.WriteRune(',')
		}
	}
	ss.WriteRune(']')
	return ss.String()
}
