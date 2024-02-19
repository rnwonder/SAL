package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rnwonder/SAL/data"
	"github.com/rnwonder/SAL/util"
	"os"
	"time"
)

func SignDataWithJWT(data *data.Merchant) string {
	var (
		key []byte
		t   *jwt.Token
	)

	secret := util.MyCmpWorkAround(os.Getenv("JWT_SECRET"), "some-secret")

	key = []byte(secret)

	t = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
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

	return encodedString
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

func LoginUser(sess *session.Session, user *data.Merchant) (fiber.Map, string) {
	token := SignDataWithJWT(user)
	sess.Set("token", token)
	sess.Set("skuId", user.SkuId)
	sess.Set("id", user.Id)
	sess.SetExpiry(time.Hour * 24 * 30)
	err := sess.Save()

	if err != nil {
		log.Error("Error saving session", err)
	}

	return fiber.Map{
		"email":     user.Email,
		"name":      user.Name,
		"skuId":     user.SkuId,
		"id":        user.Id,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	}, token
}
