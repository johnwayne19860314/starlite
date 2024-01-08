package claimx

import (
	"github.com/dgrijalva/jwt-go/v4"

	"github.startlite.cn/itapp/startlite/pkg/lines/typesx"
)

func GetClaim(userClaims *jwt.MapClaims, key string) *string {
	if value, ok := (*userClaims)[key]; ok {
		if s, ok := value.(string); ok {
			return &s
		}
	}
	return typesx.StringP("")
}

func GetUserId(userClaims *jwt.MapClaims) string {
	return *GetClaim(userClaims, "preferred_username")
}

func HasRole(userClaims *jwt.MapClaims, role string) bool {
	roles := GetRoles(userClaims)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func GetRoles(userClaims *jwt.MapClaims) []string {
	realmAccess, exist := (*userClaims)["realm_access"]
	if !exist {
		panic("no realm_access in token")
	}

	rolesInterface, exist := realmAccess.(map[string]interface{})["roles"]
	if !exist {
		panic("no roles in token")
	}
	rolesInterfaceArray, ok := rolesInterface.([]interface{})
	if !ok {
		panic("invalid roles in token")
	}
	var roles []string
	for _, role := range rolesInterfaceArray {
		roles = append(roles, role.(string))
	}
	return roles
}

func GetGroups(userClaims *jwt.MapClaims) []string {
	groupsInterface, exist := (*userClaims)["userGroups"]
	if !exist {
		panic("no groups in token")
	}

	groupsInterfaceArray, ok := groupsInterface.([]interface{})
	if !ok {
		panic("invalid groups in token")
	}
	var groups []string
	for _, group := range groupsInterfaceArray {
		groups = append(groups, group.(string))
	}
	return groups
}

func HasGroup(userClaims *jwt.MapClaims, group string) bool {
	groups := GetGroups(userClaims)
	for _, r := range groups {
		if r == group {
			return true
		}
	}
	return false
}

func GetUserName(userClaims *jwt.MapClaims) string {
	return *GetClaim(userClaims, "name")
}

func GetUserEmail(userClaims *jwt.MapClaims) string {
	return *GetClaim(userClaims, "email")
}

func GetKeycloakUserId(userClaims *jwt.MapClaims) string {
	return *GetClaim(userClaims, "sub")
}
