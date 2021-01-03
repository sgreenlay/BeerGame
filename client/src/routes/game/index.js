import { useEffect } from 'preact/hooks';

import { useQuery, useMutation } from '@apollo/react-hooks';

import { GameQueries } from '../../queries/game'

function Game({ id }) {
    const { loading, error, data } = useQuery(GameQueries.getState, {
        variables: { gameId: id },
        skip: !id,
        pollInterval: 1000,
    });
    const [joinGame] = useMutation(GameQueries.joinGame, {
        variables: { 
            gameId: id
        },
        refetchQueries: [{
            query: GameQueries.getState,
            variables: { gameId: id }
        }],
        awaitRefetchQueries: true,
    });
    const [leaveGame] = useMutation(GameQueries.leaveGame, {
        variables: { 
            gameId: id
        },
        refetchQueries: [{
            query: GameQueries.getState,
            variables: { gameId: id }
        }],
        awaitRefetchQueries: true,
    });

    if (loading) return 'Loading...';
    if (error) return "Error!";

    useEffect(() => {
        joinGame({ variables: { playerId: this.props.user.id }});
    }, [this.props.user.id]);

    return (
        <div>
            <h1>'{this.props.id}'</h1>
            <h2>Players</h2>
            <ul>
                {data.game.players.map(player => (
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

export default Game;