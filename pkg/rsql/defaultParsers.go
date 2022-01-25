package rsql

import (
	"fmt"
	"strings"
)

var clauseStr = "{ }"

var labelList = []string{
	"==",
	"!=",
	">",
	">=",
	"<",
	"<=",
	"=gt=",
	"=ge=",
	"=lt=",
	"=le=",
	"=in=",
	"=out=",
}

// Mongo adds the default mongo operators to the parser
func Mongo() func(parser *Parser) error {
	return func(parser *Parser) error {
		for _, o := range labelList {
			var label string
			switch o {
			case ">", ">=", "<", "<=":
				continue
			case "==":
				label = "$eq"
			case "!=":
				label = "$ne"
			case "=gt=":
				label = "$gt"
			case "=ge=":
				label = "$gte"
			case "=lt=":
				label = "$lt"
			case "=le=":
				label = "$lte"
			case labelList[10]:
				label = "$in"
			case labelList[11]:
				label = "$nin"
			}

			theOperator := Operator{
				Operator: o,
				Formatter: func(key, value string) string {
					return fmt.Sprintf(`{ "%s": { "%s": %s } }`, key, label, value)
				},
			}

			if o == labelList[10] || o == labelList[11] {
				theOperator.Formatter = func(key, value string) string {
					// remove parentheses
					value = value[1 : len(value)-1]
					return fmt.Sprintf(`{ "%s": { "%s": %s } }`, key, label, value)
				}
			}

			parser.operators = append(parser.operators, theOperator)
		}

		// AND formatter
		parser.andFormatter = func(ss []string) string {
			if len(ss) > 1 {
				return fmt.Sprintf(`{ "$and": [ %s ] }`, strings.Join(ss, ", "))
			}

			if len(ss) == 0 {
				return ""
			}

			return ss[0]
		}

		// OR formatter
		parser.orFormatter = func(ss []string) string {
			if len(ss) > 1 {
				return fmt.Sprintf(`{ "$or": [ %s ] }`, strings.Join(ss, ", "))
			}

			if len(ss) == 0 {
				return clauseStr
			}

			return ss[0]
		}

		return nil
	}
}

func Mysql() func(parser *Parser) error {
	return func(parser *Parser) error {
		// operators
		for _, o := range labelList {
			var label string
			switch o {
			case "==":
				label = "="
			case "!=":
				label = o
			case "=gt=", ">":
				label = ">"
			case "=ge=", ">=":
				label = ">="
			case "=lt=", "<":
				label = "<"
			case "=le=", "<=":
				label = "<="
			case labelList[10]:
				label = "in"
			case labelList[11]:
				label = "not in"
			}

			theOperator := Operator{
				Operator: o,
				Formatter: func(key, value string) string {
					return fmt.Sprintf("`%s` %s %s", key, label, value)
				},
			}

			if o == "=in=" || o == "=out=" {
				theOperator.Formatter = func(key, value string) string {
					value = value[1 : len(value)-1]
					return fmt.Sprintf("`%s` %s (%s)", key, label, value)
				}
			}

			parser.operators = append(parser.operators, theOperator)
		}

		// AND formatter
		parser.andFormatter = mysqlAndFormatter

		// OR formatter
		parser.orFormatter = mysqlOrFormatter

		return nil
	}
}

func MysqlPre(nameTransfer func(s string) string) func(parser *PreParser) error {
	return func(parser *PreParser) error {
		// operators
		for _, o := range labelList {
			var label string
			switch o {
			case labelList[0]:
				label = "="
			case labelList[1]:
				label = o
			case labelList[6], labelList[2]:
				label = ">"
			case labelList[7], labelList[3]:
				label = ">="
			case labelList[8], labelList[4]:
				label = "<"
			case labelList[9], labelList[5]:
				label = "<="
			case labelList[10]:
				label = "in"
			case labelList[11]:
				label = "not in"
			}

			parser.operators = append(parser.operators, PreOperator{
				Operator: o,
				Formatter: func(key string) string {
					return fmt.Sprintf("`%s` %s ?", nameTransfer(key), label)
				},
			})
		}

		// AND formatter
		parser.andFormatter = mysqlAndFormatter

		// OR formatter
		parser.orFormatter = mysqlOrFormatter

		return nil
	}
}

func mysqlAndFormatter(ss []string) string {
	if len(ss) > 1 {
		return fmt.Sprintf(`(%s)`, strings.Join(ss, " and "))
	}

	if len(ss) == 0 {
		return ""
	}

	return ss[0]
}

func mysqlOrFormatter(ss []string) string {
	if len(ss) > 1 {
		return fmt.Sprintf(`(%s)`, strings.Join(ss, " or "))
	}

	if len(ss) == 0 {
		return clauseStr
	}

	return ss[0]
}
