package permissions_test

import (
	"testing"

	"github.com/casbin/casbin/v2/model"
	"github.com/mikestefanello/pagoda/pkg/repos/permissions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func (mc *MockCache) Put(key string, value bool) error {
	args := mc.Called(key, value)
	return args.Error(0)
}

func (mc *MockCache) Get(key string) (bool, bool, error) {
	args := mc.Called(key)
	return args.Bool(0), args.Bool(1), args.Error(2)
}

// TestCheckPermission updates the test to use mock functions for cache operations.
func TestCheckPermission(t *testing.T) {
	mockCache := &MockCache{}
	mockCache.On("Get", mock.Anything).Return(false, false, nil)
	mockCache.On("Put", mock.Anything, mock.Anything).Return(nil)

	// Initialize the model from code
	m := model.NewModel()
	m.LoadModelFromText(`
	[request_definition]
	r = sub, obj, act

	[policy_definition]
	p = sub, obj, act

	[role_definition]
	g = _, _

	[policy_effect]
	e = some(where (p.eft == allow))

	[matchers]
	m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
	`)
	policies := ""
	sa := permissions.NewStringAdapter(policies)
	client, err := permissions.NewPermissionClient(sa, mockCache.Get, mockCache.Put)
	assert.NoError(t, err)

	t.Run("Permission Allowed", func(t *testing.T) {
		client.AddPermission("tenant1", "alice", "data1", "read")
		allowed, err := client.CheckPermission("tenant1", "alice", "data1", "read")
		assert.NoError(t, err)
		assert.True(t, allowed)

		mockCache.AssertCalled(t, "Get", "perm|tenant1|alice|data1|read")
		mockCache.AssertCalled(t, "Put", "perm|tenant1|alice|data1|read", true)
	})

	t.Run("Permission Denied", func(t *testing.T) {
		allowed, err := client.CheckPermission("tenant1", "alice", "data2", "write")
		assert.NoError(t, err)
		assert.False(t, allowed)

		mockCache.AssertCalled(t, "Get", "perm|tenant1|alice|data2|write")
		mockCache.AssertCalled(t, "Put", "perm|tenant1|alice|data2|write", false)
	})
}
