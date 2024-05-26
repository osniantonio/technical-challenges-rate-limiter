package configs

import "testing"

func TestLoadConfig(t *testing.T) {
	path := "."
	config, err := LoadConfig(path)
	expected := &Config{
		DBProtocol:     "redis",
		DBHost:         "redis",
		DBPort:         "6379",
		DBUser:         "",
		DBPassword:     "lI6sI5dS8nZ0lG6p",
		DBDatabase:     "0",
		LimitByToken:   true,
		RateLimit:      10,
		ExpirationTime: 60,
	}
	if *config != *expected || err != nil {
		t.Errorf("LoadConfig(%v) = (%v, %v), want (%v, %v)", path, config, err, expected, nil)
	}
}
