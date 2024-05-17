package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/tjfleming0101/dreampicai/types"
)

func GetAccountbyUserID(userID uuid.UUID) (types.Account, error) {
	var account types.Account
	err := Bun.NewSelect().Model(&account).Where("user_id = ?", userID).Scan(context.Background())
	return account, err
}

func CreateAccount(account *types.Account) error {
	_, err := Bun.NewInsert().Model(account).Exec(context.Background())
	return err
}
