package model

import "gorm.io/gorm"

// F5 Loadbalance

type Service struct {
	Model     `json:"-"`
	Svc       string `json:"svc"`
	Port      int    `json:"port"`
	Weight    int    `json:"weight"`
	ListnerID uint   `json:"-"`
}

type Listner struct {
	Model         `json:"-"`
	Name          string    `json:"name" gorm:"not null;uniqueIndex"`
	Port          int       `json:"port"`
	Protocal      string    `json:"protocal"`
	Session       string    `json:"session"`
	Services      []Service `json:"services" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LoadbalanceID uint      `json:"-"`
}

type LoadbalanceStatus struct {
	Model         `json:"-"`
	Process       string `json:"process"`
	LoadbalanceID uint   `json:"-"`
}

type Loadbalance struct {
	Model        `json:"-"`
	Name         string            `json:"name" gorm:"not null;uniqueIndex"`
	ProviderName string            `json:"providerName"`
	Namespace    string            `json:"namespace"`
	LtmIP        string            `json:"ltmip"`
	VeIP         string            `json:"veip"`
	Kind         string            `json:"kind"`
	Listners     []Listner         `json:"listners" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status       LoadbalanceStatus `json:"status" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func InitLoadbalance(db *gorm.DB) error {
	return db.AutoMigrate(&Loadbalance{}, &LoadbalanceStatus{}, &Listner{}, &Service{})
}
