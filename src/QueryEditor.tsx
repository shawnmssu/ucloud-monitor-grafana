import _, { defaults } from 'lodash';

import React, { ChangeEvent, PureComponent, useEffect, useState } from 'react'; //useEffect
import { Collapse, LegacyForms } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, MyDataSourceOptions, MyQuery } from './types';
import { getTemplateSrv } from '@grafana/runtime';

const { FormField, Select } = LegacyForms; //FormField

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  render() {
    return (
      <div className="gf-form-inline">
        <div className="gf-form">
          <ProjectIdSelect {...this.props} />
          <RegionSelect {...this.props} />
          <ResourceTypeSelect {...this.props} />
          <MetricNameSelect {...this.props} />
          <ResourceIdSelect {...this.props} />
          <QueryResourceIdCollapse {...this.props} />
        </div>
      </div>
    );
  }
}

const QueryResourceIdCollapse = (props: any) => {
  const [isOpen, setIsOpen] = useState(false);
  const query = defaults(props.query, defaultQuery);
  const { tag, limit, offset } = query;

  const onTagChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = props;
    query.tag = event.target.value || '';
    onChange({ ...query });
    // executes the query
    onRunQuery();
  };
  const onLimitChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = props;
    query.limit = event.target.value || '';
    onChange({ ...query });
    // executes the query
    onRunQuery();
  };
  const onOffsetChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = props;
    query.offset = event.target.value || '';
    onChange({ ...query });
    // executes the query
    onRunQuery();
  };

  return (
    <div className="gf-form-inline">
      <Collapse label="ResourceId query condition" isOpen={isOpen} onToggle={() => setIsOpen(!isOpen)}>
        <div className="gf-form">
          <FormField width={3} value={tag} onChange={onTagChange} label="Tag" type="string" step="0.1" inputWidth={5} />
          <FormField
            width={3}
            value={limit}
            onChange={onLimitChange}
            label="Limit"
            type="number"
            step="0.1"
            inputWidth={5}
          />
          <FormField
            width={3}
            value={offset}
            onChange={onOffsetChange}
            label="Offset"
            type="number"
            step="0.1"
            inputWidth={5}
          />
        </div>
      </Collapse>
    </div>
  );
};

const ProjectIdSelect = (props: any) => {
  const [value, setValue] = useState<SelectableValue<string>>({
    label: props.query.projectId,
    value: props.query.projectId,
  });

  const onProjectIdChange = (value: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = props;
    // query.projectId = setTemplateVariable(`${value.value}`) || '';
    query.projectId = value.value || '';
    onChange({ ...query });
    onRunQuery();
  };

  const projectIdOptions: Array<SelectableValue<string>> = [];
  useEffect(() => {
    const { datasource } = props;
    let param = {
      Action: 'GetProjectId',
    };
    datasource.getResource('generic_api', param).then((response: string[]) => {
      Array.prototype.forEach.call(response || [], (v) => {
        projectIdOptions.push({ label: v, value: v });
      });
    });
  });

  return (
    <div className="gf-form-inline">
      <label className="gf-form-label width-6">ProjectId</label>
      <Select
        isSearchable
        isClearable
        className="gf-form-input"
        options={projectIdOptions}
        value={value}
        allowCustomValue
        onChange={(v) => {
          setValue(v);
          onProjectIdChange(v);
        }}
      />
    </div>
  );
};

const RegionSelect = (props: any) => {
  const [value, setValue] = useState<SelectableValue<string>>({ label: props.query.region, value: props.query.region });

  const onRegionChange = (value: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = props;
    // query.region = setTemplateVariable(`${value.value}`) || '';
    query.region = value.value || '';
    onChange({ ...query });
    onRunQuery();
  };

  const regionOptions: Array<SelectableValue<string>> = [];
  useEffect(() => {
    const { datasource } = props;
    let param = {
      Action: 'GetRegion',
    };
    datasource.getResource('generic_api', param).then((response: string[]) => {
      Array.prototype.forEach.call(response || [], (v) => {
        regionOptions.push({ label: v, value: v });
      });
    });
  });

  return (
    <div className="gf-form-inline">
      <label className="gf-form-label width-6">Region</label>
      <Select
        isSearchable
        isClearable
        className="gf-form-input"
        options={regionOptions}
        value={value}
        allowCustomValue
        onChange={(v) => {
          setValue(v);
          onRegionChange(v);
        }}
      />
    </div>
  );
};

