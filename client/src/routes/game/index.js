import { useEffect } from 'preact/hooks';

import { useQuery, useMutation } from '@apollo/react-hooks';

import Lobby from "./lobby"

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

    if (loading) return 'Loading...';
    if (error) return "Error!";

    useEffect(() => {
        joinGame({ variables: { playerId: this.props.user.id }});
    }, [this.props.user.id]);

    if (data.game.state.name == "lobby") {
        return (
            <div>
                <h1>'{this.props.id}'</h1>
                <Lobby user={this.props.user} game={data.game} />
            </div>
        );
    }

    return (
        <div>
            <h1>'{this.props.id}'</h1>
            TODO
        </div>
    )
}

export default Game;