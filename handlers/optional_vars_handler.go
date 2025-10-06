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

type OptionalVarsHandler struct {
	db *gorm.DB
}

func NewOptionalVarsHandler() *OptionalVarsHandler {
	return &OptionalVarsHandler{
		db: database.GetDB(),
	}
}

// applyFilters applies comprehensive filtering to a query based on URL parameters
func (h *OptionalVarsHandler) applyFilters(db *gorm.DB, c *gin.Context) *gorm.DB {
	// Check if any filtering parameters are provided
	hasFilters := c.Query("start_date") != "" || c.Query("end_date") != "" ||
		c.Query("region") != "" || c.Query("district") != "" || c.Query("subcounty") != "" ||
		c.Query("facility") != "" || c.Query("level") != "" || c.Query("ownership") != ""

	if hasFilters {
		// Join with patients table to filter by patient parameters
		db = db.Joins("JOIN patients ON patients.key = optional_vars.parent_key")

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

// GetOptionalVars godoc
// @Summary Get all optional variables with optional filtering
// @Description Get a list of optional variables with optional filtering by prescriber type, intravenous type, etc.
// @Tags optional-vars
// @Accept json
// @Produce json
// @Param prescriber_type query string false "Filter by prescriber type"
// @Param intravenous_type query string false "Filter by intravenous type"
// @Param guidelines_compliance query string false "Filter by guidelines compliance"
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
// @Router /api/v1/optional-vars [get]
func (h *OptionalVarsHandler) GetOptionalVars(c *gin.Context) {
	var optionalVars []models.OptionalVar
	query := h.db.Model(&models.OptionalVar{})

	// Apply comprehensive filtering (date and geographic)
	query = h.applyFilters(query, c)

	// Apply optional_vars-specific filters
	if prescriberType := c.Query("prescriber_type"); prescriberType != "" {
		query = query.Where("prescriber_type = ?", prescriberType)
	}
	if intravenousType := c.Query("intravenous_type"); intravenousType != "" {
		query = query.Where("intravenous_type = ?", intravenousType)
	}
	if guidelinesCompliance := c.Query("guidelines_compliance"); guidelinesCompliance != "" {
		query = query.Where("guidelines_compliance = ?", guidelinesCompliance)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&optionalVars).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch optional variables"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": optionalVars,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetOptionalVar godoc
// @Summary Get a specific optional variable by ID
// @Description Get detailed information about a specific optional variable
// @Tags optional-vars
// @Accept json
// @Produce json
// @Param id path string true "Optional Variable ID"
// @Success 200 {object} models.OptionalVar
// @Router /api/v1/optional-vars/{id} [get]
func (h *OptionalVarsHandler) GetOptionalVar(c *gin.Context) {
	id := c.Param("id")
	var optionalVar models.OptionalVar

	if err := h.db.First(&optionalVar, "key = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Optional variable not found"})
		return
	}

	c.JSON(http.StatusOK, optionalVar)
}

// CreateOptionalVar godoc
// @Summary Create a new optional variable record
// @Description Create a new optional variable record with the provided data
// @Tags optional-vars
// @Accept json
// @Produce json
// @Param optional_var body models.OptionalVar true "Optional Variable data"
// @Success 201 {object} models.OptionalVar
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/optional-vars [post]
func (h *OptionalVarsHandler) CreateOptionalVar(c *gin.Context) {
	var optionalVar models.OptionalVar

	if err := c.ShouldBindJSON(&optionalVar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&optionalVar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create optional variable"})
		return
	}

	c.JSON(http.StatusCreated, optionalVar)
}

// UpdateOptionalVar godoc
// @Summary Update an optional variable record
// @Description Update an existing optional variable record by ID
// @Tags optional-vars
// @Accept json
// @Produce json
// @Param id path string true "Optional Variable ID"
// @Param optional_var body models.OptionalVar true "Updated Optional Variable data"
// @Success 200 {object} models.OptionalVar
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/optional-vars/{id} [put]
func (h *OptionalVarsHandler) UpdateOptionalVar(c *gin.Context) {
	id := c.Param("id")
	var optionalVar models.OptionalVar

	if err := h.db.First(&optionalVar, "key = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Optional variable not found"})
		return
	}

	if err := c.ShouldBindJSON(&optionalVar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&optionalVar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update optional variable"})
		return
	}

	c.JSON(http.StatusOK, optionalVar)
}

