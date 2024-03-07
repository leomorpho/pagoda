package permissions

import (
	"fmt"

	pgadapter "github.com/casbin/casbin-pg-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

// PermissionClient struct to interact with Casbin
type PermissionClient struct {
	enforcer *casbin.Enforcer
	putCache func(key string, value bool) error
	getCache func(key string) (bool, bool, error) // Returns value, found, error
}

func CacheKey(tenantID, sub, obj, act string) string {
	return fmt.Sprintf("perm|%s|%s|%s", sub, obj, act)
}

// NewPermissionClient creates a new PermissionClient
func NewPermissionClient(
	adapter CasbinAdapter,
	getCache func(key string) (bool, bool, error),
	putCache func(key string, value bool) error,
) (*PermissionClient, error) {

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
		m = g(r.sub, r.dom, p.sub) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
		`
	m, err := model.NewModelFromString(modelText)
	if err != nil {
		return nil, fmt.Errorf("failed to create model from string: %w", err)
	}

	// Initialize the Casbin enforcer using the adapter
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Casbin enforcer: %w", err)
	}

	return &PermissionClient{
		enforcer: enforcer,
		putCache: putCache,
		getCache: getCache,
	}, nil
}

// LoadTenantPolicies loads policies for a specific tenant.
func (s *PermissionClient) LoadTenantPolicies(tenantID string) error {
	filter := &pgadapter.Filter{
		P: []string{tenantID}, // Assuming the tenant_id is the first field in the policy
		G: []string{tenantID}, // Same assumption for the grouping policies
	}

	return s.enforcer.LoadFilteredPolicy(filter)
}

// EnsureTenantPolicyLoaded ensures the policy for a given tenant is loaded.
func (s *PermissionClient) EnsureTenantPolicyLoaded(tenantID string) error {
	cacheKey := fmt.Sprintf("policy-loaded|%s", tenantID)

	// Check if policy is already loaded
	_, found, err := s.getCache(cacheKey)
	if err != nil {
		return err // Handle error from cache check
	}
	if found {
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
	if err := s.EnsureTenantPolicyLoaded(tenantID); err != nil {
		return false, err
	}

	cacheKey := fmt.Sprintf("%s|%s|%s|%s", tenantID, sub, obj, act)
	// Attempt to retrieve from cache
	if value, found, err := s.getCache(cacheKey); err == nil && found {
		return value, nil
	} else if err != nil {
		return false, err
	}

	// Proceed with enforcing and update the cache
	allowed, err := s.enforcer.Enforce(tenantID, sub, obj, act)
	if err == nil {
		// Update cache
		s.putCache(cacheKey, allowed)
	}
	return allowed, err
}

// AddPermission adds a new permission to the policy
func (s *PermissionClient) AddPermission(tenantID, sub, obj, act string) (bool, error) {
	if err := s.EnsureTenantPolicyLoaded(tenantID); err != nil {
		return false, err
	}

	// Attempt to add the policy
	added, err := s.enforcer.AddPolicy(tenantID, sub, obj, act)
	if err != nil {
		return false, err
	}
	if added {
		// Invalidate cache because the policy changed
		if err := s.InvalidateTenantPolicyCache(tenantID); err != nil {
			return false, err
		}
	}
	return added, nil
}

// RemovePermission removes a permission from the policy
func (s *PermissionClient) RemovePermission(tenantID, sub, obj, act string) (bool, error) {
	if err := s.EnsureTenantPolicyLoaded(tenantID); err != nil {
		return false, err
	}

	// Attempt to remove the policy
	removed, err := s.enforcer.RemovePolicy(tenantID, sub, obj, act)
	if err != nil {
		return false, err
	}
	if removed {
		// Invalidate cache because the policy changed
		if err := s.InvalidateTenantPolicyCache(tenantID); err != nil {
			return false, err
		}
	}
	return removed, nil
}

// InvalidateTenantPolicyCache invalidates the cache for a tenant's policies.
func (s *PermissionClient) InvalidateTenantPolicyCache(tenantID string) error {
	cacheKey := fmt.Sprintf("policy-loaded|%s", tenantID)
	return s.putCache(cacheKey, false) // Mark policy as not loaded
}
