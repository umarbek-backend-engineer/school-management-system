package models

import "database/sql"

type Exec struct {
	ID                        int            `json:"id,omitempty" db:"id,omitempty"`
	FirstName                 string         `json:"first_name,omitempty" db:"first_name,omitempty"`
	LastName                  string         `json:"last_name,omitempty" db:"last_name,omitempty"`
	Email                     string         `json:"email,omitempty" db:"email,omitempty"`
	Username                  string         `json:"username,omitempty" db:"username,omitempty"`
	Password                  string         `json:"password,omitempty" db:"password,omitempty"`
	PasswordChangedAt         sql.NullString `json:"password_changed_at,omitempty" db:"password_changed_at,omitempty"`
	UserCreatedAT             sql.NullString `json:"user_created_at,omitempty" db:"user_created_at,omitempty"`
	PasswordResetToken        sql.NullString `json:"password_reset_token,omitempty" db:"password_reset_token,omitempty"`
	PasswordResetTokenExpires sql.NullString `json:"password_reset_token_expires,omitempty" db:"password_reset_token_expires,omitempty"`
	InacvtiveStatus           bool           `json:"inactivestatus,omitempty" db:"inactivestatus,omitempty"`
	Role                      string         `json:"role,omitempty" db:"role,omitempty"`
}

type Exec_Update_password_request struct {
	Current_Password      string `json:"current_password"`
	New_Password          string `json:"new_password"`
	Confurmation_Password string `json:"conform_password"`
}


type Exec_Update_password_response struct {
	Token string `json:"token"`
	Password_Updated bool `json:"password_updated"`
}