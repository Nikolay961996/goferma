package models

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type WithdrawnRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
