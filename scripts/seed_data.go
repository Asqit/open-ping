package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "openping.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	targets := []struct {
		name string
		url  string
	}{
		{"google", "https://www.google.com"},
		{"github", "https://github.com"},
		{"aws", "https://aws.amazon.com"},
		{"cloudflare", "https://www.cloudflare.com"},
		{"vercel", "https://vercel.com"},
		{"netlify", "https://www.netlify.com"},
	}

	now := time.Now()
	fmt.Println("Seeding 90 days of fake data...")

	for _, target := range targets {
		for day := 0; day < 90; day++ {
			date := now.AddDate(0, 0, -day)
			
			pingsPerDay := rand.Intn(40) + 10
			
			for i := 0; i < pingsPerDay; i++ {
				hour := rand.Intn(24)
				minute := rand.Intn(60)
				timestamp := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, time.UTC)
				
				success := rand.Float64() < 0.95
				if day%15 == 0 {
					success = rand.Float64() < 0.70
				}
				
				status := 200
				if !success {
					status = 500
				}
				
				latency := rand.Intn(180) + 20
				
				_, err := db.Exec(`
					INSERT INTO pings (target, target_website, status, success, latency, timestamp)
					VALUES (?, ?, ?, ?, ?, ?)`,
					target.name, target.url, status, success, latency, timestamp,
				)
				if err != nil {
					log.Printf("Failed to insert: %v", err)
				}
			}
		}
		fmt.Printf("✓ Seeded %s\n", target.name)
	}

	fmt.Println("✅ Done! Seeded 90 days of data.")
}
