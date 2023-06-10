package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"profile/internal/profile"
)

type Service interface {
	GetProfile(ctx context.Context, id string) (profile.Profile, error)
	PostProfile(ctx context.Context, profile profile.Profile) (profile.Profile, error)
	UpdateProfile(ctx context.Context, profile profile.Profile) (profile.Profile, error)
	DeleteProfile(ctx context.Context, id string) error
	GetProfiles(ctx context.Context) ([]profile.Profile, error)
}

func (h *Handler) GetProfileByUser(w http.ResponseWriter, r *http.Request) {
	profile, err := h.Service.GetProfile(r.Context(), r.Header.Get("user_id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		log.Errorf("Error getting profile: %v", err)
	}

}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	profileIDStr := chi.URLParam(r, "id")
	p, err := h.Service.GetProfile(r.Context(), profileIDStr)
	if errors.Is(err, profile.ErrFetchingProfile) {
		log.Errorf("Error getting profile: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Errorf("Error getting profile: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Errorf("Error encoding profile: %v", err)
	}
}

// CreateProfile godoc
// @Summary Create a new profile
// @Description Create a new profile
// @Tags profiles
// @Accept  json
// @Produce  json
// @Param profile body profile.Profile true "Profile info"
// @Success 200 {object} profile.Profile
// @Router /profiles [post]
func (h *Handler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var p profile.Profile
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Errorf("Error decoding profile: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p.UserId = r.Context().Value("user_id").(string)
	savedProfile, err := h.Service.PostProfile(r.Context(), p)
	if err != nil {
		log.Errorf("Error posting profile: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(savedProfile); err != nil {
		log.Errorf("Error getting profile: %v", err)
	}
}

func (h *Handler) GetProfiles(w http.ResponseWriter, r *http.Request) {
	profile, err := h.Service.GetProfiles(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		log.Errorf("Error getting profile: %v", err)
	}
}
