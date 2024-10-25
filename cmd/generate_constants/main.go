package main

import (
	"fmt"
	"log"

	"github.com/adamararcane/d2-loot-backend/internal/generator"
)

func main() {
	err := generator.GenerateWeaponData()
	if err != nil {
		log.Fatalf("Error generating weapon data: %v", err)
	}
	fmt.Println("weapon_data.go file generated successfully!")
}
