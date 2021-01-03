import { Router } from 'preact-router';
import { useState } from 'preact/hooks';

import ApolloClient from 'apollo-boost';
import { ApolloProvider } from '@apollo/react-hooks';

import gql from 'graphql-tag';
import { useQuery, useMutation } from '@apollo/react-hooks';

import { getCookie } from '../utils/cookie';

import Home from '../routes/home';
import Game from '../routes/game';
import Perferences from '../routes/preferences';

import { PlayerQueries } from '../queries/player'

const client = new ApolloClient();

function AppRoot() {
    const [user, setUser] = useState({ 
        id: getCookie("user-id")
    });

    const { loading, error, data } = useQuery(PlayerQueries.getState, {
        variables: { playerId: user.id },
        skip: !user.id
    });

    if (loading) return 'Loading...';
    if (error) return "Error!";

    const [userPreferences, setUserPreferences] = useState({ 
        show: (user.id == "") || (data.player == null)
    });

    if (userPreferences.show)
    {
        return (
            <Perferences setUser={setUser} setUserPreferences={setUserPreferences} />
        )
    }
    return (
        <div>
            <a href="#" onClick={e => {
                e.preventDefault();
                setUserPreferences({ show: true });
            }}>Change Name</a>
            <Router onChange={e => { this.currentUrl = e.url; }}>
                <Home path="/" user={data.player} />
                <Game path="/game/:id" user={data.player} />
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