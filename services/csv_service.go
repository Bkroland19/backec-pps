package services

import (
	"encoding/csv"
	"fmt"
	"io"
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

func NewCSVService() *CSVService {
	return &CSVService{
		db: database.GetDB(),
	}
}

func (s *CSVService) ImportPatients(file multipart.File) error {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %v", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file must have at least a header row and one data row")
	}

	// Skip header row
	for i, record := range records[1:] {
		if len(record) < 50 { // Ensure we have enough columns
			log.Printf("Skipping row %d: insufficient columns", i+2)
			continue
		}

		patient := s.parsePatientRecord(record)
		if patient.ID == "" {
			log.Printf("Skipping row %d: missing patient ID", i+2)
			continue
		}

		// Create or update patient
		if err := s.db.Save(&patient).Error; err != nil {
			log.Printf("Error saving patient %s: %v", patient.ID, err)
			continue
		}
	}

	return nil
}

func (s *CSVService) ImportAntibiotics(file multipart.File) error {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %v", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file must have at least a header row and one data row")
	}

	// Skip header row
	for i, record := range records[1:] {
		if len(record) < 15 {
			log.Printf("Skipping row %d: insufficient columns", i+2)
			continue
		}

		antibiotic := s.parseAntibioticRecord(record)
		if antibiotic.ID == "" {
			log.Printf("Skipping row %d: missing antibiotic ID", i+2)
			continue
		}

		// Create or update antibiotic
		if err := s.db.Save(&antibiotic).Error; err != nil {
			log.Printf("Error saving antibiotic %s: %v", antibiotic.ID, err)
			continue
		}
	}

	return nil
}

func (s *CSVService) ImportIndications(file multipart.File) error {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %v", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file must have at least a header row and one data row")
	}

	// Skip header row
	for i, record := range records[1:] {
		if len(record) < 9 {
			log.Printf("Skipping row %d: insufficient columns", i+2)
			continue
		}

		indication := s.parseIndicationRecord(record)
		if indication.ID == "" {
			log.Printf("Skipping row %d: missing indication ID", i+2)
			continue
		}

		// Create or update indication
		if err := s.db.Save(&indication).Error; err != nil {
			log.Printf("Error saving indication %s: %v", indication.ID, err)
			continue
		}
	}

	return nil
}

func (s *CSVService) ImportOptionalVars(file multipart.File) error {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %v", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file must have at least a header row and one data row")
	}

	// Skip header row
	for i, record := range records[1:] {
		if len(record) < 8 {
			log.Printf("Skipping row %d: insufficient columns", i+2)
			continue
		}

		optionalVar := s.parseOptionalVarRecord(record)
		if optionalVar.ID == "" {
			log.Printf("Skipping row %d: missing optional var ID", i+2)
			continue
		}

		// Create or update optional var
		if err := s.db.Save(&optionalVar).Error; err != nil {
			log.Printf("Error saving optional var %s: %v", optionalVar.ID, err)
			continue
		}
	}

	return nil
}

func (s *CSVService) ImportSpecimens(file multipart.File) error {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %v", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file must have at least a header row and one data row")
	}

	// Skip header row
	for i, record := range records[1:] {
		if len(record) < 6 {
			log.Printf("Skipping row %d: insufficient columns", i+2)
			continue
		}

		specimen := s.parseSpecimenRecord(record)
		if specimen.ID == "" {
			log.Printf("Skipping row %d: missing specimen ID", i+2)
			continue
		}

		// Create or update specimen
		if err := s.db.Save(&specimen).Error; err != nil {
			log.Printf("Error saving specimen %s: %v", specimen.ID, err)
			continue
		}
	}

	return nil
}

