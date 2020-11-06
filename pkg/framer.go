package main

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/itchyny/gojq"
)

var (
	ErrorUnrecognizedType     = errors.New("unrecognized type")
	ErrorUnsupportedOperation = errors.New("unsupported operation")
)

func mapToFields(m map[string]interface{}) data.Fields {
	f := data.Fields{}

	for k, v := range m {
		f = append(f, data.NewField(k, nil, v))
	}

	return f
}

// {
//   rows: ["time", "value", "name"]
//   data: {
//      "time": [100, 110, 111],
//      "value": [0, 1, 2],
//      ..
//    }
// }
func frameJqValue(value interface{}) (interface{}, error) {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Map:
		return nil, fmt.Errorf("Nested maps are unsupported: %w", ErrorUnsupportedOperation)
	case reflect.Slice:
		return nil, fmt.Errorf("Nested arrays are unsupported: %w", ErrorUnsupportedOperation)
	}

	switch v := value.(type) {
	case float64:
		return v, nil
	case string:
		return v, nil
	case int64:
		return v, nil
	default:
		return nil, fmt.Errorf("type: %s. error: %w", reflect.TypeOf(value).Kind(), ErrorUnrecognizedType)
	}
}

func frameJqResult(value interface{}) (map[string][]interface{}, error) {
	data := map[string][]interface{}{}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Map:
		// If the result is a map, then we need to traverse through each key
		// and append the results to the field that corresponds with that key.
		// Nested maps are unsupported and will return an error from frameJqValue
		m := value.(map[string]interface{})

		for k, v := range m {
			if v == nil {
				data[k] = append(data[k], "")
				continue
			}
			val, err := frameJqValue(v)
			if err != nil {
				return nil, err
			}

			data[k] = append(data[k], val)
		}
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		// If the data is a slice, it could be a slice of single primitive values, like:
		// []float64{0.1, 0.2, 0.3, ...}
		// Or it could be a slice of maps (more common):
		// []map[string]string{{"country": "USA"}, {"country": "UK"}, {"country": "Spain"}, ...}
		slice := value.([]interface{})
		for _, v := range slice {
			f, err := frameJqResult(v)
			if err != nil {
				return nil, err
			}

			for k, v := range f {
				data[k] = append(data[k], v...)
			}
		}
	default:
		v, err := frameJqValue(value)
		if err != nil {
			return nil, err
		}

		data["value"] = []interface{}{v}
	}

	return data, nil
}

func jqFields(value interface{}) (data.Fields, error) {
	results, err := frameJqResult(value)
	if err != nil {
		return nil, err
	}

	// In the results map, each field represents a key
	// and each key contains a list...

	// Assumptions:
	// * Each field in the list is of the same type
	// * Each field in the list is of the same length
	// * There are no "nil" items in the list. (?)

	fields := make(data.Fields, len(results))
	i := 0
	for k, v := range results {
		var field *data.Field
		if len(v) == 0 {
			return nil, errors.New("0 length results")
		}
		// TODO: Should labels be supported via some nested map?
		switch v[0].(type) {
		case string:
			field = data.NewField(k, nil, make([]string, len(v)))
		case int64:
			field = data.NewField(k, nil, make([]int64, len(v)))
		case time.Time:
			field = data.NewField(k, nil, make([]*time.Time, len(v)))
		case *time.Time:
			field = data.NewField(k, nil, make([]*time.Time, len(v)))
		case float64:
			field = data.NewField(k, nil, make([]float64, len(v)))
		default:
			return nil, fmt.Errorf("type: %s. error: %w", reflect.TypeOf(v).Kind(), ErrorUnrecognizedType)
		}

		for i, v := range v {
			field.Set(i, v)
		}

		fields[i] = field
		i++
	}

	return fields, nil
}

// FrameJSON takes a graphql query response and turns it into a dataframe
// query {
//   countries {
//     name
//   }
// }
// Every requested field should represent a single field, and data that corresponds to that field shoudl be appended to that field
func FrameJSON(query *Query, res map[string]interface{}) (data.Frames, error) {
	frameMap := map[string]*data.Frame{}

	for _, field := range query.Fields {
		if field.Name == "" {
			continue
		}

		// fields := data.Fields{}
		q, err := gojq.Parse(field.JQ)
		if err != nil {
			return nil, err
		}

		fields := data.Fields{}

		iter := q.Run(res)
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}

			if err, ok := v.(error); ok {
				return nil, err
			}

			f, err := jqFields(v)
			if err != nil {
				return nil, err
			}

			fields = append(fields, f...)
		}

		frameMap[field.Name] = data.NewFrame(field.Name, fields...)
	}

	frames := data.Frames{}

	for _, v := range frameMap {
		frames = append(frames, v)
	}

	return frames, nil
}
