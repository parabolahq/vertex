package communication

import "gopkg.in/olahol/melody.v1"

var Sessions = map[string][]*melody.Session{}

func StoreSession(userId string, session *melody.Session) {
	_, exists := Sessions[userId]
	if !exists {
		Sessions[userId] = []*melody.Session{session}
	} else {
		Sessions[userId] = append(Sessions[userId], session)
	}
}

func RemoveSession(userId string, session *melody.Session) {
	_, exists := Sessions[userId]
	if exists {
		for i, s := range Sessions[userId] {
			if s == session {
				Sessions[userId] = append(Sessions[userId][:i], Sessions[userId][i+1:]...)
				break
			}
		}
	}
}

func GetSessions(userId string) []*melody.Session {
	_, exists := Sessions[userId]
	if !exists {
		return []*melody.Session{}
	}
	return Sessions[userId]
}
