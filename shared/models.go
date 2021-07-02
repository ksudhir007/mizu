package shared

import "time"

type WebSocketMessageType string

const (
	WebSocketMessageTypeEntry        WebSocketMessageType = "entry"
	WebSocketMessageTypeTappedEntry  WebSocketMessageType = "tappedEntry"
	WebSocketMessageTypeUpdateStatus WebSocketMessageType = "status"
)

type WebSocketMessageMetadata struct {
	MessageType WebSocketMessageType `json:"messageType,omitempty"`
}

type WebSocketStatusMessage struct {
	*WebSocketMessageMetadata
	TappingStatus TapStatus `json:"tappingStatus"`
}

type TapStatus struct {
	Pods []PodInfo `json:"pods"`
}

type PodInfo struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func CreateWebSocketStatusMessage(tappingStatus TapStatus) WebSocketStatusMessage {
	return WebSocketStatusMessage{
		WebSocketMessageMetadata: &WebSocketMessageMetadata{
			MessageType: WebSocketMessageTypeUpdateStatus,
		},
		TappingStatus: tappingStatus,
	}
}

type TrafficFilteringOptions struct {
	PlainTextMaskingRegexes []*SerializableRegexp
}

// Enforce Policy

type RulesPolicy struct {
	Rules []RulePolicy `yaml:"rules"`
}

type RulePolicy struct {
	Type    string `yaml:"type"`
	Service string `yaml:"service"`
	Path    string `yaml:"path"`
	Method  string `yaml:"method"`
	Key     string `yaml:"key"`
	Value   string `yaml:"value"`
	Latency int    `yaml:"latency"`
	Name    string `yaml:"name"`
}

func (r RulePolicy) validateType() bool {
	permitedTypes := []string{"json", "header", "latency"}
	_, found := Find(permitedTypes, r.Type)
	return found
}

func (rules RulesPolicy) ValidateRulesPolicy() []int {
	invalidIndex := make([]int, 0)
	for i := range rules.Rules {
		validated := rules.Rules[i].validateType()
		if !validated {
			invalidIndex = append(invalidIndex, i)
		}
	}
	return invalidIndex
}

func (rules *RulesPolicy) RemoveNotValidPolicy(idx int) {
	rules.Rules = append(rules.Rules[:idx], rules.Rules[idx+1:]...)
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

type HARFetched struct {
	Log struct {
		Version string `json:"version"`
		Creator struct {
			Name    string `json:"name"`
			Version string `json:"version"`
			Source  string `json:"_source"`
		} `json:"creator"`
		Entries []struct {
			ID              string    `json:"_id"`
			StartedDateTime time.Time `json:"startedDateTime"`
			Time            int       `json:"time"`
			Request         struct {
				Method      string        `json:"method"`
				URL         string        `json:"url"`
				HTTPVersion string        `json:"httpVersion"`
				Cookies     []interface{} `json:"cookies"`
				Headers     []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"headers"`
				QueryString []interface{} `json:"queryString"`
				HeadersSize int           `json:"headersSize"`
				BodySize    int           `json:"bodySize"`
			} `json:"request"`
			Response struct {
				Status      int           `json:"status"`
				StatusText  string        `json:"statusText"`
				HTTPVersion string        `json:"httpVersion"`
				Cookies     []interface{} `json:"cookies"`
				Headers     []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"headers"`
				Content struct {
					Size     int    `json:"size"`
					MimeType string `json:"mimeType"`
					Text     string `json:"text"`
					Encoding string `json:"encoding"`
				} `json:"content"`
				RedirectURL string `json:"redirectURL"`
				HeadersSize int    `json:"headersSize"`
				BodySize    int    `json:"bodySize"`
			} `json:"response"`
			Cache struct {
			} `json:"cache"`
			Timings struct {
				Send    int `json:"send"`
				Wait    int `json:"wait"`
				Receive int `json:"receive"`
			} `json:"timings"`
		} `json:"entries"`
	} `json:"log"`
}
