package auth

type JWTService struct {
	secretKey []byte
}

func NewJWTService(secretKey []byte) *JWTService {
	return &JWTService{
		secretKey: secretKey,
	}
}