// Helper functions to parse CSV records
func (s *CSVService) parsePatientRecord(record []string) models.Patient {
	patient := models.Patient{}
	
	if len(record) > 0 && record[0] != "" {
		patient.ID = record[0]
	}
	if len(record) > 1 && record[1] != "" {
		if t, err := time.Parse("2006-01-02T15:04:05.000Z", record[1]); err == nil {
			patient.SubmissionDate = t
		}
	}
	if len(record) > 2 && record[2] != "" {
		patient.Region = record[2]
	}
	if len(record) > 3 && record[3] != "" {
		patient.District = record[3]
	}
	if len(record) > 4 && record[4] != "" {
		patient.Subcounty = record[4]
	}
	if len(record) > 5 && record[5] != "" {
		patient.Facility = record[5]
	}
	if len(record) > 6 && record[6] != "" {
		patient.LevelOfCare = record[6]
	}
	if len(record) > 7 && record[7] != "" {
		patient.Ownership = record[7]
	}
	if len(record) > 8 && record[8] != "" {
		patient.WardName = record[8]
	}
	if len(record) > 9 && record[9] != "" {
		if val, err := strconv.Atoi(record[9]); err == nil {
			patient.WardTotalPatients = val
		}
	}
	if len(record) > 10 && record[10] != "" {
		if val, err := strconv.Atoi(record[10]); err == nil {
			patient.WardEligiblePatients = val
		}
	}
	if len(record) > 11 && record[11] != "" {
		if t, err := time.Parse("2006-01-02", record[11]); err == nil {
			patient.SurveyDate = t
		}
	}
	if len(record) > 12 && record[12] != "" {
		patient.PatientInitials = record[12]
	}
	if len(record) > 13 && record[13] != "" {
		patient.Code = record[13]
	}
	if len(record) > 15 && record[15] != "" {
		if val, err := strconv.Atoi(record[15]); err == nil {
			patient.RandNum = val
		}
	}
	if len(record) > 16 && record[16] != "" {
		patient.PatientCode = record[16]
	}
	if len(record) > 17 && record[17] != "" {
		patient.ShowCode = record[17]
	}
	if len(record) > 18 && record[18] != "" {
		patient.IsThePatientAnInfant = record[18]
	}
	if len(record) > 19 && record[19] != "" {
		if val, err := strconv.Atoi(record[19]); err == nil {
			patient.AgeMonths = val
		}
	}
	if len(record) > 20 && record[20] != "" {
		if val, err := strconv.Atoi(record[20]); err == nil {
			patient.AgeYears = val
		}
	}
	if len(record) > 21 && record[21] != "" {
		patient.PreTermBirth = record[21]
	}
	if len(record) > 22 && record[22] != "" {
		patient.Gender = record[22]
	}
	if len(record) > 23 && record[23] != "" {
		if val, err := strconv.ParseFloat(record[23], 64); err == nil {
			patient.Weight = val
		}
	}
	if len(record) > 24 && record[24] != "" {
		if val, err := strconv.ParseFloat(record[24], 64); err == nil {
			patient.WeightBirthKg = val
		}
	}
	if len(record) > 25 && record[25] != "" {
		if t, err := time.Parse("2006-01-02", record[25]); err == nil {
			patient.AdmissionDate = t
		}
	}
	if len(record) > 26 && record[26] != "" {
		patient.SurgerySinceAdmission = record[26]
	}
	if len(record) > 27 && record[27] != "" {
		patient.UrinaryCatheter = record[27]
	}
	if len(record) > 28 && record[28] != "" {
		patient.PeripheralVascularCatheter = record[28]
	}
	if len(record) > 29 && record[29] != "" {
		patient.CentralVascularCatheter = record[29]
	}
	if len(record) > 30 && record[30] != "" {
		patient.Intubation = record[30]
	}
	if len(record) > 31 && record[31] != "" {
		patient.PatientOnAntibiotic = record[31]
	}
	if len(record) > 32 && record[32] != "" {
		if val, err := strconv.Atoi(record[32]); err == nil {
			patient.PatientNumberAntibiotics = val
		}
	}
	if len(record) > 33 && record[33] != "" {
		patient.MalariaStatus = record[33]
	}
	if len(record) > 34 && record[34] != "" {
		patient.TuberculosisStatus = record[34]
	}
	if len(record) > 35 && record[35] != "" {
		patient.HIVStatus = record[35]
	}
	if len(record) > 36 && record[36] != "" {
		patient.HIVOnART = record[36]
	}
	if len(record) > 37 && record[37] != "" {
		patient.HIVCD4Count = record[37]
	}
	if len(record) > 38 && record[38] != "" {
		patient.HIVViralLoad = record[38]
	}
	if len(record) > 39 && record[39] != "" {
		patient.Diabetes = record[39]
	}
	if len(record) > 40 && record[40] != "" {
		patient.MalnutritionStatus = record[40]
	}
	if len(record) > 41 && record[41] != "" {
		patient.Hypertension = record[41]
	}
	if len(record) > 42 && record[42] != "" {
		patient.ReferredFrom = record[42]
	}
	if len(record) > 43 && record[43] != "" {
		patient.Hospitalization90Days = record[43]
	}
	if len(record) > 44 && record[44] != "" {
		patient.TypeSurgerySinceAdmission = record[44]
	}
	if len(record) > 45 && record[45] != "" {
		patient.AdditionalComment = record[45]
	}
	if len(record) > 46 && record[46] != "" {
		patient.Comments = record[46]
	}
	if len(record) > 47 && record[47] != "" {
		patient.InstanceID = record[47]
	}
	if len(record) > 48 && record[48] != "" {
		patient.SubmitterID = record[48]
	}
	if len(record) > 49 && record[49] != "" {
		patient.SubmitterName = record[49]
	}
	if len(record) > 50 && record[50] != "" {
		patient.AttachmentsPresent = record[50]
	}
	if len(record) > 51 && record[51] != "" {
		patient.AttachmentsExpected = record[51]
	}
	if len(record) > 52 && record[52] != "" {
		patient.Status = record[52]
	}
	if len(record) > 53 && record[53] != "" {
		patient.ReviewState = record[53]
	}
	if len(record) > 54 && record[54] != "" {
		patient.DeviceID = record[54]
	}
	if len(record) > 55 && record[55] != "" {
		patient.Edits = record[55]
	}
	if len(record) > 56 && record[56] != "" {
		patient.FormVersion = record[56]
	}

	return patient
}

