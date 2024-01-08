package featurex

import (
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.startlite.cn/itapp/startlite/pkg/lines/typesx"
)

/*
	Fill in your own ldap info and params for each test case
*/

func TestSearchEntryByDN(t *testing.T) {
	err := flag.Set("fx", "_test/app.local.yaml")
	assert.NoError(t, err)

	var cfg typesx.LdapConfig

	cl := NewConfigLoader()
	cl.Load(&cfg)
	ldapService := NewLDAPService(cfg)
	res, err := ldapService.SearchEntryByDN("dummy")
	assert.NoError(t, err)
	entryRes := ConvertLdapCommonUser(res)
	fmt.Printf("%v", entryRes)
}

func TestSearchUserByEmail(t *testing.T) {
	err := flag.Set("fx", "_test/app.local.yaml")
	assert.NoError(t, err)

	var cfg typesx.LdapConfig

	cl := NewConfigLoader()
	cl.Load(&cfg)
	ldapService := NewLDAPService(cfg)
	res, err := ldapService.SearchUserByEmail("dummy")
	assert.NoError(t, err)
	userRes := ConvertLdapCommonUser(res)
	fmt.Printf("%v", userRes)
}

func TestSearchUsersByManager(t *testing.T) {
	err := flag.Set("fx", "_test/app.local.yaml")
	assert.NoError(t, err)

	var cfg typesx.LdapConfig

	cl := NewConfigLoader()
	cl.Load(&cfg)
	ldapService := NewLDAPService(cfg)
	res, err := ldapService.SearchUsersByManager("dummy")
	assert.NoError(t, err)
	for _, user := range res {
		userRes := ConvertLdapCommonUser(user)
		fmt.Printf("%v", userRes)
	}

}

func TestSearchManagerByEmail(t *testing.T) {
	err := flag.Set("fx", "_test/app.local.yaml")
	assert.NoError(t, err)

	var cfg typesx.LdapConfig

	cl := NewConfigLoader()
	cl.Load(&cfg)
	ldapService := NewLDAPService(cfg)
	res, err := ldapService.SearchManagerByEmail("dummy")
	assert.NoError(t, err)
	userRes := ConvertLdapCommonUser(res)
	fmt.Printf("%v", userRes)
}

func TestSearchUsersByDeptNumber(t *testing.T) {
	err := flag.Set("fx", "_test/app.local.yaml")
	assert.NoError(t, err)

	var cfg typesx.LdapConfig

	cl := NewConfigLoader()
	cl.Load(&cfg)
	ldapService := NewLDAPService(cfg)

	res, err := ldapService.SearchUsersByDeptNumber("dummy")
	assert.NoError(t, err)
	for _, user := range res {
		userRes := ConvertLdapCommonUser(user)
		fmt.Printf("%v", userRes)
	}
}

func TestSearchUserByEmpNumber(t *testing.T) {
	err := flag.Set("fx", "_test/app.local.yaml")
	assert.NoError(t, err)

	var cfg typesx.LdapConfig

	cl := NewConfigLoader()
	cl.Load(&cfg)
	ldapService := NewLDAPService(cfg)
	res, err := ldapService.SearchUserByEmpNumber("dummy")
	assert.NoError(t, err)
	userRes := ConvertLdapCommonUser(res)
	fmt.Printf("%v", userRes)
}

func TestSearchUsersByCond(t *testing.T) {
	err := flag.Set("fx", "_test/app.local.yaml")
	assert.NoError(t, err)

	var cfg typesx.LdapConfig

	cl := NewConfigLoader()
	cl.Load(&cfg)
	ldapService := NewLDAPService(cfg)
	res, err := ldapService.SearchUsersByCond([]SearchCond{
		{
			Prop:  PropMail,
			Cond:  CondEqual,
			Value: "dummy",
		},
	})

	assert.NoError(t, err)
	for _, user := range res {
		userRes := ConvertLdapCommonUser(user)
		fmt.Printf("%v", userRes)
	}
}
