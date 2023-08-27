package controller

import (
	"context"
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
	//TODO: improve error handling
	var person model.Person
	if err := ctx.ShouldBindJSON(&person); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	if err := p.service.Add(person); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, person)

}

func (p PersonController) FindAscendantsById(ctx *gin.Context) {
	//TODO: improve error handling
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ascendants, err := p.service.FindAscendantsById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, ascendants)
}

func (p PersonController) FindAll(ctx *gin.Context) {
	withContext(ctx)
	people, err := p.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, people)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func withContext(ctx context.Context) {

}