import { useEffect } from 'preact/hooks';

import { useQuery, useMutation, useSubscription } from '@apollo/react-hooks';

import Lobby from "./lobby"
import Play from "./play"

import { GameQueries, GameSubscriptions } from '../../gql/game'

function Game({ id }) {
    const { loading, error, data } = useSubscription(GameSubscriptions.gameState, {
        variables: { gameId: id },
        shouldResubscribe: true
    });
    const [joinGame] = useMutation(GameQueries.joinGame, {
        variables: { 
            gameId: id
        }
    });

    if (loading) return 'Loading...';
    if (error) {
        console.log(error);
        return "Error!";
    }

    useEffect(() => {
        joinGame({ variables: { playerId: this.props.user.id }});
    }, [this.props.user.id]);

    if (data.game.state.name == "lobby") {
        return (
            <Lobby user={this.props.user} game={data.game} />
        );
    } else if (data.game.state.name == "playing") {
        return (
            <Play user={this.props.user} game={data.game} />
        );
    } else if (data.game.state.name == "finished") {
        <div>
            <h1>'{this.props.id}'</h1>
            TODO: Finished
        </div>
    }
    return (
        <div>
            <h1>'{this.props.id}'</h1>
            ERROR: unknown state '{data.game.state.name}''
        </div>
    )
}

export default Game;