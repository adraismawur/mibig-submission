package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/adraismawur/mibig-submission/models/entry/gene"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"strconv"
)

func init() {
	RegisterEndpointGenerator(GeneInformationEndpoint)
}

func GeneInformationEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: http.MethodGet,
				Path:   "/entry/:accession/gene_information",
				Handler: func(c *gin.Context) {
					getEntryGeneInformation(db, c)
				},
			},
			{
				Method: http.MethodGet,
				Path:   "/entry/:accession/gene_information/to_add/:addition_id",
				Handler: func(c *gin.Context) {
					getEntryGeneAddition(db, c)
				},
			},
			{
				Method: http.MethodGet,
				Path:   "/entry/:accession/gene_information/to_delete/:deletion_id",
				Handler: func(c *gin.Context) {
					getEntryGeneDeletion(db, c)
				},
			},
			{
				Method: http.MethodGet,
				Path:   "/entry/:accession/gene_information/annotation/:annotation_id",
				Handler: func(c *gin.Context) {
					getEntryGeneAnnotation(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/gene_information/to_add",
				Handler: func(c *gin.Context) {
					updateOrCreateEntryGeneAddition(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/gene_information/to_delete",
				Handler: func(c *gin.Context) {
					updateOrCreateEntryGeneDeletion(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/gene_information/annotation",
				Handler: func(c *gin.Context) {
					updateOrCreateEntryGeneAnnotation(db, c)
				},
			},
			{
				Method: http.MethodDelete,
				Path:   "/entry/:accession/gene_information/to_add/:addition_id",
				Handler: func(c *gin.Context) {
					deleteEntryGeneAddition(db, c)
				},
			},
			{
				Method: http.MethodDelete,
				Path:   "/entry/:accession/gene_information/to_delete/:deletion_id",
				Handler: func(c *gin.Context) {
					deleteEntryGeneDeletion(db, c)
				},
			},
			{
				Method: http.MethodDelete,
				Path:   "/entry/:accession/gene_information/annotation/:annotation_id",
				Handler: func(c *gin.Context) {
					deleteEntryGeneAnnotation(db, c)
				},
			},
		},
	}
}

func getEntryGeneInformation(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	exists, err := entry.GetEntryExists(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "entry not found"})
		return
	}

	entryGene, err := gene.GetEntryGeneInformation(db, accession)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entryGene)
}

func getEntryGeneAddition(db *gorm.DB, c *gin.Context) {
	additionId, err := strconv.Atoi(c.Param("addition_id"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Could not convert parameter to integer"})
		return
	}

	formatJson := c.Query("pretty") == "true"

	additionExists, err := gene.GetGeneAdditionExists(db, additionId)

	if err != nil {
		slog.Error("[endpoints] [genes] Error finding gene addition", "addition_id", additionId, "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !additionExists {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	addition, err := gene.GetGeneAddition(db, additionId)

	if err != nil {
		slog.Error("[endpoints] [genes] Could not retrieve gene addition", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !formatJson {
		c.JSON(http.StatusOK, addition)
		return
	}

	formattedJson, err := json.MarshalIndent(addition, "", "  ")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		slog.Error("[endpoints] [entry] Failed to marshal gene addition", "error", err.Error())
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(formattedJson))
}

func getEntryGeneDeletion(db *gorm.DB, c *gin.Context) {
	deletionId, err := strconv.Atoi(c.Param("deletion_id"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Could not convert parameter to integer"})
		return
	}

	formatJson := c.Query("pretty") == "true"

	deletionExists, err := gene.GetGeneDeletionExists(db, deletionId)

	if err != nil {
		slog.Error("[endpoints] [genes] Error finding gene deletion", "deletion_id", deletionId, "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !deletionExists {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	deletion, err := gene.GetGeneDeletion(db, deletionId)

	if err != nil {
		slog.Error("[endpoints] [genes] Could not retrieve gene deletion", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !formatJson {
		c.JSON(http.StatusOK, deletion)
		return
	}

	formattedJson, err := json.MarshalIndent(deletion, "", "  ")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		slog.Error("[endpoints] [entry] Failed to marshal gene deletion", "error", err.Error())
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(formattedJson))
}

func getEntryGeneAnnotation(db *gorm.DB, c *gin.Context) {
	annotationId, err := strconv.Atoi(c.Param("annotation_id"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Could not convert parameter to integer"})
		return
	}

	formatJson := c.Query("pretty") == "true"

	annotationExists, err := gene.GetGeneAnnotationExists(db, annotationId)

	if err != nil {
		slog.Error("[endpoints] [genes] Error finding gene annotation", "annotation_id", annotationId, "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !annotationExists {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	annotation, err := gene.GetGeneAnnotation(db, annotationId)

	if err != nil {
		slog.Error("[endpoints] [genes] Could not retrieve gene annotation", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !formatJson {
		c.JSON(http.StatusOK, annotation)
		return
	}

	formattedJson, err := json.MarshalIndent(annotation, "", "  ")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		slog.Error("[endpoints] [entry] Failed to marshal gene annotation", "error", err.Error())
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(formattedJson))
}

func updateOrCreateEntryGeneAddition(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	var addition gene.GeneAddition

	err := c.ShouldBindJSON(&addition)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newAddition, err := gene.UpdateOrCreateGeneAddition(db, accession, &addition)

	if err != nil {
		slog.Error("[endpoints] [genes] Could not update or create gene addition")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.JSON(http.StatusOK, newAddition)
}

func updateOrCreateEntryGeneDeletion(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	var deletion gene.GeneDeletion

	err := c.ShouldBindJSON(&deletion)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newDeletion, err := gene.UpdateOrCreateGeneDeletion(db, accession, &deletion)

	if err != nil {
		slog.Error("[endpoints] [genes] Could not update or create gene deletion")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.JSON(http.StatusOK, newDeletion)
}

func updateOrCreateEntryGeneAnnotation(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")

	var annotation gene.GeneAnnotation

	err := c.ShouldBindJSON(&annotation)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newDeletion, err := gene.UpdateOrCreateGeneAnnotation(db, accession, &annotation)

	if err != nil {
		slog.Error("[endpoints] [genes] Could not update or create gene deletion")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.JSON(http.StatusOK, newDeletion)
}

func deleteEntryGeneAddition(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	additionId, err := strconv.Atoi(c.Param("addition_id"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = gene.DeleteGeneAddition(db, additionId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.Status(http.StatusOK)
}

func deleteEntryGeneDeletion(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	deletionId, err := strconv.Atoi(c.Param("deletion_id"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = gene.DeleteGeneDeletion(db, deletionId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.Status(http.StatusOK)
}

func deleteEntryGeneAnnotation(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	annotationId, err := strconv.Atoi(c.Param("annotation_id"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = gene.DeleteGeneAnnotation(db, annotationId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.Status(http.StatusOK)
}
