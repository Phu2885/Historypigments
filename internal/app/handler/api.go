package handler

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"rip/internal/app/ds"
	"rip/internal/app/repository"
)

// Helper: Отправка JSON ответа
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

// Helper: Отправка JSON ошибки
func writeErrorJSON(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]interface{}{
		"status":  "fail",
		"message": message,
	})
}

// Helper: Извлечение ID из URL path (например, /api/pigments/15 или /api/users/1/pigments -> 15 / 1)
func extractIDFromPath(path, prefix string) (uint, error) {
	trimmed := strings.TrimPrefix(path, prefix)
	trimmed = strings.Trim(trimmed, "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) == 0 || parts[0] == "" {
		return 0, fmt.Errorf("invalid id")
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid id")
	}
	return uint(id), nil
}

// ==========================================
// ДОМЕН УСЛУГИ (PIGMENTS)
// ==========================================

// 1. GET /api/pigments — Список с фильтрацией (только опубликованные) + признак is_creator (0/1)
func (h *Handler) GetPigmentsListAPI(w http.ResponseWriter, r *http.Request) {
	currentUserID := repository.GetCurrentUserID()
	var absorptionFilter *float64

	filterStr := r.URL.Query().Get("filter")
	if filterStr != "" {
		val, err := strconv.ParseFloat(filterStr, 64)
		if err == nil {
			absorptionFilter = &val
		}
	}

	pigments, err := h.Repo.GetPublishedPigmentsList(absorptionFilter, currentUserID)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   pigments,
	})
}

// 2. GET /api/pigments/feed — Лента (возвращает 1 объект пигмента для рилс/ленты)
func (h *Handler) GetFeedAPI(w http.ResponseWriter, r *http.Request) {
	idRaw := r.URL.Query().Get("id")
	next := r.URL.Query().Get("next") == "true"
	id, _ := strconv.Atoi(idRaw)

	pigment, err := h.Repo.GetFeedItem(uint(id), next)
	if err != nil || pigment == nil {
		writeErrorJSON(w, http.StatusNotFound, "нет опубликованных пигментов в ленте")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   pigment,
	})
}

// 3. GET /api/pigments/draft — Получение черновика текущего пользователя (без указывания ID)
func (h *Handler) GetDraftAPI(w http.ResponseWriter, r *http.Request) {
	currentUserID := repository.GetCurrentUserID()
	draft, err := h.Repo.GetUserDraft(currentUserID)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   draft,
	})
}

// 4. POST /api/pigments — Добавление черновика + загрузка файлов фото и видео в MinIO (латинские имена)
func (h *Handler) CreateDraftAPI(w http.ResponseWriter, r *http.Request) {
	currentUserID := repository.GetCurrentUserID()

	var pigment ds.Pigment
	var imageHeader, videoHeader *multipart.FileHeader

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err == nil {
			pigment.Name = r.FormValue("name")
			pigment.ShortDescription = r.FormValue("short_description")
			if imgKey := r.FormValue("image_key"); imgKey != "" {
				pigment.ImageURL = imgKey
			}
			if vidKey := r.FormValue("video_key"); vidKey != "" {
				pigment.VideoURL = vidKey
			}

			_, headerImg, errImg := r.FormFile("image")
			if errImg == nil {
				imageHeader = headerImg
			}
			_, headerVid, errVid := r.FormFile("video")
			if errVid == nil {
				videoHeader = headerVid
			}
		}
	} else {
		_ = json.NewDecoder(r.Body).Decode(&pigment)
	}

	if strings.TrimSpace(pigment.Name) == "" {
		writeErrorJSON(w, http.StatusBadRequest, "поле 'name' обязательно")
		return
	}

	err := h.Repo.CreateDraftPigment(&pigment, currentUserID, imageHeader, videoHeader)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"data":    pigment,
		"message": "черновик пигмента успешно создан",
	})
}

// 5. PUT /api/pigments/:id/publish — Публикация черновика (смена статуса)
func (h *Handler) PublishPigmentAPI(w http.ResponseWriter, r *http.Request) {
	currentUserID := repository.GetCurrentUserID()
	path := strings.TrimSuffix(r.URL.Path, "/publish")
	id, err := extractIDFromPath(path, "/api/pigments/")
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "некорректный ID пигмента")
		return
	}

	var req struct {
		ShortDescription string  `json:"short_description"`
		Reflection       float64 `json:"reflection"`
		Absorption       float64 `json:"absorption"`
	}

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		_ = json.NewDecoder(r.Body).Decode(&req)
	} else {
		req.ShortDescription = r.FormValue("short_description")
		req.Reflection, _ = strconv.ParseFloat(r.FormValue("reflection"), 64)
		req.Absorption, _ = strconv.ParseFloat(r.FormValue("absorption"), 64)
	}

	err = h.Repo.PublishDraftPigment(id, currentUserID, req.ShortDescription, req.Reflection, req.Absorption)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, _ := h.Repo.GetPigmentByID(id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"data":    updated,
		"message": "пигмент успешно опубликован",
	})
}

