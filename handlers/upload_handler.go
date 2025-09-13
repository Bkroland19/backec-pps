package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"point-prevalence-survey/services"
	"strings"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	csvService *services.CSVService
}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{
		csvService: services.NewCSVService(),
	}
}

// validateCSVFile validates the uploaded file
func (h *UploadHandler) validateCSVFile(fileHeader *multipart.FileHeader) error {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".csv" {
		return fmt.Errorf("invalid file type. Only CSV files are allowed")
	}

	// Check file size (max 10MB)
	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if fileHeader.Size > maxFileSize {
		return fmt.Errorf("file too large. Maximum size allowed is 10MB")
	}

	return nil
}

// processUpload handles common upload logic
func (h *UploadHandler) processUpload(c *gin.Context, uploadFunc func(multipart.File) (*services.UploadResult, error)) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No file uploaded",
			"message": "Please select a CSV file to upload",
		})
		return
	}
	defer file.Close()

	// Validate file
	if err := h.validateCSVFile(fileHeader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid file",
			"message": err.Error(),
		})
		return
	}

	// Process the file
	result, err := uploadFunc(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Upload failed",
			"message": err.Error(),
		})
		return
	}

	// Return success response with statistics
	c.JSON(http.StatusOK, gin.H{
		"message":           "File uploaded and processed successfully",
		"filename":          fileHeader.Filename,
		"total_records":     result.TotalRecords,
		"processed_records": result.ProcessedRecords,
		"skipped_records":   result.SkippedRecords,
		"inserted_records":  result.InsertedRecords,
		"updated_records":   result.UpdatedRecords,
		"errors":            result.Errors,
	})
}

// UploadPatients godoc
// @Summary Upload patients CSV file
// @Description Upload and import patients data from CSV file
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file containing patients data"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/upload/patients [post]
func (h *UploadHandler) UploadPatients(c *gin.Context) {
	h.processUpload(c, h.csvService.ImportPatients)
}

// UploadAntibiotics godoc
// @Summary Upload antibiotics CSV file
// @Description Upload and import antibiotics data from CSV file
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file containing antibiotics data"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/upload/antibiotics [post]
func (h *UploadHandler) UploadAntibiotics(c *gin.Context) {
	h.processUpload(c, h.csvService.ImportAntibiotics)
}

// UploadIndications godoc
// @Summary Upload indications CSV file
// @Description Upload and import indications data from CSV file
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file containing indications data"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/upload/indications [post]
func (h *UploadHandler) UploadIndications(c *gin.Context) {
	h.processUpload(c, h.csvService.ImportIndications)
}

// UploadOptionalVars godoc
// @Summary Upload optional variables CSV file
// @Description Upload and import optional variables data from CSV file
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file containing optional variables data"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/upload/optional-vars [post]
func (h *UploadHandler) UploadOptionalVars(c *gin.Context) {
	h.processUpload(c, h.csvService.ImportOptionalVars)
}

// UploadSpecimens godoc
// @Summary Upload specimens CSV file
// @Description Upload and import specimens data from CSV file
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file containing specimens data"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/upload/specimens [post]
func (h *UploadHandler) UploadSpecimens(c *gin.Context) {
	h.processUpload(c, h.csvService.ImportSpecimens)
}

// UploadAntibioticDetails handles CSV file upload for antibiotic details
// @Summary Upload antibiotic details CSV file
// @Description Upload a CSV file containing antibiotic details data
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file containing antibiotic details data"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/upload/antibiotic-details [post]
func (h *UploadHandler) UploadAntibioticDetails(c *gin.Context) {
	h.processUpload(c, h.csvService.ImportAntibioticDetails)
}
