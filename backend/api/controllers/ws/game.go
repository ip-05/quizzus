package ws

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/jinzhu/copier"

	"github.com/ip-05/quizzus/api/middleware"
	"github.com/ip-05/quizzus/entity"
	"github.com/ip-05/quizzus/utils"
	"nhooyr.io/websocket"
)

type User struct {
	ID             uint            `json:"id"`
	Name           string          `json:"name"`
	ProfilePicture string          `json:"profile_picture"`
	ActiveGame     *Game           `json:"-"`
	Conn           *websocket.Conn `json:"-"`
}

type Game struct {
	ID            uint             `json:"id"`
	InstID        int              `json:"-"`
	Status        string           `json:"status"`
	RoundStatus   string           `json:"round_status"`
	CurrentRound  int              `json:"current_round"`
	Points        float64          `json:"points"`
	Topic         string           `json:"topic"`
	RoundTime     int              `json:"round_time"`
	QuestionCount int              `json:"question_count"`
	InviteCode    string           `json:"invite_code"`
	Members       map[uint]*User   `json:"members"`
	Owner         *User            `json:"owner"`
	Leaderboard   map[uint]float64 `json:"leaderboard"`
	Data          *entity.Game     `json:"-"`
	Rounds        map[int]*Round   `json:"-"`
}

type Round struct {
	Answers map[uint]uint
}

type GameSocketController struct {
	Users    map[uint]*User
	Games    map[string]*Game
	Game     GameService
	User     UserService
	Session  SessionService
	GameTime int
}

type GameService interface {
	CreateGame(body entity.CreateGame, ownerID uint) (*entity.Game, error)
	UpdateGame(body entity.UpdateGame, ID int, code string, ownerID uint) (*entity.Game, error)
	DeleteGame(ID int, code string, userID uint) error

	GetGame(ID int, code string) (*entity.Game, error)
	GetGamesByOwner(ID int, user int, limit int) (*[]entity.Game, error)
	GetFavoriteGames(user int) (*[]entity.Game, error)

	Favorite(ID int, userID int) bool
}

type SessionService interface {
	NewSession(ID, userID, instID int) uint
	EndSession(ID, userID, instID, questions, players int, points float64) uint
}

type UserService interface {
	CreateUser(body *entity.CreateUser) (*entity.User, error)
	UpdateUser(ID uint, body entity.UpdateUser) (*entity.User, error)
	DeleteUser(ID uint)
	GetUserById(ID uint) *entity.User
	GetUserByProvider(ID string, provider string) *entity.User
}

func NewGameSocketController(gameSvc GameService, userSvc UserService, sessionSvc SessionService) *GameSocketController {
	c := new(GameSocketController)

	c.Users = make(map[uint]*User)
	c.Games = make(map[string]*Game)
	c.Game = gameSvc
	c.User = userSvc
	c.Session = sessionSvc
	c.GameTime = 10

	return c
}

func (c *GameSocketController) InitUser(ctx context.Context) (*User, error) {
	authedUser := ctx.Value("authedUser").(middleware.AuthedUser)

	user := c.User.GetUserById(authedUser.ID)
	if user == nil {
		return nil, errors.New("no user found")
	}

	_, found := c.Users[user.ID]
	if found {
		return nil, errors.New("user already exists on another socket")
	}

	c.Users[user.ID] = &User{
		ID:             user.ID,
		Name:           user.Name,
		ProfilePicture: user.Picture,
		Conn:           ctx.Value("conn").(*websocket.Conn),
	}

	return c.Users[user.ID], nil
}

func (c *GameSocketController) CleanUser(ctx context.Context) {
	user := ctx.Value("user").(*User)
	if user.ActiveGame != nil {
		c.LeaveGame(ctx)
	}

	delete(c.Users, user.ID)
}

type JoinGameData struct {
	GameID string `json:"game_id"`
}

func (c *GameSocketController) JoinGame(ctx context.Context, msgData json.RawMessage) {
	conn := ctx.Value("conn").(*websocket.Conn)

	var data JoinGameData
	err := json.Unmarshal(msgData, &data)
	if err != nil {
		DataReply(true, "DATA_ERROR", err.Error()).Send(conn)
		return
	}

	user := ctx.Value("user").(*User)
	if user.ActiveGame != nil {
		MessageReply(true, utils.AlreadyInGame).Send(conn)
		return
	}

	// var game entity.Game
	// g.DB.Preload("Questions.Options").Where("invite_code = ?", data.GameId).First(&game)
	game, err := c.Game.GetGame(0, data.GameID)
	if err != nil {
		DataReply(true, "DATA_ERROR", err.Error()).Send(conn)
		return
	}

	if game.ID == 0 {
		MessageReply(true, utils.GameNotFound).Send(conn)
		return
	}

	value, ok := c.Games[game.InviteCode]
	if ok {
		for _, member := range value.Members {
			DataReply(false, utils.UserJoined, user).Send(member.Conn)
		}

		value.Members[user.ID] = user
		user.ActiveGame = value
		value.Leaderboard[user.ID] = 0

		DataReply(false, utils.JoinedGame, value).Send(conn)
		return
	}

	if game.Owner != user.ID {
		MessageReply(true, utils.NotOwner).Send(conn)
		return
	}

	instID, _ := rand.Int(rand.Reader, big.NewInt(100000000000))

	newGame := Game{
		ID:            game.ID,
		InstID:        int(instID.Int64()),
		Status:        utils.Standby,
		RoundStatus:   utils.RoundWaiting,
		Points:        game.Points,
		Topic:         game.Topic,
		QuestionCount: len(game.Questions),
		RoundTime:     game.RoundTime,
		InviteCode:    game.InviteCode,
		Members:       map[uint]*User{},
		Leaderboard:   map[uint]float64{},
		Rounds:        map[int]*Round{},
		Owner:         user,
		Data:          game,
	}

	newGame.Members[user.ID] = user
	newGame.Leaderboard[user.ID] = 0

	c.Games[newGame.InviteCode] = &newGame
	user.ActiveGame = c.Games[newGame.InviteCode]
	DataReply(false, utils.JoinedGame, newGame).Send(conn)
}

