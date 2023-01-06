package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"time"
)

type WebSocketController struct{}

type SocketMessage struct {
	Message string            `json:"message"`
	Args    map[string]string `json:"args"`
}

type User struct {
	Id   string
	Name string
	Conn *websocket.Conn
}

func (user *User) ReadMessages(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			{
				return nil
			}
		default:
			{
				_, message, err := user.Conn.Read(ctx)
				if err != nil {
					return err
				}

				var msg SocketMessage
				err = json.Unmarshal(message, &msg)
				if err != nil {
					return err
				}

				if msg.Message == "JOIN_GAME" {
					resp, err := http.Get(fmt.Sprintf("http://localhost:3000/games?invite_code=%s", msg.Args["game_id"]))
					if err != nil {
						user.Conn.Write(ctx, websocket.MessageText, []byte("Game does not exist"))
					}

					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Fatalln(err)
					}

					user.Conn.Write(ctx, websocket.MessageText, body)
				}
			}
		}
	}
}

func (w WebSocketController) UpgradeWS(c *gin.Context) {
	conn, err := websocket.Accept(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusBadRequest, "the sky is falling")
	}
	defer conn.Close(websocket.StatusInternalError, "")

	authedUser, _ := c.Get("authedUser")

	user := User{
		Id:   authedUser.(middleware.AuthedUser).Id,
		Name: authedUser.(middleware.AuthedUser).Name,
		Conn: conn,
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(300*time.Second))

	err = user.ReadMessages(ctx)
	if errors.Is(err, context.Canceled) {
		return
	}

	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}

	if err != nil {
		log.Printf("error: %s\n", err)
		return
	}
}
