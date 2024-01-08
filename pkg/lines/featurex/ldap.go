package featurex

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/typesx"
)

type LDAP typesx.LdapConfig

const (
	FEATURE_NAME_LDAP = "ldap"
	BaseDN            = "DC=xxxmotors,DC=com"
	UserDN            = "OU=xxx Users"
	EmployeeDN        = "OU=Employees"
	ContractorsDN     = "OU=Contractors"
	ServiceAccountDN  = "OU=Service Accounts"
	InternsDN         = "OU=Interns"
)

var EntryAvailable = []string{
	EmployeeDN,
	ContractorsDN,
	ServiceAccountDN,
	InternsDN,
}

const (
	DN          = "DN"
	CondEqual   = "="
	CondLess    = "<"
	CondGreater = ">"
	CondIn      = "in"
	WildCard    = "*"
	LogicAnd    = "&"
	LogicOr     = "|"
	LogicNon    = "!"
)

const (
	PropObjectClass      = "objectClass"
	PropMail             = "mail"
	PropCompany          = "company"
	PropManager          = "manager"
	PropMemberOf         = "memberOf"
	PropDepartment       = "department"
	PropDisplayName      = "displayName"
	PropEmployeeNumber   = "employeeNumber"
	PropDepartmentNumber = "departmentNumber"
)

var CommonField = []string{
	PropMail,
	PropCompany,
	PropManager,
	PropMemberOf,
	PropDepartment,
	PropDisplayName,
	PropEmployeeNumber,
	PropDepartmentNumber,
}

type SearchCond struct {
	Fuzzy    bool
	LogicNon bool
	Prop     string
	Cond     string
	Value    string
	Values   []string
}

type LdapCommonUser struct {
	DN               string   `json:"DN"`
	Mail             string   `json:"mail"`
	Company          string   `json:"company"`
	Manager          string   `json:"manager"`
	Department       string   `json:"department"`
	DisplayName      string   `json:"displayName"`
	EmployeeNumber   string   `json:"employeeNumber"`
	DepartmentNumber string   `json:"departmentNumber"`
	MemberOf         []string `json:"memberOf"`
}

func (cfg LDAP) Resolve(appCtx appx.AppContext) (interface{}, error) {
	return NewLDAPService(typesx.LdapConfig(cfg)), nil
}

func (cfg LDAP) FeatureName() string {
	return FEATURE_NAME_LDAP
}

type LDAPService interface {
	SearchEntryByDN(dn string) (map[string][]string, error)
	SearchUserByEmail(email string) (map[string][]string, error)
	SearchUsersByManager(managerEmail string) ([]map[string][]string, error)
	SearchManagerByEmail(email string) (map[string][]string, error)
	SearchUsersByDeptNumber(deptNumber string) ([]map[string][]string, error)
	SearchUserByEmpNumber(empNumber string) (map[string][]string, error)
	SearchUsersByCond(cond []SearchCond) ([]map[string][]string, error)
}

type LDAPServiceImpl struct {
	config typesx.LdapConfig
}

func NewLDAPService(config typesx.LdapConfig) LDAPService {
	s := &LDAPServiceImpl{
		config: config,
	}
	return s
}

func (l *LDAPServiceImpl) getLdapConn() (*ldap.Conn, error) {
	conn, err := ldap.DialURL(l.config.URL)
	if err != nil {
		return nil, errorx.Wrap(err, "can't resolve ldap")
	}

	_, err = conn.SimpleBind(&ldap.SimpleBindRequest{
		Username: l.config.Username,
		Password: l.config.Password,
	})
	if err != nil {
		return nil, errorx.New(fmt.Sprintf(": %s\\n\"", err))
	}

	return conn, errorx.Wrap(err, "Failed to bind")
}

func (l *LDAPServiceImpl) SearchEntryByDN(dn string) (map[string][]string, error) {
	conn, err := l.getLdapConn()
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	defer conn.Close()

	ldapRequest := NewSearchRequestWithScope(dn, ldap.ScopeWholeSubtree, CommonField, []SearchCond{
		{
			Prop:  PropObjectClass,
			Cond:  CondEqual,
			Value: WildCard,
		},
	})
	searchResult, err := conn.Search(ldapRequest)

	if err != nil || len(searchResult.Entries) == 0 {
		return nil, errorx.Wrap(err, "fetch ldap data failed with DN")
	}

	return ConvertResultMap(*searchResult.Entries[0]), nil
}

func (l *LDAPServiceImpl) SearchUserByEmail(email string) (map[string][]string, error) {
	conditions := []SearchCond{
		{
			Prop:  PropMail,
			Cond:  CondEqual,
			Value: email,
		},
	}
	resultMaps, err := l.SearchUsersByCond(conditions)
	if err != nil || len(resultMaps) == 0 {
		return nil, errorx.WithStack(err)
	}

	return resultMaps[0], nil
}

func (l *LDAPServiceImpl) SearchUsersByManager(managerEmail string) ([]map[string][]string, error) {
	managerMap, err := l.SearchUserByEmail(managerEmail)
	if err != nil || managerMap == nil {
		return nil, errorx.WithStack(err)
	}

	conditions := []SearchCond{
		{
			Prop:  PropManager,
			Cond:  CondEqual,
			Value: managerMap[DN][0],
		},
	}
	resultMaps, err := l.SearchUsersByCond(conditions)
	if err != nil || len(resultMaps) == 0 {
		return nil, errorx.WithStack(err)
	}

	return resultMaps, nil
}

