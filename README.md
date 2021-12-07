# UCloud 优刻得云监控插件

[![Build](https://github.com/grafana/grafana-starter-datasource-backend/workflows/CI/badge.svg)](https://github.com/grafana/grafana-datasource-backend/actions?query=workflow%3A%22CI%22)


## 简介
 [UCloud 优刻得云监控](https://docs.ucloud.cn/umon/README) 能够提供对UCloud云平台中产品及资源的监控，通过告警通知管理及监控模板的设置，使您能够实时掌握资源及应用的状态，保证您的服务及应用的良性运行。
 目前云监控已经通过 datasource 插件的形式集成到了 grafana 中，通过简单的配置，可以构建 UCloud 产品监控大盘。目前已支持了云主机（UHost）、数据库（UDB）、负载均衡（ULB）等主流产品。

## 快速开始

### 安装数据源 

#### release 安装

- 从 release 页面 https://github.com/shawnmssu/ucloud-monitor-grafana/releases 下载并解压到 grafana 的 plugin 目录中
- 修改 [configration](https://grafana.com/docs/grafana/latest/administration/configuration/) 中的 plugins 配置，允许未签名插件运行：
   allow_loading_unsigned_plugins = ucloud-monitor-datasource
- 重启 grafana

#### 源代码安装

- 代码编译
    - git clone https://github.com/shawnmssu/ucloud-monitor-grafana.git
    - 进入 ucloud-monitor-grafana 目录下, 执行 make build 命令(依赖 make golang mage yarn)。

- 部署
   - 将 dist 目录下的文件 ucloud-monitor-datasource-backend* 增加可执行权限：chmod +x ucloud-monitor-datasource-backend*
   - 在 grafana 的 plugin目录中，创建 ucloud-monitor-datasource 目录，把编译出来的 dist 目录拷贝到此
   - 修改 [configration](https://grafana.com/docs/grafana/latest/administration/configuration/) 中的 plugins 配置，允许未签名插件运行：
     allow_loading_unsigned_plugins = ucloud-monitor-datasource 
   - 重启 grafana

### 配置云监控 grafana 数据源

  - 进入 grafana 的数据源配置页面(Data Sources), 点击 Add data source 进入配置表单页面,填入数据源名称 UCloud Monitor 并选择； 
  - 填写公私钥和配置信息:
    其中 Public Key 和 Private Key 为必填，可以从 [控制台](https://console.ucloud.cn/uapi/apikey) 获取;
    如果显示 Data source is working，说明数据源配置成功，可以开始在 grafana 中访问 UCloud 云监控的数据了。
    
## 配置 Dashboard 图表

### Data Source Query 参数

   |  参数   | 说明  | 备注| 必填
   |  :----:  | :----:  | :----:|:----:|
   | ProjectId  | 项目ID | - | 是 |
   | Region | 资源所在地域 | - | 是 |
   | ResourceType  | 资源类型 | 已支持 uhost, eip, ulb, ulb-vserver, udb, umem, udpn, phost, sharebandwidth, umemcache, uredis, natgw, ufile, udisk. udisk_ssd, udisk_rssd, udisk_sys | 是 |
   | MetricName  | 监控指标 | 不同 ResourceType 支持不同的监控指标，参考 [DescribeResourceMetric](https://docs.ucloud.cn/api/umon-api/describe_resource_metric)| 是 |
   | ResourceId  | 资源ID | - | 是 |
   |  - | - | - |
   | Tag  | 查询资源的业务组名称 | Query ResourceId 相关参数 | 否 |
   | Limit  | 返回数据长度，默认为20，最大100 | Query ResourceId 相关参数 | 否 |
   | Offset  | 列表起始位置偏移量，默认为0 | Query ResourceId 相关参数 | 否 |
   | ULBId   | ULB 的资源 ID | Query ulb-vserver ResourceId 相关参数 | 否 |
   | ClassType   | UDB 的资源的类型 | Query udb ResourceId 相关参数，已支持 mysql: sql；mongo: nosql；postgresql: postgresql，参考 [DescribeUDBInstance](https://docs.ucloud.cn/api/udb-api/describe_udb_instance)| 否 |

### 配置 variables

- Variables支持 Type 类型为 Query 和 Custom，具体请参考 [grafana 官方文档](https://grafana.com/docs/grafana/latest/variables/variable-types/),
  其中配置 Type 为 Query， 支持通过自定义 json 数据来获取 variable。
  
  |  参数   | 说明  | 备注|
  |  :----:  | :----:  | :----:|
  | Action  | 获取 variable 的 API 名称 | 规则 Get + Query 参数， 例如 GetMetricName |
  | Data Source Query 相关参数  | - | - |

-  例如：
   - 查询监控指标：{ "Action": "GetMetricName","Region": "cn-bj2", "ResourceType": "uhost" }
   - 查询资源ID：{ "Action": "GetResourceId","ResourceType": "uhost", Region": "cn-bj2", "Tag": "Default" }

### 预设 Dashboard

- 可以在配置数据源时 import 预设的 Dashboard，目前已支持 UCLoud UHost

### 参考文档
- [API 文档](https://docs.ucloud.cn/api)
- [获取监控数据 - GetMetric](https://docs.ucloud.cn/api/umon-api/get_metric)
- [获取资源支持监控指标信息 - DescribeResourceMetric](https://docs.ucloud.cn/api/umon-api/describe_resource_metric)
- [获取云数据库信息 - DescribeUDBInstance](https://docs.ucloud.cn/api/udb-api/describe_udb_instance)