func (s *CSVService) parseAntibioticRecord(record []string) models.Antibiotic {
	antibiotic := models.Antibiotic{}
	
	if len(record) > 0 && record[0] != "" {
		antibiotic.AntibioticNotes = record[0]
	}
	if len(record) > 1 && record[1] != "" {
		antibiotic.AntibioticINNName = record[1]
	}
	if len(record) > 2 && record[2] != "" {
		antibiotic.OtherAntibiotic = record[2]
	}
	if len(record) > 3 && record[3] != "" {
		antibiotic.ATCCode = record[3]
	}
	if len(record) > 4 && record[4] != "" {
		antibiotic.AntibioticClass = record[4]
	}
	if len(record) > 5 && record[5] != "" {
		antibiotic.AntibioticAwareClassification = record[5]
	}
	if len(record) > 6 && record[6] != "" {
		antibiotic.AntibioticWrittenInINN = record[6]
	}
	if len(record) > 7 && record[7] != "" {
		if t, err := time.Parse("2006-01-02", record[7]); err == nil {
			antibiotic.StartDateAntibiotic = t
		}
	}
	if len(record) > 8 && record[8] != "" {
		if val, err := strconv.ParseFloat(record[8], 64); err == nil {
			antibiotic.UnitDose = val
		}
	}
	if len(record) > 9 && record[9] != "" {
		antibiotic.UnitDosesCombination = record[9]
	}
	if len(record) > 10 && record[10] != "" {
		antibiotic.UnitDoseMeasureUnit = record[10]
	}
	if len(record) > 11 && record[11] != "" {
		antibiotic.UnitDoseFrequency = record[11]
	}
	if len(record) > 12 && record[12] != "" {
		antibiotic.AdministrationRoute = record[12]
	}
	if len(record) > 13 && record[13] != "" {
		antibiotic.ParentKey = record[13]
	}
	if len(record) > 14 && record[14] != "" {
		antibiotic.ID = record[14]
	}

	return antibiotic
}

func (s *CSVService) parseIndicationRecord(record []string) models.Indication {
	indication := models.Indication{}
	
	if len(record) > 0 && record[0] != "" {
		indication.IndicationType = record[0]
	}
	if len(record) > 1 && record[1] != "" {
		indication.SurgProphDuration = record[1]
	}
	if len(record) > 2 && record[2] != "" {
		indication.SurgProphSite = record[2]
	}
	if len(record) > 3 && record[3] != "" {
		indication.Diagnosis = record[3]
	}
	if len(record) > 4 && record[4] != "" {
		if t, err := time.Parse("2006-01-02", record[4]); err == nil {
			indication.StartDateTreatment = t
		}
	}
	if len(record) > 5 && record[5] != "" {
		indication.ReasonInNotes = record[5]
	}
	if len(record) > 6 && record[6] != "" {
		indication.CultureSampleTaken = record[6]
	}
	if len(record) > 7 && record[7] != "" {
		indication.ParentKey = record[7]
	}
	if len(record) > 8 && record[8] != "" {
		indication.ID = record[8]
	}

	return indication
}

func (s *CSVService) parseOptionalVarRecord(record []string) models.OptionalVar {
	optionalVar := models.OptionalVar{}
	
	if len(record) > 0 && record[0] != "" {
		optionalVar.PrescriberType = record[0]
	}
	if len(record) > 1 && record[1] != "" {
		optionalVar.IntravenousType = record[1]
	}
	if len(record) > 2 && record[2] != "" {
		optionalVar.OralSwitch = record[2]
	}
	if len(record) > 3 && record[3] != "" {
		if val, err := strconv.Atoi(record[3]); err == nil {
			optionalVar.NumberMissedDoses = val
		}
	}
	if len(record) > 4 && record[4] != "" {
		optionalVar.MissedDosesReason = record[4]
	}
	if len(record) > 5 && record[5] != "" {
		optionalVar.GuidelinesCompliance = record[5]
	}
	if len(record) > 6 && record[6] != "" {
		optionalVar.TreatmentType = record[6]
	}
	if len(record) > 7 && record[7] != "" {
		optionalVar.ParentKey = record[7]
	}
	if len(record) > 8 && record[8] != "" {
		optionalVar.ID = record[8]
	}

	return optionalVar
}

func (s *CSVService) parseSpecimenRecord(record []string) models.Specimen {
	specimen := models.Specimen{}
	
	if len(record) > 0 && record[0] != "" {
		specimen.SpecimenType = record[0]
	}
	if len(record) > 1 && record[1] != "" {
		specimen.CultureResult = record[1]
	}
	if len(record) > 2 && record[2] != "" {
		specimen.Microorganism = record[2]
	}
	if len(record) > 3 && record[3] != "" {
		specimen.AntibioticSusceptibilityTestResults = record[3]
	}
	if len(record) > 4 && record[4] != "" {
		specimen.ResistantPhenotype = record[4]
	}
	if len(record) > 5 && record[5] != "" {
		specimen.ParentKey = record[5]
	}
	if len(record) > 6 && record[6] != "" {
		specimen.ID = record[6]
	}

	return specimen
}
