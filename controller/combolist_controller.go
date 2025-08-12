package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"searchx-indexer/entity"
	"searchx-indexer/security"
	"searchx-indexer/service"

	"github.com/gin-gonic/gin"
)

type ComboListController struct {
	Service *service.CombolistService
}

type BulkUploadRequest struct {
	Hash     string                         `json:"hash"`
	Metadata entity.CombolistMetadataEntity `json:"metadata"`
	Data     []entity.CombolistDataEntity   `json:"data"`
}

func NewComboListController(svc *service.CombolistService) *ComboListController {
	return &ComboListController{Service: svc}
}

func (c *ComboListController) BulkUpload(ctx *gin.Context) {
	authkey := ctx.GetHeader("authorization")

	_, err := security.AuthAgent(c.Service.AgentRepository, authkey, "combolist")
	if err != nil {
		log.Println("[combolist]", err.Error())
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorization"})
		return
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
	if err := c.Service.AgentRepository.UpdateActivity(authkey, ip); err != nil {
		log.Println("[combolist] failed to update agent activity")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorization"})
		return
	}

	if err := c.Service.BulkInsert(req.Hash, req.Metadata, req.Data); err != nil {
		log.Println("[combolist] failed to save bulk data:", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorization"})
		return
	}

	if err := c.Service.AgentRepository.IncrementDataProcessed(authkey, len(req.Data)); err != nil {
		log.Println("[combolist] failed to increment data processed count")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorization"})
		return
	}

	log.Println("[combolist] bulk upload successful")
	ctx.JSON(http.StatusOK, gin.H{"message": "bulk upload successful"})
}
