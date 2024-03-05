package entities

type Sight struct {
	ID          int     `json:"id"`
	Rating      float32 `json:"rating"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	City        string  `json:"city"`
	Url         string  `json:"url"`
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
		Rating:      2.1,
		Name:        "У дяди Вани",
		Description: "Ресторан с видом на Сталинскую высотку.",
		City:        "Москва",
		Url:         "1.jpg",
	},
	{
		ID:          2,
		Rating:      3.1,
		Name:        "Государственный музей изобразительных искусств имени А.С. Пушкина",
		Description: "Музей.",
		City:        "Москва",
		Url:         "2.jpg",
	},
	{
		ID:          3,
		Rating:      4.99,
		Name:        "МГТУ им. Н. Э. Баумана",
		Description: "Хороший вуз.",
		City:        "Москва",
		Url:         "3.jpg",
	},
	{
		ID:          4,
		Rating:      3.2,
		Name:        "Вкусно - и точка",
		Description: "Неплохое кафе, вызывает гастрит.",
		City:        "Москва",
		Url:         "4.jpg",
	},
	{
		ID:          5,
		Rating:      4.1,
		Name:        "Краеведческий музей",
		Description: "Один из самых больших провинциальных музеев краеведческого профиля.",
		City:        "Вольск",
		Url:         "5.jpg",
	},
}
