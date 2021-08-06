import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { MyDataSourceOptions, MyQuery, MyVariableQuery } from './types';

export class DataSource extends DataSourceWithBackend<MyQuery, MyDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
  }

  async metricFindQuery(query: MyVariableQuery, options?: any) {
    console.log('query1', query);
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
