package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aheld/listservice/db"
	"github.com/aheld/listservice/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	database, err := db.Initialize(os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		log.Fatalf("Could not connect to DB %v", err)
	}
	defer database.Conn.Close()

	listService := ListService{
		listRepo: database,
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	r.Get("/startup", func(w http.ResponseWriter, r *http.Request) {
		msg, err := database.CheckVersion()
		if err != nil {
			render.Render(w, r, ErrSrvUnavailable(err))
			return
		}
		w.Write([]byte(msg))
	})

	r.Route("/lists", func(r chi.Router) {
		r.Post("/", listService.CreateListItem)
		r.Route("/{userId}", func(r chi.Router) {
			r.Use(ListCtx)
			r.Get("/", listService.GetListItems)
			r.Put("/{itemId}", listService.UpdateListItem)
		})
	})

	http.ListenAndServe("localhost:3333", r)
}

type ListService struct {
	listRepo domain.ListRepo
}

type ListItemRequest struct {
	*domain.ListItem
}

func (a *ListItemRequest) Bind(r *http.Request) error {
	if a.ListItem == nil {
		return errors.New("missing required ListItem fields")
	}
	return nil
}

func (s *ListService) CreateListItem(w http.ResponseWriter, r *http.Request) {
	data := &ListItemRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	listItem := data.ListItem

	s.listRepo.InsertListItem(data.UserId, data.Item)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, listItem)
}

func (s *ListService) GetListItems(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIDKey).(int)

	items, err := s.listRepo.GetListItems(userId)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &items)
}

func (s *ListService) UpdateListItem(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIDKey).(int)

	data := &ListItemRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	log.Printf("posted to udpate %v", data)

	if itemId := chi.URLParam(r, "itemId"); itemId != "" {
		itemIdInt, err := strconv.Atoi(itemId)
		if err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}
		err = s.listRepo.UpdateListItem(userId, itemIdInt, data.Item)
		if err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, struct {
		message string
	}{
		message: "Stored Proc handled your update",
	})
}

func ListCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if userId := chi.URLParam(r, "userId"); userId != "" {
			userIdInt, err := strconv.Atoi(userId)
			if err == nil { //this is going to confuse someone
				ctx := context.WithValue(r.Context(), userIDKey, userIdInt)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
		render.Render(w, r, ErrInvalidRequest(errors.New("missing required user id, so keep hacking")))
	})
}

type contextKey string

const userIDKey contextKey = "userid"

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func ErrSrvUnavailable(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 503,
		StatusText:     "Server Unavailable, try again later please.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
