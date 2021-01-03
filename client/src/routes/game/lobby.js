import { useEffect } from 'preact/hooks';

import { useQuery, useMutation } from '@apollo/react-hooks';

import { GameQueries } from '../../queries/game'

function Lobby() {
    const [leaveGame] = useMutation(GameQueries.leaveGame, {
        variables: { 
            gameId: this.props.game.id
        },
        refetchQueries: [{
            query: GameQueries.getState,
            variables: { gameId: this.props.game.id }
        }],
        awaitRefetchQueries: true,
    });

    return (
        <div>
            <h2>Players</h2>
            <ul>
                {this.props.game.players.map(player => (
                    <li>
                        {player.name} [<a href="#" onClick={e => {
                            e.preventDefault();
                            leaveGame({ variables: { playerId: player.id }});
                        }}>{player.id == this.props.user.id ? (
                            'Leave'
                        ) : (
                            'Kick'
                        )}</a>]
                    </li>
                ))}
            </ul>
        </div>
    );
}

export default Lobby;