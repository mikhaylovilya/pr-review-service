package entities

type ErrorResponse struct {
	ErrorBody ErrorBody `json:"error_body"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewErrorResponse(code string, msg string) *ErrorResponse {
	errBody := ErrorBody{
		Code:    code,
		Message: msg,
	}
	return &ErrorResponse{
		ErrorBody: errBody,
	}
}

func (e *ErrorResponse) Error() string {
	return "Error: " + e.ErrorBody.Message + ", Code: " + e.ErrorBody.Code
}

func ErrTeamExists(entityName string) error {
	return NewErrorResponse("TEAM_EXISTS", entityName+" already exists")
}

func ErrPRExists(entityName string) error {
	return NewErrorResponse("PR_EXISTS", entityName+" already exists")
}

func ErrPRMerged(entityName string) error {
	return NewErrorResponse("PR_MERGED", entityName+" already merged")
}

func ErrNotAssigned(entityName string) error {
	return NewErrorResponse("NOT_ASSIGNED", entityName+" wasn't assigned")
}

func ErrNoCandidate(entityName string) error {
	return NewErrorResponse("NO_CANDIDATE", "no other candidate, except for "+entityName)
}

func ErrNotFound(entityName string) error {
	return NewErrorResponse("NOT_FOUND", entityName+" not found")
}

func ErrUserExists(entityName string) error {
	return NewErrorResponse("USER_EXISTS", entityName+" already exists")
}

func ErrGeneric(msg string) error {
	return NewErrorResponse("GENERIC", msg)
}
