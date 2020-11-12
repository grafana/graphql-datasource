import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { MyDataSourceOptions, GrapQLQuery } from './types';

export class DataSource extends DataSourceWithBackend<GrapQLQuery, MyDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
  }

  // Enables an out-of-the-box annotaiton editor
  annotations = {};
}