func (l *LDAPServiceImpl) SearchManagerByEmail(email string) (map[string][]string, error) {
	userMap, err := l.SearchUserByEmail(email)
	if err != nil || userMap == nil {
		return nil, errorx.WithStack(err)
	}
	managerDN := userMap[PropManager]
	managerUserMap, err := l.SearchEntryByDN(managerDN[0])
	return managerUserMap, errorx.WithStack(err)
}

func (l *LDAPServiceImpl) SearchUsersByDeptNumber(deptNumber string) ([]map[string][]string, error) {
	conditions := []SearchCond{
		{
			Prop:  PropDepartmentNumber,
			Cond:  CondEqual,
			Value: deptNumber,
		},
	}
	resultMaps, err := l.SearchUsersByCond(conditions)
	if err != nil || len(resultMaps) == 0 {
		return nil, errorx.WithStack(err)
	}

	return resultMaps, nil
}

func (l *LDAPServiceImpl) SearchUserByEmpNumber(empNumber string) (map[string][]string, error) {
	conditions := []SearchCond{
		{
			Prop:  PropEmployeeNumber,
			Cond:  CondEqual,
			Value: empNumber,
		},
	}
	resultMaps, err := l.SearchUsersByCond(conditions)
	if err != nil || len(resultMaps) == 0 {
		return nil, errorx.WithStack(err)
	}

	return resultMaps[0], nil
}

func (l *LDAPServiceImpl) SearchUsersByCond(cond []SearchCond) ([]map[string][]string, error) {
	conn, err := l.getLdapConn()
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	defer conn.Close()

	var resultEntries []*ldap.Entry
	var resultMaps []map[string][]string

	for _, dn := range EntryAvailable {
		ldapEmpRequest := NewSearchRequest(getUserBaseDN(dn), CommonField, cond)
		empResult, err := conn.Search(ldapEmpRequest)

		if err != nil {
			return nil, errorx.Wrap(err, "fetch ldap data failed with DN")
		}

		resultEntries = append(resultEntries, empResult.Entries...)
	}

	for _, v := range resultEntries {
		resultMap := ConvertResultMap(*v)
		resultMaps = append(resultMaps, resultMap)
	}

	return resultMaps, nil
}

// helper funcs
func NewSearchRequest(baseDN string, fields []string, queryConditions []SearchCond) *ldap.SearchRequest {
	return NewSearchRequestWithScope(baseDN, ldap.ScopeWholeSubtree, fields, queryConditions)
}

func NewSearchRequestWithScope(baseDN string, scope int, fields []string, queryConditions []SearchCond) *ldap.SearchRequest {
	return ldap.NewSearchRequest(
		baseDN, // The base dn to search
		scope, ldap.DerefAlways, 0, 0, false,
		getSearchCondExpression(queryConditions),
		fields,
		nil,
	)
}

func ConvertResultMap(result ldap.Entry) map[string][]string {
	resultMap := map[string][]string{}

	resultMap[DN] = []string{result.DN}
	for i := range result.Attributes {
		attribute := result.Attributes[i]
		resultMap[attribute.Name] = attribute.Values
	}

	return resultMap
}

func getFirstValue(attributes []string) string {
	if len(attributes) == 0 {
		return ""
	}
	return attributes[0]
}

func ConvertLdapCommonUser(resultMap map[string][]string) LdapCommonUser {
	user := LdapCommonUser{
		DN:               getFirstValue(resultMap[DN]),
		Mail:             getFirstValue(resultMap[PropMail]),
		Company:          getFirstValue(resultMap[PropCompany]),
		Manager:          getFirstValue(resultMap[PropManager]),
		Department:       getFirstValue(resultMap[PropDepartment]),
		DisplayName:      getFirstValue(resultMap[PropDisplayName]),
		EmployeeNumber:   getFirstValue(resultMap[PropEmployeeNumber]),
		DepartmentNumber: getFirstValue(resultMap[PropDepartmentNumber]),
		MemberOf:         resultMap[PropMemberOf],
	}

	return user
}

func getSearchCondExpression(cond []SearchCond) string {
	tmpl := "(&%s)"
	condStr := ""

	var condExpressions []string
	for _, v := range cond {
		condExpressions = append(condExpressions, getSingleSearchCondExpression(v))
	}

	if len(condExpressions) != 0 {
		condStr = strings.Join(condExpressions, "")
	}

	return fmt.Sprintf(tmpl, condStr)
}

func getSingleSearchCondExpression(cond SearchCond) string {
	tmpl := "(%s%s=%s)"

	if strings.EqualFold(cond.Cond, CondIn) {
		var contStrs []string
		for _, v := range cond.Values {
			condStr := fmt.Sprintf(tmpl, "", cond.Prop, v)
			contStrs = append(contStrs, condStr)
		}

		return fmt.Sprintf("(|%s)", strings.Join(contStrs, ""))
	}

	if cond.Fuzzy {
		cond.Value = WildCard + cond.Value + WildCard
	}

	if cond.LogicNon {
		return fmt.Sprintf(tmpl, LogicNon, cond.Prop, cond.Value)
	} else {
		return fmt.Sprintf(tmpl, "", cond.Prop, cond.Value)
	}
}

func getUserBaseDN(subDN string) string {
	dnSlice := []string{
		subDN,
		UserDN,
		BaseDN,
	}
	return strings.Join(dnSlice, ",")
}
