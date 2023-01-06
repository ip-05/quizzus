package ws

import (
	"context"
	"errors"
	"github.com/ip-05/quizzus/middleware"
	"github.com/ip-05/quizzus/models"
	"nhooyr.io/websocket"
	"sync"
)

type GameStatus = string

const (
	Standby    = "GAME_STANDBY"
	Starting   = "GAME_STARTING"
	InProgress = "GAME_IN_PROGRESS"
	Finished   = "GAME_FINISHED"
)

const (
	GameNotFound  = "GAME_NOT_FOUND"
	AlreadyInGame = "ALREADY_IN_GAME"
	NotInGame     = "NOT_IN_GAME"

	JoinedGame  = "JOINED_GAME"
	LeftGame    = "LEFT_GAME"
	GameDeleted = "GAME_DELETED"
)

type User struct {
	Id         string          `json:"id"`
	Name       string          `json:"name"`
	ActiveGame *Game           `json:"-"`
	Conn       *websocket.Conn `json:"-"`
}

type Game struct {
	Status     string      `json:"status"`
	InviteCode string      `json:"inviteCode"`
	Members    []*User     `json:"members"`
	Owner      *User       `json:"owner"`
	Data       models.Game `json:"data"`
}

type gameSocketController struct {
	Users map[string]*User
	Games map[string]*Game

	mu sync.Mutex
}

func NewGameSocketController() *gameSocketController {
	c := new(gameSocketController)

	c.Users = make(map[string]*User)
	c.Games = make(map[string]*Game)

	return c
}

func (g *gameSocketController) InitUser(ctx context.Context) (*User, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	user := ctx.Value("authedUser").(middleware.AuthedUser)

	_, found := g.Users[user.Id]
	if found {
		return nil, errors.New("user already exists on another socket")
	}

	g.Users[user.Id] = &User{Id: user.Id, Name: user.Name, Conn: ctx.Value("conn").(*websocket.Conn)}

	return g.Users[user.Id], nil
}

type JoinGameData struct {
	GameId string `json:"gameId"`
}

func (g *gameSocketController) JoinGame(ctx context.Context, data JoinGameData) {
	g.mu.Lock()
	defer g.mu.Unlock()

	conn := ctx.Value("conn").(*websocket.Conn)

	user := ctx.Value("user").(*User)
	if user.ActiveGame != nil {
		MessageReply(true, AlreadyInGame).Send(conn)
		return
	}

	var game models.Game
	models.DB.Preload("Questions.Options").Where("invite_code = ?", data.GameId).First(&game)

	if game.Id == 0 {
		MessageReply(true, GameNotFound).Send(conn)
	}

	value, ok := g.Games[game.InviteCode]
	if ok {
		value.Members = append(value.Members, user)
		user.ActiveGame = value

		DataReply(false, JoinedGame, value).Send(conn)
		return
	}

	newGame := Game{
		Status:     Standby,
		InviteCode: game.InviteCode,
		Members:    []*User{user},
		Owner:      user,
		Data:       game,
	}

	g.Games[newGame.InviteCode] = &newGame
	user.ActiveGame = g.Games[newGame.InviteCode]
	DataReply(false, JoinedGame, newGame).Send(conn)
}

func (g *gameSocketController) LeaveGame(ctx context.Context) {
	g.mu.Lock()
	defer g.mu.Unlock()

	conn := ctx.Value("conn").(*websocket.Conn)

	user := ctx.Value("user").(*User)
	if user.ActiveGame == nil {
		MessageReply(true, NotInGame).Send(conn)
		return
	}

	game := user.ActiveGame
	if game.Owner == user {
		for _, member := range game.Members {
			member.ActiveGame = nil
			MessageReply(false, GameDeleted).Send(member.Conn)
		}
		delete(g.Games, game.InviteCode)
	}

	user.ActiveGame = nil
}

func (g *gameSocketController) GetGame(ctx context.Context) {
	conn := ctx.Value("conn").(*websocket.Conn)

	user := ctx.Value("user").(*User)
	if user.ActiveGame == nil {
		MessageReply(true, NotInGame).Send(conn)
		return
	}

	DataReply(false, GetGame, user.ActiveGame)
}

func (g *gameSocketController) IsOwner(ctx context.Context) {
	conn := ctx.Value("conn").(*websocket.Conn)

	user := ctx.Value("user").(*User)
	if user.ActiveGame == nil {
		MessageReply(true, NotInGame).Send(conn)
		return
	}

	DataReply(false, IsOwner, user.ActiveGame.Owner == user).Send(conn)
}
