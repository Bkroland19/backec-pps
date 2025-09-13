package routes

import (
	"point-prevalence-survey/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine) {
	// Initialize handlers
	patientHandler := handlers.NewPatientHandler()
	uploadHandler := handlers.NewUploadHandler()
	antibioticHandler := handlers.NewAntibioticHandler()
	antibioticDetailsHandler := handlers.NewAntibioticDetailsHandler()
	specimenHandler := handlers.NewSpecimenHandler()
	ppsCalculationsHandler := handlers.NewPPSCalculationsHandler()

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Patient routes
		patients := v1.Group("/patients")
		{
			patients.GET("", patientHandler.GetPatients)
			patients.GET("/stats", patientHandler.GetPatientStats)
			patients.GET("/:id", patientHandler.GetPatient)
			patients.GET("/:id/antibiotics", patientHandler.GetPatientAntibiotics)
			patients.GET("/:id/indications", patientHandler.GetPatientIndications)
			patients.GET("/:id/optional-vars", patientHandler.GetPatientOptionalVars)
			patients.GET("/:id/specimens", patientHandler.GetPatientSpecimens)
		}

		// Antibiotic routes
		antibiotics := v1.Group("/antibiotics")
		{
			antibiotics.GET("", antibioticHandler.GetAntibiotics)
			antibiotics.GET("/stats", antibioticHandler.GetAntibioticStats)
			antibiotics.GET("/:id", antibioticHandler.GetAntibiotic)
			antibiotics.GET("/patient/:patient_id", antibioticHandler.GetAntibioticUsageByPatient)
		}

		// Antibiotic Details routes
		antibioticDetails := v1.Group("/antibiotic-details")
		{
			antibioticDetails.GET("", antibioticDetailsHandler.GetAntibioticDetails)
			antibioticDetails.GET("/stats", antibioticDetailsHandler.GetAntibioticDetailsStats)
			antibioticDetails.GET("/:id", antibioticDetailsHandler.GetAntibioticDetailsByID)
			antibioticDetails.GET("/parent/:parent_key", antibioticDetailsHandler.GetAntibioticDetailsByParentKey)
			antibioticDetails.POST("", antibioticDetailsHandler.CreateAntibioticDetails)
			antibioticDetails.PUT("/:id", antibioticDetailsHandler.UpdateAntibioticDetails)
			antibioticDetails.DELETE("/:id", antibioticDetailsHandler.DeleteAntibioticDetails)
		}

		// Specimen routes
		specimens := v1.Group("/specimens")
		{
			specimens.GET("", specimenHandler.GetSpecimens)
			specimens.GET("/stats", specimenHandler.GetSpecimenStats)
			specimens.GET("/:id", specimenHandler.GetSpecimen)
			specimens.GET("/patient/:patient_id", specimenHandler.GetSpecimensByPatient)
		}

		// PPS Calculations routes
		pps := v1.Group("/pps")
		{
			pps.GET("/indicators", ppsCalculationsHandler.GetAllIndicators)
			pps.GET("/basic-metrics", ppsCalculationsHandler.GetBasicMetrics)
			pps.GET("/injectable-metrics", ppsCalculationsHandler.GetInjectableMetrics)
			pps.GET("/generic-metrics", ppsCalculationsHandler.GetGenericMetrics)
			pps.GET("/guideline-metrics", ppsCalculationsHandler.GetGuidelineMetrics)
			pps.GET("/diagnosis-metrics", ppsCalculationsHandler.GetDiagnosisMetrics)
			pps.GET("/culture-metrics", ppsCalculationsHandler.GetCultureMetrics)
			pps.GET("/missed-dose-metrics", ppsCalculationsHandler.GetMissedDoseMetrics)
			pps.GET("/prescriber-metrics", ppsCalculationsHandler.GetPrescriberMetrics)
			pps.GET("/oral-switch-metrics", ppsCalculationsHandler.GetOralSwitchMetrics)
		}

		// Upload routes
		upload := v1.Group("/upload")
		{
			upload.POST("/patients", uploadHandler.UploadPatients)
			upload.POST("/antibiotics", uploadHandler.UploadAntibiotics)
			upload.POST("/antibiotic-details", uploadHandler.UploadAntibioticDetails)
			upload.POST("/indications", uploadHandler.UploadIndications)
			upload.POST("/optional-vars", uploadHandler.UploadOptionalVars)
			upload.POST("/specimens", uploadHandler.UploadSpecimens)
		}
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Point Prevalence Survey API is running",
		})
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
