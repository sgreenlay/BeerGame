import { useState } from 'preact/hooks';

import gql from 'graphql-tag';
import { useQuery } from '@apollo/react-hooks';

const EXISTING_GAME = gql`
	query Game($id: String!) {
		exists(id: $id)
	}
`;

function Home() {
	const [state, setState] = useState({ id: '' });
	const { loading, error, data, refetch } = useQuery(EXISTING_GAME,{
		variables: { id: state.id }
	});

	if (error) return `Error! ${error.message}`;

	return (
		<div>
			<form onSubmit={e => {
				e.preventDefault();
				if (state.id != '') {
					window.location.assign("/game/" + state.id);
				}
			}}>
				<input type="text" value={state.id} onInput={e => {
					const { value } = e.target;
					setState({ id: value })
					refetch();
				}} />
				<button type="submit">
					{loading || !data.exists ? (
						'Create'
					) : (
						'Join'
					)}
				</button>
			</form>
		</div>
	);
}

export default Home;