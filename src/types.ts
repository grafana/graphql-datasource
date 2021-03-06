import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface JQField {
  name: string;
  jq: string;
}

export interface GrapQLQuery extends DataQuery {
  query: string;
  fields?: JQField[];
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
