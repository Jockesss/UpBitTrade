package domain

type Token struct {
	AccessKey string
	SecretKey string
	Algorithm string
	Exchange  int
	Type      int
}
