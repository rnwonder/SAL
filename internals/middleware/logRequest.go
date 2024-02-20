package middleware

import (
	"cmp"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"strings"
)

func LogRequest(ctx *fiber.Ctx) error {
	method := ctx.Method()
	path := ctx.OriginalURL()
	ip := ctx.IP()
	body := ctx.Body()

	stringSlice := []string{"Method: ", method, " Path: ", path, " IP: ", ip, "\n", "Body: ", cmp.Or(string(body), "{}")}
	var stringBuilder strings.Builder

	for _, v := range stringSlice {
		stringBuilder.WriteString(v)
	}
	log.Info(stringBuilder.String())

	return ctx.Next()
}
