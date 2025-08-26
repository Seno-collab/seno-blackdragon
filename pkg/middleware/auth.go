package middleware

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		token := c.GetHeader("Authorization")
// 		reqTime := time.Now()
// 		traceID := c.GetHeader(HeaderKeyTraceID)
// 		if token == "" {
// 			dto.WriteJSON(c, http.StatusUnauthorized, dto.NewError(http.StatusUnauthorized, enum.AUTH_FAILED,
// 				"Invalid email or password", traceID, reqTime, enum.ErrInvalidToken))
// 			return
// 		}
// 		// verify token
// 		// redis := &
// 		c.Next()
// 	}
// }
