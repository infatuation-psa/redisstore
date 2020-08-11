package redisstore

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/sessions"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	redisAddr = "localhost:6379"
)

func TestNew(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	store, err := NewRedisStore(client)
	if err != nil {
		t.Fatal("failed to create redis store", err)
	}

	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}

	session, err := store.New(req, "hello")
	if err != nil {
		t.Fatal("failed to create session", err)
	}
	if session.IsNew == false {
		t.Fatal("session is not new")
	}
}

func TestOptions(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	store, err := NewRedisStore(client)
	if err != nil {
		t.Fatal("failed to create redis store", err)
	}

	opts := sessions.Options{
		Path:   "/path",
		MaxAge: 99999,
	}
	store.Options(opts)

	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}

	session, err := store.New(req, "hello")
	if session.Options.Path != opts.Path || session.Options.MaxAge != opts.MaxAge {
		t.Fatal("failed to set options")
	}
}

func TestSave(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	store, err := NewRedisStore(client)
	if err != nil {
		t.Fatal("failed to create redis store", err)
	}

	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}
	w := httptest.NewRecorder()

	session, err := store.New(req, "hello")
	if err != nil {
		t.Fatal("failed to create session", err)
	}

	session.Values["key"] = "value"
	err = session.Save(req, w)
	if err != nil {
		t.Fatal("failed to save: ", err)
	}
}

func TestDelete(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	store, err := NewRedisStore(client)
	if err != nil {
		t.Fatal("failed to create redis store", err)
	}

	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}
	w := httptest.NewRecorder()

	session, err := store.New(req, "hello")
	if err != nil {
		t.Fatal("failed to create session", err)
	}

	session.Values["key"] = "value"
	err = session.Save(req, w)
	if err != nil {
		t.Fatal("failed to save session: ", err)
	}

	session.Options.MaxAge = -1
	err = session.Save(req, w)
	if err != nil {
		t.Fatal("failed to delete session: ", err)
	}
}
