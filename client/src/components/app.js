import { Router } from 'preact-router';
import { useState } from 'preact/hooks';

import ApolloClient from "apollo-client";
import { WebSocketLink } from 'apollo-link-ws';
import { HttpLink } from 'apollo-link-http';
import { split } from 'apollo-link';
import { getMainDefinition } from 'apollo-utilities';
import { InMemoryCache } from 'apollo-cache-inmemory';

import { ApolloProvider } from '@apollo/react-hooks';
import { useQuery } from '@apollo/react-hooks';

import { existsCookie, getCookie } from '../utils/cookie';

import Home from '../routes/home';
import Game from '../routes/game';
import Perferences from '../routes/preferences';

import { PlayerQueries } from '../gql/player'

const httpLink = new HttpLink({
    uri: `http://${window.location.hostname}/graphql`,
});
const wsLink = new WebSocketLink({
    uri: `ws://${window.location.hostname}/wsgraphql`,
    options: {
        reconnect: true
    }
});

const splitLink = split(
    ({ query }) => {
        const definition = getMainDefinition(query);
        return (
            definition.kind === 'OperationDefinition' &&
            definition.operation === 'subscription'
        );
    },
    wsLink,
    httpLink,
);

const client = new ApolloClient({
    link: splitLink,
    cache: new InMemoryCache()
});

function AppRoot() {
    const { loading, error, data } = useQuery(PlayerQueries.getState, {
        variables: { playerId: getCookie("user-id") },
        skip: !existsCookie("user-id")
    });

    if (loading) return 'Loading...';
    if (error) {
        console.log(error);
        return "Error!";
    }

    const [userPreferences, setUserPreferences] = useState({
        showPreferences: !existsCookie("user-id") || (data.player == null)
    });

    const player = data ? data.player : null
    if (userPreferences.showPreferences) {
        return (
            <Perferences user={player} setUserPreferences={setUserPreferences} />
        )
    }

    if (player == null) return 'Loading...';

    return (
        <div>
            <a href="#" onClick={e => {
                e.preventDefault();
                setUserPreferences({ showPreferences: true });
            }}>Options</a>
            <Router onChange={e => { this.currentUrl = e.url; }}>
                <Home path="/" user={player} />
                <Game path="/game/:id" user={player} />
            </Router>
        </div>
    )
}

function App() {
    return (
        <ApolloProvider client={client}>
            <div id="app">
                <AppRoot />
            </div>
        </ApolloProvider>
    );
}

export default App;