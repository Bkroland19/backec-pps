package handlers

import (
	"net/http"
	"point-prevalence-survey/services"

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

// UploadPatients godoc
// @Summary Upload patients CSV file
// @Description Upload and import patients data from CSV file
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file containing patients data"
// @Success 200 {object} map[string]string
// @Router /api/v1/upload/patients [post]
func (h *UploadHandler) UploadPatients(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	if err := h.csvService.ImportPatients(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Patients data imported successfully"})
}

// UploadAntibiotics godoc
// @Summary Upload antibiotics CSV file
// @Description Upload and import antibiotics data from CSV file
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file containing antibiotics data"
// @Success 200 {object} map[string]string
// @Router /api/v1/upload/antibiotics [post]
func (h *UploadHandler) UploadAntibiotics(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	if err := h.csvService.ImportAntibiotics(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Antibiotics data imported successfully"})
}

// UploadIndications godoc
// @Summary Upload indications CSV file
// @Description Upload and import indications data from CSV file
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file containing indications data"
// @Success 200 {object} map[string]string
// @Router /api/v1/upload/indications [post]
func (h *UploadHandler) UploadIndications(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	if err := h.csvService.ImportIndications(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Indications data imported successfully"})
}

// UploadOptionalVars godoc
// @Summary Upload optional variables CSV file
// @Description Upload and import optional variables data from CSV file
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file containing optional variables data"
// @Success 200 {object} map[string]string
// @Router /api/v1/upload/optional-vars [post]
func (h *UploadHandler) UploadOptionalVars(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	if err := h.csvService.ImportOptionalVars(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Optional variables data imported successfully"})
}

// UploadSpecimens godoc
// @Summary Upload specimens CSV file
// @Description Upload and import specimens data from CSV file
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file containing specimens data"
// @Success 200 {object} map[string]string
// @Router /api/v1/upload/specimens [post]
func (h *UploadHandler) UploadSpecimens(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	if err := h.csvService.ImportSpecimens(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Specimens data imported successfully"})
}
