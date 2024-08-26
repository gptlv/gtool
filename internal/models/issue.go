package models

type BlockByIssuePayload struct {
	Transition struct {
		ID string `json:"id"`
	} `json:"transition"`
	Update struct {
		Issuelinks []struct {
			Add struct {
				Type struct {
					Name string `json:"name"`
				} `json:"type"`
				InwardIssue struct {
					Key string `json:"key"`
				} `json:"inwardIssue"`
			} `json:"add"`
		} `json:"issuelinks"`
	} `json:"update"`
}

type BlockUntilTomorrowPayload struct {
	Transition struct {
		ID int `json:"id"`
	} `json:"transition"`
	Fields struct {
		Customfield10253 string `json:"customfield_10253"`
	} `json:"fields"`
}

type InternalCommentPayload struct {
	Body       string `json:"body"`
	Properties []struct {
		Key   string `json:"key"`
		Value struct {
			Internal bool `json:"internal"`
		} `json:"value"`
	} `json:"properties"`
}

type UpdateSummary struct {
	Fields struct {
		Summary string `json:"summary"`
	} `json:"fields"`
}

type DeclinePayload struct {
	Transition struct {
		ID string `json:"id"`
	} `json:"transition"`
	Fields struct {
		Resolution struct {
			Name string `json:"name"`
		} `json:"resolution"`
	} `json:"fields"`
}