type ChatData struct {
	Message string `json:"message"`
}

type ChatBroadcast struct {
	Name    string `json:"name"`
	UserID  uint   `json:"user_id"`
	Message string `json:"message"`
}

func (c *GameSocketController) SendChat(ctx context.Context, msgData json.RawMessage) {
	conn := ctx.Value("conn").(*websocket.Conn)
	user := ctx.Value("user").(*User)
	if user.ActiveGame == nil {
		MessageReply(true, utils.NotInGame).Send(conn)
		return
	}

	var data ChatData
	err := json.Unmarshal(msgData, &data)
	if err != nil {
		DataReply(true, "DATA_ERROR", err.Error()).Send(conn)
		return
	}

	game := user.ActiveGame
	for _, member := range game.Members {
		DataReply(false, utils.ReceiveChat, ChatBroadcast{Name: user.Name, Message: data.Message, UserID: user.ID}).Send(member.Conn)
	}
}

func (c *GameSocketController) LeaveGame(ctx context.Context) {
	conn := ctx.Value("conn").(*websocket.Conn)
	user := ctx.Value("user").(*User)
	if user.ActiveGame == nil {
		MessageReply(true, utils.NotInGame).Send(conn)
		return
	}

	game := user.ActiveGame
	if game.Owner == user {
		for _, member := range game.Members {
			member.ActiveGame = nil
			MessageReply(false, utils.GameDeleted).Send(member.Conn)
		}
		delete(c.Games, game.InviteCode)
	}

	delete(game.Members, user.ID)
	user.ActiveGame = nil
	MessageReply(false, utils.LeftGame).Send(conn)

	for _, member := range game.Members {
		DataReply(false, utils.UserLeft, user).Send(member.Conn)
	}
}

func (c *GameSocketController) GetGame(ctx context.Context) {
	conn := ctx.Value("conn").(*websocket.Conn)

	user := ctx.Value("user").(*User)
	if user.ActiveGame == nil {
		MessageReply(true, utils.NotInGame).Send(conn)
		return
	}

	DataReply(false, utils.GetGame, user.ActiveGame).Send(conn)
}

func (c *GameSocketController) IsOwner(ctx context.Context) {
	conn := ctx.Value("conn").(*websocket.Conn)

	user := ctx.Value("user").(*User)
	if user.ActiveGame == nil {
		MessageReply(true, utils.NotInGame).Send(conn)
		return
	}

	DataReply(false, utils.IsOwner, user.ActiveGame.Owner == user).Send(conn)
}

func (c *GameSocketController) StartGame(ctx context.Context) {
	conn := ctx.Value("conn").(*websocket.Conn)

	user := ctx.Value("user").(*User)
	if user.ActiveGame == nil {
		MessageReply(true, utils.NotInGame).Send(conn)
		return
	}

	if user.ActiveGame.Owner != user {
		MessageReply(true, utils.NotOwner).Send(conn)
		return
	}

	if user.ActiveGame.Status != utils.Standby {
		MessageReply(true, user.ActiveGame.Status).Send(conn)
		return
	}

	user.ActiveGame.Status = utils.Starting

	n := c.GameTime
	for range time.Tick(time.Second * 1) {
		if n == 0 {
			break
		}
		for _, member := range user.ActiveGame.Members {
			DataReply(false, utils.Starting, n).Send(member.Conn)
		}

		n -= 1
	}
	for id, member := range user.ActiveGame.Members {
		MessageReply(false, utils.InProgress).Send(member.Conn)

		c.Session.NewSession(int(user.ActiveGame.ID), int(id), user.ActiveGame.InstID)
	}
	user.ActiveGame.Status = utils.InProgress
	go c.PlayRounds(user.ActiveGame)

}

