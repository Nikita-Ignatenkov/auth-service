package middleware_test

//func TestAuthMiddleware(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockLogger := mock.NewMockLogger(ctrl)
//	mockLogger.EXPECT().Info("Authenticating request").Times(2)
//	mockLogger.EXPECT().Warn("Token missing").Times(1)
//	mockLogger.EXPECT().Warn("Invalid token").Times(1)
//	mockLogger.EXPECT().Warn("Invalid token claims").Times(1)
//	mockLogger.EXPECT().Warn("Invalid user_id claim").Times(1)
//	mockLogger.EXPECT().Info("Request authenticated successfully").Times(1)
//
//	mockContext := mock.NewMockContext(ctrl)
//	mockContext.EXPECT().GetString("secret_key").Return("secret").AnyTimes()
//	mockContext.EXPECT().Set("user_id", uint(123)).Times(1)
//	mockContext.EXPECT().Abort().Times(4)
//	mockContext.EXPECT().Next().Times(1)
//
//	tokenString := "valid_token"
//	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, jwt.ErrSignatureInvalid
//		}
//		return []byte("secret"), nil
//	})
//
//	claims := jwt.MapClaims{
//		"user_id": 123,
//	}
//
//	token.Claims = claims
//
//	mockContext.EXPECT().GetHeader("Authorization").Return(tokenString).Times(1)
//
//	authMiddleware := AuthMiddleware()
//	authMiddleware(mockContext)
//
//	assert.Equal(t, mockContext, gin.Context{})
//}
