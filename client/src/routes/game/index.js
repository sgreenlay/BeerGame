import gql from 'graphql-tag';
import { useQuery, useMutation } from '@apollo/react-hooks';

const POLL_GAME_STATE = gql`
	query Game($id: String!) {
		game(id: $id) {
			count
		}
	}
`;

const UPDATE_GAME_STATE = gql`
	mutation Increment($id: String!) {
		increment(id: $id)
	}
`;

function Game({ id }) {
	const { loading, error, data } = useQuery(POLL_GAME_STATE, {
		variables: { id },
		skip: !id,
		pollInterval: 1000,
	});

	const [increment] = useMutation(UPDATE_GAME_STATE, {
		variables: { id },
		refetchQueries: [{
			query: POLL_GAME_STATE,
			variables: { id }
		}],
		awaitRefetchQueries: true,
	});

	if (loading) return 'Loading...';
	if (error) return `Error! ${error.message}`;

	return (
		<div>
			<h1>Joined '{this.props.id}'</h1>
			<h2>{data.game.count}</h2>
			<button onClick={e => {
				e.preventDefault();
				increment();
			}}>Increment</button>
		</div>
	);
}

export default Game;