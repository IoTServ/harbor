// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package db

import (
	"log"
	"os"
	"testing"

	"github.com/goharbor/harbor/src/common"
	"github.com/goharbor/harbor/src/common/dao"
	"github.com/goharbor/harbor/src/common/utils/test"

	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/common/utils/ldap"
	"github.com/goharbor/harbor/src/ui/auth"
	uiConfig "github.com/goharbor/harbor/src/ui/config"
)

var adminServerTestConfig = map[string]interface{}{
	common.ExtEndpoint:        "host01.com",
	common.AUTHMode:           "db_auth",
	common.DatabaseType:       "postgresql",
	common.PostGreSQLHOST:     "127.0.0.1",
	common.PostGreSQLPort:     5432,
	common.PostGreSQLUsername: "postgres",
	common.PostGreSQLPassword: "root123",
	common.PostGreSQLDatabase: "registry",
	// config.SelfRegistration: true,
	common.LDAPURL:       "ldap://127.0.0.1",
	common.LDAPSearchDN:  "cn=admin,dc=example,dc=com",
	common.LDAPSearchPwd: "admin",
	common.LDAPBaseDN:    "dc=example,dc=com",
	common.LDAPUID:       "uid",
	common.LDAPFilter:    "",
	common.LDAPScope:     3,
	common.LDAPTimeout:   30,
	//	config.TokenServiceURL:            "",
	//	config.RegistryURL:                "",
	//	config.EmailHost:                  "",
	//	config.EmailPort:                  25,
	//	config.EmailUsername:              "",
	//	config.EmailPassword:              "password",
	//	config.EmailFrom:                  "from",
	//	config.EmailSSL:                   true,
	//	config.EmailIdentity:              "",
	//	config.ProjectCreationRestriction: config.ProCrtRestrAdmOnly,
	//	config.VerifyRemoteCert:           false,
	//	config.MaxJobWorkers:              3,
	//	config.TokenExpiration:            30,
	common.CfgExpiration: 5,
	//	config.JobLogDir:                  "/var/log/jobs",
	common.AdminInitialPassword: "password",
}

func TestMain(m *testing.M) {
	server, err := test.NewAdminserver(adminServerTestConfig)
	if err != nil {
		log.Fatalf("failed to create a mock admin server: %v", err)
	}
	defer server.Close()

	if err := os.Setenv("ADMINSERVER_URL", server.URL); err != nil {
		log.Fatalf("failed to set env %s: %v", "ADMINSERVER_URL", err)
	}

	secretKeyPath := "/tmp/secretkey"
	_, err = test.GenerateKey(secretKeyPath)
	if err != nil {
		log.Fatalf("failed to generate secret key: %v", err)
		return
	}
	defer os.Remove(secretKeyPath)

	if err := os.Setenv("KEY_PATH", secretKeyPath); err != nil {
		log.Fatalf("failed to set env %s: %v", "KEY_PATH", err)
	}

	if err := uiConfig.Init(); err != nil {
		log.Fatalf("failed to initialize configurations: %v", err)
	}

	database, err := uiConfig.Database()
	if err != nil {
		log.Fatalf("failed to get database configuration: %v", err)
	}

	if err := dao.InitDatabase(database); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
}

func TestSearchUser(t *testing.T) {
	// insert user first
	user := &models.User{
		Username: "existuser",
		Email:    "existuser@placeholder.com",
		Realname: "Existing user",
	}

	err := dao.OnBoardUser(user)
	if err != nil {
		t.Fatalf("Failed to OnBoardUser %v", user)
	}

	var auth *Auth
	newUser, err := auth.SearchUser("existuser")
	if err != nil {
		t.Fatalf("Failed to search user, error %v", err)
	}
	if newUser == nil {
		t.Fatalf("Failed to search user %v", newUser)
	}

}

func TestAuthenticateHelperOnBoardUser(t *testing.T) {
	user := models.User{
		Username: "test01",
		Realname: "test01",
		Email:    "test01@example.com",
	}

	err := auth.OnBoardUser(&user)
	if err != nil {
		t.Errorf("Failed to onboard user error: %v", err)
	}

}

func TestAuthenticateHelperSearchUser(t *testing.T) {

	user, err := auth.SearchUser("admin")
	if err != nil {
		t.Error("Failed to search user, admin")
	}

	if user == nil {
		t.Error("Failed to search user admin")
	}
}

func TestLdapConnectionTest(t *testing.T) {
	var ldapConfig = models.LdapConf{
		LdapURL:               "ldap://127.0.0.1",
		LdapSearchDn:          "cn=admin,dc=example,dc=com",
		LdapSearchPassword:    "admin",
		LdapBaseDn:            "dc=example,dc=com",
		LdapFilter:            "",
		LdapUID:               "cn",
		LdapScope:             3,
		LdapConnectionTimeout: 10,
		LdapVerifyCert:        false,
	}
	// Test ldap connection under auth_mod is db_auth
	err := ldap.ConnectionTestWithConfig(ldapConfig)
	if err != nil {
		t.Fatalf("Failed to test ldap server! error %v", err)
	}
}
