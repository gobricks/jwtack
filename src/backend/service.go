package backend

import (
	"sync"
	"github.com/gobricks/jwtack/src/app"
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
)

type ServiceMW func(Service) Service

type Service interface {
	CreateToken(key string, payload map[string]interface{}, exp *time.Duration) (t string, err error)
	ParseToken(token string, key string) (payload map[string]interface{}, err error)
}

type service struct {
	app app.App
	mtx sync.RWMutex
}

func (s *service) CreateToken(key string, payload map[string]interface{}, exp *time.Duration) (t string, err error) {
	if key == "" {
		return "", fmt.Errorf("Empty required key")
	}
	var claims jwt.MapClaims
	claims = payload
	if exp != nil {
		claims["exp"] = time.Now().Add(*exp).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.SigningString()

	t, err = token.SignedString([]byte(key))
	return
}
func (s *service) ParseToken(token string, key string) (payload map[string]interface{}, err error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(key), nil
	})

	if err != nil {
		s.app.Logs.Error.Log("jwt.Parse", err.Error())
		err = fmt.Errorf("Incorrect token")
	} else {
		payload = t.Claims.(jwt.MapClaims)
		if !t.Valid {
			err = fmt.Errorf("Invalid token")
		}
	}

	return
}

func InitService(app app.App) Service {
	var svc Service
	{
		svc = &service{
			app:app,
		}
		svc = loggingMiddleware(app.Logs)(svc)
		svc = metricsMiddleware(app.Metrics)(svc)
	}
	return svc
}