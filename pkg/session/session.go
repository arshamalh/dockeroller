package session

type session map[int64]map[string]any

func New() *session {
	return &session{}
}

func (s *session) Set(chat_id int64, key string, value any) {
	(*s)[chat_id][key] = value
}

func (s *session) Get(chat_id int64, key string) any {
	return (*s)[chat_id][key]
}

func (s *session) Del(chat_id int64, key string) {
	(*s)[chat_id][key] = nil
}
