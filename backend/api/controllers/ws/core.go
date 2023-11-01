package ws

import (
	"context"
	"encoding/json"
	"go/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/utils"
	"nhooyr.io/websocket"
)

type SocketMessage struct {
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type SocketReply[D any] struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type CoreController struct {
	gameController *GameSocketController
}

func NewCoreController(game GameService, user UserService, session SessionService) *CoreController {
	return &CoreController{
		gameController: NewGameSocketController(game, user, session),
	}
}

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
					DataReply(true, utils.MessageError, err.Error()).Send(conn)
					break
				}

				switch msg.Message {
				case utils.JoinGame:
					w.gameController.JoinGame(ctx, msg.Data)
				case utils.LeaveGame:
					w.gameController.LeaveGame(ctx)
				case utils.GetGame:
					w.gameController.GetGame(ctx)
				case utils.IsOwner:
					w.gameController.IsOwner(ctx)
				case utils.StartGame:
					w.gameController.StartGame(ctx)
				case utils.ResetGame:
					w.gameController.ResetGame(ctx)
				case utils.AnswerQuestion:
					w.gameController.AnswerQuestion(ctx, msg.Data)
				case utils.SendChat:
					w.gameController.SendChat(ctx, msg.Data)
				case utils.NextRound:
					w.gameController.NextRound(ctx)
				case utils.Ping:
					MessageReply(false, utils.Pong).Send(conn)
				}
			}
		}
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