// DeleteOptionalVar godoc
// @Summary Delete an optional variable record
// @Description Delete an optional variable record by ID
// @Tags optional-vars
// @Accept json
// @Produce json
// @Param id path string true "Optional Variable ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/optional-vars/{id} [delete]
func (h *OptionalVarsHandler) DeleteOptionalVar(c *gin.Context) {
	id := c.Param("id")
	var optionalVar models.OptionalVar

	if err := h.db.First(&optionalVar, "key = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Optional variable not found"})
		return
	}

	if err := h.db.Delete(&optionalVar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete optional variable"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetOptionalVarsStats godoc
// @Summary Get optional variables statistics
// @Description Get aggregated statistics about optional variables
// @Tags optional-vars
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
// @Router /api/v1/optional-vars/stats [get]
func (h *OptionalVarsHandler) GetOptionalVarsStats(c *gin.Context) {
	var stats struct {
		TotalOptionalVars int64 `json:"total_optional_vars"`
		ByPrescriberType  []struct {
			PrescriberType string `json:"prescriber_type"`
			Count          int64  `json:"count"`
		} `json:"by_prescriber_type"`
		ByIntravenousType []struct {
			IntravenousType string `json:"intravenous_type"`
			Count           int64  `json:"count"`
		} `json:"by_intravenous_type"`
		ByGuidelinesCompliance []struct {
			GuidelinesCompliance string `json:"guidelines_compliance"`
			Count                int64  `json:"count"`
		} `json:"by_guidelines_compliance"`
		ByTreatmentType []struct {
			TreatmentType string `json:"treatment_type"`
			Count         int64  `json:"count"`
		} `json:"by_treatment_type"`
	}

	// Get base query with filtering
	baseQuery := h.applyFilters(h.db.Model(&models.OptionalVar{}), c)

	// Total optional vars
	baseQuery.Count(&stats.TotalOptionalVars)

	// By prescriber type
	h.applyFilters(h.db.Model(&models.OptionalVar{}), c).Select("prescriber_type, count(*) as count").Group("prescriber_type").Scan(&stats.ByPrescriberType)

	// By intravenous type
	h.applyFilters(h.db.Model(&models.OptionalVar{}), c).Select("intravenous_type, count(*) as count").Group("intravenous_type").Scan(&stats.ByIntravenousType)

	// By guidelines compliance
	h.applyFilters(h.db.Model(&models.OptionalVar{}), c).Select("guidelines_compliance, count(*) as count").Group("guidelines_compliance").Scan(&stats.ByGuidelinesCompliance)

	// By treatment type
	h.applyFilters(h.db.Model(&models.OptionalVar{}), c).Select("treatment_type, count(*) as count").Group("treatment_type").Scan(&stats.ByTreatmentType)

	c.JSON(http.StatusOK, stats)
}

// BulkUploadOptionalVars godoc
// @Summary Bulk upload optional variables from CSV
// @Description Upload multiple optional variables from a CSV file
// @Tags optional-vars
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file containing optional variables data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/optional-vars/upload [post]
func (h *OptionalVarsHandler) BulkUploadOptionalVars(c *gin.Context) {
	// Get the uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Check file type
	if header.Header.Get("Content-Type") != "text/csv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be a CSV"})
		return
	}

	// Process CSV file
	optionalVars, err := h.processOptionalVarsCSV(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to process CSV: " + err.Error()})
		return
	}

	// Bulk insert
	if err := h.db.CreateInBatches(optionalVars, 100).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Optional variables uploaded successfully",
		"count":   len(optionalVars),
	})
}

// processOptionalVarsCSV processes the CSV file and returns a slice of OptionalVar
func (h *OptionalVarsHandler) processOptionalVarsCSV(file interface{}) ([]models.OptionalVar, error) {
	// This is a placeholder implementation
	// In a real implementation, you would parse the CSV file
	// and convert it to OptionalVar structs

	// For now, return an empty slice
	// You can implement CSV parsing using a library like encoding/csv
	return []models.OptionalVar{}, nil
}
