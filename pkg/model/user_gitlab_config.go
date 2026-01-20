package model

type UserGitlabConfig struct {
	Model
	UserID       uint   `json:"user_id" gorm:"not null;uniqueIndex:idx_user_gitlab_config_unique"`
	GitlabHostID uint   `json:"gitlab_host_id" gorm:"not null;uniqueIndex:idx_user_gitlab_config_unique"`
	Token        string `json:"token" gorm:"not null"`
	IsValidated  bool   `json:"is_validated" gorm:"default:false"`

	// Relationships
	User       User        `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GitlabHost GitlabHosts `json:"gitlab_host" gorm:"foreignKey:GitlabHostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (UserGitlabConfig) TableName() string {
	return "user_gitlab_configs"
}
