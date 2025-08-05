package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nikolay961996/goferma/internal/models"
	"github.com/Nikolay961996/goferma/internal/storage"
	"github.com/Nikolay961996/goferma/internal/utils"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int64
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

func CreateUser(db *storage.DBContext, login string, password string) error {
	user, err := db.GetUser(login)
	if err != nil {
		utils.Log.Error(err.Error())
		return err
	}
	if user != nil {
		utils.Log.Error("User this login already exist", login)
		return &models.AlreadyExistError{Err: errors.New(fmt.Sprintf("user this '%s' already exist", login))}
	}

	pswHash := getPasswordHash(password)
	err = db.AddNewUser(login, pswHash)
	if err != nil {
		utils.Log.Error(err.Error())
		return err
	}

	return nil
}

func AuthUser(db *storage.DBContext, secretKey string, login string, password string) (string, error) {
	user, err := db.GetUser(login)
	if err != nil {
		utils.Log.Error(err.Error())
		if user == nil {
		}
		return "", err
	}

	pswHash := getPasswordHash(password)
	if user == nil || user.PasswordHash != pswHash {
		utils.Log.Error("user not found")
		return "", &models.LoginPasswordError{Err: errors.New(fmt.Sprintf("user this login/password '%s' does not exist", login))}
	}

	token, err := buildJWTToken(user.ID, secretKey)
	if err != nil {
		utils.Log.Error(err.Error())
		return "", err
	}

	return token, nil
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

/*
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
*/

func getPasswordHash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func buildJWTToken(userID int64, secretKey string) (string, error) {
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