// 6. DELETE /api/pigments/:id — Soft Delete (только пигменты этого пользователя)
func (h *Handler) DeletePigmentAPI(w http.ResponseWriter, r *http.Request) {
	currentUserID := repository.GetCurrentUserID()
	id, err := extractIDFromPath(r.URL.Path, "/api/pigments/")
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "некорректный ID пигмента")
		return
	}

	err = h.Repo.SoftDeletePigment(id, currentUserID)
	if err != nil {
		writeErrorJSON(w, http.StatusForbidden, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "пигмент успешно переведен в статус 'удалён'",
	})
}

// 7. POST /api/pigments/:id/like — Поставить (1) или отменить (0) лайк
func (h *Handler) ToggleLikeAPI(w http.ResponseWriter, r *http.Request) {
	currentUserID := repository.GetCurrentUserID()
	path := strings.TrimSuffix(r.URL.Path, "/like")
	id, err := extractIDFromPath(path, "/api/pigments/")
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "некорректный ID пигмента")
		return
	}

	value := 1
	if valueStr := r.URL.Query().Get("value"); valueStr != "" {
		value, _ = strconv.Atoi(valueStr)
	} else {
		var req struct {
			Value int `json:"value"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			value = req.Value
		}
	}

	likesCount, isLiked, err := h.Repo.ToggleLike(id, currentUserID, value)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	actionMsg := "лайк поставлен"
	if value == 0 {
		actionMsg = "лайк отменен"
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"id":      id,
		"likes":   likesCount,
		"is_like": isLiked,
		"message": actionMsg,
	})
}

// GET /api/pigments/:id — Получение подробной информации о пигменте
func (h *Handler) GetPigmentByIDAPI(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path, "/api/pigments/")
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "некорректный ID пигмента")
		return
	}

	pigment, err := h.Repo.GetPigmentByID(id)
	if err != nil {
		writeErrorJSON(w, http.StatusNotFound, fmt.Sprintf("пигмент с id %d не найден", id))
		return
	}

	fullPigment := ds.FullPigmentSerializer{
		Pigment: pigment,
		Creator: pigment.Creator,
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   fullPigment,
	})
}

// ==========================================
// ДОМЕН ПОЛЬЗОВАТЕЛЬ (USERS)
// ==========================================

// 8. POST /api/users/register — Регистрация пользователя
func (h *Handler) RegisterUserAPI(w http.ResponseWriter, r *http.Request) {
	var user ds.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "некорректный JSON запроса")
		return
	}

	if strings.TrimSpace(user.Login) == "" || strings.TrimSpace(user.Password) == "" {
		writeErrorJSON(w, http.StatusBadRequest, "поля 'login' и 'password' обязательны")
		return
	}
	if user.Name == "" {
		user.Name = user.Login
	}

	err := h.Repo.RegisterUser(&user)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	tokenStr, expiresIn, err := GenerateJWTToken(user.ID, user.Login, user.Role)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "ошибка генерации JWT токена")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"status":       "success",
		"message":      "пользователь успешно зарегистрирован",
		"user_id":      user.ID,
		"access_token": tokenStr,
		"expires_in":   int(expiresIn.Seconds()),
		"token_type":   "Bearer",
	})
}

// 9. POST /api/users/login — Аутентификация (вход в систему)
func (h *Handler) LoginUserAPI(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "некорректный JSON запроса")
		return
	}

	user, err := h.Repo.GetUserByLogin(req.Login)
	if err != nil || user.Password != req.Password {
		writeErrorJSON(w, http.StatusUnauthorized, "неверный логин или пароль")
		return
	}

	tokenStr, expiresIn, err := GenerateJWTToken(user.ID, user.Login, user.Role)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "ошибка генерации JWT токена")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":       "success",
		"message":      "успешная аутентификация",
		"access_token": tokenStr,
		"expires_in":   int(expiresIn.Seconds()),
		"token_type":   "Bearer",
	})
}

// 10. POST /api/users/logout — Деавторизация (добавление JWT токена в Redis Blacklist)
func (h *Handler) LogoutUserAPI(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, jwtPrefix) {
		writeErrorJSON(w, http.StatusBadRequest, "отсутствует заголовок 'Authorization: Bearer <token>'")
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, jwtPrefix)

	if h.RedisClient != nil {
		_ = h.RedisClient.WriteJWTToBlacklist(r.Context(), tokenStr, 24*time.Hour)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "вы успешно вышли из системы (токен внесен в blacklist в Redis)",
	})
}

// GET /api/users/:id/pigments — Вложенная сериализация профиля пользователя и созданных им пигментов
func (h *Handler) GetUserPigmentsAPI(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path, "/api/users/")
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "некорректный ID пользователя")
		return
	}

	userPigments, err := h.Repo.GetUserPigments(id)
	if err != nil {
		writeErrorJSON(w, http.StatusNotFound, fmt.Sprintf("пользователь с id %d не найден", id))
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   userPigments,
	})
}
