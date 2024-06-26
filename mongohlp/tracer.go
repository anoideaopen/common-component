package mongohlp

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	findCmdName      = "find"
	aggregateCmdName = "aggregate"
)

type findCmd struct {
	Find   string                 `json:"find"`
	Filter map[string]interface{} `json:"filter"`
	Limit  map[string]interface{} `json:"limit"`
	Sort   map[string]interface{} `json:"sort"`
}

type aggregateCmd struct {
	Aggregate string        `json:"aggregate"`
	Pipeline  []interface{} `json:"pipeline"`
}

// TraceMongoCmd traces mongo command
func TraceMongoCmd(rawCmd string) (string, error) {
	var cmd map[string]interface{}
	if err := json.Unmarshal([]byte(rawCmd), &cmd); err != nil {
		return "", fmt.Errorf("error unmarshal cmd: %w", err)
	}

	if _, ok := cmd[findCmdName]; ok {
		return traceFindCmd(rawCmd)
	}

	if _, ok := cmd[aggregateCmdName]; ok {
		return traceAggregateCmd(rawCmd)
	}

	return "", nil
}

func traceAggregateCmd(rawCmd string) (string, error) {
	var cmd aggregateCmd
	if err := json.Unmarshal([]byte(rawCmd), &cmd); err != nil {
		return "", fmt.Errorf("error unmarshal aggregateCmd: %w", err)
	}

	var err error
	for i, v := range cmd.Pipeline {
		cmd.Pipeline[i], err = normalizeMongoCmd(v)
		if err != nil {
			return "", err
		}
	}

	b, err := json.Marshal(cmd.Pipeline)
	if err != nil {
		return "", fmt.Errorf("error unmarshal pipeline: %w", err)
	}

	sb := &strings.Builder{}
	sb.WriteString("db.")
	sb.WriteString(cmd.Aggregate)
	sb.WriteString(".aggregate(")
	sb.Write(b)
	sb.WriteString(")")

	return sb.String(), nil
}

func traceFindCmd(rawCmd string) (string, error) {
	var cmd findCmd
	if err := json.Unmarshal([]byte(rawCmd), &cmd); err != nil {
		return "", fmt.Errorf("error unmarshal findCmd: %w", err)
	}
	filterStr := "{}"
	if cmd.Filter != nil {
		o, err := normalizeMongoObj(cmd.Filter)
		if err != nil {
			return "", err
		}
		b, err := json.Marshal(o)
		if err != nil {
			return "", fmt.Errorf("error marshal filter object: %w", err)
		}
		filterStr = string(b)
	}

	sortStr := ""
	if cmd.Sort != nil {
		sort, err := normalizeMongoObj(cmd.Sort)
		if err != nil {
			return "", err
		}
		if sort != nil {
			b, err := json.Marshal(sort)
			if err != nil {
				return "", fmt.Errorf("error marshal sort object: %w", err)
			}
			sortStr = fmt.Sprintf(".sort( %s )", string(b))
		}
	}

	limitStr := ""
	if cmd.Limit != nil {
		limit, err := normalizeMongoObj(cmd.Limit)
		if err != nil {
			return "", err
		}
		if limit != nil {
			limitStr = fmt.Sprintf(".limit(%v)", limit)
		}
	}

	sb := &strings.Builder{}
	sb.WriteString("db.")
	sb.WriteString(cmd.Find)
	sb.WriteString(".find(")
	sb.WriteString(filterStr)
	sb.WriteString(")")
	sb.WriteString(sortStr)
	sb.WriteString(limitStr)

	return sb.String(), nil
}

func tryToAsNumber(obj map[string]interface{}) (float64, bool, error) {
	if len(obj) != 1 {
		return 0, false, nil
	}
	for _, numKey := range []string{"$numberLong", "$numberInt", "$numberDouble"} {
		if v, ok := obj[numKey]; ok {
			if s, ok := v.(string); ok {
				const float64BitSize = 64
				f, err := strconv.ParseFloat(s, float64BitSize)
				if err != nil {
					return 0, false, fmt.Errorf("error parse num: %w", err)
				}
				return f, true, nil
			}
		}
	}
	return 0, false, nil
}

func normalizeMongoCmd(raw interface{}) (interface{}, error) {
	obj, ok := raw.(map[string]interface{})
	if ok {
		return normalizeMongoObj(obj)
	}
	arr, ok := raw.([]interface{})
	if ok {
		return normalizeMongoArr(arr)
	}

	return raw, nil
}

func normalizeMongoArr(arr []interface{}) (interface{}, error) {
	var err error
	for i, v := range arr {
		arr[i], err = normalizeMongoCmd(v)
		if err != nil {
			return nil, err
		}
	}
	return arr, nil
}

func normalizeMongoObj(obj map[string]interface{}) (interface{}, error) {
	n, ok, err := tryToAsNumber(obj)
	if err != nil {
		return nil, err
	}
	if ok {
		return n, nil
	}

	for k, v := range obj {
		obj[k], err = normalizeMongoCmd(v)
		if err != nil {
			return nil, err
		}
	}
	return obj, nil
}
