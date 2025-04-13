package UIcomponents

type HeaderData struct {
	Title string
	Name  string
}
type PageData struct {
	Header  HeaderData
	Content interface{} // now correctly named Content
}

// CoursesData holds the courses page data.
type CoursesData struct {
	Total        int
	Attempts     int
	InReview     int
	Approved     int
	DateSections []DateSection
}

// DateSection groups cards by a date label.
type DateSection struct {
	DateLabel string
	Cards     []Card
}

type Card struct {
	Title       string
	ID          int
	Date        string
	Lecturer    string
	Field       string
	Status      string
	StatusLabel string
}
