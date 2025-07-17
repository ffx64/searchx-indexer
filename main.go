package main

import (
	"log"
	"os"
	"runtime"
	"searchx-indexer/controller"
	"searchx-indexer/repository"
	"searchx-indexer/service"
	"searchx-indexer/storage"

	"github.com/gin-gonic/gin"
)

func Show() {
	isWindows := runtime.GOOS == "windows"
	isDumbTerminal := os.Getenv("TERM") == "dumb"

	useColors := isWindows && isDumbTerminal

	banner := `
███████╗███████╗ █████╗ ██████╗  ██████╗██╗  ██╗██╗  ██╗
██╔════╝██╔════╝██╔══██╗██╔══██╗██╔════╝██║  ██║╚██╗██╔╝
███████╗█████╗  ███████║██████╔╝██║     ███████║ ╚███╔╝ 
╚════██║██╔══╝  ██╔══██║██╔══██╗██║     ██╔══██║ ██╔██╗ 
███████║███████╗██║  ██║██║  ██║╚██████╗██║  ██║██╔╝ ██╗
╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝
                  Indexer v1.0.0
=======================================================
   SentielX @ 2024 | Data Processing via Socket Stream
=======================================================
`

	if useColors {
		banner = `
\033[1;36m███████╗███████╗ █████╗ ██████╗  ██████╗██╗  ██╗██╗  ██╗\033[0m
\033[1;36m██╔════╝██╔════╝██╔══██╗██╔══██╗██╔════╝██║  ██║╚██╗██╔╝\033[0m
\033[1;36m███████╗█████╗  ███████║██████╔╝██║     ███████║ ╚███╔╝ \033[0m
\033[1;36m╚════██║██╔══╝  ██╔══██║██╔══██╗██║     ██╔══██║ ██╔██╗ \033[0m
\033[1;36m███████║███████╗██║  ██║██║  ██║╚██████╗██║  ██║██╔╝ ██╗\033[0m
\033[1;36m╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝\033[0m
\033[1;33m                  Indexer v1.0.0\033[0m
\033[1;37m=======================================================\033[0m
\033[1;35m   SentielX @ 2024 | Data Processing via Socket Stream\033[0m
\033[1;37m=======================================================\033[0m
`
	}

	println(banner)
}

func main() {
	Show()

	dbManager := storage.NewDBManager()
	err := dbManager.AddDB("agent", "postgres://docker:docker@localhost:5435/searchx?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = dbManager.AddDB("combolist", "postgres://docker:docker@localhost:5433/searchx_combolist?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	agentDB, err := dbManager.GetDB("agent")
	if err != nil {
		log.Fatal(err)
	}
	combolistDB, err := dbManager.GetDB("combolist")
	if err != nil {
		log.Fatal(err)
	}

	agentRepo := repository.NewAgentRepository(agentDB)
	comboRepo := repository.NewComboListRepository(combolistDB)
	svc := service.NewComboListService(comboRepo, agentRepo)
	ctrl := controller.NewComboListController(svc)

	r := gin.Default()

	v1 := r.Group("/api/v1")

	v1.POST("/combolist/bulk", ctrl.BulkUpload)

	r.Run(":9090")
}
