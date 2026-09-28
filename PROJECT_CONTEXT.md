# PROJECT CONTEXT SUMMARY: Alchemy of Art History (Pigments Database)

---

## 1. Executive Summary & Goals
* **Project Name**: `rip` (Alchemy of Art History — Historical Pigments Database / Анализ и реконструкция пигментов)
* **Git Branch**: `SPA` (Active branch for Lab 3 REST API development)
* **Objective**: A web application designed for art history expertise and conservation analysis. It manages a catalog of historical pigments, supports multi-step pigment analysis/drafting, media integration via MinIO, absorption filtering, feed-style reconstruction views, and user feedback (likes).
* **Lab Context**: Developed as part of a Web Development / Database course laboratory work using Go, GORM, PostgreSQL, MinIO, and HTML templates.

---

## 2. Tech Stack & Dependencies
* **Programming Language**: Go 1.2x
* **Router & Server**: Standard Library `net/http` (`http.ServeMux`), `html/template`
* **ORM & Database**: GORM (`gorm.io/gorm`, `gorm.io/driver/postgres`), PostgreSQL 15
* **Object Storage**: MinIO (`historypigments` bucket, public read policy)
* **Environment Management**: `github.com/joho/godotenv` (`.env` configuration)
* **Containerization**: Docker Compose (`postgres`, `minio`, `adminer`)

---

## 3. Architecture & File Structure

```
rip/
├── cmd/
│   └── migrate/
│       └── main.go               # Database AutoMigrate script for User, Pigment, Like
├── internal/
│   └── app/
│       ├── ds/                   # Data Structures / GORM Models
│       │   ├── user.go           # User model
│       │   ├── pigment.go        # Pigment model + MinIO URL helper functions
│       │   └── like.go           # Like model (User <-> Pigment relationship)
│       ├── handler/
│       │   └── handler.go        # HTTP Handlers (Feed, Draft/Analysis, Catalog, Actions)
│       └── repository/
│           └── repository.go     # GORM & Raw SQL database operations
├── templates/                    # HTML Templates (Go html/template)
│   ├── assets/
│   │   └── css/
│   │       └── styles.css        # Core stylesheet for dark theme UI
│   ├── pigment-database-analysis.html          # Draft creation & publication form
│   ├── pigments-database-reconstruction.html  # Vertical Feed / Reconstruction view
│   └── pigments-database.html                 # Grid / Tile catalog view with absorption filter
├── scripts/
│   └── init-minio.sh             # Shell script to configure MinIO bucket & access policy
├── docker-compose.yml            # Docker setup for Postgres, MinIO, and Adminer
├── main.go                       # Web server entry point, routing, DB initialization
├── go.mod                        # Go module file
└── .env                          # Environment variables (DB credentials, etc.)
```

---

## 4. Database Schema & Data Models

### Models (`internal/app/ds/`)

1. **User** (`users` table):
   - `ID` (`uint`, primary key)
   - `Name` (`string`, `varchar(100)`, not null)
   - `Login` (`string`, `varchar(50)`, unique, not null)
   - `Password` (`string`, `varchar(100)`, not null)

2. **Pigment** (`pigments` table):
   - `ID` (`uint`, primary key)
   - `Name` (`string`, `varchar(100)`, not null)
   - `ShortDescription` (`string`, `text`)
   - `Status` (`string`, `varchar(20)`, default: `'черновик'`) — Allowed values: `'черновик'`, `'опубликован'`, `'удалён'`
   - `ImageURL` (`string`, `varchar(255)`)
   - `VideoURL` (`string`, `varchar(255)`)
   - `Color` (`string`, `varchar(50)`)
   - `Reflection` (`float64`, `decimal(5,2)`) — Reflectance spectrum value (R)
   - `Absorption` (`float64`, `decimal(5,2)`) — Absorption spectrum value (A)
   - `CreatedAt` (`time.Time`, not null)
   - `FormedAt` (`*time.Time`) — Timestamp when draft was published
   - `CreatorID` (`uint`, foreign key -> User)
   - Helpers: `GetImage()`, `GetVideo()` -> Returns full MinIO URL (`http://localhost:9000/historypigments/...`) or default placeholder.

