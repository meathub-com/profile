package transport

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Offer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Item  string `json:"item"`
	Price int    `json:"price"`
}

func (h *Handler) CreateOffer(w http.ResponseWriter, r *http.Request) {
	profileId := chi.URLParam(r, "id")

	foundProfile, err := h.Service.GetProfile(r.Context(), profileId)
	if err != nil {
		log.Errorf("Error getting profile: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Infof("Found profile: %v", foundProfile)

	// Create a new HTTP request
	url := "http://localhost:8082/offers"
	req, err := http.NewRequest("POST", url, r.Body)
	if err != nil {
		log.Errorf("Error creating request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error making request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		w.WriteHeader(http.StatusOK)
		var offerResponse OfferResponse
		err := json.NewDecoder(resp.Body).Decode(&offerResponse)
		if err != nil {
			log.Errorf("Error decoding offer: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		offerResponseJson, err := json.Marshal(offerResponse)
		if err != nil {
			log.Errorf("Error encoding offer response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(offerResponseJson)

	default:
		w.WriteHeader(http.StatusBadRequest)
	}

}

type OfferResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Item  string `json:"item"`
	Price int    `json:"price"`
}
