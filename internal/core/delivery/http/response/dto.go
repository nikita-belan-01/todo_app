package response

import "github.com/nikita-belan-01/todo_app/internal/core/domain"

type problem struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type errorBody struct {
	Message  string    `json:"message"`
	Problems []problem `json:"errors,omitempty"`
}

func buildProblems(problems []domain.Problem) []problem {
	if len(problems) == 0 {
		return nil
	}

	result := make([]problem, 0, len(problems))
	for _, p := range problems {
		result = append(result, problem{Field: p.Field, Message: p.Message})
	}

	return result
}
