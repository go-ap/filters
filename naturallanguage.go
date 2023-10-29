package filters

import (
	"strings"

	vocab "github.com/go-ap/activitypub"
	"golang.org/x/text/unicode/norm"
)

// NameIs checks an activitypub.Object's Name, or, in the case of an activitypub.Actor
// also the PreferredUsername against the "name" value.
// If any of the Language Ref map values match the value, the function returns true.
func NameIs(name string) Fn {
	return nameCheck(name, naturalLanguageValuesEquals)
}

// NameLike checks an activitypub.Object's Name, or, in the case of an activitypub.Actor
// // also the PreferredUsername against the "name" value.
// If any of the Language Ref map values contains the value as a substring,
// the function returns true.
func NameLike(name string) Fn {
	return nameCheck(name, naturalLanguageValuesLike)
}

// NameEmpty checks an activitypub.Object's Name, *and*, in the case of an activitypub.Actor
// also the PreferredUsername to be empty.
// If *all* of the values are empty, the function returns true.
//
// Please note that the logic of this check is different from NameIs and NameLike.
func NameEmpty() Fn {
	return nameCheck("", naturalLanguageEmpty)
}

// ContentIs checks an activitypub.Object's Content against the "cont" value.
// If any of the Language Ref map values match the value, the function returns true.
func ContentIs(cont string) Fn {
	return contentCheck(cont, naturalLanguageValuesEquals)
}

// ContentLike checks an activitypub.Object's Content property against the "cont" value.
// If any of the Language Ref map values contains the value as a substring,
// the function returns true.
func ContentLike(cont string) Fn {
	return contentCheck(cont, naturalLanguageValuesLike)
}

// ContentEmpty checks an activitypub.Object's Content, *and*, in the case of an activitypub.Actor
// also the PreferredUsername to be empty.
// If *all* of the values are empty, the function returns true.
//
// Please note that the logic of this check is different from ContentIs and ContentLike.
func ContentEmpty() Fn {
	return contentCheck("", naturalLanguageEmpty)
}

// SummaryIs checks an activitypub.Object's Summary against the "sum" value.
// If any of the Language Ref map values match the value, the function returns true.
func SummaryIs(sum string) Fn {
	return summaryCheck(sum, naturalLanguageValuesEquals)
}

// SummaryLike checks an activitypub.Object's Summary property against the "sum" value.
// If any of the Language Ref map values contains the value as a substring,
// the function returns true.
func SummaryLike(sum string) Fn {
	return summaryCheck(sum, naturalLanguageValuesLike)
}

// SummaryEmpty checks an activitypub.Object's Summary, *and*, in the case of an activitypub.Actor
// also the PreferredUsername to be empty.
// If *all* of the values are empty, the function returns true.
//
// Please note that the logic of this check is different from SummaryIs and SummaryLike.
func SummaryEmpty() Fn {
	return summaryCheck("", naturalLanguageEmpty)
}

func naturalLanguageValuesEquals(check vocab.NaturalLanguageValues, val string) bool {
	for _, c := range check {
		if strings.EqualFold(c.String(), val) {
			return true
		}
	}
	return false
}

func naturalLanguageEmpty(check vocab.NaturalLanguageValues, _ string) bool {
	return len(check) == 0
}

func naturalLanguageValuesLike(check vocab.NaturalLanguageValues, val string) bool {
	nfc := norm.NFC.String
	for _, c := range check {
		if strings.Contains(nfc(c.String()), nfc(val)) {
			return true
		}
	}
	return false
}

type naturalLanguageValuesCheckFn func(vocab.NaturalLanguageValues, string) bool

func nameCheck(name string, checkFn naturalLanguageValuesCheckFn) Fn {
	return func(it vocab.Item) bool {
		if vocab.IsNil(it) {
			return false
		}
		toCheck := make(vocab.NaturalLanguageValues, 0)
		_ = vocab.OnObject(it, func(ob *vocab.Object) error {
			if len(ob.Name) > 0 {
				toCheck = append(toCheck, ob.Name...)
			}
			return nil
		})
		_ = vocab.OnActor(it, func(act *vocab.Actor) error {
			if len(act.PreferredUsername) > 0 {
				toCheck = append(toCheck, act.PreferredUsername...)
			}
			return nil
		})
		return checkFn(toCheck, name)
	}
}

func contentCheck(content string, checkFn naturalLanguageValuesCheckFn) Fn {
	return func(it vocab.Item) bool {
		var toCheck vocab.NaturalLanguageValues
		_ = vocab.OnObject(it, func(ob *vocab.Object) error {
			toCheck = ob.Content
			return nil
		})
		return checkFn(toCheck, content)
	}
}

func summaryCheck(content string, checkFn naturalLanguageValuesCheckFn) Fn {
	return func(it vocab.Item) bool {
		var toCheck vocab.NaturalLanguageValues
		_ = vocab.OnObject(it, func(ob *vocab.Object) error {
			toCheck = ob.Summary
			return nil
		})
		return checkFn(toCheck, content)
	}
}
