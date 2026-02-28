package filters

import (
	"bytes"
	"encoding/json"
	"maps"
	"reflect"
	"slices"
	"strconv"

	"quamina.net/go/quamina"
)

func appendS(s *bytes.Buffer, key string) {
	s.WriteRune('"')
	s.WriteString(key)
	s.WriteRune('"')
}

type qString string

func (qs qString) MarshalJSON() ([]byte, error) {
	return []byte(`"` + qs + `"`), nil
}

type qPrefix string

func (qp qPrefix) MarshalJSON() ([]byte, error) {
	return []byte(`{"prefix":"` + string(qp) + `"}`), nil
}

type qExists bool

func (qe qExists) MarshalJSON() ([]byte, error) {
	return []byte(`{"exists":` + strconv.FormatBool(bool(qe)) + `}`), nil
}

type qAnythingBut string

func (qab qAnythingBut) MarshalJSON() ([]byte, error) {
	return []byte(`{"anything-but":"` + string(qab) + `"}`), nil
}

var (
	nlvEqCheck    = reflect.ValueOf(naturalLanguageValuesEquals)
	nlvLikeCheck  = reflect.ValueOf(naturalLanguageValuesLike)
	nlvEmptyCheck = reflect.ValueOf(naturalLanguageEmpty)
)

type qLeaf json.Marshaler

type qLeafArray []qLeaf

func (la qLeafArray) MarshalJSON() ([]byte, error) {
	ss := bytes.Buffer{}
	if len(la) > 0 {
		ss.WriteRune('[')
	}
	for i, l := range la {
		raw, _ := l.MarshalJSON()
		ss.Write(raw)
		if i < len(la)-1 {
			ss.WriteRune(',')
		}
	}
	if len(la) > 0 {
		ss.WriteRune(']')
	}
	return ss.Bytes(), nil
}

type qFullPattern map[string]json.Marshaler

func (qp qFullPattern) MarshalJSON() ([]byte, error) {
	ss := bytes.Buffer{}
	if len(qp) > 0 {
		ss.WriteRune('{')
	}

	for i, k := range slices.Sorted(maps.Keys(qp)) {
		m, ok := qp[k]
		if !ok || m == nil {
			continue
		}
		vv, _ := m.MarshalJSON()
		if k == "" || vv == nil {
			continue
		}
		if i > 0 && i < len(qp) {
			ss.WriteRune(',')
		}

		appendS(&ss, k)
		ss.WriteRune(':')
		ss.Write(vv)
	}
	if len(qp) > 0 {
		ss.WriteRune('}')
	}
	return ss.Bytes(), nil
}

func buildFullPattern(c Checks) json.Marshaler {
	fp := make(qFullPattern)
	for _, cc := range c {
		k, vv := getLeafValue(cc)
		if k == "" && vv == nil {
			continue
		}
		fp[k] = vv
	}
	if len(fp) == 1 {
		for k, v := range fp {
			if (k == "-" || k == "") && v != nil {
				return v
			}
		}
	}

	return fp
}

func getLeafValue(ff Check) (string, json.Marshaler) {
	switch c := ff.(type) {
	case itemNil:
		return "", qLeafArray{qExists(false)}
	case notCrit:
		for _, v := range c {
			if _, ok := v.(itemNil); ok {
				return "-", qLeafArray{qExists(true)}
			}
			if _, ok := v.(idNil); ok {
				return keyID, qLeafArray{qExists(true)}
			}
			if _, ok := v.(iriNil); ok {
				return keyID, qLeafArray{qExists(true)}
			}
			if s, ok := v.(iriEquals); ok {
				return keyID, qLeafArray{qAnythingBut(s)}
			}
			if s, ok := v.(idEquals); ok {
				return keyID, qLeafArray{qAnythingBut(s)}
			}
		}
		return "-", nil
	case iriNil:
		return keyID, qLeafArray{qExists(false)}
	case iriEquals:
		return keyID, qLeafArray{qString(c)}
	case iriLike:
		return keyID, qLeafArray{qPrefix(c)}
	case idNil:
		return keyID, qLeafArray{qExists(false)}
	case idEquals:
		return keyID, qLeafArray{qString(c)}
	case idLike:
		return keyID, qLeafArray{qPrefix(c)}
	case withTypes:
		r := make(qLeafArray, 0, len(c))
		for _, v := range c {
			r = append(r, qString(v))
		}
		return keyType, r
	case naturalLanguageValCheck:
		var name string
		switch c.typ {
		case byName:
			name = keyName
		case byPreferredUsername:
			name = keyPreferredUsername
		case bySummary:
			name = keySummary
		case byContent:
			name = keyContent
		}
		rCheckFn := reflect.ValueOf(c.checkFn)
		switch rCheckFn.Pointer() {
		case nlvEqCheck.Pointer():
			return name, qLeafArray{qString(c.checkValue)}
		case nlvEmptyCheck.Pointer():
			return name, qLeafArray{qExists(false)}
		case nlvLikeCheck.Pointer():
			fallthrough
		default:
			return name, qLeafArray{qPrefix(c.checkValue)}
		}
	case objectChecks:
		return keyObject, buildFullPattern(Checks(c))
	case actorChecks:
		return keyActor, buildFullPattern(Checks(c))
	case targetChecks:
		return keyTarget, buildFullPattern(Checks(c))
	case tagChecks:
		return keyTag, buildFullPattern(Checks(c))
	case checkAll:
		return "-", buildFullPattern(Checks(c))
	}
	return "", nil
}

func RawMatcher(filters Checks) func([]byte) bool {
	alwaysT := func(_ []byte) bool {
		return true
	}

	if len(filters) == 0 {
		return alwaysT
	}
	q, err := quamina.New()
	if err != nil {
		// NOTE(marius): failed to initialize quamina, fallback to other filtering methods
		return alwaysT
	}

	// NOTE(marius): we build two patters, and if none is valid we skip matching.
	noValidPattern := true

	// NOTE(marius): The first pattern considers the full set of filters and assumes the raw JSON document
	// is in a denormalized/not-flattened form.
	pattern1, _ := buildFullPattern(filters).MarshalJSON()
	if len(pattern1) > 2 {
		if err = q.AddPattern("full", string(pattern1)); err == nil {
			noValidPattern = false
		}
	}
	// NOTE(marius): the second pattern assumes the JSON is normalized,
	// therefore we convert all filters beyond depth 1 to just an "exist" check.
	// This transforms the checks for Activities containing Object, Actor, or Target filters,
	// and the ones for Objects with Tag.
	pattern2, _ := buildFullPattern(TopLevelChecks(filters...)).MarshalJSON()
	if len(pattern2) > 2 && !bytes.Equal(pattern1, pattern2) {
		if err = q.AddPattern("top-level", string(pattern2)); err == nil {
			noValidPattern = false
		}
	}

	if noValidPattern {
		return alwaysT
	}

	return func(raw []byte) bool {
		matchAny, err := q.MatchesForEvent(raw)
		return err == nil && len(matchAny) > 0
	}
}

func MatchRaw(filters Checks, raw []byte) bool {
	return RawMatcher(filters)(raw)
}
