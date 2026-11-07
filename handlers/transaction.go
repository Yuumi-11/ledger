package handlers

import (
	"ledger/models"
	"strings"
	"time"

	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	DB *gorm.DB
}
type CreateTransactionRequest struct {
	Kind        string `json:"kind" binding:"required,oneof=income expense"`
	AmountCents int64  `json:"amount_cents" binding:"required,gt=0"`
	Category    string `json:"category" binding:"required,max=50"`
	OccurredOn  string `json:"occurred_on" binding:"required"`
	Note        string `json:"note" binding:"max=255"`
}

func (h TransactionHandler) List(c *gin.Context) {
	records := make([]models.Transaction, 0)

	result := h.DB.Find(&records)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "查询账目失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "查询成功",
		"data":    records,
	})
}
func (h TransactionHandler) Create(c *gin.Context) {
	var input CreateTransactionRequest

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "输入有误，请检查收支类型、金额和字段长度",
		})
		return
	}

	input.Category = strings.TrimSpace(input.Category)
	if input.Category == "" {
		c.JSON(400, gin.H{
			"message": "分类不能全是空白",
		})
		return
	}

	date, err := time.ParseInLocation(
		"2006-01-02",
		input.OccurredOn,
		time.Local,
	)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "日期必须是有效的 YYYY-MM-DD，例如 2026-11-07",
		})
		return
	}

	record := models.Transaction{
		KIND:        input.Kind,
		AmountCents: input.AmountCents,
		Category:    input.Category,
		OccurredOn:  date,
		Note:        input.Note,
	}

	result := h.DB.Create(&record)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "保存账目失败",
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "新增成功",
		"data":    record,
	})
}
func (h TransactionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(400, gin.H{"message": "账目编号必须是正整数"})
		return
	}

	var input CreateTransactionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"message": "输入有误，请检查账目字段"})
		return
	}

	input.Category = strings.TrimSpace(input.Category)
	if input.Category == "" {
		c.JSON(400, gin.H{"message": "分类不能全是空白"})
		return
	}

	date, err := time.ParseInLocation(
		"2006-01-02", input.OccurredOn, time.Local,
	)
	if err != nil {
		c.JSON(400, gin.H{"message": "日期必须是有效的 YYYY-MM-DD"})
		return
	}

	var record models.Transaction
	result := h.DB.First(&record, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{"message": "这笔账不存在"})
		return
	}
	if result.Error != nil {
		c.JSON(500, gin.H{"message": "查询账目失败"})
		return
	}

	result = h.DB.Model(&record).Updates(map[string]any{
		"kind":         input.Kind,
		"amount_cents": input.AmountCents,
		"category":     input.Category,
		"occurred_on":  date,
		"note":         input.Note,
	})
	if result.Error != nil {
		c.JSON(500, gin.H{"message": "修改账目失败"})
		return
	}

	c.JSON(200, gin.H{
		"message": "修改成功",
		"id":      id,
	})
}
func (h TransactionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(400, gin.H{
			"message": "账目编号必须是正整数",
		})
		return
	}

	result := h.DB.Where("id = ?", id).
		Delete(&models.Transaction{})

	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "删除账目失败",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{
			"message": "这笔账不存在",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "删除成功",
		"id":      id,
	})
}
