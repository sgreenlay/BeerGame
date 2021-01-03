
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
                playerState {
                    player {
                        id
                        name
                    }
                    role {
                        name
                        value
                    }
                }
            }
        }
    `,
    getRoles: gql`
        query Roles {
            gameRoles {
                name
                value
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
    startGame: gql`
        mutation StartGame($gameId: String!) {
            startGame(gameId: $gameId)
        }
    `,
    setRole: gql`
        mutation ChangePlayerRole($gameId: String!, $playerId: String!, $role: Int!) {
            changePlayerRole(gameId: $gameId, playerId: $playerId, role: $role)
        }
    `,
    submitValue: gql`
        mutation ChangePlayerRole($gameId: String!, $playerId: String!, $value: Int!) {
            submitValue(gameId: $gameId, playerId: $playerId, value: $value)
        }
    `,
};
