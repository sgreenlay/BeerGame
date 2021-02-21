import { useState } from 'preact/hooks';

import gql from 'graphql-tag';
import { useQuery, useSubscription } from '@apollo/react-hooks';

import { GameQueries } from '../../gql/game'

function Home() {
    const [state, setState] = useState({ id: '' });
    const { loading, error, data, refetch } = useQuery(GameQueries.getExists,{
        variables: { id: state.id }
    });

    if (error) {
        console.log(error);
        return "Error!";
    }

    return (
        <div>
            <h1>Hello '{this.props.user.name}'!</h1>
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
                    {loading || !data.gameExists ? (
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