import { notFound } from 'next/navigation';
import { CodeKey } from '@inngest/components/CodeKey';

import { graphql } from '@/gql';
import graphqlAPI from '@/queries/graphqlAPI';
import { getEnvironment } from '@/queries/server-only/getEnvironment';
import { Provider } from './Context';
import DeleteKeyButton from './DeleteKeyButton';
import EditKeyName from './EditKeyName';
import FilterEvents from './FilterEvents';
import TransformEvent from './TransformEvent';

const GetKeyDocument = graphql(`
  query GetIngestKey($environmentID: ID!, $keyID: ID!) {
    environment: workspace(id: $environmentID) {
      ingestKey(id: $keyID) {
        id
        name
        createdAt
        presharedKey
        url
        filter {
          type
          ips
          events
        }
        metadata
      }
    }
  }
`);

type KeyDetailsProps = {
  params: {
    environmentSlug: string;
    ingestKeys: string;
    keyID: string;
  };
};

export const runtime = 'nodejs';

export default async function Keys({
  params: { environmentSlug, ingestKeys, keyID },
}: KeyDetailsProps) {
  const environment = await getEnvironment({
    environmentSlug: environmentSlug,
  });

  const response = await graphqlAPI.request(GetKeyDocument, {
    environmentID: environment.id,
    keyID,
  });

  if (!response) {
    notFound();
  }

  const key = response.environment.ingestKey;

  if (!key) {
    return null;
  }
  const filterType = key.filter.type;
  if (!filterType || !isFilterType(filterType)) {
    throw new Error(`invalid filter type: ${filterType}`);
  }

  let value = '',
    maskedValue = '',
    keyLabel = 'Key';
  if (ingestKeys === 'webhooks') {
    value = key.url || '';
    // Leave the base url + the beginning of the key
    maskedValue = value.replace(/(.{0,}\/e\/)(\w{0,8}).+/, '$1$2');
    keyLabel = 'Webhook URL';
  } else {
    value = key.presharedKey;
    maskedValue = value.substring(0, 8);
    keyLabel = 'Event Key';
  }

  return (
    <div className="m-6 divide-y divide-slate-100">
      <Provider initialState={key}>
        <div className="pb-8">
          <div className="mb-8 flex flex-wrap justify-between">
            <EditKeyName keyID={keyID} keyName={key.name} />
            <DeleteKeyButton
              environmentSlug={environmentSlug}
              environmentID={environment.id}
              keyID={keyID}
            />
          </div>
          <div className="w-3/5">
            <CodeKey fullKey={value} maskedKey={maskedValue} label={keyLabel} />
          </div>
        </div>
        <TransformEvent keyID={keyID} metadata={key.metadata} keyName={key.name} />
        <FilterEvents
          keyID={keyID}
          filter={{
            ...key.filter,
            type: filterType,
          }}
          keyName={key.name}
        />
      </Provider>
    </div>
  );
}

function isFilterType(value: string): value is 'allow' | 'deny' {
  return value === 'allow' || value === 'deny';
}
