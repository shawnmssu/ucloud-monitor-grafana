package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"net/http"
	"sync"
	"time"
)

// Make sure UCloudDatasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler, backend.StreamHandler interfaces. Plugin should not
// implement all these interfaces - only those which are required for a particular task.
// For example if plugin does not need streaming functionality then you are free to remove
// methods that implement backend.StreamHandler. Implementing instancemgmt.InstanceDisposer
// is useful to clean up resources used by previous datasource instance when a new datasource
// instance created upon datasource settings changed.
var (
	_ backend.QueryDataHandler   = (*UCloudDatasource)(nil)
	_ backend.CheckHealthHandler = (*UCloudDatasource)(nil)
	//_ backend.StreamHandler         = (*UCloudDatasource)(nil)
	//_ instancemgmt.InstanceDisposer = (*UCloudDatasource)(nil)
	_ backend.CallResourceHandler = (*UCloudDatasource)(nil)
)

// NewUCloudDatasource creates a new datasource instance.
func NewUCloudDatasource(backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/generic_api", GenericApi)
	return &UCloudDatasource{
		callResourceHandler: httpadapter.New(mux),
	}, nil
}

// UCloudDatasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type UCloudDatasource struct {
	callResourceHandler backend.CallResourceHandler
}

func (d *UCloudDatasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	return d.callResourceHandler.CallResource(ctx, req, sender)
}

//QueryData handles multiple queries and returns multiple responses.
//req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
//The QueryDataResponse contains a map of RefID to the response for each query, and each response
//contains Frames ([]*Frame).
func (d *UCloudDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	log.DefaultLogger.Info("QueryData called", "request", req)

	// create response struct
	response := backend.NewQueryDataResponse()

	if req.PluginContext.DataSourceInstanceSettings == nil {
		return nil, fmt.Errorf("data source setting got nil")
	}
	conf, err := getUCloudConfig(*req.PluginContext.DataSourceInstanceSettings)
	if err != nil {
		return nil, fmt.Errorf("get ucloud setting got error, %s", err)
	}
	client := conf.Client()

	// loop over queries and execute them individually.
	var wg sync.WaitGroup
	var mux sync.Mutex
	for _, q := range req.Queries {
		wg.Add(1)
		go func(q backend.DataQuery) {
			res := d.query(ctx, client.ucloudconn, q)

			// save the response in a hashmap
			// based on with RefID as identifier
			mux.Lock()
			response.Responses[q.RefID] = res
			mux.Unlock()
			wg.Done()
		}(q)

	}
	wg.Wait()

	return response, nil
}

type queryModel struct {
	ProjectId    string `json:"projectId"`
	Region       string `json:"region"`
	ResourceType string `json:"resourceType"`
	MetricName   string `json:"metricName"`
	ResourceId   string `json:"resourceId"`
}

func (d *UCloudDatasource) query(_ context.Context, client *ucloud.Client, query backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}

	// Unmarshal the JSON into our queryModel.
	var qm queryModel

	response.Error = json.Unmarshal(query.JSON, &qm)
	if response.Error != nil {
		return response
	}
	reqGet := client.NewGenericRequest()
	if qm.ProjectId != "" {
		_ = reqGet.SetProjectId(qm.ProjectId)
	}
	response.Error = reqGet.SetPayload(map[string]interface{}{
		"Action":       "GetMetric",
		"Region":       qm.Region,
		"ResourceType": qm.ResourceType,
		"MetricName":   []string{qm.MetricName},
		"ResourceId":   qm.ResourceId,
		"BeginTime":    query.TimeRange.From.Unix(),
		"EndTime":      query.TimeRange.To.Unix(),
	})

	if response.Error != nil {
		return response
	}
	respGet, err := client.GenericInvoke(reqGet)
	if err != nil {
		response.Error = err
		return response
	}

	type ResponseItem struct {
		Timestamp int64
		Value     float64
	}
	type GetMetricResponse struct {
		DataSets map[string][]ResponseItem
	}

	respGetObj := GetMetricResponse{}
	response.Error = respGet.Unmarshal(&respGetObj)
	if response.Error != nil {
		return response
	}

	for metric, items := range respGetObj.DataSets {
		frame := data.NewFrame(qm.ResourceId)
		times := make([]time.Time, 0)
		values := make([]float64, 0)
		for _, v := range items {
			times = append(times, time.Unix(v.Timestamp, 0))
			values = append(values, v.Value)
		}
		frame.Fields = append(frame.Fields,
			data.NewField("time", nil, times),
			data.NewField(metric, nil, values),
		)
		response.Frames = append(response.Frames, frame)
	}

	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *UCloudDatasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	log.DefaultLogger.Info("CheckHealth called", "request", req)

	var status = backend.HealthStatusOk
	var message = "Data source is working"

	if req.PluginContext.DataSourceInstanceSettings != nil {
		setting := req.PluginContext.DataSourceInstanceSettings
		if setting.DecryptedSecureJSONData["publicKey"] == "" {
			status = backend.HealthStatusError
			message = "Public Key must be set"
		}

		if setting.DecryptedSecureJSONData["privateKey"] == "" {
			status = backend.HealthStatusError
			message = "Private Key must be set"
		}
	}
	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}
