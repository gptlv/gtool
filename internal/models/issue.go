package models

const EMAIL_FIELD_KEY = "customfield_10145"

var internalCommentPayloadBody = `{
	"body": "%s",
	"properties": [
	  {
		"key": "sd.public.comment",
		"value": {
		   "internal": true
		}
	  }
	]
 }`

var blockByIssuePayloadBody = `
{
	"transition": {
		"id": "%v"
	},
	"update": {
		"issuelinks": [
			{
				"add": {
					"type": {
						"name": "Blocks"
					},
					"inwardIssue": {
						"key": "%v"
					}
				}
			}
		]
	}
}`

var BlockUntilTomorrowPayloadBody = `
{
    "transition": {
        "id": "%v"
    },
    "fields": {
        "customfield_10253": "%v"
    }
}
`

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
		ID string `json:"id"`
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
