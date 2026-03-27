package filters

import (
	"fmt"
	"strings"

	vocab "github.com/go-ap/activitypub"
	"github.com/leporo/sqlf"
)

func GetLimit(f ...Check) int {
	for _, check := range f {
		switch c := check.(type) {
		case *counter:
			return c.max
		}
	}
	return -1
}

func GetWhereClauses(ff ...Check) ([]string, []any) {
	s := sqlf.Select("")
	getWhereClauses(s, ff...)
	var pieces []string
	if q := strings.TrimPrefix(s.String(), "SELECT"); len(q) > 0 {
		pieces = strings.Split(strings.TrimPrefix(q, " WHERE"), "AND ")
		for i := range pieces {
			pieces[i] = strings.TrimSpace(pieces[i])
		}
	}

	args := s.Args()
	if len(args) == 0 && len(pieces) == 0 {
		return nil, nil
	}
	return pieces, args
}

type Stmt = sqlf.Stmt

func getLimit(s *Stmt, f ...Check) {
	for _, check := range f {
		switch c := check.(type) {
		case *counter:
			s.Limit(c.max)
		}
	}
}

func BuildSQL(s *Stmt, ff ...Check) error {
	getWhereClauses(s, ff...)
	//getLimit(s, ff...)
	return nil
}

func getWhereClauses(s *Stmt, f ...Check) {
	getTypeWheres(s, f...)
	getIRIWheres(s, f...)
	getNamesWheres(s, f...)
	getInReplyToWheres(s, f...)
	getAttributedToWheres(s, f...)
	getURLWheres(s, f...)
	getContextWheres(s, f...)
}

func getIRIWheres(s *Stmt, f ...Check) (string, []any) {
	strs := make(CompStrs, 0)
	for _, check := range f {
		switch i := check.(type) {
		case idEquals:
			strs = append(strs, CompStr{Str: string(i)})
		case iriEquals:
			strs = append(strs, CompStr{Str: string(i)})
		case iriLike:
			strs = append(strs, CompStr{Operator: "~", Str: string(i)})
		case idLike:
			strs = append(strs, CompStr{Operator: "~", Str: string(i)})
		case iriNil:
			strs = append(strs, CompStr{Str: string(vocab.NilIRI)})
		case idNil:
			strs = append(strs, CompStr{Str: string(vocab.NilIRI)})
		}
	}
	return getStringFieldWheres(s, "iri", strs...)
}

