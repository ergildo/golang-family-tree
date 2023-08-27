package application

import (
	"database/sql"
	datasource "golang-family-tree/infra/db/datasource"
	"golang-family-tree/internal/controller"
	"golang-family-tree/internal/domain/service"
	"golang-family-tree/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewApplication() Application {
	db, err := datasource.GetDb()
	if err != nil {
		panic(err)
	}
	personRepository := repository.NewPersonRepository(db)
	personService := service.NewPersonService(personRepository)

	return Application{db: db, personService: personService}
}

type Application struct {
	db            *sql.DB
	personService service.PersonService
}

func (p Application) Run(port string) error {
	defer p.db.Close()
	router := p.webRouter()
	return router.Run(port)
}

func (p Application) webRouter() *gin.Engine {
	router := gin.Default()
	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	personController := controller.NewPersonController(p.personService)
	api := router.Group("/api/v1")
	personController.Routes(api)
	return router
}
