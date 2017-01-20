package controllers

import (
	"fmt"
	"net/http"

	h "github.com/baelorswift/api/helpers"
	"github.com/baelorswift/api/middleware"
	m "github.com/baelorswift/api/models"

	"gopkg.in/gin-gonic/gin.v1"
)

// StudiosController ..
type StudiosController struct {
	context *m.Context
}

const studioSafeName = "studios"

// Get ..
func (ctrl StudiosController) Get(c *gin.Context) {
	var studios []m.Studio
	ctrl.context.Db.Find(&studios)
	response := make([]*m.StudioResponse, len(studios))
	for i, studio := range studios {
		response[i] = studio.Map()
	}

	c.JSON(http.StatusOK, &response)
}

// GetByIdent ..
func (ctrl StudiosController) GetByIdent(c *gin.Context) {
	var studio m.Studio
	identType, ident := h.DetectParamType(c.Param("ident"), "name")

	if ctrl.context.Db.First(&studio, fmt.Sprintf("`%s` = ?", identType), ident).RecordNotFound() {
		c.JSON(http.StatusNotFound, m.NewBaelorError("studio_not_found", nil))
	} else {
		c.JSON(http.StatusOK, studio.Map())
	}
}

// Post ..
func (ctrl StudiosController) Post(c *gin.Context) {
	// Validate Payload
	var studio m.Studio
	status, err := h.ValidateJSON(c, &studio, studioSafeName)
	if err != nil {
		c.JSON(status, &err)
		return
	}

	// Check studio is unique
	studio.NameSlug = h.GenerateSlug(studio.Name)
	if !ctrl.context.Db.First(&m.Studio{}, "name_slug = ?", studio.NameSlug).RecordNotFound() {
		c.JSON(http.StatusConflict, m.NewBaelorError("studio_already_exists", nil))
		return
	}

	// Insert into database
	studio.Init()
	if ctrl.context.Db.Create(&studio); ctrl.context.Db.NewRecord(studio) {
		c.JSON(http.StatusInternalServerError,
			m.NewBaelorError("unknown_error_creating_studio", nil))
		return
	}

	c.JSON(http.StatusCreated, studio.Map())
}

// NewStudiosController ..
func NewStudiosController(r *gin.RouterGroup, c *m.Context) {
	ctrl := new(StudiosController)
	ctrl.context = c

	r.GET("studios", ctrl.Get)
	r.GET("studios/:ident", ctrl.GetByIdent)
	r.POST("studios", middleware.BearerAuth(c), ctrl.Post)
}
