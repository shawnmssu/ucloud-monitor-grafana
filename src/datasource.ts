import { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { MyDataSourceOptions, MyQuery } from './types';
import { MetricFindValue } from '@grafana/data/types/datasource';

export class DataSource extends DataSourceWithBackend<MyQuery, MyDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
  }

  applyTemplateVariables(query: MyQuery, scopedVars: ScopedVars): Record<string, any> {
    query.projectId = getTemplateSrv().replace(query.projectId || '');
    query.region = getTemplateSrv().replace(query.region);
    query.resourceType = getTemplateSrv().replace(query.resourceType);
    query.metricName = getTemplateSrv().replace(query.metricName);
    query.resourceId = getTemplateSrv().replace(query.resourceId);
    query.tag = getTemplateSrv().replace(query.tag);
    return super.applyTemplateVariables(query, scopedVars);
  }

  async metricFindQuery(query: string, options?: any) {
    if (query) {
      let obj: any;
      try {
        obj = JSON.parse(getTemplateSrv().replace(query));
      } catch (e) {
        console.log('[Find Query error]:', e);
        return Promise.resolve([]);
      }

      let param = {
        Action: obj.Action,
        ProjectId: obj.ProjectId,
        Region: obj.Region,
        ResourceType: obj.ResourceType,
      };

      let respArr: MetricFindValue[] = [];
      await this.getResource('generic_api', param).then((response :any) => {
        if (response instanceof Array) {
          Array.prototype.forEach.call(response || [], (v) => {
            respArr.push({ text: v, value: v });
          });
        }
      });
      return respArr;
    }
    return Promise.resolve([]);
  }
}
