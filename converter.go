package DirectumLogConverter

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"strings"
)

type Converter struct {
	FixedWidth bool
}

func NewConverter(fixedWidth bool) *Converter {
	return &Converter{fixedWidth}
}

type logElementWidth struct {
	Name  string
	Width int
}

var defaultLogElementWidth = []logElementWidth{
	{Name: "pid", Width: 10},
	{Name: "l", Width: 5},
	{Name: "lg", Width: 30},
}

func (c *Converter) Convert(line *LogLine) *LogEntry {
	var result = LogEntry{}
	if line.Elements == nil {
		result.Elements = append(result.Elements, LogElement{Value: string(line.Raw)})
		return &result
	}
	for _, element := range *line.Elements {
		var key = element.Key
		var s string
		if key == "ex" || key == "st" {
			result.AdditionalElements = convertExObject(element.Value)
		} else {
			s = convertObject(element.Value, 1)
		}
		if s == "" {
			continue
		}
		if key == "span" {
			s = "Span(" + s + ")"
		}
		if key == "args" {
			s = "(" + s + ")"
		}
		if key == "cust" {
			s = "[" + s + "]"
		}
		if c.FixedWidth {
			for _, element := range defaultLogElementWidth {
				if key == element.Name {
					s = FitWidth(s, element.Width)
					break
				}
			}
		}
		result.Elements = append(result.Elements, LogElement{key, s})
	}
	return &result
}

func convertObject(object interface{}, depth int) string {
	if m, ok := object.(map[string]interface{}); ok {
		if depth <= 0 {
			var json = jsoniter.ConfigDefault
			j, _ := json.Marshal(object)
			return string(j)
		}
		var sb strings.Builder
		firstElement := true
		for key, value := range m {
			if firstElement {
				firstElement = false
			} else {
				sb.WriteString(", ")
			}
			sb.WriteString(key)
			sb.WriteString("=\"")
			sb.WriteString(convertObject(value, depth-1))
			sb.WriteRune('"')
		}
		return sb.String()
	} else {
		if s, ok := object.(string); ok {
			return s
		}
		if b, ok := object.(bool); ok {
			return strconv.FormatBool(b)
		}
		if f, ok := object.(float64); ok {
			i := int64(f)
			if f == float64(i) {
				return strconv.FormatInt(i, 10)
			} else {
				fmt.Sprintf("%.2f", f)
			}
		}
		return fmt.Sprintf("%v", object)
	}
}

func convertExObject(object interface{}) []string {
	var result []string
	if ex, ok := object.(map[string]interface{}); ok {
		var exType = ex["type"]
		var sb strings.Builder
		sb.WriteString(exType.(string))
		var exMsg = ex["m"]
		if exMsg != nil {
			sb.WriteString(": ")
			sb.WriteString(exMsg.(string))
		}
		result = append(result, sb.String())
		var exStack = ex["stack"]
		if exStack != nil {
			result = append(result, strings.Split(strings.ReplaceAll(fmt.Sprintf("   %v", exStack), "\r\n", "\n"), "\n")...)
		}
	}
	if s, ok := object.(string); ok {
		result = append(result, s)
	}
	return result
}
