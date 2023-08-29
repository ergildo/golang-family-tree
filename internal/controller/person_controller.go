package controller

import (
	"golang-family-tree/internal/domain/model"
	"golang-family-tree/internal/domain/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewPersonController(service service.PersonService) PersonController {
	return PersonController{service: service}
}

type PersonController struct {
	service service.PersonService
}

func (p PersonController) Routes(group *gin.RouterGroup) {
	personGroup := group.Group("person")
	personGroup.POST("", p.Add)
	personGroup.GET("", p.FindAll)
	personGroup.GET("/ascendants/:id", p.FindAscendantsById)
}

func (p PersonController) Add(ctx *gin.Context) {
	var person model.Person
	if err := ctx.ShouldBindJSON(&person); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	if err := p.service.Add(ctx, person); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.Status(http.StatusCreated)

}

func (p PersonController) FindAscendantsById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ascendants, err := p.service.FindAscendantsById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, ascendants)
}

func (p PersonController) FindAll(ctx *gin.Context) {
	people, err := p.service.FindAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, people)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
