package governance

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type UserStore interface {
	// GetUser 根据用户 ID 获取指定的用户。
	GetUser(ctx contractapi.TransactionContextInterface, id string) (User, error)

	// StoreUser 存储用户。
	StoreUser(ctx contractapi.TransactionContextInterface, u User) error

	// UpdateUser 更新用户信息。
	UpdateUser(ctx contractapi.TransactionContextInterface, u User) error
}

type userStore struct {
}

func NewUserStore() UserStore {
	return &userStore{}
}

func (us *userStore) GetUser(ctx contractapi.TransactionContextInterface, id string) (User, error) {
	userBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("user [%s] not found: [%s]", id, err.Error())
	}

	u := &user{}
	if err = u.Deserialize(userBytes); err != nil {
		return nil, fmt.Errorf("user [%s] not found: [%s]", id, err.Error())
	}

	return u, nil
}

func (us *userStore) StoreUser(ctx contractapi.TransactionContextInterface, u User) error {
	userBytes, err := u.Serialize()
	if err != nil {
		return fmt.Errorf("failed storing user [%s]: [%s]", u.GetUserID(), err.Error())
	}

	exists, err := ctx.GetStub().GetState(u.GetUserID())
	if len(exists) != 0 && err == nil {
		return fmt.Errorf("failed storing user [%s], because [%s] is already existed", u.GetUserID(), u.GetUserID())
	}

	if err = ctx.GetStub().PutState(u.GetUserID(), userBytes); err != nil {
		return fmt.Errorf("failed storing user [%s]: [%s]", u.GetUserID(), err.Error())
	}

	return nil
}

func (us *userStore) UpdateUser(ctx contractapi.TransactionContextInterface, u User) error {
	if oldUser, err := us.GetUser(ctx, u.GetUserID()); err != nil || oldUser == nil {
		return fmt.Errorf("cannot update user [%s], because this user may not exist ")
	}

	userBytes, err := u.Serialize()
	if err != nil {
		return fmt.Errorf("failed updating user [%s]: [%s]", u.GetUserID(), err.Error())
	}

	if err = ctx.GetStub().PutState(u.GetUserID(), userBytes); err != nil {
		return fmt.Errorf("failed updating user [%s]: [%s]", u.GetUserID(), err.Error())
	}

	return nil
}
