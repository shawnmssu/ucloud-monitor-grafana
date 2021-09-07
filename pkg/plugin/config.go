package plugin

import (
	"encoding/json"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/ucloud/ucloud-sdk-go/services/uaccount"
	"github.com/ucloud/ucloud-sdk-go/services/udb"
	"github.com/ucloud/ucloud-sdk-go/services/udpn"
	"github.com/ucloud/ucloud-sdk-go/services/uhost"
	"github.com/ucloud/ucloud-sdk-go/services/ulb"
	"github.com/ucloud/ucloud-sdk-go/services/umem"
	"github.com/ucloud/ucloud-sdk-go/services/unet"
	"github.com/ucloud/ucloud-sdk-go/services/vpc"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
	"github.com/ucloud/ucloud-sdk-go/ucloud/log"
	"time"
)

type uCloudClient struct {
	ucloudconn   *ucloud.Client
	uhostconn    *uhost.UHostClient
	unetconn     *unet.UNetClient
	ulbconn      *ulb.ULBClient
	vpcconn      *vpc.VPCClient
	udbconn      *udb.UDBClient
	umemconn     *umem.UMemClient
	udpnconn     *udpn.UDPNClient
	uaccountconn *uaccount.UAccountClient
}

type config struct {
	ProjectId  string
	PublicKey  string
	PrivateKey string
}

func getUCloudConfig(instanceSettings backend.DataSourceInstanceSettings) (*config, error) {
	var setting config
	jsonData := map[string]interface{}{}
	if err := json.Unmarshal(instanceSettings.JSONData, &jsonData); err != nil {
		return nil, err
	}

	if v, ok := jsonData["projectId"]; ok {
		setting.ProjectId = v.(string)
	}
	setting.PublicKey = instanceSettings.DecryptedSecureJSONData["publicKey"]
	setting.PrivateKey = instanceSettings.DecryptedSecureJSONData["privateKey"]

	return &setting, nil
}

func (c *config) Client() *uCloudClient {
	var client uCloudClient

	cfg := ucloud.NewConfig()
	cfg.ProjectId = c.ProjectId

	cfg.LogLevel = log.PanicLevel
	cfg.UserAgent = "UCloud-monitor-grafana"

	cred := auth.NewCredential()
	cred.PublicKey = c.PublicKey
	cred.PrivateKey = c.PrivateKey

	// initialize client connections
	client.ucloudconn = ucloud.NewClient(&cfg, &cred)
	client.unetconn = unet.NewClient(&cfg, &cred)
	client.ulbconn = ulb.NewClient(&cfg, &cred)
	client.vpcconn = vpc.NewClient(&cfg, &cred)
	client.umemconn = umem.NewClient(&cfg, &cred)
	client.udpnconn = udpn.NewClient(&cfg, &cred)
	client.uaccountconn = uaccount.NewClient(&cfg, &cred)
	longtimeCfg := cfg
	longtimeCfg.Timeout = 60 * time.Second
	client.udbconn = udb.NewClient(&longtimeCfg, &cred)
	client.uhostconn = uhost.NewClient(&longtimeCfg, &cred)
	return &client
}
