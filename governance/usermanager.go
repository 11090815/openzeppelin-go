package governance

import (
	"fmt"
	"sync"

	"github.com/11090815/openzeppelin-go/protos"
	"google.golang.org/protobuf/proto"
)

type User interface {
	// GetUserID 返回用户的唯一可区分身份标识符。
	GetUserID() string

	// CheckAuthorisedAttribute 检查用户是否拥有授权的属性。
	CheckAuthorisedAttribute(attr string) bool

	// DelegateAttributeAuthority 授予属性权力。
	DelegateAttributeAuthority(attr string) error

	// RevokeAttributeAuthority 剥夺属性权力。
	RevokeAttributeAuthority(attr string) error

	// Serialize 序列化成字节。
	Serialize() ([]byte, error)

	// Deserialize 将字节反序列化成结构体。
	Deserialize([]byte) error
}

type user struct {
	*protos.User

	mutex sync.RWMutex
}

func NewUser(userID string) User {
	return &user{
		User: &protos.User{
			Id:    userID,
			Attrs: make(map[string][]byte),
		},
	}
}

func (u *user) GetUserID() string {
	return u.Id
}

func (u *user) CheckAuthorisedAttribute(attr string) bool {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	_, ok := u.Attrs[attr]

	return ok
}

func (u *user) DelegateAttributeAuthority(attr string) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if _, exists := u.Attrs[attr]; exists {
		return fmt.Errorf("user [%s] has been already delegated attribute [%s]", u.Id, attr)
	}

	u.Attrs[attr] = []byte{0x01}

	return nil
}

func (u *user) RevokeAttributeAuthority(attr string) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if _, exists := u.Attrs[attr]; !exists {
		return fmt.Errorf("user [%s] has not been delegated attribute [%s]", u.Id, attr)
	}

	delete(u.Attrs, attr)

	return nil
}

func (u *user) Serialize() ([]byte, error) {
	return proto.Marshal(u.User)
}

func (u *user) Deserialize(raw []byte) error {
	if u.User == nil {
		u.User = new(protos.User)
	}

	return proto.Unmarshal(raw, u.User)
}
