import { useState } from 'preact/hooks';

import { useQuery, useMutation, useSubscription } from '@apollo/react-hooks';

import { GameQueries, GameSubscriptions } from '../../gql/game'

function Play() {
    const { loading, error, data } = useSubscription(GameSubscriptions.playerState, {
        variables: {
            gameId: this.props.game.id,
            playerId: this.props.user.id
        },
        shouldResubscribe: true
    });

    if (loading) return 'Loading...';
    if (error) {
        console.log(error);
        return "Error!";
    }

    const [state, setState] = useState({
        value: '',
        valid: false
    });
    const [setOutgoing] = useMutation(GameQueries.submitOutgoing, {
        variables: {
            gameId: this.props.game.id,
            playerId: this.props.user.id
        },
    });

    return (
        <div>
            <h1>'{this.props.game.id}'</h1>

            <div class="player-state">
                {this.props.game.playerState.sort(function(a, b) {
                    return a.role.value - b.role.value;
                }).map(state => (
                    <div class={"block " + (state.outgoing == -1 ? "waiting" : "done")}>
                        {state.player.name}
                        <div class="role">{state.role.name}</div>
                    </div>
                ))}
            </div>

            <div class="game-state">
                <div class="block incoming">
                    <span class="title">Incoming</span>
                    <span class="value">{ data.playerState.incoming }</span>
                </div>
                <div class="block outgoing">
                    <span class="title">Outgoing</span>
                    <form class="value" onSubmit={e => {
                        e.preventDefault();
                        setOutgoing({ variables: { outgoing: state.value } });
                        setState({ value: '', valid: false });
                    }}>
                        <input type="text" value={state.value} class={
                            state.valid ? "input valid" : "input invalid"
                        } onInput={e => {
                            const { value } = e.target;

                            const intValue = Number(value)
                            var isValid = (value.length > 0) && (intValue != NaN && intValue >= 0 && intValue < 2147483647);

                            setState({ value: value, valid: isValid });
                        }} />
                    </form>
                </div>
                <div class="block backlog">
                    <span class="title">Backlog</span>
                    <span class="value">{ data.playerState.backlog }</span>
                </div>
                <div class="block stock">
                    <span class="title">Stock</span>
                    <span class="value">{ data.playerState.stock }</span>
                </div>
                <div class="indicator last-sent-indicator">
                    <span class="title">&#9654;</span>
                </div>
                <div class="block last-sent">
                    <span class="title">Last Sent</span>
                    <span class="value">{ data.playerState.lastsent }</span>
                </div>
                <div class="indicator pending-next-indicator">
                    <span class="title">&#9664;</span>
                </div>
                <div class="block pending-next">
                    <span class="title">Pending</span>
                    <span class="value">{ data.playerState.pending0 }</span>
                </div>
                <div class="block pending-all">
                    <span class="title">Outstanding</span>
                    <span class="value">{ data.playerState.outstanding }</span>
                </div>
            </div>
        </div>
    );
}

export default Play;