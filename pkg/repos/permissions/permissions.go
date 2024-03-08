package permissions

import (
	"fmt"

	pgadapter "github.com/casbin/casbin-pg-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

// CasbinAdapter defines the interface for loading, saving, and managing policies
type CasbinAdapter interface {
	persist.Adapter
}

func NewPostgresCasbinAdapter(dsn string) (CasbinAdapter, error) {
	return pgadapter.NewAdapter(dsn)
}

// PermissionClient struct to interact with Casbin
type PermissionClient struct {
	enforcer            *casbin.Enforcer
	useFilteredPolicies bool
	putCache            func(key string, value bool) error
	getCache            func(key string) (bool, bool, error) // Returns value, found, error
}

// CacheKey constructs a unique key for caching, based on policy components
func CacheKey(tenantID, sub, obj, act string) string {
	return fmt.Sprintf("perm|%s|%s|%s|%s", tenantID, sub, obj, act)
}

// NewPermissionClient creates a new PermissionClient
func NewPermissionClient(
	modelStr string,
	adapter CasbinAdapter,
	useFilteredPolicies bool,
	getCache func(key string) (bool, bool, error),
	putCache func(key string, value bool) error,
) (*PermissionClient, error) {

	m, err := model.NewModelFromString(modelStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create model from string: %w", err)
	}

	// Initialize the Casbin enforcer using the adapter
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Casbin enforcer: %w", err)
	}

	return &PermissionClient{
		enforcer:            enforcer,
		useFilteredPolicies: useFilteredPolicies,
		putCache:            putCache,
		getCache:            getCache,
	}, nil
}

// LoadTenantPolicies loads policies for a specific tenant.
func (s *PermissionClient) LoadTenantPolicies(tenantID string) error {
	if s.useFilteredPolicies {
		filter := &pgadapter.Filter{
			P: []string{tenantID}, // Assuming the tenant_id is the first field in the policy
			G: []string{tenantID}, // Same assumption for the grouping policies
		}
		return s.enforcer.LoadFilteredPolicy(filter)
	}

	// For adapters not supporting filtered policies or when filtered policies are not required
	return s.enforcer.LoadPolicy()
}

// EnsureTenantPolicyLoaded ensures the policy for a given tenant is loaded.
func (s *PermissionClient) EnsureTenantPolicyLoaded(tenantID string) error {
	cacheKey := fmt.Sprintf("policy-loaded|%s", tenantID)

	// Check if policy is already loaded
	valid, found, err := s.getCache(cacheKey)
	if err != nil {
		return err // Handle error from cache check
	}
	if found && valid {
		return nil // Policy already loaded
	}

	// Load policy for the tenant
	if err := s.LoadTenantPolicies(tenantID); err != nil {
		return err // Handle error from policy loading
	}

	// Mark policy as loaded in cache
	return s.putCache(cacheKey, true)
}

// CheckPermission checks if a user has permission to perform an action on an object
func (s *PermissionClient) CheckPermission(tenantID, sub, obj, act string) (bool, error) {
	cacheKey := CacheKey(tenantID, sub, obj, act)
	value, found, _ := s.getCache(cacheKey)
	if found {
		return value, nil
	}

	if err := s.EnsureTenantPolicyLoaded(tenantID); err != nil {
		return false, err
	}

	// Proceed with enforcing and update the cache
	allowed, err := s.enforcer.Enforce(sub, tenantID, obj, act)
	if err == nil {
		// Update cache
		s.putCache(cacheKey, allowed)
	}
	return allowed, err
}

// AddPermission adds a new permission to the policy
func (s *PermissionClient) AddPolicy(tenantID, sub, obj, act string) (bool, error) {
	if err := s.EnsureTenantPolicyLoaded(tenantID); err != nil {
		return false, err
	}

	// Attempt to add the policy
	added, err := s.enforcer.AddPolicy(sub, tenantID, obj, act)
	if err != nil {
		return false, err
	}
	if added {
		s.enforcer.SavePolicy()
		err = s.updateCacheForPermission(tenantID, sub, obj, act)
	}
	return added, err
}

func (s *PermissionClient) AddGroupingPolicy(tenantID, sub, group string) (bool, error) {
	if err := s.EnsureTenantPolicyLoaded(tenantID); err != nil {
		return false, err
	}

	// TODO: I suspect there may be stale permissions until the cache
	// gets dropped. Not sure how to cleanly invalidate cached permissions for a specific
	// tenant or subject when adding roles like that. It might not be worth it,
	// and waiting for the cache to become stale might be easiest (by keeping it
	// to 20-30min or so).
	added, err := s.enforcer.AddGroupingPolicy(sub, group, tenantID)
	if err != nil {
		return false, err
	}
	if added {
		s.enforcer.SavePolicy()
	}
	return added, nil

}

// RemovePolicy removes a permission from the policy
func (s *PermissionClient) RemovePolicy(tenantID, sub, obj, act string) (bool, error) {
	if err := s.EnsureTenantPolicyLoaded(tenantID); err != nil {
		return false, err
	}

	// Attempt to remove the policy
	removed, err := s.enforcer.RemovePolicy(sub, tenantID, obj, act)
	if err != nil {
		return false, err
	}
	if removed {
		s.enforcer.SavePolicy()
		err = s.updateCacheForPermission(tenantID, sub, obj, act)
	}
	return removed, err
}

func (s *PermissionClient) updateCacheForPermission(tenantID, sub, obj, act string) error {
	cacheKey := CacheKey(tenantID, sub, obj, act)
	// Do not read from cache but directly form the casbin model, since we
	// are updating the cache based on its results.
	allowed, err := s.enforcer.Enforce(sub, tenantID, obj, act)
	if err != nil {
		return err
	}
	s.putCache(cacheKey, allowed)
	return nil
}

func (s *PermissionClient) GetPolicies() [][]string {
	return s.enforcer.GetPolicy()
}
