package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/rs/cors"
)

type NameValueMapping struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var Players = map[string]*Player{}

const (
	LOBBY = iota
	PLAYING
	FINISHED
)

var GameStateMappings = []NameValueMapping{
	NameValueMapping{
		Name:  "lobby",
		Value: LOBBY,
	},
	NameValueMapping{
		Name:  "playing",
		Value: PLAYING,
	},
	NameValueMapping{
		Name:  "finished",
		Value: FINISHED,
	},
}

const (
	NONE = iota
	RETAILER
	WHOLESALER
	DISTRIBUTER
	MANUFACTURER
)

var GameRoleMappings = []NameValueMapping{
	NameValueMapping{
		Name:  "none",
		Value: NONE,
	},
	NameValueMapping{
		Name:  "retailer",
		Value: RETAILER,
	},
	NameValueMapping{
		Name:  "wholesaler",
		Value: WHOLESALER,
	},
	NameValueMapping{
		Name:  "distributer",
		Value: DISTRIBUTER,
	},
	NameValueMapping{
		Name:  "manufacturer",
		Value: MANUFACTURER,
	},
}

type PlayerState struct {
	PlayerID string `json:"playerId"`
	Role     int    `json:"role"`
	Incoming int    `json:"incoming"`
	Outgoing int    `json:"outgoing"`
	Stock    int    `json:"stock"`
	Backlog  int    `json:"backlog"`
}

type Game struct {
	ID          string         `json:"id"`
	State       int            `json:"state"`
	PlayerState []*PlayerState `json:"playerState"`
}

var Games = map[string]*Game{}

func FindGame(id string) *Game {
	game, _ := Games[id]
	return game
}

func ExistsGame(id string) bool {
	_, found := Games[id]
	return found
}

func FindOrCreateGame(id string) *Game {
	game, found := Games[id]
	if !found {
		newGame := &Game{
			ID:          id,
			State:       LOBBY,
			PlayerState: []*PlayerState{},
		}
		Games[id] = newGame
		return newGame
	}
	return game
}

func FindPlayer(id string) *Player {
	player, _ := Players[id]
	return player
}

func FindOrCreatePlayer(id string, name string) *Player {
	player, found := Players[id]
	if !found {
		newPlayer := &Player{
			ID:   id,
			Name: name,
		}
		Players[id] = newPlayer
		return newPlayer
	}
	player.Name = name
	return player
}

func (game *Game) AddPlayer(id string) bool {
	for _, value := range game.PlayerState {
		if value.PlayerID == id {
			return false
		}
	}
	newPlayerState := &PlayerState{
		PlayerID: id,
		Incoming: -1,
		Outgoing: -1,
		Stock:    -1,
		Backlog:  -1,
	}
	game.PlayerState = append(game.PlayerState, newPlayerState)
	return true
}

func (game *Game) RemovePlayer(id string) bool {
	for index, playerState := range game.PlayerState {
		if playerState.PlayerID == id {
			game.PlayerState = append(game.PlayerState[:index], game.PlayerState[index+1:]...)
			return true
		}
	}
	return false
}

func (game *Game) FindPlayerState(id string) *PlayerState {
	for _, playerState := range game.PlayerState {
		if playerState.PlayerID == id {
			return playerState
		}
	}
	return nil
}

func (game *Game) Start() bool {
	if game.State == LOBBY {
		// TODO: Validation

		game.State = PLAYING

		return true
	}
	return false
}

func (game *Game) TryStep() bool {
	// TODO: Game logic

	for _, playerState := range game.PlayerState {
		playerState.Outgoing = -1
	}
	return false
}

var nameValueType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NameValue",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"value": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var playerType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Player",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var playerStateType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PlayerState",
	Fields: graphql.Fields{
		"player": &graphql.Field{
			Type: playerType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				playerState := p.Source.(*PlayerState)
				return FindPlayer(playerState.PlayerID), nil
			},
		},
		"incoming": &graphql.Field{
			Type: graphql.Int,
		},
		"stock": &graphql.Field{
			Type: graphql.Int,
		},
		"backlog": &graphql.Field{
			Type: graphql.Int,
		},
		"role": &graphql.Field{
			Type: nameValueType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				playerState := p.Source.(*PlayerState)
				return GameRoleMappings[playerState.Role], nil
			},
		},
	},
})

