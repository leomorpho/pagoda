package permissions_test

import (
	"sync"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/davecgh/go-spew/spew"
	"github.com/mikestefanello/pagoda/pkg/repos/permissions"
	"github.com/mikestefanello/pagoda/pkg/repos/tester"
	"github.com/stretchr/testify/assert"
)

// SimpleInMemoryCache provides a basic in-memory cache implementation.
type SimpleInMemoryCache struct {
	mutex sync.RWMutex
	data  map[string]bool
}

// NewSimpleInMemoryCache creates a new instance of SimpleInMemoryCache.
func NewSimpleInMemoryCache() *SimpleInMemoryCache {
	return &SimpleInMemoryCache{
		data: make(map[string]bool),
	}
}

// Put stores a key-value pair in the cache.
func (c *SimpleInMemoryCache) Put(key string, value bool) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
	return nil
}

// Get retrieves a value by key from the cache.
// Returns the value, whether the key was found, and an error if occurred.
func (c *SimpleInMemoryCache) Get(key string) (bool, bool, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, found := c.data[key]
	return value, found, nil
}

// TestCheckPermission updates the test to use mock functions for cache operations.
func TestCheckPermission(t *testing.T) {

	dsn, _ := tester.CreateTestContainerPostgresConnStr(t)
	adapter, err := permissions.NewPostgresCasbinAdapter(dsn)
	assert.NoError(t, err)
	modelText :=
		`
		[request_definition]
		r = sub, dom, obj, act

		[policy_definition]
		p = sub, dom, obj, act

		[role_definition]
		g = _, _, _

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
		`
	m, err := model.NewModelFromString(modelText)

	// Initialize the Casbin enforcer using the adapter
	enforcer, err := casbin.NewEnforcer(m, adapter)

	t.Run("Permission Allowed", func(t *testing.T) {
		// inMemoryCache := NewSimpleInMemoryCache()
		// client, err := permissions.NewPermissionClient(adapter, false, inMemoryCache.Get, inMemoryCache.Put)
		assert.NoError(t, err)
		enforcer.AddPolicy("tenant1", "alice", "data1", "read")
		// added, err := client.AddPermission("tenant1", "alice", "data1", "read")
		// assert.NoError(t, err)
		// assert.True(t, added)
		// assert.Len(t, inMemoryCache.data, 1)
		// // Check that the policy is invalidated in cache since we added a new permission
		// assert.Equal(t, inMemoryCache.data["policy-loaded|tenant1"], false)

		// policy := client.GetPolicies()
		// assert.Len(t, policy, 1)
		// assert.Len(t, policy[0], 4)
		// assert.Equal(t, policy[0], []string{"tenant1", "alice", "data1", "read"})

		allowed, err := enforcer.Enforce("tenant1", "alice", "data1", "read")
		assert.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("Permission Denied", func(t *testing.T) {
		inMemoryCache := NewSimpleInMemoryCache()
		client, err := permissions.NewPermissionClient(adapter, false, inMemoryCache.Get, inMemoryCache.Put)
		assert.NoError(t, err)

		allowed, err := client.CheckPermission("tenant1", "alice", "data2", "write")
		assert.NoError(t, err)
		assert.False(t, allowed)
	})
}
func TestAccessControlDefault(t *testing.T) {

	e, err := casbin.NewEnforcer("test/rbac_model.conf", "test/rbac_policy.csv")
	assert.NoError(t, err)

	// Tests
	tests := []struct {
		sub      string
		obj      string
		act      string
		expected bool
	}{
		{"alice", "data1", "read", true},
		{"alice", "data1", "write", false},
		{"alice", "data2", "read", true},
		{"bob", "data2", "write", true},
		{"bob", "data1", "read", false},
		{"bob", "data2", "read", false},
		{"bob", "data2", "write", true},
	}

	for _, test := range tests {
		t.Run(test.sub, func(t *testing.T) {
			spew.Dump(e.GetPolicy())
			roles, _ := e.GetRolesForUser("alice")

			spew.Dump(roles)
			result, _ := e.Enforce(test.sub, test.obj, test.act)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestAccessControlWithAdapterAndModel(t *testing.T) {
	// Define model text
	modelText := `
	[request_definition]
	r = sub, obj, act

	[policy_definition]
	p = sub, obj, act

	[role_definition]
	g = _, _
	g2 = _, _

	[policy_effect]
	e = some(where (p.eft == allow))

	[matchers]
	m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
	`

	dsn, _ := tester.CreateTestContainerPostgresConnStr(t)
	adapter, _ := permissions.NewPostgresCasbinAdapter(dsn)
	m, _ := model.NewModelFromString(modelText)
	e, err := casbin.NewEnforcer(m, adapter)
	// e, err := casbin.NewEnforcer("test/rbac_model.conf", adapter)
	assert.NoError(t, err)

	// Load policies and roles into Casbin
	e.AddPolicy("alice", "data1", "read")
	e.AddPolicy("bob", "data2", "write")
	e.AddPolicy("data2_admin", "data2", "read")
	e.AddPolicy("data2_admin", "data2", "write")
	e.AddGroupingPolicy("alice", "data2_admin")

	e.SavePolicy()
	e.LoadPolicy()

	// Tests
	tests := []struct {
		sub      string
		obj      string
		act      string
		expected bool
	}{
		{"alice", "data1", "read", true},
		{"alice", "data1", "write", false},
		{"alice", "data2", "read", true},
		{"bob", "data2", "write", true},
		{"bob", "data1", "read", false},
		{"bob", "data2", "read", false},
		{"bob", "data2", "write", true},
	}

	for _, test := range tests {
		t.Run(test.sub, func(t *testing.T) {
			spew.Dump(e.GetPolicy())
			roles, _ := e.GetRolesForUser("alice")

			spew.Dump(roles)
			result, _ := e.Enforce(test.sub, test.obj, test.act)
			assert.Equal(t, test.expected, result)
		})
	}
}
