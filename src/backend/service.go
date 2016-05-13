package backend

import (
	"sync"
	logs "github.com/gobricks/jwtack/src/loggers"
	metrics "github.com/gobricks/jwtack/src/metrics"
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
	logger logs.AppLogs
	mtx    sync.RWMutex
}

func (s *service) CreateToken(key string, payload map[string]interface{}, exp *time.Duration) (t string, err error) {
	if key == "" {
		return "", fmt.Errorf("Empty required key")
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = payload
	if exp != nil {
		token.Claims["exp"] = time.Now().Add(*exp).Unix()
	}
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
		s.logger.Error.Log("jwt.Parse", err.Error())
		err = fmt.Errorf("Incorrect token")
	} else {
		payload = t.Claims
		if !t.Valid {
			err = fmt.Errorf("Invalid token")
		}
	}

	return
}

func InitService(appLogs logs.AppLogs, appMetrics metrics.AppMetrics) Service {
	var svc Service
	{
		svc = &service{
			logger:appLogs,
		}
		svc = loggingMiddleware(appLogs)(svc)
		svc = metricsMiddleware(appMetrics)(svc)
	}
	return svc
}