import React, { useState } from 'react';

import {  QueryEditorProps } from '@grafana/data';
import { CodeEditor } from '@grafana/ui';
import { DataSource } from './DataSource';
import { MyDataSourceOptions, MyQuery } from './types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export default (props: Props) => {
  const { query } = props.query;

  return (
    <>
      <CodeEditor
        language={'graphql'}
        value={query}
        showMiniMap={false}
        showLineNumbers={true}
        height={'250px'}
        onBlur={(value) => props.onChange({
          ...props.query,
          query: value,
        })}
      />
    </>
  );
}
