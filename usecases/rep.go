package usecases

type Database interface {
	GetResponse(intent string) (string, error)
	FindIntentByKeywords(query string) (string, string, error)
	GetAllIntents() ([]string, error)
}
