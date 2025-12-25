package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const baseURL = "http://localhost:8080/api/v1"

func main() {
	fmt.Println("Starting Integration Tests...")

	// 1. Health Check
	assertStatus("GET", "/health", nil, 200)

	// 2. Signup
	email := fmt.Sprintf("test-%d@example.com", time.Now().Unix())
	password := "password123"
	signupPayload := map[string]string{"email": email, "password": password}
	assertStatus("POST", "/auth/signup", signupPayload, 201)

	// 3. Login
	loginPayload := map[string]string{"email": email, "password": password}
	resp := request("POST", "/auth/login", loginPayload)
	if resp.StatusCode != 200 {
		fatal("Login failed")
	}
	var authResp struct {
		AccessToken string `json:"access_token"`
	}
	decodeJSON(resp, &authResp)
	token := authResp.AccessToken
	fmt.Println("Login successful, token retrieved.")

	// 4. Create Trip
	tripPayload := map[string]interface{}{
		"location":   "Paris",
		"start_date": time.Now(),
		"end_date":   time.Now().Add(24 * time.Hour),
	}
	resp = requestWithAuth("POST", "/trips", tripPayload, token)
	if resp.StatusCode != 201 {
		fatal(fmt.Sprintf("Create Trip failed: %d", resp.StatusCode))
	}
	var tripResp struct {
		ID string `json:"id"`
	}
	decodeJSON(resp, &tripResp)
	tripID := tripResp.ID
	fmt.Printf("Trip created: %s\n", tripID)

	// 5. Create Itinerary
	itinPayload := map[string]string{"name": "Day 1"}
	resp = requestWithAuth("POST", fmt.Sprintf("/trips/%s/itineraries", tripID), itinPayload, token)
	if resp.StatusCode != 201 {
		fatal("Create Itinerary failed")
	}
	var itinResp struct {
		ID string `json:"id"`
	}
	decodeJSON(resp, &itinResp)
	itinID := itinResp.ID
	fmt.Printf("Itinerary created: %s\n", itinID)

	// 6. Create Activity
	actPayload := map[string]interface{}{
		"name":       "Eiffel Tower",
		"type":       "sightseeing",
		"start_time": time.Now(),
		"end_time":   time.Now().Add(time.Hour),
		"doba":       100.0,
	}
	resp = requestWithAuth("POST", fmt.Sprintf("/itineraries/%s/activities", itinID), actPayload, token)
	if resp.StatusCode != 201 {
		fatal("Create Activity failed")
	}
	var actResp struct {
		ID string `json:"id"`
	}
	decodeJSON(resp, &actResp)
	actID := actResp.ID
	fmt.Printf("Activity created: %s\n", actID)

	// 7. Upload Media (Multipart)
	uploadMedia(token, actID)

	// 8. Share Trip
	resp = requestWithAuth("POST", fmt.Sprintf("/trips/%s/share", tripID), nil, token)
	if resp.StatusCode != 200 {
		fatal("Share Trip failed")
	}
	var shareResp struct {
		Token string `json:"share_token"`
	}
	decodeJSON(resp, &shareResp)
	shareToken := shareResp.Token
	fmt.Printf("Share token: %s\n", shareToken)

	// 9. Public Preview
	resp = request("GET", fmt.Sprintf("/preview/%s", shareToken), nil)
	if resp.StatusCode != 200 {
		fatal(fmt.Sprintf("Preview failed: %d", resp.StatusCode))
	}
	fmt.Println("Preview accessed successfully.")

	// 10. Users Device Token
	devicePayload := map[string]string{"token": "fcm-fake-token"}
	assertStatusWithAuth("POST", "/users/device-token", devicePayload, token, 200)

	fmt.Println("ALL TESTS PASSED!")
}

func request(method, path string, payload interface{}) *http.Response {
	return doRequest(method, path, payload, "")
}

func requestWithAuth(method, path string, payload interface{}, token string) *http.Response {
	return doRequest(method, path, payload, token)
}

func doRequest(method, path string, payload interface{}, token string) *http.Response {
	var body io.Reader
	if payload != nil {
		b, _ := json.Marshal(payload)
		body = bytes.NewBuffer(b)
	}
	req, _ := http.NewRequest(method, baseURL+path, body)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fatal(err.Error())
	}
	return resp
}

func assertStatus(method, path string, payload interface{}, status int) {
	resp := request(method, path, payload)
	if resp.StatusCode != status {
		fatal(fmt.Sprintf("%s %s expected %d got %d", method, path, status, resp.StatusCode))
	}
	fmt.Printf("PASS: %s %s\n", method, path)
}

func assertStatusWithAuth(method, path string, payload interface{}, token string, status int) {
	resp := requestWithAuth(method, path, payload, token)
	if resp.StatusCode != status {
		fatal(fmt.Sprintf("%s %s expected %d got %d", method, path, status, resp.StatusCode))
	}
	fmt.Printf("PASS: %s %s\n", method, path)
}

func uploadMedia(token, activityID string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("dummy content"))

	writer.WriteField("activity_id", activityID)
	writer.Close()

	req, _ := http.NewRequest("POST", baseURL+"/media/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fatal(err.Error())
	}
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		fatal(fmt.Sprintf("Upload Media failed: %d %s", resp.StatusCode, string(b)))
	}
	fmt.Println("PASS: Upload Media")
}

func decodeJSON(resp *http.Response, target interface{}) {
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(target)
}

func fatal(msg string) {
	fmt.Println("FAIL:", msg)
	os.Exit(1)
}