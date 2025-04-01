package delivery

import (
	"encoding/json"
	"net/http"
	"service/usecases"
)

type QueryHandler struct {
	queryProcessor *usecases.QueryProcessor
}

func NewQueryHandler(qp *usecases.QueryProcessor) *QueryHandler {
	return &QueryHandler{queryProcessor: qp}
}

func (qh *QueryHandler) HandleQuery(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Text string `json:"text"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	response, err := qh.queryProcessor.ProcessQuery(request.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"response": response.Response})
}
