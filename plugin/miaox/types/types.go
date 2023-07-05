package types

type JsonRest struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

type Page struct {
	Total int   `json:"total"`
	List  []any `json:"list"`
}

type ModelService interface {
	NewModel() interface{}

	Find(model interface{}) Page
	Get(key string) interface{}
	Edit(model interface{}) bool
	Del(key string) bool
}

// ====

type ConversationContextArgs struct {
	Current  string
	Nickname string
}
