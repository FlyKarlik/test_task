package apiserver

type User struct {
	Id           string `bson:"id"`
	RefreshToken string `bson:"refresh_token"`
}
