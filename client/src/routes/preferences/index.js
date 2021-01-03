import { useState } from 'preact/hooks';

import gql from 'graphql-tag';
import { useQuery, useMutation } from '@apollo/react-hooks';

import { existsCookie, setCookie, getCookie } from '../../utils/cookie';
import { generateUUID } from '../../utils/uuid';

import { PlayerQueries } from '../../queries/player'

function Preferences() {
    const [state, setState] = useState({ name: '' });
    const [createPlayer] = useMutation(PlayerQueries.createPlayer);

    return (
        <div>
            <form onSubmit={e => {
                e.preventDefault();

                const userId = existsCookie("user-id") ? getCookie("user-id") : generateUUID();
                const userName = state.name;

                setCookie("user-id", userId);

                createPlayer({
                    variables: {
                        playerId: userId,
                        playerName: userName
                    },
                    refetchQueries: [{
                        query: PlayerQueries.getState,
                        variables: { playerId: userId }
                    }],
                    awaitRefetchQueries: true,
                });

                this.props.setUser({
                    id: userId
                });
                this.props.setUserPreferences({
                    show: false
                });
            }}>
                <input type="text" value={state.name} onInput={e => {
                    const { value } = e.target;
                    setState({ name: value })
                }} />
                <button type="submit">Set Name</button>
            </form>
        </div>
    );
}

export default Preferences;