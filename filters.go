package filters

import (
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/processing"
)

type Fn = processing.FilterFn

func Type(ty ...vocab.ActivityVocabularyType) func(it vocab.Item) bool {
	types := vocab.ActivityVocabularyTypes(ty)
	return func(it vocab.Item) bool {
		result := false
		vocab.OnObject(it, func(object *vocab.Object) error {
			result = types.Contains(it.GetType())
			return nil
		})
		return result
	}
}
