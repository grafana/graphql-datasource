import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import ConfigEditor from './ConfigEditor';
import QueryEditor from './QueryEditor';
import { GrapQLQuery, MyDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, GrapQLQuery, MyDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
