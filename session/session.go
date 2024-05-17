package session

type Session interface {
	Get(userID int64) *UserData
}

type session struct {
	userData map[int64]*UserData
}

func New() *session {
	return &session{
		userData: make(map[int64]*UserData),
	}
}

func (e *session) Get(userID int64) *UserData {
	if ud := e.userData[userID]; ud == nil {
		e.init(userID)
	}
	return e.userData[userID]
}

func (e *session) init(userID int64) {
	e.userData[userID] = &UserData{
		userID: userID,
		scene:  0,
	}
}
