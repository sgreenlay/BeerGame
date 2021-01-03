
import gql from 'graphql-tag';

export const PlayerQueries = {
    getState: gql`
        query Player($playerId: String!) {
            player(playerId: $playerId) {
                id
                name
            }
        }
    `,
    createPlayer: gql`
        mutation CreatePlayer($playerId: String!, $playerName: String!) {
            createPlayer(playerId: $playerId, playerName: $playerName)
        }
    `,
};
