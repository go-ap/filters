package filters

import (
	"strings"

	vocab "github.com/go-ap/activitypub"
)

// NilID checks if the activitypub.Object's ID property matches any of the two magic values
// that denote an empty value: activitypub.NilID = "-", or activitypub.EmptyID = ""
func NilID(it vocab.Item) bool {
	return Any(IRI(vocab.NilIRI), IRI(vocab.EmptyIRI))(it.GetLink())
}

// NotNilID checks if the activitypub.Object's ID property doesn't match any of the two magic values
// that denote an empty value: activitypub.NilID = "-", or activitypub.EmptyID = ""
func NotNilID(it vocab.Item) bool {
	return !NilID(it)
}

// ID checks an activitypub.Object's ID property against the received iri.
func ID(iri vocab.IRI) Fn {
	return func(item vocab.Item) bool {
		return item.GetID().Equals(iri, true)
	}
}

// IRI checks an activitypub.Object's IRI
func IRI(iri vocab.IRI) Fn {
	return func(item vocab.Item) bool {
		return item.GetLink().Equals(iri, true)
	}
}

// Type checks an activitypub.Object's Type property against a series of values.
// If any of the ty values matches, the function returns true.
func Type(ty ...vocab.ActivityVocabularyType) Fn {
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

// NameIs checks an activitypub.Object's Name property against the name value.
// If any of the Language Ref map values match the name, the function returns true.
func NameIs(name string) Fn {
	return func(it vocab.Item) bool {
		valid := false
		vocab.OnObject(it, func(ob *vocab.Object) error {
			for _, nn := range ob.Name {
				if strings.EqualFold(nn.String(), name) {
					valid = true
					break
				}
			}
			return nil
		})
		if valid {
			return valid
		}
		vocab.OnActor(it, func(act *vocab.Actor) error {
			for _, nn := range act.PreferredUsername {
				if strings.EqualFold(nn.String(), name) {
					valid = true
					break
				}
			}
			return nil
		})
		return valid
	}
}
