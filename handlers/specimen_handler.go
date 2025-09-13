package handlers

import (
	"net/http"
	"point-prevalence-survey/database"
	"point-prevalence-survey/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SpecimenHandler struct {
	db *gorm.DB
}

func NewSpecimenHandler() *SpecimenHandler {
	return &SpecimenHandler{
		db: database.GetDB(),
	}
}

// GetSpecimens godoc
// @Summary Get all specimens with optional filtering
// @Description Get a list of specimens with optional filtering by type, result, etc.
// @Tags specimens
// @Accept json
// @Produce json
// @Param type query string false "Filter by specimen type"
// @Param result query string false "Filter by culture result"
// @Param patient_id query string false "Filter by patient ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/specimens [get]
func (h *SpecimenHandler) GetSpecimens(c *gin.Context) {
	var specimens []models.Specimen
	query := h.db.Model(&models.Specimen{})

	// Apply filters
	if specimenType := c.Query("type"); specimenType != "" {
		query = query.Where("specimen_type = ?", specimenType)
	}
	if result := c.Query("result"); result != "" {
		query = query.Where("culture_result = ?", result)
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

	if err := query.Offset(offset).Limit(limit).Find(&specimens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch specimens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": specimens,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetSpecimen godoc
// @Summary Get a specific specimen by ID
// @Description Get detailed information about a specific specimen
// @Tags specimens
// @Accept json
// @Produce json
// @Param id path string true "Specimen ID"
// @Success 200 {object} models.Specimen
// @Router /api/v1/specimens/{id} [get]
func (h *SpecimenHandler) GetSpecimen(c *gin.Context) {
	id := c.Param("id")
	var specimen models.Specimen

	if err := h.db.First(&specimen, "key = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Specimen not found"})
		return
	}

	c.JSON(http.StatusOK, specimen)
}

// GetSpecimenStats godoc
// @Summary Get specimen statistics
// @Description Get aggregated statistics about specimens
// @Tags specimens
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/specimens/stats [get]
func (h *SpecimenHandler) GetSpecimenStats(c *gin.Context) {
	var stats struct {
		TotalSpecimens int64 `json:"total_specimens"`
		ByType         []struct {
			Type  string `json:"type"`
			Count int64  `json:"count"`
		} `json:"by_type"`
		ByResult []struct {
			Result string `json:"result"`
			Count  int64  `json:"count"`
		} `json:"by_result"`
		ByMicroorganism []struct {
			Microorganism string `json:"microorganism"`
			Count         int64  `json:"count"`
		} `json:"by_microorganism"`
		ByResistantPhenotype []struct {
			ResistantPhenotype string `json:"resistant_phenotype"`
			Count              int64  `json:"count"`
		} `json:"by_resistant_phenotype"`
	}

	// Total specimens
	h.db.Model(&models.Specimen{}).Count(&stats.TotalSpecimens)

	// By type
	h.db.Model(&models.Specimen{}).Select("specimen_type as type, count(*) as count").Group("specimen_type").Scan(&stats.ByType)

	// By result
	h.db.Model(&models.Specimen{}).Select("culture_result as result, count(*) as count").Group("culture_result").Scan(&stats.ByResult)

	// By microorganism
	h.db.Model(&models.Specimen{}).Select("microorganism, count(*) as count").Where("microorganism != ''").Group("microorganism").Scan(&stats.ByMicroorganism)

	// By resistant phenotype
	h.db.Model(&models.Specimen{}).Select("resistant_phenotype, count(*) as count").Where("resistant_phenotype != ''").Group("resistant_phenotype").Scan(&stats.ByResistantPhenotype)

	c.JSON(http.StatusOK, stats)
}

// GetSpecimensByPatient godoc
// @Summary Get specimens by patient
// @Description Get all specimens associated with a specific patient
// @Tags specimens
// @Accept json
// @Produce json
// @Param patient_id path string true "Patient ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/specimens/patient/{patient_id} [get]
func (h *SpecimenHandler) GetSpecimensByPatient(c *gin.Context) {
	patientID := c.Param("patient_id")

	var patient models.Patient
	if err := h.db.First(&patient, "key = ?", patientID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	var specimens []models.Specimen
	if err := h.db.Where("parent_key = ?", patientID).Find(&specimens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch specimens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"patient":   patient,
		"specimens": specimens,
	})
}
