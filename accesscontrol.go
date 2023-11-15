package openzeppelingo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/11090815/openzeppelin-go/governance"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type AccessControl struct {
	contractapi.Contract
	store governance.UserStore
}

func (ac *AccessControl) InitLedger(ctx contractapi.TransactionContextInterface) error {
	ac.store = governance.NewUserStore()
	return nil
}

func (ac *AccessControl) RegisterUser(ctx contractapi.TransactionContextInterface, userID string) error {
	if userID == "" {
		return errors.New("failed registering new user, user id should not be empty")
	}

	user := governance.NewUser(userID)

	if err := ac.store.StoreUser(ctx, user); err != nil {
		return fmt.Errorf("failed registering new user [%s]: [%s]", userID, err.Error())
	}

	return nil
}

func (ac *AccessControl) DelegateUserAttributes(ctx contractapi.TransactionContextInterface, userID string, attrs string) error {
	if userID == "" {
		return errors.New("failed delegating attributes to user, user id should not be empty")
	}

	attrArr := make([]string, 0)
	if strings.Contains(attrs, ";") {
		attrArr = strings.Split(attrs, ";")
	} else {
		attrArr = append(attrArr, attrs)
	}

	user, err := ac.store.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed delegating attributes to user [%s]: [%s]", userID, err.Error())
	}

	for _, attr := range attrArr {
		if err = user.DelegateAttributeAuthority(attr); err != nil {
			return fmt.Errorf("failed delegating attribute [%s] to user [%s]: [%s]", attr, userID, err.Error())
		}
	}

	if err = ac.store.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed delegating attributes to user [%s]: [%s]", userID, err.Error())
	}

	return nil
}

func (ac *AccessControl) RevokeUserAttributes(ctx contractapi.TransactionContextInterface, userID string, attrs string) error {
	if userID == "" {
		return errors.New("failed revoking attributes from user, user id should not be empty")
	}

	attrArr := make([]string, 0)
	if strings.Contains(attrs, ";") {
		attrArr = strings.Split(attrs, ";")
	} else {
		attrArr = append(attrArr, attrs)
	}

	user, err := ac.store.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed revoking attributes from user [%s]: [%s]", userID, err.Error())
	}

	for _, attr := range attrArr {
		if err = user.RevokeAttributeAuthority(attr); err != nil {
			return fmt.Errorf("failed revoking attribute [%s] from user [%s]: [%s]", attr, userID, err.Error())
		}
	}

	if err = ac.store.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed delegating attributes to user [%s]: [%s]", userID, err.Error())
	}

	return nil
}

func (ac *AccessControl) Authentication(ctx contractapi.TransactionContextInterface, userID string, attr string) error {
	if userID == "" {
		return errors.New("failed verifying user access rights, user id should not be empty")
	}

	if attr == "" {
		return errors.New("failed verifying user access rights, attribute should not be empty")
	}

	user, err := ac.store.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed verifying user [%s] access rights: [%s]", userID, err.Error())
	}

	if user.CheckAuthorisedAttribute(attr) {
		return nil
	} else {
		return fmt.Errorf("the user [%s] has not been delegated the attribute [%s] and therefore access is denied to him", userID, attr)
	}
}
