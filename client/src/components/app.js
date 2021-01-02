import { Router } from 'preact-router';

import ApolloClient from 'apollo-boost';
import { ApolloProvider } from '@apollo/react-hooks';

import Home from '../routes/home';
import Game from '../routes/game';

const client = new ApolloClient();

function App() {
	return (
		<ApolloProvider client={client}>
			<div id="app">
				<Router onChange={e => {
					this.currentUrl = e.url;
				}}>
					<Home path="/" />
					<Game path="/game/:id" />
				</Router>
			</div>
		</ApolloProvider>
	);
}

export default App;