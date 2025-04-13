package models

type Book struct {
	ID           int     `db:"id" json:"id"`
	Nama         string  `db:"nama" json:"nama" binding:"required"`
	KategoriID   int     `db:"kategori_id" json:"kategori_id" binding:"required"`
	Description  string  `db:"description" json:"description"`
	Rating       float64 `db:"rating" json:"rating" binding:"lte=100"`
	KategoriNama string  `db:"kategori_nama" json:"kategori_nama"`
}

type Category struct {
	ID          int    `db:"id" json:"id"`
	Nama        string `db:"nama" json:"nama" binding:"required"`
	Description string `db:"description" json:"description"`
}
