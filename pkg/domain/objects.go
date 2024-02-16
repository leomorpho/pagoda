package domain

type Profile struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Age          int     `json:"age"`
	Bio          string  `json:"bio"`
	ProfileImage *Photo  `json:"profileImage"`
	Gender       string  `json:"gender"`
	Km           int     `json:"km"`
	Miles        int     `json:"miles"`
	Photos       []Photo `json:"photos"`
}

type Photo struct {
	Key string `json:"key"`
	Url string `json:"url"`
	Alt string `json:"alt"`
}
