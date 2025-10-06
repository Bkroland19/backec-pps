package handlers

import (
	"net/http"
	"point-prevalence-survey/database"
	"point-prevalence-survey/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PPSCalculationsHandler struct {
	db *gorm.DB
}

func NewPPSCalculationsHandler() *PPSCalculationsHandler {
	return &PPSCalculationsHandler{
		db: database.GetDB(),
	}
}

// applyFilters applies comprehensive filtering to a query based on URL parameters
func applyFilters(db *gorm.DB, c *gin.Context, dateColumn string) *gorm.DB {
	// Date filtering
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			// Use PostgreSQL DATE() function to extract just the date part for comparison
			db = db.Where("DATE("+dateColumn+") >= ?", startDate.Format("2006-01-02"))
		}
	}

	if endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// Use PostgreSQL DATE() function to extract just the date part for comparison
			db = db.Where("DATE("+dateColumn+") <= ?", endDate.Format("2006-01-02"))
		}
	}

	// Geographic and facility filtering
	if region := c.Query("region"); region != "" {
		db = db.Where("region = ?", region)
	}

	if district := c.Query("district"); district != "" {
		db = db.Where("district = ?", district)
	}

	if subcounty := c.Query("subcounty"); subcounty != "" {
		db = db.Where("subcounty = ?", subcounty)
	}

	if facility := c.Query("facility"); facility != "" {
		db = db.Where("facility = ?", facility)
	}

	if level := c.Query("level"); level != "" {
		db = db.Where("level_of_care = ?", level)
	}

	if ownership := c.Query("ownership"); ownership != "" {
		db = db.Where("ownership = ?", ownership)
	}

	return db
}

// getFilteredPatientQuery returns a query for patients filtered by date and geographic/facility parameters
func (h *PPSCalculationsHandler) getFilteredPatientQuery(c *gin.Context) *gorm.DB {
	return applyFilters(h.db.Model(&models.Patient{}), c, "submission_date")
}

// getFilteredAntibioticQuery returns a query for antibiotics filtered by parent patient's parameters
func (h *PPSCalculationsHandler) getFilteredAntibioticQuery(c *gin.Context) *gorm.DB {
	query := h.db.Model(&models.Antibiotic{})

	// Check if any filtering parameters are provided
	hasFilters := c.Query("start_date") != "" || c.Query("end_date") != "" ||
		c.Query("region") != "" || c.Query("district") != "" || c.Query("subcounty") != "" ||
		c.Query("facility") != "" || c.Query("level") != "" || c.Query("ownership") != ""

	if hasFilters {
		// Join with patients table to filter by patient parameters
		query = query.Joins("JOIN patients ON patients.key = antibiotics.parent_key")

		// Apply date filtering
		startDateStr := c.Query("start_date")
		endDateStr := c.Query("end_date")

		if startDateStr != "" {
			if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
				query = query.Where("DATE(patients.submission_date) >= ?", startDate.Format("2006-01-02"))
			}
		}

		if endDateStr != "" {
			if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
				query = query.Where("DATE(patients.submission_date) <= ?", endDate.Format("2006-01-02"))
			}
		}

		// Apply geographic and facility filtering
		if region := c.Query("region"); region != "" {
			query = query.Where("patients.region = ?", region)
		}

		if district := c.Query("district"); district != "" {
			query = query.Where("patients.district = ?", district)
		}

		if subcounty := c.Query("subcounty"); subcounty != "" {
			query = query.Where("patients.subcounty = ?", subcounty)
		}

		if facility := c.Query("facility"); facility != "" {
			query = query.Where("patients.facility = ?", facility)
		}

		if level := c.Query("level"); level != "" {
			query = query.Where("patients.level_of_care = ?", level)
		}

		if ownership := c.Query("ownership"); ownership != "" {
			query = query.Where("patients.ownership = ?", ownership)
		}
	}

	return query
}

// applyPatientFiltersToAntibioticDetailsQuery applies comprehensive filtering to antibiotic details query
func (h *PPSCalculationsHandler) applyPatientFiltersToAntibioticDetailsQuery(query *gorm.DB, c *gin.Context) *gorm.DB {
	return h.applyPatientFiltersToQuery(query, c, "antibiotic_details")
}

