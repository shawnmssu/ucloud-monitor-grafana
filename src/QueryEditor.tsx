import { defaults } from 'lodash';

import React, { ChangeEvent, PureComponent, useState } from 'react'; //useEffect
import { LegacyForms } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, MyDataSourceOptions, MyQuery } from './types';

const { FormField, AsyncSelect, Select } = LegacyForms;

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
  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { projectId, region, metricName } = query;

    console.log('query', query);

    return (
      <div className="gf-form">
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
          value={metricName}
          onChange={this.onMetricNameChange}
          label="MetricName"
          type="string"
          step="0.1"
        />
        <ResourceTypeSelect {...this.props} />
        <ResourceIdSelect {...this.props} />
      </div>
    );
  }
}

const ResourceTypeSelect = (props: any) => {
  const [value, setValue] = useState<SelectableValue<string>>();

  const onResourceTypeChange = (value: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = props;
    onChange({ ...query, resourceType: value.value || '' });
    onRunQuery();
  };
  // useEffect(() => {
  //
  // })
  const resourceTypeOptions = [
    { label: 'uhost', value: 'uhost' },
    { label: 'eip', value: 'eip' },
    { label: 'ulb', value: 'ulb' },
    { label: 'udb', value: 'udb' },
    { label: 'umem', value: 'umem' },
  ];

  return (
    <div className="gf-form-inline">
      <label className="gf-form-label width-6">ResourceType</label>
      <Select
        isSearchable
        isClearable
        className="gf-form-input"
        options={resourceTypeOptions}
        value={value}
        allowCustomValue
        onChange={v => {
          setValue(v);
          onResourceTypeChange(v);
        }}
      />
    </div>
  );
};

const ResourceIdSelect = (props: any) => {
  const [value, setValue] = useState<SelectableValue<string>>();

  const onResourceIdChange = (value: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = props;
    onChange({ ...query, resourceId: value.value || '' });
    onRunQuery();
  };
  // useEffect(() => {
  //
  // })
  const loadAsyncOptions = () => {
    return new Promise<Array<SelectableValue<string>>>(resolve => {
      const { query, datasource } = props;
      let param = {
        Action: 'GetResourceId',
        ProjectId: query.projectId,
        Region: query.region,
        ResourceType: query.resourceType,
      };

      let selectValues: Array<SelectableValue<string>> = [];
      datasource.getResource('generic_api', param).then((response: string[]) => {
        response.forEach(function(v) {
          selectValues.push({ label: v, value: v });
        });
      });

      resolve(selectValues);
    });
  };

  return (
    <div className="gf-form-inline">
      <label className="gf-form-label width-6">ResourceId</label>
      <AsyncSelect
        isSearchable
        isClearable
        className="gf-form-input"
        loadOptions={loadAsyncOptions}
        value={value}
        defaultOptions
        onChange={v => {
          setValue(v);
          onResourceIdChange(v);
        }}
      />
    </div>
  );
};
