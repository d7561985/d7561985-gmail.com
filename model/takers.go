package model

type TakerList []Taker

type Taker struct {
	Login     string `json:"login" csv:"login"`
	Password  string `json:"password" csv:"password"`
	Title     string `json:"title" csv:"title"`
	LastName  string `json:"lastname" csv:"lastname"`
	FirstName string `json:"firstname" csv:"firstname"`
	Gender    string `json:"gender" csv:"gender"`
	Email     string `json:"email" csv:"email"`
	Picture   string `json:"picture" csv:"picture"`
	Address   string `json:"address" csv:"address"`
}
