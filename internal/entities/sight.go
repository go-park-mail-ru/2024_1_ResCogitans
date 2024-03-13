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
		Url:         "public/1.jpg",
	},
	{
		ID:          2,
		Rating:      3.1,
		Name:        "Государственный музей изобразительных искусств имени А.С. Пушкина",
		Description: "Музей.",
		City:        "Москва",
		Url:         "public/2.jpg",
	},
	{
		ID:          3,
		Rating:      4.99,
		Name:        "МГТУ им. Н. Э. Баумана",
		Description: "Хороший вуз.",
		City:        "Москва",
		Url:         "public/3.jpg",
	},
	{
		ID:          4,
		Rating:      3.2,
		Name:        "Вкусно - и точка",
		Description: "Неплохое кафе, вызывает гастрит.",
		City:        "Москва",
		Url:         "public/4.jpg",
	},
	{
		ID:          5,
		Rating:      4.1,
		Name:        "Краеведческий музей",
		Description: "Один из самых больших провинциальных музеев краеведческого профиля.",
		City:        "Вольск",
		Url:         "public/5.jpg",
	},
	{
		ID:          6,
		Rating:      4.3,
		Name:        "Спасо-Преображенский кафедральный собор",
		Description: "Спасо-Преображенский кафедральный собор расположен в центре города и является первым каменным храмом Тамбова и старейшим в Тамбовской обл.",
		City:        "Тамбов",
		Url:         "public/6.jpg",
	},
	{
		ID:          7,
		Rating:      3.9,
		Name:        "Мирский замок",
		Description: "Памятник архитектуры, внесён в список Всемирного наследия ЮНЕСКО.",
		City:        "Мир",
		Url:         "public/7.jpg",
	},
	{
		ID:          8,
		Rating:      4.9,
		Name:        "Чуфут-Кале",
		Description: "Пещерный город в Крыму. Топ.",
		City:        "Бахчисарай",
		Url:         "public/8.jpg",
	},
	{
		ID:          9,
		Rating:      3.5,
		Name:        "Сасык-Сиваш",
		Description: "Розовое озеро. Оно реально розовое.",
		City:        "Евпатория",
		Url:         "public/9.jpg",
	},
	{
		ID:          10,
		Rating:      4.7,
		Name:        "Крепость Чембело",
		Description: "Остатки крепости.",
		City:        "Балаклава",
		Url:         "public/10.jpg",
	},
	{
		ID:          11,
		Rating:      4.0,
		Name:        "Мечеть Кул Шариф",
		Description: "Главная джума-мечеть республики Татарстан и города Казани.",
		City:        "Казань",
		Url:         "public/11.jpg",
	},
	{
		ID:          12,
		Rating:      4.5,
		Name:        "Салтинский Подземный Водопад",
		Description: "Единственный в России подземный водопад.",
		City:        "Салта",
		Url:         "public/12.jpg",
	},
}
