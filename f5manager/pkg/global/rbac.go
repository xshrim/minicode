package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var rbac_policy = `
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && keyMatch(r.dom, p.dom) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
`

func NewEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	m, err := model.NewModelFromString(rbac_policy)
	if err != nil {
		return nil, err
	}

	a, err := gormadapter.NewAdapterByDB(db) // Your driver and data source.
	if err != nil {
		return nil, err
	}

	return casbin.NewEnforcer(m, a)
}
