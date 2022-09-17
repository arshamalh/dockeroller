package contracts

type Session interface {
	Set(chat_id int64, key string, value any)
	Get(chat_id int64, key string) any
	Del(chat_id int64, key string)
}
