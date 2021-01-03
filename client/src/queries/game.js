
import gql from 'graphql-tag';

export const GameQueries = {
    getExists: gql`
        query Game($id: String!) {
            gameExists(gameId: $id)
        }
    `,
    getState: gql`
        query Game($gameId: String!) {
            game(gameId: $gameId) {
                id
                players {
                    id
                    name
                }
                state {
                    name
                }
            }
        }
    `,
    joinGame: gql`
        mutation JoinGame($gameId: String!, $playerId: String!) {
            addPlayer(gameId: $gameId, playerId: $playerId)
        }
    `,
    leaveGame: gql`
        mutation LeaveGame($gameId: String!, $playerId: String!) {
            removePlayer(gameId: $gameId, playerId: $playerId)
        }
    `,
};
