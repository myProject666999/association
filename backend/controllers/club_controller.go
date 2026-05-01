package controllers

import (
	"association/database"
	"association/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateClubRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Logo        string `json:"logo"`
}

type UpdateClubRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Logo        string `json:"logo"`
}

type MemberReviewRequest struct {
	Status   int    `json:"status" binding:"required"`
	Position string `json:"position"`
}

func GetClubs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")
	category := c.Query("category")
	status := c.Query("status")

	var clubs []models.Club
	var total int64

	query := database.DB.Model(&models.Club{})

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if status != "" {
		statusInt, _ := strconv.Atoi(status)
		query = query.Where("status = ?", statusInt)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&clubs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取社团列表失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":      clubs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetClubByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var club models.Club
	if err := database.DB.Preload("Members").Preload("Members.User").First(&club, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "社团不存在",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    club,
	})
}

func CreateClub(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateClubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var existingClub models.Club
	if err := database.DB.Where("name = ?", req.Name).First(&existingClub).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "社团名称已存在",
			"data":    nil,
		})
		return
	}

	club := models.Club{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Logo:        req.Logo,
		FoundedAt:   time.Now(),
		Status:      1,
		CreatedBy:   userID.(uint),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := database.DB.Create(&club).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建社团失败",
			"data":    nil,
		})
		return
	}

	now := time.Now()
	member := models.ClubMember{
		ClubID:   club.ID,
		UserID:   userID.(uint),
		Position: "president",
		Status:   1,
		JoinedAt: &now,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	database.DB.Create(&member)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    club,
	})
}

func UpdateClub(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var req UpdateClubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var club models.Club
	if err := database.DB.First(&club, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "社团不存在",
			"data":    nil,
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Logo != "" {
		updates["logo"] = req.Logo
	}
	updates["updated_at"] = time.Now()

	if err := database.DB.Model(&club).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新社团失败",
			"data":    nil,
		})
		return
	}

	database.DB.First(&club, id)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    club,
	})
}

func DeleteClub(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var club models.Club
	if err := database.DB.First(&club, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "社团不存在",
			"data":    nil,
		})
		return
	}

	if err := database.DB.Delete(&club).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除社团失败",
			"data":    nil,
		})
		return
	}

	database.DB.Where("club_id = ?", id).Delete(&models.ClubMember{})

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data":    nil,
	})
}

func ApplyJoinClub(c *gin.Context) {
	userID, _ := c.Get("user_id")
	clubID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var club models.Club
	if err := database.DB.First(&club, clubID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "社团不存在",
			"data":    nil,
		})
		return
	}

	var existingMember models.ClubMember
	if err := database.DB.Where("club_id = ? AND user_id = ?", clubID, userID).First(&existingMember).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "已提交申请或已是社团成员",
			"data":    nil,
		})
		return
	}

	member := models.ClubMember{
		ClubID:    uint(clubID),
		UserID:    userID.(uint),
		Position:  "member",
		Status:    0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "申请失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "申请成功，等待审核",
		"data":    member,
	})
}

func GetClubMembers(c *gin.Context) {
	clubID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	status := c.Query("status")

	var members []models.ClubMember
	query := database.DB.Where("club_id = ?", clubID).Preload("User")

	if status != "" {
		statusInt, _ := strconv.Atoi(status)
		query = query.Where("status = ?", statusInt)
	}

	if err := query.Order("created_at DESC").Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取成员列表失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    members,
	})
}

func ReviewMemberApplication(c *gin.Context) {
	memberID, err := strconv.ParseUint(c.Param("member_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var req MemberReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var member models.ClubMember
	if err := database.DB.First(&member, memberID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "成员申请不存在",
			"data":    nil,
		})
		return
	}

	member.Status = req.Status
	if req.Position != "" {
		member.Position = req.Position
	}
	if req.Status == 1 {
		now := time.Now()
		member.JoinedAt = &now
	}
	member.UpdatedAt = time.Now()

	if err := database.DB.Save(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "审核失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "审核成功",
		"data":    member,
	})
}

func RemoveMember(c *gin.Context) {
	memberID, err := strconv.ParseUint(c.Param("member_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var member models.ClubMember
	if err := database.DB.First(&member, memberID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "成员不存在",
			"data":    nil,
		})
		return
	}

	if err := database.DB.Delete(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "移除成员失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "移除成功",
		"data":    nil,
	})
}
