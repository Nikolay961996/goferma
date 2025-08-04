package services

import (
	"encoding/json"
	"github.com/Nikolay961996/goferma/internal/models"
	"github.com/Nikolay961996/goferma/internal/utils"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func BuildJWTToken(userID int, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(models.JWT_EXPIRE_TIME)),
		},
		UserID: userID,
	})

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		utils.Log.Error("Error signing token")
		return "", err
	}

	return signedToken, nil
}

func ReadAuthModel(r *http.Request) (*models.AuthRequest, error) {
	bytes, err := utils.ReadJSONBody(r)
	if err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}

	var model models.AuthRequest
	err = json.Unmarshal(bytes, &model)
	if err != nil {
		utils.Log.Error(err.Error())
		return nil, &models.FormatError{Err: err}
	}

	return &model, nil
}

/*
func GetUserID(jwtToken string, secretKey string) (int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		utils.Log.Error("Error parsing token")
		return 0, err
	}

	if !token.Valid {
		utils.Log.Error("Invalid token")
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}
*/

func IsJWBTokenValid(jwtToken string, secretKey string) bool {
	token, err := jwt.ParseWithClaims(jwtToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		utils.Log.Error("Error parsing token")
		return false
	}
	return token.Valid
}
