package main

import (
	"net/http"
	"net/http/httptest"

	"github.com/mrjones/oauth"
)

type Token struct {
	consumerKey    string
	consumerSecret string

	consumer *oauth.Consumer
	server   *httptest.Server
	reqToken *oauth.RequestToken
	ch       chan *oauth.AccessToken
}

func NewToken(ck, cs string) *Token {
	token := &Token{
		consumerKey:    ck,
		consumerSecret: cs,

		ch: make(chan *oauth.AccessToken),
		consumer: oauth.NewConsumer(ck, cs, oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		}),
	}

	return token
}

func (t *Token) URL() string {
	t.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		code := values.Get("oauth_verifier")

		aToken, err := t.consumer.AuthorizeToken(t.reqToken, code)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(400)
			return
		}

		w.Write([]byte("Success!"))
		t.ch <- aToken
	}))

	callbackURL := t.server.URL
	token, reqURL, err := t.consumer.GetRequestTokenAndUrl(callbackURL)
	if err != nil {
		panic(err)
	}
	t.reqToken = token

	return reqURL
}

func (t *Token) AccessToken() *oauth.AccessToken {
	defer t.server.Close()
	return <-t.ch
}
