package entities

type Sight struct {
	ID          int      `json:"id"`
	Rating      float32  `json:"rating"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	City        string   `json:"city"`
	Images      []string `json:"images"`
}

func GetSightsList() []Sight {
	return sights
}

func (h Sight) Validate() error {
	return nil
}

var sights = []Sight{
	{
		ID:          1,
		Rating:      4.3434,
		Name:        "Парижская башня",
		Description: "Самая высокая башня в мире.",
		City:        "Париж",
		Images:      []string{"path/to/image1.jpg", "path/to/image2.jpg"},
	},
	{
		ID:          2,
		Rating:      4.9,
		Name:        "Колосс Родосский",
		Description: "Один из семи чудес древнего мира.",
		City:        "Родос",
		Images:      []string{"path/to/image3.jpg"},
	},
	{
		ID:          2,
		Rating:      4.99,
		Name:        "Дом-музей Паисия Мальцева",
		Description: "Просто потому что находится рядом с Вольском.",
		City:        "Балаково",
		Images:      []string{"path/to/image4.jpg"},
	},
}
