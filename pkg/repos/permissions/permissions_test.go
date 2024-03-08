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

// TestAccessControlWithAdapterAndModel demonstrates how to use Casbin directly for RBAC
func TestAccessControlWithAdapterAndModel(t *testing.T) {
	// Define model text
	modelText := `
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

	dsn, _ := tester.CreateTestContainerPostgresConnStr(t)
	adapter, _ := permissions.NewPostgresCasbinAdapter(dsn)
	m, _ := model.NewModelFromString(modelText)
	e, err := casbin.NewEnforcer(m, adapter)
	assert.NoError(t, err)

	// Load policies and roles into Casbin
	e.AddPolicy("alice", "tenant1", "data1", "read")
	e.AddPolicy("bob", "tenant1", "data2", "write")
	e.AddPolicy("data2_admin", "tenant1", "data2", "read")
	e.AddPolicy("data2_admin", "tenant1", "data2", "write")
	e.AddGroupingPolicy("alice", "data2_admin", "tenant1")

	// Tests
	tests := []struct {
		sub      string
		tenant   string
		obj      string
		act      string
		expected bool
	}{
		{"alice", "tenant1", "data1", "read", true},
		{"alice", "tenant1", "data1", "write", false},
		{"alice", "tenant1", "data2", "read", true},
		{"bob", "tenant1", "data2", "write", true},
		{"bob", "tenant1", "data1", "read", false},
		{"bob", "tenant1", "data2", "read", false},
		{"bob", "tenant1", "data2", "write", true},
	}

	for _, test := range tests {
		t.Run(test.sub, func(t *testing.T) {
			spew.Dump(e.GetPolicy())
			roles, _ := e.GetRolesForUser("alice")

			spew.Dump(roles)
			result, _ := e.Enforce(test.sub, test.tenant, test.obj, test.act)
			assert.Equal(t, test.expected, result)
		})
	}
}

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

func TestAccessControlWithPermissionClient(t *testing.T) {
	// Define model text
	modelText := `
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

	dsn, _ := tester.CreateTestContainerPostgresConnStr(t)
	adapter, _ := permissions.NewPostgresCasbinAdapter(dsn)

	// Tests
	tests := []struct {
		sub      string
		tenant   string
		obj      string
		act      string
		expected bool
	}{
		{"alice", "tenant1", "data1", "read", true},
		{"alice", "tenant1", "data1", "write", false},
		{"alice", "tenant1", "data2", "read", true},
		{"bob", "tenant1", "data2", "write", true},
		{"bob", "tenant1", "data1", "read", false},
		{"bob", "tenant1", "data2", "read", false},
		{"bob", "tenant1", "data2", "write", true},
	}

	for _, test := range tests {
		t.Run(test.sub, func(t *testing.T) {
			inMemoryCache := NewSimpleInMemoryCache()
			client, err := permissions.NewPermissionClient(modelText, adapter, false, inMemoryCache.Get, inMemoryCache.Put)
			client.AddPolicy("tenant1", "alice", "data1", "read")
			client.AddPolicy("tenant1", "bob", "data2", "write")
			client.AddPolicy("tenant1", "data2_admin", "data2", "read")
			client.AddPolicy("tenant1", "data2_admin", "data2", "write")
			client.AddGroupingPolicy("tenant1", "alice", "data2_admin")

			assert.NoError(t, err)

			result, _ := client.CheckPermission(test.tenant, test.sub, test.obj, test.act)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestAccessControlCaching(t *testing.T) {
	// Define model text
	modelText := `
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

	dsn, _ := tester.CreateTestContainerPostgresConnStr(t)
	adapter, _ := permissions.NewPostgresCasbinAdapter(dsn)

	t.Run("Remove and add permission", func(t *testing.T) {
		inMemoryCache := NewSimpleInMemoryCache()
		client, err := permissions.NewPermissionClient(modelText, adapter, false, inMemoryCache.Get, inMemoryCache.Put)
		client.AddPolicy("tenant1", "alice", "data1", "read")
		client.AddPolicy("tenant1", "bob", "data2", "write")
		client.AddPolicy("tenant1", "data2_admin", "data2", "read")
		client.AddPolicy("tenant1", "data2_admin", "data2", "write")
		client.AddGroupingPolicy("tenant1", "alice", "data2_admin")

		assert.NoError(t, err)

		result, err := client.CheckPermission("tenant1", "alice", "data1", "read")
		assert.NoError(t, err)
		assert.Equal(t, true, result)

		removed, err := client.RemovePolicy("tenant1", "alice", "data1", "read")
		assert.NoError(t, err)
		assert.Equal(t, true, removed)

		result, err = client.CheckPermission("tenant1", "alice", "data1", "read")
		assert.NoError(t, err)
		assert.Equal(t, false, result)

		added, err := client.AddPolicy("tenant1", "alice", "data1", "read")
		assert.NoError(t, err)
		assert.Equal(t, true, added)

		result, err = client.CheckPermission("tenant1", "alice", "data1", "read")
		assert.NoError(t, err)
		assert.Equal(t, true, result)
	})
}
