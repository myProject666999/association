package controllers

import (
	"association/database"
	"association/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateActivityRequest struct {
	Title           string    `json:"title" binding:"required"`
	Description     string    `json:"description"`
	Location        string    `json:"location"`
	StartTime       time.Time `json:"start_time" binding:"required"`
	EndTime         time.Time `json:"end_time" binding:"required"`
	RegistrationStart time.Time `json:"registration_start" binding:"required"`
	RegistrationEnd   time.Time `json:"registration_end" binding:"required"`
	MaxParticipants int       `json:"max_participants"`
	ClubID          uint      `json:"club_id"`
}

type UpdateActivityRequest struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Location        string    `json:"location"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	RegistrationStart time.Time `json:"registration_start"`
	RegistrationEnd   time.Time `json:"registration_end"`
	MaxParticipants int       `json:"max_participants"`
}

type RegistrationReviewRequest struct {
	Status int `json:"status" binding:"required"`
}

type CommentRequest struct {
	Content string `json:"content" binding:"required"`
}

func GetActivities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	clubID := c.Query("club_id")

	var activities []models.Activity
	var total int64

	query := database.DB.Model(&models.Activity{}).Preload("Club").Preload("Organizer")

	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	if status != "" {
		statusInt, _ := strconv.Atoi(status)
		query = query.Where("status = ?", statusInt)
	}

	if clubID != "" {
		clubIDInt, _ := strconv.ParseUint(clubID, 10, 64)
		query = query.Where("club_id = ?", clubIDInt)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取活动列表失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"list":      activities,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetActivityByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var activity models.Activity
	if err := database.DB.Preload("Club").Preload("Organizer").First(&activity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "活动不存在",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    activity,
	})
}

func CreateActivity(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	activity := models.Activity{
		Title:           req.Title,
		Description:     req.Description,
		Location:        req.Location,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		RegistrationStart: req.RegistrationStart,
		RegistrationEnd:   req.RegistrationEnd,
		MaxParticipants: req.MaxParticipants,
		CurrentParticipants: 0,
		ClubID:          &req.ClubID,
		OrganizerID:     userID.(uint),
		Status:          0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if req.ClubID == 0 {
		activity.ClubID = nil
	}

	if err := database.DB.Create(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建活动失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功，等待审核",
		"data":    activity,
	})
}

func UpdateActivity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var req UpdateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var activity models.Activity
	if err := database.DB.First(&activity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "活动不存在",
			"data":    nil,
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if !req.StartTime.IsZero() {
		updates["start_time"] = req.StartTime
	}
	if !req.EndTime.IsZero() {
		updates["end_time"] = req.EndTime
	}
	if !req.RegistrationStart.IsZero() {
		updates["registration_start"] = req.RegistrationStart
	}
	if !req.RegistrationEnd.IsZero() {
		updates["registration_end"] = req.RegistrationEnd
	}
	if req.MaxParticipants > 0 {
		updates["max_participants"] = req.MaxParticipants
	}
	updates["updated_at"] = time.Now()

	if err := database.DB.Model(&activity).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新活动失败",
			"data":    nil,
		})
		return
	}

	database.DB.First(&activity, id)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    activity,
	})
}

func DeleteActivity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var activity models.Activity
	if err := database.DB.First(&activity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "活动不存在",
			"data":    nil,
		})
		return
	}

	if err := database.DB.Delete(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除活动失败",
			"data":    nil,
		})
		return
	}

	database.DB.Where("activity_id = ?", id).Delete(&models.ActivityRegistration{})
	database.DB.Where("activity_id = ?", id).Delete(&models.ActivityComment{})

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data":    nil,
	})
}

func ReviewActivity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var req RegistrationReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var activity models.Activity
	if err := database.DB.First(&activity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "活动不存在",
			"data":    nil,
		})
		return
	}

	activity.Status = req.Status
	activity.UpdatedAt = time.Now()

	if err := database.DB.Save(&activity).Error; err != nil {
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
		"data":    activity,
	})
}

func RegisterActivity(c *gin.Context) {
	userID, _ := c.Get("user_id")
	activityID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var activity models.Activity
	if err := database.DB.First(&activity, activityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "活动不存在",
			"data":    nil,
		})
		return
	}

	if activity.Status != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "活动未审核通过，无法报名",
			"data":    nil,
		})
		return
	}

	now := time.Now()
	if now.Before(activity.RegistrationStart) || now.After(activity.RegistrationEnd) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不在报名时间范围内",
			"data":    nil,
		})
		return
	}

	if activity.MaxParticipants > 0 && activity.CurrentParticipants >= activity.MaxParticipants {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "活动报名人数已满",
			"data":    nil,
		})
		return
	}

	var existingRegistration models.ActivityRegistration
	if err := database.DB.Where("activity_id = ? AND user_id = ?", activityID, userID).First(&existingRegistration).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "已报名或已提交报名申请",
			"data":    nil,
		})
		return
	}

	registration := models.ActivityRegistration{
		ActivityID:    uint(activityID),
		UserID:        userID.(uint),
		Status:        0,
		RegistratedAt: time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := database.DB.Create(&registration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "报名失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "报名成功，等待审核",
		"data":    registration,
	})
}

func GetActivityRegistrations(c *gin.Context) {
	activityID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	status := c.Query("status")

	var registrations []models.ActivityRegistration
	query := database.DB.Where("activity_id = ?", activityID).Preload("User")

	if status != "" {
		statusInt, _ := strconv.Atoi(status)
		query = query.Where("status = ?", statusInt)
	}

	if err := query.Order("created_at DESC").Find(&registrations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取报名列表失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    registrations,
	})
}

func ReviewRegistration(c *gin.Context) {
	registrationID, err := strconv.ParseUint(c.Param("registration_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var req RegistrationReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var registration models.ActivityRegistration
	if err := database.DB.First(&registration, registrationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "报名记录不存在",
			"data":    nil,
		})
		return
	}

	oldStatus := registration.Status
	registration.Status = req.Status
	registration.UpdatedAt = time.Now()

	if err := database.DB.Save(&registration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "审核失败",
			"data":    nil,
		})
		return
	}

	if oldStatus != 1 && req.Status == 1 {
		database.DB.Model(&models.Activity{}).Where("id = ?", registration.ActivityID).
			UpdateColumn("current_participants", gorm.Expr("current_participants + 1"))
	} else if oldStatus == 1 && req.Status != 1 {
		database.DB.Model(&models.Activity{}).Where("id = ?", registration.ActivityID).
			UpdateColumn("current_participants", gorm.Expr("current_participants - 1"))
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "审核成功",
		"data":    registration,
	})
}

func AddComment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	activityID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var activity models.Activity
	if err := database.DB.First(&activity, activityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "活动不存在",
			"data":    nil,
		})
		return
	}

	comment := models.ActivityComment{
		ActivityID: uint(activityID),
		UserID:     userID.(uint),
		Content:    req.Content,
		Status:     1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "评论失败",
			"data":    nil,
		})
		return
	}

	database.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "评论成功",
		"data":    comment,
	})
}

func GetActivityComments(c *gin.Context) {
	activityID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var comments []models.ActivityComment
	if err := database.DB.Where("activity_id = ? AND status = 1", activityID).
		Preload("User").Order("created_at DESC").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取评论失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    comments,
	})
}

func DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	var comment models.ActivityComment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "评论不存在",
			"data":    nil,
		})
		return
	}

	if err := database.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除评论失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data":    nil,
	})
}