// applyPatientFiltersToIndicationQuery applies comprehensive filtering to indication query
func (h *PPSCalculationsHandler) applyPatientFiltersToIndicationQuery(query *gorm.DB, c *gin.Context) *gorm.DB {
	return h.applyPatientFiltersToQuery(query, c, "indications")
}

// applyPatientFiltersToSpecimenQuery applies comprehensive filtering to specimen query
func (h *PPSCalculationsHandler) applyPatientFiltersToSpecimenQuery(query *gorm.DB, c *gin.Context) *gorm.DB {
	return h.applyPatientFiltersToQuery(query, c, "specimens")
}

// applyPatientFiltersToQuery applies comprehensive filtering to a query that needs to join with patients table
func (h *PPSCalculationsHandler) applyPatientFiltersToQuery(query *gorm.DB, c *gin.Context, tableName string) *gorm.DB {
	// Check if any filtering parameters are provided
	hasFilters := c.Query("start_date") != "" || c.Query("end_date") != "" ||
		c.Query("region") != "" || c.Query("district") != "" || c.Query("subcounty") != "" ||
		c.Query("facility") != "" || c.Query("level") != "" || c.Query("ownership") != ""

	if hasFilters {
		// Join with patients table to filter by patient parameters
		// For optional_vars, use the key field instead of parent_key
		if tableName == "optional_vars" {
			query = query.Joins("JOIN patients ON patients.key = " + tableName + ".key")
		} else {
			query = query.Joins("JOIN patients ON patients.key = " + tableName + ".parent_key")
		}

		// Apply date filtering
		startDateStr := c.Query("start_date")
		endDateStr := c.Query("end_date")

		if startDateStr != "" {
			if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
				query = query.Where("DATE(patients.submission_date) >= ?", startDate.Format("2006-01-02"))
			}
		}

		if endDateStr != "" {
			if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
				query = query.Where("DATE(patients.submission_date) <= ?", endDate.Format("2006-01-02"))
			}
		}

		// Apply geographic and facility filtering
		if region := c.Query("region"); region != "" {
			query = query.Where("patients.region = ?", region)
		}

		if district := c.Query("district"); district != "" {
			query = query.Where("patients.district = ?", district)
		}

		if subcounty := c.Query("subcounty"); subcounty != "" {
			query = query.Where("patients.subcounty = ?", subcounty)
		}

		if facility := c.Query("facility"); facility != "" {
			query = query.Where("patients.facility = ?", facility)
		}

		if level := c.Query("level"); level != "" {
			query = query.Where("patients.level_of_care = ?", level)
		}

		if ownership := c.Query("ownership"); ownership != "" {
			query = query.Where("patients.ownership = ?", ownership)
		}
	}

	return query
}

// getFilteredAntibioticDetailsQuery returns a query for antibiotic details filtered by parent patient's parameters
func (h *PPSCalculationsHandler) getFilteredAntibioticDetailsQuery(c *gin.Context) *gorm.DB {
	query := h.db.Model(&models.AntibioticDetails{})
	return h.applyPatientFiltersToAntibioticDetailsQuery(query, c)
}

// getFilteredIndicationQuery returns a query for indications filtered by parent patient's parameters
func (h *PPSCalculationsHandler) getFilteredIndicationQuery(c *gin.Context) *gorm.DB {
	query := h.db.Model(&models.Indication{})
	return h.applyPatientFiltersToIndicationQuery(query, c)
}

// getFilteredSpecimenQuery returns a query for specimens filtered by parent patient's parameters
func (h *PPSCalculationsHandler) getFilteredSpecimenQuery(c *gin.Context) *gorm.DB {
	query := h.db.Model(&models.Specimen{})
	return h.applyPatientFiltersToSpecimenQuery(query, c)
}

// getFilteredOptionalVarsQuery returns a query for optional_vars filtered by parent patient's parameters
func (h *PPSCalculationsHandler) getFilteredOptionalVarsQuery(c *gin.Context) *gorm.DB {
	query := h.db.Model(&models.OptionalVar{})

	// Check if any filtering parameters are provided
	hasFilters := c.Query("start_date") != "" || c.Query("end_date") != "" ||
		c.Query("region") != "" || c.Query("district") != "" || c.Query("subcounty") != "" ||
		c.Query("facility") != "" || c.Query("level") != "" || c.Query("ownership") != ""

	if hasFilters {
		return h.applyPatientFiltersToQuery(query, c, "optional_vars")
	}

	// If no filters, return the base query
	return query
}

