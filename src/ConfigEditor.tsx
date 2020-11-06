import React from 'react';
import { DataSourceHttpSettings } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions } from './types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> {}

export default (props: Props) => {
  return (
    <DataSourceHttpSettings
      defaultUrl="http://localhost:9999"
      dataSourceConfig={props.options}
      onChange={props.onOptionsChange}
      showAccessOptions={true}
    />
  );
};
