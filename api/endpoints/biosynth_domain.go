package endpoints

import (
	"errors"
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
				Path:   "/entry/:accession/biosynth/:domain_type/list/:module_id",
				Handler: func(c *gin.Context) {
					getModuleDomains(db, c)
				},
			},
			{
				Method: http.MethodGet,
				Path:   "/entry/:accession/biosynth/:domain_type/:module_id/:domain_id",
				Handler: func(c *gin.Context) {
					getModuleDomain(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/biosynth/:domain_type/create/:module_id",
				Handler: func(c *gin.Context) {
					createModuleDomain(db, c)
				},
			},
			{
				Method: http.MethodPost,
				Path:   "/entry/:accession/biosynth/:domain_type/:module_id/:domain_id",
				Handler: func(c *gin.Context) {
					updateModuleDomain(db, c)
				},
			},
			{
				Method: http.MethodDelete,
				Path:   "/entry/:accession/biosynth/:domain_type/:module_id/:domain_id",
				Handler: func(c *gin.Context) {
					deleteModuleDomain(db, c)
				},
			},
		},
	}
}

func getModuleDomains(db *gorm.DB, c *gin.Context) {
	moduleId, err := strconv.Atoi(c.Param("module_id"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not convert ID parameter to int"})
		return
	}

	modificationDomains, err := biosynthesis.GetModificationDomains(db, moduleId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, modificationDomains)
}

func getModuleDomain(db *gorm.DB, c *gin.Context) {
	//accession := c.Param("accession")
	domainType := c.Param("domain_type")
	moduleId, err := strconv.Atoi(c.Param("module_id"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not convert ID parameter to int"})
		return
	}

	domainId, err := strconv.Atoi(c.Param("domain_id"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not convert ID parameter to int"})
		return
	}

	if domainType == "modification_domain" {
		modificationDomain, err := biosynthesis.GetModificationDomain(db, domainId)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, modificationDomain)
		return
	}

	module, err := biosynthesis.GetEntryBiosynthesisModule(db, moduleId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch domainType {
	case "adenylation":
		fallthrough
	case "a_domain":
		aDomain := module.ADomain
		c.JSON(http.StatusOK, aDomain)
		return
	case "condensation":
		fallthrough
	case "c_domain":
		cDomain := module.CDomain
		c.JSON(http.StatusOK, cDomain)
		return
	case "ks_domain":
		fallthrough
	case "ketosynthase":
		ksDomain := module.KSDomain
		c.JSON(http.StatusOK, ksDomain)
		return
	case "at_domain":
		fallthrough
	case "acyltransferase":
		atDomain := module.ATDomain
		c.JSON(http.StatusOK, atDomain)
		return
	}

	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "domain type not supported"})
}

func createModuleDomain(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	domainType := c.Param("domain_type")
	moduleId, err := strconv.Atoi(c.Param("module_id"))

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if domainType == "modification_domain" {
		var newModificationDomain biosynthesis.ModificationDomain
		err = c.ShouldBindJSON(&newModificationDomain)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = biosynthesis.CreateModificationDomain(db, moduleId, newModificationDomain)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	switch domainType {

	}

	entry.AddContributor(db, accession, user.ID)

	c.Status(http.StatusOK)
}

func updateModuleDomain(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	domainType := c.Param("domain_type")
	moduleId, err := strconv.Atoi(c.Param("module_id"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not convert ID parameter to int"})
		return
	}

	_, err = strconv.Atoi(c.Param("domain_id"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not convert ID parameter to int"})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if domainType == "modification_domain" {
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
		return
	} else {
		module, err := biosynthesis.GetEntryBiosynthesisModule(db, moduleId)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		switch domainType {
		case "adenylation":
			fallthrough
		case "a_domain":
			var newADomain biosynthesis.AdenylationDomain
			err := c.ShouldBindJSON(&newADomain)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			err = db.Model(&module.ADomain).Save(newADomain).Error

			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			err = db.Model(&module.ADomain).Association("evidence").Replace(newADomain.Evidence)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			err = db.Model(&module.ADomain).Association("substrates").Replace(newADomain.Substrates)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, newADomain)
		default:
			err = errors.New("domain type not supported")
		}

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	entry.AddContributor(db, accession, user.ID)

}

func deleteModuleDomain(db *gorm.DB, c *gin.Context) {
	accession := c.Param("accession")
	domainType := c.Param("domain_type")
	moduleId, err := strconv.Atoi(c.Param("module_id"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not convert ID parameter to int"})
		return
	}

	domainId, err := strconv.Atoi(c.Param("domain_id"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not convert ID parameter to int"})
		return
	}

	user, err := models.GetUserFromContext(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if domainType == "modification_domain" {
		err = biosynthesis.DeleteModificationDomain(db, domainId)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		module, err := biosynthesis.GetEntryBiosynthesisModule(db, moduleId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		switch domainType {
		case "adenylation":
			err = db.Delete(&module.ADomain).Error
		default:
			err = errors.New("domain type not supported")
		}

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	entry.AddContributor(db, accession, user.ID)

	c.Status(http.StatusOK)
}
