import { DataSourceInstanceSettings, MetricFindValue } from '@grafana/data';
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
    const url = 'generic_api';
    let values: MetricFindValue[] = [{ text: 'qwe' }];
    this.getResource(url, param).then((response: string[]) => {
      console.log('response', response);
      response.forEach(function(v) {
        console.log('value', v);
        values.push({ text: v });
      });
    });
    console.log('values', values);
    console.log('values2', [{ text: 'sdfas' }]);
    return values;
    // return  [
    //   {
    //     text: "abcd",
    //   },
    // ];
  }
}
