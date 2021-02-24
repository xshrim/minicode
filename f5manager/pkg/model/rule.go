package model

import "gorm.io/gorm"

// F5 Rule

type Value struct {
	Model     `json:"-"`
	Op        string `json:"op"`
	Val       string `json:"val"`
	MatcherID uint   `json:"-"`
}

type Matcher struct {
	Model  `json:"-"`
	Kind   string  `json:"kind"`
	Key    string  `json:"key"`
	Values []Value `json:"values" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RuleID uint    `json:"-"`
}

type RuleService struct {
	Model  `json:"-"`
	Svc    string `json:"svc"`
	Port   int    `json:"port"`
	Weight int    `json:"weight"`
	RuleID uint   `json:"-"`
}

type Rule struct {
	Model           `json:"-"`
	Name            string        `json:"name" gorm:"not null;uniqueIndex"`
	LoadbalanceName string        `json:"loadbalanceName"`
	ListenerName    string        `json:"listenerName"`
	Services        []RuleService `json:"services" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Matchers        []Matcher     `json:"matchers" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func InitRule(db *gorm.DB) error {
	return db.AutoMigrate(&Rule{}, &RuleService{}, &Matcher{}, &Value{})
}
