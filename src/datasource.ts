import { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { MyDataSourceOptions, MyQuery, MyVariableQuery } from './types';
import { setTemplateVariable } from './QueryEditor';

export class DataSource extends DataSourceWithBackend<MyQuery, MyDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
  }

  applyTemplateVariables(query: MyQuery, scopedVars: ScopedVars): Record<string, any> {
    query.projectId = setTemplateVariable(query.projectId || '');
    query.region = setTemplateVariable(query.region);
    query.resourceType = setTemplateVariable(query.resourceType);
    query.metricName = setTemplateVariable(query.metricName);
    query.resourceId = setTemplateVariable(query.resourceId);
    query.tag = setTemplateVariable(query.tag);
    return super.applyTemplateVariables(query, scopedVars);
  }

  async metricFindQuery(query: MyVariableQuery, options?: any) {
    console.log('options', options);
    const obj = JSON.parse(query.query);
    let param = {
      Action: obj.Action,
      ProjectId: obj.ProjectId,
      Region: obj.Region,
      ResourceType: obj.ResourceType,
    };

    const response = await this.getResource('generic_api', param);
    return response.map((v: string) => ({ text: v }));
  }
}
