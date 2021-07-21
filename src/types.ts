import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface MyQuery extends DataQuery {
  region: string;
  resourceType: string;
  metricName: string;
  resourceId: string;
}

export const defaultQuery: Partial<MyQuery> = {
};

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
