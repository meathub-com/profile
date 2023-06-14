package transport

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Offer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Item  string `json:"item"`
	Price int    `json:"price"`
}

func (h *Handler) TestOfferService(w http.ResponseWriter, r *http.Request) {
	getEndpoints := []string{
		"http://offers:8082/offers/1", // test GetOffer
		"http://offers:8082/offers",   // test GetOffers
	}

	postEndpoints := map[string]string{
		"http://offers:8082/offers": `{
			"offerName": "Organic Free Range Chicken",
			"item": "fewfeff",
			"price": 22
		}`,
	}

	for _, endpoint := range getEndpoints {
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			log.Errorf("Error creating request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Errorf("Error making request to %s: %v", endpoint, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Errorf("Error: status code for %s is %d", endpoint, resp.StatusCode)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Infof("Successfully connected to %s with status code: %d", endpoint, resp.StatusCode)
	}

	for endpoint, body := range postEndpoints {
		req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(body))
		if err != nil {
			log.Errorf("Error creating request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Errorf("Error making request to %s: %v", endpoint, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Errorf("Error: status code for %s is %d", endpoint, resp.StatusCode)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Infof("Successfully connected to %s with status code: %d", endpoint, resp.StatusCode)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All endpoints returned status code 200"))
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
	url := "http://offers:8082/offers"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
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

	// Read the response
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

type OfferResponse struct {
	ID    string `json:"id"`
	Name  string `json:"offerName"`
	Item  string `json:"item"`
	Price int    `json:"price"`
}
