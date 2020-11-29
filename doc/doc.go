// Package classification Reverse Market API.
//
// Documentation for Reverse Market API.
//
//     Schemes: https
//     BasePath: /
//     Version: 1.0.0
//     Host: localhost
//     Contact: Sergey Kozhin<sergeykozhin@gmail.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//     - application/txt
//
//     Security:
//     - basic
//
//     Security:
//     - oauth2:
//
//     SecurityDefinitions:
//     oauth2:
//         type: oauth2
//         tokenUrl: /auth/sign_in
//         in: header
//
// swagger:meta
package doc

// swagger:route POST /auth/sign_in auth signIn
// Consumes Google ID token and produces JWT token, if ID token is valid
// responses:
//   200: SignInResponse
//   400: BadRequestError

// JSON containing JWT token.
// swagger:response SignInResponse
type SignInResponse struct {
	// JWT token
	// in:body
	Body struct {
		// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
		JwtToken string `json:"jwt_token"`
	}
}

// swagger:parameters signIn
type SignInParams struct {
	// Google Id token
	// in:body
	Body struct {
		//required: true
		// example: eyJhbGciOiJSUzI1NiIsImtpZCI6IjJlMzAyNWYyNmI1OTVmOTZlYWM5MDdjYzJiOTQ3MTQyMmJjYWViOTMiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI0MDc0MDg3MTgxOTIuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI0MDc0MDg3MTgxOTIuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMDIzNjk3NjA1MTAzMjYxNjkyNzEiLCJlbWFpbCI6Imt6aGludmxkbXJAZ21haWwuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiI1alROei1NSzVmQUIzdmNpY3cwU1ZBIiwibmFtZSI6ItCh0LXRgNCz0LXQuSDQmtC-0LbQuNC9IiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hLS9BT2gxNEdqQi1sakU2QzVIYWIySW03YVB4bG1ROW9HUlpwSExaODdNd215eT1zOTYtYyIsImdpdmVuX25hbWUiOiLQodC10YDQs9C10LkiLCJmYW1pbHlfbmFtZSI6ItCa0L7QttC40L0iLCJsb2NhbGUiOiJydSIsImlhdCI6MTYwNjY1NTU2OSwiZXhwIjoxNjA2NjU5MTY5fQ.X_B9DZXuwKV4A5KMBVtuCPERKCZh6ugH0s4AYzwbLzzouMzoUZCDMz0EjVVikTBjiPG86uMYUojWimjxcwZJ0TJKnspJ3BnFBBncHT4uxZtXik3MZdlUA0q0jIusXaR-GEXZGlzUVXYlfs8ppWP964k2a2NYxczJWFfvXRNi7SC9CSmvLUc4km6GG7abbqRtI-EJR61cWXjxNJwA9EeJz7n-rynKSEekTCa0ZJCYUtZ5aMJ7TKGM2tYlwy3s04-MMGuhodgk62Jfa0FoKAimnPDDljF5Q2A-t-pJRqO3x9IXOICf52aZQbFjzmFu8c7XgfPLiyiC4tO4RVa9R1dCYQ
		IDToken string `json:"id_token"`
	}
}

// swagger:route GET /auth/check auth check
// Check endpoint for JWT validation
// Security:
//   oauth2: read
// responses:
//   200: OkResponse
//   401: UnauthorizedError
//   400: BadRequestError

// swagger:route GET /user user getUser
// Provides information about current user
// Security:
//   oauth2: read
// responses:
//   200: UserResponse
//   401: UnauthorizedError
//   400: BadRequestError

// JSON containing user info.
// swagger:response UserResponse
type UserResponse struct {
	// User info
	// in:body
	Body struct {
		// example: 1
		ID int `json:"id"`
		// example: Иван Иванов
		Name string `json:"name"`
		// example: test@test.com
		Email string `json:"email"`
		// example: /avatars/1.png
		Avatar *string `json:"avatar"`
		//example: 1
		DefaultAddressID *int `json:"default_address_id"`
	}
}

// OK Response
// swagger:response OkResponse
type OkResponse struct{}

// Error message
// swagger:response BadRequestError
type BadRequestError struct{}

// Unauthorized error message
// swagger:response UnauthorizedError
type UnauthorizedError struct{}
