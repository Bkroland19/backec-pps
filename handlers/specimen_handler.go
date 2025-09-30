package handlers

import (
	"net/http"
	"point-prevalence-survey/database"
	"point-prevalence-survey/models"
	"strconv"
	"time"

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

// applyFilters applies comprehensive filtering to a query based on URL parameters
func (h *SpecimenHandler) applyFilters(db *gorm.DB, c *gin.Context) *gorm.DB {
	// Check if any filtering parameters are provided
	hasFilters := c.Query("start_date") != "" || c.Query("end_date") != "" ||
		c.Query("region") != "" || c.Query("district") != "" || c.Query("subcounty") != "" ||
		c.Query("facility") != "" || c.Query("level") != "" || c.Query("ownership") != ""

	if hasFilters {
		// Join with patients table to filter by patient parameters
		db = db.Joins("JOIN patients ON patients.key = specimens.parent_key")

		// Apply date filtering
		startDateStr := c.Query("start_date")
		endDateStr := c.Query("end_date")

		if startDateStr != "" {
			if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
				db = db.Where("DATE(patients.submission_date) >= ?", startDate.Format("2006-01-02"))
			}
		}

		if endDateStr != "" {
			if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
				db = db.Where("DATE(patients.submission_date) <= ?", endDate.Format("2006-01-02"))
			}
		}

		// Apply geographic and facility filtering
		if region := c.Query("region"); region != "" {
			db = db.Where("patients.region = ?", region)
		}

		if district := c.Query("district"); district != "" {
			db = db.Where("patients.district = ?", district)
		}

		if subcounty := c.Query("subcounty"); subcounty != "" {
			db = db.Where("patients.subcounty = ?", subcounty)
		}

		if facility := c.Query("facility"); facility != "" {
			db = db.Where("patients.facility = ?", facility)
		}

		if level := c.Query("level"); level != "" {
			db = db.Where("patients.level_of_care = ?", level)
		}

		if ownership := c.Query("ownership"); ownership != "" {
			db = db.Where("patients.ownership = ?", ownership)
		}
	}

	return db
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
// @Param start_date query string false "Start date for filtering (YYYY-MM-DD)"
// @Param end_date query string false "End date for filtering (YYYY-MM-DD)"
// @Param region query string false "Region for filtering"
// @Param district query string false "District for filtering"
// @Param subcounty query string false "Subcounty for filtering"
// @Param facility query string false "Facility for filtering"
// @Param level query string false "Level of care for filtering"
// @Param ownership query string false "Ownership for filtering"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/specimens [get]
func (h *SpecimenHandler) GetSpecimens(c *gin.Context) {
	var specimens []models.Specimen
	query := h.db.Model(&models.Specimen{})

	// Apply comprehensive filtering (date and geographic)
	query = h.applyFilters(query, c)

	// Apply specimen-specific filters
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
// @Param start_date query string false "Start date for filtering (YYYY-MM-DD)"
// @Param end_date query string false "End date for filtering (YYYY-MM-DD)"
// @Param region query string false "Region for filtering"
// @Param district query string false "District for filtering"
// @Param subcounty query string false "Subcounty for filtering"
// @Param facility query string false "Facility for filtering"
// @Param level query string false "Level of care for filtering"
// @Param ownership query string false "Ownership for filtering"
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

	// Get base query with filtering
	baseQuery := h.applyFilters(h.db.Model(&models.Specimen{}), c)

	// Total specimens
	baseQuery.Count(&stats.TotalSpecimens)

	// By type
	h.applyFilters(h.db.Model(&models.Specimen{}), c).Select("specimen_type as type, count(*) as count").Group("specimen_type").Scan(&stats.ByType)

	// By result
	h.applyFilters(h.db.Model(&models.Specimen{}), c).Select("culture_result as result, count(*) as count").Group("culture_result").Scan(&stats.ByResult)

	// By microorganism
	h.applyFilters(h.db.Model(&models.Specimen{}), c).Select("microorganism, count(*) as count").Where("microorganism != ''").Group("microorganism").Scan(&stats.ByMicroorganism)

	// By resistant phenotype
	h.applyFilters(h.db.Model(&models.Specimen{}), c).Select("resistant_phenotype, count(*) as count").Where("resistant_phenotype != ''").Group("resistant_phenotype").Scan(&stats.ByResistantPhenotype)

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
