package repositories

import (
    "context"

    "deployflow/user-service/internal/models"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

type UserRepository interface {
    CreateOrGetByRegistrationID(ctx context.Context, registrationID string, user *models.User) (uint, error)
    FindByID(ctx context.Context, userID uint) (*models.User, error)
    FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (repo *userRepository) CreateOrGetByRegistrationID(
    ctx context.Context,
    registrationID string,
    user *models.User,
) (uint, error) {
    user.RegistrationID = registrationID

    err := repo.db.WithContext(ctx).
        Clauses(clause.OnConflict{
            Columns: []clause.Column{
                {Name: "registration_id"},
            },
            DoNothing: true,
        }).
        Create(user).Error

    if err != nil {
        return 0, err
    }

    var existingUser models.User

    if err := repo.db.WithContext(ctx).
        Where("registration_id = ?", registrationID).
        First(&existingUser).Error; err != nil {
        return 0, err
    }

    return existingUser.ID, nil
}

func (repo *userRepository) FindByID(
    ctx context.Context,
    userID uint,
) (*models.User, error) {
    var user models.User

    if err := repo.db.WithContext(ctx).
        First(&user, userID).Error; err != nil {
        return nil, err
    }

    return &user, nil
}

func (repo *userRepository) FindByEmail(
    ctx context.Context,
    email string,
) (*models.User, error) {
    var user models.User

    if err := repo.db.WithContext(ctx).
        Where("email = ?", email).
        First(&user).Error; err != nil {
        return nil, err
    }

    return &user, nil
}

var _ UserRepository = (*userRepository)(nil)