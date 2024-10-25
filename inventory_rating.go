package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"path/filepath"
	"sort"

	"github.com/adamararcane/d2-loot-backend/cmd/constants"
)

// WeaponDefinition represents the structure of a weapon from the JSON file.
type WeaponDefinition struct {
	WeaponName   string   `json:"weaponName"`
	DesiredPerks []string `json:"desiredPerks"`
	Description  string   `json:"description"`
	Source       string   `json:"source"`
	Bucket       string   `json:"bucket"`
	Rank         int      `json:"rank,string"` // Parses "rank": "1" as integer 1
}

// PerkWeights assigns weights to desired perks (optional).
var PerkWeights = map[string]float64{
	"PerkA": 10.0,
	"PerkB": 15.0,
	"PerkC": 12.0,
	"PerkD": 8.0,
	// Add other perks and their weights as needed
}

// loadWeaponDefinitions loads weapon definitions from a JSON file.
func loadWeaponDefinitions(jsonPath string) ([]WeaponDefinition, error) {
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read weapons JSON file: %w", err)
	}

	var weapons []WeaponDefinition
	err = json.Unmarshal(data, &weapons)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal weapons JSON: %w", err)
	}

	err = validateWeaponDefinitions(weapons)
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return weapons, nil
}

// validateWeaponDefinitions ensures that each weapon has necessary fields.
func validateWeaponDefinitions(weapons []WeaponDefinition) error {
	for _, weapon := range weapons {
		if weapon.WeaponName == "" {
			return fmt.Errorf("weapon with empty name found")
		}
		if weapon.Bucket == "" {
			return fmt.Errorf("weapon '%s' has empty bucket", weapon.WeaponName)
		}
		if weapon.Rank < 1 {
			return fmt.Errorf("weapon '%s' has invalid rank '%d'", weapon.WeaponName, weapon.Rank)
		}
		// Add more validation rules as needed
	}
	return nil
}

// buildHashToWeaponMap creates a map from item hash to WeaponDefinition.
func buildHashToWeaponMap(weapons []WeaponDefinition) (map[int64]WeaponDefinition, error) {
	hashToWeapon := make(map[int64]WeaponDefinition)
	for _, weapon := range weapons {
		hashes, exists := constants.WeaponHashes[weapon.WeaponName]
		if !exists {
			log.Printf("Warning: No hashes found for weapon '%s'", weapon.WeaponName)
			continue // Skip weapons without defined hashes
		}
		for _, hash := range hashes {
			hashToWeapon[hash] = weapon
		}
	}
	return hashToWeapon, nil
}

// findBucketIndex returns the index of a bucket in BucketPoints based on its name.
func findBucketIndex(bucketName string, bucketPoints []constants.BucketPoint) int {
	for i, bp := range bucketPoints {
		if bp.BucketName == bucketName {
			return i
		}
	}
	return -1
}

// calculateBucketPoints calculates the total points for a bucket based on its weapons.
func calculateBucketPoints(weapons []WeaponDefinition, bp constants.BucketPoint) float64 {
	// Sort weapons by Rank ascending
	sort.Slice(weapons, func(i, j int) bool {
		return weapons[i].Rank < weapons[j].Rank
	})

	if len(weapons) == 0 {
		return 0.0
	}

	totalPoints := 0.0

	// Top-tier weapon
	topWeapon := weapons[0]
	deduction := 0.2 * float64(topWeapon.Rank-1) * bp.MaxPoints
	topPoints := bp.MaxPoints - deduction
	if topPoints < 0 {
		topPoints = 0
	}
	totalPoints += topPoints

	// Additional weapons with diminishing returns
	for i := 1; i < len(weapons); i++ {
		bonus := bp.AdditionalWeaponPts * math.Pow(bp.DiminishingFactor, float64(i-1))
		totalPoints += bonus
	}

	return totalPoints
}

