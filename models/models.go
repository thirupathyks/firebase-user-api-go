package models

//User - as defined in Firebase SDK
type User struct {
	UID           string `json:"uid"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	PhoneNumber   string `json:"phoneNumber"`
	Password      string `json:"password"`
	DisplayName   string `json:"displayName"`
	PhotoURL      string `json:"photoURL"`
	Disabled      bool   `json:"disable"`
}

//SignInUserRequest - as defined in Firebase REST API
type SignInUserRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}
