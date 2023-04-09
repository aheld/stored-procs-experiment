package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/aheld/listservice/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (s *Server) MountHandlers() {
	listService := ListService{
		listDb: s.Database,
	}

	s.Router.Route("/lists", func(r chi.Router) {
		r.Use(BannerCtx)
		r.Post("/", listService.CreateListItem)
		r.Route("/{userId}", func(r chi.Router) {
			r.Use(ListCtx)
			r.Get("/", listService.GetListItems)
			r.Put("/{itemId}", listService.UpdateListItem)
		})
	})
}

type ListService struct {
	listDb domain.ListDb
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
	bannerId := r.Context().Value(bannerIDKey).(string)

	data := &ListItemRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	listItem := data.ListItem

	s.listDb.InsertListItem(bannerId, data.UserId, data.Item)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, listItem)
}

func (s *ListService) GetListItems(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIDKey).(int)
	bannerId := r.Context().Value(bannerIDKey).(string)

	items, err := s.listDb.GetListItems(bannerId, userId)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &items)
}

func (s *ListService) UpdateListItem(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIDKey).(int)
	bannerId := r.Context().Value(bannerIDKey).(string)

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
		err = s.listDb.UpdateListItem(bannerId, userId, itemIdInt, data.Item)
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

func BannerCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), bannerIDKey, "f4bd6cdc-eb4b-4f74-8565-c243d3fdf20x")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type contextKey string

const userIDKey contextKey = "userid"
const bannerIDKey contextKey = "userid"
