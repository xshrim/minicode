package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//	DeletedAt DeletedAt `gorm:"index"`
}

// F5 Provider

type LTM struct {
	Model      `json:"-"`
	Url        string `json:"url"`
	Token      string `json:"token"`
	CIDR       string `json:"cidr"`
	ProviderID uint   `json:"-"`
}

type VE struct {
	Model      `json:"-"`
	Addr       string `json:"addr"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	CIDR       string `json:"cidr"`
	ProviderID uint   `json:"-"`
}

type IPPair struct {
	Model            `json:"-"`
	LtmIP            string `json:"ltmip"`
	VeIP             string `json:"veip"`
	Status           string `json:"status"`
	ProviderStatusID uint   `json:"-"`
}

type ProviderStatus struct {
	Model      `json:"-"`
	IPPairs    []IPPair `json:"ippairs" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProviderID uint     `json:"-"`
}

type Provider struct {
	Model    `json:"-"`
	Name     string         `json:"name" gorm:"not null;uniqueIndex;"`
	Ltm      LTM            `json:"ltm" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Ve       VE             `json:"ve" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Cluster  string         `json:"cluster"`
	Selector string         `json:"selector"`
	Status   ProviderStatus `json:"status" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func InitProvider(db *gorm.DB) error {
	return db.AutoMigrate(&Provider{}, &ProviderStatus{}, &LTM{}, &VE{}, &IPPair{})
}
