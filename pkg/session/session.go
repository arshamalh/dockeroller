package session

type session map[int64]map[string]any

func New() *session {
	ssn := make(session)
	return &ssn
}

func (s *session) Set(chat_id int64, key string, value any) {
	if (*s)[chat_id] == nil {
		(*s)[chat_id] = make(map[string]any)
	}
	(*s)[chat_id][key] = value
}

func (s *session) Get(chat_id int64, key string) any {
	return (*s)[chat_id][key]
}

func (s *session) Del(chat_id int64, key string) {
	(*s)[chat_id][key] = nil
}
