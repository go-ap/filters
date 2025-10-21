package filters

import (
	"bytes"
	"net/url"

	vocab "github.com/go-ap/activitypub"
	"golang.org/x/text/unicode/norm"
)

type nlvType uint8

const (
	byName nlvType = iota
	byPreferredUsername
	bySummary
	byContent
)

type naturalLanguageValCheck struct {
	checkValue string
	checkFn    naturalLanguageValuesCheckFn
	accumFn    func(vocab.Item) []vocab.NaturalLanguageValues
	typ        nlvType
}

func (n naturalLanguageValCheck) Match(it vocab.Item) bool {
	return n.checkFn(n.accumFn(it), n.checkValue)
}

// NameIs checks an [vocab.Object]'s Name, or, in the case of an [vocab.Actor]
// also the PreferredUsername against the "name" value.
// If any of the Language Ref map values match the value, the function returns true.
func NameIs(name string) Check {
	return nameCheck(name, naturalLanguageValuesEquals)
}

// NameLike checks an [vocab.Object]'s Name, or, in the case of an [vocab.Actor]
// // also the PreferredUsername against the "name" value.
// If any of the Language Ref map values contains the value as a substring,
// the function returns true.
func NameLike(name string) Check {
	return nameCheck(name, naturalLanguageValuesLike)
}

// NameEmpty checks an [vocab.Object]'s Name, *and*, in the case of an [vocab.Actor]
// also its PreferredUsername to be empty.
// If *all* of the values are empty, the function returns true.
//
// Please note that the logic of this check is different from NameIs and NameLike.
var NameEmpty = nameCheck("", naturalLanguageEmpty)

// ContentIs checks an [vocab.Object]'s Content against the "cont" value.
// If any of the Language Ref map values match the value, the function returns true.
func ContentIs(cont string) Check {
	return contentCheck(cont, naturalLanguageValuesEquals)
}

// ContentLike checks an [vocab.Object]'s Content property against the "cont" value.
// If any of the Language Ref map values contains the value as a substring,
// the function returns true.
func ContentLike(cont string) Check {
	return contentCheck(cont, naturalLanguageValuesLike)
}

// ContentEmpty checks an [vocab.Object]'s Content, *and*, in the case of an [vocab.Actor]
// also the PreferredUsername to be empty.
// If *all* of the values are empty, the function returns true.
//
// Please note that the logic of this check is different from ContentIs and ContentLike.
var ContentEmpty = contentCheck("", naturalLanguageEmpty)

// SummaryIs checks an [vocab.Object]'s Summary against the "sum" value.
// If any of the Language Ref map values match the value, the function returns true.
func SummaryIs(sum string) Check {
	return summaryCheck(sum, naturalLanguageValuesEquals)
}

// SummaryLike checks an [vocab.Object]'s Summary property against the "sum" value.
// If any of the Language Ref map values contains the value as a substring,
// the function returns true.
func SummaryLike(sum string) Check {
	return summaryCheck(sum, naturalLanguageValuesLike)
}

// SummaryEmpty checks an [vocab.Object]'s Summary, *and*, in the case of an [vocab.Actor]
// also the PreferredUsername to be empty.
// If *all* of the values are empty, the function returns true.
//
// Please note that the logic of this check is different from SummaryIs and SummaryLike.
var SummaryEmpty = summaryCheck("", naturalLanguageEmpty)

func naturalLanguageValuesEquals(check []vocab.NaturalLanguageValues, val string) bool {
	nfc := norm.NFC.Bytes

	val, _ = url.QueryUnescape(val)
	for _, nlv := range check {
		for _, c := range nlv {
			if c != nil && bytes.EqualFold(nfc(c), nfc([]byte(val))) {
				return true
			}
		}
	}
	return false
}

func naturalLanguageEmpty(check []vocab.NaturalLanguageValues, _ string) bool {
	cnt := 0
	for _, nlv := range check {
		cnt += len(nlv)
	}
	return cnt == 0
}

func naturalLanguageValuesLike(check []vocab.NaturalLanguageValues, val string) bool {
	nfc := norm.NFC.Bytes

	val, _ = url.QueryUnescape(val)
	for _, nlv := range check {
		for _, c := range nlv {
			if c != nil && bytes.Contains(nfc(c), nfc([]byte(val))) {
				return true
			}
		}
	}
	return false
}

type naturalLanguageValuesCheckFn func([]vocab.NaturalLanguageValues, string) bool

func nameCheck(name string, checkFn naturalLanguageValuesCheckFn) Check {
	return naturalLanguageValCheck{
		checkValue: name,
		checkFn:    checkFn,
		accumFn:    loadName,
		typ:        byName,
	}
}

func loadName(it vocab.Item) []vocab.NaturalLanguageValues {
	if vocab.IsNil(it) {
		return nil
	}
	toCheck := make([]vocab.NaturalLanguageValues, 0)
	_ = vocab.OnObject(it, func(ob *vocab.Object) error {
		if len(ob.Name) > 0 {
			toCheck = append(toCheck, ob.Name)
		}
		return nil
	})
	_ = vocab.OnActor(it, func(act *vocab.Actor) error {
		if len(act.PreferredUsername) > 0 {
			toCheck = append(toCheck, act.PreferredUsername)
		}
		return nil
	})
	return toCheck
}

func contentCheck(content string, checkFn naturalLanguageValuesCheckFn) Check {
	return naturalLanguageValCheck{
		checkValue: content,
		checkFn:    checkFn,
		accumFn:    loadContent,
		typ:        byContent,
	}
}

func loadContent(it vocab.Item) []vocab.NaturalLanguageValues {
	if vocab.IsNil(it) {
		return nil
	}
	toCheck := make([]vocab.NaturalLanguageValues, 0)
	_ = vocab.OnObject(it, func(ob *vocab.Object) error {
		toCheck = append(toCheck, ob.Content)
		return nil
	})
	return toCheck
}

func summaryCheck(summary string, checkFn naturalLanguageValuesCheckFn) Check {
	return naturalLanguageValCheck{
		checkValue: summary,
		checkFn:    checkFn,
		accumFn:    loadSummary,
		typ:        bySummary,
	}
}

func loadSummary(it vocab.Item) []vocab.NaturalLanguageValues {
	if vocab.IsNil(it) {
		return nil
	}
	toCheck := make([]vocab.NaturalLanguageValues, 0)
	_ = vocab.OnObject(it, func(ob *vocab.Object) error {
		toCheck = append(toCheck, ob.Summary)
		return nil
	})
	return toCheck
}
