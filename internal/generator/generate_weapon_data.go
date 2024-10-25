package generator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type WeaponPerkInput struct {
	WeaponName   string   `json:"weaponName"`
	DesiredPerks []string `json:"desiredPerks"`
	Description  string   `json:"description"`
	Source       string   `json:"source"`
	Bucket       string   `json:"bucket"`
	Rank         string   `json:"rank"`
}

type ItemDefinition struct {
	DisplayProperties struct {
		Name        string `json:"name"`
		Icon        string `json:"icon"`
		Description string `json:"description"`
	} `json:"displayProperties"`
	ItemTypeDisplayName string `json:"itemTypeDisplayName"`
	ItemSubType         int    `json:"itemSubType"`
	ItemType            int    `json:"itemType"`
	Hash                int64  `json:"hash"`
	Redacted            bool   `json:"redacted"`
	Sockets             struct {
		SocketEntries []struct {
			SingleInitialItemHash int64 `json:"singleInitialItemHash"`
			ReusablePlugSetHash   int64 `json:"reusablePlugSetHash"`
			RandomizedPlugSetHash int64 `json:"randomizedPlugSetHash"`
		} `json:"socketEntries"`
	} `json:"sockets"`
}

type PlugSetDefinition struct {
	ReusablePlugItems []struct {
		PlugItemHash int64 `json:"plugItemHash"`
	} `json:"reusablePlugItems"`
}

func GenerateWeaponData() error {
	// Paths to files
	itemDefPath := filepath.Join("..", "..", "DestinyInventoryItemDefinition.json")
	plugSetDefPath := filepath.Join("..", "..", "DestinyPlugSetDefinition.json")
	weaponsPerksFilePath := filepath.Join("..", "..", "cmd", "generate_constants", "weapons_and_perks.json")
	outputPath := filepath.Join("..", "constants", "weapon_data.go")

	// Read the item definitions JSON file
	itemDefinitions, err := readItemDefinitions(itemDefPath)
	if err != nil {
		return fmt.Errorf("error reading item definitions: %v", err)
	}

	// Read the plug set definitions JSON file
	plugSetDefinitions, err := readPlugSetDefinitions(plugSetDefPath)
	if err != nil {
		return fmt.Errorf("error reading plug set definitions: %v", err)
	}

	// Read the weapons and perks input file
	weaponInputs, err := readWeaponPerkInputs(weaponsPerksFilePath)
	if err != nil {
		return fmt.Errorf("error reading weapons and perks input: %v", err)
	}

	// Find item hashes for the weapons (including all versions)
	weaponHashesMap, weaponDefinitions, err := findWeaponHashes(itemDefinitions, weaponInputs)
	if err != nil {
		return fmt.Errorf("error finding weapon hashes: %v", err)
	}

	// Collect all desired perk names
	desiredPerkNames := []string{}
	for _, input := range weaponInputs {
		desiredPerkNames = append(desiredPerkNames, input.DesiredPerks...)
	}

	// Find perk hashes and build weapon possible perks map
	perkHashesMap, perkHashesReverseMap, weaponPossiblePerksMap, err := findWeaponPossiblePerkHashes(weaponDefinitions, itemDefinitions, plugSetDefinitions, desiredPerkNames)
	if err != nil {
		return fmt.Errorf("error finding weapon possible perk hashes: %v", err)
	}

	// Generate weapon_data.go file
	err = generateWeaponDataFile(weaponInputs, weaponHashesMap, perkHashesMap, perkHashesReverseMap, weaponPossiblePerksMap, weaponDefinitions, outputPath)
	if err != nil {
		return fmt.Errorf("error generating weapon data file: %v", err)
	}

	return nil
}

// Helper functions

