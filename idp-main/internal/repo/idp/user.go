package idprepo

import (
	idpbiz "application/internal/biz/idp"
	"application/internal/datasource"
	idpentities "application/internal/entity/idp"
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type idpRepository struct {
	logger *slog.Logger
	db     *sql.DB
}

func NewUserRepo(logger *slog.Logger, ds *datasource.Datasource) idpbiz.UserRepositoryInterface {
	return &idpRepository{
		logger: logger.With("layer", "IDPRepository"),
		db:     ds.DBpsql,
	}
}
func (r *idpRepository) CreateUser(ctx context.Context, user *idpentities.User) (string, *idpentities.User, error) {
	logger := r.logger.With("method", "CreateUser")
	logger.Debug("repository CreateUser")

	user.UserID = uuid.New().String()

	row := r.db.QueryRowContext(
		ctx,
		`
        INSERT INTO users 
            (msisdn, msisdn_verified, encrypted_password, contractor_name, contractor_code, registration_number,
            contact_person, ceo_name, authorized_signatories, phone_number, email, address, contract_type, 
            bank_account_number, description,role,created_at,updated_at) 
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,$15,$16,$17,$18) 
            RETURNING user_id`,
		user.Msisdn,
		user.MsisdnVerified,
		user.EncryptedPassword,
		user.ContractorName,
		user.ContractorCode,
		user.RegistrationNumber,
		user.ContactPerson,
		user.CeoName,
		user.AuthorizedSignatories,
		user.PhoneNumber,
		user.Email,
		user.Address,
		user.ContractType,
		user.BankAccountNumber,
		user.Description,
		user.Role,
		time.Now().Unix(),
		time.Now().Unix(),
	)

	var id string
	if err := row.Scan(&id); err != nil {
		logger.Error("error creating user", "error", err)
		return "", nil, err
	}

	logger.Debug("repository CreateUser successful", "userID", id)
	return id, user, nil
}

func (r *idpRepository) UpdateUser(ctx context.Context, userID string, updatedUser *idpentities.User) error {
	logger := r.logger.With("method", "UpdateUser")
	logger.Debug("repository UpdateUser", "userID", userID, "updatedUser", updatedUser)

	query := `
		UPDATE users
		SET msisdn = COALESCE($2, msisdn),
		    msisdn_verified = COALESCE($3, msisdn_verified),
		    encrypted_password = COALESCE($4, encrypted_password),
		    contractor_name = COALESCE($5, contractor_name),
		    contractor_code = COALESCE($6, contractor_code),
		    registration_number = COALESCE($7, registration_number),
		    contact_person = COALESCE($8, contact_person),
		    ceo_name = COALESCE($9, ceo_name),
		    authorized_signatories = COALESCE($10, authorized_signatories),
		    phone_number = COALESCE($11, phone_number),
		    email = COALESCE($12, email),
		    address = COALESCE($13, address),
		    contract_type = COALESCE($14, contract_type),
		    bank_account_number = COALESCE($15, bank_account_number),
		    description = COALESCE($16, description),
		    role = COALESCE($16, role),
		    updated_at = COALESCE($17, updated_at)
		WHERE user_id = $1
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		userID,
		updatedUser.Msisdn,
		updatedUser.MsisdnVerified,
		updatedUser.EncryptedPassword,
		updatedUser.ContractorName,
		updatedUser.ContractorCode,
		updatedUser.RegistrationNumber,
		updatedUser.ContactPerson,
		updatedUser.CeoName,
		updatedUser.AuthorizedSignatories,
		updatedUser.PhoneNumber,
		updatedUser.Email,
		updatedUser.Address,
		updatedUser.ContractType,
		updatedUser.BankAccountNumber,
		updatedUser.Description,
		updatedUser.Role,
		time.Now().Unix(),
	)
	if err != nil {
		logger.Error("error updating user", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Warn("error checking rows affected", "error", err)
		return err
	}
	if rowsAffected == 0 {
		logger.Error("no user updated", "userID", userID)
		return idpbiz.ErrorUserNotFount
	}

	logger.Debug("repository UpdateUser successful", "userID", userID)
	return nil
}

func (r *idpRepository) ListUser(ctx context.Context) ([]idpentities.User, error) {
	logger := r.logger.With("method", "ListUsers")
	logger.Debug("repository ListUsers")

	var users []idpentities.User

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT user_id, msisdn, msisdn_verified, encrypted_password, contractor_name, contractor_code, "+
			"registration_number, contact_person, ceo_name, authorized_signatories, phone_number, email, address, "+
			"contract_type, bank_account_number, description,role, created_at, updated_at FROM users",
	)
	if err != nil {
		logger.Error("error listing users", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user idpentities.User
		err = rows.Scan(
			&user.UserID,
			&user.Msisdn,
			&user.MsisdnVerified,
			&user.EncryptedPassword,
			&user.ContractorName,
			&user.ContractorCode,
			&user.RegistrationNumber,
			&user.ContactPerson,
			&user.CeoName,
			&user.AuthorizedSignatories,
			&user.PhoneNumber,
			&user.Email,
			&user.Address,
			&user.ContractType,
			&user.BankAccountNumber,
			&user.Description,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			logger.Error("error scanning user row", "error", err)
			return nil, err
		}
		users = append(users, user)
	}

	logger.Debug("repository ListUsers successful", "users", users)
	return users, nil
}

func (r *idpRepository) GetUserByID(ctx context.Context, userID string) (*idpentities.User, error) {
	logger := r.logger.With("method", "GetUserByID")
	logger.Debug("repository GetUserByID", "userID", userID)

	var user idpentities.User

	row := r.db.QueryRowContext(
		ctx,
		`
		SELECT 
			user_id, msisdn, msisdn_verified, encrypted_password, contractor_name, contractor_code, 
			registration_number, contact_person, ceo_name, authorized_signatories, phone_number, 
			email, address, contract_type, bank_account_number, description,role, created_at, updated_at 
		FROM users 
		WHERE user_id = $1`,
		userID,
	)

	err := row.Scan(
		&user.UserID,
		&user.Msisdn,
		&user.MsisdnVerified,
		&user.EncryptedPassword,
		&user.ContractorName,
		&user.ContractorCode,
		&user.RegistrationNumber,
		&user.ContactPerson,
		&user.CeoName,
		&user.AuthorizedSignatories,
		&user.PhoneNumber,
		&user.Email,
		&user.Address,
		&user.ContractType,
		&user.BankAccountNumber,
		&user.Description,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn("user not found", "userID", userID)
			return nil, idpbiz.ErrorUserNotFount
		}
		logger.Error("error retrieving user", "error", err)
		return nil, err
	}

	logger.Debug("repository GetUserByID successful", "user", user)
	return &user, nil
}

func (r *idpRepository) DeleteUser(ctx context.Context, userID string) error {
	logger := r.logger.With("method", "DeleteUser")
	logger.Debug("repository DeleteUser", "userID", userID)

	result, err := r.db.ExecContext(
		ctx,
		"DELETE FROM users WHERE user_id = $1",
		userID,
	)
	if err != nil {
		logger.Error("error deleting user", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Warn("error checking rows affected", "error", err)
		return err
	}
	if rowsAffected == 0 {
		logger.Error("no user deleted", "userID", userID)
		return idpbiz.ErrorUserNotFount
	}

	logger.Debug("repository DeleteUser successful", "userID", userID)
	return nil
}

func (r *idpRepository) GetUserByMsisdn(ctx context.Context, msisdn string) (*idpentities.User, error) {
	logger := r.logger.With("method", "GetUserByMsisdn")
	logger.Debug("repository GetUserByMsisdn", "msisdn", msisdn)

	var user idpentities.User

	row := r.db.QueryRowContext(
		ctx,
		"SELECT user_id, msisdn, msisdn_verified, encrypted_password, contractor_name, contractor_code,"+
			" registration_number, contact_person, ceo_name, authorized_signatories, phone_number, email, address, "+
			"contract_type, bank_account_number, description,role, created_at, updated_at FROM users WHERE msisdn = $1",
		msisdn,
	)

	err := row.Scan(
		&user.UserID,
		&user.Msisdn,
		&user.MsisdnVerified,
		&user.EncryptedPassword,
		&user.ContractorName,
		&user.ContractorCode,
		&user.RegistrationNumber,
		&user.ContactPerson,
		&user.CeoName,
		&user.AuthorizedSignatories,
		&user.PhoneNumber,
		&user.Email,
		&user.Address,
		&user.ContractType,
		&user.BankAccountNumber,
		&user.Description,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		logger.Error("error retrieving user by msisdn", "error", err)
		return nil, idpbiz.ErrorNotFound
	}

	logger.Debug("repository GetUserByMsisdn successful", "user", user)
	return &user, nil
}
