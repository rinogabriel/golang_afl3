package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Post represents a sample data structure for the API.
type Stock struct {
	ID         int    `json:"id"`
	Namabarang string `json:"namabarang"`
	Beratisi   string `json:"beratisi"`
	Harga      string `json:"harga"`
	Stoock     int    `json:"stoock"`
}

var db *sql.DB

func main() {
	connect()

	r := gin.Default()

	r.GET("/api/stock", getStock)
	r.GET("/api/stock/:id", getStockById)
	r.POST("/api/stock", createStock)
	r.PUT("/api//stock/:id", updateStock)
	r.DELETE("/api/stock/:id", deleteStock)

	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

func connect() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/stockbarang")
	if err != nil {
		fmt.Println(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
}

func getStock(c *gin.Context) {
	rows, err := db.Query("SELECT id, namabarang, beratisi, harga, stoock FROM stock")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan mencari barang"})
		return
	}
	defer rows.Close()

	var listbarang []Stock
	for rows.Next() {
		var barang Stock
		err := rows.Scan(&barang.ID, &barang.Namabarang, &barang.Beratisi, &barang.Harga, &barang.Stoock)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan scan barang"})
			return
		}
		listbarang = append(listbarang, barang)
	}

	c.JSON(http.StatusOK, listbarang)
}

func getStockById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id barang tidak valid"})
		return
	}

	var barang Stock
	err = db.QueryRow("SELECT id, namabarang, beratisi, harga, stoock FROM stock WHERE id=?", id).Scan(&barang.ID, &barang.Namabarang, &barang.Beratisi, &barang.Harga, &barang.Stoock)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, barang)
}

func createStock(c *gin.Context) {
	var barang Stock
	if err := c.ShouldBindJSON(&barang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kesalahn input kata"})
		return
	}

	result, err := db.Exec("INSERT INTO stock (namabarang, beratisi, harga, stoock) VALUES (?, ?, ?, ?)", barang.Namabarang, barang.Beratisi, barang.Harga, barang.Stoock)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan input barang"})
		return
	}

	id, _ := result.LastInsertId()
	barang.ID = int(id)

	c.JSON(http.StatusCreated, barang)
}

func updateStock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Barang tidak ditemukan"})
		return
	}

	var barang Stock
	if err := c.ShouldBindJSON(&barang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kesalahan input kata"})
		return
	}

	_, err = db.Exec("UPDATE stock SET namabarang=?, beratisi=?, harga=?, stoock=? WHERE id=?", barang.Namabarang, barang.Beratisi, barang.Harga, barang.Stoock, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan saat update barang"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Barang telah terupdate"})
}

func deleteStock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Barang tidak ditemukan"})
		return
	}

	_, err = db.Exec("DELETE FROM stock WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi error saat menghapus barang"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Barang telah dihapus"})
}
