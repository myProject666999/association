package database

import (
	"log"
	"time"

	"association/config"
	"association/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/glebarez/sqlite"
)

var DB *gorm.DB

func InitDatabase() {
	var err error

	if config.AppConfig.Database.Driver == "sqlite" {
		DB, err = gorm.Open(sqlite.Open(config.AppConfig.Database.DSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Club{},
		&models.ClubMember{},
		&models.Activity{},
		&models.ActivityRegistration{},
		&models.ActivityComment{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	seedData()

	log.Println("Database initialized successfully")
}

func seedData() {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	universityAdmin := &models.User{
		Username:   "admin",
		Password:   string(hashedPassword),
		Name:       "校级管理员",
		Email:      "admin@university.edu",
		Phone:      "13800138000",
		StudentID:  "U000001",
		Department: "校团委",
		Role:       models.RoleUniversityAdmin,
		Status:     1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	DB.Create(universityAdmin)

	deptAdmin := &models.User{
		Username:   "dept_admin",
		Password:   string(hashedPassword),
		Name:       "院级管理员",
		Email:      "dept_admin@university.edu",
		Phone:      "13800138001",
		StudentID:  "D000001",
		Department: "计算机学院",
		Role:       models.RoleDeptAdmin,
		Status:     1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	DB.Create(deptAdmin)

	student := &models.User{
		Username:   "student",
		Password:   string(hashedPassword),
		Name:       "测试学生",
		Email:      "student@university.edu",
		Phone:      "13800138002",
		StudentID:  "S2024001",
		Department: "计算机学院",
		Major:      "软件工程",
		Grade:      "2024级",
		Role:       models.RoleStudent,
		Status:     1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	DB.Create(student)

	log.Println("Default users created successfully")
}