3. **Like** (`likes` table):
   - `ID` (`uint`, primary key)
   - `UserID` (`uint`, foreign key -> User)
   - `PigmentID` (`uint`, foreign key -> Pigment)

---

## 5. Endpoints & Route Map

| Method | Path | Handler | Description |
|---|---|---|---|
| `GET` | `/` | Redirect | Redirects to `/pigments-database-reconstruction` |
| `GET` | `/pigments-database-reconstruction` | `h.GetFeed` | Renders Feed view for published pigments with video/image, description, likes count, and "NEXT" navigation (`?id=X&next=true`). |
| `GET` | `/pigment-database-analysis` | `h.GetAnalysis` | Multi-step view: Step 1 (create draft if none exists), Step 2 (publish draft with description, reflection, absorption). |
| `GET` | `/pigments-database` | `h.GetDatabase` | Grid/Tile view of published pigments with range slider filter by absorption (`?filter=0.5`). |
| `POST` | `/create-draft` | `h.CreateDraft` | Creates a new draft pigment record (`status = "черновик"`). Currently hardcoded to `CreatorID = 1`. |
| `POST` | `/publish-pigment` | `h.Publish` | Updates draft status to `"опубликован"`, sets description, reflection, absorption, and `formed_at` timestamp. |
| `POST` | `/delete-pigment` | `h.Delete` | Performs **logical deletion** (`status = "удалён"`) via **Raw SQL UPDATE** (without ORM). |
| `GET` | `/api/pigments` | `h.GetPigmentsAPI` | **[Lab 3]** REST API: Returns JSON list of all pigments with author info. |
| `GET` | `/api/pigments/{id}` | `h.GetPigmentByIDAPI` | **[Lab 3]** REST API: Returns single pigment JSON using `FullPigmentSerializer` (nested author JSON). |
| `POST` | `/api/pigments` | `h.AddPigmentAPI` | **[Lab 3]** REST API: Creates new pigment from JSON or multipart form + optional image upload to MinIO. |
| `PUT` | `/api/pigments/{id}` | `h.ModifyPigmentAPI` | **[Lab 3]** REST API: Modifies existing pigment by ID + optional image upload to MinIO. |
| `DELETE` | `/api/pigments/{id}` | `h.DeletePigmentAPI` | **[Lab 3]** REST API: Hard deletes pigment by ID. |
| `GET` | `/api/users/{id}/pigments` | `h.GetUserPigmentsAPI` | **[Lab 3]** REST API: Returns user details and an array of their created pigments (Nested Serialization). |

---

## 6. Key Technical Decisions & Lab Requirements

1. **Standard Library Routing**: Chosen `net/http` (`http.ServeMux`) and standard `html/template` over framework wrappers (e.g., Gin) for lightness and standard Go idioms.
2. **Hybrid DB Approach (GORM + Raw SQL)**:
   - GORM is used for model queries, filtering (`Where`), and updates (`Updates`).
   - Raw SQL query (`r.db.DB().Exec("UPDATE pigments SET status = $1 WHERE id = $2", "удалён", id)`) is intentionally used in `LogicalDelete` to satisfy specific lab requirements (logical delete without ORM).
3. **MinIO Integration & SDK**:
   - Integrated MinIO Go SDK (`github.com/minio/minio-go/v7`) in Repository (`AddOrReplacePigmentImage`).
   - Same-origin reverse proxy `/minio/` in `main.go` for seamless CORS-free media delivery.
4. **Lab 3 REST API & Nested Serialization**:
   - Models annotated with `json:"..."` tags (`User`, `Pigment`, `UserPigments`, `FullPigmentSerializer`).
   - Preloading relationships in GORM (`r.db.Preload("Creator")`).

