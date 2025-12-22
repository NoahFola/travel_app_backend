package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/NoahFola/travel_app_backend/internal/repository"
)

type LocationService struct {
	Repo *repository.LocationRepository
}

type GooglePlaceResult struct {
	PlaceID          string `json:"place_id"`
	Name             string `json:"name"`
	FormattedAddress string `json:"formatted_address"`
	Geometry         struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
	} `json:"geometry"`
}

type GooglePlacesResponse struct {
	Results []GooglePlaceResult `json:"results"`
	Status  string              `json:"status"`
}

// SearchPlaces proxies the request to Google Places API
func (s *LocationService) SearchPlaces(query string) ([]GooglePlaceResult, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		return nil, errors.New("GOOGLE_MAPS_API_KEY is not set")
	}

	baseURL := "https://maps.googleapis.com/maps/api/place/textsearch/json"
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("query", query)
	q.Set("key", apiKey)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google places api returned status: %d", resp.StatusCode)
	}

	var placesResp GooglePlacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&placesResp); err != nil {
		return nil, err
	}

	if placesResp.Status != "OK" && placesResp.Status != "ZERO_RESULTS" {
		return nil, fmt.Errorf("google places api error: %s", placesResp.Status)
	}

	return placesResp.Results, nil
}

// GetOrCreateLocation checks if a location exists by PlaceID, otherwise creates it from provided data
func (s *LocationService) GetOrCreateLocation(ctx context.Context, placeData GooglePlaceResult) (*repository.Location, error) {
	// 1. Check if exists
	existing, err := s.Repo.GetByPlaceID(ctx, placeData.PlaceID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	// 2. Create new
	newLoc := &repository.Location{
		Name:          placeData.Name,
		Address:       placeData.FormattedAddress,
		Latitude:      placeData.Geometry.Location.Lat,
		Longitude:     placeData.Geometry.Location.Lng,
		GooglePlaceID: &placeData.PlaceID,
	}

	if err := s.Repo.Create(ctx, newLoc); err != nil {
		return nil, err
	}

	return newLoc, nil
}
