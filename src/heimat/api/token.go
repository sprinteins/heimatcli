package api

import (
	jwt "github.com/golang-jwt/jwt/v4"
)

// ExtractEmployeeID _
func ExtractEmployeeID(tokenString string) int {
	claims := &heimatClaim{}
	// TODO: Check for errors and ignore missing key function error
	jwt.ParseWithClaims(tokenString, claims, nil)
	// jwt.Parse(tokenString, nil)

	return claims.EmployeeID
}

type heimatClaim struct {
	EmployeeID int `json:"employeeId"`
	jwt.StandardClaims
}
