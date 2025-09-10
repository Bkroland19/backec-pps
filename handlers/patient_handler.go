package handlers

import (
	"net/http"
	"point-prevalence-survey/database"
	"point-prevalence-survey/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	db *gorm.DB
}

func NewPatientHandler() *PatientHandler {
	return &PatientHandler{
		db: database.GetDB(),
	}
}

// GetPatients godoc
// @Summary Get all patients with optional filtering
// @Description Get a list of patients with optional filtering by region, district, facility, etc.
// @Tags patients
// @Accept json
// @Produce json
// @Param region query string false "Filter by region"
// @Param district query string false "Filter by district"
// @Param facility query string false "Filter by facility"
// @Param ward query string false "Filter by ward name"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/patients [get]
func (h *PatientHandler) GetPatients(c *gin.Context) {
	var patients []models.Patient
	query := h.db.Model(&models.Patient{})

	// Apply filters
	if region := c.Query("region"); region != "" {
		query = query.Where("region = ?", region)
	}
	if district := c.Query("district"); district != "" {
		query = query.Where("district = ?", district)
	}
	if facility := c.Query("facility"); facility != "" {
		query = query.Where("facility = ?", facility)
	}
	if ward := c.Query("ward"); ward != "" {
		query = query.Where("ward_name = ?", ward)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&patients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch patients"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": patients,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetPatient godoc
// @Summary Get a specific patient by ID
// @Description Get detailed information about a specific patient including all related data
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Success 200 {object} models.Patient
// @Router /api/v1/patients/{id} [get]
func (h *PatientHandler) GetPatient(c *gin.Context) {
	id := c.Param("id")
	var patient models.Patient

	if err := h.db.Preload("Antibiotics").Preload("Indications").Preload("OptionalVars").Preload("Specimens").First(&patient, "key = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	c.JSON(http.StatusOK, patient)
}

// GetPatientAntibiotics godoc
// @Summary Get antibiotics for a specific patient
// @Description Get all antibiotics associated with a specific patient
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Success 200 {object} []models.Antibiotic
// @Router /api/v1/patients/{id}/antibiotics [get]
func (h *PatientHandler) GetPatientAntibiotics(c *gin.Context) {
	id := c.Param("id")
	var antibiotics []models.Antibiotic

	if err := h.db.Where("parent_key = ?", id).Find(&antibiotics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch antibiotics"})
		return
	}

	c.JSON(http.StatusOK, antibiotics)
}

// GetPatientIndications godoc
// @Summary Get indications for a specific patient
// @Description Get all indications associated with a specific patient
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Success 200 {object} []models.Indication
// @Router /api/v1/patients/{id}/indications [get]
func (h *PatientHandler) GetPatientIndications(c *gin.Context) {
	id := c.Param("id")
	var indications []models.Indication

	if err := h.db.Where("parent_key = ?", id).Find(&indications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch indications"})
		return
	}

	c.JSON(http.StatusOK, indications)
}

// GetPatientOptionalVars godoc
// @Summary Get optional variables for a specific patient
// @Description Get all optional variables associated with a specific patient
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Success 200 {object} []models.OptionalVar
// @Router /api/v1/patients/{id}/optional-vars [get]
func (h *PatientHandler) GetPatientOptionalVars(c *gin.Context) {
	id := c.Param("id")
	var optionalVars []models.OptionalVar

	if err := h.db.Where("parent_key = ?", id).Find(&optionalVars).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch optional variables"})
		return
	}

	c.JSON(http.StatusOK, optionalVars)
}

// GetPatientSpecimens godoc
// @Summary Get specimens for a specific patient
// @Description Get all specimens associated with a specific patient
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Success 200 {object} []models.Specimen
// @Router /api/v1/patients/{id}/specimens [get]
func (h *PatientHandler) GetPatientSpecimens(c *gin.Context) {
	id := c.Param("id")
	var specimens []models.Specimen

	if err := h.db.Where("parent_key = ?", id).Find(&specimens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch specimens"})
		return
	}

	c.JSON(http.StatusOK, specimens)
}

// GetPatientStats godoc
// @Summary Get patient statistics
// @Description Get aggregated statistics about patients
// @Tags patients
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/patients/stats [get]
func (h *PatientHandler) GetPatientStats(c *gin.Context) {
	var stats struct {
		TotalPatients        int64 `json:"total_patients"`
		PatientsOnAntibiotic int64 `json:"patients_on_antibiotic"`
		ByRegion             []struct {
			Region string `json:"region"`
			Count  int64  `json:"count"`
		} `json:"by_region"`
		ByFacility []struct {
			Facility string `json:"facility"`
			Count    int64  `json:"count"`
		} `json:"by_facility"`
		ByWard []struct {
			Ward  string `json:"ward"`
			Count int64  `json:"count"`
		} `json:"by_ward"`
	}

	// Total patients
	h.db.Model(&models.Patient{}).Count(&stats.TotalPatients)

	// Patients on antibiotic
	h.db.Model(&models.Patient{}).Where("patient_on_antibiotic = ?", "yes").Count(&stats.PatientsOnAntibiotic)

	// By region
	h.db.Model(&models.Patient{}).Select("region, count(*) as count").Group("region").Scan(&stats.ByRegion)

	// By facility
	h.db.Model(&models.Patient{}).Select("facility, count(*) as count").Group("facility").Scan(&stats.ByFacility)

	// By ward
	h.db.Model(&models.Patient{}).Select("ward_name as ward, count(*) as count").Group("ward_name").Scan(&stats.ByWard)

	c.JSON(http.StatusOK, stats)
}
