package model

import (
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db
	err = DB.AutoMigrate(&GitlabHosts{})
	if err != nil {
		panic(err)
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func TestSeedGitlabHosts(t *testing.T) {
	setupTestDB()

	tests := []struct {
		name          string
		envValue      string
		expectedHosts []GitlabHosts
	}{
		{
			name:     "Single HTTPS host with scheme",
			envValue: "https://gitlab.com",
			expectedHosts: []GitlabHosts{
				{Host: "gitlab.com", IsHTTPS: boolPtr(true)},
			},
		},
		{
			name:     "Single HTTP host with scheme",
			envValue: "http://local.gitlab",
			expectedHosts: []GitlabHosts{
				{Host: "local.gitlab", IsHTTPS: boolPtr(false)},
			},
		},
		{
			name:     "Host without scheme (defaults to HTTPS)",
			envValue: "gitlab.company.com",
			expectedHosts: []GitlabHosts{
				{Host: "gitlab.company.com", IsHTTPS: boolPtr(true)},
			},
		},
		{
			name:     "Multiple hosts mixed",
			envValue: "https://gitlab.com, http://insecure.gitlab, other.gitlab",
			expectedHosts: []GitlabHosts{
				{Host: "gitlab.com", IsHTTPS: boolPtr(true)},
				{Host: "insecure.gitlab", IsHTTPS: boolPtr(false)},
				{Host: "other.gitlab", IsHTTPS: boolPtr(true)},
			},
		},
		{
			name:     "Hosts with paths and trailing slashes",
			envValue: "https://gitlab.com/group/project, http://server:8080/",
			expectedHosts: []GitlabHosts{
				{Host: "gitlab.com", IsHTTPS: boolPtr(true)},
				{Host: "server:8080", IsHTTPS: boolPtr(false)},
			},
		},
		{
			name:          "Empty env var",
			envValue:      "",
			expectedHosts: []GitlabHosts{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear table
			DB.Exec("DELETE FROM gitlab_hosts")

			err := os.Setenv("GITLAB_HOSTS", tt.envValue)
			assert.NoError(t, err)
			seedGitlabHosts()

			var hosts []GitlabHosts
			DB.Find(&hosts)

			assert.Equal(t, len(tt.expectedHosts), len(hosts))

			// loose check for existence
			for _, expected := range tt.expectedHosts {
				found := false
				for _, actual := range hosts {
					if actual.Host == expected.Host {
						if actual.IsHTTPS != nil && expected.IsHTTPS != nil && *actual.IsHTTPS == *expected.IsHTTPS {
							found = true
							break
						}
					}
				}
				assert.True(t, found, "Expected host %s (https=%v) not found", expected.Host, *expected.IsHTTPS)
			}
		})
	}
}
