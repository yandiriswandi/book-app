package controllers

import (
	"bioskop-app/config"
	"bioskop-app/models"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateBook(c *gin.Context) {
	var book models.Book
	// Validasi input
	if err := c.ShouldBindJSON(&book); err != nil {
		// Tapi jangan langsung return, simpan errornya dulu
		validationError := err.Error()

		// Cek manual field kosong
		if book.Nama == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "Nama tidak boleh kosong",
			})
			return
		}

		if book.Rating > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "Rating tidak boleh lebih dari 100",
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

	query := `INSERT INTO book (nama, lokasi,kategori, rating,description) VALUES ($1, $2, $3,$4) RETURNING id`
	err := config.DB.QueryRow(query, book.Nama, book.Kategori, book.Rating, book.Description).Scan(&book.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Gagal menambahkan book",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Data berhasil dibuat!",
		"data":    book,
	})
}

func UpdateBook(c *gin.Context) {
	var book models.Book
	id := c.Param("id") // ambil id dari URL parameter

	// Validasi input
	if err := c.ShouldBindJSON(&book); err != nil {
		validationError := err.Error()

		if book.Nama == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "Nama tidak boleh kosong",
			})
			return
		}
		if book.Rating > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  400,
				"message": "Rating tidak boleh lebih dari 100",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": validationError,
		})
		return
	}

	query := `UPDATE book SET nama = $1, kategori = $2, rating = $3,description=$4 WHERE id = $5`
	res, err := config.DB.Exec(query, book.Nama, book.Kategori, book.Rating, book.Description, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Gagal memperbarui book",
		})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "Book tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Data berhasil diperbarui!",
		"data":    book,
	})
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id") // ambil id dari URL parameter

	query := `DELETE FROM book WHERE id = $1`
	res, err := config.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Gagal menghapus book",
		})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "Book tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Book berhasil dihapus",
	})
}

func GetBookList(c *gin.Context) {
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
	err = config.DB.QueryRow(`SELECT COUNT(*) FROM book`).Scan(&totalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Gagal mengambil total data",
		})
		return
	}

	// Ambil data book
	rows, err := config.DB.Query(`SELECT id, nama, kategori, rating,description FROM book ORDER BY id LIMIT $1 OFFSET $2`, size, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Gagal mengambil data book",
		})
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Nama, &book.Kategori, &book.Rating, &book.Description); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  500,
				"message": "Gagal membaca data book",
			})
			return
		}
		books = append(books, book)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"message":   "Berhasil mengambil list book",
		"data":      books,
		"page":      page,
		"size":      size,
		"totalData": totalData,
	})
}

func GetBookByID(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	query := `SELECT id, nama, kategori, rating,description FROM book WHERE id = $1`
	err := config.DB.QueryRow(query, id).Scan(&book.ID, &book.Nama, &book.Kategori, &book.Rating, &book.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  404,
				"message": "Book tidak ditemukan",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  500,
				"message": "Gagal mengambil detail book",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Berhasil mengambil detail book",
		"data":    book,
	})
}
