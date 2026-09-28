package handler

import (
	"html/template"
	"net/http"
	"strconv"
	"rip/internal/app/repository"
	"rip/internal/app/ds"
)

type Handler struct {
	Repo *repository.Repository
	Tmpl *template.Template
}

func New(repo *repository.Repository, tmpl ...*template.Template) *Handler {
	var t *template.Template
	if len(tmpl) > 0 {
		t = tmpl[0]
	}
	return &Handler{Repo: repo, Tmpl: t}
}

// 1. GET: Лента (Reconstruction)
func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	idRaw := r.URL.Query().Get("id")
	next := r.URL.Query().Get("next") == "true"

	published, _ := h.Repo.GetAllPublished()
	if len(published) == 0 {
		h.Tmpl.ExecuteTemplate(w, "pigments-database-reconstruction.html", nil)
		return
	}

	var current ds.Pigment
	id, _ := strconv.Atoi(idRaw)

	if next {
		found := false
		for _, p := range published {
			if p.ID > uint(id) {
				current = p
				found = true
				break
			}
		}
		if !found { current = published[0] }
	} else {
		found := false
		for _, p := range published {
			if p.ID == uint(id) {
				current = p
				found = true
				break
			}
		}
		if !found { current = published[0] }
	}

	likes, _ := h.Repo.GetLikesCount(current.ID)
	data := struct {
		Pigment ds.Pigment
		Likes   int64
		Active  string
	}{Pigment: current, Likes: likes, Active: "lenta"}
	h.Tmpl.ExecuteTemplate(w, "pigments-database-reconstruction.html", data)
}

// 2. GET: Анализ/Черновик
func (h *Handler) GetAnalysis(w http.ResponseWriter, r *http.Request) {
	draft, _ := h.Repo.GetUserDraft(1) // ID=1 захардкожен
	data := struct {
		Pigment *ds.Pigment
		Active  string
	}{Pigment: draft, Active: "add"}
	h.Tmpl.ExecuteTemplate(w, "pigment-database-analysis.html", data)
}

// 3. GET: База/Плитка
func (h *Handler) GetDatabase(w http.ResponseWriter, r *http.Request) {
	filterStr := r.URL.Query().Get("filter")
	var pigments []ds.Pigment
	if filterStr != "" {
		val, _ := strconv.ParseFloat(filterStr, 64)
		pigments, _ = h.Repo.FilterByAbsorption(val)
	} else {
		pigments, _ = h.Repo.GetAllPublished()
	}

	data := struct {
		Pigments []ds.Pigment
		Filter   string
		Active   string
	}{Pigments: pigments, Filter: filterStr, Active: "tiles"}
	h.Tmpl.ExecuteTemplate(w, "pigments-database.html", data)
}

// 4. POST: Создание черновика (через ORM)
func (h *Handler) CreateDraft(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	img := r.FormValue("image_key")
	vid := r.FormValue("video_key")

	draft, _ := h.Repo.GetUserDraft(1)
	if draft == nil {
		newP := &ds.Pigment{
			Name: name,
			ImageURL: img,
			VideoURL: vid,
			CreatorID: 1,
		}
		h.Repo.CreateDraft(newP)
	}
	http.Redirect(w, r, "/pigment-database-analysis", http.StatusFound)
}

// 5. POST: Публикация (через ORM)
func (h *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	desc := r.FormValue("description")
	refl, _ := strconv.ParseFloat(r.FormValue("reflection"), 64)
	abs, _ := strconv.ParseFloat(r.FormValue("absorption"), 64)

	h.Repo.PublishPigment(uint(id), desc, refl, abs)
	http.Redirect(w, r, "/pigments-database", http.StatusFound)
}

// 6. POST: Удаление (через SQL UPDATE)
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	h.Repo.LogicalDelete(uint(id))
	http.Redirect(w, r, "/pigments-database", http.StatusFound)
}
