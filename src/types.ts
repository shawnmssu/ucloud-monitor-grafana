import { DataQuery, DataSourceJsonData, SelectableValue } from '@grafana/data';

export interface MyQuery extends DataQuery {
  projectId?: string;
  region: string;
  resourceType: string;
  metricName: string;
  resourceId: string;
  tag: string;
  limit: number;
  offset: number;
  ulbId: string;
  classType: string;
}

export type SelectableStrings = Array<SelectableValue<string>>;

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  projectId?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  publicKey: string;
  privateKey: string;
}
