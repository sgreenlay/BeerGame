package main

import (
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/websocket"

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
	PlayerID      string `json:"playerId"`
	Role          int    `json:"role"`
	Incoming      int    `json:"incoming"`
	Outgoing      int    `json:"outgoing"`
	Outstanding   int    `json:"outstanding"`
	LastSent      int    `json:"lastsent"`
	Stock         int    `json:"stock"`
	Backlog       int    `json:"backlog"`
	Pending0      int    `json:"pending0"`
	Pending1      int    `json:"pending1"`
	Costs         int    `json:"costs"`
	OutgoingPrev  []int  `json:"outgoingprev"`
	StockBackPrev []int  `json:"stockbackprev"`
	CostPrev      []int  `json:"costprev"`
}

type Game struct {
	ID          string         	`json:"id"`
	State       int            	`json:"state"`
	PlayerState []*PlayerState 	`json:"playerState"`
	Week        int            	`json:"week"`
	LastWeek    int        		`json:"lastweek"`
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
			Week:        0,
			LastWeek:	 50,
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
	if game.State != LOBBY {
		return false
	}
	for _, value := range game.PlayerState {
		if value.PlayerID == id {
			return false
		}
	}
	newPlayerState := &PlayerState{
		PlayerID:      id,
		Incoming:      0,
		Outgoing:      -1,
		Outstanding:   0,
		LastSent:      0,
		Stock:         15,
		Backlog:       0,
		Pending0:      0,
		Pending1:      0,
		Costs:         0,
		OutgoingPrev:  []int{},
		StockBackPrev: []int{},
		CostPrev:      []int{},
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
		var p1 *PlayerState = nil
		var p2 *PlayerState = nil
		var p3 *PlayerState = nil
		var p4 *PlayerState = nil

		for _, playerState := range game.PlayerState {
			if playerState.Role == RETAILER {
				p1 = playerState
			}
			if playerState.Role == WHOLESALER {
				p2 = playerState
			}
			if playerState.Role == DISTRIBUTER {
				p3 = playerState
			}
			if playerState.Role == MANUFACTURER {
				p4 = playerState
			}
		}

		if p1 == nil || p2 == nil || p3 == nil || p4 == nil || len(game.PlayerState) != 4 {
			return false
		}

		game.State = PLAYING

		return true
	}
	return false
}

func (game *Game) TryStep() bool {
	if game.State != PLAYING {
		return false
	}

	var p1 *PlayerState = nil
	var p2 *PlayerState = nil
	var p3 *PlayerState = nil
	var p4 *PlayerState = nil

	for _, playerState := range game.PlayerState {
		if playerState.Outgoing == -1 {
			return false
		}
		if playerState.Role == RETAILER {
			p1 = playerState
		}
		if playerState.Role == WHOLESALER {
			p2 = playerState
		}
		if playerState.Role == DISTRIBUTER {
			p3 = playerState
		}
		if playerState.Role == MANUFACTURER {
			p4 = playerState
		}
	}

	if p1 == nil || p2 == nil || p3 == nil || p4 == nil || len(game.PlayerState) != 4 {
		return false
	}

	p1.Incoming = rand.Intn(20) // Customers
	p2.Incoming = p1.Outgoing
	p3.Incoming = p2.Outgoing
	p4.Incoming = p3.Outgoing

	for _, p := range game.PlayerState {
		p.Backlog = p.Backlog + p.Incoming
		p.Outstanding = p.Outstanding + p.Outgoing - p.Pending0
		p.Stock = p.Stock + p.Pending0
		p.Pending0 = p.Pending1

		if p.Stock > p.Backlog {
			p.LastSent = p.Backlog
			p.Stock = p.Stock - p.Backlog
			p.Backlog = 0
		} else {
			p.LastSent = p.Stock
			p.Backlog = p.Backlog - p.Stock
			p.Stock = 0
		}
		p.Costs = p.Costs + p.Stock + p.Backlog*2
	}

	p1.Pending1 = p2.LastSent
	p2.Pending1 = p3.LastSent
	p3.Pending1 = p4.LastSent
	p4.Pending1 = p4.Outgoing

	for _, p := range game.PlayerState {
		p.OutgoingPrev = append(p.OutgoingPrev, p.Outgoing)
		p.StockBackPrev = append(p.StockBackPrev, p.Stock-p.Backlog)
		p.CostPrev = append(p.CostPrev, p.Costs)
		p.Outgoing = -1
	}

	if game.Week >= game.LastWeek - 1 {
		game.State = FINISHED
	} else {
		game.Week = game.Week + 1
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

var publicPlayerStateType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PublicPlayerState",
	Fields: graphql.Fields{
		"player": &graphql.Field{
			Type: playerType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				playerState := p.Source.(*PlayerState)
				return FindPlayer(playerState.PlayerID), nil
			},
		},
		"outgoing": &graphql.Field{
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

var privatePlayerStateType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PrivatePlayerState",
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
		"lastsent": &graphql.Field{
			Type: graphql.Int,
		},
		"pending0": &graphql.Field{
			Type: graphql.Int,
		},
		"costs": &graphql.Field{
			Type: graphql.Int,
		},
		"outstanding": &graphql.Field{
			Type: graphql.Int,
		},
		"outgoingprev": &graphql.Field{
			Type: graphql.NewList(graphql.Int),
		},
		"stockbackprev": &graphql.Field{
			Type: graphql.NewList(graphql.Int),
		},
		"costprev": &graphql.Field{
			Type: graphql.NewList(graphql.Int),
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
				Type: graphql.NewList(publicPlayerStateType),
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
			Type: privatePlayerStateType,
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
				Subscriptions.broadcast()
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
				Subscriptions.broadcast()
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
				Subscriptions.broadcast()
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
				started := game.Start()
				Subscriptions.broadcast()
				return started, nil
			},
		},
		"submitLastWeek": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"gameId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"lastWeek": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				gameId, _ := p.Args["gameId"].(string)
				game := FindGame(gameId)
				if game == nil {
					return false, nil
				}

				if game.State != LOBBY {
					return false, nil
				}

				lastWeek, validLastWeek := p.Args["lastWeek"].(int)
				if !validLastWeek {
					return false, nil
				}

				game.LastWeek = lastWeek
				Subscriptions.broadcast()
				return true, nil
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
				Subscriptions.broadcast()
				return true, nil
			},
		},
	},
})

var subscriptionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Subscription",
	Fields: graphql.Fields{
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
			Type: privatePlayerStateType,
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

type Subscriber struct {
	ID            int
	Conn          *websocket.Conn
	RequestString string
	Variables 	  map[string]interface{}
	OperationID   string
}

type SubscriptionHandler struct {
	Schema	*graphql.Schema
	NextID	int
	Subscribers	[]Subscriber
}
var Subscriptions SubscriptionHandler

type SubscriptionMessage struct {
	OperationID string `json:"id,omitempty"`
	Type        string `json:"type"`
	Payload     struct {
		Query 		string `json:"query"`
		Variables 	map[string]interface{} `json:"variables"`
	} `json:"payload,omitempty"`
}

func (h *SubscriptionHandler) uniqueId() int {
	id := h.NextID
	h.NextID += 1
	return id
}

func (h *SubscriptionHandler) handler(ws *websocket.Conn) {
	for {
		var msg SubscriptionMessage
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			break
		}

		switch msg.Type {
			case "connection_init":
			case "start":
				subscriber := Subscriber{
					ID:            h.uniqueId(),
					Conn:          ws,
					RequestString: 	msg.Payload.Query,
					Variables: msg.Payload.Variables,
					OperationID:   msg.OperationID,
				}
				h.Subscribers = append(h.Subscribers, subscriber)
				go h.initilizeSubscriber(&subscriber)
			case "stop":
			default:
				println("Unknown message:", msg.Type)
		}
	}
}

func (h *SubscriptionHandler) removeSubscriber(id int) {
	for index, subscriber := range h.Subscribers {
		if subscriber.ID == id {
			h.Subscribers = append(h.Subscribers[:index], h.Subscribers[index+1:]...)
			return
		}
	}
}

func (subscriber *Subscriber) broadcast(schema *graphql.Schema) bool {
	payload := graphql.Do(graphql.Params{
		Schema: *schema,
		RequestString: subscriber.RequestString,
		VariableValues: subscriber.Variables,
	})
	msg := map[string]interface{}{
		"type":    "data",
		"id":      subscriber.OperationID,
		"payload": payload,
	}
	if err := websocket.JSON.Send(subscriber.Conn, msg); err != nil {
		return false
	}
	return true
}

func (h *SubscriptionHandler) initilizeSubscriber(subscriber *Subscriber) {
	time.Sleep(100 * time.Millisecond)
	succeeded := subscriber.broadcast(h.Schema)
	if !succeeded {
		h.removeSubscriber(subscriber.ID)
	}
}

func (h *SubscriptionHandler) broadcast() {
	invalidSubscriptions := []int{}
	for _, subscriber := range h.Subscribers {
		succeeded := subscriber.broadcast(h.Schema)
		if !succeeded {
			invalidSubscriptions = append(invalidSubscriptions, subscriber.ID)
		}
	}

	for _, id := range invalidSubscriptions {
		h.removeSubscriber(id)
	}
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
		Subscription: subscriptionType,
	})

	graphqlHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	mux.Handle("/graphql", graphqlHandler)
	mux.Handle("/graphql/", graphqlHandler)

	Subscriptions = SubscriptionHandler{
		Schema:   &schema,
	}
	mux.Handle("/wsgraphql", websocket.Handler(Subscriptions.handler))

	handler := cors.Default().Handler(mux)
	http.ListenAndServe("0.0.0.0:80", handler)
}
