package models

type Book struct {
	ID          int     `db:"id" json:"id"`
	Nama        string  `db:"nama" json:"nama" binding:"required"`
	Kategori    string  `db:"kategori" json:"kategori" binding:"required"`
	Description string  `db:"description" json:"descriptrion"`
	Rating      float64 `db:"rating" json:"rating" binding:"lte=100"`
}

type Category struct {
	ID          int    `db:"id" json:"id"`
	Nama        string `db:"nama" json:"nama" binding:"required"`
	Description string `db:"description" json:"descriptrion"`
}