// PPSIndicators represents all the calculated indicators
type PPSIndicators struct {
	// Grey/White Section - Basic Metrics
	TotalAntibioticsPrescribed   int     `json:"total_antibiotics_prescribed"`
	TotalPatientsWithAntibiotics int     `json:"total_patients_with_antibiotics"`
	AverageAntibioticsPerPatient float64 `json:"average_antibiotics_per_patient"`
	TotalPatientsOnWard          int     `json:"total_patients_on_ward"`

	// Pink Section - Encounter Metrics
	PercentageEncounterWithAntibiotic float64 `json:"percentage_encounter_with_antibiotic"`
	TotalInjectablePrescriptions      int     `json:"total_injectable_prescriptions"`

	// Light Blue Section - Injectable Metrics
	PercentageInjectablePrescriptions float64 `json:"percentage_injectable_prescriptions"`
	TotalGenericPrescriptions         int     `json:"total_generic_prescriptions"`

	// Light Green Section - Generic Name Metrics
	PercentageGenericPrescriptions float64 `json:"percentage_generic_prescriptions"`

	// Yellow Section - Treatment Guidelines
	TotalGuidelineCompliant      int     `json:"total_guideline_compliant"`
	PercentageGuidelineCompliant float64 `json:"percentage_guideline_compliant"`

	// Orange Section - Appropriate Diagnosis
	TotalAppropriateDiagnosis      int     `json:"total_appropriate_diagnosis"`
	PercentageAppropriateDiagnosis float64 `json:"percentage_appropriate_diagnosis"`

	// Dark Green Section - Culture and Sensitivity
	TotalCultureBasedPrescriptions      int     `json:"total_culture_based_prescriptions"`
	PercentageCultureBasedPrescriptions float64 `json:"percentage_culture_based_prescriptions"`
	TotalCultureSamplesTaken            int     `json:"total_culture_samples_taken"`

	// Purple Section - Missed Doses
	TotalMissedDoses int `json:"total_missed_doses"`

	// Light Purple Section - Missed Dose Reasons
	MissedDoseReasons map[string]int `json:"missed_dose_reasons"`
}

