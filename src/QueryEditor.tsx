import { defaults } from 'lodash';

import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, MyDataSourceOptions, MyQuery } from './types';

const { FormField } = LegacyForms;

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  onProjectIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, projectId: event.target.value });
    // executes the query
    onRunQuery();
  };

  onRegionChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, region: event.target.value });
    // executes the query
    onRunQuery();
  };

  onResourceTypeChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, resourceType: event.target.value });
    // executes the query
    onRunQuery();
  };
  onMetricNameChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, metricName: event.target.value });
    // executes the query
    onRunQuery();
  };
  onResourceIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, resourceId: event.target.value });
    // executes the query a
    onRunQuery();
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { projectId, region, resourceType, metricName, resourceId } = query;

    console.log('query', query);

    return (
      <div className="gf-form">
        {/*<AsyncSelect*/}
        {/*    loadOptions={() => {*/}
        {/*        console.log('asd');*/}
        {/*        return new Promise<Array<SelectableValue<string>>>((resolve) => {*/}
        {/*            setTimeout(() => {*/}
        {/*                resolve([*/}
        {/*                    { label: 'Basic option', value: "0" },*/}
        {/*                    { label: 'Option with description', value: "1", description: 'this is a description' },*/}
        {/*                    {*/}
        {/*                        label: 'Option with description and image',*/}
        {/*                        value: "1",*/}
        {/*                        description: 'This is a very elaborate description, describing all the wonders in the world.',*/}
        {/*                        imgUrl: 'https://placekitten.com/40/40',*/}
        {/*                    },*/}
        {/*                ]);*/}
        {/*            }, 2000);*/}
        {/*        });*/}
        {/*    }*/}
        {/*    }*/}
        {/*    defaultOptions*/}
        {/*    value={{label:"fasdf", value: "afsdf"}}*/}
        {/*    onChange={this.onProjectIdChange}*/}
        {/*/>*/}
        {/*===================================================*/}
        <FormField
          width={3}
          value={projectId}
          onChange={this.onProjectIdChange}
          label="ProjectId"
          type="string"
          step="0.1"
          required={true}
        />
        <FormField width={3} value={region} onChange={this.onRegionChange} label="Region" type="string" step="0.1" />
        <FormField
          width={3}
          value={resourceType}
          onChange={this.onResourceTypeChange}
          label="ResourceType"
          type="string"
          step="0.1"
        />
        <FormField
          width={3}
          value={metricName}
          onChange={this.onMetricNameChange}
          label="MetricName"
          type="string"
          step="0.1"
        />
        <FormField
          width={3}
          value={resourceId}
          onChange={this.onResourceIdChange}
          label="ResourceId"
          type="string"
          step="0.1"
        />
      </div>
    );
  }
}
