package handler

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cscercel/volleyball-tracker/internal/service"
	webassets "github.com/cscercel/volleyball-tracker/web"
)

type PageHandler struct {
	userService   *service.UserService
	playerService *service.PlayerService
	matchService  *service.MatchService
	jwtSecret     string
	secureCookies bool
}

func NewPageHandler(
	userService *service.UserService,
	playerService *service.PlayerService,
	matchService *service.MatchService,
	jwtSecret string,
	secureCookies bool,
) *PageHandler {
	return &PageHandler{
		userService:   userService,
		playerService: playerService,
		matchService:  matchService,
		jwtSecret:     jwtSecret,
		secureCookies: secureCookies,
	}
}

func (h *PageHandler) RegisterRoutes(r chi.Router) {
	staticSub, err := fs.Sub(webassets.StaticFS, "static")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(staticSub))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Public page routes
	r.Get("/", h.handleHomePage)
	r.Get("/login", h.handleLoginPage)
	r.Post("/login", h.handleLoginSubmit)
	r.Post("/logout", h.handleLogout)
	r.Get("/register", h.handleRegisterPage)
	r.Post("/register", h.handleRegisterSubmit)
	r.Get("/players", h.handlePlayersPage)
	r.Get("/matches", h.handleMatchesPage)

	// Protected: only actions that mutate data require login.
	r.Group(func(r chi.Router) {
		r.Use(PageAuthMiddleware(h.jwtSecret))
		r.Post("/players/add", h.handleAddPlayerSubmit)
		r.Post("/players/rename", h.handleRenamePlayerSubmit)
		r.Post("/players/delete", h.handleDeletePlayerSubmit)

		r.Post("/matches/create/update", h.handleCreateTeamUpdate)
		r.Post("/matches/create/submit", h.handleCreateMatchSubmit)
		r.Post("/matches/drafts/submit", h.handleDraftSubmit)
		r.Post("/matches/drafts/delete", h.handleDraftDelete)
	})
}
