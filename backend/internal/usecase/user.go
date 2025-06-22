package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kaitoyama/kaitoyama-server-template/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	querier db.Querier
}

func NewUserUsecase(querier db.Querier) *UserUsecase {
	return &UserUsecase{
		querier: querier,
	}
}

type CreateUserRequest struct {
	Username string
	Password string
}

type CreateUserResponse struct {
	UserID   int64
	Username string
}

func (u *UserUsecase) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// ユーザーを作成
	result, err := u.querier.CreateUserWithPassword(ctx, db.CreateUserWithPasswordParams{
		Username:     req.Username,
		PasswordHash: sql.NullString{String: string(hashedPassword), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{
		UserID:   userID,
		Username: req.Username,
	}, nil
}

func (u *UserUsecase) CreateUserWithoutPassword(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	// パスワードなしでユーザーを作成（X-Forwarded-User用）
	result, err := u.querier.CreateUser(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{
		UserID:   userID,
		Username: req.Username,
	}, nil
}

func (u *UserUsecase) GetUserByUsername(ctx context.Context, username string) (*db.User, error) {
	user, err := u.querier.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (u *UserUsecase) AuthenticateUser(ctx context.Context, username, password string) (*db.User, error) {
	// ユーザー情報を取得
	user, err := u.querier.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// パスワードを検証
	if !user.PasswordHash.Valid {
		return nil, errors.New("user has no password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}
