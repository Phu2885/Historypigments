package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": ["http"],
    "swagger": "2.0",
    "info": {
        "description": "Historical Pigments Database REST API with JWT Authentication and MinIO Object Storage",
        "title": "Alchemy of Art History API",
        "version": "1.0"
    },
    "host": "localhost:8080",
    "basePath": "/",
    "paths": {
        "/api/users/login": {
            "post": {
                "tags": ["Users / Auth"],
                "summary": "Аутентификация пользователя (Вход)",
                "description": "Возвращает JWT access_token при корректных учетных данных",
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "login": { "type": "string", "example": "test" },
                                "password": { "type": "string", "example": "test123" }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": { "description": "Успешная аутентификация" },
                    "401": { "description": "Неверный логин или пароль" }
                }
            }
        },
        "/api/pigments": {
            "get": {
                "tags": ["Pigments"],
                "summary": "Список опубликованных пигментов (Без авторизации)",
                "parameters": [
                    {
                        "name": "filter",
                        "in": "query",
                        "type": "number",
                        "description": "Фильтр по поглощению (absorption <= filter)"
                    }
                ],
                "responses": {
                    "200": { "description": "Успешный список пигментов" }
                }
            },
            "post": {
                "tags": ["Pigments"],
                "summary": "Создание черновика (Требуется JWT Авторизация)",
                "security": [{ "BearerAuth": [] }],
                "responses": {
                    "201": { "description": "Черновик создан" },
                    "401": { "description": "Требуется авторизация" }
                }
            }
        },
        "/api/pigments/{id}/publish": {
            "put": {
                "tags": ["Pigments"],
                "summary": "Публикация черновика пигмента (Требуется JWT Авторизация)",
                "security": [{ "BearerAuth": [] }],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "integer"
                    }
                ],
                "responses": {
                    "200": { "description": "Пигмент опубликован" },
                    "401": { "description": "Неавторизован" }
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type s struct{}

func (s *s) ReadDoc() string {
	return docTemplate
}

func init() {
	swag.Register(swag.Name, &s{})
}
