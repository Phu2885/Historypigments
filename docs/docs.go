package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "swagger": "2.0",
    "schemes": [
        "http"
    ],
    "info": {
        "title": "Alchemy of Art History API",
        "version": "1.0",
        "description": "REST API каталога исторических пигментов: JWT, Redis blacklist, MinIO"
    },
    "host": "localhost:8080",
    "basePath": "/",
    "tags": [
        {
            "name": "pigments",
            "description": "Домен услуг — исторические пигменты"
        },
        {
            "name": "users",
            "description": "Домен пользователей — регистрация и авторизация"
        }
    ],
    "paths": {
        "/api/users/register": {
            "post": {
                "tags": [
                    "users"
                ],
                "operationId": "registerUser",
                "summary": "Регистрация пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/UserRegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Создано",
                        "schema": {
                            "$ref": "#/definitions/AuthResponse"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/users/login": {
            "post": {
                "tags": [
                    "users"
                ],
                "operationId": "loginUser",
                "summary": "Аутентификация, выдача JWT",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/UserLoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешно",
                        "schema": {
                            "$ref": "#/definitions/AuthResponse"
                        }
                    },
                    "401": {
                        "description": "Требуется авторизация",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/users/logout": {
            "post": {
                "tags": [
                    "users"
                ],
                "operationId": "logoutUser",
                "summary": "Деавторизация, JWT в blacklist Redis",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "responses": {
                    "200": {
                        "description": "Успешно",
                        "schema": {
                            "$ref": "#/definitions/MessageResponse"
                        }
                    },
                    "401": {
                        "description": "Требуется авторизация",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/users/{id}/pigments": {
            "get": {
                "tags": [
                    "users"
                ],
                "operationId": "getUserPigments",
                "summary": "Пользователь и его пигменты",
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "integer"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешно",
                        "schema": {
                            "$ref": "#/definitions/UserPigmentsResponse"
                        }
                    }
                }
            }
        },
        "/api/pigments": {
            "get": {
                "tags": [
                    "pigments"
                ],
                "operationId": "getPigmentsList",
                "summary": "Список опубликованных пигментов с фильтром по поглощению",
                "description": "Гость и создатель. С токеном возвращается is_creator=1 для своих пигментов (кнопка удаления).",
                "security": [
                    {},
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "name": "filter",
                        "in": "query",
                        "type": "number",
                        "required": false,
                        "description": "Максимальный коэффициент поглощения (absorption <= filter)"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешно",
                        "schema": {
                            "$ref": "#/definitions/PigmentsListResponse"
                        }
                    }
                }
            },
            "post": {
                "tags": [
                    "pigments"
                ],
                "operationId": "createPigmentDraft",
                "summary": "Создание черновика пигмента с изображением и видео (MinIO)",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "name": "name",
                        "in": "formData",
                        "type": "string",
                        "required": true,
                        "description": "Название пигмента"
                    },
                    {
                        "name": "short_description",
                        "in": "formData",
                        "type": "string",
                        "required": false
                    },
                    {
                        "name": "image_key",
                        "in": "formData",
                        "type": "string",
                        "required": false,
                        "description": "Ключ готового изображения в MinIO или локальный путь /assets/... (если файл не загружается)"
                    },
                    {
                        "name": "video_key",
                        "in": "formData",
                        "type": "string",
                        "required": false,
                        "description": "Ключ готового видео в MinIO или локальный путь /assets/..."
                    },
                    {
                        "name": "image",
                        "in": "formData",
                        "type": "file",
                        "required": false,
                        "description": "Изображение пигмента"
                    },
                    {
                        "name": "video",
                        "in": "formData",
                        "type": "file",
                        "required": false,
                        "description": "Видео пигмента"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Создано",
                        "schema": {
                            "$ref": "#/definitions/PigmentResponse"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Требуется авторизация",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/pigments/feed": {
            "get": {
                "tags": [
                    "pigments"
                ],
                "operationId": "getPigmentsFeed",
                "summary": "Лента: первый пигмент, пигмент по id или следующий",
                "security": [
                    {},
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "name": "id",
                        "in": "query",
                        "type": "integer",
                        "required": false
                    },
                    {
                        "name": "next",
                        "in": "query",
                        "type": "boolean",
                        "required": false,
                        "description": "true — следующий после id"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешно",
                        "schema": {
                            "$ref": "#/definitions/PigmentResponse"
                        }
                    },
                    "404": {
                        "description": "Не найдено",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/pigments/draft": {
            "get": {
                "tags": [
                    "pigments"
                ],
                "operationId": "getPigmentDraft",
                "summary": "Черновик пигмента текущего пользователя",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "responses": {
                    "200": {
                        "description": "Успешно",
                        "schema": {
                            "$ref": "#/definitions/PigmentResponse"
                        }
                    },
                    "401": {
                        "description": "Требуется авторизация",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/pigments/{id}": {
            "get": {
                "tags": [
                    "pigments"
                ],
                "operationId": "getPigmentById",
                "summary": "Пигмент по ID",
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "integer",
                        "description": "ID пигмента"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешно",
                        "schema": {
                            "$ref": "#/definitions/PigmentResponse"
                        }
                    },
                    "404": {
                        "description": "Не найдено",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "tags": [
                    "pigments"
                ],
                "operationId": "deletePigment",
                "summary": "Логическое удаление своего пигмента",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "integer",
                        "description": "ID пигмента"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешно",
                        "schema": {
                            "$ref": "#/definitions/MessageResponse"
                        }
                    },
                    "403": {
                        "description": "Нет доступа (чужой пигмент)",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/pigments/{id}/publish": {
            "put": {
                "tags": [
                    "pigments"
                ],
                "operationId": "publishPigmentDraft",
                "summary": "Публикация своего черновика пигмента",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "integer",
                        "description": "ID пигмента"
                    },
                    {
                        "name": "publish_data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/PigmentPublishRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешно",
                        "schema": {
                            "$ref": "#/definitions/PigmentResponse"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/pigments/{id}/like": {
            "post": {
                "tags": [
                    "pigments"
                ],
                "operationId": "likePigment",
                "summary": "Поставить (1) или отменить (0) лайк",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "required": true,
                        "type": "integer",
                        "description": "ID пигмента"
                    },
                    {
                        "name": "like",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/PigmentLikeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешно",
                        "schema": {
                            "$ref": "#/definitions/PigmentLikeResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "Pigment": {
            "type": "object",
            "required": [
                "id",
                "name"
            ],
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "short_description": {
                    "type": "string"
                },
                "status": {
                    "type": "string",
                    "description": "черновик / опубликован / удалён"
                },
                "image_url": {
                    "type": "string",
                    "description": "Ключ изображения в MinIO"
                },
                "video_url": {
                    "type": "string",
                    "description": "Ключ видео в MinIO"
                },
                "reflection": {
                    "type": "number"
                },
                "absorption": {
                    "type": "number"
                },
                "created_at": {
                    "type": "string",
                    "format": "date-time"
                },
                "formed_at": {
                    "type": "string",
                    "format": "date-time"
                },
                "creator_id": {
                    "type": "integer"
                },
                "is_creator": {
                    "type": "integer",
                    "description": "1 — пигмент текущего пользователя (можно удалить)"
                },
                "likes": {
                    "type": "integer"
                },
                "is_like": {
                    "type": "boolean"
                }
            }
        },
        "PigmentsListResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Pigment"
                    }
                }
            }
        },
        "PigmentResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "data": {
                    "$ref": "#/definitions/Pigment"
                }
            }
        },
        "PigmentPublishRequest": {
            "type": "object",
            "required": [
                "short_description",
                "reflection",
                "absorption"
            ],
            "properties": {
                "short_description": {
                    "type": "string",
                    "example": "Лазурно-синий медный пигмент"
                },
                "reflection": {
                    "type": "number",
                    "example": 0.75
                },
                "absorption": {
                    "type": "number",
                    "example": 0.25
                }
            }
        },
        "PigmentLikeRequest": {
            "type": "object",
            "required": [
                "value"
            ],
            "properties": {
                "value": {
                    "type": "integer",
                    "enum": [
                        0,
                        1
                    ]
                }
            }
        },
        "PigmentLikeResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "likes": {
                    "type": "integer"
                },
                "is_like": {
                    "type": "boolean"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "User": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "login": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                }
            }
        },
        "UserPigmentsResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                },
                "data": {
                    "type": "object",
                    "properties": {
                        "user": {
                            "$ref": "#/definitions/User"
                        },
                        "pigments": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/Pigment"
                            }
                        }
                    }
                }
            }
        },
        "UserRegisterRequest": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "login": {
                    "type": "string",
                    "example": "expert"
                },
                "password": {
                    "type": "string",
                    "example": "pass123"
                }
            }
        },
        "UserLoginRequest": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string",
                    "example": "admin"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "AuthResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "access_token": {
                    "type": "string"
                },
                "expires_in": {
                    "type": "integer"
                },
                "token_type": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "MessageResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "ErrorResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header",
            "description": "Bearer <access_token>"
        }
    },
    "produces": [
        "application/json"
    ],
    "consumes": [
        "application/json"
    ]
}`

type s struct{}

func (s *s) ReadDoc() string {
	return docTemplate
}

func init() {
	swag.Register(swag.Name, &s{})
}
