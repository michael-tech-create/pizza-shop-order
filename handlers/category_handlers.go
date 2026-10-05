package handlers

import (
	"log"
	"net/http"
	"strconv"

	"pizza-app/models"
	"pizza-app/repositories"

	"github.com/gin-gonic/gin"
)

// GET /api/categories — public, used to populate the menu filter
func GetCategoriesHandler(c *gin.Context) {
	categories, err := repositories.GetAllCategories()
	if err != nil {
		log.Printf("GetAllCategories error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// POST /api/categories — admin only
func CreateCategoryHandler(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category name is required"})
		return
	}

	cat, err := repositories.CreateCategory(req.Name, req.ImageURL)
	if err != nil {
		log.Printf("CreateCategory error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create category — name may already exist"})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

// PUT /api/categories/:id — admin only
func UpdateCategoryHandler(c *gin.Context) {
	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category name is required"})
		return
	}

	cat, err := repositories.UpdateCategory(idInt, req.Name, req.ImageURL)
	if err != nil {
		if err == repositories.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		log.Printf("UpdateCategory error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update category"})
		return
	}
	c.JSON(http.StatusOK, cat)
}

// DELETE /api/categories/:id — admin only
func DeleteCategoryHandler(c *gin.Context) {
	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	if err := repositories.DeleteCategory(idInt); err != nil {
		if err == repositories.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		log.Printf("DeleteCategory error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}
