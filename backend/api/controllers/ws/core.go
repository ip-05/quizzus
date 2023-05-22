package ws

import (
	"context"
	"encoding/json"
	"go/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/app/game"
	"nhooyr.io/websocket"
)

type CoreController struct {
	gameController *GameSocketController
}

type SocketMessage struct {
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type SocketReply[D any] struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

const (
	JoinGame       = "JOIN_GAME"
	GetGame        = "GET_GAME"
	LeaveGame      = "LEAVE_GAME"
	IsOwner        = "IS_OWNER"
	StartGame      = "START_GAME"
	ResetGame      = "RESET_GAME"
	AnswerQuestion = "ANSWER_QUESTION"
	NextRound      = "NEXT_ROUND"
	Ping           = "PING"
	Pong           = "PONG"
)

func MessageReply(error bool, message string) SocketReply[types.Nil] {
	return SocketReply[types.Nil]{
		Error:   error,
		Message: message,
	}
}

func DataReply[D any](error bool, message string, data D) SocketReply[D] {
	return SocketReply[D]{
		Error:   error,
		Message: message,
		Data:    data,
	}
}

func (s SocketReply[D]) Send(conn *websocket.Conn) {
	bytes, _ := json.Marshal(s)

	conn.Write(context.Background(), websocket.MessageText, bytes)
}

func (w CoreController) messageHandler(ctx context.Context, conn *websocket.Conn) error {
	for {
		select {
		case <-ctx.Done():
			{
				return nil
			}
		default:
			{
				_, message, err := conn.Read(ctx)
				if err != nil {
					return err
				}

				var msg SocketMessage
				err = json.Unmarshal(message, &msg)
				if err != nil {
					DataReply(true, "MESSAGE_ERROR", err.Error()).Send(conn)
					break
				}

				switch msg.Message {
				case JoinGame:
					w.gameController.JoinGame(ctx, msg.Data)
				case LeaveGame:
					w.gameController.LeaveGame(ctx)
				case GetGame:
					w.gameController.GetGame(ctx)
				case IsOwner:
					w.gameController.IsOwner(ctx)
				case StartGame:
					w.gameController.StartGame(ctx)
				case ResetGame:
					w.gameController.ResetGame(ctx)
				case AnswerQuestion:
					w.gameController.AnswerQuestion(ctx, msg.Data)
				case NextRound:
					w.gameController.NextRound(ctx)
				case Ping:
					MessageReply(false, Pong).Send(conn)
				}
			}
		}
	}
}

func NewCoreController(game game.IService) *CoreController {
	return &CoreController{
		gameController: NewGameSocketController(game),
	}
}

func (w CoreController) HandleWS(c *gin.Context) {
	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		c.String(http.StatusBadRequest, "the sky is falling")
	}
	defer conn.Close(websocket.StatusInternalError, "")

	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}

	authedUser, _ := c.Get("authedUser")

	ctx := context.WithValue(context.Background(), "authedUser", authedUser)
	ctx = context.WithValue(ctx, "conn", conn)

	user, err := w.gameController.InitUser(ctx)
	if err != nil {
		DataReply(true, "INIT_ERROR", err.Error()).Send(conn)
		return
	}

	ctx = context.WithValue(ctx, "user", user)

	defer w.gameController.CleanUser(ctx)

	err = w.messageHandler(ctx, conn)
	if err != nil {
		DataReply(true, "HANDLER_ERROR", err.Error()).Send(conn)
		return
	}
}
