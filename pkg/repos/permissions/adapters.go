package permissions

import (
	"bytes"
	"errors"
	"strings"

	pgadapter "github.com/casbin/casbin-pg-adapter"
	"github.com/casbin/casbin/util"
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

// StringAdapter implements CasbinAdapter using a string to store policies.
type StringAdapter struct {
	Line string
}

// NewStringAdapter creates a new StringAdapter instance.
func NewStringAdapter(line string) StringAdapter {
	return StringAdapter{
		Line: line,
	}
}

// LoadPolicy loads policy from the adapter's string.
func (a StringAdapter) LoadPolicy(model model.Model) error {
	if a.Line == "" {
		return errors.New("invalid line, line cannot be empty")
	}
	strs := strings.Split(a.Line, "\n")
	for _, str := range strs {
		if str == "" {
			continue
		}
		persist.LoadPolicyLine(str, model)
	}

	return nil
}

// SavePolicy saves policy to the adapter's string.
func (a StringAdapter) SavePolicy(model model.Model) error {
	var tmp bytes.Buffer
	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			tmp.WriteString(ptype + ", ")
			tmp.WriteString(util.ArrayToString(rule))
			tmp.WriteString("\n")
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			tmp.WriteString(ptype + ", ")
			tmp.WriteString(util.ArrayToString(rule))
			tmp.WriteString("\n")
		}
	}
	a.Line = strings.TrimRight(tmp.String(), "\n")
	return nil
}

// AddPolicy adds a policy rule to the storage.
func (a StringAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	// Since the StringAdapter is not designed to dynamically add policies, return an error or implement as needed.
	return errors.New("AddPolicy not implemented")
}

// RemovePolicy removes a policy rule from the storage.
func (a StringAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	// This simplistic approach clears all policies. Consider implementing policy removal logic.
	a.Line = ""
	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a StringAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	// Since the StringAdapter is not designed to dynamically remove filtered policies, return an error or implement as needed.
	return errors.New("RemoveFilteredPolicy not implemented")
}
