package DirectumLogConverter

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Converter struct {
	FixedWidth bool
}

func NewConverter(fixedWidth bool) *Converter {
	return &Converter{fixedWidth}
}

type fieldWidth struct {
	Name  string
	Width int
}

var defaultFieldWidth = []fieldWidth{
	{Name: "pid", Width: 10},
	{Name: "l", Width: 5},
	{Name: "lg", Width: 30},
}

func (c *Converter) Convert(line *LogLine) *LogEntry {
	var result = LogEntry{}
	if line.Elements == nil {
		result.Elements = append(result.Elements, string(line.Raw))
		return &result
	}
	for _, element := range line.Elements {
		var key = element.Key.(string)
		var s string
		if key == "ex" {
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
			for _, element := range defaultFieldWidth {
				if key == element.Name {
					s = FitWidth(s, element.Width)
				}
			}
		}
		result.Elements = append(result.Elements, s)
	}
	return &result
}

func convertObject(object interface{}, depth int) string {
	if m, ok := object.(map[string]interface{}); ok {
		if depth <= 0 {
			j, _ := json.Marshal(object)
			return string(j)
		}
		var s []string
		for key, value := range m {
			s = append(s, fmt.Sprintf("%s=\"%s\"", key, convertObject(value, depth-1)))
		}
		return strings.Join(s, ", ")
	} else {
		if f, ok := object.(float64); ok {
			if f == float64(int64(f)) {
				return fmt.Sprintf("%.0f", f)
			} else {
				return fmt.Sprintf("%.2f", f)
			}
		}
		return fmt.Sprintf("%v", object)
	}
}

func convertExObject(object interface{}) []string {
	var result []string
	if ex, ok := object.(map[string]interface{}); ok {
		var exType = ex["type"]
		var s = fmt.Sprintf("%v", exType)
		var exMsg = ex["m"]
		if exMsg != nil {
			s += fmt.Sprintf(": %v", exMsg)
		}
		result = append(result, s)
		var exStack = ex["stack"]
		if exStack != nil {
			result = append(result, strings.Split(strings.ReplaceAll(fmt.Sprintf("   %v", exStack), "\r\n", "\n"), "\n")...)
		}
	}
	return result
}