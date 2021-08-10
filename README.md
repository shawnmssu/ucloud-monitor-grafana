# UCloud 优刻得云监控插件

[![Build](https://github.com/grafana/grafana-starter-datasource-backend/workflows/CI/badge.svg)](https://github.com/grafana/grafana-datasource-backend/actions?query=workflow%3A%22CI%22)


## 简介
 [UCloud 优刻得云监控](https://docs.ucloud.cn/umon/README) 能够提供对UCloud云平台中产品及资源的监控，通过告警通知管理及监控模板的设置，使您能够实时掌握资源及应用的状态，保证您的服务及应用的良性运行。
 目前云监控已经通过 datasource 插件的形式集成到了 grafana 中，通过简单的配置，可以构建 UCloud 产品监控大盘。目前已支持了云主机（UHost）、数据库（UDB）、负载均衡（ULB）等主流产品。

## 快速开始

### 1、直接安装云监控 grafana 数据源
    a. 直接 从release 页面 https://github.com/shawnmssu/ucloud-umon-grafana-datasource/releases 里面下载 ucloud-umon-datasource_v0.1.tar.gz
    b. 下载到 grafan的plugin目录中，解压缩 tar -xzf ucloud-umon-datasource_v0.1.tar.gz
    c. 修改 conf/defaults.ini 允许未签名插件运行
        allow_loading_unsigned_plugins = ucloud-umon-datasource
    d. 重启grafana
### 2、源代码安装
    a. 代码编译
        - git clone https://github.com/shawnmssu/ucloud-monitor-grafana.git
        - 进入 ucloud-umon-grafana-datasource 目录下, 执行 make build 命令(依赖 make golang mage yarn)。

    b. 部署
        1）按照上面顺序编译完成后，代码都会到 dist下面。包括前端文件和二进制可执行文件 cms-datasource*。
        2）保证 ucloud-umon-datasource* 都具有可执行权限。chmod +x ucloud-umon-datasource*
        3) 在 grafana 的 plugin目录中，创建 ucloud-umon-datasource 目录，把编译出来的dist目录拷贝到此
        4) 修改 conf/defaults.ini 允许未签名插件运行
            allow_loading_unsigned_plugins = aliyun_cms_grafana_datasource
        5) 重启grafana

### 3、配置云监控 grafana 数据源
    a.进入 grafana 的数据源配置页面(Data Sources), 点击 Add data source 进入配置表单页面,填入数据源名称 ucloud-umon-datasource 并选择；
    b. 填写公私钥和配置信息
        其中 Public Key 和 Private Key 为必填，可以从[控制台](https://console.ucloud.cn/uapi/apikey)获取;
        如果显示 Data source is working,说明数据源配置成功,可以开始在 grafana 中访问 UCloud 云监控的数据了。
    
## 配置 Dashboard 图表
### Data Source Query 参数
   |  参数   | 说明  | 备注|
   |  :----:  | :----:  | :----:|
   | ProjectId  | 项目ID | - |
   | Region | 资源所在地域 | - |
   | ResourceType  | 资源类型 | 已支持 uhost, eip, ulb, udb, umem |
   | MetricName  | 监控指标 | - |
   | ResourceId  | 资源ID | - |
   |  - | - | - |
   | Tag  | 查询资源的业务组名称 | Query ResourceId 相关参数|
   | Limit  | 返回数据长度，默认为20，最大100 | Query ResourceId 相关参数 |
   | Offset  | 列表起始位置偏移量，默认为0 | Query ResourceId 相关参数 |

### 配置 variables
- Variables支持 Type 类型为 Query 和 Custom，具体请参考 [grafana 官方文档](https://grafana.com/docs/grafana/latest/variables/variable-types/),
  其中配置 Type 为 Query， 支持通过自定义 json 数据来获取 variable。
  
  |  参数   | 说明  | 备注|
  |  :----:  | :----:  | :----:|
  | Action  | 获取 variable 的 API 名称 | 规则 Get + Query 参数， 例如 GetMetricName |
  | Data Source Query 相关参数  | - | - |

-  例如：
   - 查询监控指标：{ "Action": "GetMetricName","Region": "cn-bj2", "ResourceType": "uhost" }
   - 查询资源ID：{ "Action": "GetResourceId","Region": "cn-bj2", "Tag": "Default" }
