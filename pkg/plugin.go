package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	glog "github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

// newDatasource returns datasource.ServeOpts.
func newDatasource() datasource.ServeOpts {
	// creates a instance manager for your plugin. The function passed
	// into `NewInstanceManger` is called when the instance is created
	// for the first time or when a datasource configuration changed.
	im := datasource.NewInstanceManager(newDataSourceInstance)

	ds := &SampleDatasource{
		im: im,
	}

	return datasource.ServeOpts{
		QueryDataHandler:   ds,
		CheckHealthHandler: ds,
	}
}

// SampleDatasource is an example datasource used to scaffold
// new datasource plugins with an backend.
type SampleDatasource struct {
	// The instance manager can help with lifecycle management
	// of datasource instances in plugins. It's not a requirements
	// but a best practice that we recommend that you follow.
	im instancemgmt.InstanceManager
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifer).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (td *SampleDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	glog.DefaultLogger.Info("QueryData", "request", req)

	h, err := td.im.Get(req.PluginContext)
	if err != nil {
		return nil, err
	}

	settings := h.(*instanceSettings)

	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := td.query(ctx, settings, q)
		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

func (td *SampleDatasource) query(ctx context.Context, settings *instanceSettings, query backend.DataQuery) backend.DataResponse {
	// Parse query from JSON
	q := &Query{}

	if err := json.Unmarshal(query.JSON, q); err != nil {
		return backend.DataResponse{
			Error: err,
		}
	}
	// Send request to GraphQL API

	u := url.Values{
		"query": []string{q.Query},
	}

	queryURL, err := url.Parse(settings.URL)
	if err != nil {
		return backend.DataResponse{
			Error: err,
		}
	}

	queryURL.RawQuery = u.Encode()

	res, err := settings.httpClient.Get(queryURL.String())

	if err != nil {
		return backend.DataResponse{
			Error: err,
		}
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return backend.DataResponse{
			Error: err,
		}
	}

	dataResponse := map[string]interface{}{}

	if err := json.Unmarshal(b, &dataResponse); err != nil {
		return backend.DataResponse{
			Error: err,
		}
	}

	frames, err := FrameJSON(q, dataResponse)

	if err != nil {
		return backend.DataResponse{
			Error: err,
		}
	}

	return backend.DataResponse{
		Frames: frames,
	}
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (td *SampleDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Data source is working"
	glog.DefaultLogger.Info("CheckHealth", "request", req)

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

type instanceSettings struct {
	backend.DataSourceInstanceSettings
	httpClient *http.Client
}

func newDataSourceInstance(setting backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	return &instanceSettings{
		DataSourceInstanceSettings: setting,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}, nil
}

func (s *instanceSettings) Dispose() {
	// Called before creatinga a new instance to allow plugin authors
	// to cleanup.
}
