package accounts

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/w0ikid/highload-auth-go/internal/repository"
	"github.com/w0ikid/highload-auth-go/internal/usecase"
)

type GetProfileUsecase struct {
	usecase.BaseUsecase
	userRepo repository.IUserRepo
}

func NewGetProfileUsecase(base usecase.BaseUsecase, userRepo repository.IUserRepo) GetProfileUsecase {
	return GetProfileUsecase{
		BaseUsecase: base,
		userRepo:    userRepo,
	}
}

type ProfileResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

func (uc *GetProfileUsecase) Execute(ctx context.Context, userIDStr string) (*ProfileResponse, error) {
	uc.Logger.Infow("starting GetProfileUsecase execution", "user_id", userIDStr)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user id format")
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.Logger.Errorw("user not found", "user_id", userID, "error", err)
		return nil, errors.New("user not found")
	}

	uc.Logger.Infow("GetProfileUsecase executed successfully", "user_id", userID)
	return &ProfileResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
