package api

import "github.com/Dokhoyan/go-messenger-microservices/user_service/internal/model"

// RouteAccesses - содержит название ручки и роли, которые имеют к ней доступ
var RouteAccesses = map[string][]model.UserRole{
	"/chat_v1.ChatV1/SendMessage": {model.USER},
	"/chat_v1.ChatV1/Create":      {model.USER},
	"/chat_v1.ChatV1/Delete":      {model.ADMIN},
}