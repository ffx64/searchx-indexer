package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"searchx-indexer/entity"
	"searchx-indexer/service"

	"github.com/gin-gonic/gin"
)

type ComboListController struct {
	Service *service.ComboListService
}

type BulkUploadRequest struct {
	Hash     string                         `json:"hash"`
	Metadata entity.ComboListMetadataEntity `json:"metadata"`
	Data     []entity.ComboListDataEntity   `json:"data"`
}

func NewComboListController(svc *service.ComboListService) *ComboListController {
	return &ComboListController{Service: svc}
}

func (c *ComboListController) BulkUpload(ctx *gin.Context) {
	authKey := ctx.GetHeader("authorization")
	if authKey == "" {
		log.Println("[combolist] not authorization")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorization"})
		return
	}

	agent, err := c.Service.AgentRepo.FindByAuthKey(authKey)
	if err != nil {
		log.Println("[combolist] failed to find agent with provided authorization:", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorization"})
		return
	}

	if agent.AgentStatus != "active" {
		log.Println("[combolist] agent is not active")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorization"})
		return
	}

	if agent.Platform != "combolist" {
		log.Println("[combolist] not found platform to agent")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorization"})
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("[combolist] failed to read request body")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorization"})
		return
	}

	var req BulkUploadRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Println("[combolist] invalid json payload")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorization"})
		return
	}

	ip := ctx.ClientIP()
	if err := c.Service.AgentRepo.UpdateActivity(authKey, ip); err != nil {
		log.Println("[combolist] failed to update agent activity")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorization"})
		return
	}

	if err := c.Service.BulkInsert(req.Hash, req.Metadata, req.Data); err != nil {
		log.Println("[combolist] failed to save bulk data:", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorization"})
		return
	}

	if err := c.Service.AgentRepo.IncrementDataProcessed(authKey, len(req.Data)); err != nil {
		log.Println("[combolist] failed to increment data processed count")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorization"})
		return
	}

	log.Println("[combolist] bulk upload successful")
	ctx.JSON(http.StatusOK, gin.H{"message": "bulk upload successful"})
}
