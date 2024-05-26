package ratelimiter

type Settings struct {
	Ratelimit      int
	ExpirationTime int
	LimitByToken   bool
}

func NewSettings(ratelimit, expirationTime int, limitByToken bool) *Settings {
	return &Settings{
		Ratelimit:      ratelimit,
		ExpirationTime: expirationTime,
		LimitByToken:   limitByToken}
}
