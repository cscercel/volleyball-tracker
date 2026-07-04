package handler

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cscercel/volleyball-tracker/internal/service"
	webassets "github.com/cscercel/volleyball-tracker/web"
)

type PageHandler struct {
	userService 	*service.UserService
	playerService	*service.PlayerService
	jwtSecret 		string
	secureCookies	bool // http in dev vs https in prod
}

func NewPageHandler(
	userService *service.UserService, 
	playerService *service.PlayerService,
	jwtSecret string, 
	secureCookies bool,
) *PageHandler {
	return &PageHandler{
		userService: userService, 
		playerService: playerService,
		jwtSecret: jwtSecret, 
		secureCookies: secureCookies,
	}
}

func (h *PageHandler) RegisterRoutes(r chi.Router) {
	staticSub, err := fs.Sub(webassets.StaticFS, "static")
	if err != nil {
		panic(err) // crash at startup if embed path is wrong
	}

	fileServer := http.FileServer(http.FS(staticSub))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Public page routes
	r.Get("/login", h.handleLoginPage)
	r.Post("/login", h.handleLoginSubmit)
	r.Post("/logout", h.handleLogout)
	r.Get("/register", h.handleRegisterPage)
	r.Post("/register", h.handleRegisterSubmit)

	// Protected page routes
	r.Group(func(r chi.Router) {
		r.Use(PageAuthMiddleware(h.jwtSecret))
		r.Get("/", h.handleHomePage)
		r.Get("/players", h.handlePlayersPage)
		r.Post("/players/add", h.handleAddPlayerSubmit)
		r.Post("/players/rename", h.handleRenamePlayerSubmit)
		r.Post("/players/delete", h.handleDeletePlayerSubmit)
	})
}
