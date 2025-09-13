package handlers

import (
	"net/http"
	"point-prevalence-survey/database"
	"point-prevalence-survey/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AntibioticHandler struct {
	db *gorm.DB
}

func NewAntibioticHandler() *AntibioticHandler {
	return &AntibioticHandler{
		db: database.GetDB(),
	}
}

// GetAntibiotics godoc
// @Summary Get all antibiotics with optional filtering
// @Description Get a list of antibiotics with optional filtering by class, classification, etc.
// @Tags antibiotics
// @Accept json
// @Produce json
// @Param class query string false "Filter by antibiotic class"
// @Param classification query string false "Filter by antibiotic aware classification"
// @Param patient_id query string false "Filter by patient ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/antibiotics [get]
func (h *AntibioticHandler) GetAntibiotics(c *gin.Context) {
	var antibiotics []models.Antibiotic
	query := h.db.Model(&models.Antibiotic{})

	// Apply filters
	if class := c.Query("class"); class != "" {
		query = query.Where("antibiotic_class = ?", class)
	}
	if classification := c.Query("classification"); classification != "" {
		query = query.Where("antibiotic_aware_classification = ?", classification)
	}
	if patientID := c.Query("patient_id"); patientID != "" {
		query = query.Where("parent_key = ?", patientID)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&antibiotics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch antibiotics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": antibiotics,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetAntibiotic godoc
// @Summary Get a specific antibiotic by ID
// @Description Get detailed information about a specific antibiotic
// @Tags antibiotics
// @Accept json
// @Produce json
// @Param id path string true "Antibiotic ID"
// @Success 200 {object} models.Antibiotic
// @Router /api/v1/antibiotics/{id} [get]
func (h *AntibioticHandler) GetAntibiotic(c *gin.Context) {
	id := c.Param("id")
	var antibiotic models.Antibiotic

	if err := h.db.First(&antibiotic, "key = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Antibiotic not found"})
		return
	}

	c.JSON(http.StatusOK, antibiotic)
}

// GetAntibioticStats godoc
// @Summary Get antibiotic statistics
// @Description Get aggregated statistics about antibiotics usage
// @Tags antibiotics
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/antibiotics/stats [get]
func (h *AntibioticHandler) GetAntibioticStats(c *gin.Context) {
	var stats struct {
		TotalAntibiotics int64 `json:"total_antibiotics"`
		ByClass          []struct {
			Class string `json:"class"`
			Count int64  `json:"count"`
		} `json:"by_class"`
		ByClassification []struct {
			Classification string `json:"classification"`
			Count          int64  `json:"count"`
		} `json:"by_classification"`
		ByRoute []struct {
			Route string `json:"route"`
			Count int64  `json:"count"`
		} `json:"by_route"`
		ByFrequency []struct {
			Frequency string `json:"frequency"`
			Count     int64  `json:"count"`
		} `json:"by_frequency"`
	}

	// Total antibiotics
	h.db.Model(&models.Antibiotic{}).Count(&stats.TotalAntibiotics)

	// By class
	h.db.Model(&models.Antibiotic{}).Select("antibiotic_class as class, count(*) as count").Group("antibiotic_class").Scan(&stats.ByClass)

	// By classification
	h.db.Model(&models.Antibiotic{}).Select("antibiotic_aware_classification as classification, count(*) as count").Group("antibiotic_aware_classification").Scan(&stats.ByClassification)

	// By route
	h.db.Model(&models.Antibiotic{}).Select("administration_route as route, count(*) as count").Group("administration_route").Scan(&stats.ByRoute)

	// By frequency
	h.db.Model(&models.Antibiotic{}).Select("unit_dose_frequency as frequency, count(*) as count").Group("unit_dose_frequency").Scan(&stats.ByFrequency)

	c.JSON(http.StatusOK, stats)
}

// GetAntibioticUsageByPatient godoc
// @Summary Get antibiotic usage by patient
// @Description Get detailed antibiotic usage information for a specific patient
// @Tags antibiotics
// @Accept json
// @Produce json
// @Param patient_id path string true "Patient ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/antibiotics/patient/{patient_id} [get]
func (h *AntibioticHandler) GetAntibioticUsageByPatient(c *gin.Context) {
	patientID := c.Param("patient_id")

	var patient models.Patient
	if err := h.db.First(&patient, "key = ?", patientID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	var antibiotics []models.Antibiotic
	if err := h.db.Where("parent_key = ?", patientID).Find(&antibiotics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch antibiotics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"patient":     patient,
		"antibiotics": antibiotics,
	})
}
