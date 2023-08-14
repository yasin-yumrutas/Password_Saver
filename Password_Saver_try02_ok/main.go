package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// PostgreSQL veritabanına bağlanma

	connStr := "user=postgres password=321654 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Veritabanına bağlanırken hata oluştu:", err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	server := &http.Server{Addr: ":8082"}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP sunucusu çalışırken hata oluştu:", err)
		}
	}()

	// Belirli bir süre sonra sunucuyu kapat
	time.Sleep(time.Second * 60) // 60 saniye (1 dakika) bekleyin

	log.Println("HTTP sunucusu kapanıyor...")
	if err := server.Shutdown(nil); err != nil {
		log.Fatal("HTTP sunucusunu kapatırken hata oluştu:", err)
	}

	// time.Sleep(time.Second * 60)
	// JSON dosyasını okuma
	// jsonFile, err := os.Open("data.json")
	// if err != nil {
	// 	log.Fatal("JSON dosyasını açarken hata oluştu:", err)
	// }
	// defer jsonFile.Close()

	jsonFilePath := "C:/Users/yasin/Downloads/data.json" // Yolunuzu buraya ekleyin
	jsonFile, err := os.Open(jsonFilePath)
	// jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		log.Fatal("JSON dosyasını açarken hata oluştu:", err)
	}
	defer jsonFile.Close()

	var users []User
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(&users); err != nil {
		log.Fatal("JSON dosyasını ayrıştırırken hata oluştu:", err)
	}

	// JSON verilerini veritabanına ekleme
	for _, user := range users {
		insertQuery := "INSERT INTO users (username, password) VALUES ($1, $2)"
		_, err := db.Exec(insertQuery, user.Username, user.Password)
		if err != nil {
			log.Printf("Veritabanına ekleme hatası (%s, %s): %v", user.Username, user.Password, err)
		} else {
			fmt.Printf("Veri eklendi: %s\n", user.Username)
		}
	}
}
