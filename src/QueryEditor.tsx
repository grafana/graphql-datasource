import React, { PureComponent } from 'react';

import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './DataSource';
import { MyDataSourceOptions, MyQuery } from './types';

import GraphiQL from 'graphiql';
import 'graphiql/graphiql.min.css';
import { FetcherParams } from 'graphiql/dist/components/GraphiQL';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  async getData(graphQLParams: FetcherParams) {
    const response = await fetch(
      'https://countries.trevorblades.com/',
      {
        method: 'POST',
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(graphQLParams),
        credentials: 'same-origin',
      }
    );

    const queryResults = await response.json().catch(() => response.text());

    return queryResults;
  }

  render() {
    const divStyle = {
      height: '100vh',
    };

    return (
      <div style={divStyle}>
        <GraphiQL fetcher={this.getData} />
      </div>
    );
  }
}
