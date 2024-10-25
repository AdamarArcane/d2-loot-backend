package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (api *apiConfig) getMembershipData(client *http.Client) (string, int, error) {
	req, err := http.NewRequest("GET", "https://www.bungie.net/Platform/User/GetMembershipsForCurrentUser/", nil)
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("X-API-Key", api.API_KEY)

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Response struct {
			DestinyMemberships []struct {
				MembershipID   string `json:"membershipId"`
				MembershipType int    `json:"membershipType"`
			} `json:"destinyMemberships"`
		} `json:"Response"`
		ErrorCode   int    `json:"ErrorCode"`
		ErrorStatus string `json:"ErrorStatus"`
		Message     string `json:"Message"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", 0, err
	}

	if result.ErrorCode != 1 {
		return "", 0, fmt.Errorf("API error: %s", result.Message)
	}

	if len(result.Response.DestinyMemberships) == 0 {
		return "", 0, fmt.Errorf("no Destiny memberships found")
	}

	// Use the first membership
	return result.Response.DestinyMemberships[0].MembershipID, result.Response.DestinyMemberships[0].MembershipType, nil
}

func (api *apiConfig) getPlayerProfile(client *http.Client, membershipType int, membershipID string) (*ProfileData, error) {
	url := fmt.Sprintf("https://www.bungie.net/Platform/Destiny2/%d/Profile/%s/?components=100,102,103,200,201,205,300,305", membershipType, membershipID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", api.API_KEY)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response into ProfileResponse
	var profile ProfileData
	err = json.Unmarshal(body, &profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}
