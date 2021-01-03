package main

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/rs/cors"
	"os"
	"strings"
	"net/http"
	"path/filepath"
)

type Player struct {
	ID		string	`json:"id"`
	Name	string	`json:"name"`
}

type Game struct {
	ID    	string		`json:"id"`
	Players []string	`json:"players"`
}

var Games = map[string]*Game{}
var Players = map[string]*Player{}

func ExistingGame(id string) bool {
	_, found := Games[id]
	return found
}

func FindOrCreateGame(id string) *Game {
	game, found := Games[id]
	if (!found) {
		newGame := &Game{
			ID: id,
			Players: []string{},
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
	if (!found) {
		newPlayer := &Player{
			ID: id,
			Name: name,
		}
		Players[id] = newPlayer
		return newPlayer
	}
	player.Name = name
	return player
}

func (game *Game) AddPlayer(id string) bool {
	for _, value := range game.Players {
		if value == id {
			return false
		}
	}
	game.Players = append(game.Players, id)
	return true
}

func (game *Game) RemovePlayer(id string) bool {
	for index, value := range game.Players {
		if value == id {
			game.Players = append(game.Players[:index], game.Players[index+1:]...)
			return true
		}
	}
	return false
}

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
					for _, playerId := range game.Players {
						player := FindPlayer(playerId)
						if player != nil {
							players = append(players, player)
						}
					}
					return players, nil
				},
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
				return ExistingGame(id), nil
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
				if (player == nil) {
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
	},
})

type SinglePageAppHandler struct {
	Directory string
	IndexFile  string
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
		Query: queryType,
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
