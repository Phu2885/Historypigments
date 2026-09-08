package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const minioBaseURL = "http://localhost:9000/historypigments"

type Pigment struct {
	ID          int
	Name        string
	Description string
	Color       string
	Reflection  float64
	Absorption  float64
	ImageKey    string
	VideoKey    string
	Status      string
	UserLikes   []int // Collection of user IDs who liked this
}

// LikesCount calculates likes from the collection
func (p Pigment) LikesCount() int {
	return len(p.UserLikes)
}

func (p Pigment) ImageURL() string {
	if strings.TrimSpace(p.ImageKey) == "" {
		return "/assets/images/placeholder.jpg"
	}
	return minioBaseURL + "/" + p.ImageKey
}

func (p Pigment) VideoURL() string {
	if strings.TrimSpace(p.VideoKey) == "" {
		return ""
	}
	return minioBaseURL + "/" + p.VideoKey
}

var pigments = []Pigment{
	{
		ID:          1,
		Name:        "Киноварь",
		Description: "Исторический красный пигмент, использовавшийся с древности.",
		Color:       "Красный",
		Reflection:  0.72,
		Absorption:  0.28,
		ImageKey:    "cinnabar.jpg",
		VideoKey:    "cinnabar.mp4",
		Status:      "опубликован",
		UserLikes:   []int{1, 2, 3, 4, 5},
	},
	{
		ID:          2,
		Name:        "Жёлтая охра",
		Description: "Природный земляной пигмент, смесь гидратов оксида железа с глиной.",
		Color:       "Жёлтый",
		Reflection:  0.68,
		Absorption:  0.32,
		ImageKey:    "yellow-ochre.jpg",
		VideoKey:    "yellow-ochre.mp4",
		Status:      "опубликован",
		UserLikes:   []int{6, 7},
	},
	{
		ID:          3,
		Name:        "Малахит",
		Description: "Зелёный пигмент, получаемый из размолотого минерала малахита.",
		Color:       "Зелёный",
		Reflection:  0.64,
		Absorption:  0.36,
		ImageKey:    "malachite.jpg",
		VideoKey:    "malachite.mp4",
		Status:      "опубликован",
		UserLikes:   []int{8, 9, 10},
	},
	{
		ID:          4,
		Name:        "Лазурит",
		Description: "Синий пигмент исключительной чистоты và яркости.",
		Color:       "Синий",
		Reflection:  0.62,
		Absorption:  0.38,
		ImageKey:    "lapis-lazuli.jpg",
		VideoKey:    "lapis-lazuli.mp4",
		Status:      "опубликован",
		UserLikes:   []int{11},
	},
	{
		ID:          5,
		Name:        "Ультрамарин",
		Description: "Ярко-синий пигмент, исторически ценившийся на вес золота.",
		Color:       "Синий",
		Reflection:  0.60,
		Absorption:  0.40,
		ImageKey:    "ultramarine.jpg",
		VideoKey:    "ultramarine.mp4",
		Status:      "опубликован",
		UserLikes:   []int{12, 13, 14, 15, 16},
	},
	{
		ID:          6,
		Name:        "Умбра",
		Description: "Коричневый земляной пигмент из оксидов железа и марганца.",
		Color:       "Коричневый",
		Reflection:  0.58,
		Absorption:  0.42,
		ImageKey:    "umber.jpg",
		VideoKey:    "umber.mp4",
		Status:      "опубликован",
		UserLikes:   []int{17},
	},
	{
		ID:          7,
		Name:        "Новый образец (Черновик)",
		Description: "Данные для нового исследования фрагмента живописи.",
		Color:       "Не определен",
		Reflection:  0.0,
		Absorption:  0.0,
		ImageKey:    "draft.jpg",
		VideoKey:    "",
		Status:      "черновик",
		UserLikes:   []int{},
	},
}

func main() {
	tmpl := template.Must(template.ParseGlob(filepath.Join("./templates", "*.html")))

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./templates/assets"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/reconstruction-pigments", http.StatusFound)
	})

	mux.HandleFunc("/reconstruction-pigments", func(w http.ResponseWriter, r *http.Request) {
		pigment, ok := findPigment(r)
		if !ok {
			for _, p := range pigments {
				if p.Status == "опубликован" {
					pigment = p
					break
				}
			}
		}
		data := struct {
			Pigment Pigment
			Active  string
		}{Pigment: pigment, Active: "lenta"}
		if err := tmpl.ExecuteTemplate(w, "reconstruction-pigments.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/pigment-analysis-request", func(w http.ResponseWriter, r *http.Request) {
		var draft Pigment
		for _, p := range pigments {
			if p.Status == "черновик" {
				draft = p
				break
			}
		}
		data := struct {
			Pigment Pigment
			Active  string
		}{Pigment: draft, Active: "add"}
		if err := tmpl.ExecuteTemplate(w, "pigment-analysis-request.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/pigments-database", func(w http.ResponseWriter, r *http.Request) {
		filterStr := r.URL.Query().Get("filter")
		filtered := filterPigments(filterStr)

		data := struct {
			Pigments []Pigment
			Filter   string
			Active   string
		}{
			Pigments: filtered,
			Filter:   filterStr,
			Active:   "tiles",
		}
		if err := tmpl.ExecuteTemplate(w, "pigments-database.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Printf("Alchemy of Art History serving on http://localhost:8080\n")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}

func findPigment(r *http.Request) (Pigment, bool) {
	idRaw := r.URL.Query().Get("id")
	if idRaw == "" {
		return Pigment{}, false
	}
	id, _ := strconv.Atoi(idRaw)
	next := r.URL.Query().Get("next") == "true"

	var published []Pigment
	for _, p := range pigments {
		if p.Status == "опубликован" {
			published = append(published, p)
		}
	}

	for i, p := range published {
		if p.ID == id {
			if next {
				return published[(i+1)%len(published)], true
			}
			return p, true
		}
	}
	return Pigment{}, false
}

func filterPigments(filter string) []Pigment {
	var out []Pigment
	if filter == "" {
		for _, p := range pigments {
			if p.Status == "опубликован" {
				out = append(out, p)
			}
		}
		sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
		return out
	}

	fVal, err := strconv.ParseFloat(filter, 64)
	if err != nil {
		for _, p := range pigments {
			if p.Status == "опубликован" {
				out = append(out, p)
			}
		}
	} else {
		for _, p := range pigments {
			if p.Status == "опубликован" && p.Absorption <= fVal {
				out = append(out, p)
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
