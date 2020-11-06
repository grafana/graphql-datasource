import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface jqField {
  name: string;
  jq: string;
}

export interface MyQuery extends DataQuery {
  query: string;
  fields?: jqField[];
}

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  url?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  apiKey?: string;
}