// GetAllIndicators calculates all PPS indicators
// @Summary Get all PPS indicators
// @Description Calculate and return all Point Prevalence Survey indicators with optional filtering
// @Tags pps-calculations
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
// @Success 200 {object} PPSIndicators
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/indicators [get]
func (h *PPSCalculationsHandler) GetAllIndicators(c *gin.Context) {
	indicators := PPSIndicators{}

	// Basic Metrics (Grey/White Section) - using filtered queries
	var totalAntibiotics int64
	h.getFilteredAntibioticQuery(c).Count(&totalAntibiotics)
	indicators.TotalAntibioticsPrescribed = int(totalAntibiotics)

	var totalPatientsWithAntibiotics int64
	h.getFilteredAntibioticQuery(c).Distinct("parent_key").Count(&totalPatientsWithAntibiotics)
	indicators.TotalPatientsWithAntibiotics = int(totalPatientsWithAntibiotics)

	var totalPatients int64
	h.getFilteredPatientQuery(c).Count(&totalPatients)
	indicators.TotalPatientsOnWard = int(totalPatients)

	if indicators.TotalPatientsWithAntibiotics > 0 {
		indicators.AverageAntibioticsPerPatient = float64(indicators.TotalAntibioticsPrescribed) / float64(indicators.TotalPatientsWithAntibiotics)
	}

	// Encounter Metrics (Pink Section)
	// Total patients on ward should be total patients, not distinct wards
	// indicators.TotalPatientsOnWard is already set above as total patients

	if indicators.TotalPatientsOnWard > 0 {
		indicators.PercentageEncounterWithAntibiotic = (float64(indicators.TotalPatientsWithAntibiotics) / float64(indicators.TotalPatientsOnWard)) * 100
	}

	// Injectable Prescriptions
	var injectableCount int64
	h.getFilteredAntibioticDetailsQuery(c).Where("intraveno = ? OR intraveno LIKE ?", "iv", "iv%").Count(&injectableCount)
	indicators.TotalInjectablePrescriptions = int(injectableCount)

	// Injectable Percentage
	if indicators.TotalAntibioticsPrescribed > 0 {
		indicators.PercentageInjectablePrescriptions = (float64(indicators.TotalInjectablePrescriptions) / float64(indicators.TotalAntibioticsPrescribed)) * 100
	}

	// Generic Prescriptions (Light Blue/Green Section)
	var genericCount int64
	h.getFilteredAntibioticQuery(c).Where("antibiotic_inn_name != '' AND antibiotic_inn_name IS NOT NULL").Count(&genericCount)
	indicators.TotalGenericPrescriptions = int(genericCount)

	if indicators.TotalAntibioticsPrescribed > 0 {
		indicators.PercentageGenericPrescriptions = (float64(indicators.TotalGenericPrescriptions) / float64(indicators.TotalAntibioticsPrescribed)) * 100
	}

	// Treatment Guidelines (Yellow Section)
	var guidelineCompliant int64
	h.getFilteredAntibioticDetailsQuery(c).Where("guideline = ?", "yes").Count(&guidelineCompliant)
	indicators.TotalGuidelineCompliant = int(guidelineCompliant)

	if indicators.TotalAntibioticsPrescribed > 0 {
		indicators.PercentageGuidelineCompliant = (float64(indicators.TotalGuidelineCompliant) / float64(indicators.TotalAntibioticsPrescribed)) * 100
	}

	// Appropriate Diagnosis (Orange Section)
	var appropriateDiagnosis int64
	h.getFilteredIndicationQuery(c).Where("indication_type != '' AND indication_type IS NOT NULL").Count(&appropriateDiagnosis)
	indicators.TotalAppropriateDiagnosis = int(appropriateDiagnosis)

	if indicators.TotalAntibioticsPrescribed > 0 {
		indicators.PercentageAppropriateDiagnosis = (float64(indicators.TotalAppropriateDiagnosis) / float64(indicators.TotalAntibioticsPrescribed)) * 100
	}

	// Culture and Sensitivity (Dark Green Section)
	var cultureBased int64
	h.getFilteredIndicationQuery(c).Where("culture_sample_taken = ?", "yes").Count(&cultureBased)
	indicators.TotalCultureBasedPrescriptions = int(cultureBased)

	if indicators.TotalAntibioticsPrescribed > 0 {
		indicators.PercentageCultureBasedPrescriptions = (float64(indicators.TotalCultureBasedPrescriptions) / float64(indicators.TotalAntibioticsPrescribed)) * 100
	}

	var cultureSamples int64
	h.getFilteredSpecimenQuery(c).Where("specimen_type != '' AND specimen_type IS NOT NULL").Count(&cultureSamples)
	indicators.TotalCultureSamplesTaken = int(cultureSamples)

	// Missed Doses (Purple Section)
	var missedDoses int64
	h.getFilteredAntibioticDetailsQuery(c).Where("number_missed != '' AND number_missed != '0' AND number_missed IS NOT NULL").Count(&missedDoses)
	indicators.TotalMissedDoses = int(missedDoses)

	// Missed Dose Reasons (Light Purple Section)
	var missedDoseReasons []struct {
		MissedDose string `json:"missed_dose"`
		Count      int64  `json:"count"`
	}
	h.getFilteredAntibioticDetailsQuery(c).
		Select("missed_dose, count(*) as count").
		Where("missed_dose != '' AND missed_dose IS NOT NULL").
		Group("missed_dose").
		Find(&missedDoseReasons)

	indicators.MissedDoseReasons = make(map[string]int)
	for _, reason := range missedDoseReasons {
		indicators.MissedDoseReasons[reason.MissedDose] = int(reason.Count)
	}

	c.JSON(http.StatusOK, indicators)
}

