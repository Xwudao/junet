package jwtx

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	//key     = "junet jwt key"
	secret  = "junet jwt secret"
	issuer  = "junet"
	subject = "junet"
)

var config = Config{
	//Key:     key,
	Secret:  secret,
	Expire:  time.Hour * 24,
	Method:  jwt.SigningMethodHS256,
	Issuer:  issuer,
	Subject: subject,
}

type JPayload map[string]interface{}

type Opt func(*Config)
type Config struct {
	//Key     string        `json:"key"`
	Secret  string        `json:"secret"`
	Expire  time.Duration `json:"expire"`
	Issuer  string        `json:"issuer"`
	Subject string        `json:"subject"`

	Method jwt.SigningMethod
}

func SetIssuer(i string) Opt {
	return func(config *Config) {
		config.Issuer = i
	}
}
func SetMethod(m jwt.SigningMethod) Opt {
	return func(config *Config) {
		config.Method = m
	}
}
func SetExpire(expire time.Duration) Opt {
	return func(config *Config) {
		config.Expire = expire
	}
}
func SetSecret(secret string) Opt {
	return func(config *Config) {
		config.Secret = secret
	}
}
func SetSubject(sub string) Opt {
	return func(config *Config) {
		config.Subject = sub
	}
}

//func SetKey(key string) Opt {
//	return func(config *Config) {
//		config.Key = key
//	}
//}

type JClaims struct {
	Payload JPayload `json:"payload,omitempty"`
	jwt.StandardClaims
}

func Init(opts ...Opt) {
	for _, opt := range opts {
		opt(&config)
	}
}

func Generate(payload JPayload) (string, error) {
	now := time.Now()
	var claims = JClaims{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(config.Expire).Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    config.Issuer,
			NotBefore: now.Unix(),
			Subject:   config.Subject,
		},
	}
	token := jwt.NewWithClaims(config.Method, &claims)
	signedToken, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ParseToken(str string) (JPayload, error) {
	token, err := jwt.ParseWithClaims(str, &JClaims{}, func(token *jwt.Token) (i interface{}, err error) { // 解析token
		return []byte(config.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("token not valid")
	}
	if claims, ok := token.Claims.(*JClaims); ok {
		return claims.Payload, nil
	}
	return nil, fmt.Errorf("illegal token")
}
