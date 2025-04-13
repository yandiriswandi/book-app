package controllers

import (
	"bioskop-app/config"
	"bioskop-app/models"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	var category models.Category
	// Validasi input
	if err := c.ShouldBindJSON(&category); err != nil {
		// Tapi jangan langsung return, simpan errornya dulu
		validationError := err.Error()

		// Cek manual field kosong
		if category.Nama == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "Nama tidak boleh kosong",
			})
			return
		}

		// Kalau field udah aman, baru kalau masih ada error di ShouldBindJSON, kirim errornya
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": validationError,
		})
		return
	}

	query := `INSERT INTO category (nama, description) VALUES ($1, $2) RETURNING id`
	err := config.DB.QueryRow(query, category.Nama, category.Description).Scan(&category.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Gagal menambahkan category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Data berhasil dibuat!",
		"data":    category,
	})
}

func UpdateCategory(c *gin.Context) {
	var category models.Category
	id := c.Param("id") // ambil id dari URL parameter

	// Validasi input
	if err := c.ShouldBindJSON(&category); err != nil {
		validationError := err.Error()

		if category.Nama == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "Nama tidak boleh kosong",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": validationError,
		})
		return
	}

	query := `UPDATE category SET nama = $1, description=$4 WHERE id = $3`
	res, err := config.DB.Exec(query, category.Nama, category.Description, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Gagal memperbarui category",
		})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "Category tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Data berhasil diperbarui!",
		"data":    category,
	})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id") // ambil id dari URL parameter

	query := `DELETE FROM category WHERE id = $1`
	res, err := config.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Gagal menghapus category",
		})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "Category tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Category berhasil dihapus",
	})
}

func GetCategoryList(c *gin.Context) {
	// Ambil query parameter page dan size, default kalau kosong
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 10
	}

	offset := (page - 1) * size

	// Hitung total data
	var totalData int
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM category`).Scan(&totalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Gagal mengambil total data",
		})
		return
	}

	// Ambil data category
	rows, err := config.DB.Query(`SELECT id, nama, description FROM category ORDER BY id LIMIT $1 OFFSET $2`, size, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Gagal mengambil data category",
		})
		return
	}
	defer rows.Close()

	var categorys []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Nama, &category.Description); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  500,
				"message": "Gagal membaca data category",
			})
			return
		}
		categorys = append(categorys, category)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"message":   "Berhasil mengambil list category",
		"data":      categorys,
		"page":      page,
		"size":      size,
		"totalData": totalData,
	})
}

func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	var category models.Category
	query := `SELECT id, nama, description FROM category WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(&category.ID, &category.Nama, &category.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  404,
				"message": "Category tidak ditemukan",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  500,
				"message": "Gagal mengambil detail category",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Berhasil mengambil detail category",
		"data":    category,
	})
}