// GetBasicMetrics returns basic PPS metrics
// @Summary Get basic PPS metrics
// @Description Calculate basic Point Prevalence Survey metrics with optional filtering
// @Tags pps-calculations
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
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/basic-metrics [get]
func (h *PPSCalculationsHandler) GetBasicMetrics(c *gin.Context) {
	var totalAntibiotics int64
	var totalPatientsWithAntibiotics int64
	var totalPatients int64

	h.getFilteredAntibioticQuery(c).Count(&totalAntibiotics)
	h.getFilteredAntibioticQuery(c).Distinct("parent_key").Count(&totalPatientsWithAntibiotics)
	// Total patients should be all patients, not distinct wards
	h.getFilteredPatientQuery(c).Count(&totalPatients)

	averageAntibiotics := 0.0
	if totalPatientsWithAntibiotics > 0 {
		averageAntibiotics = float64(totalAntibiotics) / float64(totalPatientsWithAntibiotics)
	}

	percentageEncounter := 0.0
	if totalPatients > 0 {
		percentageEncounter = (float64(totalPatientsWithAntibiotics) / float64(totalPatients)) * 100
	}

	metrics := gin.H{
		"total_antibiotics_prescribed":         totalAntibiotics,
		"total_patients_with_antibiotics":      totalPatientsWithAntibiotics,
		"total_patients_on_ward":               totalPatients,
		"average_antibiotics_per_patient":      averageAntibiotics,
		"percentage_encounter_with_antibiotic": percentageEncounter,
	}

	c.JSON(http.StatusOK, metrics)
}

// GetInjectableMetrics returns injectable antibiotic metrics
// @Summary Get injectable antibiotic metrics
// @Description Calculate injectable antibiotic prescription metrics
// @Tags pps-calculations
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/injectable-metrics [get]
func (h *PPSCalculationsHandler) GetInjectableMetrics(c *gin.Context) {
	var totalInjectable int64
	var totalAntibiotics int64

	h.getFilteredAntibioticDetailsQuery(c).Where("intraveno = ? OR intraveno LIKE ?", "iv", "iv%").Count(&totalInjectable)
	h.getFilteredAntibioticQuery(c).Count(&totalAntibiotics)

	percentageInjectable := 0.0
	if totalAntibiotics > 0 {
		percentageInjectable = (float64(totalInjectable) / float64(totalAntibiotics)) * 100
	}

	metrics := gin.H{
		"total_injectable_prescriptions":      totalInjectable,
		"total_antibiotics_prescribed":        totalAntibiotics,
		"percentage_injectable_prescriptions": percentageInjectable,
	}

	c.JSON(http.StatusOK, metrics)
}

// GetGenericMetrics returns generic name prescription metrics
// @Summary Get generic name prescription metrics
// @Description Calculate generic name antibiotic prescription metrics
// @Tags pps-calculations
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
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/generic-metrics [get]
func (h *PPSCalculationsHandler) GetGenericMetrics(c *gin.Context) {
	var totalGeneric int64
	var totalAntibiotics int64

	h.getFilteredAntibioticQuery(c).Where("antibiotic_inn_name != '' AND antibiotic_inn_name IS NOT NULL").Count(&totalGeneric)
	h.getFilteredAntibioticQuery(c).Count(&totalAntibiotics)

	percentageGeneric := 0.0
	if totalAntibiotics > 0 {
		percentageGeneric = (float64(totalGeneric) / float64(totalAntibiotics)) * 100
	}

	metrics := gin.H{
		"total_generic_prescriptions":      totalGeneric,
		"total_antibiotics_prescribed":     totalAntibiotics,
		"percentage_generic_prescriptions": percentageGeneric,
	}

	c.JSON(http.StatusOK, metrics)
}