// rateInventory calculates the inventory rating based on the player's profile data.
func (api *apiConfig) rateInventory(profileData ProfileData) (ResponseData, error) {
	// Step 1: Load weapon definitions from JSON
	jsonPath := filepath.Join("cmd", "generate_constants", "weapons_and_perks.json") // Adjust the path as necessary
	weapons, err := loadWeaponDefinitions(jsonPath)
	if err != nil {
		return ResponseData{}, err
	}

	// Step 2: Build hashToWeapon map
	hashToWeapon, err := buildHashToWeaponMap(weapons)
	if err != nil {
		return ResponseData{}, err
	}

	// Step 3: Collect all inventory items
	allItems := []InventoryItem{}
	for _, character := range profileData.Response.CharacterInventories.Data {
		allItems = append(allItems, character.Items...)
	}
	allItems = append(allItems, profileData.Response.ProfileInventory.Data.Items...)
	for _, character := range profileData.Response.CharacterEquipment.Data {
		allItems = append(allItems, character.Items...)
	}

	// Step 4: Initialize max possible points
	maxPossiblePoints := 0.0
	for _, bp := range constants.BucketPoints {
		maxPossiblePoints += bp.MaxPoints + (bp.AdditionalWeaponPts * 5) // Assuming a max of 5 additional weapons for max potential
	}

	// Step 5: Build a map of desired perk hashes per weapon
	desiredPerkHashesMap := make(map[string]map[int64]struct{}) // weaponName -> set of desired perk hashes
	for _, weapon := range weapons {
		desiredPerkHashesMap[weapon.WeaponName] = make(map[int64]struct{})
		for _, perk := range weapon.DesiredPerks {
			perkHashes, exists := constants.PerkHashes[perk]
			if exists {
				for _, perkHash := range perkHashes {
					desiredPerkHashesMap[weapon.WeaponName][perkHash] = struct{}{}
				}
			} else {
				log.Printf("Warning: Perk '%s' not found in PerkHashes map for weapon '%s'", perk, weapon.WeaponName)
			}
		}
	}

	// Step 6: Initialize bucket ownership map
	ownedWeaponsPerBucket := make(map[string][]WeaponDefinition)

	/// var iconUrls []string

	// Step 7: Process each inventory item
	for _, item := range allItems {
		weaponDef, isValuable := hashToWeapon[item.ItemHash]
		if !isValuable {
			continue // Not a valuable weapon, skip
		}

		bucketName := weaponDef.Bucket
		if bucketName == "" {
			log.Printf("Warning: Weapon '%s' does not have a bucket assigned", weaponDef.WeaponName)
			continue // Weapon not assigned to any bucket
		}

		// Extract perks
		socketsData, exists := profileData.Response.ItemComponents.Sockets.Data[item.ItemInstanceID]
		if !exists {
			continue // No socket data for this item
		}

		// Collect perk hashes from the item
		itemPerks := make(map[int64]struct{})
		for _, socket := range socketsData.Sockets {
			perkHash := socket.PlugHash
			if _, exists := constants.PerkHashesReverse[perkHash]; exists {
				itemPerks[perkHash] = struct{}{}
			}
		}

		// Get desired perks for this weapon
		weaponDesiredPerks, exists := desiredPerkHashesMap[weaponDef.WeaponName]
		if !exists || len(weaponDesiredPerks) == 0 {
			log.Printf("Warning: No desired perks defined for weapon '%s'", weaponDef.WeaponName)
			continue // No desired perks defined for this weapon, skip
		}

		// Count matching perks for this instance
		matchingPerkCount := 0
		for perkHash := range itemPerks {
			if _, desired := weaponDesiredPerks[perkHash]; desired {
				matchingPerkCount++
			}
		}

		// If this instance has at least two matching perks, consider it
		if matchingPerkCount >= 2 {
			ownedWeaponsPerBucket[bucketName] = append(ownedWeaponsPerBucket[bucketName], weaponDef)
		}
	}

	// Step 8: Assign current points based on owned weapons per bucket
	currentBucketPoints := make(map[string]float64)
	for _, bp := range constants.BucketPoints {
		bucketName := bp.BucketName
		ownedWeapons, owned := ownedWeaponsPerBucket[bucketName]
		if owned && len(ownedWeapons) > 0 {
			currentPoints := calculateBucketPoints(ownedWeapons, bp)
			currentBucketPoints[bucketName] = currentPoints
		} else {
			currentBucketPoints[bucketName] = 0.0
		}
	}

	// Step 9: Prepare weapon details with potential points
	weaponDetails := []WeaponDetail{}
	var weaponsToGet []WeaponDefinition
	for _, weapon := range weapons {
		obtained := false
		if ownedWeapons, ok := ownedWeaponsPerBucket[weapon.Bucket]; ok {
			for _, ow := range ownedWeapons {
				if ow.WeaponName == weapon.WeaponName {
					obtained = true
					break
				}
			}
		}

		// Initialize WeaponDetail
		detail := WeaponDetail{
			WeaponName:       weapon.WeaponName,
			Icon:             "https://bungie.net" + constants.WeaponIcons[weapon.WeaponName],
			WeaponBucket:     weapon.Bucket,
			WeaponType:       constants.WeaponTypes[weapon.WeaponName],
			Points:           0.0,      // To be calculated
			Perks:            []Perk{}, // Populate if needed
			RecommendedPerks: getRecommendedPerkNames(weapon.WeaponName),
			Obtained:         obtained,
			Description:      weapon.Description,
			Source:           weapon.Source,
		}

		if detail.WeaponType == "" {
			detail.WeaponType = "Unknown"
		}

		if obtained {
			// Current contribution of the weapon
			ownedWeapons := ownedWeaponsPerBucket[weapon.Bucket]
			// Retrieve current points from the map
			// No need to declare currentPoints here

			// Sort ownedWeapons by Rank ascending
			sort.Slice(ownedWeapons, func(i, j int) bool {
				return ownedWeapons[i].Rank < ownedWeapons[j].Rank
			})

			// Points for top-tier weapon
			topWeapon := ownedWeapons[0]
			bpIndex := findBucketIndex(weapon.Bucket, constants.BucketPoints)
			if bpIndex == -1 {
				log.Printf("Warning: Bucket '%s' not found for weapon '%s'", weapon.Bucket, weapon.WeaponName)
				continue
			}
			bp := constants.BucketPoints[bpIndex]
			deduction := 0.2 * float64(topWeapon.Rank-1) * bp.MaxPoints
			topPoints := bp.MaxPoints - deduction
			if topPoints < 0 {
				topPoints = 0
			}

			// Points for additional weapons
			additionalPoints := 0.0
			for i := 1; i < len(ownedWeapons); i++ {
				bonus := bp.AdditionalWeaponPts * math.Pow(bp.DiminishingFactor, float64(i-1))
				additionalPoints += bonus
			}

			// Total points for the bucket (not used further)
			// totalBucketPoints := topPoints + additionalPoints // Not used

			// Weapon's current contribution
			if weapon.WeaponName == topWeapon.WeaponName {
				detail.Points = topPoints
			} else {
				// Calculate additional weapon points
				for i, ow := range ownedWeapons {
					if ow.WeaponName == weapon.WeaponName {
						bonus := bp.AdditionalWeaponPts * math.Pow(bp.DiminishingFactor, float64(i-1))
						detail.Points = bonus
						break
					}
				}
			}

			// Add perk points
			perkPoints := 0.0
			for _, perk := range weapon.DesiredPerks {
				if weight, exists := PerkWeights[perk]; exists {
					perkPoints += weight
					// Optionally, add the perk to the Perks slice
					detail.Perks = append(detail.Perks, Perk{Name: perk})
				}
			}
			detail.Points += perkPoints

		} else {
			// Potential points if the weapon is obtained
			bucketName := weapon.Bucket
			bpIndex := findBucketIndex(bucketName, constants.BucketPoints)
			if bpIndex == -1 {
				log.Printf("Warning: Bucket '%s' not found for weapon '%s'", bucketName, weapon.WeaponName)
				continue
			}
			bp := constants.BucketPoints[bpIndex]

			// Current owned weapons in the bucket
			currentOwned := ownedWeaponsPerBucket[bucketName]
			currentPoints := currentBucketPoints[bucketName]

			// Simulate adding the weapon
			simulatedOwned := append([]WeaponDefinition(nil), currentOwned...) // Clone the slice
			simulatedOwned = append(simulatedOwned, weapon)

			// Calculate new bucket points
			newPoints := calculateBucketPoints(simulatedOwned, bp)

			// Potential points added by obtaining this weapon
			potentialPoints := newPoints - currentPoints

			// Ensure that potentialPoints are not negative
			if potentialPoints < 0 {
				potentialPoints = 0.0
			}

			// Add perk points
			perkPoints := 0.0
			for _, perk := range weapon.DesiredPerks {
				if weight, exists := PerkWeights[perk]; exists {
					perkPoints += weight
				}
			}
			potentialPoints += perkPoints

			detail.Points = potentialPoints

			// Add to weaponsToGet for potential next important gun
			weaponsToGet = append(weaponsToGet, weapon)
		}

		weaponDetails = append(weaponDetails, detail)
	}

	// Step 10: Determine the next important gun to acquire based on potential points
	var nextImportantGun WeaponDefinition
	maxPotentialPoints := -1.0
	for _, weapon := range weaponsToGet {
		// Find the bucket and rank of the weapon
		bpIndex := findBucketIndex(weapon.Bucket, constants.BucketPoints)
		if bpIndex == -1 {
			log.Printf("Warning: Bucket '%s' not found for weapon '%s'", weapon.Bucket, weapon.WeaponName)
			continue // Bucket not found
		}
		bp := constants.BucketPoints[bpIndex]

		// Current owned weapons in the bucket
		currentOwned := ownedWeaponsPerBucket[weapon.Bucket]
		currentPoints := currentBucketPoints[weapon.Bucket]

		// Simulate adding the weapon
		simulatedOwned := append([]WeaponDefinition(nil), currentOwned...)
		simulatedOwned = append(simulatedOwned, weapon)

		// Calculate new bucket points
		newPoints := calculateBucketPoints(simulatedOwned, bp)

		// Potential points added by obtaining this weapon
		potentialPoints := newPoints - currentPoints

		// Add perk points
		for _, perk := range weapon.DesiredPerks {
			if weight, exists := PerkWeights[perk]; exists {
				potentialPoints += weight
			}
		}

		if potentialPoints > maxPotentialPoints {
			maxPotentialPoints = potentialPoints
			nextImportantGun = weapon
		}
	}

	// Step 11: Prepare the inventory rating
	inventoryRating := InventoryRating{
		TotalPoints:       0.0, // Will be recalculated below
		MaxPossiblePoints: maxPossiblePoints,
		WeeklyChange:      0.0, // Update this if needed
	}

	// Recalculate total points based on current bucket points
	for _, bp := range constants.BucketPoints {
		inventoryRating.TotalPoints += currentBucketPoints[bp.BucketName]
	}

	// Step 12: Prepare bucket details
	bucketDetails := []BucketDetail{}
	for _, bp := range constants.BucketPoints {
		bucketName := bp.BucketName
		maxPoints := bp.MaxPoints

		// Get total weapons in this bucket
		totalOptions := 0
		for _, weapon := range weapons {
			if weapon.Bucket == bucketName {
				totalOptions++
			}
		}

		// Get obtained weapons in this bucket
		obtainedCount := 0
		additionalCount := 0
		currentPoints := 0.0
		additionalPoints := 0.0

		if ownedWeapons, owned := ownedWeaponsPerBucket[bucketName]; owned && len(ownedWeapons) > 0 {
			obtainedCount = len(ownedWeapons)
			if obtainedCount > 0 {
				// Sort ownedWeapons by Rank ascending (lower rank first)
				sort.Slice(ownedWeapons, func(i, j int) bool {
					return ownedWeapons[i].Rank < ownedWeapons[j].Rank
				})

				// Calculate points for the top-tier weapon
				firstWeapon := ownedWeapons[0]
				deduction := 0.2 * float64(firstWeapon.Rank-1) * maxPoints
				points := maxPoints - deduction
				if points < 0 {
					points = 0
				}
				currentPoints += points

				// Calculate points for additional weapons with diminishing returns
				if obtainedCount > 1 {
					additionalCount = obtainedCount - 1
					for i := 1; i < obtainedCount; i++ {
						bonus := bp.AdditionalWeaponPts * math.Pow(bp.DiminishingFactor, float64(i-1))
						additionalPoints += bonus
					}
					currentPoints += additionalPoints
				}
			}
		}

		bucketDetails = append(bucketDetails, BucketDetail{
			Name:             bucketName,
			TotalOptions:     totalOptions,
			ObtainedCount:    obtainedCount,
			MaxPoints:        maxPoints,
			CurrentPoints:    currentPoints,
			AdditionalCount:  additionalCount,
			AdditionalPoints: additionalPoints,
		})
	}

	// Step 13: Prepare the next important gun detail
	var nextGun NextImportantGun
	if nextImportantGun.WeaponName != "" {
		nextGun = NextImportantGun{
			Name:        nextImportantGun.WeaponName,
			Icon:        "https://bungie.net" + constants.WeaponIcons[nextImportantGun.WeaponName],
			WeaponType:  nextImportantGun.Bucket,
			Description: nextImportantGun.Description,
			Source:      nextImportantGun.Source,
			Points:      maxPotentialPoints,
		}
	}

	// Step 14: Prepare the response data
	responseData := ResponseData{
		Username:         profileData.Response.Profile.Data.UserInfo.BungieGlobalDisplayName,
		InventoryRating:  inventoryRating,
		NextImportantGun: nextGun,
		WeaponDetails:    weaponDetails,
		BucketDetails:    bucketDetails,
	}

	return responseData, nil
}
