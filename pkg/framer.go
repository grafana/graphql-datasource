package main

import (
	"fmt"
	"reflect"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// {
//   rows: ["time", "value", "name"]
//   data: {
//      "time": [100, 110, 111],
//      "value": [0, 1, 2],
//      ..
//    }
// }

func frame(value interface{}) (data.Fields, error) {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Float64:
		log.DefaultLogger.Info("this is a float")
	case reflect.String:
		log.DefaultLogger.Info("this is a string")
	case reflect.Bool:
		log.DefaultLogger.Info("this is a bool")
	case reflect.Map:
		return frame(value)
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		log.DefaultLogger.Info("this is an array / slice")
	default:
		return nil, fmt.Errorf("unrecognized type %s / %v / %v", reflect.TypeOf(value).Kind(), reflect.TypeOf(value), value)
	}
	return nil, nil
}

// FrameJSON takes a graphql query response and turns it into a dataframe
// query {
//   countries {
//     name
//   }
// }
// Every requested field should represent a single field, and data that corresponds to that field shoudl be appended to that field
func FrameJSON(query string, res map[string]interface{}) (data.Frames, error) {
	frames := data.Frames{}

	res = res["data"].(map[string]interface{})

	for k, v := range res {
		fields, err := frame(v)
		if err != nil {
			return nil, err
		}

		frames = append(frames, data.NewFrame(k, fields...))
	}

	return frames, nil
}