---

## 7. Current Project Status & Completed Features

- [x] **Lab 3 REST API Implementation**:
  - Full REST API endpoints returning JSON (`/api/pigments`, `/api/pigments/{id}`, `/api/users/{id}/pigments`).
  - MinIO SDK client initialization and file upload handler (`AddOrReplacePigmentImage`).
  - Nested JSON serialization (`FullPigmentSerializer` and `UserPigments`).
  - Seeding default user in `cmd/migrate/main.go`.
- [x] **Feed View (Reconstruction)**:
  - Video/image media background rendering fixed (rendered conditionally via `{{ if .Pigment.VideoURL }}` so `<img>` element doesn't overlay and block MinIO video playback).
  - Collapsible short description with hidden checkbox and styled "БОЛЬШЕ" / "МЕНЬШЕ" toggle button (`.desc-label::after`).
  - Removed "ЦВЕТ" parameter label from right side panel.
  - Live Likes count display (`GetLikesCount`).
  - Next item vertical navigation button (`next=true`).
- [x] **Draft & Analysis View (Addressing Lab 2 Feedback)**:
  - Vertical media previews stacked ("друг под другом") for photo and video at the top of the form (for both Step 1 & Step 2).
  - Preset `<select>` dropdown lists with default local media (`/assets/images/placeholder.svg` & `/assets/videos/placeholder.mp4`) placed at the top (selected by default on first load before clicking "Далее").
  - Options list populated with exact MinIO bucket contents (`cinnabar.jpg`/`.mp4`, `lapis-lazuli.jpg`/`.mp4`, `malachite.jpg`/`.mp4`, `ultramarine.jpg`/`.mp4`, `umber.jpg`/`.mp4`, `yellow-ochre.jpg`/`.mp4`) and local fallback static assets.
  - Instant client-side JS preview update when selecting options in Step 1.
  - Step 2 publication form with preview, description, reflection (R), absorption (A).
- [x] **Local Assets Fallback (No Server Uploads)**:
  - Local static media assets placed in `templates/assets/images/` and `templates/assets/videos/`.
  - MinIO remains primary, with automatic local fallback if MinIO is offline or image fails to load (`onerror` fallback handling).
- [x] **Grid/Tile View (Catalog)**:
  - Responsive card grid.
  - Absorption range slider filter (`absorption <= val`).
  - Delete button styled uniformly with theme (`.card-delete-btn` in dark/gold theme instead of raw red style).
- [x] **Logical Deletion**:
  - Soft delete updating status to `"удалён"` via SQL query.

---

## 8. Known Issues & Tech Debt

1. **Hardcoded User Authentication**:
   - `CreatorID` is hardcoded to `1` in `GetAnalysis` and `CreateDraft`.
   - No login/registration flow or JWT/session middleware implemented yet.
2. **Missing Like Interactivity**:
   - `GetLikesCount` counts records in `likes` table, but there is no handler or UI button to add/remove a like.
3. **MinIO Direct Upload**:
   - Form currently requires manual key entry (e.g. `cinnabar.jpg`). File upload directly to MinIO via multipart form is not yet connected.
4. **Error Handling & Validation**:
   - Form inputs lack validation for empty strings or invalid float ranges.
   - Database errors in handlers are largely ignored (`_`) instead of returning 400/500 HTTP status pages.

---

## 9. Next Steps / TODOs for Future Development

1. **User Authentication & Authorization**:
   - Implement login/register pages or endpoints.
   - Store user identity in JWT/Session and replace hardcoded `userID = 1`.
2. **Like System**:
   - Add a `POST /like-pigment` handler and toggle button on the feed and tile views.
3. **MinIO Direct File Upload**:
   - Integrate MinIO SDK in Go backend to upload media files directly via the analysis form.
4. **Enhanced Error Handling & Validation**:
   - Return proper user feedback/toast messages on invalid inputs or DB errors.
