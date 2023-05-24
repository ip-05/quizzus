package ws

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jinzhu/copier"

	"github.com/ip-05/quizzus/api/controllers/web"
	"github.com/ip-05/quizzus/api/middleware"
	"github.com/ip-05/quizzus/entity"
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
	RoundFinished   = "ROUND_FINISHED"
	AnswerAccepted  = "ANSWER_ACCEPTED"
)

const (
	GameNotFound  = "GAME_NOT_FOUND"
	AlreadyInGame = "ALREADY_IN_GAME"
	NotInGame     = "NOT_IN_GAME"
	NotOwner      = "NOT_OWNER"

	JoinedGame   = "JOINED_GAME"
	LeftGame     = "LEFT_GAME"
	GameDeleted  = "GAME_DELETED"
	UserLeft     = "USER_LEFT"
	UserJoined   = "USER_JOINED"
	UserAnswered = "USER_ANSWERED"
)

type User struct {
	Id             string          `json:"id"`
	Name           string          `json:"name"`
	ProfilePicture string          `json:"profilePicture"`
	ActiveGame     *Game           `json:"-"`
	Conn           *websocket.Conn `json:"-"`
}

type Game struct {
	Status        string             `json:"status"`
	RoundStatus   string             `json:"roundStatus"`
	CurrentRound  int                `json:"currentRound"`
	Points        float64            `json:"points"`
	Topic         string             `json:"topic"`
	RoundTime     int                `json:"roundTime"`
	QuestionCount int                `json:"questionCount"`
	InviteCode    string             `json:"inviteCode"`
	Members       map[string]*User   `json:"members"`
	Owner         *User              `json:"owner"`
	Leaderboard   map[string]float64 `json:"leaderboard"`
	Data          *entity.Game       `json:"-"`
	Rounds        map[int]*Round     `json:"-"`
}

type Round struct {
	Answers map[string]uint
}

type GameSocketController struct {
	Users    map[string]*User
	Games    map[string]*Game
	Game     web.IGameService
	GameTime int
}

func NewGameSocketController(game web.IGameService) *GameSocketController {
	c := new(GameSocketController)

	c.Users = make(map[string]*User)
	c.Games = make(map[string]*Game)
	c.Game = game
	c.GameTime = 10

	return c
}

func (g *GameSocketController) InitUser(ctx context.Context) (*User, error) {
	user := ctx.Value("authedUser").(middleware.AuthedUser)

	_, found := g.Users[user.Id]
	if found {
		return nil, errors.New("user already exists on another socket")
	}

	g.Users[user.Id] = &User{
		Id:             user.Id,
		Name:           user.Name,
		ProfilePicture: user.ProfilePicture,
		Conn:           ctx.Value("conn").(*websocket.Conn),
	}

	return g.Users[user.Id], nil
}

func (g *GameSocketController) CleanUser(ctx context.Context) {
	user := ctx.Value("user").(*User)
	if user.ActiveGame != nil {
		g.LeaveGame(ctx)
	}

	delete(g.Users, user.Id)
}

type JoinGameData struct {
	GameId string `json:"gameId"`
}

func (g *GameSocketController) JoinGame(ctx context.Context, msgData json.RawMessage) {
	conn := ctx.Value("conn").(*websocket.Conn)

	var data JoinGameData
	err := json.Unmarshal(msgData, &data)
	if err != nil {
		DataReply(true, "DATA_ERROR", err.Error()).Send(conn)
		return
	}

	user := ctx.Value("user").(*User)
	if user.ActiveGame != nil {
		MessageReply(true, AlreadyInGame).Send(conn)
		return
	}

	// var game entity.Game
	// g.DB.Preload("Questions.Options").Where("invite_code = ?", data.GameId).First(&game)
	game, err := g.Game.GetGame(0, data.GameId)
	if err != nil {
		DataReply(true, "DATA_ERROR", err.Error()).Send(conn)
		return
	}

	if game.Id == 0 {
		MessageReply(true, GameNotFound).Send(conn)
		return
	}

	value, ok := g.Games[game.InviteCode]
	if ok {
		for _, member := range value.Members {
			DataReply(false, UserJoined, user).Send(member.Conn)
		}

		value.Members[user.Id] = user
		user.ActiveGame = value
		value.Leaderboard[user.Id] = 0

		DataReply(false, JoinedGame, value).Send(conn)
		return
	}

	if game.Owner != user.Id {
		MessageReply(true, NotOwner).Send(conn)
		return
	}

	newGame := Game{
		Status:        Standby,
		RoundStatus:   RoundWaiting,
		Points:        game.Points,
		Topic:         game.Topic,
		QuestionCount: len(game.Questions),
		RoundTime:     game.RoundTime,
		InviteCode:    game.InviteCode,
		Members:       map[string]*User{},
		Leaderboard:   map[string]float64{},
		Rounds:        map[int]*Round{},
		Owner:         user,
		Data:          game,
	}

	newGame.Members[user.Id] = user
	newGame.Leaderboard[user.Id] = 0

	g.Games[newGame.InviteCode] = &newGame
	user.ActiveGame = g.Games[newGame.InviteCode]
	DataReply(false, JoinedGame, newGame).Send(conn)
}

