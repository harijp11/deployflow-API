
package repositories

import (
	"context"
	"errors"
	"time"

	"deployflow/auth-service/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthRepository interface {
	CreateCredential(ctx context.Context, credential *models.Credential) error
	FindCredentialByEmail(ctx context.Context, email string) (*models.Credential, error)
	FindCredentialByID(ctx context.Context, credentialID uint) (*models.Credential, error)
	UpdateCredential(ctx context.Context, credential *models.Credential) error

	CreateSession(ctx context.Context, session *models.AuthSession) error
	FindSessionByTokenHash(ctx context.Context, tokenHash string) (*models.AuthSession, error)
	RevokeSession(ctx context.Context, sessionID uint) error
	RevokeAllSessions(ctx context.Context, credentialID uint) error
	RotateSession(ctx context.Context, oldSessionID uint, newSession *models.AuthSession) error

	SavePasswordResetToken(ctx context.Context, token *models.PasswordResetToken) error
	FindPasswordResetToken(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error)
	MarkPasswordResetTokenUsed(ctx context.Context, tokenID uint) error
	CompletePasswordReset(ctx context.Context, tokenID uint, credential *models.Credential) error

	CreateOrUpdateRegistrationIntent(ctx context.Context, intent *models.RegistrationIntent) error
	FindRegistrationIntentByTokenHash(ctx context.Context, tokenHash string) (*models.RegistrationIntent, error)
	MarkRegistrationEmailSent(ctx context.Context, intentID uint) error
	MarkRegistrationEmailFailed(ctx context.Context, intentID uint) error
	CompleteRegistration(ctx context.Context, intentID uint, credential *models.Credential) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (repo *authRepository) CreateCredential(ctx context.Context, credential *models.Credential) error {
	return repo.db.WithContext(ctx).Create(credential).Error
}

func (repo *authRepository) FindCredentialByEmail(ctx context.Context, email string) (*models.Credential, error) {
	var credential models.Credential

	err := repo.db.WithContext(ctx).Where("LOWER(email) = LOWER(?)", email).First(&credential).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &credential, nil
}

func (repo *authRepository) FindCredentialByID(ctx context.Context, credentialID uint) (*models.Credential, error) {
	var credential models.Credential

	err := repo.db.WithContext(ctx).First(&credential, "id = ?", credentialID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &credential, nil
}

func (repo *authRepository) UpdateCredential(ctx context.Context, credential *models.Credential) error {
	return repo.db.WithContext(ctx).Model(&models.Credential{}).
		Where("id = ?", credential.ID).
		Select("PasswordHash", "PasswordChangedAt", "UpdatedAt").
		Updates(credential).Error
}

func (repo *authRepository) CreateSession(ctx context.Context, session *models.AuthSession) error {
	return repo.db.WithContext(ctx).Create(session).Error
}

func (repo *authRepository) FindSessionByTokenHash(ctx context.Context, tokenHash string) (*models.AuthSession, error) {
	var session models.AuthSession

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

	result := repo.db.WithContext(ctx).Model(&models.AuthSession{}).
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

	return repo.db.WithContext(ctx).Model(&models.AuthSession{}).
		Where("credential_id = ? AND revoked_at IS NULL", credentialID).
		Update("revoked_at", now).Error
}

func (repo *authRepository) RotateSession(ctx context.Context, oldSessionID uint, newSession *models.AuthSession) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		result := tx.Model(&models.AuthSession{}).
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

func (repo *authRepository) SavePasswordResetToken(ctx context.Context, token *models.PasswordResetToken) error {
	return repo.db.WithContext(ctx).Create(token).Error
}

func (repo *authRepository) FindPasswordResetToken(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error) {
	var token models.PasswordResetToken

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

	result := repo.db.WithContext(ctx).Model(&models.PasswordResetToken{}).
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

func (repo *authRepository) CompletePasswordReset(ctx context.Context, tokenID uint, credential *models.Credential) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		tokenResult := tx.Model(&models.PasswordResetToken{}).
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

		credentialResult := tx.Model(&models.Credential{}).
			Where("id = ?", credential.ID).
			Select("PasswordHash", "PasswordChangedAt", "UpdatedAt").
			Updates(credential)

		if credentialResult.Error != nil {
			return credentialResult.Error
		}
		if credentialResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return tx.Model(&models.AuthSession{}).
			Where("credential_id = ? AND revoked_at IS NULL", credential.ID).
			Update("revoked_at", now).Error
	})
}

func (repo *authRepository) CreateOrUpdateRegistrationIntent(ctx context.Context, intent *models.RegistrationIntent) error {
	return repo.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"first_name",
			"last_name",
			"password_hash",
			"company_name",
			"phone_number",
			"verification_token_hash",
			"expires_at",
			"email_send_status",
			"email_send_attempts",
			"last_email_sent_at",
			"updated_at",
		}),
	}).Create(intent).Error
}

func (repo *authRepository) FindRegistrationIntentByTokenHash(ctx context.Context, tokenHash string) (*models.RegistrationIntent, error) {
	var intent models.RegistrationIntent

	err := repo.db.WithContext(ctx).
		Where("verification_token_hash = ? AND expires_at > ?", tokenHash, time.Now().UTC()).
		First(&intent).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &intent, nil
}

func (repo *authRepository) MarkRegistrationEmailSent(ctx context.Context, intentID uint) error {
	now := time.Now().UTC()

	return repo.db.WithContext(ctx).Model(&models.RegistrationIntent{}).
		Where("id = ?", intentID).
		Updates(map[string]interface{}{
			"email_send_status":   "sent",
			"email_send_attempts": gorm.Expr("email_send_attempts + 1"),
			"last_email_sent_at":  now,
			"updated_at":          now,
		}).Error
}

func (repo *authRepository) MarkRegistrationEmailFailed(ctx context.Context, intentID uint) error {
	now := time.Now().UTC()

	return repo.db.WithContext(ctx).Model(&models.RegistrationIntent{}).
		Where("id = ?", intentID).
		Updates(map[string]interface{}{
			"email_send_status":   "failed",
			"email_send_attempts": gorm.Expr("email_send_attempts + 1"),
			"updated_at":          now,
		}).Error
}

func (repo *authRepository) CompleteRegistration(ctx context.Context, intentID uint, credential *models.Credential) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var intent models.RegistrationIntent

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND expires_at > ?", intentID, time.Now().UTC()).
			First(&intent).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		if err != nil {
			return err
		}

		if intent.Email != credential.Email {
			return errors.New("registration email mismatch")
		}

		if err := tx.Create(credential).Error; err != nil {
			return err
		}

		return tx.Delete(&intent).Error
	})
}

var _ AuthRepository = (*authRepository)(nil)
