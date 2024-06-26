package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       string  `gorm:"type:string;primary_key"`
	Name     string  `gorm:"type:varchar(100);not null"`
	Email    string  `gorm:"type:varchar(100);unique;not null"`
	Password string  `gorm:"type:varchar(100);not null"`
	Role     string  `gorm:"type:varchar(50);not null"`
	Provider *string `gorm:"type:varchar(50);default:'local';not null"`
	Photo    *string `gorm:"not null;default:'default.png'"`
	Verified *bool   `gorm:"not null;default:false"`
	// VerificationCode string    `json:"verification_code,omitempty"` // ? This is for email verification
	VerificationCode   string    `gorm:"type:varchar(100);"`
	PasswordResetToken string    `gorm:"type:varchar(100);"`
	PasswordResetAt    time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	CreatedAt          time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	Task               []Task
}

type Task struct {
	ID          string    `gorm:"type:string;primary_key"`
	UserID      string    `gorm:"type:string;not null"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`                 // User ID of the assigned user
	Status      string    `gorm:"type:varchar(50);not null"` // Status: To Do, In Progress, Completed
	Deadline    time.Time `gorm:"type:datetime"`
	CreatedAt   time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
}

func (u *Task) BeforeCreate(tx *gorm.DB) error {
	// Generate a new UUID
	uuid, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	// Set the UUID as the primary key
	// u.ID = uuid.String()
	u.ID = uuid.String()

	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Generate a new UUID
	uuid, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	// Set the UUID as the primary key
	// u.ID = uuid.String()
	u.ID = uuid.String()

	return nil
}

type SignUpInput struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,min=8"`
	Photo           string `json:"photo"`
	Role            string `json:"role" validate:"required"`
}

type SignInInput struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TaskInput struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"required"`
	// Deadline    time.Time `json:"deadline" validate:"required"`
}

// user response
type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type taskResponse struct {
	ID          uuid.UUID `json:"id,omitempty"`
	UserID      uuid.UUID `json:"user_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	Deadline    time.Time `json:"deadline,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// filter task
func FilterTaskRecord(task *Task) taskResponse {
	id := task.ID
	userID := task.UserID
	return taskResponse{
		ID:          uuid.MustParse(id),
		UserID:      uuid.MustParse(userID),
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Deadline:    task.Deadline,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

// filter user
func FilterUserRecord(user *User) UserResponse {
	id := user.ID
	return UserResponse{
		ID:        uuid.MustParse(id),
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Photo:     *user.Photo,
		Provider:  *user.Provider,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

// VALIDATE STRUCT
func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

// ForgotPasswordInput struct
type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

// ResetPasswordInput struct
type ResetPasswordInput struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}
