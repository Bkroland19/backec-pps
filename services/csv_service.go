package services

import (
	"encoding/csv"
	"fmt"
	"log"
	"mime/multipart"
	"point-prevalence-survey/database"
	"point-prevalence-survey/models"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type CSVService struct {
	db *gorm.DB
}

// UploadResult contains statistics about the upload process
type UploadResult struct {
	TotalRecords     int      `json:"total_records"`
	ProcessedRecords int      `json:"processed_records"`
	SkippedRecords   int      `json:"skipped_records"`
	InsertedRecords  int      `json:"inserted_records"`
	UpdatedRecords   int      `json:"updated_records"`
	Errors           []string `json:"errors"`
}

func NewCSVService() *CSVService {
	return &CSVService{
		db: database.GetDB(),
	}
}

func (s *CSVService) ImportPatients(file multipart.File) (*UploadResult, error) {
	result := &UploadResult{
		Errors: make([]string, 0),
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return result, fmt.Errorf("error reading CSV file: %v", err)
	}

	if len(records) < 2 {
		return result, fmt.Errorf("CSV file must have at least a header row and one data row")
	}

	result.TotalRecords = len(records) - 1 // Exclude header row

	// Skip header row
	for i, record := range records[1:] {
		rowNum := i + 2 // Account for header row
		result.ProcessedRecords++

		if len(record) < 56 { // Ensure we have enough columns (56 fields total)
			errorMsg := fmt.Sprintf("Row %d: insufficient columns (expected at least 56, got %d)", rowNum, len(record))
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		patient := s.parsePatientRecord(record)
		if patient.ID == "" {
			errorMsg := fmt.Sprintf("Row %d: missing patient ID", rowNum)
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		// Check if patient already exists
		var existingPatient models.Patient
		err := s.db.Where("key = ?", patient.ID).First(&existingPatient).Error
		if err == nil {
			// Patient exists, skip
			result.SkippedRecords++
			log.Printf("Skipping patient %s: already exists", patient.ID)
			continue
		} else if err != gorm.ErrRecordNotFound {
			// Database error
			errorMsg := fmt.Sprintf("Row %d: database error checking patient %s: %v", rowNum, patient.ID, err)
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		// Create new patient
		if err := s.db.Create(&patient).Error; err != nil {
			errorMsg := fmt.Sprintf("Row %d: error creating patient %s: %v", rowNum, patient.ID, err)
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		result.InsertedRecords++
		log.Printf("Successfully created patient %s", patient.ID)
	}

	return result, nil
}

func (s *CSVService) ImportAntibiotics(file multipart.File) (*UploadResult, error) {
	result := &UploadResult{
		Errors: make([]string, 0),
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return result, fmt.Errorf("error reading CSV file: %v", err)
	}

	if len(records) < 2 {
		return result, fmt.Errorf("CSV file must have at least a header row and one data row")
	}

	result.TotalRecords = len(records) - 1 // Exclude header row

	// Skip header row
	for i, record := range records[1:] {
		rowNum := i + 2 // Account for header row
		result.ProcessedRecords++

		if len(record) < 15 {
			errorMsg := fmt.Sprintf("Row %d: insufficient columns (expected at least 15, got %d)", rowNum, len(record))
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		antibiotic := s.parseAntibioticRecord(record)
		if antibiotic.ID == "" {
			errorMsg := fmt.Sprintf("Row %d: missing antibiotic ID", rowNum)
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		// Check if antibiotic already exists
		var existingAntibiotic models.Antibiotic
		err := s.db.Where("key = ?", antibiotic.ID).First(&existingAntibiotic).Error
		if err == nil {
			// Antibiotic exists, skip
			result.SkippedRecords++
			log.Printf("Skipping antibiotic %s: already exists", antibiotic.ID)
			continue
		} else if err != gorm.ErrRecordNotFound {
			// Database error
			errorMsg := fmt.Sprintf("Row %d: database error checking antibiotic %s: %v", rowNum, antibiotic.ID, err)
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		// Check if parent patient exists
		if antibiotic.ParentKey != "" {
			var parentPatient models.Patient
			err := s.db.Where("key = ?", antibiotic.ParentKey).First(&parentPatient).Error
			if err != nil {
				// Log the ParentKey being searched for debugging
				log.Printf("Looking for parent patient with key: '%s'", antibiotic.ParentKey)

				// Also log a few existing patient keys for comparison
				var samplePatients []models.Patient
				s.db.Select("key").Limit(3).Find(&samplePatients)
				log.Printf("Sample patient keys in database: %v", samplePatients)

				errorMsg := fmt.Sprintf("Row %d: parent patient %s not found for antibiotic %s", rowNum, antibiotic.ParentKey, antibiotic.ID)
				result.Errors = append(result.Errors, errorMsg)
				result.SkippedRecords++
				continue
			}
		}

		// Create new antibiotic
		if err := s.db.Create(&antibiotic).Error; err != nil {
			errorMsg := fmt.Sprintf("Row %d: error creating antibiotic %s: %v", rowNum, antibiotic.ID, err)
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		result.InsertedRecords++
		log.Printf("Successfully created antibiotic %s", antibiotic.ID)
	}

	return result, nil
}

func (s *CSVService) ImportAntibioticDetails(file multipart.File) (*UploadResult, error) {
	result := &UploadResult{
		Errors: make([]string, 0),
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to read CSV file: %v", err))
		return result, err
	}

	if len(records) < 2 {
		result.Errors = append(result.Errors, "CSV file must contain at least a header row and one data row")
		return result, nil
	}

	// Skip header row
	for i, record := range records[1:] {
		rowNum := i + 2 // Account for header row

		if len(record) < 8 {
			errorMsg := fmt.Sprintf("Row %d: insufficient columns (expected at least 8, got %d)", rowNum, len(record))
			result.Errors = append(result.Errors, errorMsg)
			continue
		}

		antibioticDetails := s.parseAntibioticDetailsRecord(record)
		if antibioticDetails.ID == "" {
			errorMsg := fmt.Sprintf("Row %d: missing antibiotic details ID", rowNum)
			result.Errors = append(result.Errors, errorMsg)
			continue
		}

		// Check if parent patient exists
		if antibioticDetails.ParentKey != "" {
			var parentPatient models.Patient
			err := s.db.Where("key = ?", antibioticDetails.ParentKey).First(&parentPatient).Error
			if err != nil {
				// Log the ParentKey being searched for debugging
				log.Printf("Looking for parent patient with key: '%s'", antibioticDetails.ParentKey)

				// Also log a few existing patient keys for comparison
				var samplePatients []models.Patient
				s.db.Select("key").Limit(3).Find(&samplePatients)
				log.Printf("Sample patient keys in database: %v", samplePatients)

				errorMsg := fmt.Sprintf("Row %d: parent patient %s not found for antibiotic details %s", rowNum, antibioticDetails.ParentKey, antibioticDetails.ID)
				result.Errors = append(result.Errors, errorMsg)
				continue
			}
		}

		// Check if record already exists
		var existingRecord models.AntibioticDetails
		err = s.db.Where("key = ?", antibioticDetails.ID).First(&existingRecord).Error
		if err == nil {
			// Record exists, skip it
			continue
		}

		// Create new record
		if err := s.db.Create(&antibioticDetails).Error; err != nil {
			errorMsg := fmt.Sprintf("Row %d: failed to create antibiotic details record: %v", rowNum, err)
			result.Errors = append(result.Errors, errorMsg)
			continue
		}
	}

	return result, nil
}

func (s *CSVService) ImportIndications(file multipart.File) (*UploadResult, error) {
	result := &UploadResult{
		Errors: make([]string, 0),
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return result, fmt.Errorf("error reading CSV file: %v", err)
	}

	if len(records) < 2 {
		return result, fmt.Errorf("CSV file must have at least a header row and one data row")
	}

	result.TotalRecords = len(records) - 1 // Exclude header row

	// Skip header row
	for i, record := range records[1:] {
		rowNum := i + 2 // Account for header row
		result.ProcessedRecords++

		if len(record) < 9 {
			errorMsg := fmt.Sprintf("Row %d: insufficient columns (expected at least 9, got %d)", rowNum, len(record))
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		indication := s.parseIndicationRecord(record)
		if indication.ID == "" {
			errorMsg := fmt.Sprintf("Row %d: missing indication ID", rowNum)
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		// Check if indication already exists
		var existingIndication models.Indication
		err := s.db.Where("key = ?", indication.ID).First(&existingIndication).Error
		if err == nil {
			// Indication exists, skip
			result.SkippedRecords++
			log.Printf("Skipping indication %s: already exists", indication.ID)
			continue
		} else if err != gorm.ErrRecordNotFound {
			// Database error
			errorMsg := fmt.Sprintf("Row %d: database error checking indication %s: %v", rowNum, indication.ID, err)
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		// Check if parent patient exists
		if indication.ParentKey != "" {
			var parentPatient models.Patient
			err := s.db.Where("key = ?", indication.ParentKey).First(&parentPatient).Error
			if err != nil {
				errorMsg := fmt.Sprintf("Row %d: parent patient %s not found for indication %s", rowNum, indication.ParentKey, indication.ID)
				result.Errors = append(result.Errors, errorMsg)
				result.SkippedRecords++
				continue
			}
		}

		// Create new indication
		if err := s.db.Create(&indication).Error; err != nil {
			errorMsg := fmt.Sprintf("Row %d: error creating indication %s: %v", rowNum, indication.ID, err)
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		result.InsertedRecords++
		log.Printf("Successfully created indication %s", indication.ID)
	}

	return result, nil
}

func (s *CSVService) ImportOptionalVars(file multipart.File) (*UploadResult, error) {
	result := &UploadResult{
		Errors: make([]string, 0),
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return result, fmt.Errorf("error reading CSV file: %v", err)
	}

	if len(records) < 2 {
		return result, fmt.Errorf("CSV file must have at least a header row and one data row")
	}

	result.TotalRecords = len(records) - 1 // Exclude header row

	// Skip header row
	for i, record := range records[1:] {
		rowNum := i + 2 // Account for header row
		result.ProcessedRecords++

		if len(record) < 8 {
			errorMsg := fmt.Sprintf("Row %d: insufficient columns (expected at least 8, got %d)", rowNum, len(record))
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		optionalVar := s.parseOptionalVarRecord(record)
		if optionalVar.ID == "" {
			errorMsg := fmt.Sprintf("Row %d: missing optional var ID", rowNum)
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		// Check if optional var already exists
		var existingOptionalVar models.OptionalVar
		err := s.db.Where("key = ?", optionalVar.ID).First(&existingOptionalVar).Error
		if err == nil {
			// Optional var exists, skip
			result.SkippedRecords++
			log.Printf("Skipping optional var %s: already exists", optionalVar.ID)
			continue
		} else if err != gorm.ErrRecordNotFound {
			// Database error
			errorMsg := fmt.Sprintf("Row %d: database error checking optional var %s: %v", rowNum, optionalVar.ID, err)
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		// Check if parent patient exists
		if optionalVar.ParentKey != "" {
			var parentPatient models.Patient
			err := s.db.Where("key = ?", optionalVar.ParentKey).First(&parentPatient).Error
			if err != nil {
				errorMsg := fmt.Sprintf("Row %d: parent patient %s not found for optional var %s", rowNum, optionalVar.ParentKey, optionalVar.ID)
				result.Errors = append(result.Errors, errorMsg)
				result.SkippedRecords++
				continue
			}
		}

		// Create new optional var
		if err := s.db.Create(&optionalVar).Error; err != nil {
			errorMsg := fmt.Sprintf("Row %d: error creating optional var %s: %v", rowNum, optionalVar.ID, err)
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		result.InsertedRecords++
		log.Printf("Successfully created optional var %s", optionalVar.ID)
	}

	return result, nil
}

func (s *CSVService) ImportSpecimens(file multipart.File) (*UploadResult, error) {
	result := &UploadResult{
		Errors: make([]string, 0),
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return result, fmt.Errorf("error reading CSV file: %v", err)
	}

	if len(records) < 2 {
		return result, fmt.Errorf("CSV file must have at least a header row and one data row")
	}

	result.TotalRecords = len(records) - 1 // Exclude header row

	// Skip header row
	for i, record := range records[1:] {
		rowNum := i + 2 // Account for header row
		result.ProcessedRecords++

		if len(record) < 7 {
			errorMsg := fmt.Sprintf("Row %d: insufficient columns (expected at least 7, got %d)", rowNum, len(record))
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		specimen := s.parseSpecimenRecord(record)
		if specimen.ID == "" {
			errorMsg := fmt.Sprintf("Row %d: missing specimen ID", rowNum)
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		// Check if specimen already exists
		var existingSpecimen models.Specimen
		err := s.db.Where("key = ?", specimen.ID).First(&existingSpecimen).Error
		if err == nil {
			// Specimen exists, skip
			result.SkippedRecords++
			log.Printf("Skipping specimen %s: already exists", specimen.ID)
			continue
		} else if err != gorm.ErrRecordNotFound {
			// Database error
			errorMsg := fmt.Sprintf("Row %d: database error checking specimen %s: %v", rowNum, specimen.ID, err)
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		// Check if parent patient exists
		if specimen.ParentKey != "" {
			var parentPatient models.Patient
			err := s.db.Where("key = ?", specimen.ParentKey).First(&parentPatient).Error
			if err != nil {
				errorMsg := fmt.Sprintf("Row %d: parent patient %s not found for specimen %s", rowNum, specimen.ParentKey, specimen.ID)
				result.Errors = append(result.Errors, errorMsg)
				result.SkippedRecords++
				continue
			}
		}

		// Create new specimen
		if err := s.db.Create(&specimen).Error; err != nil {
			errorMsg := fmt.Sprintf("Row %d: error creating specimen %s: %v", rowNum, specimen.ID, err)
			result.Errors = append(result.Errors, errorMsg)
			result.SkippedRecords++
			continue
		}

		result.InsertedRecords++
		log.Printf("Successfully created specimen %s", specimen.ID)
	}

	return result, nil
}

// Helper function to parse dates with multiple format support
func (s *CSVService) parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}

	// Try different date formats
	formats := []string{
		"2006-01-02T15:04:05.000Z",    // ISO 8601 with milliseconds
		"2006-01-02T15:04:05Z",        // ISO 8601 without milliseconds
		"2006-01-02T15:04:05",         // ISO 8601 without timezone
		"2006-01-02 15:04:05",         // Standard datetime
		"2006-01-02",                  // Date only
		"01/02/2006",                  // US format
		"02/01/2006",                  // European format
		"2006/01/02",                  // Alternative format
		"2006-01-02 15:04:05.000",     // With milliseconds
		"2006-01-02 15:04:05.000000",  // With microseconds
		"2006-01-02T15:04:05.000000Z", // ISO with microseconds
		"2006-01-02T15:04:05.000000",  // ISO with microseconds no timezone
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}

	// If all formats fail, return zero time
	return time.Time{}
}

// Helper functions to parse CSV records
func (s *CSVService) parsePatientRecord(record []string) models.Patient {
	patient := models.Patient{}

	// CSV Column 0 = submission_date
	if len(record) > 0 && record[0] != "" {
		patient.SubmissionDate = s.parseDate(record[0])
	}
	// CSV Column 1 = region
	if len(record) > 1 && record[1] != "" {
		patient.Region = record[1]
	}
	// CSV Column 2 = district
	if len(record) > 2 && record[2] != "" {
		patient.District = record[2]
	}
	// CSV Column 3 = subcounty
	if len(record) > 3 && record[3] != "" {
		patient.Subcounty = record[3]
	}
	// CSV Column 4 = facility
	if len(record) > 4 && record[4] != "" {
		patient.Facility = record[4]
	}
	// CSV Column 5 = level_of_care
	if len(record) > 5 && record[5] != "" {
		patient.LevelOfCare = record[5]
	}
	// CSV Column 6 = ownership
	if len(record) > 6 && record[6] != "" {
		patient.Ownership = record[6]
	}
	// CSV Column 7 = ward_name
	if len(record) > 7 && record[7] != "" {
		patient.WardName = record[7]
	}
	// CSV Column 8 = ward_total_patients
	if len(record) > 8 && record[8] != "" {
		if val, err := strconv.Atoi(record[8]); err == nil {
			patient.WardTotalPatients = val
		}
	}
	// CSV Column 9 = ward_eligible_patients
	if len(record) > 9 && record[9] != "" {
		if val, err := strconv.Atoi(record[9]); err == nil {
			patient.WardEligiblePatients = val
		}
	}
	// CSV Column 10 = survey_date
	if len(record) > 10 && record[10] != "" {
		patient.SurveyDate = s.parseDate(record[10])
	}
	// CSV Column 11 = patient_initials
	if len(record) > 11 && record[11] != "" {
		patient.PatientInitials = record[11]
	}
	// CSV Column 12 = code
	if len(record) > 12 && record[12] != "" {
		patient.Code = record[12]
	}
	// CSV Column 13 = rand_num
	if len(record) > 13 && record[13] != "" {
		if val, err := strconv.Atoi(record[13]); err == nil {
			patient.RandNum = val
		}
	}
	// CSV Column 14 = patient_code
	if len(record) > 14 && record[14] != "" {
		patient.PatientCode = record[14]
	}
	// CSV Column 15 = show_code
	if len(record) > 15 && record[15] != "" {
		patient.ShowCode = record[15]
	}
	// CSV Column 16 = is_the_patient_an_infant
	if len(record) > 16 && record[16] != "" {
		patient.IsThePatientAnInfant = record[16]
	}
	// CSV Column 17 = age_months
	if len(record) > 17 && record[17] != "" {
		if val, err := strconv.Atoi(record[17]); err == nil {
			patient.AgeMonths = val
		}
	}
	// CSV Column 18 = age_years
	if len(record) > 18 && record[18] != "" {
		if val, err := strconv.Atoi(record[18]); err == nil {
			patient.AgeYears = val
		}
	}
	// CSV Column 19 = pre_term_birth
	if len(record) > 19 && record[19] != "" {
		patient.PreTermBirth = record[19]
	}
	// CSV Column 20 = gender
	if len(record) > 20 && record[20] != "" {
		patient.Gender = record[20]
	}
	// CSV Column 21 = weight
	if len(record) > 21 && record[21] != "" {
		if val, err := strconv.ParseFloat(record[21], 64); err == nil {
			patient.Weight = val
		}
	}
	// CSV Column 22 = weight_birth_kg
	if len(record) > 22 && record[22] != "" {
		if val, err := strconv.ParseFloat(record[22], 64); err == nil {
			patient.WeightBirthKg = val
		}
	}
	// CSV Column 23 = admission_date
	if len(record) > 23 && record[23] != "" {
		patient.AdmissionDate = s.parseDate(record[23])
	}
	// CSV Column 24 = surgery_since_admission
	if len(record) > 24 && record[24] != "" {
		patient.SurgerySinceAdmission = record[24]
	}
	// CSV Column 25 = urinary_catheter
	if len(record) > 25 && record[25] != "" {
		patient.UrinaryCatheter = record[25]
	}
	// CSV Column 26 = peripheral_vascular_catheter
	if len(record) > 26 && record[26] != "" {
		patient.PeripheralVascularCatheter = record[26]
	}
	// CSV Column 27 = central_vascular_catheter
	if len(record) > 27 && record[27] != "" {
		patient.CentralVascularCatheter = record[27]
	}
	// CSV Column 28 = intubation
	if len(record) > 28 && record[28] != "" {
		patient.Intubation = record[28]
	}
	// CSV Column 29 = patient_on_antibiotic
	if len(record) > 29 && record[29] != "" {
		patient.PatientOnAntibiotic = record[29]
	}
	// CSV Column 30 = patient_number_antibiotics
	if len(record) > 30 && record[30] != "" {
		if val, err := strconv.Atoi(record[30]); err == nil {
			patient.PatientNumberAntibiotics = val
		}
	}
	// CSV Column 31 = malaria_status
	if len(record) > 31 && record[31] != "" {
		patient.MalariaStatus = record[31]
	}
	// CSV Column 32 = tuberculosis_status
	if len(record) > 32 && record[32] != "" {
		patient.TuberculosisStatus = record[32]
	}
	// CSV Column 33 = hiv_status
	if len(record) > 33 && record[33] != "" {
		patient.HIVStatus = record[33]
	}
	// CSV Column 34 = hiv_on_art
	if len(record) > 34 && record[34] != "" {
		patient.HIVOnART = record[34]
	}
	// CSV Column 35 = hiv_cd4_count
	if len(record) > 35 && record[35] != "" {
		patient.HIVCD4Count = record[35]
	}
	// CSV Column 36 = hiv_viral_load
	if len(record) > 36 && record[36] != "" {
		patient.HIVViralLoad = record[36]
	}
	// CSV Column 37 = diabetes
	if len(record) > 37 && record[37] != "" {
		patient.Diabetes = record[37]
	}
	// CSV Column 38 = malnutrition_status
	if len(record) > 38 && record[38] != "" {
		patient.MalnutritionStatus = record[38]
	}
	// CSV Column 39 = hypertension
	if len(record) > 39 && record[39] != "" {
		patient.Hypertension = record[39]
	}
	// CSV Column 40 = referred_from
	if len(record) > 40 && record[40] != "" {
		patient.ReferredFrom = record[40]
	}
	// CSV Column 41 = hospitalization_90_days
	if len(record) > 41 && record[41] != "" {
		patient.Hospitalization90Days = record[41]
	}
	// CSV Column 42 = type_surgery_since_admission
	if len(record) > 42 && record[42] != "" {
		patient.TypeSurgerySinceAdmission = record[42]
	}
	// CSV Column 43 = additional_comment
	if len(record) > 43 && record[43] != "" {
		patient.AdditionalComment = record[43]
	}
	// CSV Column 44 = comments
	if len(record) > 44 && record[44] != "" {
		patient.Comments = record[44]
	}
	// CSV Column 45 = instance_id (this is the key field!)
	if len(record) > 45 && record[45] != "" {
		patient.ID = record[45]
		patient.InstanceID = record[45]
	}
	// CSV Column 46 = submitter_id
	if len(record) > 46 && record[46] != "" {
		patient.SubmitterID = record[46]
	}
	// CSV Column 47 = submitter_name
	if len(record) > 47 && record[47] != "" {
		patient.SubmitterName = record[47]
	}
	// CSV Column 48 = attachments_present
	if len(record) > 48 && record[48] != "" {
		patient.AttachmentsPresent = record[48]
	}
	// CSV Column 49 = attachments_expected
	if len(record) > 49 && record[49] != "" {
		patient.AttachmentsExpected = record[49]
	}
	// CSV Column 50 = status
	if len(record) > 50 && record[50] != "" {
		patient.Status = record[50]
	}
	// CSV Column 51 = review_state
	if len(record) > 51 && record[51] != "" {
		patient.ReviewState = record[51]
	}
	// CSV Column 52 = device_id
	if len(record) > 52 && record[52] != "" {
		patient.DeviceID = record[52]
	}
	// CSV Column 53 = edits
	if len(record) > 53 && record[53] != "" {
		patient.Edits = record[53]
	}
	// CSV Column 54 = form_version
	if len(record) > 54 && record[54] != "" {
		patient.FormVersion = record[54]
	}

	return patient
}

func (s *CSVService) parseAntibioticRecord(record []string) models.Antibiotic {
	antibiotic := models.Antibiotic{}

	// Debug: Log the first few columns to understand the structure
	first10 := record
	if len(record) > 10 {
		first10 = record[:10]
	}
	log.Printf("Antibiotic CSV Record structure - First 10 columns: %v", first10)
	log.Printf("Antibiotic Total columns: %d", len(record))

	// Based on your CSV structure:
	// Column 14 (index 13): PARENT_K
	// Column 15 (index 14): KEY (this should be the antibiotic ID)
	
	// Set the antibiotic ID from the KEY column (column 15, index 14)
	if len(record) > 14 && record[14] != "" {
		// Extract the UUID part from the KEY (remove /Antibioticform/Core_variables[X])
		key := record[14]
		if strings.Contains(key, "/Antibioticform/") {
			// Extract just the UUID part before /Antibioticform/
			parts := strings.Split(key, "/Antibioticform/")
			if len(parts) > 0 {
				antibiotic.ID = parts[0]
			}
		} else {
			antibiotic.ID = key
		}
	}
	
	// Set the ParentKey from the PARENT_K column (column 14, index 13)
	if len(record) > 13 && record[13] != "" {
		// Extract the UUID part from the ParentKey (remove /Antibioticform/Core_variables[X])
		parentKey := record[13]
		if strings.Contains(parentKey, "/Antibioticform/") {
			// Extract just the UUID part before /Antibioticform/
			parts := strings.Split(parentKey, "/Antibioticform/")
			if len(parts) > 0 {
				antibiotic.ParentKey = parts[0]
			}
		} else {
			antibiotic.ParentKey = parentKey
		}
	}
	// Based on your CSV structure:
	// Column 1: Empty
	// Column 2: Antibiotic name (e.g., "Ceftriaxone")
	// Column 3: Other_An (empty)
	// Column 4: atc_code
	// Column 5: antibiotic class
	// Column 6: antibiotic classification
	// Column 7: Empty
	// Column 8: StartDateAntibiotic
	// Column 9: UnitDose
	// Column 10: UnitDoses
	// Column 11: UnitDosel (frequency)
	// Column 12: UnitDosef (administration route)
	// Column 13: Empty
	// Column 14: PARENT_K
	// Column 15: KEY
	
	if len(record) > 1 && record[1] != "" {
		antibiotic.AntibioticINNName = record[1] // Column 2: Antibiotic name
	}
	if len(record) > 2 && record[2] != "" {
		antibiotic.OtherAntibiotic = record[2] // Column 3: Other_An
	}
	if len(record) > 3 && record[3] != "" {
		antibiotic.ATCCode = record[3] // Column 4: atc_code
	}
	if len(record) > 4 && record[4] != "" {
		antibiotic.AntibioticClass = record[4] // Column 5: antibiotic class
	}
	if len(record) > 5 && record[5] != "" {
		antibiotic.AntibioticAwareClassification = record[5] // Column 6: antibiotic classification
	}
	if len(record) > 6 && record[6] != "" {
		antibiotic.AntibioticWrittenInINN = record[6] // Column 7: Empty in your data
	}
	if len(record) > 7 && record[7] != "" {
		antibiotic.StartDateAntibiotic = s.parseDate(record[7]) // Column 8: StartDateAntibiotic
	}
	if len(record) > 8 && record[8] != "" {
		if val, err := strconv.ParseFloat(record[8], 64); err == nil {
			antibiotic.UnitDose = val // Column 9: UnitDose
		}
	}
	if len(record) > 9 && record[9] != "" {
		antibiotic.UnitDosesCombination = record[9] // Column 10: UnitDoses
	}
	if len(record) > 10 && record[10] != "" {
		antibiotic.UnitDoseMeasureUnit = record[10] // Column 11: UnitDosel (frequency)
	}
	if len(record) > 11 && record[11] != "" {
		antibiotic.UnitDoseFrequency = record[11] // Column 12: UnitDosef (administration route)
	}
	if len(record) > 12 && record[12] != "" {
		antibiotic.AdministrationRoute = record[12] // Column 13: Empty in your data
	}
	// ParentKey is now set as the ID above, so we don't need to set it again

	return antibiotic
}

func (s *CSVService) parseAntibioticDetailsRecord(record []string) models.AntibioticDetails {
	antibioticDetails := models.AntibioticDetails{}

	// Debug: Log the first few columns to understand the structure
	first10 := record
	if len(record) > 10 {
		first10 = record[:10]
	}
	log.Printf("AntibioticDetails CSV Record structure - First 10 columns: %v", first10)
	log.Printf("AntibioticDetails Total columns: %d", len(record))

	// Based on the image, the fields are in these positions:
	// Column 1: Prescriber
	// Column 2: Intraveno
	// Column 3: OralSwitch
	// Column 4: NumberMissed
	// Column 5: MissedDose
	// Column 6: Guideline
	// Column 7: Treatment
	// Column 8: ParentKey (used as ID)

	if len(record) > 0 && record[0] != "" {
		antibioticDetails.Prescriber = record[0]
	}
	if len(record) > 1 && record[1] != "" {
		antibioticDetails.Intraveno = record[1]
	}
	if len(record) > 2 && record[2] != "" {
		antibioticDetails.OralSwitch = record[2]
	}
	if len(record) > 3 && record[3] != "" {
		antibioticDetails.NumberMissed = record[3]
	}
	if len(record) > 4 && record[4] != "" {
		antibioticDetails.MissedDose = record[4]
	}
	if len(record) > 5 && record[5] != "" {
		antibioticDetails.Guideline = record[5]
	}
	if len(record) > 6 && record[6] != "" {
		antibioticDetails.Treatment = record[6]
	}
	if len(record) > 7 && record[7] != "" {
		// Extract the UUID part from the ParentKey (remove /Antibioticform/Core_variables[X])
		parentKey := record[7]
		if strings.Contains(parentKey, "/Antibioticform/") {
			// Extract just the UUID part before /Antibioticform/
			parts := strings.Split(parentKey, "/Antibioticform/")
			if len(parts) > 0 {
				antibioticDetails.ID = parts[0]
				antibioticDetails.ParentKey = parts[0]
			}
		} else {
			antibioticDetails.ID = parentKey
			antibioticDetails.ParentKey = parentKey
		}
	}

	return antibioticDetails
}

func (s *CSVService) parseIndicationRecord(record []string) models.Indication {
	indication := models.Indication{}

	// Debug: Log the first few columns to understand the structure
	first10 := record
	if len(record) > 10 {
		first10 = record[:10]
	}
	log.Printf("Indication CSV Record structure - First 10 columns: %v", first10)
	log.Printf("Indication Total columns: %d", len(record))

	// Based on the CSV structure:
	// Column 1 (index 0): Indication
	// Column 2 (index 1): Surg.Prop (empty)
	// Column 3 (index 2): Surg.Prop (empty)
	// Column 4 (index 3): Diagnosis
	// Column 5 (index 4): StartDateTreatment
	// Column 6 (index 5): ReasonIn!
	// Column 7 (index 6): CultureSa
	// Column 8 (index 7): PARENT_k
	// Column 9 (index 8): KEY

	if len(record) > 0 && record[0] != "" {
		indication.IndicationType = record[0] // Indication field
	}
	if len(record) > 1 && record[1] != "" {
		indication.SurgProphDuration = record[1] // First Surg.Prop
	}
	if len(record) > 2 && record[2] != "" {
		indication.SurgProphSite = record[2] // Second Surg.Prop
	}
	if len(record) > 3 && record[3] != "" {
		indication.Diagnosis = record[3]
	}
	if len(record) > 4 && record[4] != "" {
		indication.StartDateTreatment = s.parseDate(record[4])
	}
	if len(record) > 5 && record[5] != "" {
		indication.ReasonInNotes = record[5] // ReasonIn!
	}
	if len(record) > 6 && record[6] != "" {
		indication.CultureSampleTaken = record[6] // CultureSa
	}
	if len(record) > 7 && record[7] != "" {
		// Extract the UUID part from the ParentKey (remove /Antibioticform/Core_variables[X])
		parentKey := record[7]
		if strings.Contains(parentKey, "/Antibioticform/") {
			// Extract just the UUID part before /Antibioticform/
			parts := strings.Split(parentKey, "/Antibioticform/")
			if len(parts) > 0 {
				indication.ParentKey = parts[0]
			}
		} else {
			indication.ParentKey = parentKey
		}
	}
	if len(record) > 8 && record[8] != "" {
		// Extract the UUID part from the KEY (remove /Antibioticform/Core_variables[X])
		key := record[8]
		if strings.Contains(key, "/Antibioticform/") {
			// Extract just the UUID part before /Antibioticform/
			parts := strings.Split(key, "/Antibioticform/")
			if len(parts) > 0 {
				indication.ID = parts[0]
			}
		} else {
			indication.ID = key
		}
	}

	return indication
}

func (s *CSVService) parseOptionalVarRecord(record []string) models.OptionalVar {
	optionalVar := models.OptionalVar{}

	// First column is the key field
	if len(record) > 0 && record[0] != "" {
		optionalVar.ID = record[0]
	}
	if len(record) > 1 && record[1] != "" {
		optionalVar.PrescriberType = record[1]
	}
	if len(record) > 2 && record[2] != "" {
		optionalVar.IntravenousType = record[2]
	}
	if len(record) > 3 && record[3] != "" {
		optionalVar.OralSwitch = record[3]
	}
	if len(record) > 4 && record[4] != "" {
		if val, err := strconv.Atoi(record[4]); err == nil {
			optionalVar.NumberMissedDoses = val
		}
	}
	if len(record) > 5 && record[5] != "" {
		optionalVar.MissedDosesReason = record[5]
	}
	if len(record) > 6 && record[6] != "" {
		optionalVar.GuidelinesCompliance = record[6]
	}
	if len(record) > 7 && record[7] != "" {
		optionalVar.TreatmentType = record[7]
	}
	if len(record) > 8 && record[8] != "" {
		optionalVar.ParentKey = record[8]
	}

	return optionalVar
}

func (s *CSVService) parseSpecimenRecord(record []string) models.Specimen {
	specimen := models.Specimen{}

	// Debug: Log the first few columns to understand the structure
	first10 := record
	if len(record) > 10 {
		first10 = record[:10]
	}
	log.Printf("Specimen CSV Record structure - First 10 columns: %v", first10)
	log.Printf("Specimen Total columns: %d", len(record))

	// Based on the CSV structure:
	// Column 1 (index 0): Specimen
	// Column 2 (index 1): CultureRe
	// Column 3 (index 2): Microorga
	// Column 4 (index 3): Antibiotic
	// Column 5 (index 4): Resistant
	// Column 6 (index 5): PARENT_k
	// Column 7 (index 6): KEY

	if len(record) > 0 && record[0] != "" {
		specimen.SpecimenType = record[0] // Specimen field
	}
	if len(record) > 1 && record[1] != "" {
		specimen.CultureResult = record[1] // CultureRe
	}
	if len(record) > 2 && record[2] != "" {
		specimen.Microorganism = record[2] // Microorga
	}
	if len(record) > 3 && record[3] != "" {
		specimen.AntibioticSusceptibilityTestResults = record[3] // Antibiotic
	}
	if len(record) > 4 && record[4] != "" {
		specimen.ResistantPhenotype = record[4] // Resistant
	}
	if len(record) > 5 && record[5] != "" {
		// Extract the UUID part from the ParentKey (remove /Antibioticform/Core_variables[X])
		parentKey := record[5]
		if strings.Contains(parentKey, "/Antibioticform/") {
			// Extract just the UUID part before /Antibioticform/
			parts := strings.Split(parentKey, "/Antibioticform/")
			if len(parts) > 0 {
				specimen.ParentKey = parts[0]
			}
		} else {
			specimen.ParentKey = parentKey
		}
	}
	if len(record) > 6 && record[6] != "" {
		// Extract the UUID part from the KEY (remove /Antibioticform/Core_variables[X])
		key := record[6]
		if strings.Contains(key, "/Antibioticform/") {
			// Extract just the UUID part before /Antibioticform/
			parts := strings.Split(key, "/Antibioticform/")
			if len(parts) > 0 {
				specimen.ID = parts[0]
			}
		} else {
			specimen.ID = key
		}
	}

	return specimen
}
