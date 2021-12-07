import React, { PureComponent, useEffect, useState, FunctionComponent, InputHTMLAttributes } from 'react';
const { Input } = LegacyForms;
import { Collapse, LegacyForms, InlineFormLabel, Segment, SegmentAsync } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './datasource';
import { MyDataSourceOptions, MyQuery, SelectableStrings } from './types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;
export class QueryEditor extends PureComponent<Props> {
  render() {
    return (
      <>
        <MetricsQueryFieldsEditor {...this.props} />
        <QueryResourceIdCollapse {...this.props} />
      </>
    );
  }
}

interface State {
  projectIds: SelectableStrings;
  regions: SelectableStrings;
  resourceTypes: SelectableStrings;
}
const MetricsQueryFieldsEditor = (props: any) => {
  const { onChange, query, onRunQuery, datasource } = props;
  const onQueryChange = (query: MyQuery) => {
    onChange(query);
    onRunQuery();
  };

  const [state, setState] = useState<State>({
    projectIds: [],
    regions: [],
    resourceTypes: [],
  });

  useEffect(() => {
    let projectParam = {
      Action: 'GetProjectId',
    };
    let regionParam = {
      Action: 'GetRegion',
    };
    let resourceTypeParam = {
      Action: 'GetResourceType',
    };
    Promise.all([
      datasource.metricFindQuery(JSON.stringify(projectParam)),
      datasource.metricFindQuery(JSON.stringify(regionParam)),
      datasource.metricFindQuery(JSON.stringify(resourceTypeParam)),
    ]).then(([projectIds, regions, resourceTypes]) => {
      setState((prevState) => ({
        ...prevState,
        projectIds: projectIds,
        regions: regions,
        resourceTypes: resourceTypes,
      }));
    });
  }, [datasource]);

  const loadMetricNames = async () => {
    return datasource
      .metricFindQuery(
        JSON.stringify({
          Action: 'GetMetricName',
          ResourceType: query.resourceType,
        })
      )
      .then((value: SelectableValue[]) => value);
  };

  const loadResourceIds = async () => {
    return datasource
      .metricFindQuery(
        JSON.stringify({
          Action: 'GetResourceId',
          ProjectId: query.projectId,
          Region: query.region,
          ResourceType: query.resourceType,
          Tag: query.tag,
          Limit: query.limit,
          Offset: query.offset,
          ULBId: query.ulbId,
          ClassType: query.classType,
        })
      )
      .then((value: SelectableValue[]) => value);
  };

  const { projectIds, regions, resourceTypes } = state;
  console.log('projectIds:///', projectIds);
  return (
    <>
      <QueryInlineField label="ProjectId">
        <Segment
          value={query.projectId}
          placeholder="Select project"
          options={projectIds}
          allowCustomValue
          onChange={({ value: projectId }) => onQueryChange({ ...query, projectId: projectId! })}
        />
      </QueryInlineField>
      <QueryInlineField label="Region">
        <Segment
          value={query.region}
          placeholder="Select region"
          options={regions}
          allowCustomValue
          onChange={({ value: region }) => onQueryChange({ ...query, region: region! })}
        />
      </QueryInlineField>
      <QueryInlineField label="ResourceType">
        <Segment
          value={query.resourceType}
          placeholder="Select resourceType"
          options={resourceTypes}
          allowCustomValue
          onChange={({ value: resourceType }) => onQueryChange({ ...query, resourceType: resourceType! })}
        />
      </QueryInlineField>
      <QueryInlineField label="MetricName">
        <SegmentAsync
          value={query.metricName}
          placeholder="Select metric name"
          loadOptions={loadMetricNames}
          allowCustomValue
          onChange={({ value: metricName }) => onQueryChange({ ...query, metricName: metricName! })}
        />
      </QueryInlineField>
      <QueryInlineField label="ResourceId">
        <SegmentAsync
          value={query.resourceId}
          placeholder="Select resourceId"
          loadOptions={loadResourceIds}
          allowCustomValue
          onChange={({ value: resourceId }) => onQueryChange({ ...query, resourceId: resourceId! })}
        />
      </QueryInlineField>
    </>
  );
};

const QueryResourceIdCollapse = (props: any) => {
  const [isOpen, setIsOpen] = useState(false);
  const { onChange, query, onRunQuery } = props;
  const { tag, limit, offset, ulbId, classType } = query;
  const onQueryChange = (query: MyQuery) => {
    onChange(query);
    onRunQuery();
  };

  const resourceType = query.resourceType;

  return (
    <div className="gf-form gf-form--grow">
      <Collapse label="ResourceId query condition" isOpen={isOpen} onToggle={() => setIsOpen(!isOpen)}>
        <div className="gf-form-inline">
          <div>
            {resourceType === 'ulb-vserver' ? (
              <div className="gf-form">
                <QueryField label="ULBId">
                  <Input
                    className="gf-form-input width-6"
                    onBlur={onRunQuery}
                    value={ulbId}
                    onChange={(v) => onQueryChange({ ...query, ulbId: v.target.value! })}
                  />
                </QueryField>
              </div>
            ) : null}
          </div>
          <div>
            {resourceType === 'udb' ? (
              <div className="gf-form">
                <QueryField label="ClassType">
                  <Input
                    className="gf-form-input width-6"
                    onBlur={onRunQuery}
                    value={classType}
                    onChange={(v) => onQueryChange({ ...query, classType: v.target.value! })}
                  />
                </QueryField>
              </div>
            ) : null}
          </div>
          <div className="gf-form">
            <QueryField label="Offset">
              <Input
                className="gf-form-input width-6"
                onBlur={onRunQuery}
                value={offset}
                type="number"
                onChange={(v) => onQueryChange({ ...query, offset: v.target.value! })}
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
                onChange={(v) => onQueryChange({ ...query, limit: v.target.value! })}
              />
            </QueryField>
          </div>
          <div className="gf-form gf-form--grow">
            <QueryField className="gf-form--grow" label="Tag">
              <Input
                className="gf-form-input"
                onBlur={onRunQuery}
                value={tag}
                onChange={(v) => onQueryChange({ ...query, tag: v.target.value! })}
              />
            </QueryField>
          </div>
        </div>
      </Collapse>
    </div>
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
