package endpoints

import (
	"github.com/adraismawur/mibig-submission/models"
	"github.com/adraismawur/mibig-submission/models/entry"
	"github.com/adraismawur/mibig-submission/models/entry/biosynthesis"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func init() {
	RegisterEndpointGenerator(BiosynthDomainEndpoint)
}

func BiosynthDomainEndpoint(db *gorm.DB) Endpoint {
	return Endpoint{
		Routes: []Route{
			{
				Method: http.MethodGet,
				Path:   "/entry/:accession/biosynth/modification_domain/list/:id",
				Handler: func(c *gin.Context) {
					getModuleModificationDomains(db, c)
				},
			},
			{
				Method: http.MethodGet,
				Path:   "/entry/:accession/biosynth/modification_domain/:modification_domain_id",
				Handler: func(c *gin.Context) {
					getModuleModificationDomain(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/biosynth/modification_domain/add/:id",
				Handler: func(c *gin.Context) {
					createModuleModificationDomain(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/biosynth/modification_domain/:modification_domain_id",
				Handler: func(c *gin.Context) {
					updateModuleModificationDomain(db, c)
				},
			},
			{
				Method: http.MethodDelete,
				Path:   "/entry/:accession/biosynth/modification_domain/:modification_domain_id",
				Handler: func(c *gin.Context) {
					deleteModuleModificationDomain(db, c)
				},
			},
		},
	}
}

func getModuleModificationDomains(db *gorm.DB, c *gin.Context) {
	moduleId := c.Param("id")
	iModuleId, err := strconv.Atoi(moduleId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not convert ID parameter to int"})
		return
	}

	modificationDomains, err := biosynthesis.GetModificationDomains(db, iModuleId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, modificationDomains)
}

func getModuleModificationDomain(db *gorm.DB, c *gin.Context) {
	modificationDomainId := c.Param("modification_domain_id")
	iModificationDomainId, err := strconv.Atoi(modificationDomainId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not convert ID parameter to int"})
		return
	}

	modificationDomain, err := biosynthesis.GetModificationDomain(db, iModificationDomainId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, modificationDomain)
}

func createModuleModificationDomain(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	moduleId := c.Param("id")
	iModuleId, err := strconv.Atoi(moduleId)

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newModificationDomain biosynthesis.ModificationDomain
	err = c.ShouldBindJSON(&newModificationDomain)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = biosynthesis.CreateModificationDomain(db, iModuleId, newModificationDomain)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.Status(http.StatusOK)
}

func updateModuleModificationDomain(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	modificationDomainId := c.Param("modification_domain_id")
	_, err := strconv.Atoi(modificationDomainId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not convert ID parameter to int"})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newModificationDomain biosynthesis.ModificationDomain
	err = c.ShouldBindJSON(&newModificationDomain)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = biosynthesis.UpdateModificationDomain(db, newModificationDomain)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.JSON(http.StatusOK, newModificationDomain)
}

func deleteModuleModificationDomain(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	modificationDomainId := c.Param("modification_domain_id")
	iModificationDomainId, err := strconv.Atoi(modificationDomainId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not convert ID parameter to int"})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = biosynthesis.DeleteModificationDomain(db, iModificationDomainId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.AddContributor(db, accession, user.ID)

	c.Status(http.StatusOK)
}
