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
    Update(ctx context.Context, userID uint, updates map[string]interface{}) error
    FindAll(ctx context.Context, limit, offset int) ([]models.User, int64, error)
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (repo *userRepository) CreateOrGetByRegistrationID(ctx context.Context, registrationID string, user *models.User) (uint, error) {
    user.RegistrationID = registrationID

    err := repo.db.WithContext(ctx).
        Clauses(clause.OnConflict{
            Columns: []clause.Column{{Name: "registration_id"}},
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

func (repo *userRepository) FindByID(ctx context.Context, userID uint) (*models.User, error) {
    var user models.User

    err := repo.db.WithContext(ctx).First(&user, userID).Error
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (repo *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User

    err := repo.db.WithContext(ctx).
        Where("email = ?", email).
        First(&user).Error

    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (repo *userRepository) Update(ctx context.Context, userID uint, updates map[string]interface{}) error {
    result := repo.db.WithContext(ctx).
        Model(&models.User{}).
        Where("id = ?", userID).
        Updates(updates)

    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }

    return nil
}

func (repo *userRepository) FindAll(ctx context.Context, limit, offset int) ([]models.User, int64, error) {
    var users []models.User
    var total int64

    query := repo.db.WithContext(ctx).Model(&models.User{})

    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    if err := query.
        Limit(limit).
        Offset(offset).
        Order("id DESC").
        Find(&users).Error; err != nil {
        return nil, 0, err
    }

    return users, total, nil
}

var _ UserRepository = (*userRepository)(nil)