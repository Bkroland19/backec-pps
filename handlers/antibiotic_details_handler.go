package handlers

import (
	"net/http"
	"point-prevalence-survey/database"
	"point-prevalence-survey/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AntibioticDetailsHandler struct {
	db *gorm.DB
}

func NewAntibioticDetailsHandler() *AntibioticDetailsHandler {
	return &AntibioticDetailsHandler{
		db: database.GetDB(),
	}
}

// GetAntibioticDetails retrieves all antibiotic details
// @Summary Get all antibiotic details
// @Description Retrieve all antibiotic details from the database
// @Tags antibiotic-details
// @Accept json
// @Produce json
// @Success 200 {array} models.AntibioticDetails
// @Failure 500 {object} map[string]string
// @Router /antibiotic-details [get]
func (h *AntibioticDetailsHandler) GetAntibioticDetails(c *gin.Context) {
	var antibioticDetails []models.AntibioticDetails

	if err := h.db.Find(&antibioticDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve antibiotic details"})
		return
	}

	c.JSON(http.StatusOK, antibioticDetails)
}

// GetAntibioticDetailsByID retrieves a specific antibiotic detail by ID
// @Summary Get antibiotic detail by ID
// @Description Retrieve a specific antibiotic detail by its ID
// @Tags antibiotic-details
// @Accept json
// @Produce json
// @Param id path string true "Antibiotic Detail ID"
// @Success 200 {object} models.AntibioticDetails
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /antibiotic-details/{id} [get]
func (h *AntibioticDetailsHandler) GetAntibioticDetailsByID(c *gin.Context) {
	id := c.Param("id")

	var antibioticDetails models.AntibioticDetails
	if err := h.db.Where("key = ?", id).First(&antibioticDetails).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Antibiotic detail not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve antibiotic detail"})
		return
	}

	c.JSON(http.StatusOK, antibioticDetails)
}

// GetAntibioticDetailsByParentKey retrieves antibiotic details by parent key
// @Summary Get antibiotic details by parent key
// @Description Retrieve antibiotic details for a specific parent patient
// @Tags antibiotic-details
// @Accept json
// @Produce json
// @Param parent_key path string true "Parent Key"
// @Success 200 {array} models.AntibioticDetails
// @Failure 500 {object} map[string]string
// @Router /antibiotic-details/parent/{parent_key} [get]
func (h *AntibioticDetailsHandler) GetAntibioticDetailsByParentKey(c *gin.Context) {
	parentKey := c.Param("parent_key")

	var antibioticDetails []models.AntibioticDetails
	if err := h.db.Where("parent_key = ?", parentKey).Find(&antibioticDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve antibiotic details"})
		return
	}

	c.JSON(http.StatusOK, antibioticDetails)
}

// CreateAntibioticDetails creates a new antibiotic detail
// @Summary Create antibiotic detail
// @Description Create a new antibiotic detail record
// @Tags antibiotic-details
// @Accept json
// @Produce json
// @Param antibioticDetails body models.AntibioticDetails true "Antibiotic Details"
// @Success 201 {object} models.AntibioticDetails
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /antibiotic-details [post]
func (h *AntibioticDetailsHandler) CreateAntibioticDetails(c *gin.Context) {
	var antibioticDetails models.AntibioticDetails

	if err := c.ShouldBindJSON(&antibioticDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&antibioticDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create antibiotic detail"})
		return
	}

	c.JSON(http.StatusCreated, antibioticDetails)
}

// UpdateAntibioticDetails updates an existing antibiotic detail
// @Summary Update antibiotic detail
// @Description Update an existing antibiotic detail record
// @Tags antibiotic-details
// @Accept json
// @Produce json
// @Param id path string true "Antibiotic Detail ID"
// @Param antibioticDetails body models.AntibioticDetails true "Antibiotic Details"
// @Success 200 {object} models.AntibioticDetails
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /antibiotic-details/{id} [put]
func (h *AntibioticDetailsHandler) UpdateAntibioticDetails(c *gin.Context) {
	id := c.Param("id")

	var antibioticDetails models.AntibioticDetails
	if err := h.db.Where("key = ?", id).First(&antibioticDetails).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Antibiotic detail not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve antibiotic detail"})
		return
	}

	if err := c.ShouldBindJSON(&antibioticDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	antibioticDetails.ID = id // Ensure ID doesn't change

	if err := h.db.Save(&antibioticDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update antibiotic detail"})
		return
	}

	c.JSON(http.StatusOK, antibioticDetails)
}

// DeleteAntibioticDetails deletes an antibiotic detail
// @Summary Delete antibiotic detail
// @Description Delete an antibiotic detail record
// @Tags antibiotic-details
// @Accept json
// @Produce json
// @Param id path string true "Antibiotic Detail ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /antibiotic-details/{id} [delete]
func (h *AntibioticDetailsHandler) DeleteAntibioticDetails(c *gin.Context) {
	id := c.Param("id")

	var antibioticDetails models.AntibioticDetails
	if err := h.db.Where("key = ?", id).First(&antibioticDetails).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Antibiotic detail not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve antibiotic detail"})
		return
	}

	if err := h.db.Delete(&antibioticDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete antibiotic detail"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAntibioticDetailsStats retrieves statistics about antibiotic details
// @Summary Get antibiotic details statistics
// @Description Retrieve statistics about antibiotic details in the database
// @Tags antibiotic-details
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /antibiotic-details/stats [get]
func (h *AntibioticDetailsHandler) GetAntibioticDetailsStats(c *gin.Context) {
	var totalCount int64
	var prescriberStats []struct {
		Prescriber string `json:"prescriber"`
		Count      int64  `json:"count"`
	}
	var guidelineStats []struct {
		Guideline string `json:"guideline"`
		Count     int64  `json:"count"`
	}

	// Get total count
	if err := h.db.Model(&models.AntibioticDetails{}).Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total count"})
		return
	}

	// Get prescriber statistics
	if err := h.db.Model(&models.AntibioticDetails{}).
		Select("prescriber, count(*) as count").
		Group("prescriber").
		Find(&prescriberStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get prescriber stats"})
		return
	}

	// Get guideline statistics
	if err := h.db.Model(&models.AntibioticDetails{}).
		Select("guideline, count(*) as count").
		Group("guideline").
		Find(&guidelineStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get guideline stats"})
		return
	}

	stats := gin.H{
		"total_count":      totalCount,
		"prescriber_stats": prescriberStats,
		"guideline_stats":  guidelineStats,
	}

	c.JSON(http.StatusOK, stats)
}
