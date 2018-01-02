package chat

import (
	"context"
	"net/http"
	"time"

	gttp "github.com/benjamw/gogame/http"
	"github.com/benjamw/gogame/model"
	"github.com/benjamw/gogame/session"
)

func init() {
	gttp.R.Path("/room/{id:[0-9]+}").
		Methods("POST").
		Handler(&gttp.PlayerJSONHandler{handleAdd})

	gttp.R.Path("/room/{id:[0-9]+}").
		Methods("GET").
		Handler(&gttp.PlayerJSONHandler{handleRead})

	gttp.R.Path("/room/{id:[0-9]+}/after/{time:[0-9]+}").
		Methods("GET").
		Handler(&gttp.PlayerJSONHandler{handleLatest})

	gttp.R.Path("/muted").
		Methods("GET").
		Handler(&gttp.PlayerJSONHandler{handleMuted})

	gttp.R.Path("/mute").
		Methods("POST").
		Handler(&gttp.PlayerJSONHandler{handleMute})

	gttp.R.Path("/unmute").
		Methods("POST").
		Handler(&gttp.PlayerJSONHandler{handleUnmute})
}

type Reply struct {
	gttp.Response
	ChatID   string    `json:"chat_id"`
	RoomID   string    `json:"room_id"`
	PlayerID string    `json:"player_id"`
	Message  string    `json:"message"`
	Created  time.Time `json:"created"`
}

func (r *Reply) Set(m model.Chat) {
	r.ChatID = m.GetKey().Encode()
	r.RoomID = m.RoomKey.Encode()
	r.PlayerID = m.PlayerKey.Encode()
	r.Message = m.Message
	r.Created = m.Created
}

func handleAdd(ctx context.Context, s session.Data, w http.ResponseWriter, r *http.Request) (replyRaw interface{}, errReply error) {
	roomID := gttp.GetURLValue(r, "id")

	message := r.FormValue("message")
	if message == "" {
		errReply = &gttp.MissingRequiredError{FormElement: "message"}
		return
	}

	var chat model.Chat
	chat, errReply = AddChat(ctx, roomID, s.PlayerID, message)
	if errReply != nil {
		return
	}

	reply := Reply{
		Response: gttp.Response{
			Success: true,
		},
	}
	reply.Set(chat)

	replyRaw = reply

	return
}

type RoomReply struct {
	gttp.Response
	RoomID string  `json:"room_id"`
	Name   string  `json:"name"`
	Chats  []Reply `json:"chats"`
}

func (r *RoomReply) Set(m model.Room) {
	r.RoomID = m.GetKey().Encode()
	r.Name = m.Name
}

func (r *RoomReply) SetChats(l model.ChatList) {
	r.Chats = make([]Reply, 0)

	var reply Reply

	for k := range l {
		reply.Set(l[k])
		r.Chats = append(r.Chats, reply)
	}
}

func handleRead(ctx context.Context, s session.Data, w http.ResponseWriter, r *http.Request) (replyRaw interface{}, errReply error) {
	roomID := gttp.GetURLValue(r, "id")

	var room model.Room
	var chats model.ChatList
	room, chats, errReply = GetChats(ctx, roomID)
	if errReply != nil {
		return
	}

	reply := RoomReply{
		Response: gttp.Response{
			Success: true,
		},
	}
	reply.Set(room)
	reply.SetChats(chats)

	replyRaw = reply

	return
}

func handleLatest(ctx context.Context, s session.Data, w http.ResponseWriter, r *http.Request) (replyRaw interface{}, errReply error) {
	roomID := gttp.GetURLValue(r, "id")
	seen := gttp.GetURLValue(r, "time")

	var after time.Time
	after, errReply = time.Parse("20060102150405", seen)
	if errReply != nil {
		return
	}

	var room model.Room
	var chats model.ChatList
	room, chats, errReply = GetChatsAfter(ctx, roomID, after)
	if errReply != nil {
		return
	}

	reply := RoomReply{
		Response: gttp.Response{
			Success: true,
		},
	}
	reply.Set(room)
	reply.SetChats(chats)

	replyRaw = reply

	return
}

type MuteReply struct {
	gttp.Response
	MutedID string `json:"muted_id"`
}

func (r *MuteReply) Set(mute model.Mute) {
	r.MutedID = mute.MutedKey.Encode()
}

type MutedReply struct {
	Muted []MuteReply `json:"muted"`
}

func (r *MutedReply) Set(muted model.MuteList) {
	r.Muted = make([]MuteReply, len(muted))

	for k, v := range muted {
		r.Muted[k].Set(v)
	}
}

func handleMuted(ctx context.Context, s session.Data, w http.ResponseWriter, r *http.Request) (replyRaw interface{}, errReply error) {
	var mutes model.MuteList
	mutes, errReply = GetMuted(ctx, s.PlayerID)
	if errReply != nil {
		return
	}

	reply := MutedReply{}
	reply.Set(mutes)

	replyRaw = reply

	return
}

func handleMute(ctx context.Context, s session.Data, w http.ResponseWriter, r *http.Request) (replyRaw interface{}, errReply error) {
	mutedID := r.FormValue("muted_id")
	if mutedID == "" {
		errReply = &gttp.MissingRequiredError{FormElement: "muted_id"}
		return
	}

	var mute model.Mute
	mute, errReply = Mute(ctx, s.PlayerID, mutedID)
	if errReply != nil {
		return
	}

	reply := MuteReply{
		Response: gttp.Response{
			Success: true,
		},
	}
	reply.Set(mute)

	replyRaw = reply

	return
}

func handleUnmute(ctx context.Context, s session.Data, w http.ResponseWriter, r *http.Request) (replyRaw interface{}, errReply error) {
	mutedID := r.FormValue("muted_id")
	if mutedID == "" {
		errReply = &gttp.MissingRequiredError{FormElement: "muted_id"}
		return
	}

	errReply = Unmute(ctx, s.PlayerID, mutedID)
	if errReply != nil {
		return
	}

	reply := MuteReply{
		Response: gttp.Response{
			Success: true,
		},
	}

	replyRaw = reply

	return
}
