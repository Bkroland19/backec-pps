package handlers

import (
	"net/http"
	"point-prevalence-survey/database"
	"point-prevalence-survey/models"

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
// @Description Calculate and return all Point Prevalence Survey indicators
// @Tags pps-calculations
// @Accept json
// @Produce json
// @Success 200 {object} PPSIndicators
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/indicators [get]
func (h *PPSCalculationsHandler) GetAllIndicators(c *gin.Context) {
	indicators := PPSIndicators{}

	// Basic Metrics (Grey/White Section)
	var totalAntibiotics int64
	h.db.Model(&models.Antibiotic{}).Count(&totalAntibiotics)
	indicators.TotalAntibioticsPrescribed = int(totalAntibiotics)

	var totalPatientsWithAntibiotics int64
	h.db.Model(&models.Antibiotic{}).Distinct("parent_key").Count(&totalPatientsWithAntibiotics)
	indicators.TotalPatientsWithAntibiotics = int(totalPatientsWithAntibiotics)

	var totalPatients int64
	h.db.Model(&models.Patient{}).Count(&totalPatients)
	indicators.TotalPatientsOnWard = int(totalPatients)

	if indicators.TotalPatientsWithAntibiotics > 0 {
		indicators.AverageAntibioticsPerPatient = float64(indicators.TotalAntibioticsPrescribed) / float64(indicators.TotalPatientsWithAntibiotics)
	}

	// Encounter Metrics (Pink Section)
	if indicators.TotalPatientsOnWard > 0 {
		indicators.PercentageEncounterWithAntibiotic = (float64(indicators.TotalPatientsWithAntibiotics) / float64(indicators.TotalPatientsOnWard)) * 100
	}

	// Injectable Prescriptions
	var injectableCount int64
	h.db.Model(&models.AntibioticDetails{}).Where("intraveno = ? OR intraveno LIKE ?", "iv", "iv%").Count(&injectableCount)
	indicators.TotalInjectablePrescriptions = int(injectableCount)

	// Injectable Percentage
	if indicators.TotalAntibioticsPrescribed > 0 {
		indicators.PercentageInjectablePrescriptions = (float64(indicators.TotalInjectablePrescriptions) / float64(indicators.TotalAntibioticsPrescribed)) * 100
	}

	// Generic Prescriptions (Light Blue/Green Section)
	var genericCount int64
	h.db.Model(&models.Antibiotic{}).Where("antibiotic_inn_name != '' AND antibiotic_inn_name IS NOT NULL").Count(&genericCount)
	indicators.TotalGenericPrescriptions = int(genericCount)

	if indicators.TotalAntibioticsPrescribed > 0 {
		indicators.PercentageGenericPrescriptions = (float64(indicators.TotalGenericPrescriptions) / float64(indicators.TotalAntibioticsPrescribed)) * 100
	}

	// Treatment Guidelines (Yellow Section)
	var guidelineCompliant int64
	h.db.Model(&models.AntibioticDetails{}).Where("guideline = ?", "yes").Count(&guidelineCompliant)
	indicators.TotalGuidelineCompliant = int(guidelineCompliant)

	if indicators.TotalAntibioticsPrescribed > 0 {
		indicators.PercentageGuidelineCompliant = (float64(indicators.TotalGuidelineCompliant) / float64(indicators.TotalAntibioticsPrescribed)) * 100
	}

	// Appropriate Diagnosis (Orange Section)
	var appropriateDiagnosis int64
	h.db.Model(&models.Indication{}).Where("indication_type != '' AND indication_type IS NOT NULL").Count(&appropriateDiagnosis)
	indicators.TotalAppropriateDiagnosis = int(appropriateDiagnosis)

	if indicators.TotalAntibioticsPrescribed > 0 {
		indicators.PercentageAppropriateDiagnosis = (float64(indicators.TotalAppropriateDiagnosis) / float64(indicators.TotalAntibioticsPrescribed)) * 100
	}

	// Culture and Sensitivity (Dark Green Section)
	var cultureBased int64
	h.db.Model(&models.Indication{}).Where("culture_sample_taken = ?", "yes").Count(&cultureBased)
	indicators.TotalCultureBasedPrescriptions = int(cultureBased)

	if indicators.TotalAntibioticsPrescribed > 0 {
		indicators.PercentageCultureBasedPrescriptions = (float64(indicators.TotalCultureBasedPrescriptions) / float64(indicators.TotalAntibioticsPrescribed)) * 100
	}

	var cultureSamples int64
	h.db.Model(&models.Specimen{}).Where("specimen_type != '' AND specimen_type IS NOT NULL").Count(&cultureSamples)
	indicators.TotalCultureSamplesTaken = int(cultureSamples)

	// Missed Doses (Purple Section)
	var missedDoses int64
	h.db.Model(&models.AntibioticDetails{}).Where("number_missed != '' AND number_missed != '0' AND number_missed IS NOT NULL").Count(&missedDoses)
	indicators.TotalMissedDoses = int(missedDoses)

	// Missed Dose Reasons (Light Purple Section)
	var missedDoseReasons []struct {
		MissedDose string `json:"missed_dose"`
		Count      int64  `json:"count"`
	}
	h.db.Model(&models.AntibioticDetails{}).
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
// @Description Calculate basic Point Prevalence Survey metrics
// @Tags pps-calculations
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/basic-metrics [get]
func (h *PPSCalculationsHandler) GetBasicMetrics(c *gin.Context) {
	var totalAntibiotics int64
	var totalPatientsWithAntibiotics int64
	var totalPatients int64

	h.db.Model(&models.Antibiotic{}).Count(&totalAntibiotics)
	h.db.Model(&models.Antibiotic{}).Distinct("parent_key").Count(&totalPatientsWithAntibiotics)
	h.db.Model(&models.Patient{}).Count(&totalPatients)

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

	h.db.Model(&models.AntibioticDetails{}).Where("intraveno = ? OR intraveno LIKE ?", "iv", "iv%").Count(&totalInjectable)
	h.db.Model(&models.Antibiotic{}).Count(&totalAntibiotics)

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
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/generic-metrics [get]
func (h *PPSCalculationsHandler) GetGenericMetrics(c *gin.Context) {
	var totalGeneric int64
	var totalAntibiotics int64

	h.db.Model(&models.Antibiotic{}).Where("antibiotic_inn_name != '' AND antibiotic_inn_name IS NOT NULL").Count(&totalGeneric)
	h.db.Model(&models.Antibiotic{}).Count(&totalAntibiotics)

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
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/guideline-metrics [get]
func (h *PPSCalculationsHandler) GetGuidelineMetrics(c *gin.Context) {
	var totalGuidelineCompliant int64
	var totalAntibiotics int64

	h.db.Model(&models.AntibioticDetails{}).Where("guideline = ?", "yes").Count(&totalGuidelineCompliant)
	h.db.Model(&models.Antibiotic{}).Count(&totalAntibiotics)

	percentageGuideline := 0.0
	if totalAntibiotics > 0 {
		percentageGuideline = (float64(totalGuidelineCompliant) / float64(totalAntibiotics)) * 100
	}

	metrics := gin.H{
		"total_guideline_compliant":      totalGuidelineCompliant,
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

	h.db.Model(&models.Indication{}).Where("indication_type != '' AND indication_type IS NOT NULL").Count(&totalAppropriateDiagnosis)
	h.db.Model(&models.Antibiotic{}).Count(&totalAntibiotics)

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
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/culture-metrics [get]
func (h *PPSCalculationsHandler) GetCultureMetrics(c *gin.Context) {
	var totalCultureBased int64
	var totalCultureSamples int64
	var totalAntibiotics int64

	h.db.Model(&models.Indication{}).Where("culture_sample_taken = ?", "yes").Count(&totalCultureBased)
	h.db.Model(&models.Specimen{}).Where("specimen_type != '' AND specimen_type IS NOT NULL").Count(&totalCultureSamples)
	h.db.Model(&models.Antibiotic{}).Count(&totalAntibiotics)

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
	h.db.Model(&models.AntibioticDetails{}).Where("number_missed != '' AND number_missed != '0' AND number_missed IS NOT NULL").Count(&totalMissedDoses)

	var missedDoseReasons []struct {
		MissedDose string `json:"missed_dose"`
		Count      int64  `json:"count"`
	}
	h.db.Model(&models.AntibioticDetails{}).
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
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/prescriber-metrics [get]
func (h *PPSCalculationsHandler) GetPrescriberMetrics(c *gin.Context) {
	var prescriberStats []struct {
		Prescriber string `json:"prescriber"`
		Count      int64  `json:"count"`
	}
	h.db.Model(&models.AntibioticDetails{}).
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
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/pps/oral-switch-metrics [get]
func (h *PPSCalculationsHandler) GetOralSwitchMetrics(c *gin.Context) {
	var oralSwitchStats []struct {
		OralSwitch string `json:"oral_switch"`
		Count      int64  `json:"count"`
	}
	h.db.Model(&models.AntibioticDetails{}).
		Select("oral_switch, count(*) as count").
		Where("oral_switch != '' AND oral_switch IS NOT NULL").
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
