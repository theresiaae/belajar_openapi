package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	Age       uint      `gorm:"column:age"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func main() {
	dsn := "root:@tcp(localhost:3306)/db_openapi?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database")
	}

	router := gin.Default()

	router.GET("/users", func(c *gin.Context) {
		var users []User
		db.Find(&users)
		c.JSON(http.StatusOK, gin.H{"data": users})
	})

	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": user})
	})

	router.POST("/users", func(c *gin.Context) {
		var input User
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user := User{Name: input.Name, Email: input.Email, Age: input.Age}
		db.Create(&user)
		c.JSON(http.StatusCreated, gin.H{"data": user})
	})

	router.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		var input User
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		db.Model(&user).Updates(User{Name: input.Name, Email: input.Email, Age: input.Age})
		c.JSON(http.StatusOK, gin.H{"data": user})
	})

	router.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var user User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		db.Delete(&user)
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	})

	router.Run(":3000")
}