// GetGuidelineMetrics returns treatment guideline compliance metrics
// @Summary Get treatment guideline compliance metrics
// @Description Calculate treatment guideline compliance metrics
// @Tags pps-calculations
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
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/guideline-metrics [get]
func (h *PPSCalculationsHandler) GetGuidelineMetrics(c *gin.Context) {
	var totalGuidelineCompliant int64
	var totalOptionalVars int64

	// Check if any filtering parameters are provided
	hasFilters := c.Query("start_date") != "" || c.Query("end_date") != "" ||
		c.Query("region") != "" || c.Query("district") != "" || c.Query("subcounty") != "" ||
		c.Query("facility") != "" || c.Query("level") != "" || c.Query("ownership") != ""

	if hasFilters {
		// Use filtered queries
		compliantQuery := h.getFilteredOptionalVarsQuery(c)
		compliantQuery.Where("guidelines_compliance = ?", "y").Count(&totalGuidelineCompliant)

		totalQuery := h.getFilteredOptionalVarsQuery(c)
		totalQuery.Count(&totalOptionalVars)
	} else {
		// Use direct queries without filters
		h.db.Model(&models.OptionalVar{}).Where("guidelines_compliance = ?", "y").Count(&totalGuidelineCompliant)
		h.db.Model(&models.OptionalVar{}).Count(&totalOptionalVars)
	}

	percentageGuideline := 0.0
	if totalOptionalVars > 0 {
		percentageGuideline = (float64(totalGuidelineCompliant) / float64(totalOptionalVars)) * 100
	}

	// Also get total antibiotics for reference
	var totalAntibiotics int64
	h.getFilteredAntibioticQuery(c).Count(&totalAntibiotics)

	metrics := gin.H{
		"total_guideline_compliant":      totalGuidelineCompliant,
		"total_optional_vars":            totalOptionalVars,
		"total_antibiotics_prescribed":   totalAntibiotics,
		"percentage_guideline_compliant": percentageGuideline,
	}

	c.JSON(http.StatusOK, metrics)
}

// GetDiagnosisMetrics returns appropriate diagnosis metrics
// @Summary Get appropriate diagnosis metrics
// @Description Calculate appropriate diagnosis metrics
// @Tags pps-calculations
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/diagnosis-metrics [get]
func (h *PPSCalculationsHandler) GetDiagnosisMetrics(c *gin.Context) {
	var totalAppropriateDiagnosis int64
	var totalAntibiotics int64

	h.getFilteredIndicationQuery(c).Where("indication_type != '' AND indication_type IS NOT NULL").Count(&totalAppropriateDiagnosis)
	h.getFilteredAntibioticQuery(c).Count(&totalAntibiotics)

	percentageAppropriate := 0.0
	if totalAntibiotics > 0 {
		percentageAppropriate = (float64(totalAppropriateDiagnosis) / float64(totalAntibiotics)) * 100
	}

	metrics := gin.H{
		"total_appropriate_diagnosis":      totalAppropriateDiagnosis,
		"total_antibiotics_prescribed":     totalAntibiotics,
		"percentage_appropriate_diagnosis": percentageAppropriate,
	}

	c.JSON(http.StatusOK, metrics)
}

// GetCultureMetrics returns culture and sensitivity metrics
// @Summary Get culture and sensitivity metrics
// @Description Calculate culture and sensitivity metrics
// @Tags pps-calculations
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
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/culture-metrics [get]
func (h *PPSCalculationsHandler) GetCultureMetrics(c *gin.Context) {
	var totalCultureBased int64
	var totalCultureSamples int64
	var totalAntibiotics int64

	h.getFilteredIndicationQuery(c).Where("culture_sample_taken = ?", "yes").Count(&totalCultureBased)
	h.getFilteredSpecimenQuery(c).Where("specimen_type != '' AND specimen_type IS NOT NULL").Count(&totalCultureSamples)
	h.getFilteredAntibioticQuery(c).Count(&totalAntibiotics)

	percentageCultureBased := 0.0
	if totalAntibiotics > 0 {
		percentageCultureBased = (float64(totalCultureBased) / float64(totalAntibiotics)) * 100
	}

	metrics := gin.H{
		"total_culture_based_prescriptions":      totalCultureBased,
		"total_culture_samples_taken":            totalCultureSamples,
		"total_antibiotics_prescribed":           totalAntibiotics,
		"percentage_culture_based_prescriptions": percentageCultureBased,
	}

	c.JSON(http.StatusOK, metrics)
}

