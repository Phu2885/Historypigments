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
    Wavelength  float64
    Reflection  float64
    Absorption  float64
    ImageKey    string
    VideoKey    string
    Status      string
    Likes       int
}

func (p Pigment) ImageURL() string {
    // If no ImageKey provided, fall back to local placeholder so templates always have a valid image
    if strings.TrimSpace(p.ImageKey) == "" {
        return "/assets/local.jpg"
    }
    return minioBaseURL + "/" + p.ImageKey
}

func (p Pigment) VideoURL() string {
    // Return empty when no video key to avoid rendering broken <video> tags
    if strings.TrimSpace(p.VideoKey) == "" {
        return ""
    }
    return minioBaseURL + "/" + p.VideoKey
}

var pigments = []Pigment{
    {
        ID:          1,
        Name:        "Киноварь",
        Description: "Исторический красный пигмент",
        Color:       "Красный",
        Wavelength:  610,
        Reflection:  0.72,
        Absorption:  0.28,
        ImageKey:    "cinnabar.jpg",
        VideoKey:    "cinnabar.mp4",
        Status:      "опубликован",
        Likes:       124,
    },
    {
        ID:          2,
        Name:        "Жёлтая охра",
        Description: "Жёлтый исторический пигмент",
        Color:       "Жёлтый",
        Wavelength:  580,
        Reflection:  0.68,
        Absorption:  0.32,
        ImageKey:    "yellow-ochre.jpg",
        VideoKey:    "yellow-ochre.mp4",
        Status:      "опубликован",
        Likes:       92,
    },
    {
        ID:          3,
        Name:        "Малахит",
        Description: "Зелёный исторический пигмент",
        Color:       "Зелёный",
        Wavelength:  530,
        Reflection:  0.64,
        Absorption:  0.36,
        ImageKey:    "malachite.jpg",
        VideoKey:    "malachite.mp4",
        Status:      "опубликован",
        Likes:       71,
    },
    {
        ID:          4,
        Name:        "Лазурит",
        Description: "Синий исторический пигмент",
        Color:       "Синий",
        Wavelength:  470,
        Reflection:  0.62,
        Absorption:  0.38,
        ImageKey:    "lapis-lazuli.jpg",
        VideoKey:    "lapis-lazuli.mp4",
        Status:      "опубликован",
        Likes:       154,
    },
    {
        ID:          5,
        Name:        "Ультрамарин",
        Description: "Синий и фиолетовый пигмент",
        Color:       "Синий",
        Wavelength:  450,
        Reflection:  0.60,
        Absorption:  0.40,
        ImageKey:    "ultramarine.jpg",
        VideoKey:    "ultramarine.mp4",
        Status:      "опубликован",
        Likes:       236,
    },
    {
        ID:          6,
        Name:        "Умбра",
        Description: "Коричневый исторический пигмент",
        Color:       "Коричневый",
        Wavelength:  490,
        Reflection:  0.58,
        Absorption:  0.42,
        ImageKey:    "umber.jpg",
        VideoKey:    "umber.mp4",
        Status:      "опубликован",
        Likes:       61,
    },
    {
        ID:          7,
        Name:        "Киноварь черновик",
        Description: "Черновик для добавления пигмента",
        Color:       "Красный",
        Wavelength:  610,
        Reflection:  0.72,
        Absorption:  0.28,
        ImageKey:    "cinnabar-draft.jpg",
        VideoKey:    "cinnabar-draft.mp4",
        Status:      "черновик",
        Likes:       0,
    },
}

func main() {
    // Resolve templates and assets directories with fallbacks so running from repo root or historypigments/ works
    templatesDir, assetsDir := resolveDataDirs()

    tmpl := template.Must(template.ParseGlob(filepath.Join(templatesDir, "*.html")))

    mux := http.NewServeMux()
    mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsDir))))
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        http.Redirect(w, r, "/lenta?id=1", http.StatusFound)
    })
    mux.HandleFunc("/lenta", func(w http.ResponseWriter, r *http.Request) {
        pigment, ok := findPigment(r)
        if !ok {
            pigment = pigments[0]
        }
        data := struct {
            Pigment Pigment
            Active  string
        }{Pigment: pigment, Active: "lenta"}
        if err := tmpl.ExecuteTemplate(w, "lenta.html", data); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })
    mux.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
        pigment := getDraft()
        data := struct {
            Pigment Pigment
            Active  string
        }{Pigment: pigment, Active: "add"}
        if err := tmpl.ExecuteTemplate(w, "add.html", data); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })
    mux.HandleFunc("/tiles", func(w http.ResponseWriter, r *http.Request) {
        parsed := r.URL.Query().Get("wavelength")
        filtered := filterPigmentsByWavelength(parsed)
        data := struct {
            Pigments []Pigment
            Filter   string
            Active   string
        }{
            Pigments: filtered,
            Filter:   parsed,
            Active:   "tiles",
        }
        if err := tmpl.ExecuteTemplate(w, "tiles.html", data); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })

    fmt.Println("MinIO base:", minioBaseURL)
    fmt.Println("Serving on http://localhost:8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        panic(err)
    }
}

func findPigment(r *http.Request) (Pigment, bool) {
    idRaw := strings.TrimSpace(r.URL.Query().Get("id"))
    if idRaw == "" {
        return pigments[0], true
    }
    id, err := strconv.Atoi(idRaw)
    if err != nil {
        return Pigment{}, false
    }
    if strings.TrimSpace(r.URL.Query().Get("next")) == "true" {
        for i, p := range pigments {
            if p.ID == id {
                if i+1 < len(pigments) {
                    return pigments[i+1], true
                }
                return pigments[0], true
            }
        }
    }
    for _, p := range pigments {
        if p.ID == id {
            return p, true
        }
    }
    return Pigment{}, false
}

func getDraft() Pigment {
    for _, p := range pigments {
        if p.Status == "черновик" {
            return p
        }
    }
    return pigments[0]
}

func filterPigmentsByWavelength(value string) []Pigment {
    q := strings.TrimSpace(value)
    out := make([]Pigment, 0)
    for _, p := range pigments {
        if p.Status != "опубликован" {
            continue
        }
        if q == "" {
            out = append(out, p)
            continue
        }
        filter, err := strconv.Atoi(q)
        if err != nil {
            out = append(out, p)
            continue
        }
        if int(p.Wavelength) == filter {
            out = append(out, p)
        }
    }
    return sortedPigments(out)
}

func sortedPigments(items []Pigment) []Pigment {
    out := append([]Pigment(nil), items...)
    sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
    return out
}

// resolveDataDirs tries a few candidate locations for templates/ and assets/ so the
// server can be started from the repo root (go run ./historypigments) or from
// inside the historypigments/ directory. It falls back to the naive relative
// paths if nothing else is found.
func resolveDataDirs() (templatesDir, assetsDir string) {
    candidatesTemplates := []string{"./templates", "./historypigments/templates"}
    candidatesAssets := []string{"./assets", "./historypigments/assets"}

    for _, c := range candidatesTemplates {
        if fi, err := os.Stat(c); err == nil && fi.IsDir() {
            templatesDir = c
            break
        }
    }
    if templatesDir == "" {
        templatesDir = "./templates"
    }

    for _, c := range candidatesAssets {
        if fi, err := os.Stat(c); err == nil && fi.IsDir() {
            assetsDir = c
            break
        }
    }
    if assetsDir == "" {
        assetsDir = "./assets"
    }

    fmt.Printf("Using templates dir: %s, assets dir: %s\n", templatesDir, assetsDir)
    return
}
