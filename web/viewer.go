package main

type Viewer struct {
	Token  *Token
	UserID int

	RequestHost   string
	RequestScheme string
}
