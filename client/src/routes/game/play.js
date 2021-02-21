import { useState } from 'preact/hooks';

import { useQuery, useMutation } from '@apollo/react-hooks';

import { GameQueries } from '../../queries/game'

function Play() {
    const { loading, error, data } = useQuery(GameQueries.getPlayerState, {
        variables: { 
            gameId: this.props.game.id, 
            playerId: this.props.user.id 
        },
        pollInterval: 1000,
    });

    if (loading) return 'Loading...';
    if (error) return "Error!";
    
    const [state, setState] = useState({ 
        value: '',
        valid: false
    });
    const [setOutgoing] = useMutation(GameQueries.submitOutgoing, {
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

            {this.props.game.playerState.map(state => (
               <div>{state.player.name} : {state.outgoing == -1 ? "Waiting" : "Submitted"}</div> 
            ))}

            <div>Stock: { data.playerState.stock }</div>
            <div>Backlog: { data.playerState.backlog }</div>
            <div>Incoming: { data.playerState.incoming }</div>
            <div>Sent: { data.playerState.lastsent }</div>
            <div>Pending: { data.playerState.pending0 }</div>
            <div>Outstanding: { data.playerState.outstanding }</div>
            <div>OutgoingPrev: { data.playerState.outgoingprev.join(',') }</div>
            <div>Outgoing: <form style="display: inline;" onSubmit={e => {
                e.preventDefault();
                setOutgoing({ variables: { outgoing: state.value } });
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
        </div>
    );
}

export default Play;