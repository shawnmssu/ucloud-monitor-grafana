package plugin

import (
	"encoding/json"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"net/http"
)

const (
	ResourceTypeUHost = "uhost"
	ResourceTypeEIP   = "eip"
	ResourceTypeULB   = "ulb"
	ResourceTypeUDB   = "udb"
	ResourceTypeUMem  = "umem"
)

func GenericApi(rw http.ResponseWriter, req *http.Request) {
	//parse param map
	params, err := parseRequestParams(req)
	if err != nil {
		handleResponse(rw, nil, err)
		return
	}
	datasourceInstanceSettings := httpadapter.PluginConfigFromContext(req.Context()).DataSourceInstanceSettings
	conf, err := getUCloudConfig(*datasourceInstanceSettings)
	if err != nil {
		handleResponse(rw, nil, fmt.Errorf("get ucloud setting got error, %s", err))
		return
	}
	client := conf.Client()

	// redirect by action
	switch params["Action"] {
	case "GetResourceId":
		switch params["ResourceType"] {
		case ResourceTypeUHost:
			client.proxyDescribeUHostInstance(params, rw)
			break
		case ResourceTypeEIP:
			client.proxyDescribeEIP(params, rw)
			break
		case ResourceTypeULB:
			client.proxyDescribeULB(params, rw)
			break
		case ResourceTypeUDB:
			client.proxyDescribeUDBInstance(params, rw)
			break
		case ResourceTypeUMem:
			client.proxyDescribeUMem(params, rw)
			break
		}
		break
	case "GetMetricName":
		client.proxyDescribeResourceMetric(params, rw)
		break
	case "GetProjectId":
		client.proxyGetProjectList(params, rw)
		break
	case "GetRegion":
		client.proxyGetRegion(params, rw)
		break
	case "GetResourceType":
		client.proxyResourceType(params, rw)
	}
}

func (client *uCloudClient) proxyResourceType(params map[string]string, rw http.ResponseWriter) {
	var ids = []string{
		ResourceTypeUHost,
		ResourceTypeEIP,
		ResourceTypeULB,
		ResourceTypeUDB,
		ResourceTypeUMem,
	}
	d, err := json.Marshal(ids)
	log.DefaultLogger.Debug(string(d))
	handleResponse(rw, d, err)
}

func (client *uCloudClient) proxyGetRegion(params map[string]string, rw http.ResponseWriter) {
	request := client.uaccountconn.NewGetRegionRequest()

	response, err := client.uaccountconn.GetRegion(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.Regions {
			ids = append(ids, instance.Region)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}

func (client *uCloudClient) proxyGetProjectList(params map[string]string, rw http.ResponseWriter) {
	request := client.uaccountconn.NewGetProjectListRequest()

	response, err := client.uaccountconn.GetProjectList(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.ProjectSet {
			ids = append(ids, instance.ProjectId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}

func (client *uCloudClient) proxyDescribeResourceMetric(params map[string]string, rw http.ResponseWriter) {
	request := client.ucloudconn.NewGenericRequest()

	if v, ok := params["ProjectId"]; ok {
		_ = request.SetProjectId(v)
	}

	if v, ok := params["Region"]; ok {
		_ = request.SetRegion(v)
	}

	var resourceType string
	if v, ok := params["ResourceType"]; ok {
		resourceType = v
	} else {
		const message = "must set ResourceType"
		log.DefaultLogger.Error(fmt.Sprintf(message))
		handleResponse(rw, nil, fmt.Errorf(message))
		return
	}

	err := request.SetPayload(map[string]interface{}{
		"Action":       "DescribeResourceMetric",
		"ResourceType": resourceType,
	})
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
		return
	}
	resp, err := client.ucloudconn.GenericInvoke(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
		return
	}

	type ResponseItem struct {
		MetricName string
	}
	type DescribeResourceMetricResponse struct {
		DataSet []ResponseItem
	}

	respObj := DescribeResourceMetricResponse{}
	err = resp.Unmarshal(&respObj)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
		return
	}

	var names []string
	for _, instance := range respObj.DataSet {
		names = append(names, instance.MetricName)
	}

	d, err := json.Marshal(names)
	log.DefaultLogger.Debug(string(d))
	handleResponse(rw, d, err)
}

func (client *uCloudClient) proxyDescribeUHostInstance(params map[string]string, rw http.ResponseWriter) {
	request := client.uhostconn.NewDescribeUHostInstanceRequest()

	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}

	response, err := client.uhostconn.DescribeUHostInstance(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.UHostSet {
			ids = append(ids, instance.UHostId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}

func (client *uCloudClient) proxyDescribeEIP(params map[string]string, rw http.ResponseWriter) {
	request := client.unetconn.NewDescribeEIPRequest()

	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}

	response, err := client.unetconn.DescribeEIP(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.EIPSet {
			ids = append(ids, instance.EIPId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}
func (client *uCloudClient) proxyDescribeULB(params map[string]string, rw http.ResponseWriter) {
	request := client.ulbconn.NewDescribeULBRequest()

	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}

	response, err := client.ulbconn.DescribeULB(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.DataSet {
			ids = append(ids, instance.ULBId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}
func (client *uCloudClient) proxyDescribeUDBInstance(params map[string]string, rw http.ResponseWriter) {
	request := client.udbconn.NewDescribeUDBInstanceRequest()

	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}

	response, err := client.udbconn.DescribeUDBInstance(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.DataSet {
			ids = append(ids, instance.DBId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}

func (client *uCloudClient) proxyDescribeUMem(params map[string]string, rw http.ResponseWriter) {
	request := client.umemconn.NewDescribeUMemRequest()

	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}

	response, err := client.umemconn.DescribeUMem(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.DataSet {
			ids = append(ids, instance.ResourceId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}

func handleResponse(rw http.ResponseWriter, data []byte, err error) {
	if err != nil {
		rw.Write([]byte(err.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(data)
	}
}

func parseRequestParams(req *http.Request) (map[string]string, error) {
	result := map[string]string{}
	for k, values := range req.URL.Query() {
		if len(values) > 0 {
			result[k] = values[0]
		}
	}
	d, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	log.DefaultLogger.Debug("request_params: ", string(d))
	_, hasAction := result["Action"]
	if !hasAction {
		return result, fmt.Errorf("Action parameter is missing")
	}
	return result, nil
}
