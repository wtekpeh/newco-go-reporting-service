package middleware

import (
	"context"
	"fmt"
	"strings"

	"newco-go-reporting-service/internal/config"
	"newco-go-reporting-service/internal/dto"
	"newco-go-reporting-service/internal/services"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const AccessContextKey = "access_context"

type AuthMiddleware struct {
	cfg           *config.Config
	accessService services.AccessService
	jwks          keyfunc.Keyfunc
	issuer        string
}

func NewAuthMiddleware(cfg *config.Config, accessService services.AccessService) (*AuthMiddleware, error) {
	jwksURL := fmt.Sprintf(
		"%s/realms/%s/protocol/openid-connect/certs",
		strings.TrimRight(cfg.KeycloakServerURL, "/"),
		cfg.KeycloakRealm,
	)

	jwks, err := keyfunc.NewDefaultCtx(context.Background(), []string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("failed to create JWKS client: %w", err)
	}

	issuer := fmt.Sprintf(
		"%s/realms/%s",
		strings.TrimRight(cfg.KeycloakServerURL, "/"),
		cfg.KeycloakRealm,
	)

	return &AuthMiddleware{
		cfg:           cfg,
		accessService: accessService,
		jwks:          jwks,
		issuer:        issuer,
	}, nil
}

func (m *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"detail": "missing Authorization header",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"detail": "invalid Authorization header",
			})
		}

		tokenString := strings.TrimSpace(parts[1])

		token, err := jwt.Parse(tokenString, m.jwks.Keyfunc, jwt.WithIssuer(m.issuer))
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"detail": "invalid or expired token",
			})
		}

		claimsMap, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"detail": "invalid token claims",
			})
		}

		sub, ok := claimsMap["sub"].(string)
		if !ok || strings.TrimSpace(sub) == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"detail": "token missing sub claim",
			})
		}

		access, err := m.accessService.ResolveAccessContext(c.UserContext(), sub)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"detail": "no active staff profile found for this token",
			})
		}

		c.Locals(AccessContextKey, access)
		return c.Next()
	}
}

func (m *AuthMiddleware) RequireExecutive() fiber.Handler {
	return func(c *fiber.Ctx) error {
		value := c.Locals(AccessContextKey)
		access, ok := value.(*dto.AccessContext)
		if !ok || access == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"detail": "access context missing",
			})
		}

		if !m.accessService.IsExecutive(access) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"detail": "executive access required",
			})
		}

		return c.Next()
	}
}

func (m *AuthMiddleware) RequireBranchManager() fiber.Handler {
	return func(c *fiber.Ctx) error {
		value := c.Locals(AccessContextKey)
		access, ok := value.(*dto.AccessContext)
		if !ok || access == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"detail": "access context missing",
			})
		}

		if !m.accessService.IsBranchManager(access) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"detail": "branch manager access required",
			})
		}

		return c.Next()
	}
}
