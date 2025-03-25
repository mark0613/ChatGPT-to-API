package chatgpt

type ChatGPTResponse struct {
	Message        Message     `json:"message"`
	ConversationID string      `json:"conversation_id"`
	Error          interface{} `json:"error"`
}
type ChatGPTWSSResponse struct {
	WssUrl         string `json:"wss_url"`
	ConversationId string `json:"conversation_id,omitempty"`
	ResponseId     string `json:"response_id,omitempty"`
}

type WSSMsgResponse struct {
	SequenceId int                `json:"sequenceId"`
	Type       string             `json:"type"`
	From       string             `json:"from"`
	DataType   string             `json:"dataType"`
	Data       WSSMsgResponseData `json:"data"`
}

type WSSMsgResponseData struct {
	Type           string `json:"type"`
	Body           string `json:"body"`
	MoreBody       bool   `json:"more_body"`
	ResponseId     string `json:"response_id"`
	ConversationId string `json:"conversation_id"`
}

type Message struct {
	ID         string      `json:"id"`
	Author     Author      `json:"author"`
	CreateTime float64     `json:"create_time"`
	UpdateTime interface{} `json:"update_time"`
	Content    Content     `json:"content"`
	EndTurn    interface{} `json:"end_turn"`
	Weight     float64     `json:"weight"`
	Metadata   Metadata    `json:"metadata"`
	Recipient  string      `json:"recipient"`
}

type Content struct {
	ContentType string        `json:"content_type"`
	Parts       []interface{} `json:"parts"`
}

type Author struct {
	Role     string                 `json:"role"`
	Name     interface{}            `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Metadata struct {
	Citations         []Citation         `json:"citations,omitempty"`
	ContentReferences []ContentReference `json:"content_references,omitempty"`
	SearchResultGroups []SearchResultGroup `json:"search_result_groups,omitempty"`
	SafeURLs          []string           `json:"safe_urls,omitempty"`
	MessageType       string             `json:"message_type"`
	FinishDetails     *FinishDetails     `json:"finish_details"`
	ModelSlug         string             `json:"model_slug"`
}

type ContentReference struct {
	MatchedText string   `json:"matched_text"`
	StartIdx    int      `json:"start_idx"`
	EndIdx      int      `json:"end_idx"`
	Refs        []string `json:"refs"`
	Type        string   `json:"type"`
	Invalid     bool     `json:"invalid"`
}

type SearchResultGroup struct {
	Type    string         `json:"type"`
	Domain  string         `json:"domain"`
	Entries []SearchResult `json:"entries"`
}

type SearchResult struct {
	Type        string      `json:"type"`
	URL         string      `json:"url"`
	Title       string      `json:"title"`
	Snippet     string      `json:"snippet"`
	RefID       RefID       `json:"ref_id"`
	ContentType interface{} `json:"content_type"`
	PubDate     float64     `json:"pub_date,omitempty"`
	Attribution string      `json:"attribution"`
}

type RefID struct {
	TurnIndex int    `json:"turn_index"`
	RefType   string `json:"ref_type"`
	RefIndex  int    `json:"ref_index"`
}
type Citation struct {
	Metadata CitaMeta `json:"metadata"`
	StartIx  int      `json:"start_ix"`
	EndIx    int      `json:"end_ix"`
}
type CitaMeta struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}
type FinishDetails struct {
	Type string `json:"type"`
	Stop string `json:"stop"`
}
type DalleContent struct {
	AssetPointer string `json:"asset_pointer"`
	Metadata     struct {
		Dalle struct {
			Prompt string `json:"prompt"`
		} `json:"dalle"`
	} `json:"metadata"`
}
