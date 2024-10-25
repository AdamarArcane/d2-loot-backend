// constants/bucket_constants.go

package constants

type BucketPoint struct {
	BucketName          string
	MaxPoints           float64 // Maximum points for the bucket
	AdditionalWeaponPts float64 // Fixed bonus for additional weapons (e.g., 10% of MaxPoints)
	DiminishingFactor   float64 // Scaling factor for diminishing returns (e.g., 0.5)
}

var BucketPoints = []BucketPoint{
	{
		BucketName:          "Orb Generation",
		MaxPoints:           10.0,
		AdditionalWeaponPts: 1.0, // 10% of 10.0
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Kinetic Rocket Sidearm",
		MaxPoints:           10.0,
		AdditionalWeaponPts: 1.0,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Energy Rocket Sidearm",
		MaxPoints:           10.0,
		AdditionalWeaponPts: 1.0,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "DPS Heavy Grenade Launcher",
		MaxPoints:           10.0,
		AdditionalWeaponPts: 1.0,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Exotic Energy Primary",
		MaxPoints:           10.0,
		AdditionalWeaponPts: 1.0,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Exotic DPS (Consistent)",
		MaxPoints:           9.0,
		AdditionalWeaponPts: 0.9, // 10% of 9.0
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Exotic DPS (Total Damage)",
		MaxPoints:           9.0,
		AdditionalWeaponPts: 0.9,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Exotic Debuff",
		MaxPoints:           10.0,
		AdditionalWeaponPts: 1.0,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Team Support Weapon",
		MaxPoints:           7.0,
		AdditionalWeaponPts: 0.7, // 10% of 7.0
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Add Clear with Damage Resistance",
		MaxPoints:           7.0,
		AdditionalWeaponPts: 0.7,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Kinetic One-Two Punch",
		MaxPoints:           7.0,
		AdditionalWeaponPts: 0.7,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Weaken on Demand",
		MaxPoints:           6.0,
		AdditionalWeaponPts: 0.6, // 10% of 6.0
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Hitscan Overload Stun",
		MaxPoints:           6.0,
		AdditionalWeaponPts: 0.6,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Kinetic Sniper",
		MaxPoints:           6.0,
		AdditionalWeaponPts: 0.6,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Transcendance Generation",
		MaxPoints:           6.0,
		AdditionalWeaponPts: 0.6,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Machine Gun",
		MaxPoints:           6.0,
		AdditionalWeaponPts: 0.6,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Movement Sword",
		MaxPoints:           6.0,
		AdditionalWeaponPts: 0.6,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Energy Primary",
		MaxPoints:           5.0,
		AdditionalWeaponPts: 0.5, // 10% of 5.0
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Exotic Heavy Burst",
		MaxPoints:           5.0,
		AdditionalWeaponPts: 0.5,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Exotic Add Clear",
		MaxPoints:           5.0,
		AdditionalWeaponPts: 0.5,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Ammoless DPS",
		MaxPoints:           4.0,
		AdditionalWeaponPts: 0.4, // 10% of 4.0
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Movement Grenade Launcher",
		MaxPoints:           4.0,
		AdditionalWeaponPts: 0.4,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Energy One-Two Punch",
		MaxPoints:           3.0,
		AdditionalWeaponPts: 0.3, // 10% of 3.0
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Kinetic Burst Damage",
		MaxPoints:           3.0,
		AdditionalWeaponPts: 0.3,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Energy Damage Sniper",
		MaxPoints:           3.0,
		AdditionalWeaponPts: 0.3,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Kinetic Fusion",
		MaxPoints:           3.0,
		AdditionalWeaponPts: 0.3,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Energy Fusion",
		MaxPoints:           3.0,
		AdditionalWeaponPts: 0.3,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Energy Wave-Frame",
		MaxPoints:           3.0,
		AdditionalWeaponPts: 0.3,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Kinetic Wave-Frame",
		MaxPoints:           3.0,
		AdditionalWeaponPts: 0.3,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Kinetic Blind",
		MaxPoints:           2.0,
		AdditionalWeaponPts: 0.2, // 10% of 2.0
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Energy Blind",
		MaxPoints:           2.0,
		AdditionalWeaponPts: 0.2,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Glaive",
		MaxPoints:           1.0,
		AdditionalWeaponPts: 0.1, // 10% of 1.0
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Energy Trace",
		MaxPoints:           2.0,
		AdditionalWeaponPts: 0.2,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "DPS Sword",
		MaxPoints:           3.0,
		AdditionalWeaponPts: 0.3,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "DPS Rocket",
		MaxPoints:           4.0,
		AdditionalWeaponPts: 0.4,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Linear",
		MaxPoints:           2.0,
		AdditionalWeaponPts: 0.2,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Kinetic Primary",
		MaxPoints:           2.0,
		AdditionalWeaponPts: 0.2,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Kinetic Hand Cannon (Lucky Pants)",
		MaxPoints:           2.0,
		AdditionalWeaponPts: 0.2,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Energy Hand Cannon (Lucky Pants)",
		MaxPoints:           2.0,
		AdditionalWeaponPts: 0.2,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Exotic Kinetic Primary",
		MaxPoints:           3.0,
		AdditionalWeaponPts: 0.3,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Super Generation",
		MaxPoints:           2.0,
		AdditionalWeaponPts: 0.2,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Kinetic SMG (Peacekeepers)",
		MaxPoints:           2.0,
		AdditionalWeaponPts: 0.2,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Exotic Energy Burst",
		MaxPoints:           2.0,
		AdditionalWeaponPts: 0.2,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Survivability",
		MaxPoints:           2.0,
		AdditionalWeaponPts: 0.2,
		DiminishingFactor:   0.5,
	},
	{
		BucketName:          "Exotic Special DPS (Total)",
		MaxPoints:           1.0,
		AdditionalWeaponPts: 0.1,
		DiminishingFactor:   0.5,
	},
}
