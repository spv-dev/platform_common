package prettier

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// PlaceholderDollar знак доллара
	PlaceholderDollar = "$"
	// PlaceholderQuestion знак вопроса
	PlaceholderQuestion = "?"
)

// Pretty форматирует строку запроса
func Pretty(query string, placeholder string, args ...any) string {
	for i, param := range args {
		var value string
		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("'%v'", v)
		case []byte:
			value = fmt.Sprintf("'%v'", string(v))
		default:
			value = fmt.Sprintf("%v", v)
		}

		query = strings.Replace(query, fmt.Sprintf("%s%s", placeholder, strconv.Itoa(i+1)), value, -1)
	}

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")

	return strings.TrimSpace(query)
}
