package telegram

// Session is not multiuser, it can be if we make it a map[user(string)]map[key(string)]any

var session = map[string]any{}

func SetSession(key string, value any) {
	session[key] = value
}

func GetSession(key string) any {
	return session[key]
}

func DelSession(key string) {
	session[key] = nil
}