func readItemDefinitions(filePath string) (map[int64]ItemDefinition, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Remove the outer wrapper (Bungie's JSON files sometimes have additional metadata)
	raw := json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	itemDefs := make(map[int64]ItemDefinition)
	if err := json.Unmarshal(raw, &itemDefs); err != nil {
		return nil, err
	}

	return itemDefs, nil
}

func readPlugSetDefinitions(filePath string) (map[int64]PlugSetDefinition, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	raw := json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	plugSetDefs := make(map[int64]PlugSetDefinition)
	if err := json.Unmarshal(raw, &plugSetDefs); err != nil {
		return nil, err
	}

	return plugSetDefs, nil
}

func readWeaponPerkInputs(filePath string) ([]WeaponPerkInput, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var inputs []WeaponPerkInput
	if err := json.Unmarshal(data, &inputs); err != nil {
		return nil, err
	}

	return inputs, nil
}

func findWeaponHashes(itemDefs map[int64]ItemDefinition, weaponInputs []WeaponPerkInput) (map[string][]int64, map[int64]ItemDefinition, error) {
	weaponHashes := make(map[string][]int64)
	weaponNameSet := make(map[string]string) // Map normalized name to original name
	for _, input := range weaponInputs {
		normalizedWeaponName := strings.ToLower(strings.TrimSpace(input.WeaponName))
		weaponNameSet[normalizedWeaponName] = input.WeaponName
	}

	weaponDefinitions := make(map[int64]ItemDefinition)

	for hash, item := range itemDefs {
		if item.Redacted {
			continue
		}
		itemNameLower := strings.ToLower(strings.TrimSpace(item.DisplayProperties.Name))
		if originalName, exists := weaponNameSet[itemNameLower]; exists {
			// Add the item's hash to the list for this weapon name
			weaponHashes[originalName] = append(weaponHashes[originalName], hash)
			weaponDefinitions[hash] = item
		}
	}

	return weaponHashes, weaponDefinitions, nil
}

func findWeaponPossiblePerkHashes(
	weaponDefs map[int64]ItemDefinition,
	itemDefs map[int64]ItemDefinition,
	plugSetDefs map[int64]PlugSetDefinition,
	desiredPerkNames []string,
) (map[string][]int64, map[int64]string, map[int64][]int64, error) {
	perkHashes := make(map[string][]int64)      // Map perk name to list of hashes
	perkHashesReverse := make(map[int64]string) // Map perk hash to perk name
	weaponPossiblePerksMap := make(map[int64][]int64)

	// Normalize desired perk names for matching
	desiredPerkNameSet := make(map[string]string) // Map normalized name to original name
	for _, name := range desiredPerkNames {
		normalizedName := strings.ToLower(strings.TrimSpace(name))
		desiredPerkNameSet[normalizedName] = name
	}

	for weaponHash, weaponDef := range weaponDefs {
		possiblePerks := make(map[int64]struct{})
		if weaponDef.Sockets.SocketEntries != nil {
			for _, socket := range weaponDef.Sockets.SocketEntries {
				plugSetHash := socket.RandomizedPlugSetHash
				if plugSetHash == 0 {
					plugSetHash = socket.ReusablePlugSetHash
				}
				if plugSetHash == 0 {
					continue
				}
				plugSet, exists := plugSetDefs[plugSetHash]
				if !exists {
					continue
				}
				for _, plugItem := range plugSet.ReusablePlugItems {
					perkHash := plugItem.PlugItemHash
					possiblePerks[perkHash] = struct{}{}
				}
			}
		}
		// Convert to slice
		var perkHashList []int64
		for perkHash := range possiblePerks {
			perkHashList = append(perkHashList, perkHash)
		}
		weaponPossiblePerksMap[weaponHash] = perkHashList

		// Populate perkHashes and perkHashesReverse maps with perk names
		for _, perkHash := range perkHashList {
			perkDef, exists := itemDefs[perkHash]
			if !exists {
				continue
			}
			perkName := strings.ToLower(strings.TrimSpace(perkDef.DisplayProperties.Name))
			if perkName == "" {
				continue
			}
			// Check if perkName matches any of the desired perks or their enhanced versions
			for normalizedDesiredName, originalDesiredName := range desiredPerkNameSet {
				if perkName == normalizedDesiredName || perkName == "enhanced "+normalizedDesiredName {
					// Add the perk hash to the list for this perk name
					perkHashes[originalDesiredName] = append(perkHashes[originalDesiredName], perkHash)
					perkHashesReverse[perkHash] = perkDef.DisplayProperties.Name // Use original name
					break
				}
			}
		}
	}

	return perkHashes, perkHashesReverse, weaponPossiblePerksMap, nil
}

func generateWeaponDataFile(
	weaponInputs []WeaponPerkInput,
	weaponHashesMap map[string][]int64,
	perkHashesMap map[string][]int64,
	perkHashesReverseMap map[int64]string,
	weaponPossiblePerksMap map[int64][]int64,
	weaponDefinitions map[int64]ItemDefinition,
	outputPath string,
) error {
	// Ensure the output directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Create or overwrite the weapon_data.go file
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the package declaration
	fmt.Fprintln(file, "package constants")

	// Write WeaponHashes map
	fmt.Fprintln(file, "// WeaponHashes maps weapon names to their item hashes")
	fmt.Fprintln(file, "var WeaponHashes = map[string][]int64{")
	for weaponName, hashes := range weaponHashesMap {
		fmt.Fprintf(file, "    \"%s\": {", escapeString(weaponName))
		for _, hash := range hashes {
			fmt.Fprintf(file, "%d, ", hash)
		}
		fmt.Fprintln(file, "},")
	}
	fmt.Fprintln(file, "}")

	// Write WeaponTypes map
	fmt.Fprintln(file, "// WeaponTypes maps weapon names to their types")
	fmt.Fprintln(file, "var WeaponTypes = map[string]string{")
	for weaponName, hashes := range weaponHashesMap {
		if len(hashes) > 0 {
			weaponDef := weaponDefinitions[hashes[0]]
			fmt.Fprintf(file, "    \"%s\": \"%s\",\n", escapeString(weaponName), escapeString(weaponDef.ItemTypeDisplayName))
		}
	}
	fmt.Fprintln(file, "}")

	// Write WeaponIcons map
	fmt.Fprintln(file, "// WeaponIcons maps weapon names to their icons")
	fmt.Fprintln(file, "var WeaponIcons = map[string]string{")
	for weaponName, hashes := range weaponHashesMap {
		if len(hashes) > 0 {
			weaponDef := weaponDefinitions[hashes[0]]
			fmt.Fprintf(file, "    \"%s\": \"%s\",\n", escapeString(weaponName), escapeString(weaponDef.DisplayProperties.Icon))
		}
	}
	fmt.Fprintln(file, "}")

	// Write PerkHashes map
	fmt.Fprintln(file, "// PerkHashes maps perk names to their hashes")
	fmt.Fprintln(file, "var PerkHashes = map[string][]int64{")
	for perkName, hashes := range perkHashesMap {
		fmt.Fprintf(file, "    \"%s\": {", escapeString(perkName))
		for _, hash := range hashes {
			fmt.Fprintf(file, "%d, ", hash)
		}
		fmt.Fprintln(file, "},")
	}
	fmt.Fprintln(file, "}")

	// Write PerkDescriptions map
	fmt.Fprintln(file, "// PerkDescriptions maps perk names to their descriptions")
	fmt.Fprintln(file, "var PerkDescriptions = map[string]string{")
	for perkName, hashes := range perkHashesMap {
		if len(hashes) > 0 {
			weaponDef := weaponDefinitions[hashes[0]]
			fmt.Fprintf(file, "    \"%s\": \"%s\",\n", escapeString(perkName), escapeString(weaponDef.DisplayProperties.Description))
		}
	}
	fmt.Fprintln(file, "}")

	// Write PerkHashesReverse map
	fmt.Fprintln(file, "// PerkHashesReverse maps perk hashes to their names")
	fmt.Fprintln(file, "var PerkHashesReverse = map[int64]string{")
	for hash, perkName := range perkHashesReverseMap {
		fmt.Fprintf(file, "    %d: \"%s\",\n", hash, escapeString(perkName))
	}
	fmt.Fprintln(file, "}")

	// Write WeaponDesiredPerks map
	fmt.Fprintln(file, "// WeaponDesiredPerks maps weapon names to their desired perk hashes")
	fmt.Fprintln(file, "var WeaponDesiredPerks = map[string][]int64{")
	for _, input := range weaponInputs {
		desiredPerkHashes := []int64{}
		for _, perkName := range input.DesiredPerks {
			if hashes, exists := perkHashesMap[perkName]; exists {
				desiredPerkHashes = append(desiredPerkHashes, hashes...)
			}
		}
		fmt.Fprintf(file, "    \"%s\": {", escapeString(input.WeaponName))
		for _, hash := range desiredPerkHashes {
			fmt.Fprintf(file, "%d, ", hash)
		}
		fmt.Fprintln(file, "},")
	}
	fmt.Fprintln(file, "}")

	// Write WeaponBuckets map
	fmt.Fprintln(file, "// WeaponBuckets maps weapon names to their buckets")
	fmt.Fprintln(file, "var WeaponBuckets = map[string]string{")
	for _, input := range weaponInputs {
		fmt.Fprintf(file, "    \"%s\": \"%s\",\n", escapeString(input.WeaponName), input.Bucket)
	}
	fmt.Fprintln(file, "}")

	// Write WeaponBuckets map
	fmt.Fprintln(file, "// WeaponRanks maps weapon names to their bucket rank")
	fmt.Fprintln(file, "var WeaponRanks = map[string]string{")
	for _, input := range weaponInputs {
		fmt.Fprintf(file, "    \"%s\": \"%s\",\n", escapeString(input.WeaponName), input.Rank)
	}
	fmt.Fprintln(file, "}")

	// Write WeaponDescription map
	fmt.Fprintln(file, "// WeaponDescriptions maps weapon names to their descriptions")
	fmt.Fprintln(file, "var WeaponDescriptions = map[string]string{")
	for _, input := range weaponInputs {
		fmt.Fprintf(file, "    \"%s\": \"%s\",\n", escapeString(input.WeaponName), input.Description)
	}
	fmt.Fprintln(file, "}")

	// Write WeaponSource map
	fmt.Fprintln(file, "// WeaponSource maps weapon names to their source in-game")
	fmt.Fprintln(file, "var WeaponSource = map[string]string{")
	for _, input := range weaponInputs {
		fmt.Fprintf(file, "    \"%s\": \"%s\",\n", escapeString(input.WeaponName), input.Source)
	}
	fmt.Fprintln(file, "}")

	// Write WeaponPossiblePerks map
	fmt.Fprintln(file, "// WeaponPossiblePerks maps weapon hashes to their possible perk hashes")
	fmt.Fprintln(file, "var WeaponPossiblePerks = map[int64][]int64{")
	for weaponHash, perkHashes := range weaponPossiblePerksMap {
		fmt.Fprintf(file, "    %d: {", weaponHash)
		for _, hash := range perkHashes {
			fmt.Fprintf(file, "%d, ", hash)
		}
		fmt.Fprintln(file, "},")
	}
	fmt.Fprintln(file, "}")

	return nil
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}
