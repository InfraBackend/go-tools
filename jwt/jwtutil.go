package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var Secret = []byte("what a secret")

// jwt过期时间, 按照实际环境设置
const expiration = 2 * time.Minute

type Claims struct {
	// 自定义字段, 可以存在用户名, 用户ID, 用户角色等等
	Username string
	// jwt.StandardClaims包含了官方定义的字段
	jwt.StandardClaims
}

func GenToken(username string) (string, error) {
	// 创建声明
	a := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiration).Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(),                 // 签发时间
			Issuer:    "gin-jwt-demo",                    // 签发者
			Id:        "",                                // 按需求选这个, 有些实现中, 会控制这个ID是不是在黑/白名单来判断是否还有效
			NotBefore: 0,                                 // 生效起始时间
			Subject:   "",                                // 主题
		},
	}

	// 用指定的哈希方法创建签名对象
	tt := jwt.NewWithClaims(jwt.SigningMethodHS256, a)
	// 用上面的声明和签名对象签名字符串token
	// 1. 先对Header和PayLoad进行Base64URL转换
	// 2. Header和PayLoadBase64URL转换后的字符串用.拼接在一起
	// 3. 用secret对拼接在一起之后的字符串进行HASH加密
	// 4. 连在一起返回
	return tt.SignedString(Secret)
}

func ParseToken(tokenStr string) (*Claims, error) {
	// 第三个参数: 提供一个回调函数用于提供要选择的秘钥, 回调函数里面的token参数,是已经解析但未验证的,可以根据token里面的值做一些逻辑, 如`kid`的判断
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	// 校验token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