const MetricNameSelect = (props: any) => {
  const [value, setValue] = useState<SelectableValue<string>>({
    label: props.query.metricName,
    value: props.query.metricName,
  });

  const onMetricNameChange = (value: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = props;
    // query.metricName = setTemplateVariable(`${value.value}`) || '';
    query.metricName = value.value || '';
    onChange({ ...query });
    onRunQuery();
  };

  const metricNameOptions: Array<SelectableValue<string>> = [];
  useEffect(() => {
    const { query, datasource } = props;
    let param = {
      Action: 'GetMetricName',
      ResourceType: setTemplateVariable(query.resourceType),
    };
    datasource.getResource('generic_api', param).then((response: string[]) => {
      Array.prototype.forEach.call(response || [], (v) => {
        metricNameOptions.push({ label: v, value: v });
      });
    });
  });

  return (
    <div className="gf-form-inline">
      <label className="gf-form-label width-6">MetricName</label>
      <Select
        isSearchable
        isClearable
        className="gf-form-input"
        options={metricNameOptions}
        value={value}
        allowCustomValue
        onChange={(v) => {
          setValue(v);
          onMetricNameChange(v);
        }}
      />
    </div>
  );
};

const ResourceTypeSelect = (props: any) => {
  const [value, setValue] = useState<SelectableValue<string>>({
    label: props.query.resourceType,
    value: props.query.resourceType,
  });

  const onResourceTypeChange = (value: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = props;
    // query.resourceType = setTemplateVariable(`${value.value}`) || '';
    query.resourceType = value.value || '';
    onChange({ ...query });
    onRunQuery();
  };

  const resourceTypeOptions: Array<SelectableValue<string>> = [];
  useEffect(() => {
    const { datasource } = props;
    let param = {
      Action: 'GetResourceType',
    };
    datasource.getResource('generic_api', param).then((response: string[]) => {
      Array.prototype.forEach.call(response || [], (v) => {
        resourceTypeOptions.push({ label: v, value: v });
      });
    });
  });

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
        onChange={(v) => {
          setValue(v);
          onResourceTypeChange(v);
        }}
      />
    </div>
  );
};

const ResourceIdSelect = (props: any) => {
  const [value, setValue] = useState<SelectableValue<string>>({
    label: props.query.resourceId,
    value: props.query.resourceId,
  });

  const onResourceIdChange = (value: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = props;
    // query.resourceId = setTemplateVariable(`${value.value}`) || '';
    query.resourceId = value.value || '';
    onChange({ ...query });
    onRunQuery();
  };
  const resourceIdOptions: Array<SelectableValue<string>> = [];
  useEffect(() => {
    const { query, datasource } = props;
    let param = {
      Action: 'GetResourceId',
      ProjectId: setTemplateVariable(query.projectId),
      Region: setTemplateVariable(query.region),
      ResourceType: setTemplateVariable(query.resourceType),
      Tag: setTemplateVariable(query.tag),
      Limit: setTemplateVariable(query.limit),
      Offset: setTemplateVariable(query.offset),
    };

    datasource.getResource('generic_api', param).then((response: string[]) => {
      Array.prototype.forEach.call(response || [], (v) => {
        resourceIdOptions.push({ label: v, value: v });
      });
    });
  });
  return (
    <div className="gf-form-inline">
      <label className="gf-form-label width-6">ResourceId</label>
      <Select
        isSearchable
        isClearable
        className="gf-form-input"
        options={resourceIdOptions}
        value={value}
        allowCustomValue
        onChange={(v) => {
          setValue(v);
          onResourceIdChange(v);
        }}
      />
    </div>
  );
};

export const setTemplateVariable = (value: string) => {
  if (value?.includes('$')) {
    Array.prototype.forEach.call(getTemplateSrv().getVariables() || [], (v) => {
      value = value.replace(`$${_.get(v, 'name')}`, _.get(v, ['current', 'value']));
    });
  }

  return value;
};
