package plugin

import (
	"encoding/json"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"net/http"
	"strconv"
)

const (
	ResourceTypeUHost      = "uhost"
	ResourceTypeEIP        = "eip"
	ResourceTypeULB        = "ulb"
	ResourceTypeUDB        = "udb"
	ResourceTypeUMem       = "umem"
	ResourceTypeUDPN       = "udpn"
	ResourceTypePHost      = "phost"
	ResourceTypeShareBW    = "sharebandwidth"
	ResourceTypeUMemCache  = "umemcache"
	ResourceTypeURedis     = "uredis"
	ResourceTypeNatGW      = "natgw"
	ResourceTypeUFile      = "ufile"
	ResourceTypeULBVServer = "ulb-vserver"
	ResourceTypeUDisk      = "udisk"
	ResourceTypeUDiskSSD   = "udisk_ssd"
	ResourceTypeUDiskRSSD  = "udisk_rssd"
	ResourceTypeUDiskSys   = "udisk_sys"

	//todo
	//ResourceTypeUHadoopHost = "uhadoop_host"
	//ResourceTypeUKafkaHost  = "ukafka_host"
	//ResourceTypeUdwNode     = "udw_node"
	//ResourceTypeUHadoop     = "uhadoop"
	//ResourceTypeUKafka      = "ukafka"
	//ResourceTypeUdw         = "udw"
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
		case ResourceTypeUDPN:
			client.proxyDescribeUDPN(params, rw)
			break
		case ResourceTypePHost:
			client.proxyDescribePHost(params, rw)
			break
		case ResourceTypeShareBW:
			client.proxyDescribeShareBW(params, rw)
			break
		case ResourceTypeUMemCache:
			client.proxyDescribeUMemCache(params, rw)
			break
		case ResourceTypeURedis:
			client.proxyDescribeURedis(params, rw)
			break
		case ResourceTypeNatGW:
			client.proxyDescribeNatGW(params, rw)
			break
		case ResourceTypeUFile:
			client.proxyDescribeUFile(params, rw)
			break
		case ResourceTypeULBVServer:
			client.proxyDescribeULBVServer(params, rw)
			break
		case ResourceTypeUDisk:
			client.proxyDescribeUDisk(params, rw)
			break
		case ResourceTypeUDiskSSD:
			client.proxyDescribeUDiskSSD(params, rw)
			break
		case ResourceTypeUDiskRSSD:
			client.proxyDescribeUDiskRSSD(params, rw)
			break
		case ResourceTypeUDiskSys:
			client.proxyDescribeUDiskSys(params, rw)
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
		ResourceTypeUDPN,
		ResourceTypePHost,
		ResourceTypeShareBW,
		ResourceTypeUMemCache,
		ResourceTypeURedis,
		ResourceTypeNatGW,
		ResourceTypeUFile,
		ResourceTypeULBVServer,
		ResourceTypeUDisk,
		ResourceTypeUDiskSSD,
		ResourceTypeUDiskRSSD,
		ResourceTypeUDiskSys,
	}
	d, err := json.Marshal(ids)
	log.DefaultLogger.Debug(string(d))
	handleResponse(rw, d, err)
}