func (c *GameSocketController) ResetGame(ctx context.Context) {
	conn := ctx.Value("conn").(*websocket.Conn)

	user := ctx.Value("user").(*User)
	if user.ActiveGame == nil {
		MessageReply(true, utils.NotInGame).Send(conn)
		return
	}

	if user.ActiveGame.Owner != user {
		MessageReply(true, utils.NotOwner).Send(conn)
		return
	}

	if user.ActiveGame.Status != utils.Finished {
		MessageReply(true, user.ActiveGame.Status).Send(conn)
		return
	}

	user.ActiveGame.Status = utils.Standby
	user.ActiveGame.RoundStatus = utils.RoundWaiting
	user.ActiveGame.CurrentRound = 0
	user.ActiveGame.Leaderboard = map[uint]float64{}
	user.ActiveGame.Rounds = map[int]*Round{}

	DataReply(false, utils.ResetGame, user.ActiveGame).Send(conn)
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
	Correct     bool             `json:"correct"`
	Options     []*entity.Option `json:"options"`
	Leaderboard map[uint]float64 `json:"leaderboard"`
}

func (c *GameSocketController) PlayRounds(game *Game) {
	game.RoundStatus = utils.RoundInProgress
	for {
		if game.Status == utils.InProgress && game.RoundStatus == utils.RoundInProgress {
			n := game.Data.RoundTime
			question := Question{}
			copier.Copy(&question, game.Data.Questions[game.CurrentRound])
			game.Rounds[game.CurrentRound] = &Round{Answers: map[uint]uint{}}

			for range time.Tick(time.Second * 1) {
				if n == 0 {
					for _, member := range game.Members {
						choice := game.Rounds[game.CurrentRound].Answers[member.ID]
						correct := game.Data.Questions[game.CurrentRound].Options[choice].Correct

						if correct {
							game.Leaderboard[member.ID] += game.Data.Points
						}
					}

					for _, member := range game.Members {
						choice := game.Rounds[game.CurrentRound].Answers[member.ID]
						correct := game.Data.Questions[game.CurrentRound].Options[choice].Correct

						DataReply(false, utils.RoundFinished, FinishedReply{Correct: correct, Options: game.Data.Questions[game.CurrentRound].Options, Leaderboard: game.Leaderboard}).Send(member.Conn)
					}

					game.RoundStatus = utils.RoundWaiting
					game.CurrentRound += 1

					break
				}
				for _, member := range game.Members {
					if game.Owner == member {
						DataReply(false, utils.RoundInProgress, RoundData[entity.Question]{Question: game.Data.Questions[game.CurrentRound], Timer: n}).Send(member.Conn)
					} else {
						DataReply(false, utils.RoundInProgress, RoundData[Question]{Question: &question, Timer: n}).Send(member.Conn)
					}
				}

				n -= 1
			}
		} else {
			if game.CurrentRound >= len(game.Data.Questions) {
				game.Status = utils.Finished

				for id, member := range game.Members {
					DataReply(false, utils.Finished, game.Leaderboard).Send(member.Conn)

					c.Session.EndSession(int(game.ID), int(id), game.InstID, game.QuestionCount, len(game.Members)-1, game.Leaderboard[id])
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
	UserID uint `json:"user"`
	Option uint `json:"option"`
}

func (c *GameSocketController) AnswerQuestion(ctx context.Context, msgData json.RawMessage) {
	user := ctx.Value("user").(*User)
	conn := ctx.Value("conn").(*websocket.Conn)

	var data AnswerData
	err := json.Unmarshal(msgData, &data)
	if err != nil {
		DataReply(true, "DATA_ERROR", err.Error()).Send(conn)
		return
	}

	if user.ActiveGame == nil {
		MessageReply(true, utils.NotInGame).Send(conn)
		return
	}

	if user.ActiveGame.RoundStatus == utils.RoundInProgress && user.ActiveGame.Status == utils.InProgress {
		user.ActiveGame.Rounds[user.ActiveGame.CurrentRound].Answers[user.ID] = data.Option
		DataReply(false, utils.AnswerAccepted, user.ActiveGame).Send(conn)
		DataReply(false, utils.UserAnswered, AnswerResponse{
			UserID: user.ID,
			Option: data.Option,
		}).Send(user.ActiveGame.Owner.Conn)
	} else {
		MessageReply(true, utils.RoundWaiting).Send(conn)
	}
}

func (c *GameSocketController) NextRound(ctx context.Context) {
	conn := ctx.Value("conn").(*websocket.Conn)

	user := ctx.Value("user").(*User)
	if user.ActiveGame == nil {
		MessageReply(true, utils.NotInGame).Send(conn)
		return
	}

	if user.ActiveGame.Owner != user {
		MessageReply(true, utils.NotOwner).Send(conn)
		return
	}

	if user.ActiveGame.Status != utils.InProgress {
		MessageReply(true, user.ActiveGame.Status).Send(conn)
		return
	}

	if user.ActiveGame.RoundStatus != utils.RoundWaiting {
		MessageReply(true, utils.InProgress).Send(conn)
		return
	}

	user.ActiveGame.RoundStatus = utils.RoundInProgress

	MessageReply(false, utils.RoundInProgress).Send(conn)
}
