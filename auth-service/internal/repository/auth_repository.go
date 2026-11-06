
package repository

import (
	"context"
	"errors"
	"time"

	"deployflow/auth-service/internal/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateCredential(ctx context.Context, credential *model.Credential) error
	FindCredentialByEmail(ctx context.Context, email string) (*model.Credential, error)
	FindCredentialByID(ctx context.Context, credentialID uint) (*model.Credential, error)
	UpdateCredential(ctx context.Context, credential *model.Credential) error

	CreateSession(ctx context.Context, session *model.AuthSession) error
	FindSessionByTokenHash(ctx context.Context, tokenHash string) (*model.AuthSession, error)
	RevokeSession(ctx context.Context, sessionID uint) error
	RevokeAllSessions(ctx context.Context, credentialID uint) error
	RotateSession(ctx context.Context, oldSessionID uint, newSession *model.AuthSession) error

	SavePasswordResetToken(ctx context.Context, token *model.PasswordResetToken) error
	FindPasswordResetToken(ctx context.Context, tokenHash string) (*model.PasswordResetToken, error)
	MarkPasswordResetTokenUsed(ctx context.Context, tokenID uint) error
	CompletePasswordReset(ctx context.Context, tokenID uint, credential *model.Credential) error

	SaveEmailVerificationToken(ctx context.Context, token *model.EmailVerificationToken) error
	FindEmailVerificationToken(ctx context.Context, tokenHash string) (*model.EmailVerificationToken, error)
	CompleteEmailVerification(ctx context.Context, tokenID uint, credential *model.Credential) error
	SetEmailVerified(ctx context.Context, credential *model.Credential) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (repo *authRepository) CreateCredential(ctx context.Context, credential *model.Credential) error {
	return repo.db.WithContext(ctx).Create(credential).Error
}

func (repo *authRepository) FindCredentialByEmail(ctx context.Context, email string) (*model.Credential, error) {
	var credential model.Credential

	err := repo.db.WithContext(ctx).Where("LOWER(email) = LOWER(?)", email).First(&credential).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &credential, nil
}

func (repo *authRepository) FindCredentialByID(ctx context.Context, credentialID uint) (*model.Credential, error) {
	var credential model.Credential

	err := repo.db.WithContext(ctx).First(&credential, "id = ?", credentialID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &credential, nil
}

func (repo *authRepository) UpdateCredential(ctx context.Context, credential *model.Credential) error {
	return repo.db.WithContext(ctx).Model(&model.Credential{}).
		Where("id = ?", credential.ID).
		Select("PasswordHash", "PasswordChangedAt", "UpdatedAt").
		Updates(credential).Error
}

func (repo *authRepository) CreateSession(ctx context.Context, session *model.AuthSession) error {
	return repo.db.WithContext(ctx).Create(session).Error
}

func (repo *authRepository) FindSessionByTokenHash(ctx context.Context, tokenHash string) (*model.AuthSession, error) {
	var session model.AuthSession

	err := repo.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (repo *authRepository) RevokeSession(ctx context.Context, sessionID uint) error {
	now := time.Now().UTC()

	result := repo.db.WithContext(ctx).Model(&model.AuthSession{}).
		Where("id = ? AND revoked_at IS NULL", sessionID).
		Update("revoked_at", now)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (repo *authRepository) RevokeAllSessions(ctx context.Context, credentialID uint) error {
	now := time.Now().UTC()

	return repo.db.WithContext(ctx).Model(&model.AuthSession{}).
		Where("credential_id = ? AND revoked_at IS NULL", credentialID).
		Update("revoked_at", now).Error
}

func (repo *authRepository) RotateSession(ctx context.Context, oldSessionID uint, newSession *model.AuthSession) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		result := tx.Model(&model.AuthSession{}).
			Where("id = ? AND revoked_at IS NULL AND expires_at > ?", oldSessionID, now).
			Update("revoked_at", now)

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return tx.Create(newSession).Error
	})
}

func (repo *authRepository) SavePasswordResetToken(ctx context.Context, token *model.PasswordResetToken) error {
	return repo.db.WithContext(ctx).Create(token).Error
}

func (repo *authRepository) FindPasswordResetToken(ctx context.Context, tokenHash string) (*model.PasswordResetToken, error) {
	var token model.PasswordResetToken

	err := repo.db.WithContext(ctx).
		Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", tokenHash, time.Now().UTC()).
		First(&token).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (repo *authRepository) MarkPasswordResetTokenUsed(ctx context.Context, tokenID uint) error {
	now := time.Now().UTC()

	result := repo.db.WithContext(ctx).Model(&model.PasswordResetToken{}).
		Where("id = ? AND used_at IS NULL AND expires_at > ?", tokenID, now).
		Update("used_at", now)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (repo *authRepository) CompletePasswordReset(ctx context.Context, tokenID uint, credential *model.Credential) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		tokenResult := tx.Model(&model.PasswordResetToken{}).
			Where("id = ? AND credential_id = ? AND used_at IS NULL AND expires_at > ?", tokenID, credential.ID, now).
			Update("used_at", now)

		if tokenResult.Error != nil {
			return tokenResult.Error
		}
		if tokenResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		credential.PasswordChangedAt = now
		credential.UpdatedAt = now

		credentialResult := tx.Model(&model.Credential{}).
			Where("id = ?", credential.ID).
			Select("PasswordHash", "PasswordChangedAt", "UpdatedAt").
			Updates(credential)

		if credentialResult.Error != nil {
			return credentialResult.Error
		}
		if credentialResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return tx.Model(&model.AuthSession{}).
			Where("credential_id = ? AND revoked_at IS NULL", credential.ID).
			Update("revoked_at", now).Error
	})
}

func (repo *authRepository) SaveEmailVerificationToken(ctx context.Context, token *model.EmailVerificationToken) error {
	return repo.db.WithContext(ctx).Create(token).Error
}

func (repo *authRepository) FindEmailVerificationToken(ctx context.Context, tokenHash string) (*model.EmailVerificationToken, error) {
	var token model.EmailVerificationToken

	err := repo.db.WithContext(ctx).
		Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", tokenHash, time.Now().UTC()).
		First(&token).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (repo *authRepository) CompleteEmailVerification(ctx context.Context, tokenID uint, credential *model.Credential) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		tokenResult := tx.Model(&model.EmailVerificationToken{}).
			Where("id = ? AND credential_id = ? AND used_at IS NULL AND expires_at > ?", tokenID, credential.ID, now).
			Update("used_at", now)

		if tokenResult.Error != nil {
			return tokenResult.Error
		}
		if tokenResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		credential.EmailVerifiedAt = &now
		credential.UpdatedAt = now

		credentialResult := tx.Model(&model.Credential{}).
			Where("id = ? AND email_verified_at IS NULL", credential.ID).
			Select("EmailVerifiedAt", "UpdatedAt").
			Updates(credential)

		if credentialResult.Error != nil {
			return credentialResult.Error
		}
		if credentialResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}

func (repo *authRepository) SetEmailVerified(ctx context.Context, credential *model.Credential) error {
	now := time.Now().UTC()
	credential.EmailVerifiedAt = &now
	credential.UpdatedAt = now

	result := repo.db.WithContext(ctx).Model(&model.Credential{}).
		Where("id = ? AND email_verified_at IS NULL", credential.ID).
		Select("EmailVerifiedAt", "UpdatedAt").
		Updates(credential)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

var _ AuthRepository = (*authRepository)(nil)
