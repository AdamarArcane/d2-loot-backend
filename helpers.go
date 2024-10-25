package main

import (
	"strings"

	"github.com/adamararcane/d2optifarm/backend/cmd/constants"
)

// Helper function to get recommended perk names for a weapon
func getRecommendedPerkNames(weaponName string) []string {
	desiredPerkHashes := constants.WeaponDesiredPerks[weaponName]
	perkHashToName := constants.PerkHashesReverse

	recommendedPerks := []string{}
	uniquePerkNames := make(map[string]struct{})
	for _, perkHash := range desiredPerkHashes {
		if perkName, exists := perkHashToName[perkHash]; exists {
			// Remove "Enhanced " prefix if present for consistency
			normalizedPerkName := strings.Replace(perkName, "Enhanced ", "", 1)
			if _, exists := uniquePerkNames[normalizedPerkName]; !exists {
				uniquePerkNames[normalizedPerkName] = struct{}{}
				recommendedPerks = append(recommendedPerks, normalizedPerkName)
			}
		}
	}
	return recommendedPerks
}
