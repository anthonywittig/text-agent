package types

// https://docs.aws.amazon.com/bedrock/latest/userguide/agents-lambda.html
type AgentRequest struct {
	MessageVersion string `json:"messageVersion"`
	Function       string `json:"function"`
	Parameters     []struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"parameters"`
	InputText string `json:"inputText"`
	SessionId string `json:"sessionId"`
	Agent     struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Id      string `json:"id"`
		Alias   string `json:"alias"`
	} `json:"agent"`
	ActionGroup string `json:"actionGroup"`
	// Not sure what this looks like in practice.
	SessionAttributes interface{} `json:"sessionAttributes,omitempty"`
	// Not sure what this looks like in practice.
	PromptSessionAttributes interface{} `json:"promptSessionAttributes,omitempty"`
}

// https://docs.aws.amazon.com/bedrock/latest/userguide/agents-lambda.html
type AgentResponse struct {
	MessageVersion string                `json:"messageVersion"`
	Response       AgentResponseResponse `json:"response"`
	// SessionAttributes           interface{} `json:"sessionAttributes,omitempty"`
	// PromptSessionAttributes     interface{} `json:"promptSessionAttributes,omitempty"`
	// KnowledgeBasesConfiguration []struct {
	// 	KnowledgeBaseId        string `json:"knowledgeBaseId"`
	// 	RetrievalConfiguration struct {
	// 		VectorSearchConfiguration struct {
	// 			NumberOfResults int `json:"numberOfResults"`
	// 			Filter          struct {
	// 				RetrievalFilter struct {
	// 					Field    string `json:"field"`
	// 					Operator string `json:"operator"`
	// 					Value    string `json:"value"`
	// 				} `json:"filter"`
	// 			} `json:"vectorSearchConfiguration"`
	// 		} `json:"retrievalConfiguration"`
	// 	} `json:"retrievalConfiguration"`
	// } `json:"knowledgeBasesConfiguration"`
}

type AgentResponseResponse struct {
	ActionGroup      string                                `json:"actionGroup"`
	Function         string                                `json:"function"`
	FunctionResponse AgentResponseResponseFunctionResponse `json:"functionResponse"`
}

type AgentResponseResponseFunctionResponse struct {
	ResponseState string                                            `json:"responseState"`
	ResponseBody  AgentResponseResponseFunctionResponseResponseBody `json:"responseBody"`
}

type AgentResponseResponseFunctionResponseResponseBody struct {
	ContentType AgentResponseResponseFunctionResponseResponseBodyContentType `json:"TEXT"`
}

type AgentResponseResponseFunctionResponseResponseBodyContentType struct {
	Body string `json:"body"` // This should be a JSON string.
}
