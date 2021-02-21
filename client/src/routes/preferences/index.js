import { useState } from 'preact/hooks';

import { useQuery, useMutation } from '@apollo/react-hooks';

import { existsCookie, setCookie, getCookie } from '../../utils/cookie';
import { generateUUID } from '../../utils/uuid';

import { PlayerQueries } from '../../gql/player'

function Preferences() {
    const [state, setState] = useState({ 
        name: this.props.user ? this.props.user.name : '' 
    });
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

                this.props.setUserPreferences({
                    showPreferences: false
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