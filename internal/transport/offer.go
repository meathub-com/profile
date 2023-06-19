package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

const ()
const (
	ContentType = "Content-Type"
	JsonType    = "application/json"
)

const (
	OfferServiceUrl = "http://offers:8082/offers"
)

type OfferResponse struct {
	ID    string `json:"id"`
	Name  string `json:"offerName"`
	Item  string `json:"item"`
	Price int    `json:"price"`
}

func (h *Handler) CreateOfferV2(w http.ResponseWriter, r *http.Request) {
	//Get the profile id from the URL
	profileId := chi.URLParam(r, "id")
	if profileId == "" {
		log.Errorf("Error getting profile id")
		writeError(w, http.StatusBadRequest, "Error getting profile id", nil)
		return
	}
	//Find profile by id
	foundProfile, err := h.Service.GetProfile(r.Context(), profileId)
	if err != nil {
		log.Errorf("Error getting profile: %v", err)
		writeError(w, http.StatusNotFound, fmt.Sprintf("Error fetching profile with id %v", profileId), err)
		return
	}
	log.Infof("Found profile: %v", foundProfile)

	//Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error reading request body: %v", err)
		writeError(w, http.StatusBadRequest, "Error reading request body", err)
		return
	}
	//Create a new offer by company profile
	url := OfferServiceUrl
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Error creating request", err)
		return
	}
	req.Header.Set(ContentType, JsonType)

	resp, err := h.Client.Do(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error creating offer", err)
		return
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error reading response body", err)
		return
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var offerResponse OfferResponse
		err := json.Unmarshal(response, &offerResponse)
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
		writeResponse(w, http.StatusOK, string(offerResponseJson))
	default:
		errorMessage := "Unexpected status code: " + strconv.Itoa(resp.StatusCode)
		writeError(w, http.StatusInternalServerError, errorMessage, nil)
	}

}
func writeError(w http.ResponseWriter, httpStatus int, message string, err error) {
	log.Errorf("%s: %v", message, err)
	http.Error(w, fmt.Sprintf("%s: %v", message, err), httpStatus)
}
func writeResponse(w http.ResponseWriter, httpStatus int, message string) {
	w.WriteHeader(httpStatus)
	write, err := w.Write([]byte(message))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error writing response", err)
	}
	log.Infof("Wrote %v bytes", write)
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

	// Read the body content
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error reading request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create a new HTTP request
	url := OfferServiceUrl
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Errorf("Error creating request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Header.Set(ContentType, JsonType)

	resp, err := h.Client.Do(req)
	if err != nil {
		log.Errorf("Error making request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var offerResponse OfferResponse
		err := json.Unmarshal(response, &offerResponse)
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
		errorMessage := "Unexpected status code: " + strconv.Itoa(resp.StatusCode)
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte(errorMessage))
	}
}
