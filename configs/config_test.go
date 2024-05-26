package configs

import "testing"

func TestLoadConfig(t *testing.T) {
	path := "."
	config, err := LoadConfig(path)
	expected := &Config{
		DBProtocol:     "redis",
		DBHost:         "redis",
		DBPort:         "6379",
		DBPassword:     "4K!iHeNgR32",
		DBDatabase:     "0",
		LimitByToken:   true,
		RateLimit:      10,
		ExpirationTime: 60,
	}
	if *config != *expected || err != nil {
		t.Errorf("LoadConfig(%v) = (%v, %v), want (%v, %v)", path, config, err, expected, nil)
	}
}
