package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"assignment4/repository"

	"github.com/google/uuid"
)

type UserHandler struct {
	repo *repository.Repository
}

func NewUserHandler(r *repository.Repository) *UserHandler {
	return &UserHandler{repo: r}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	page := 1
	pageSize := 10
	var err error
	if query.Get("page") != "" {
		page, err = strconv.Atoi(query.Get("page"))
		if err != nil || page < 1 {
			http.Error(w, "invalid page", http.StatusBadRequest)
			return
		}
	}
	if query.Get("page_size") != "" {
		pageSize, err = strconv.Atoi(query.Get("page_size"))
		if err != nil || pageSize < 1 {
			http.Error(w, "invalid page_size", http.StatusBadRequest)
			return
		}
	}

	filters := map[string]string{}
	for _, key := range []string{"id", "name", "email", "gender", "birth_date"} {
		if v := query.Get(key); v != "" {
			filters[key] = v
		}
	}

	orderBy := query.Get("order_by")

	resp, err := h.repo.GetPaginatedUsers(page, pageSize, filters, orderBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UserHandler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	u1 := query.Get("user1")
	u2 := query.Get("user2")
	if u1 == "" || u2 == "" {
		http.Error(w, "missing user1 or user2", http.StatusBadRequest)
		return
	}
	userID1, err := uuid.Parse(u1)
	if err != nil {
		http.Error(w, "invalid user1", http.StatusBadRequest)
		return
	}
	userID2, err := uuid.Parse(u2)
	if err != nil {
		http.Error(w, "invalid user2", http.StatusBadRequest)
		return
	}

	users, err := h.repo.GetCommonFriends(userID1, userID2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
