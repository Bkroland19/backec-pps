package models

import (
	"time"
)

// Patient represents the main patient record
type Patient struct {
	ID                         string    `json:"id" gorm:"primaryKey;column:key"`
	SubmissionDate             time.Time `json:"submission_date" gorm:"column:submission_date"`
	Region                     string    `json:"region"`
	District                   string    `json:"district"`
	Subcounty                  string    `json:"subcounty"`
	Facility                   string    `json:"facility"`
	LevelOfCare                string    `json:"level_of_care" gorm:"column:level_of_care"`
	Ownership                  string    `json:"ownership"`
	WardName                   string    `json:"ward_name" gorm:"column:ward_name"`
	WardTotalPatients          int       `json:"ward_total_patients" gorm:"column:ward_total_patients"`
	WardEligiblePatients       int       `json:"ward_eligible_patients" gorm:"column:ward_eligible_patients"`
	SurveyDate                 time.Time `json:"survey_date" gorm:"column:survey_date"`
	PatientInitials            string    `json:"patient_initials" gorm:"column:patient_initials"`
	Code                       string    `json:"code"`
	RandNum                    int       `json:"rand_num" gorm:"column:rand_num"`
	PatientCode                string    `json:"patient_code" gorm:"column:patient_code"`
	ShowCode                   string    `json:"show_code" gorm:"column:show_code"`
	IsThePatientAnInfant       string    `json:"is_the_patient_an_infant" gorm:"column:is_the_patient_an_infant"`
	AgeMonths                  int       `json:"age_months" gorm:"column:age_months"`
	AgeYears                   int       `json:"age_years" gorm:"column:age_years"`
	PreTermBirth               string    `json:"pre_term_birth" gorm:"column:pre_term_birth"`
	Gender                     string    `json:"gender"`
	Weight                     float64   `json:"weight"`
	WeightBirthKg              float64   `json:"weight_birth_kg" gorm:"column:weight_birth_kg"`
	AdmissionDate              time.Time `json:"admission_date" gorm:"column:admission_date"`
	SurgerySinceAdmission      string    `json:"surgery_since_admission" gorm:"column:surgery_since_admission"`
	UrinaryCatheter            string    `json:"urinary_catheter" gorm:"column:urinary_catheter"`
	PeripheralVascularCatheter string    `json:"peripheral_vascular_catheter" gorm:"column:peripheral_vascular_catheter"`
	CentralVascularCatheter    string    `json:"central_vascular_catheter" gorm:"column:central_vascular_catheter"`
	Intubation                 string    `json:"intubation"`
	PatientOnAntibiotic        string    `json:"patient_on_antibiotic" gorm:"column:patient_on_antibiotic"`
	PatientNumberAntibiotics   int       `json:"patient_number_antibiotics" gorm:"column:patient_number_antibiotics"`
	MalariaStatus              string    `json:"malaria_status" gorm:"column:malaria_status"`
	TuberculosisStatus         string    `json:"tuberculosis_status" gorm:"column:tuberculosis_status"`
	HIVStatus                  string    `json:"hiv_status" gorm:"column:hiv_status"`
	HIVOnART                   string    `json:"hiv_on_art" gorm:"column:hiv_on_art"`
	HIVCD4Count                string    `json:"hiv_cd4_count" gorm:"column:hiv_cd4_count"`
	HIVViralLoad               string    `json:"hiv_viral_load" gorm:"column:hiv_viral_load"`
	Diabetes                   string    `json:"diabetes"`
	MalnutritionStatus         string    `json:"malnutrition_status" gorm:"column:malnutrition_status"`
	Hypertension               string    `json:"hypertension"`
	ReferredFrom               string    `json:"referred_from" gorm:"column:referred_from"`
	Hospitalization90Days      string    `json:"hospitalization_90_days" gorm:"column:hospitalization_90_days"`
	TypeSurgerySinceAdmission  string    `json:"type_surgery_since_admission" gorm:"column:type_surgery_since_admission"`
	AdditionalComment          string    `json:"additional_comment" gorm:"column:additional_comment"`
	Comments                   string    `json:"comments"`
	InstanceID                 string    `json:"instance_id" gorm:"column:instance_id"`
	SubmitterID                string    `json:"submitter_id" gorm:"column:submitter_id"`
	SubmitterName              string    `json:"submitter_name" gorm:"column:submitter_name"`
	AttachmentsPresent         string    `json:"attachments_present" gorm:"column:attachments_present"`
	AttachmentsExpected        string    `json:"attachments_expected" gorm:"column:attachments_expected"`
	Status                     string    `json:"status"`
	ReviewState                string    `json:"review_state" gorm:"column:review_state"`
	DeviceID                   string    `json:"device_id" gorm:"column:device_id"`
	Edits                      string    `json:"edits"`
	FormVersion                string    `json:"form_version" gorm:"column:form_version"`

	// Relationships
	Antibiotics  []Antibiotic  `json:"antibiotics" gorm:"foreignKey:ParentKey;references:ID"`
	Indications  []Indication  `json:"indications" gorm:"foreignKey:ParentKey;references:ID"`
	OptionalVars []OptionalVar `json:"optional_vars" gorm:"foreignKey:ParentKey;references:ID"`
	Specimens    []Specimen    `json:"specimens" gorm:"foreignKey:ParentKey;references:ID"`
}

