package repository

import (
	"context"
	"github.com/vlaship/book-catalog-go/internal/app/model"
	"github.com/vlaship/book-catalog-go/internal/app/types"
	"github.com/vlaship/book-catalog-go/internal/database"
	"github.com/vlaship/book-catalog-go/internal/logger"
	"github.com/vlaship/go-mask"
	"strings"
)

// UserRepository is an interface for user repository
type UserRepository struct {
	pool database.ConnPool
	log  logger.Logger
}

// NewUserRepository creates new user repository
func NewUserRepository(pool database.ConnPool, log logger.Logger) *UserRepository {
	return &UserRepository{
		pool: pool,
		log:  log.New("UserRepository"),
	}
}

func (r *UserRepository) l() logger.Logger {
	return r.log
}

func (r *UserRepository) p() database.ConnPool {
	return r.pool
}

const (
	userGetByUsername = `
	SELECT user_id, username, password, user_data FROM catalog.users
	WHERE username = $1 AND deleted = FALSE;
`
	userGetByID = `
	SELECT user_id, username, password, user_data FROM catalog.users
	WHERE user_id = $1 AND deleted = FALSE;
`
	userCreate = `
	INSERT INTO catalog.users (user_id, username, password, user_data)
	VALUES ($1, $2, $3, $4)
	RETURNING user_id, username, user_data;
`
	userUpdateStatus = `
	UPDATE catalog.users SET user_data = jsonb_set(user_data, '{status}',  to_jsonb($1::text), true)
 	WHERE username = $2 AND deleted = FALSE;
`
	userUpdatePassword = `
	UPDATE catalog.users SET password = $1 WHERE username = $2 AND deleted = FALSE AND status = '';
`
	userUpdateInfo = `
	UPDATE catalog.users SET user_data = $1 WHERE user_id = $2 AND deleted = FALSE;
`
)

const (
	entityNameUser = "user"
)

// GetUserByUsername get user by username
func (r *UserRepository) GetUserByUsername(ctx context.Context, username types.Username) (*model.User, error) {
	r.log.Dbg().Ctx(ctx).Values("username", mask.String(string(username))).Msg("GetUserByUsername")

	req := entity[model.User]{
		query:      userGetByUsername,
		entityName: "user",
		args:       []any{r.lower(username)},
		destinations: func(out *model.User) []any {
			return []any{
				&out.ID,
				&out.Username,
				&out.Password,
				&out.Data,
			}
		},
	}

	return getOne(ctx, r, req)
}

// GetUserByID get user by signin
func (r *UserRepository) GetUserByID(ctx context.Context, userID types.UserID) (*model.User, error) {
	r.log.Dbg().Ctx(ctx).Values("userID", userID).Msg("GetUserBySignin")

	req := entity[model.User]{
		query:      userGetByID,
		entityName: "user",
		args:       []any{userID},
		destinations: func(out *model.User) []any {
			return []any{
				&out.ID,
				&out.Username,
				&out.Password,
				&out.Data,
			}
		},
	}

	return getOne(ctx, r, req)
}

// Create create user
func (r *UserRepository) Create(ctx context.Context, input model.User) (*model.User, error) {
	r.log.Dbg().Ctx(ctx).Values("username", mask.String(string(input.Username))).Msg("CreateUser")

	req := entity[model.User]{
		query:      userCreate,
		entityName: entityNameUser,
		args: []any{
			input.ID,
			r.lower(input.Username),
			input.Password,
			input.Data,
		},
		destinations: func(output *model.User) []any {
			return []any{
				&output.ID,
				&output.Username,
				&output.Data,
			}
		},
	}

	return create(ctx, r, req)
}

// UpdateStatus activate user
func (r *UserRepository) UpdateStatus(ctx context.Context, user model.User) error {
	r.log.Dbg().Ctx(ctx).Values("username", mask.String(string(user.Username))).Msg("Activate")

	req := execRequest{
		query:      userUpdateStatus,
		entityName: entityNameUser,
		args:       []any{user.Data.Status, r.lower(user.Username)},
	}

	return exec(ctx, r, req)
}

// UpdatePassword update user password
func (r *UserRepository) UpdatePassword(ctx context.Context, user model.User) error {
	r.log.Dbg().Ctx(ctx).Values("username", mask.String(string(user.Username))).Msg("ReplacePassword")

	req := execRequest{
		query:      userUpdatePassword,
		entityName: entityNameUser,
		args:       []any{user.Password, r.lower(user.Username)},
	}

	return exec(ctx, r, req)
}

// UpdateInfo update user
func (r *UserRepository) UpdateInfo(ctx context.Context, user model.User, userID types.UserID) error {
	r.log.Dbg().Ctx(ctx).Values("user", user).Msg("UpdateInfo")

	req := execRequest{
		query:      userUpdateInfo,
		entityName: entityNameUser,
		args: []any{
			user.Data,
			userID,
		},
	}

	return exec(ctx, r, req)
}

func (r *UserRepository) lower(username types.Username) string {
	return strings.ToLower(string(username))
}
