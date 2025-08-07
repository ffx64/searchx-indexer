package main

import (
	"log"
	"searchx-indexer/controller"
	"searchx-indexer/internal/banner"
	"searchx-indexer/internal/multidb"
	"searchx-indexer/repository"
	"searchx-indexer/service"

	"github.com/gin-gonic/gin"
)

func init() {
	banner.Show()
}

func main() {
	database := multidb.New()

	if err := database.ConnectDatabase("agent", "postgres://docker:docker@localhost:5435/searchx?sslmode=disable"); err != nil {
		log.Fatal(err)
	}

	if err := database.ConnectDatabase("combolist", "postgres://docker:docker@localhost:5433/searchx_combolist?sslmode=disable"); err != nil {
		log.Fatal(err)
	}

	agentdb, err := database.GetConnection("agent")
	if err != nil {
		log.Fatal(err)
	}

	combolistdb, err := database.GetConnection("combolist")
	if err != nil {
		log.Fatal(err)
	}

	agentrepository := repository.NewAgentRepository(agentdb)
	combolistrepository := repository.NewComboListRepository(combolistdb)

	combolistservice := service.NewComboListService(combolistrepository, agentrepository)
	combolistcontroller := controller.NewComboListController(combolistservice)

	r := gin.Default()
	v1 := r.Group("/api/v1")

	v1.POST("/combolist/bulk", combolistcontroller.BulkUpload)

	r.Run(":9090")
}
