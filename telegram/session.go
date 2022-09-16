package telegram

// Session is not multiuser, it can be if we make it a map[user(string)]map[key(string)]interface{}{}

var session = map[string]interface{}{}

func SetSession(key string, value interface{}) {
	session[key] = value
}

func GetSession(key string) interface{} {
	return session[key]
}

func DelSession(key string) {
	session[key] = nil
}
