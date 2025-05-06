package UIcomponents

// SyllabusRow represents one row of the syllabus table.
type SyllabusRow struct {
	LessonNumber    string
	MainTopic       string
	LessonTopics    string
	Subtopics       string
	ReadingMaterial string
}

// GradeComponent represents one row of the grade composition.
type GradeComponent struct {
	PartName   string
	Percentage string
}

type Comments struct {
	Name          string
	Message       string
	Time          string
	IsCurrentUser bool
}

type Draft struct {
	ID                      int              `json:"ID"`
	LecturerName            string           `json:"lecturerName"`
	LecturerEmail           string           `json:"lecturerEmail"`
	OfficeDay               string           `json:"officeDay"`
	OfficeStart             string           `json:"officeStart"`
	OfficeEnd               string           `json:"officeEnd"`
	SyllabusDepartment      string           `json:"syllabusDepartment"` // Selected department
	Departments             []string         `json:"departments"`
	SelectedCourse          string           `json:"selectedCourse"` // Selected course
	Courses                 []string         `json:"courses"`        // List of available courses
	CourseRequirements      []string         `json:"courseRequirements"`
	LearningOutcomes        []string         `json:"learningOutcomes"`
	CourseObjectives        []string         `json:"courseObjectives"`
	Credits                 string           `json:"credits"`
	WeeklyHours             string           `json:"weeklyHours"`
	Year                    string           `json:"year"`
	Semester                string           `json:"semester"`
	Prerequisites           string           `json:"prerequisites"`
	CourseStructure         []string         `json:"courseStructure"`
	OtherCourseStructure    string           `json:"otherCourseStructure"`
	ActiveLearning1         string           `json:"activeLearning1"`
	ActiveLearning2         string           `json:"activeLearning2"`
	ActiveLearning3         string           `json:"activeLearning3"`
	ActiveLearning4         string           `json:"activeLearning4"`
	SyllabusRows            []SyllabusRow    `json:"syllabusRows"`
	GradeComponents         []GradeComponent `json:"gradeComponents"`
	AssignmentsStructure    []string         `json:"assignmentsStructure"`
	BibliographyRequired    []string         `json:"bibliographyRequired"`
	BibliographyRecommended []string         `json:"bibliographyRecommended"`
}
