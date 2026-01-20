package model

type GitlabHosts struct {
	Model
	Host    string `gorm:"not null;uniqueIndex:idx_user_host" json:"gitlab_host"`
	IsHTTPS *bool  `gorm:"default:true" json:"is_https"`
}

func (GitlabHosts) TableName() string {
	return "gitlab_hosts"
}
