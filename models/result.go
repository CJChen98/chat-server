package models

type JSON struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Token string `json:"token"`
	Data  Data   `json:"data"`
}
type Data struct {
	Messages      []Message      `json:"messages"`
	Users         []User         `json:"users"`
	Conversations []Conversation `json:"conversations"`
	Rooms         []Room         `json:"rooms"`
}
