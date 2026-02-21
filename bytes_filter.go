package filters

import (
	"bytes"
	"fmt"
	"strconv"

	"quamina.net/go/quamina"
)

func appendS(s *bytes.Buffer, key string) {
	s.WriteRune('"')
	s.WriteString(key)
	s.WriteRune('"')
}

func appendV(s *bytes.Buffer, v any) {
	if vv, ok := v.([]any); ok {
		s.WriteRune('[')
		for i, v := range vv {
			appendLiteral(s, v)
			if i < len(vv)-1 {
				s.WriteRune(',')
			}
		}
		s.WriteRune(']')
	} else {
		appendLiteral(s, v)
	}
}

type qPrefix string

func formatQPrefix(c qPrefix) string {
	return `{"prefix":"` + string(c) + `"}`
}

type qExists bool

func formatQExists(c qExists) string {
	return `{"exists":` + strconv.FormatBool(bool(c)) + `}`
}

type qAnythingBut string

func formatQAnythingBut(c qAnythingBut) string {
	return `{"anything-but":"` + string(c) + `"}`
}

type qPattern string

func appendLiteral(s *bytes.Buffer, v any) {
	if v == nil {
		s.WriteString("null")
		return
	}

	switch vv := v.(type) {
	case qPattern:
		s.WriteString(string(vv))
	case qPrefix:
		s.WriteString(formatQPrefix(vv))
	case qExists:
		s.WriteString(formatQExists(vv))
	case qAnythingBut:
		s.WriteString(formatQAnythingBut(vv))
	case bool:
		s.WriteString(strconv.FormatBool(vv))
	case string:
		appendS(s, vv)
	case int:
		s.WriteString(strconv.FormatInt(int64(vv), 10))
	case int8:
		s.WriteString(strconv.FormatInt(int64(vv), 10))
	case int16:
		s.WriteString(strconv.FormatInt(int64(vv), 10))
	case int32:
		s.WriteString(strconv.FormatInt(int64(vv), 10))
	case int64:
		s.WriteString(strconv.FormatInt(vv, 10))
	case uint:
		s.WriteString(strconv.FormatUint(uint64(vv), 10))
	case uint8:
		s.WriteString(strconv.FormatUint(uint64(vv), 10))
	case uint16:
		s.WriteString(strconv.FormatUint(uint64(vv), 10))
	case uint32:
		s.WriteString(strconv.FormatUint(uint64(vv), 10))
	case uint64:
		s.WriteString(strconv.FormatUint(vv, 10))
	case float32:
		s.WriteString(strconv.FormatFloat(float64(vv), 10, -1, 32))
	case float64:
		s.WriteString(strconv.FormatFloat(vv, 10, -1, 64))
	}
}

func quaminaPattern(c Checks) []byte {
	if len(c) == 0 {
		return nil
	}

	s := bytes.Buffer{}
	s.WriteRune('{')
	for i, ff := range c {
		appendS(&s, checkName(ff))
		s.WriteRune(':')
		appendV(&s, checkValue(ff))
		if i < len(c)-1 {
			s.WriteRune(',')
		}
	}
	s.WriteRune('}')
	return s.Bytes()
}

func checkValue(ff Check) any {
	switch c := ff.(type) {
	case notCrit:
		r := make([]any, 0, len(c))
		for _, v := range c {
			if _, ok := v.(idNil); ok {
				r = append(r, qExists(true))
			}
			if _, ok := v.(iriNil); ok {
				r = append(r, qExists(true))
			}
			if s, ok := v.(iriEquals); ok {
				r = append(r, qAnythingBut(s))
			}
			if s, ok := v.(idEquals); ok {
				r = append(r, qAnythingBut(s))
			}
		}
		return r
	case iriNil:
		return []any{qExists(false)}
	case iriEquals:
		return []any{string(c)}
	case iriLike:
		return []any{qPrefix(c)}
	case idNil:
		return []any{qExists(false)}
	case idEquals:
		return []any{string(c)}
	case idLike:
		return []any{qPrefix(c)}
	case withTypes:
		r := make([]any, 0, len(c))
		for _, v := range c {
			r = append(r, string(v))
		}
		return r
	case objectChecks:
		return qPattern(quaminaPattern(Checks(c)))
	case actorChecks:
		return qPattern(quaminaPattern(Checks(c)))
	case targetChecks:
		return qPattern(quaminaPattern(Checks(c)))
	case tagChecks:
		vv := make([]any, 0, len(c))
		for _, v := range c {
			vv = append(vv, checkValue(v))
		}
		return vv
	case *beforeCrit:
		vv := make([]any, 0, len(c.fns))
		for _, v := range c.fns {
			vv = append(vv, checkValue(v))
		}
		return vv
	case *afterCrit:
		vv := make([]any, 0, len(c.fns))
		for _, v := range c.fns {
			vv = append(vv, checkValue(v))
		}
		return vv
	case *counter:
		return []any{c.max}
	}
	return nil
}

func checkName(ff Check) string {
	switch c := ff.(type) {
	case notCrit:
		for _, v := range c {
			if _, ok := v.(idNil); ok {
				return keyID
			}
			if _, ok := v.(iriNil); ok {
				return keyID
			}
			if _, ok := v.(iriEquals); ok {
				return keyID
			}
			if _, ok := v.(idEquals); ok {
				return keyID
			}
		}
	case idNil:
		return keyID
	case idEquals:
		return keyID
	case idLike:
		return keyID
	case iriNil:
		return keyID
	case iriEquals:
		return keyID
	case iriLike:
		return keyID
	case withTypes:
		return keyType
	case objectChecks:
		return keyObject
	case actorChecks:
		return keyActor
	case targetChecks:
		return keyTarget
	case tagChecks:
		return keyTag
	case *beforeCrit:
		return keyBefore
	case *afterCrit:
		return keyAfter
	case *counter:
		return keyMaxItems
	}
	return ""
}

func rawFilter(filters Checks) (*quamina.Quamina, error) {
	q, err := quamina.New()
	if err != nil {
		return nil, err
	}
	pattern := quaminaPattern(filters)
	if len(pattern) == 0 {
		return nil, fmt.Errorf("unable to create quamina pattern")
	}
	if err = q.AddPattern("filter", string(pattern)); err != nil {
		return nil, err
	}
	return q, nil
}

func MatchRaw(filters Checks, raw []byte) bool {
	q, err := rawFilter(filters)
	if err != nil {
		return false
	}
	matchAny, err := q.MatchesForEvent(raw)
	if err != nil {
		return false
	}
	return len(matchAny) > 0
}
