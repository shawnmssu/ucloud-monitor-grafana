import _, { defaults } from 'lodash';

import React, { ChangeEvent, PureComponent, useEffect, useState, FunctionComponent, InputHTMLAttributes } from 'react';
const { Input } = LegacyForms;
import { Collapse, LegacyForms, InlineFormLabel, Segment } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, MyDataSourceOptions, MyQuery } from './types';
import { getTemplateSrv } from '@grafana/runtime';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  render() {
    return (
      <>
        <ProjectIdSelect {...this.props} />
        <RegionSelect {...this.props} />
        <ResourceTypeSelect {...this.props} />
        <MetricNameSelect {...this.props} />
        <ResourceIdSelect {...this.props} />
        <QueryResourceIdCollapse {...this.props} />
      </>
    );
  }
}

const QueryResourceIdCollapse = (props: any) => {
  const [isOpen, setIsOpen] = useState(false);
  const query = defaults(props.query, defaultQuery);
  const { tag, limit, offset, onRunQuery } = query;

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
    onRunQuery();
  };
  const onOffsetChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = props;
    query.offset = event.target.value || '';
    onChange({ ...query });
    onRunQuery();
  };

  return (
    <div className="gf-form gf-form--grow">
      <Collapse label="ResourceId query condition" isOpen={isOpen} onToggle={() => setIsOpen(!isOpen)}>
        <div className="gf-form-inline">
          <div className="gf-form">
            <QueryField label="Offset">
              <Input
                className="gf-form-input width-6"
                onBlur={onRunQuery}
                value={offset}
                type="number"
                onChange={onOffsetChange}
              />
            </QueryField>
          </div>
          <div className="gf-form">
            <QueryField label="Limit">
              <Input
                className="gf-form-input width-6"
                onBlur={onRunQuery}
                value={limit}
                type="number"
                onChange={onLimitChange}
              />
            </QueryField>
          </div>
          <div className="gf-form gf-form--grow">
            <QueryField className="gf-form--grow" label="Tag">
              <Input className="gf-form-input" onBlur={onRunQuery} value={tag} onChange={onTagChange} />
            </QueryField>
          </div>
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
    <QueryInlineField label="ProjectId">
      <Segment
        value={value}
        placeholder="Select project"
        options={projectIdOptions}
        allowCustomValue
        onChange={(v) => {
          setValue(v);
          onProjectIdChange(v);
        }}
      />
    </QueryInlineField>
  );
};

const RegionSelect = (props: any) => {
  const [value, setValue] = useState<SelectableValue<string>>({ label: props.query.region, value: props.query.region });

  const onRegionChange = (value: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = props;
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
    <QueryInlineField label="Region">
      <Segment
        value={value}
        placeholder="Select region"
        options={regionOptions}
        allowCustomValue
        onChange={(v) => {
          setValue(v);
          onRegionChange(v);
        }}
      />
    </QueryInlineField>
  );
};

const MetricNameSelect = (props: any) => {
  const [value, setValue] = useState<SelectableValue<string>>({
    label: props.query.metricName,
    value: props.query.metricName,
  });

  const onMetricNameChange = (value: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = props;
    query.metricName = value.value || '';
    onChange({ ...query });
    onRunQuery();
  };

  const metricNameOptions: Array<SelectableValue<string>> = [];
  useEffect(() => {
    const { query, datasource } = props;
    let param = {
      Action: 'GetMetricName',
      ResourceType: getTemplateSrv().replace(query.resourceType),
    };
    datasource.getResource('generic_api', param).then((response: string[]) => {
      Array.prototype.forEach.call(response || [], (v) => {
        metricNameOptions.push({ label: v, value: v });
      });
    });
  });

  return (
    <QueryInlineField label="MetricName">
      <Segment
        value={value}
        placeholder="Select metricName"
        options={metricNameOptions}
        allowCustomValue
        onChange={(v) => {
          setValue(v);
          onMetricNameChange(v);
        }}
      />
    </QueryInlineField>
  );
};

const ResourceTypeSelect = (props: any) => {
  const [value, setValue] = useState<SelectableValue<string>>({
    label: props.query.resourceType,
    value: props.query.resourceType,
  });

  const onResourceTypeChange = (value: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = props;
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
    <QueryInlineField label="ResourceType">
      <Segment
        value={value}
        placeholder="Select resourceType"
        options={resourceTypeOptions}
        allowCustomValue
        onChange={(v) => {
          setValue(v);
          onResourceTypeChange(v);
        }}
      />
    </QueryInlineField>
  );
};

const ResourceIdSelect = (props: any) => {
  const [value, setValue] = useState<SelectableValue<string>>({
    label: props.query.resourceId,
    value: props.query.resourceId,
  });

  const onResourceIdChange = (value: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = props;
    query.resourceId = value.value || '';
    onChange({ ...query });
    onRunQuery();
  };
  const resourceIdOptions: Array<SelectableValue<string>> = [];
  useEffect(() => {
    const { query, datasource } = props;
    let param = {
      Action: 'GetResourceId',
      ProjectId: getTemplateSrv().replace(query.projectId),
      Region: getTemplateSrv().replace(query.region),
      ResourceType: getTemplateSrv().replace(query.resourceType),
      Tag: getTemplateSrv().replace(query.tag),
      Limit: getTemplateSrv().replace(query.limit),
      Offset: getTemplateSrv().replace(query.offset),
    };

    datasource.getResource('generic_api', param).then((response: string[]) => {
      Array.prototype.forEach.call(response || [], (v) => {
        resourceIdOptions.push({ label: v, value: v });
      });
    });
  });
  return (
    <QueryInlineField label="ResourceId">
      <Segment
        value={value}
        placeholder="Select resourceId"
        options={resourceIdOptions}
        allowCustomValue
        onChange={(v) => {
          setValue(v);
          onResourceIdChange(v);
        }}
      />
    </QueryInlineField>
  );
};

interface PropsField extends InputHTMLAttributes<HTMLInputElement> {
  label: string;
  tooltip?: string;
  children?: React.ReactNode;
}

const QueryField: FunctionComponent<Partial<PropsField>> = ({ label, tooltip, children }) => (
  <>
    <InlineFormLabel width={8} className="query-keyword" tooltip={tooltip}>
      {label}
    </InlineFormLabel>
    {children}
  </>
);

const QueryInlineField: FunctionComponent<PropsField> = ({ ...props }) => {
  return (
    <div className={'gf-form-inline'}>
      <QueryField {...props} />
      <div className="gf-form gf-form--grow">
        <div className="gf-form-label gf-form-label--grow" />
      </div>
    </div>
  );
};
