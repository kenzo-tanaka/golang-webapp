package main

import "errors"

// アバターのURLを返すことができないときに発生するエラー
var ErroNoAvatarURL = errors.New("chat: Unable to get an avatar URL.")

type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}
