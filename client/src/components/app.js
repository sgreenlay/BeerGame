import { Router } from 'preact-router';
import { useState } from 'preact/hooks';

import ApolloClient from 'apollo-boost';
import { ApolloProvider } from '@apollo/react-hooks';

import gql from 'graphql-tag';
import { useQuery, useMutation } from '@apollo/react-hooks';

import { existsCookie, getCookie } from '../utils/cookie';

import Home from '../routes/home';
import Game from '../routes/game';
import Perferences from '../routes/preferences';

import { PlayerQueries } from '../queries/player'

const client = new ApolloClient();

function AppRoot() {
    const { loading, error, data } = useQuery(PlayerQueries.getState, {
        variables: { playerId: getCookie("user-id") },
        skip: !existsCookie("user-id")
    });

    if (loading) return 'Loading...';
    if (error) return "Error!";

    const [userPreferences, setUserPreferences] = useState({ 
        showPreferences: !existsCookie("user-id") || (data.player == null)
    });

    const player = data ? data.player : null
    if (userPreferences.showPreferences)
    {
        return (
            <Perferences user={player} setUserPreferences={setUserPreferences} />
        )
    }
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