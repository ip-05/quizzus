package ws

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"sync"
	"time"

	"github.com/ip-05/quizzus/middleware"
	"github.com/ip-05/quizzus/models"
	"nhooyr.io/websocket"
)

const (
	Standby    = "GAME_STANDBY"
	Starting   = "GAME_STARTING"
	InProgress = "GAME_IN_PROGRESS"
	Finished   = "GAME_FINISHED"
)

const (
	RoundInProgress = "ROUND_IN_PROGRESS"
	RoundWaiting    = "ROUND_WAITING"
)

const (
	GameNotFound  = "GAME_NOT_FOUND"
	AlreadyInGame = "ALREADY_IN_GAME"
	NotInGame     = "NOT_IN_GAME"
	NotOwner      = "NOT_OWNER"

	JoinedGame  = "JOINED_GAME"
	LeftGame    = "LEFT_GAME"
	GameDeleted = "GAME_DELETED"
	UserLeft    = "USER_LEFT"
	UserJoined  = "USER_JOINED"
)

type User struct {
	Id         string          `json:"id"`
	Name       string          `json:"name"`
	ActiveGame *Game           `json:"-"`
	Conn       *websocket.Conn `json:"-"`
}

type Game struct {
	Status       string             `json:"status"`
	RoundStatus  string             `json:"roundStatus"`
	CurrentRound int                `json:"currentRound"`
	InviteCode   string             `json:"inviteCode"`
	Members      map[string]*User   `json:"members"`
	Owner        *User              `json:"owner"`
	Leaderboard  map[string]float64 `json:"leaderboard"`
	Data         models.Game        `json:"-"`
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

func (g *gameSocketController) CleanUser(ctx context.Context) {
	g.mu.Lock()
	defer g.mu.Unlock()

	user := ctx.Value("user").(*User)
	if user.ActiveGame != nil {
		g.LeaveGame(ctx)
	}

	delete(g.Users, user.Id)
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
		return
	}

	value, ok := g.Games[game.InviteCode]
	if ok {
		for _, member := range value.Members {
			DataReply(false, UserJoined, member).Send(member.Conn)
		}

		value.Members[user.Id] = user
		user.ActiveGame = value
		value.Leaderboard[user.Id] = 0

		DataReply(false, JoinedGame, value).Send(conn)
		return
	}

	newGame := Game{
		Status:      Standby,
		InviteCode:  game.InviteCode,
		Members:     map[string]*User{},
		Leaderboard: map[string]float64{},
		Owner:       user,
		Data:        game,
	}

	newGame.Members[user.Id] = user
	newGame.Leaderboard[user.Id] = 0

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

	delete(game.Members, user.Id)
	user.ActiveGame = nil
	MessageReply(false, LeftGame).Send(conn)

	for _, member := range game.Members {
		DataReply(false, UserLeft, user).Send(member.Conn)
	}
}

func (g *gameSocketController) GetGame(ctx context.Context) {
	conn := ctx.Value("conn").(*websocket.Conn)

	user := ctx.Value("user").(*User)
	if user.ActiveGame == nil {
		MessageReply(true, NotInGame).Send(conn)
		return
	}

	DataReply(false, GetGame, user.ActiveGame).Send(conn)
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

func (g *gameSocketController) StartGame(ctx context.Context) {
	g.mu.Lock()
	defer g.mu.Unlock()

	conn := ctx.Value("conn").(*websocket.Conn)

	user := ctx.Value("user").(*User)
	if user.ActiveGame == nil {
		MessageReply(true, NotInGame).Send(conn)
		return
	}

	if user.ActiveGame.Owner != user {
		MessageReply(true, NotOwner).Send(conn)
		return
	}

	user.ActiveGame.Status = Starting

	n := 1
	for range time.Tick(time.Second * 1) {
		if n == 0 {
			for _, member := range user.ActiveGame.Members {
				MessageReply(false, InProgress).Send(member.Conn)
			}
			user.ActiveGame.Status = InProgress
			go g.PlayRound(user.ActiveGame)
			return
		}
		for _, member := range user.ActiveGame.Members {
			DataReply(false, Starting, n).Send(member.Conn)
		}

		n -= 1
	}
}

func (g *gameSocketController) ResetGame(ctx context.Context) {
	g.mu.Lock()
	defer g.mu.Unlock()

	conn := ctx.Value("conn").(*websocket.Conn)

	user := ctx.Value("user").(*User)
	if user.ActiveGame == nil {
		MessageReply(true, NotInGame).Send(conn)
		return
	}

	if user.ActiveGame.Owner != user {
		MessageReply(true, NotOwner).Send(conn)
		return
	}

	user.ActiveGame.Status = Standby

	DataReply(false, ResetGame, user.ActiveGame).Send(conn)
}

type RoundData struct {
	Timer    int       `json:"timer"`
	Question *Question `json:"question"`
}

type Option struct {
	Name string `json:"name"`
}

type Question struct {
	Name    string   `json:"name"`
	Options []Option `json:"options"`
}

func (g *gameSocketController) PlayRound(game *Game) {
	g.mu.Lock()
	defer g.mu.Unlock()

	game.RoundStatus = RoundInProgress
	for {
		if game.Status == InProgress && game.RoundStatus == RoundInProgress {
			n := game.Data.RoundTime
			question := Question{}
			copier.Copy(&question, game.Data.Questions[game.CurrentRound])

			for range time.Tick(time.Second * 1) {
				if n == 0 {
					for _, member := range game.Members {
						MessageReply(false, RoundWaiting).Send(member.Conn)
					}

					game.RoundStatus = RoundWaiting
					game.CurrentRound += 1

					break
				}
				for _, member := range game.Members {
					if game.Owner == member {
						DataReply(false, RoundInProgress, game.Data.Questions[game.CurrentRound]).Send(member.Conn)
					} else {
						DataReply(false, RoundInProgress, RoundData{Question: &question, Timer: n}).Send(member.Conn)
					}
				}

				n -= 1
			}
		} else {
			if game.CurrentRound >= len(game.Data.Questions) {
				break
			}
		}
	}
}

type AnswerData struct {
	Option uint `json:"option"`
}

func (g *gameSocketController) AnswerQuestion(ctx context.Context, data AnswerData) {
	g.mu.Lock()
	defer g.mu.Unlock()

	user := ctx.Value("user").(*User)
	conn := ctx.Value("conn").(*websocket.Conn)

	if user.ActiveGame == nil {
		MessageReply(true, NotInGame).Send(conn)
		return
	}

	if user.ActiveGame.RoundStatus == RoundInProgress {
		question := user.ActiveGame.Data.Questions[user.ActiveGame.CurrentRound]
		if question.Options[data.Option].Correct {
			user.ActiveGame.Leaderboard[user.Id] += user.ActiveGame.Data.Points
		}
		return
	} else {
		MessageReply(true, RoundWaiting).Send(conn)
	}
}

// {
// 	"message": "JOIN_GAME",
// 	"data": {
// 		"gameId": "9dbc-b4a1"
// 	}
// }
