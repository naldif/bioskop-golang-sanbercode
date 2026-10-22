package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Bioskop struct {
	ID     uint    `gorm:"primaryKey" json:"id"`
	Nama   string  `gorm:"not null" json:"nama" binding:"required"`
	Lokasi string  `gorm:"not null" json:"lokasi" binding:"required"`
	Rating float64 `json:"rating"`
}

type CreateBioskopInput struct {
	Nama   string  `json:"nama" binding:"required"`
	Lokasi string  `json:"lokasi" binding:"required"`
	Rating float64 `json:"rating"`
}

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

var DB *gorm.DB

func initDB() {

	dsn := "host=localhost user=postgres password=postgres dbname=bioskop port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	err = DB.AutoMigrate(&Bioskop{})
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi tabel: %v", err)
	}

	fmt.Println("Koneksi database dan migrasi berhasil!")
}

// POST /bioskop
func createBioskop(c *gin.Context) {
	var input CreateBioskopInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validasi gagal: Nama dan Lokasi wajib diisi",
			"error":   err.Error(),
		})
		return
	}

	bioskop := Bioskop{
		Nama:   input.Nama,
		Lokasi: input.Lokasi,
		Rating: input.Rating,
	}

	if result := DB.Create(&bioskop); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan data bioskop",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Status:  "success",
		Message: "Bioskop berhasil ditambahkan",
		Data:    bioskop,
	})
}

func main() {

	initDB()

	r := gin.Default()

	r.POST("/bioskop", createBioskop)

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}