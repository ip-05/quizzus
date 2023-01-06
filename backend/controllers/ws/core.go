package ws

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go/types"
	"log"
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

func (s SocketReply[D]) Send(ctx context.Context) {
	conn := ctx.Value("conn").(*websocket.Conn)

	bytes, _ := json.Marshal(s)

	conn.Write(ctx, websocket.MessageText, bytes)
}

//
//type SocketMessage struct {
//	Message string            `json:"message"`
//	Args    map[string]string `json:"args"`
//}
//
//type User struct {
//	Id           string          `json:"-"`
//	Name         string          `json:"name"`
//	CurrentLobby *Lobby          `json:"-"`
//	Conn         *websocket.Conn `json:"-"`
//}
//
//type LobbyStatus = string
//
//const (
//	STANDBY     = "STANDBY"
//	STARTING    = "STARTING"
//	IN_PROGRESS = "IN_PROGRESS"
//)
//
//type Lobby struct {
//	Status  LobbyStatus `json:"status"`
//	Owner   *User       `json:"owner"`
//	Members []*User     `json:"members"`
//	Game    models.Game `json:"game"`
//}
//
//func NewWSController() *WebSocketController {
//	ws := new(WebSocketController)
//	ws.Lobbies = make(map[string]*Lobby)
//	return ws
//}
//
//func (user *User) ReadMessages(ctx context.Context, w *WebSocketController) error {
//	for {
//		select {
//		case <-ctx.Done():
//			{
//				return nil
//			}
//		default:
//			{
//				_, message, err := user.Conn.Read(ctx)
//				if err != nil {
//					return err
//				}
//
//				var msg SocketMessage
//				err = json.Unmarshal(message, &msg)
//				if err != nil {
//					return err
//				}
//
//				if msg.Message == "JOIN_GAME" {
//					if user.CurrentLobby != nil {
//						user.Conn.Write(ctx, websocket.MessageText, []byte("ALREADY_IN_GAME"))
//						break
//					}
//
//					var game models.Game
//
//					models.DB.Preload("Questions.Options").Where("invite_code = ?", msg.Args["game_id"]).First(&game)
//
//					if game.Id == 0 {
//						user.Conn.Write(ctx, websocket.MessageText, []byte("GAME_NOT_FOUND"))
//						break
//					}
//
//					value, ok := w.Lobbies[game.InviteCode]
//					if ok {
//						value.Members = append(value.Members, user)
//						user.CurrentLobby = value
//						lobbyBytes, _ := json.Marshal(value)
//
//						user.Conn.Write(ctx, websocket.MessageText, lobbyBytes)
//						break
//					}
//
//					lobby := Lobby{
//						Status:  STANDBY,
//						Owner:   user,
//						Members: []*User{user},
//						Game:    game,
//					}
//
//					w.Lobbies[game.InviteCode] = &lobby
//					user.CurrentLobby = &lobby
//
//					lobbyBytes, _ := json.Marshal(lobby)
//
//					user.Conn.Write(ctx, websocket.MessageText, lobbyBytes)
//				} else if msg.Message == "GET_GAME" {
//					if user.CurrentLobby == nil {
//						user.Conn.Write(ctx, websocket.MessageText, []byte("NOT_IN_GAME"))
//						break
//					}
//
//					lobbyBytes, _ := json.Marshal(user.CurrentLobby)
//
//					user.Conn.Write(ctx, websocket.MessageText, lobbyBytes)
//				} else if msg.Message == "GAME_IS_OWNER" {
//					if user.CurrentLobby == nil {
//						user.Conn.Write(ctx, websocket.MessageText, []byte("NOT_IN_GAME"))
//						break
//					}
//
//					if user.CurrentLobby.Owner == user {
//						user.Conn.Write(ctx, websocket.MessageText, []byte("true"))
//					} else {
//						user.Conn.Write(ctx, websocket.MessageText, []byte("false"))
//					}
//				} else if msg.Message == "DELETE_GAME" {
//					if user.CurrentLobby == nil {
//						user.Conn.Write(ctx, websocket.MessageText, []byte("NOT_IN_GAME"))
//						break
//					}
//
//					if user.CurrentLobby.Owner != user {
//						user.Conn.Write(ctx, websocket.MessageText, []byte("NOT_OWNER"))
//						break
//					}
//
//					delete(w.Lobbies, user.CurrentLobby.Game.InviteCode)
//
//					for _, member := range user.CurrentLobby.Members {
//						member.Conn.Write(ctx, websocket.MessageText, []byte("GAME_DELETED"))
//						member.CurrentLobby = nil
//					}
//				} else if msg.Message == "START_GAME" {
//					if user.CurrentLobby == nil {
//						user.Conn.Write(ctx, websocket.MessageText, []byte("NOT_IN_GAME"))
//						break
//					}
//
//					if user.CurrentLobby.Owner != user {
//						user.Conn.Write(ctx, websocket.MessageText, []byte("NOT_OWNER"))
//						break
//					}
//
//					user.CurrentLobby.Status = STARTING
//
//					n := 10
//					for range time.Tick(time.Second) {
//						for _, member := range user.CurrentLobby.Members {
//							if n == 0 {
//								member.Conn.Write(ctx, websocket.MessageText, []byte("GAME_STARTED"))
//							} else {
//								member.Conn.Write(ctx, websocket.MessageText, []byte("GAME_STARTING"))
//							}
//						}
//
//						if n == 0 {
//							user.CurrentLobby.Status = IN_PROGRESS
//							break
//						}
//						n -= 1
//					}
//				}
//			}
//		}
//	}
//}

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
					DataReply(true, "MESSAGE_ERROR", err.Error()).Send(ctx)
					break
				}

				if msg.Message == "JOIN_GAME" {
					var data JoinGameData
					err = json.Unmarshal(msg.Data, &data)
					if err != nil {
						DataReply(true, "DATA_ERROR", err.Error()).Send(ctx)
						break
					}

					gameController.JoinGame(ctx, data)
				}
			}
		}
	}
}

func (w CoreController) HandleWS(c *gin.Context) {
	conn, err := websocket.Accept(c.Writer, c.Request, nil)
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
		DataReply(true, "INIT_ERROR", err.Error()).Send(ctx)
		return
	}

	ctx = context.WithValue(ctx, "user", user)

	err = messageHandler(ctx, conn)
	if err != nil {
		log.Printf("error: %s\n", err)
		DataReply(true, "HANDLER_ERROR", err.Error()).Send(ctx)
		return
	}
}