var gameType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Game",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"players": &graphql.Field{
				Type: graphql.NewList(playerType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					game := p.Source.(*Game)
					players := []*Player{}
					for _, playerState := range game.PlayerState {
						player := FindPlayer(playerState.PlayerID)
						if player != nil {
							players = append(players, player)
						}
					}
					return players, nil
				},
			},
			"state": &graphql.Field{
				Type: nameValueType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					game := p.Source.(*Game)
					return GameStateMappings[game.State], nil
				},
			},
			"playerState": &graphql.Field{
				Type: graphql.NewList(playerStateType),
			},
		},
	},
)

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"gameExists": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"gameId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["gameId"].(string)
				return ExistsGame(id), nil
			},
		},
		"game": &graphql.Field{
			Type: gameType,
			Args: graphql.FieldConfigArgument{
				"gameId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["gameId"].(string)
				return FindOrCreateGame(id), nil
			},
		},
		"playerState": &graphql.Field{
			Type: playerStateType,
			Args: graphql.FieldConfigArgument{
				"gameId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"playerId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["gameId"].(string)

				game := FindGame(id)
				if game == nil {
					return nil, nil
				}

				playerId, _ := p.Args["playerId"].(string)
				playerState := game.FindPlayerState(playerId)
				if playerState == nil {
					return nil, nil
				}

				return playerState, nil
			},
		},
		"player": &graphql.Field{
			Type: playerType,
			Args: graphql.FieldConfigArgument{
				"playerId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				playerId, _ := p.Args["playerId"].(string)
				return FindPlayer(playerId), nil
			},
		},
		"gameStates": &graphql.Field{
			Type: graphql.NewList(nameValueType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return GameStateMappings, nil
			},
		},
		"gameRoles": &graphql.Field{
			Type: graphql.NewList(nameValueType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return GameRoleMappings, nil
			},
		},
	},
})

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createPlayer": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"playerId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"playerName": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				playerId, _ := p.Args["playerId"].(string)
				playerName, _ := p.Args["playerName"].(string)
				return FindOrCreatePlayer(playerId, playerName), nil
			},
		},
		"addPlayer": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"gameId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"playerId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				gameId, _ := p.Args["gameId"].(string)
				game := FindOrCreateGame(gameId)
				playerId, _ := p.Args["playerId"].(string)
				player := FindPlayer(playerId)
				if player == nil {
					return false, nil
				}
				added := game.AddPlayer(playerId)
				return added, nil
			},
		},
		"removePlayer": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"gameId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"playerId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				gameId, _ := p.Args["gameId"].(string)
				game := FindOrCreateGame(gameId)
				playerId, _ := p.Args["playerId"].(string)
				removed := game.RemovePlayer(playerId)
				return removed, nil
			},
		},
		"changePlayerRole": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"gameId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"playerId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"role": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				gameId, _ := p.Args["gameId"].(string)
				game := FindGame(gameId)
				if game == nil {
					return false, nil
				}

				playerId, _ := p.Args["playerId"].(string)
				playerState := game.FindPlayerState(playerId)
				if playerState == nil {
					return false, nil
				}

				role, _ := p.Args["role"].(int)
				playerState.Role = role

				return true, nil
			},
		},
		"startGame": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"gameId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				gameId, _ := p.Args["gameId"].(string)
				game := FindOrCreateGame(gameId)
				return game.Start(), nil
			},
		},
		"submitOutgoing": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"gameId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"playerId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"outgoing": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				gameId, _ := p.Args["gameId"].(string)
				game := FindGame(gameId)
				if game == nil {
					return false, nil
				}

				playerId, _ := p.Args["playerId"].(string)
				playerState := game.FindPlayerState(playerId)
				if playerState == nil {
					return false, nil
				}

				outgoing, validOutgoing := p.Args["outgoing"].(int)
				if !validOutgoing {
					return false, nil
				}

				playerState.Outgoing = outgoing
				game.TryStep()

				return true, nil
			},
		},
	},
})

type SinglePageAppHandler struct {
	Directory string
	IndexFile string
}

func (h SinglePageAppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
		r.URL.Path = path
	}

	path = filepath.Join(h.Directory, r.URL.Path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.Directory, h.IndexFile))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.FileServer(http.Dir(h.Directory)).ServeHTTP(w, r)
}

func main() {
	mux := http.NewServeMux()

	appHandler := SinglePageAppHandler{
		Directory: "static",
		IndexFile: "index.html",
	}
	mux.Handle("/", appHandler)

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})

	graphqlHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	mux.Handle("/graphql", graphqlHandler)
	mux.Handle("/graphql/", graphqlHandler)

	handler := cors.Default().Handler(mux)
	http.ListenAndServe("0.0.0.0:80", handler)
}
