package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rnwonder/SAL/data"
	"github.com/rnwonder/SAL/util"
	"os"
	"time"
)

func SignDataWithJWT(data *data.Merchant) (string, time.Time) {
	var (
		key []byte
		t   *jwt.Token
	)

	secret := util.MyCmpWorkAround(os.Getenv("JWT_SECRET"), "some-secret")

	key = []byte(secret)

	expiresAt := time.Now().Add(time.Hour * 24 * 30)

	t = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expiresAt.Unix(),
		"iss": "SAL",
		"data": map[string]string{
			"email":    data.Email,
			"name":     data.Name,
			"Id":       data.Id,
			"password": data.Password,
			"skuId":    data.SkuId,
		},
	})

	encodedString, err := t.SignedString(key)

	if err != nil {
		log.Error("Error signing token", err)
	}

	return encodedString, expiresAt
}

func DecodeJWTData(encodedString string) (data.Merchant, error) {
	var (
		key  []byte
		t    *jwt.Token
		user data.Merchant
	)

	secret := util.MyCmpWorkAround(os.Getenv("JWT_SECRET"), "some-secret")

	key = []byte(secret)

	t, err := jwt.Parse(encodedString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return user, err
	}

	claims := t.Claims.(jwt.MapClaims)

	userData := claims["data"].(map[string]interface{})
	userId := userData["Id"].(string)

	for _, value := range data.MerchantData {
		if value.Id == userId {
			return value, nil
		}
	}

	return user, errors.New("User not found")
}

func LoginUser(user *data.Merchant) (fiber.Map, string, time.Time) {
	token, expiredAt := SignDataWithJWT(user)
	return fiber.Map{
		"email":     user.Email,
		"name":      user.Name,
		"skuId":     user.SkuId,
		"id":        user.Id,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	}, token, expiredAt
}