// GetMissedDoseMetrics returns missed dose metrics
// @Summary Get missed dose metrics
// @Description Calculate missed dose metrics and reasons
// @Tags pps-calculations
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/missed-dose-metrics [get]
func (h *PPSCalculationsHandler) GetMissedDoseMetrics(c *gin.Context) {
	var totalMissedDoses int64
	h.getFilteredAntibioticDetailsQuery(c).Where("number_missed != '' AND number_missed != '0' AND number_missed IS NOT NULL").Count(&totalMissedDoses)

	var missedDoseReasons []struct {
		MissedDose string `json:"missed_dose"`
		Count      int64  `json:"count"`
	}
	h.getFilteredAntibioticDetailsQuery(c).
		Select("missed_dose, count(*) as count").
		Where("missed_dose != '' AND missed_dose IS NOT NULL").
		Group("missed_dose").
		Find(&missedDoseReasons)

	reasonsMap := make(map[string]int64)
	for _, reason := range missedDoseReasons {
		reasonsMap[reason.MissedDose] = reason.Count
	}

	metrics := gin.H{
		"total_missed_doses":  totalMissedDoses,
		"missed_dose_reasons": reasonsMap,
	}

	c.JSON(http.StatusOK, metrics)
}

// GetPrescriberMetrics returns prescriber metrics
// @Summary Get prescriber metrics
// @Description Calculate prescriber-related metrics
// @Tags pps-calculations
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
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/prescriber-metrics [get]
func (h *PPSCalculationsHandler) GetPrescriberMetrics(c *gin.Context) {
	var prescriberStats []struct {
		Prescriber string `json:"prescriber"`
		Count      int64  `json:"count"`
	}
	h.getFilteredAntibioticDetailsQuery(c).
		Select("prescriber, count(*) as count").
		Where("prescriber != '' AND prescriber IS NOT NULL").
		Group("prescriber").
		Find(&prescriberStats)

	prescriberMap := make(map[string]int64)
	for _, stat := range prescriberStats {
		prescriberMap[stat.Prescriber] = stat.Count
	}

	metrics := gin.H{
		"prescriber_stats": prescriberMap,
	}

	c.JSON(http.StatusOK, metrics)
}

// GetOralSwitchMetrics returns oral switch metrics
// @Summary Get oral switch metrics
// @Description Calculate oral switch metrics
// @Tags pps-calculations
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
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/oral-switch-metrics [get]
func (h *PPSCalculationsHandler) GetOralSwitchMetrics(c *gin.Context) {
	var oralSwitchStats []struct {
		OralSwitch string `json:"oral_switch"`
		Count      int64  `json:"count"`
	}
	h.getFilteredAntibioticDetailsQuery(c).
		Select("oral_switch, count(*) as count").
		Where("oral_switch != '' AND oral_switch IS NOT NULL AND oral_switch = 'yes'").
		Group("oral_switch").
		Find(&oralSwitchStats)

	oralSwitchMap := make(map[string]int64)
	for _, stat := range oralSwitchStats {
		oralSwitchMap[stat.OralSwitch] = stat.Count
	}

	metrics := gin.H{
		"oral_switch_stats": oralSwitchMap,
	}

	c.JSON(http.StatusOK, metrics)
}

// GetLongStayPatients calculates the number of patients staying longer than 7 days
// @Summary Get patients staying longer than 7 days
// @Description Calculate the number of patients who have been staying for more than 7 days based on survey date minus admission date
// @Tags pps-calculations
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
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/long-stay-patients [get]
func (h *PPSCalculationsHandler) GetLongStayPatients(c *gin.Context) {
	var longStayCount int64

	// Calculate patients staying longer than 7 days
	// Using SQL to calculate the difference between survey_date and admission_date
	err := h.getFilteredPatientQuery(c).
		Where("survey_date - admission_date > INTERVAL '7 days'").
		Count(&longStayCount).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate long stay patients"})
		return
	}

	// Also get total patients for percentage calculation
	var totalPatients int64
	err = h.getFilteredPatientQuery(c).Count(&totalPatients).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total patients count"})
		return
	}

	// Calculate percentage
	percentage := 0.0
	if totalPatients > 0 {
		percentage = (float64(longStayCount) / float64(totalPatients)) * 100
	}

	metrics := gin.H{
		"patients_staying_longer_than_7_days": longStayCount,
		"total_patients":                      totalPatients,
		"percentage_long_stay":                percentage,
		"description":                         "Patients staying longer than 7 days (survey_date - admission_date > 7 days)",
	}

	c.JSON(http.StatusOK, metrics)
}

