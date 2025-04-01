package usecases

import (
	"service/entities"
)

type NLPService interface {
	AnalyzeIntent(query string, db Database) (string, string, error)
}

type QueryProcessor struct {
	nlp NLPService
	db  Database
}

func NewQueryProcessor(nlp NLPService, db Database) *QueryProcessor {
	return &QueryProcessor{nlp: nlp, db: db}
}

func (qp *QueryProcessor) ProcessQuery(queryText string) (*entities.Query, error) {
	intent, response, err := qp.nlp.AnalyzeIntent(queryText, qp.db)
	if err != nil {
		return nil, err
	}

	return &entities.Query{Text: queryText, Intent: intent, Response: response}, nil
}
