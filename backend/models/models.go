package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleStudent       Role = "student"
	RoleDeptAdmin     Role = "dept_admin"
	RoleUniversityAdmin Role = "university_admin"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password     string         `gorm:"size:255;not null" json:"-"`
	Name         string         `gorm:"size:50;not null" json:"name"`
	Email        string         `gorm:"size:100" json:"email"`
	Phone        string         `gorm:"size:20" json:"phone"`
	StudentID    string         `gorm:"uniqueIndex;size:50" json:"student_id"`
	Department   string         `gorm:"size:100" json:"department"`
	Major        string         `gorm:"size:100" json:"major"`
	Grade        string         `gorm:"size:20" json:"grade"`
	Role         Role           `gorm:"size:20;default:'student'" json:"role"`
	Status       int            `gorm:"default:1" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Club struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Category    string         `gorm:"size:50" json:"category"`
	FoundedAt   time.Time      `json:"founded_at"`
	Logo        string         `gorm:"size:255" json:"logo"`
	Status      int            `gorm:"default:1" json:"status"`
	CreatedBy   uint           `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Members     []ClubMember   `gorm:"foreignKey:ClubID" json:"members,omitempty"`
}

type ClubMember struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ClubID    uint           `gorm:"not null" json:"club_id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	Position  string         `gorm:"size:50;default:'member'" json:"position"`
	Status    int            `gorm:"default:0" json:"status"`
	JoinedAt  *time.Time     `json:"joined_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Club      Club           `gorm:"foreignKey:ClubID" json:"club,omitempty"`
}

type Activity struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Title           string         `gorm:"size:200;not null" json:"title"`
	Description     string         `gorm:"type:text" json:"description"`
	Location        string         `gorm:"size:200" json:"location"`
	StartTime       time.Time      `json:"start_time"`
	EndTime         time.Time      `json:"end_time"`
	RegistrationStart time.Time    `json:"registration_start"`
	RegistrationEnd   time.Time    `json:"registration_end"`
	MaxParticipants int            `gorm:"default:0" json:"max_participants"`
	CurrentParticipants int        `gorm:"default:0" json:"current_participants"`
	ClubID          *uint          `json:"club_id"`
	OrganizerID     uint           `json:"organizer_id"`
	Status          int            `gorm:"default:0" json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Club            *Club          `gorm:"foreignKey:ClubID" json:"club,omitempty"`
	Organizer       User           `gorm:"foreignKey:OrganizerID" json:"organizer,omitempty"`
	Registrations   []ActivityRegistration `gorm:"foreignKey:ActivityID" json:"registrations,omitempty"`
}

type ActivityRegistration struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ActivityID uint          `gorm:"not null" json:"activity_id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	Status    int            `gorm:"default:0" json:"status"`
	RegistratedAt time.Time   `json:"registrated_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Activity  Activity       `gorm:"foreignKey:ActivityID" json:"activity,omitempty"`
}

type ActivityComment struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ActivityID uint           `gorm:"not null" json:"activity_id"`
	UserID     uint           `gorm:"not null" json:"user_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Status     int            `gorm:"default:1" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	User       User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
