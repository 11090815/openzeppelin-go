package governance_test

import (
	"testing"

	"github.com/11090815/openzeppelin-go/governance"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	user := governance.NewUser("alice")

	serializedUser, err := user.Serialize()
	require.NoError(t, err)
	t.Log(serializedUser)

	newUser := governance.NewUser("bob")
	newUser.Deserialize(serializedUser)

	t.Log(newUser.GetUserID())
}

func TestHello(t *testing.T) {
	t.Log("hello")	
}