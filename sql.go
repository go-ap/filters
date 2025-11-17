package filters

import (
	"fmt"
	"strings"

	vocab "github.com/go-ap/activitypub"
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

func GetWhereClauses(f ...Check) ([]string, []any) {
	var clauses = make([]string, 0)
	var values = make([]any, 0)

	if typClause, typValues := getTypeWheres(f...); len(typClause) > 0 {
		values = append(values, typValues...)
		clauses = append(clauses, typClause)
	}

	if iriClause, iriValues := getIRIWheres(f...); len(iriClause) > 0 {
		values = append(values, iriValues...)
		clauses = append(clauses, iriClause)
	}

	if nameClause, nameValues := getNamesWheres(f...); len(nameClause) > 0 {
		values = append(values, nameValues...)
		clauses = append(clauses, nameClause)
	}

	if replClause, replValues := getInReplyToWheres(f...); len(replClause) > 0 {
		values = append(values, replValues...)
		clauses = append(clauses, replClause)
	}

	if authorClause, authorValues := getAttributedToWheres(f...); len(authorClause) > 0 {
		values = append(values, authorValues...)
		clauses = append(clauses, authorClause)
	}

	if urlClause, urlValues := getURLWheres(f...); len(urlClause) > 0 {
		values = append(values, urlValues...)
		clauses = append(clauses, urlClause)
	}

	if ctxtClause, ctxtValues := getContextWheres(f...); len(ctxtClause) > 0 {
		values = append(values, ctxtValues...)
		clauses = append(clauses, ctxtClause)
	}

	return clauses, values
}

func getIRIWheres(f ...Check) (string, []any) {
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
	return getStringFieldWheres(strs, "iri")
}

const MaxItems int = 100

func getStringFieldInJSONWheres(strs CompStrs, props ...string) (string, []any) {
	if len(strs) == 0 {
		return "", nil
	}
	var values = make([]any, 0)
	keyWhere := make([]string, 0)
	for _, n := range strs {
		switch n.Operator {
		case "!":
			for _, prop := range props {
				if len(n.Str) == 0 || n.Str == vocab.NilLangRef.String() {
					keyWhere = append(keyWhere, fmt.Sprintf(`json_extract("raw", '$.%s') IS NOT NULL`, prop))
				} else {
					keyWhere = append(keyWhere, fmt.Sprintf(`json_extract("raw", '$.%s') NOT LIKE ?`, prop))
					values = append(values, any("%"+n.Str+"%"))
				}
			}
		case "~":
			for _, prop := range props {
				keyWhere = append(keyWhere, fmt.Sprintf(`json_extract("raw", '$.%s') LIKE ?`, prop))
				values = append(values, any("%"+n.Str+"%"))
			}
		case "", "=":
			fallthrough
		default:
			for _, prop := range props {
				if len(n.Str) == 0 || n.Str == vocab.NilLangRef.String() {
					keyWhere = append(keyWhere, fmt.Sprintf(`json_extract("raw", '$.%s') IS NULL`, prop))
				} else {
					keyWhere = append(keyWhere, fmt.Sprintf(`json_extract("raw", '$.%s') = ?`, prop))
					values = append(values, any(n.Str))
				}
			}
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(keyWhere, " OR ")), values
}

func getStringFieldWheres(strs CompStrs, fields ...string) (string, []any) {
	if len(strs) == 0 {
		return "", nil
	}
	var values = make([]any, 0)
	keyWhere := make([]string, 0)
	for _, t := range strs {
		switch t.Operator {
		case "!":
			for _, field := range fields {
				if len(t.Str) == 0 || t.Str == vocab.NilLangRef.String() {
					keyWhere = append(keyWhere, fmt.Sprintf(`"%s" IS NOT NULL`, field))
				} else {
					keyWhere = append(keyWhere, fmt.Sprintf(`"%s" NOT LIKE ?`, field))
					values = append(values, any("%"+t.Str+"%"))
				}
			}
		case "~":
			for _, field := range fields {
				keyWhere = append(keyWhere, fmt.Sprintf(`"%s" LIKE ?`, field))
				values = append(values, any("%"+t.Str+"%"))
			}
		case "", "=":
			for _, field := range fields {
				if len(t.Str) == 0 || t.Str == vocab.NilLangRef.String() {
					keyWhere = append(keyWhere, fmt.Sprintf(`"%s" IS NULL`, field))
				} else {
					keyWhere = append(keyWhere, fmt.Sprintf(`"%s" = ?`, field))
					values = append(values, any(t.Str))
				}
			}
		}
	}

	return fmt.Sprintf("(%s)", strings.Join(keyWhere, " OR ")), values
}
func getTypeWheres(f ...Check) (string, []any) {
	types := make(CompStrs, 0)
	for _, check := range f {
		if c, ok := check.(withTypes); ok {
			for _, t := range c {
				types = append(types, CompStr{Str: string(t)})
			}
		}
	}
	return getStringFieldWheres(types, "type")
}

func getContextWheres(f ...Check) (string, []any) {
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
	return getStringFieldInJSONWheres(strs, "context")
}

func getURLWheres(f ...Check) (string, []any) {
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
	clause, values := getStringFieldWheres(strs, "url")
	jClause, jValues := getStringFieldInJSONWheres(strs, "url")
	if len(jClause) > 0 {
		if len(clause) > 0 {
			clause += " OR "
		}
		clause += jClause
	}
	values = append(values, jValues...)
	return clause, values
}

func getNamesWheres(f ...Check) (string, []any) {
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
	return getStringFieldInJSONWheres(strs, "name", "preferredUsername")
}

func getInReplyToWheres(f ...Check) (string, []any) {
	strs := make(CompStrs, 0)
	for _, check := range f {
		switch c := check.(type) {
		case iriEquals:
			strs = append(strs, CompStr{Str: string(c)})
		}
	}
	return getStringFieldInJSONWheres(strs, "inReplyTo")
}

func getAttributedToWheres(f ...Check) (string, []any) {
	strs := make(CompStrs, 0)
	for _, check := range f {
		switch c := check.(type) {
		case iriEquals:
			strs = append(strs, CompStr{Str: string(c)})
		}
	}
	return getStringFieldInJSONWheres(strs, "attributedTo")
}
