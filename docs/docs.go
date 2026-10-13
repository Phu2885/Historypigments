package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": ["http"],
    "swagger": "2.0",
    "info": {
        "description": "Historical Pigments Database REST API with JWT Authentication, Redis Blacklist, and MinIO SDK",
        "title": "Alchemy of Art History API (ЛР4)",
        "version": "1.0"
    },
    "host": "localhost:8080",
    "basePath": "/",
    "paths": {
        "/api/users/register": {
            "post": {
                "tags": ["Users / Auth"],
                "summary": "Регистрация нового пользователя",
                "responses": {
                    "201": { "description": "Успешно зарегистрирован" }
                }
            }
        },
        "/api/users/login": {
            "post": {
                "tags": ["Users / Auth"],
                "summary": "Аутентификация пользователя (Вход)",
                "description": "Возвращает JWT access_token при корректных учетных данных",
                "responses": {
                    "200": { "description": "Успешная аутентификация" },
                    "401": { "description": "Неверный логин или пароль" }
                }
            }
        },
        "/api/users/logout": {
            "post": {
                "tags": ["Users / Auth"],
                "summary": "Деавторизация (Выход)",
                "description": "Заносит JWT токен в Redis Blacklist",
                "security": [{ "BearerAuth": [] }],
                "responses": {
                    "200": { "description": "Успешный выход" }
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
        "/api/pigments/feed": {
            "get": {
                "tags": ["Pigments"],
                "summary": "Лента опубликованных пигментов (Без авторизации)",
                "responses": {
                    "200": { "description": "Успешный объект ленты" }
                }
            }
        },
        "/api/pigments/draft": {
            "get": {
                "tags": ["Pigments"],
                "summary": "Черновик текущего пользователя (Без авторизации)",
                "responses": {
                    "200": { "description": "Черновик пользователя" }
                }
            }
        },
        "/api/pigments/{id}": {
            "get": {
                "tags": ["Pigments"],
                "summary": "Получить пигмент по ID",
                "parameters": [{ "name": "id", "in": "path", "required": true, "type": "integer" }],
                "responses": { "200": { "description": "Успешный пигмент" } }
            },
            "delete": {
                "tags": ["Pigments"],
                "summary": "Удаление пигмента (Soft Delete - Требуется JWT)",
                "security": [{ "BearerAuth": [] }],
                "parameters": [{ "name": "id", "in": "path", "required": true, "type": "integer" }],
                "responses": { "200": { "description": "Удалено" } }
            }
        },
        "/api/pigments/{id}/publish": {
            "put": {
                "tags": ["Pigments"],
                "summary": "Публикация черновика пигмента (Требуется JWT Авторизация)",
                "security": [{ "BearerAuth": [] }],
                "parameters": [{ "name": "id", "in": "path", "required": true, "type": "integer" }],
                "responses": { "200": { "description": "Пигмент опубликован" } }
            }
        },
        "/api/pigments/{id}/like": {
            "post": {
                "tags": ["Pigments"],
                "summary": "Лайк от текущего пользователя (Требуется JWT)",
                "security": [{ "BearerAuth": [] }],
                "parameters": [{ "name": "id", "in": "path", "required": true, "type": "integer" }],
                "responses": { "200": { "description": "Лайк обновлен" } }
            }
        },
        "/api/users/{id}/pigments": {
            "get": {
                "tags": ["Users / Auth"],
                "summary": "Вложенная сериализация пользователя и его пигментов",
                "parameters": [{ "name": "id", "in": "path", "required": true, "type": "integer" }],
                "responses": { "200": { "description": "Пользователь и пигменты" } }
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
