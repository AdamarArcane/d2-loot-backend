package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Struct to hold the manifest response with component paths
type ManifestResponse struct {
	Response struct {
		JsonWorldComponentContentPaths struct {
			English map[string]string `json:"en"`
		} `json:"jsonWorldComponentContentPaths"`
	} `json:"Response"`
}

// Structs for item definitions and plug definitions
type ItemDefinition struct {
	DisplayProperties struct {
		Name string `json:"name"`
	} `json:"displayProperties"`
	ItemType     int     `json:"itemType"`
	DefaultPerks []int64 `json:"defaultPerks"`
}

type PlugSetDefinition struct {
	DisplayProperties struct {
		Name string `json:"name"`
	} `json:"displayProperties"`
}

var items map[string]ItemDefinition
var perks map[string]PlugSetDefinition

// ManageManifest handles downloading and parsing the manifest
func ManageManifest(client *http.Client, apiKey string) error {
	// Step 1: Download the manifest metadata
	manifestMetadata, err := downloadManifestMetadata(client, apiKey)
	if err != nil {
		return fmt.Errorf("failed to download manifest metadata: %w", err)
	}

	// Step 2: Retrieve the relevant manifest content URLs
	itemManifestPath, ok := manifestMetadata.Response.JsonWorldComponentContentPaths.English["DestinyInventoryItemDefinition"]
	if !ok {
		return fmt.Errorf("item definition URL not found in the manifest metadata")
	}

	plugManifestPath, ok := manifestMetadata.Response.JsonWorldComponentContentPaths.English["DestinyPlugSetDefinition"]
	if !ok {
		return fmt.Errorf("plug set definition URL not found in the manifest metadata")
	}

	// Step 3: Download the manifest content (JSON)
	itemManifestURL := "https://www.bungie.net" + itemManifestPath
	plugManifestURL := "https://www.bungie.net" + plugManifestPath

	log.Printf("Downloading item manifest from: %s\n", itemManifestURL)
	log.Printf("Downloading plug manifest from: %s\n", plugManifestURL)

	outputItemFile := "DestinyInventoryItemDefinition.json"
	outputPlugFile := "DestinyPlugSetDefinition.json"

	// Download item manifest content
	err = downloadManifestContent(client, itemManifestURL, outputItemFile)
	if err != nil {
		return fmt.Errorf("failed to download item manifest content: %w", err)
	}

	// Download plug manifest content
	err = downloadManifestContent(client, plugManifestURL, outputPlugFile)
	if err != nil {
		return fmt.Errorf("failed to download plug manifest content: %w", err)
	}

	// Step 4: Load and parse the JSON files
	items, err = loadItemManifestContent(outputItemFile)
	if err != nil {
		return fmt.Errorf("failed to load item manifest content: %w", err)
	}

	perks, err = loadPlugManifestContent(outputPlugFile)
	if err != nil {
		return fmt.Errorf("failed to load plug manifest content: %w", err)
	}

	log.Println("Manifest loaded successfully.")
	return nil
}

// GetItemInfo retrieves an item's name and perks by its hash from the loaded manifest
func GetItemInfo(itemHash string) (string, error) {
	if item, found := items[itemHash]; found {
		if item.ItemType == 3 { // Assuming itemType 3 is a weapon
			// Get item name
			itemName := item.DisplayProperties.Name
			return itemName, nil
		}
		return "", fmt.Errorf("item with hash %s is not a weapon", itemHash)
	}
	return "", fmt.Errorf("item with hash %s not found", itemHash)
}

// GetPerkName retrieves a perk's name by its hash from the loaded manifest
func GetPerkName(perkHash string) (string, error) {
	if perk, found := perks[perkHash]; found {
		return perk.DisplayProperties.Name, nil
	}
	return "", fmt.Errorf("perk with hash %s not found", perkHash)
}

// Step 1: Download Manifest Metadata
func downloadManifestMetadata(client *http.Client, apiKey string) (*ManifestResponse, error) {
	url := "https://www.bungie.net/Platform/Destiny2/Manifest/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create manifest metadata request: %w", err)
	}
	req.Header.Set("X-API-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest metadata response: %w", err)
	}

	var manifestMetadata ManifestResponse
	err = json.Unmarshal(body, &manifestMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to parse manifest metadata: %w", err)
	}

	return &manifestMetadata, nil
}

// Step 3: Download Manifest Content (JSON)
func downloadManifestContent(client *http.Client, manifestURL, outputFile string) error {
	req, err := http.NewRequest("GET", manifestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create manifest content request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download manifest content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write manifest content to file: %w", err)
	}

	log.Printf("Manifest content saved to %s\n", outputFile)
	return nil
}

// Step 4: Load Manifest Content and Parse JSON

// Load Item Manifest Content and Parse JSON
func loadItemManifestContent(filePath string) (map[string]ItemDefinition, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open manifest file: %w", err)
	}
	defer file.Close()

	var items map[string]ItemDefinition
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&items); err != nil {
		return nil, fmt.Errorf("failed to parse manifest JSON: %w", err)
	}

	return items, nil
}

// Load Plug Manifest Content and Parse JSON
func loadPlugManifestContent(filePath string) (map[string]PlugSetDefinition, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open manifest file: %w", err)
	}
	defer file.Close()

	var plugs map[string]PlugSetDefinition
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&plugs); err != nil {
		return nil, fmt.Errorf("failed to parse manifest JSON: %w", err)
	}

	return plugs, nil
}
