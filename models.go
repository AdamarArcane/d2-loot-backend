package main

type ProfileData struct {
	Response struct {
		Profile struct {
			Data struct {
				UserInfo struct {
					BungieGlobalDisplayName string `json:"bungieGlobalDisplayName"`
				} `json:"userInfo"`
			} `json:"data"`
		} `json:"profile"`
		ProfileInventory struct {
			Data struct {
				Items []InventoryItem `json:"items"`
			} `json:"data"`
		} `json:"profileInventory"`
		CharacterInventories struct {
			Data map[string]struct {
				Items []InventoryItem `json:"items"`
			} `json:"data"`
		} `json:"characterInventories"`
		CharacterEquipment struct {
			Data map[string]struct {
				Items []InventoryItem `json:"items"`
			} `json:"data"`
		} `json:"characterEquipment"`
		ItemComponents struct {
			Sockets struct {
				Data map[string]ItemSockets `json:"data"`
			} `json:"sockets"`
		} `json:"itemComponents"`
	} `json:"Response"`
	// ... other fields ...
}

type InventoryItem struct {
	ItemHash       int64  `json:"itemHash"`
	ItemInstanceID string `json:"itemInstanceId"`
	BucketHash     int64  `json:"bucketHash"`
}

type ItemSockets struct {
	Sockets []Socket `json:"sockets"`
}

type Socket struct {
	PlugHash  int64 `json:"plugHash"`
	IsEnabled bool  `json:"isEnabled"`
	IsVisible bool  `json:"isVisible"`
}

type SocketsData struct {
	Data map[string]ItemSockets `json:"data"`
}

type Perk struct {
	Name        string `json:"name"`
	Obtained    bool   `json:"obtained"`
	Description string `json:"description"`
}

type WeaponDetail struct {
	WeaponName       string   `json:"weaponName"`
	Icon             string   `json:"icon"`
	WeaponType       string   `json:"weaponType"`
	WeaponBucket     string   `json:"weaponBucket"`
	Points           float64  `json:"points"`
	Perks            []Perk   `json:"perks"`
	RecommendedPerks []string `json:"recommendedPerks"`
	Obtained         bool     `json:"obtained"`
	Description      string   `json:"description"`
	Source           string   `json:"source"`
}

type NextImportantGun struct {
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	WeaponType  string  `json:"weaponType"`
	Description string  `json:"description"`
	Source      string  `json:"source"`
	Points      float64 `json:"points"`
}

type ResponseData struct {
	Username         string           `json:"username"`         // User's display name
	InventoryRating  InventoryRating  `json:"inventoryRating"`  // Overall inventory rating
	NextImportantGun NextImportantGun `json:"nextImportantGun"` // Next weapon to acquire
	WeaponDetails    []WeaponDetail   `json:"weaponDetails"`    // Detailed information about each weapon
	BucketDetails    []BucketDetail   `json:"bucketDetails"`    // Detailed information about each bucket
}

type InventoryRating struct {
	TotalPoints       float64 `json:"totalPoints"`
	MaxPossiblePoints float64 `json:"maxPossiblePoints"`
	WeeklyChange      float64 `json:"weeklyChange"`
}

type BucketDetail struct {
	Name             string  `json:"name"`             // Name of the bucket
	TotalOptions     int     `json:"totalOptions"`     // Total number of weapons available in the bucket
	ObtainedCount    int     `json:"obtainedCount"`    // Number of weapons obtained by the user in the bucket
	MaxPoints        float64 `json:"maxPoints"`        // Maximum possible points for the bucket
	CurrentPoints    float64 `json:"currentPoints"`    // Current points based on obtained weapons
	AdditionalCount  int     `json:"additionalCount"`  // Number of additional weapons obtained beyond the first
	AdditionalPoints float64 `json:"additionalPoints"` // Points from additional weapons
}