// Antibiotic represents antibiotic data
type Antibiotic struct {
	ID                            string    `json:"id" gorm:"primaryKey;column:key"`
	AntibioticNotes               string    `json:"antibiotic_notes" gorm:"column:antibiotic_notes"`
	AntibioticINNName             string    `json:"antibiotic_inn_name" gorm:"column:antibiotic_inn_name"`
	OtherAntibiotic               string    `json:"other_antibiotic" gorm:"column:other_antibiotic"`
	ATCCode                       string    `json:"atc_code" gorm:"column:atc_code"`
	AntibioticClass               string    `json:"antibiotic_class" gorm:"column:antibiotic_class"`
	AntibioticAwareClassification string    `json:"antibiotic_aware_classification" gorm:"column:antibiotic_aware_classification"`
	AntibioticWrittenInINN        string    `json:"antibiotic_written_in_inn" gorm:"column:antibiotic_written_in_inn"`
	StartDateAntibiotic           time.Time `json:"start_date_antibiotic" gorm:"column:start_date_antibiotic"`
	UnitDose                      float64   `json:"unit_dose" gorm:"column:unit_dose"`
	UnitDosesCombination          string    `json:"unit_doses_combination" gorm:"column:unit_doses_combination"`
	UnitDoseMeasureUnit           string    `json:"unit_dose_measure_unit" gorm:"column:unit_dose_measure_unit"`
	UnitDoseFrequency             string    `json:"unit_dose_frequency" gorm:"column:unit_dose_frequency"`
	AdministrationRoute           string    `json:"administration_route" gorm:"column:administration_route"`
	ParentKey                     string    `json:"parent_key" gorm:"column:parent_key"`
}

// AntibioticDetails represents additional antibiotic details
type AntibioticDetails struct {
	ID           string `json:"id" gorm:"primaryKey;column:key"`
	Prescriber   string `json:"prescriber" gorm:"column:prescriber"`
	Intraveno    string `json:"intraveno" gorm:"column:intraveno"`
	OralSwitch   string `json:"oral_switch" gorm:"column:oral_switch"`
	NumberMissed string `json:"number_missed" gorm:"column:number_missed"`
	MissedDose   string `json:"missed_dose" gorm:"column:missed_dose"`
	Guideline    string `json:"guideline" gorm:"column:guideline"`
	Treatment    string `json:"treatment" gorm:"column:treatment"`
	ParentKey    string `json:"parent_key" gorm:"column:parent_key"`
}

// Indication represents indication data
type Indication struct {
	ID                 string    `json:"id" gorm:"primaryKey;column:key"`
	IndicationType     string    `json:"indication_type" gorm:"column:indication_type"`
	SurgProphDuration  string    `json:"surg_proph_duration" gorm:"column:surg_proph_duration"`
	SurgProphSite      string    `json:"surg_proph_site" gorm:"column:surg_proph_site"`
	Diagnosis          string    `json:"diagnosis"`
	StartDateTreatment time.Time `json:"start_date_treatment" gorm:"column:start_date_treatment"`
	ReasonInNotes      string    `json:"reason_in_notes" gorm:"column:reason_in_notes"`
	CultureSampleTaken string    `json:"culture_sample_taken" gorm:"column:culture_sample_taken"`
	ParentKey          string    `json:"parent_key" gorm:"column:parent_key"`
}

// OptionalVar represents optional variables data
type OptionalVar struct {
	ID                   string `json:"id" gorm:"primaryKey;column:key"`
	PrescriberType       string `json:"prescriber_type" gorm:"column:prescriber_type"`
	IntravenousType      string `json:"intravenous_type" gorm:"column:intravenous_type"`
	OralSwitch           string `json:"oral_switch" gorm:"column:oral_switch"`
	NumberMissedDoses    int    `json:"number_missed_doses" gorm:"column:number_missed_doses"`
	MissedDosesReason    string `json:"missed_doses_reason" gorm:"column:missed_doses_reason"`
	GuidelinesCompliance string `json:"guidelines_compliance" gorm:"column:guidelines_compliance"`
	TreatmentType        string `json:"treatment_type" gorm:"column:treatment_type"`
	ParentKey            string `json:"parent_key" gorm:"column:parent_key"`
}

// Specimen represents specimen data
type Specimen struct {
	ID                                  string `json:"id" gorm:"primaryKey;column:key"`
	SpecimenType                        string `json:"specimen_type" gorm:"column:specimen_type"`
	CultureResult                       string `json:"culture_result" gorm:"column:culture_result"`
	Microorganism                       string `json:"microorganism"`
	AntibioticSusceptibilityTestResults string `json:"antibiotic_susceptibility_test_results" gorm:"column:antibiotic_susceptibility_test_results"`
	ResistantPhenotype                  string `json:"resistant_phenotype" gorm:"column:resistant_phenotype"`
	ParentKey                           string `json:"parent_key" gorm:"column:parent_key"`
}

// TableName methods to specify table names
func (Patient) TableName() string {
	return "patients"
}

func (Antibiotic) TableName() string {
	return "antibiotics"
}

func (Indication) TableName() string {
	return "indications"
}

func (OptionalVar) TableName() string {
	return "optional_vars"
}

func (AntibioticDetails) TableName() string {
	return "antibiotic_details"
}

func (Specimen) TableName() string {
	return "specimens"
}
