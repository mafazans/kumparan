package common

import (
	"bytes"
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const maxLimit, defaultLimit int64 = 50, 10

type SQLClauseBuilder struct {
	paramTag    string
	colTag      string
	suffixQuery string
	values      map[string]reflect.Value
	page        int64
	limit       int64
}

func NewSQLClauseBuilder(paramTag, colTag string, suffix string, page, limit int64) *SQLClauseBuilder {
	return &SQLClauseBuilder{
		paramTag:    paramTag,
		colTag:      colTag,
		suffixQuery: suffix,
		values:      make(map[string]reflect.Value),
		page:        page,
		limit:       limit,
	}
}

func (qb *SQLClauseBuilder) AliasPrefix(alias string, ptr interface{}) *SQLClauseBuilder {
	p := reflect.ValueOf(ptr)
	if p.Kind() != reflect.Ptr {
		panic(errors.New("passed interface{} should be a pointer"))
	}
	v := p.Elem()
	qb.values[alias] = v
	return qb
}

const (
	unknown = iota
	one
	many
	like
	lte
	lt
	gte
	gt
)

func (qb *SQLClauseBuilder) Build() (string, []string, []interface{}, error) {
	var (
		args   []interface{}
		sortBy []string
	)

	// sortByDisplay needs to be initialized to return empty array when empty
	sortByDisplay := []string{}

	mapDBcolsByParam := make(map[string]string)
	buff := bytes.NewBufferString("")

	// build map by paramTag where the value is colDBTag value
	for table, v := range qb.values {
		var alias string
		if table == "-" || len(table) < 1 {
			alias = ""
		} else {
			alias = table + "."
		}
		for i := 0; i < v.NumField(); i++ {
			tag := v.Type().Field(i).Tag
			if tag.Get(qb.paramTag) != "-" && tag.Get(qb.paramTag) != "" {
				mapDBcolsByParam[alias+tag.Get(qb.paramTag)] = tag.Get(qb.colTag)
			}
		}
	}

	buff.WriteString(" WHERE 1=1")
	if len(qb.suffixQuery) > 0 {
		buff.WriteString(" AND " + qb.suffixQuery)
	}

	for table, v := range qb.values {
		var alias string
		if table == "-" || len(table) < 1 {
			alias = ""
		} else {
			alias = table + "."
		}
		// iterate over the structs to get the value
		for i := 0; i < v.NumField(); i++ {
			var arg interface{}

			// if only the struct declares colTag we will build the filter query
			// if not then skip
			tag := v.Type().Field(i).Tag
			colTag := tag.Get(qb.colTag)
			if colTag == "" || colTag == "-" {
				continue
			}

			vFieldItf := v.Field(i).Interface()
			qType := unknown
			skip := true
			switch f := vFieldItf.(type) {
			case []int64:
				if len(f) > 0 {
					arg = f
					qType = many
					skip = false
				}
			case []string:
				if len(f) > 0 {
					arg = f
					qType = many
					skip = false
				}
			case []float64:
				if len(f) > 0 {
					arg = f
					qType = many
					skip = false
				}
			case []bool:
				if len(f) > 0 {
					arg = f
					qType = many
					skip = false
				}
			case []time.Time:
				if len(f) > 0 {
					arg = f
					qType = many
					skip = false
				}
			case int64:
				if f > 0 {
					arg = f
					qType = qb.getOperator(tag.Get(qb.paramTag))
					skip = false
				}
			case string:
				if len(f) > 0 {
					arg = f
					qType = one
					skip = false
					if strings.Contains(f, "%") {
						qType = like
					}
				}
			case float64:
				if f > 0 {
					arg = f
					qType = qb.getOperator(tag.Get(qb.paramTag))
					skip = false
				}
			case bool:
				if f {
					arg = f
					qType = one
					skip = false
				}
			case sql.NullBool:
				if f.Valid {
					arg = f.Bool
					qType = one
					skip = false
				}
			case sql.NullInt64:
				if f.Valid {
					arg = f.Int64
					qType = qb.getOperator(tag.Get(qb.paramTag))
					skip = false
				}
			case sql.NullFloat64:
				if f.Valid {
					arg = f.Float64
					qType = qb.getOperator(tag.Get(qb.paramTag))
					skip = false
				}
			case sql.NullString:
				if f.Valid {
					arg = f.String
					qType = one
					skip = false
				}
			case time.Time:
				if !f.Equal(time.Time{}) {
					arg = f
					qType = qb.getOperator(tag.Get(qb.paramTag))
					skip = false
				}
			case sql.NullTime:
				if f.Valid {
					arg = f.Time
					qType = qb.getOperator(tag.Get(qb.paramTag))
					skip = false
				}
			default:
				continue
			}

			if !skip {
				switch colTag {
				case "sortby", "orderby", "sort_by", "order_by", "sort-by", "order-by":
					if v, ok := arg.([]string); ok {
						for _, s := range v {
							// put validation on regex
							reg := regexp.MustCompile("(?P<sign>-)?(?P<col>[a-zA-Z_]+),?")
							if reg.MatchString(s) {
								for _, _s := range strings.Split(s, ",") {
									var col string
									sort := "asc"
									match := reg.FindStringSubmatch(_s)
									for i, name := range reg.SubexpNames() {
										if i == 0 || name == "" {
											continue
										}
										if match != nil {
											if name == "sign" && match[i] == "-" {
												sort = "desc"
											} else if name == "col" {
												// match param tag to db column tag
												if _col, ok := mapDBcolsByParam[alias+match[i]]; ok {
													// col will be added with table alias name
													col = alias + _col
													// sortByDisplay contains values displayed to users
													sortByDisplay = append(sortByDisplay, alias+match[i]+" "+sort)
												}
											}
										}
									}
									// append to sort array
									if col != "" {
										sortBy = append(sortBy, col+" "+sort)
									}
								}
							}
						}
					}
					continue

				case "page", "size", "limit", "offset":
					//skip on tag field with these values
					continue

				default:
					switch qType {
					case one:
						buff.WriteString(" AND " + alias + colTag + "=?")
					case gte:
						buff.WriteString(" AND " + alias + colTag + ">=?")
					case gt:
						buff.WriteString(" AND " + alias + colTag + ">?")
					case lte:
						buff.WriteString(" AND " + alias + colTag + "<=?")
					case lt:
						buff.WriteString(" AND " + alias + colTag + "<?")
					case like:
						buff.WriteString(" AND " + alias + colTag + " LIKE ?")
					case many:
						buff.WriteString(" AND " + alias + colTag + " IN (?)")
					default:
						return buff.String(), sortByDisplay, args, errors.New(`unknown query type`)
					}
					args = append(args, arg)
				}
			}
		}
	}

	// handle sort
	if len(sortBy) > 0 {
		buff.WriteString(" ORDER BY " + strings.Join(sortBy, ", "))
	}

	// validation
	qb.limit = ValidateLimit(qb.limit)
	qb.page = ValidatePage(qb.page)

	// build limit offset query
	if qb.page > 0 || qb.limit > 0 {
		var offsetStr string
		offset := getOffset(qb.page, qb.limit)
		offsetStr = strconv.FormatInt(offset, 10)
		limitStr := strconv.FormatInt(qb.limit, 10)
		buff.WriteString(" LIMIT " + offsetStr + ", " + limitStr)
	}

	buff.WriteString(";")
	return buff.String(), sortByDisplay, args, nil
}

func ValidateLimit(l int64) int64 {
	if l < 1 {
		return defaultLimit
	} else if l > maxLimit {
		return maxLimit
	}
	return l
}

func ValidatePage(p int64) int64 {
	if p < 1 {
		return 1
	}
	return p
}

func getOffset(p, l int64) int64 {
	return (p - 1) * l
}

func (qb *SQLClauseBuilder) getOperator(paramTag string) int {
	if strings.Contains(paramTag, "__gte") {
		return gte
	} else if strings.Contains(paramTag, "__lte") {
		return lte
	} else if strings.Contains(paramTag, "__lt") {
		return lt
	} else if strings.Contains(paramTag, "__gt") {
		return gt
	}
	return one
}