func (g *GameSocketController) LeaveGame(ctx context.Context) {
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

func (g *GameSocketController) GetGame(ctx context.Context) {
	conn := ctx.Value("conn").(*websocket.Conn)

	user := ctx.Value("user").(*User)
	if user.ActiveGame == nil {
		MessageReply(true, NotInGame).Send(conn)
		return
	}

	DataReply(false, GetGame, user.ActiveGame).Send(conn)
}

func (g *GameSocketController) IsOwner(ctx context.Context) {
	conn := ctx.Value("conn").(*websocket.Conn)

	user := ctx.Value("user").(*User)
	if user.ActiveGame == nil {
		MessageReply(true, NotInGame).Send(conn)
		return
	}

	DataReply(false, IsOwner, user.ActiveGame.Owner == user).Send(conn)
}

func (g *GameSocketController) StartGame(ctx context.Context) {
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

	if user.ActiveGame.Status != Standby {
		MessageReply(true, user.ActiveGame.Status).Send(conn)
		return
	}

	user.ActiveGame.Status = Starting

	n := g.GameTime
	for range time.Tick(time.Second * 1) {
		if n == 0 {
			break
		}
		for _, member := range user.ActiveGame.Members {
			DataReply(false, Starting, n).Send(member.Conn)
		}

		n -= 1
	}
	for _, member := range user.ActiveGame.Members {
		MessageReply(false, InProgress).Send(member.Conn)
	}
	user.ActiveGame.Status = InProgress
	go g.PlayRounds(user.ActiveGame)

}

func (g *GameSocketController) ResetGame(ctx context.Context) {
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

	if user.ActiveGame.Status != Finished {
		MessageReply(true, user.ActiveGame.Status).Send(conn)
		return
	}

	user.ActiveGame.Status = Standby
	user.ActiveGame.RoundStatus = RoundWaiting
	user.ActiveGame.CurrentRound = 0
	user.ActiveGame.Leaderboard = map[string]float64{}
	user.ActiveGame.Rounds = map[int]*Round{}

	DataReply(false, ResetGame, user.ActiveGame).Send(conn)
}

type RoundData[T entity.Question | Question] struct {
	Timer    int `json:"timer"`
	Question *T  `json:"question"`
}

type Option struct {
	Name string `json:"name"`
}

type Question struct {
	Name    string   `json:"name"`
	Options []Option `json:"options"`
}

type FinishedReply struct {
	Correct     bool               `json:"correct"`
	Options     []*entity.Option   `json:"options"`
	Leaderboard map[string]float64 `json:"leaderboard"`
}

func (g *GameSocketController) PlayRounds(game *Game) {
	game.RoundStatus = RoundInProgress
	for {
		if game.Status == InProgress && game.RoundStatus == RoundInProgress {
			n := game.Data.RoundTime
			question := Question{}
			copier.Copy(&question, game.Data.Questions[game.CurrentRound])
			game.Rounds[game.CurrentRound] = &Round{Answers: map[string]uint{}}

			for range time.Tick(time.Second * 1) {
				if n == 0 {
					for _, member := range game.Members {
						choice := game.Rounds[game.CurrentRound].Answers[member.Id]
						correct := game.Data.Questions[game.CurrentRound].Options[choice].Correct

						if correct {
							game.Leaderboard[member.Id] += game.Data.Points
						}
					}

					for _, member := range game.Members {
						choice := game.Rounds[game.CurrentRound].Answers[member.Id]
						correct := game.Data.Questions[game.CurrentRound].Options[choice].Correct

						DataReply(false, RoundFinished, FinishedReply{Correct: correct, Options: game.Data.Questions[game.CurrentRound].Options, Leaderboard: game.Leaderboard}).Send(member.Conn)
					}

					game.RoundStatus = RoundWaiting
					game.CurrentRound += 1

					break
				}
				for _, member := range game.Members {
					if game.Owner == member {
						DataReply(false, RoundInProgress, RoundData[entity.Question]{Question: game.Data.Questions[game.CurrentRound], Timer: n}).Send(member.Conn)
					} else {
						DataReply(false, RoundInProgress, RoundData[Question]{Question: &question, Timer: n}).Send(member.Conn)
					}
				}

				n -= 1
			}
		} else {
			if game.CurrentRound >= len(game.Data.Questions) {
				game.Status = Finished

				for _, member := range game.Members {
					DataReply(false, Finished, game.Leaderboard).Send(member.Conn)
				}

				break
			}
		}
	}
}

type AnswerData struct {
	Option uint `json:"option"`
}

type AnswerResponse struct {
	UserId string `json:"user"`
	Option uint   `json:"option"`
}

func (g *GameSocketController) AnswerQuestion(ctx context.Context, msgData json.RawMessage) {
	user := ctx.Value("user").(*User)
	conn := ctx.Value("conn").(*websocket.Conn)

	var data AnswerData
	err := json.Unmarshal(msgData, &data)
	if err != nil {
		DataReply(true, "DATA_ERROR", err.Error()).Send(conn)
		return
	}

	if user.ActiveGame == nil {
		MessageReply(true, NotInGame).Send(conn)
		return
	}

	if user.ActiveGame.RoundStatus == RoundInProgress && user.ActiveGame.Status == InProgress {
		user.ActiveGame.Rounds[user.ActiveGame.CurrentRound].Answers[user.Id] = data.Option
		DataReply(false, AnswerAccepted, user.ActiveGame).Send(conn)
		DataReply(false, UserAnswered, AnswerResponse{
			UserId: user.Id,
			Option: data.Option,
		}).Send(user.ActiveGame.Owner.Conn)
	} else {
		MessageReply(true, RoundWaiting).Send(conn)
	}
}

func (g *GameSocketController) NextRound(ctx context.Context) {
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

	if user.ActiveGame.Status != InProgress {
		MessageReply(true, user.ActiveGame.Status).Send(conn)
		return
	}

	if user.ActiveGame.RoundStatus != RoundWaiting {
		MessageReply(true, InProgress).Send(conn)
		return
	}

	user.ActiveGame.RoundStatus = RoundInProgress

	MessageReply(false, RoundInProgress).Send(conn)
}
