package handler

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"rip/internal/app/ds"
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

// Helper: Валидация загружаемого файла-изображения
func isImage(contentType string) bool {
	imageTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/svg+xml",
	}
	for _, t := range imageTypes {
		if contentType == t {
			return true
		}
	}
	return false
}

func validateFileUpload(header *multipart.FileHeader) (int, error) {
	file, err := header.Open()
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("не удалось открыть файл")
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("не удалось прочитать файл")
	}

	contentType := http.DetectContentType(buffer)
	if !isImage(contentType) && !strings.Contains(header.Filename, ".svg") {
		return http.StatusBadRequest, fmt.Errorf("файл должен быть изображением")
	}

	return 0, nil
}

// Helper: Извлечение ID из URL path (например, /api/pigments/15 -> 15)
func extractIDFromPath(path, prefix string) (uint, error) {
	idStr := strings.TrimPrefix(path, prefix)
	idStr = strings.Trim(idStr, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid id")
	}
	return uint(id), nil
}

// 1. GET /api/pigments — получение списка всех пигментов
func (h *Handler) GetPigmentsAPI(w http.ResponseWriter, r *http.Request) {
	pigments, err := h.Repo.GetAllPigmentsAPI()
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   pigments,
	})
}

// 2. GET /api/pigments/{id} — получение пигмента по ID с вложенным автором
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

// 3. POST /api/pigments — добавление пигмента
func (h *Handler) AddPigmentAPI(w http.ResponseWriter, r *http.Request) {
	var pigment ds.Pigment
	var fileHeader *multipart.FileHeader
	fileFound := false

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(2 << 20); err == nil {
			pigment.Name = r.FormValue("name")
			pigment.ShortDescription = r.FormValue("short_description")
			pigment.Color = r.FormValue("color")
			pigment.Status = r.FormValue("status")
			if pigment.Status == "" {
				pigment.Status = "опубликован"
			}
			pigment.Reflection, _ = strconv.ParseFloat(r.FormValue("reflection"), 64)
			pigment.Absorption, _ = strconv.ParseFloat(r.FormValue("absorption"), 64)

			file, header, err := r.FormFile("pic")
			if err == nil {
				file.Close()
				fileHeader = header
				fileFound = true
			}
		}
	} else if strings.Contains(contentType, "application/json") {
		_ = json.NewDecoder(r.Body).Decode(&pigment)
	} else {
		pigment.Name = r.FormValue("name")
		pigment.ShortDescription = r.FormValue("short_description")
	}

	if strings.TrimSpace(pigment.Name) == "" {
		writeErrorJSON(w, http.StatusBadRequest, "поле 'name' обязательно")
		return
	}

	if fileFound {
		code, err := validateFileUpload(fileHeader)
		if err != nil {
			writeErrorJSON(w, code, err.Error())
			return
		}
	}

	err := h.Repo.AddPigment(&pigment)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if fileFound {
		_ = h.Repo.AddOrReplacePigmentImage(pigment.ID, fileHeader)
		updated, err := h.Repo.GetPigmentByID(pigment.ID)
		if err == nil {
			pigment = updated
		}
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"data":    pigment,
		"message": "пигмент успешно добавлен",
	})
}

// 4. PUT /api/pigments/{id} — обновление пигмента
func (h *Handler) ModifyPigmentAPI(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path, "/api/pigments/")
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "некорректный ID пигмента")
		return
	}

	existing, err := h.Repo.GetPigmentByID(id)
	if err != nil {
		writeErrorJSON(w, http.StatusNotFound, fmt.Sprintf("пигмент с id %d не найден", id))
		return
	}

	var fileHeader *multipart.FileHeader
	fileFound := false

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(2 << 20); err == nil {
			if r.FormValue("name") != "" {
				existing.Name = r.FormValue("name")
			}
			if r.FormValue("short_description") != "" {
				existing.ShortDescription = r.FormValue("short_description")
			}
			if r.FormValue("color") != "" {
				existing.Color = r.FormValue("color")
			}
			if r.FormValue("status") != "" {
				existing.Status = r.FormValue("status")
			}
			if r.FormValue("reflection") != "" {
				existing.Reflection, _ = strconv.ParseFloat(r.FormValue("reflection"), 64)
			}
			if r.FormValue("absorption") != "" {
				existing.Absorption, _ = strconv.ParseFloat(r.FormValue("absorption"), 64)
			}

			file, header, err := r.FormFile("pic")
			if err == nil {
				file.Close()
				fileHeader = header
				fileFound = true
			}
		}
	} else if strings.Contains(contentType, "application/json") {
		var input ds.Pigment
		_ = json.NewDecoder(r.Body).Decode(&input)
		if input.Name != "" {
			existing.Name = input.Name
		}
		if input.ShortDescription != "" {
			existing.ShortDescription = input.ShortDescription
		}
		if input.Reflection != 0 {
			existing.Reflection = input.Reflection
		}
		if input.Absorption != 0 {
			existing.Absorption = input.Absorption
		}
	}

	if fileFound {
		code, err := validateFileUpload(fileHeader)
		if err != nil {
			writeErrorJSON(w, code, err.Error())
			return
		}
	}

	err = h.Repo.ModifyPigment(id, &existing)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if fileFound {
		_ = h.Repo.AddOrReplacePigmentImage(id, fileHeader)
	}

	updated, _ := h.Repo.GetPigmentByID(id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"data":    updated,
		"message": "запись успешно обновлена",
	})
}

// 5. DELETE /api/pigments/{id} — удаление пигмента
func (h *Handler) DeletePigmentAPI(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path, "/api/pigments/")
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "некорректный ID пигмента")
		return
	}

	err = h.Repo.HardDeletePigment(id)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "пигмент успешно удален",
	})
}

// 6. GET /api/users/{id}/pigments — вложенная сериализация пользователя и его пигментов
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
