import { useState } from 'preact/hooks';

import { useQuery, useMutation } from '@apollo/react-hooks';

import { GameQueries } from '../../queries/game'

function Play() {
    const [state, setState] = useState({ 
        value: '',
        valid: false
    });
    const [setValue] = useMutation(GameQueries.submitValue, {
        variables: { 
            gameId: this.props.game.id,
            playerId: this.props.user.id
        },
        refetchQueries: [{
            query: GameQueries.getState,
            variables: { gameId: this.props.game.id, }
        }],
        awaitRefetchQueries: true,
    });

    return (
        <div>
            <h1>'{this.props.game.id}'</h1>
            <form onSubmit={e => {
                e.preventDefault();
                setValue({ variables: { value: state.value } });
                setState({ value: '', valid: false });
            }}>
                <input type="text" value={state.value} style={state.valid ? "border: 1px solid black" : "border: 1px solid red"} onInput={e => {
                    const { value } = e.target;

                    const intValue = Number(value)
                    var isValid = (intValue != NaN && intValue >= 0 && intValue < 2147483647);

                    setState({ value: value, valid: isValid });
                }} />
                <button disabled={!state.valid} type="submit">Set</button>
            </form>
        </div>
    );
}

export default Play;