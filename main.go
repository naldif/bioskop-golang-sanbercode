package main

import (
	"os"
	"errors"
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
	Error   string `json:"error,omitempty"`
}

var DB *gorm.DB

func initDB() {

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=bioskop port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}

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

// GET /bioskop (Get All)
func getAllBioskop(c *gin.Context) {
	var daftarBioskop []Bioskop

	if err := DB.Find(&daftarBioskop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Gagal mengambil daftar bioskop",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Daftar bioskop berhasil diambil",
		Data:    daftarBioskop,
	})
}

// GET /bioskop/:id (Get Detail by ID)
func getBioskopByID(c *gin.Context) {
	id := c.Param("id")
	var bioskop Bioskop

	if err := DB.First(&bioskop, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, APIResponse{
				Status:  "error",
				Message: fmt.Sprintf("Bioskop dengan ID %s tidak ditemukan", id),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Terjadi kesalahan pada server",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Detail bioskop berhasil ditemukan",
		Data:    bioskop,
	})
}

// PUT /bioskop/:id (Update)
func updateBioskop(c *gin.Context) {
	id := c.Param("id")
	var bioskop Bioskop

	if err := DB.First(&bioskop, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, APIResponse{
				Status:  "error",
				Message: fmt.Sprintf("Bioskop dengan ID %s tidak ditemukan", id),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Terjadi kesalahan pada server",
			Error:   err.Error(),
		})
		return
	}

	var input CreateBioskopInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Validasi gagal: Nama dan Lokasi wajib diisi",
			Error:   err.Error(),
		})
		return
	}

	bioskop.Nama = input.Nama
	bioskop.Lokasi = input.Lokasi
	bioskop.Rating = input.Rating

	if err := DB.Save(&bioskop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Gagal memperbarui data bioskop",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Data bioskop berhasil diperbarui",
		Data:    bioskop,
	})
}

// DELETE /bioskop/:id (Delete)
func deleteBioskop(c *gin.Context) {
	id := c.Param("id")
	var bioskop Bioskop

	if err := DB.First(&bioskop, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, APIResponse{
				Status:  "error",
				Message: fmt.Sprintf("Bioskop dengan ID %s tidak ditemukan", id),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Terjadi kesalahan pada server",
			Error:   err.Error(),
		})
		return
	}

	if err := DB.Delete(&bioskop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Gagal menghapus data bioskop",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: fmt.Sprintf("Bioskop dengan ID %s berhasil dihapus", id),
	})
}

func main() {

	initDB()

	r := gin.Default()

	r.POST("/bioskop", createBioskop)
	r.GET("/bioskop", getAllBioskop)
	r.GET("/bioskop/:id", getBioskopByID)
	r.PUT("/bioskop/:id", updateBioskop)
	r.DELETE("/bioskop/:id", deleteBioskop)	

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}