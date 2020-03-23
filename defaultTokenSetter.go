package sniproxy

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type defaultAuthTokenType int

const (
	COOKIE defaultAuthTokenType = iota
	HEADER
	EITHER
)

var defaultAuthTokenTypes = [...]string{
	"COOKIE",
	"HEADER",
	"EITHER",
}

func (a defaultAuthTokenType) String() string {
	return defaultAuthTokenTypes[a]
}

func enumerateDefaultAuthTokenType(authTokenString string) (defaultAuthTokenType, error) {
	authTokenString = strings.ToUpper(authTokenString)
	for i, v := range defaultAuthTokenTypes {
		if authTokenString == v {
			return defaultAuthTokenType(i), nil
		}
	}
	return EITHER, fmt.Errorf("Sniproxy error - Unsupported AuthTokenString passed for the route")
}

type defaultTokenSetter struct {
	tokenValue string
	tokenName  string
}

func (c *defaultTokenSetter) SetToken(w http.ResponseWriter, r *http.Request,
	tokenTypeString string) {
	tokenType, enumerationError := enumerateDefaultAuthTokenType(tokenTypeString)
	// if the tokenTypeString cannot be enumerated, we will check everywhere
	if enumerationError != nil {
		tokenType = EITHER
	}
	switch tokenType {
	case COOKIE:
		cookie := &http.Cookie{
			Name:     c.tokenName,
			Value:    c.tokenValue,
			Expires:  time.Now(),
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	case HEADER:
		w.Header().Add(c.tokenName, c.tokenValue)
	case EITHER:
		c.SetToken(w, r, "Cookie")
		c.SetToken(w, r, "Header")
	}
}