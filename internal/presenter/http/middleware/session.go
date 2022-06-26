package middleware

type Session interface {
	IsExpired() bool
}

type Sessions interface {
	New(user, ip string) Session
	Add(session Session)
	Get(user, ip string) (session Session, ok bool)
	Remove(session Session)
}
