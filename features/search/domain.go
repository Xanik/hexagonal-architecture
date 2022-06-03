package search

//Usecase interface represents search usecases
type Usecase interface {
	// search Usecases
	IndexAccounts(string) (string, error)
	IndexContents(string) (string, error)
	IndexCourses(string) (string, error)
	IndexInterests(string) (string, error)
	SearchElastic(string, string, string) (interface{}, error)
	SearchInterest(string, string) (interface{}, error)
	SearchContent(string, string, string) (interface{}, error)
	SearchAccount(string, string, string) (interface{}, error)
	SearchCourse(string, string, string) (interface{}, error)
}
