package ws

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go/types"
	"net/http"
	"nhooyr.io/websocket"
	"time"
)

type CoreController struct{}

var gameController = NewGameSocketController()

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
)

func MessageReply[D types.Nil](error bool, message string) SocketReply[D] {
	return SocketReply[D]{
		Error:   error,
		Message: message,
		Data:    nil,
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

func messageHandler(ctx context.Context, conn *websocket.Conn) error {
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

				if msg.Message == JoinGame {
					var data JoinGameData
					err = json.Unmarshal(msg.Data, &data)
					if err != nil {
						DataReply(true, "DATA_ERROR", err.Error()).Send(conn)
						break
					}

					gameController.JoinGame(ctx, data)
				} else if msg.Message == LeaveGame {
					gameController.LeaveGame(ctx)
				} else if msg.Message == GetGame {
					gameController.GetGame(ctx)
				} else if msg.Message == IsOwner {
					gameController.IsOwner(ctx)
				} else if msg.Message == StartGame {
					gameController.StartGame(ctx)
				} else if msg.Message == ResetGame {
					gameController.ResetGame(ctx)
				} else if msg.Message == AnswerQuestion {
					var data AnswerData
					err = json.Unmarshal(msg.Data, &data)
					if err != nil {
						DataReply(true, "DATA_ERROR", err.Error()).Send(conn)
						break
					}

					gameController.AnswerQuestion(ctx, data)
				} else if msg.Message == NextRound {
					gameController.NextRound(ctx)
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

	ctx, _ := context.WithTimeout(context.Background(), 300*time.Second)
	ctx = context.WithValue(ctx, "authedUser", authedUser)
	ctx = context.WithValue(ctx, "conn", conn)

	user, err := gameController.InitUser(ctx)
	if err != nil {
		DataReply(true, "INIT_ERROR", err.Error()).Send(conn)
		return
	}

	ctx = context.WithValue(ctx, "user", user)

	defer gameController.CleanUser(ctx)

	err = messageHandler(ctx, conn)
	if err != nil {
		DataReply(true, "HANDLER_ERROR", err.Error()).Send(conn)
		return
	}
}
