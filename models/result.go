package models

type JSON struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}
type Data struct {
	User     User      `json:"user"`
	Messages []Message `json:"messages"`
}