// GetAWaReCategorization calculates WHO AWaRe antibiotic categorization
// @Summary Get WHO AWaRe antibiotic categorization
// @Description Calculate the distribution of antibiotics based on WHO AWaRe classification (Access, Watch, Reserve, Unclassified) with optional filtering
// @Tags pps-calculations
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
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/aware-categorization [get]
func (h *PPSCalculationsHandler) GetAWaReCategorization(c *gin.Context) {
	// Get total antibiotics count
	var totalAntibiotics int64
	err := h.getFilteredAntibioticQuery(c).Count(&totalAntibiotics).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total antibiotics count"})
		return
	}

	// Get Access antibiotics count (first choice antibiotics)
	var accessCount int64
	err = h.getFilteredAntibioticQuery(c).
		Where("antibiotic_aware_classification = ?", "Access").
		Count(&accessCount).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Access antibiotics count"})
		return
	}

	// Get Watch antibiotics count (second choice antibiotics)
	var watchCount int64
	err = h.getFilteredAntibioticQuery(c).
		Where("antibiotic_aware_classification = ?", "Watch").
		Count(&watchCount).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Watch antibiotics count"})
		return
	}

	// Get Reserve antibiotics count (last resort antibiotics)
	var reserveCount int64
	err = h.getFilteredAntibioticQuery(c).
		Where("antibiotic_aware_classification = ?", "Reserve").
		Count(&reserveCount).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Reserve antibiotics count"})
		return
	}

	// Compute Unclassified as everything not Access/Watch/Reserve (case-insensitive)
	unclassifiedCount := totalAntibiotics - (accessCount + watchCount + reserveCount)

	// Calculate percentages
	accessPercentage := 0.0
	watchPercentage := 0.0
	reservePercentage := 0.0
	unclassifiedPercentage := 0.0

	if totalAntibiotics > 0 {
		accessPercentage = (float64(accessCount) / float64(totalAntibiotics)) * 100
		watchPercentage = (float64(watchCount) / float64(totalAntibiotics)) * 100
		reservePercentage = (float64(reserveCount) / float64(totalAntibiotics)) * 100
		unclassifiedPercentage = (float64(unclassifiedCount) / float64(totalAntibiotics)) * 100
	}

	metrics := gin.H{
		"total_antibiotics": totalAntibiotics,
		"access": gin.H{
			"count":       accessCount,
			"percentage":  accessPercentage,
			"description": "First choice antibiotics",
		},
		"watch": gin.H{
			"count":       watchCount,
			"percentage":  watchPercentage,
			"description": "Second choice antibiotics",
		},
		"reserve": gin.H{
			"count":       reserveCount,
			"percentage":  reservePercentage,
			"description": "Last resort antibiotics",
		},
		"unclassified": gin.H{
			"count":       unclassifiedCount,
			"percentage":  unclassifiedPercentage,
			"description": "Not categorized",
		},
	}

	c.JSON(http.StatusOK, metrics)
}

// GetAppropriateDiagnosisPercentage returns percentage of "yes" responses for appropriate diagnosis
// @Summary Get appropriate diagnosis percentage
// @Description Percentage = count of "yes" responses divided by total responses (yes/no) from indications.reason_in_notes
// @Tags pps-calculations
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
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/appropriate-diagnosis [get]
func (h *PPSCalculationsHandler) GetAppropriateDiagnosisPercentage(c *gin.Context) {
	var totalResponses int64
	var yesResponses int64

	if err := h.getFilteredIndicationQuery(c).
		Where("reason_in_notes IS NOT NULL AND reason_in_notes != ''").
		Count(&totalResponses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count total responses"})
		return
	}

	if err := h.getFilteredIndicationQuery(c).
		Where("LOWER(reason_in_notes) = ?", "yes").
		Count(&yesResponses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count yes responses"})
		return
	}

	percentage := 0.0
	if totalResponses > 0 {
		percentage = (float64(yesResponses) / float64(totalResponses)) * 100
	}

	metrics := gin.H{
		"appropriate_diagnosis_yes":        yesResponses,
		"appropriate_diagnosis_total":      totalResponses,
		"percentage_appropriate_diagnosis": percentage,
	}

	c.JSON(http.StatusOK, metrics)
}
