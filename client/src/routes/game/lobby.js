import { useQuery, useMutation } from '@apollo/react-hooks';

import { GameQueries } from '../../gql/game'

function Lobby() {
    const { loading, error, data } = useQuery(GameQueries.getRoles);

    if (loading) return 'Loading...';
    if (error) {
        console.log(error);
        return "Error!";
    }

    const [leaveGame] = useMutation(GameQueries.leaveGame, {
        variables: { 
            gameId: this.props.game.id
        },
    });
    const [startGame] = useMutation(GameQueries.startGame, {
        variables: { 
            gameId: this.props.game.id
        },
    });
    const [changeRole] = useMutation(GameQueries.setRole, {
        variables: { 
            gameId: this.props.game.id
        },
    });

    return (
        <div>
        <h1>'{this.props.game.id}'</h1>
            <ul>
                {this.props.game.playerState.map(state => (
                    <li>
                        <span>{state.player.name}</span>
                        &nbsp;
                        <span>
                            {state.player.id == this.props.user.id ? (
                                <select value={state.role.value} onChange={e => {
                                    e.preventDefault();
                                    changeRole({ variables: { playerId: state.player.id, role: e.target.value }});
                                }}>{data.gameRoles.map(role => (
                                    <option value={role.value}>{role.name}</option>
                                ))}</select>
                            ) : (
                                <span>{state.role.name}</span>
                            )}
                        </span>
                        &nbsp;
                        <span>
                            [<a href="#" onClick={e => {
                                e.preventDefault();
                                leaveGame({ variables: { playerId: state.player.id }});
                            }}>{state.player.id == this.props.user.id ? (
                                'Leave'
                            ) : (
                                'Kick'
                            )}</a>]
                        </span>
                    </li>
                ))}
            </ul>
            <a href="#" onClick={e => {
                e.preventDefault();
                startGame();
            }}>Start</a>
        </div>
    );
}

export default Lobby;