func (client *uCloudClient) proxyDescribeULBVServer(params map[string]string, rw http.ResponseWriter) {
	request := client.ulbconn.NewDescribeVServerRequest()
	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}
	if v, ok := params["ULBId"]; ok {
		request.ULBId = ucloud.String(v)
	}
	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		request.Limit = ucloud.Int(limit)
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		request.Offset = ucloud.Int(offset)
	}

	response, err := client.ulbconn.DescribeVServer(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.DataSet {
			// todo
			//if tag, ok := params["Tag"]; ok {
			//	if instance.Tag != tag {
			//		continue
			//	}
			//}
			ids = append(ids, instance.VServerId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}

func (client *uCloudClient) proxyDescribeUDisk(params map[string]string, rw http.ResponseWriter) {
	request := client.udiskconn.NewDescribeUDiskRequest()
	request.DiskType = ucloud.String("DataDisk")
	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}
	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		request.Limit = ucloud.Int(limit)
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		request.Offset = ucloud.Int(offset)
	}

	response, err := client.udiskconn.DescribeUDisk(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.DataSet {
			if tag, ok := params["Tag"]; ok {
				if instance.Tag != tag {
					continue
				}
			}
			ids = append(ids, instance.UDiskId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}

func (client *uCloudClient) proxyDescribeUDiskSSD(params map[string]string, rw http.ResponseWriter) {
	request := client.udiskconn.NewDescribeUDiskRequest()
	request.ProtocolVersion = ucloud.Int(1)
	request.IsBoot = ucloud.String("False")
	request.DiskType = ucloud.String("CLOUD_SSD")
	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}
	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		request.Limit = ucloud.Int(limit)
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		request.Offset = ucloud.Int(offset)
	}

	response, err := client.udiskconn.DescribeUDisk(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.DataSet {
			if tag, ok := params["Tag"]; ok {
				if instance.Tag != tag {
					continue
				}
			}
			ids = append(ids, instance.UDiskId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}
func (client *uCloudClient) proxyDescribeUDiskRSSD(params map[string]string, rw http.ResponseWriter) {
	request := client.udiskconn.NewDescribeUDiskRequest()
	request.ProtocolVersion = ucloud.Int(1)
	request.IsBoot = ucloud.String("False")
	request.DiskType = ucloud.String("CLOUD_RSSD")
	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}
	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		request.Limit = ucloud.Int(limit)
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		request.Offset = ucloud.Int(offset)
	}

	response, err := client.udiskconn.DescribeUDisk(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.DataSet {
			if tag, ok := params["Tag"]; ok {
				if instance.Tag != tag {
					continue
				}
			}
			ids = append(ids, instance.UDiskId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}

func (client *uCloudClient) proxyDescribeUDiskSys(params map[string]string, rw http.ResponseWriter) {
	request := client.udiskconn.NewDescribeUDiskRequest()
	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}
	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		request.Limit = ucloud.Int(limit)
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		request.Offset = ucloud.Int(offset)
	}

	response, err := client.udiskconn.DescribeUDisk(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.DataSet {
			if tag, ok := params["Tag"]; ok {
				if instance.Tag != tag {
					continue
				}
			}
			if instance.IsBoot != "True" {
				continue
			}

			ids = append(ids, instance.UDiskId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}

func (client *uCloudClient) proxyDescribeURedis(params map[string]string, rw http.ResponseWriter) {
	request := client.umemconn.NewDescribeURedisGroupRequest()
	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}
	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		request.Limit = ucloud.Int(limit)
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		request.Offset = ucloud.Int(offset)
	}

	response, err := client.umemconn.DescribeURedisGroup(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.DataSet {
			if tag, ok := params["Tag"]; ok {
				if instance.Tag != tag {
					continue
				}
			}
			ids = append(ids, instance.GroupId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}

func (client *uCloudClient) proxyDescribeNatGW(params map[string]string, rw http.ResponseWriter) {
	request := client.vpcconn.NewDescribeNATGWRequest()
	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}
	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		request.Limit = ucloud.Int(limit)
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		request.Offset = ucloud.Int(offset)
	}

	response, err := client.vpcconn.DescribeNATGW(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.DataSet {
			if tag, ok := params["Tag"]; ok {
				if instance.Tag != tag {
					continue
				}
			}
			ids = append(ids, instance.NATGWId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}

func (client *uCloudClient) proxyDescribeUFile(params map[string]string, rw http.ResponseWriter) {
	request := client.ufileconn.NewDescribeBucketRequest()
	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}
	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		request.Limit = ucloud.Int(limit)
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		request.Offset = ucloud.Int(offset)
	}

	response, err := client.ufileconn.DescribeBucket(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.DataSet {
			if tag, ok := params["Tag"]; ok {
				if instance.Tag != tag {
					continue
				}
			}
			ids = append(ids, instance.BucketId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}

func (client *uCloudClient) proxyDescribeUMemCache(params map[string]string, rw http.ResponseWriter) {
	request := client.umemconn.NewDescribeUMemRequest()
	request.Protocol = ucloud.String("memcache")
	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}
	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		request.Limit = ucloud.Int(limit)
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		request.Offset = ucloud.Int(offset)
	}

	response, err := client.umemconn.DescribeUMem(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.DataSet {
			if tag, ok := params["Tag"]; ok {
				if instance.Tag != tag {
					continue
				}
			}
			ids = append(ids, instance.ResourceId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}

func (client *uCloudClient) proxyDescribeShareBW(params map[string]string, rw http.ResponseWriter) {
	req := client.ucloudconn.NewGenericRequest()
	reqMap := map[string]interface{}{
		"Action": "DescribeShareBandwidth",
	}

	if v, ok := params["Region"]; ok {
		reqMap["Region"] = v
	}
	if v, ok := params["ProjectId"]; ok {
		reqMap["ProjectId"] = v
	}

	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		reqMap["Limit"] = limit
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		reqMap["Offset"] = offset
	}

	err := req.SetPayload(reqMap)
	if err != nil {
		handleResponse(rw, nil, fmt.Errorf("set DescribeShareBandwidth requset got err, %s", err))
		return
	}

	genericResp, err := client.ucloudconn.GenericInvoke(req)
	if err != nil {
		handleResponse(rw, nil, fmt.Errorf("do DescribeShareBandwidth got err, %s", err))
		return
	}

	type DescribeShareBandwidthResponse struct {
		DataSet []struct {
			ShareBandwidthId string
		}
	}
	respDescribe := &DescribeShareBandwidthResponse{}
	if err = genericResp.Unmarshal(respDescribe); err != nil {
		handleResponse(rw, nil, fmt.Errorf("unmarshal DescribeShareBandwidth resp got err, %s", err))
		return
	}

	var ids []string
	for _, instance := range respDescribe.DataSet {
		// todo
		//if tag, ok := params["Tag"]; ok {
		//	if instance.Tag != tag {
		//		continue
		//	}
		//}
		ids = append(ids, instance.ShareBandwidthId)
	}

	d, err := json.Marshal(ids)
	log.DefaultLogger.Debug(string(d))
	handleResponse(rw, d, err)
}

func (client *uCloudClient) proxyDescribePHost(params map[string]string, rw http.ResponseWriter) {
	request := client.uphostconn.NewDescribePHostRequest()
	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}
	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		request.Limit = ucloud.Int(limit)
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		request.Offset = ucloud.Int(offset)
	}

	response, err := client.uphostconn.DescribePHost(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.PHostSet {
			if tag, ok := params["Tag"]; ok {
				if instance.Tag != tag {
					continue
				}
			}
			ids = append(ids, instance.PHostId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
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
			var isRepeat bool
			for _, id := range ids {
				if instance.Region == id {
					isRepeat = true
					break
				}
			}
			if !isRepeat {
				ids = append(ids, instance.Region)
			}
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
	if v, ok := params["Tag"]; ok {
		request.Tag = ucloud.String(v)
	}
	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		request.Limit = ucloud.Int(limit)
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		request.Offset = ucloud.Int(offset)
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

	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		request.Limit = ucloud.Int(limit)
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		request.Offset = ucloud.Int(offset)
	}

	response, err := client.unetconn.DescribeEIP(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.EIPSet {
			if tag, ok := params["Tag"]; ok {
				if instance.Tag != tag {
					continue
				}
			}
			ids = append(ids, instance.EIPId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}
func (client *uCloudClient) proxyDescribeULB(params map[string]string, rw http.ResponseWriter) {
	request := client.ulbconn.NewDescribeULBSimpleRequest()

	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}
	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		request.Limit = ucloud.Int(limit)
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		request.Offset = ucloud.Int(offset)
	}

	response, err := client.ulbconn.DescribeULBSimple(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.DataSet {
			if tag, ok := params["Tag"]; ok {
				if instance.Tag != tag {
					continue
				}
			}
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
	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		request.Limit = ucloud.Int(limit)
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		request.Offset = ucloud.Int(offset)
	}

	response, err := client.udbconn.DescribeUDBInstance(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.DataSet {
			if tag, ok := params["Tag"]; ok {
				if instance.Tag != tag {
					continue
				}
			}
			ids = append(ids, instance.DBId)
		}

		d, err := json.Marshal(ids)
		log.DefaultLogger.Debug(string(d))
		handleResponse(rw, d, err)
	}
}

func (client *uCloudClient) proxyDescribeUDPN(params map[string]string, rw http.ResponseWriter) {
	request := client.udpnconn.NewDescribeUDPNRequest()

	if v, ok := params["ProjectId"]; ok {
		request.ProjectId = ucloud.String(v)
	}
	if v, ok := params["Region"]; ok {
		request.Region = ucloud.String(v)
	}
	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		request.Limit = ucloud.Int(limit)
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		request.Offset = ucloud.Int(offset)
	}

	response, err := client.udpnconn.DescribeUDPN(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.DataSet {
			ids = append(ids, instance.UDPNId)
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
	if v, ok := params["Limit"]; ok {
		limit, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Limit must set to int value"))
			return
		}
		request.Limit = ucloud.Int(limit)
	}
	if v, ok := params["Offset"]; ok {
		offset, err := strconv.Atoi(v)
		if err != nil {
			handleResponse(rw, nil, fmt.Errorf("type is invalid, Offset must set to int value"))
			return
		}
		request.Offset = ucloud.Int(offset)
	}

	response, err := client.umemconn.DescribeUMem(request)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		handleResponse(rw, nil, err)
	} else {
		var ids []string
		for _, instance := range response.DataSet {
			if tag, ok := params["Tag"]; ok {
				if instance.Tag != tag {
					continue
				}
			}
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
