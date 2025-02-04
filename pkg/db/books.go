package db

type Book struct {
	ID       int    `json:"id"`
	BOOK     string `json:"book"`
	AUTHOR   string `json:"author"`
	PAGES    int    `json:"pages"`
	TYPE     string `json:"type"`
	LANGUAGE string `json:"language"`
	DATE     string `json:"date"`
}

type Books struct {
	BOOKS []Book `json:"books"`
}

func Insert() {

}