func getStringFieldInJSONWheres(s *Stmt, prop string, strs ...CompStr) (string, []any) {
	if len(strs) == 0 {
		return "", nil
	}
	var values = make([]any, 0)
	keyWhere := make([]string, 0)

	isPg := stmtIsPostgres(s)
	for _, n := range strs {
		switch n.Operator {
		case "!":
			if len(n.Str) == 0 || n.Str == vocab.NilLangRef.String() {
				if isPg {
					s.Where(fmt.Sprintf(`json_extract("raw", '$.%s') IS NOT NULL`, prop), n.Str)
				} else {
					s.Where(fmt.Sprintf(`json_extract("raw", '$.%s') IS NOT NULL`, prop), n.Str)
				}
				keyWhere = append(keyWhere, fmt.Sprintf(`json_extract("raw", '$.%s') IS NOT NULL`, prop))
			} else {
				keyWhere = append(keyWhere, fmt.Sprintf(`json_extract("raw", '$.%s') NOT LIKE ?`, prop))
				values = append(values, any("%"+n.Str+"%"))
			}
		case "~":
			keyWhere = append(keyWhere, fmt.Sprintf(`json_extract("raw", '$.%s') LIKE ?`, prop))
			values = append(values, any("%"+n.Str+"%"))
		case "", "=":
			fallthrough
		default:
			if len(n.Str) == 0 || n.Str == vocab.NilLangRef.String() {
				keyWhere = append(keyWhere, fmt.Sprintf(`json_extract("raw", '$.%s') IS NULL`, prop))
			} else {
				keyWhere = append(keyWhere, fmt.Sprintf(`json_extract("raw", '$.%s') = ?`, prop))
				values = append(values, any(n.Str))
			}
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(keyWhere, " OR ")), values
}

func getTypeWheres(s *Stmt, f ...Check) (string, []any) {
	types := make(CompStrs, 0)
	for _, check := range f {
		if c, ok := check.(withTypes); ok {
			for _, t := range c {
				types = append(types, CompStr{Str: string(t)})
			}
		}
	}
	return getStringFieldWheres(s, "type", types...)
}

func getStringFieldWheres(s *Stmt, field string, strs ...CompStr) (string, []any) {
	if len(strs) == 0 {
		return "", nil
	}
	var values = make([]any, 0)
	keyWhere := make([]string, 0)

	stmtIsPostgres(s)
	inVal := make([]any, 0, len(strs))
	for _, t := range strs {
		switch t.Operator {
		case "!":
			if len(t.Str) == 0 || t.Str == vocab.NilLangRef.String() {
				s.Where(fmt.Sprintf(`%s IS NOT NULL`, field))
				keyWhere = append(keyWhere, fmt.Sprintf(`"%s" IS NOT NULL`, field))
			} else {
				s.Where(fmt.Sprintf(`%s NOT LIKE ?`, field), t.Str)
				keyWhere = append(keyWhere, fmt.Sprintf(`"%s" NOT LIKE ?`, field))
				values = append(values, any("%"+t.Str+"%"))
			}
		case "~":
			s.Where(fmt.Sprintf(`"%s" LIKE ?`, field), t.Str)
			keyWhere = append(keyWhere, fmt.Sprintf(`"%s" LIKE ?`, field))
			values = append(values, any("%"+t.Str+"%"))
		case "", "=":
			if len(t.Str) == 0 || t.Str == vocab.NilLangRef.String() {
				s.Where(fmt.Sprintf(`"%s" IS NULL`, field))
				keyWhere = append(keyWhere, fmt.Sprintf(`"%s" IS NULL`, field))
			} else {
				inVal = append(inVal, t.Str)
				keyWhere = append(keyWhere, fmt.Sprintf(`"%s" = ?`, field))
				values = append(values, any(t.Str))
			}
		}
	}
	if lv := len(inVal); lv > 0 {
		if lv == 1 {
			s.Where(fmt.Sprintf(`"%s" = ?`, field), inVal[0])
		} else {
			s.Where(field).In(inVal...)
		}
	}

	return fmt.Sprintf("(%s)", strings.Join(keyWhere, " OR ")), values
}

func stmtIsPostgres(s *Stmt) bool {
	sc := s.Clone()
	sc.Where("t = ?", 1)
	isPg := strings.Contains(sc.String(), "$1")
	return isPg
}

func getContextWheres(s *Stmt, f ...Check) (string, []any) {
	strs := make(CompStrs, 0)
	for _, check := range f {
		switch c := check.(type) {
		case contextEquals:
			strs = append(strs, CompStr{Str: string(c)})
		case contextLike:
			strs = append(strs, CompStr{Operator: "~", Str: string(c)})
		case contextNil:
			strs = append(strs, CompStr{Operator: "=", Str: vocab.NilIRI.String()})
		}
	}
	return getStringFieldInJSONWheres(s, "context", strs...)
}

func getURLWheres(s *Stmt, f ...Check) (string, []any) {
	strs := make(CompStrs, 0)
	for _, check := range f {
		switch c := check.(type) {
		case iriEquals:
			strs = append(strs, CompStr{Str: string(c)})
		case iriLike:
			strs = append(strs, CompStr{Operator: "~", Str: string(c)})
		case idNil:
			strs = append(strs, CompStr{Operator: "=", Str: vocab.NilIRI.String()})
		}
	}
	clause, values := getStringFieldWheres(s, "url", strs...)
	jClause, jValues := getStringFieldInJSONWheres(s, "url", strs...)
	if len(jClause) > 0 {
		if len(clause) > 0 {
			clause += " OR "
		}
		clause += jClause
	}
	values = append(values, jValues...)
	return clause, values
}

func getNamesWheres(s *Stmt, f ...Check) (string, []any) {
	strs := make(CompStrs, 0)
	for _, check := range f {
		switch c := check.(type) {
		case naturalLanguageValCheck:
			strs = append(strs, CompStr{
				Operator: "", // TODO(marius): we probably need to change the API of the naturalLanguageValCheck
				Str:      c.checkValue,
			})
		}
	}
	ns, np := getStringFieldInJSONWheres(s, "name", strs...)
	pus, pup := getStringFieldInJSONWheres(s, "preferredUsername", strs...)
	return strings.Join([]string{ns, pus}, " OR "), append(np, pup)
}

func getInReplyToWheres(s *Stmt, f ...Check) (string, []any) {
	strs := make(CompStrs, 0)
	for _, check := range f {
		switch c := check.(type) {
		case iriEquals:
			strs = append(strs, CompStr{Str: string(c)})
		}
	}
	return getStringFieldInJSONWheres(s, "inReplyTo", strs...)
}

func getAttributedToWheres(s *Stmt, f ...Check) (string, []any) {
	strs := make(CompStrs, 0)
	for _, check := range f {
		switch c := check.(type) {
		case iriEquals:
			strs = append(strs, CompStr{Str: string(c)})
		}
	}
	return getStringFieldInJSONWheres(s, "attributedTo", strs...)
}
