package types

import "time"

type DismissalRecord struct {
	ISC      string
	Flaw     string
	Decision string
}

type LaptopDescription struct {
	Serial      string
	ISC         string
	Name        string
	InventoryID string
	Cost        string
}

type DismissalDocument struct {
	*DismissalRecord
	*LaptopDescription
	Date string
	Boss string
	Lead string
	ID   int
}

type InsightUserAttributesPayload struct {
	Attributes []struct {
		ObjectTypeAttributeID int `json:"objectTypeAttributeId"`
		ObjectAttributeValues []struct {
			Value string `json:"value"`
		} `json:"objectAttributeValues"`
	} `json:"attributes"`
}

type InsightAttachment struct {
	ID            int       `json:"id"`
	Author        string    `json:"author"`
	MimeType      string    `json:"mimeType"`
	Filename      string    `json:"filename"`
	Filesize      string    `json:"filesize"`
	Created       time.Time `json:"created"`
	Comment       string    `json:"comment"`
	CommentOutput string    `json:"commentOutput"`
	URL           string    `json:"url"`
}

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

type ObjectSearchRes struct {
	ObjectEntries         []ObjectEntry         `json:"objectEntries"`
	ObjectTypeAttributes  []ObjectTypeAttribute `json:"objectTypeAttributes"`
	ObjectTypeID          int                   `json:"objectTypeId"`
	ObjectTypeIsInherited bool                  `json:"objectTypeIsInherited"`
	AbstractObjectType    bool                  `json:"abstractObjectType"`
	TotalFilterCount      int                   `json:"totalFilterCount"`
	StartIndex            int                   `json:"startIndex"`
	ToIndex               int                   `json:"toIndex"`
	PageObjectSize        int                   `json:"pageObjectSize"`
	PageNumber            int                   `json:"pageNumber"`
	OrderWay              string                `json:"orderWay"`
	QlQuery               string                `json:"qlQuery"`
	QlQuerySearchResult   bool                  `json:"qlQuerySearchResult"`
	ConversionPossible    bool                  `json:"conversionPossible"`
	IqlSearchResult       bool                  `json:"iqlSearchResult"`
	Iql                   string                `json:"iql"`
	PageSize              int                   `json:"pageSize"`
}
type Avatar struct {
	URL16    string `json:"url16"`
	URL48    string `json:"url48"`
	URL72    string `json:"url72"`
	URL144   string `json:"url144"`
	URL288   string `json:"url288"`
	ObjectID int    `json:"objectId"`
}
type Icon struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	URL16 string `json:"url16"`
	URL48 string `json:"url48"`
}
type ObjectType struct {
	ID                        int       `json:"id"`
	Name                      string    `json:"name"`
	Type                      int       `json:"type"`
	Description               string    `json:"description"`
	Icon                      Icon      `json:"icon"`
	Position                  int       `json:"position"`
	Created                   time.Time `json:"created"`
	Updated                   time.Time `json:"updated"`
	ObjectCount               int       `json:"objectCount"`
	ParentObjectTypeID        int       `json:"parentObjectTypeId"`
	ObjectSchemaID            int       `json:"objectSchemaId"`
	Inherited                 bool      `json:"inherited"`
	AbstractObjectType        bool      `json:"abstractObjectType"`
	ParentObjectTypeInherited bool      `json:"parentObjectTypeInherited"`
}
type ObjectAttributeValues struct {
	Value          string `json:"value"`
	DisplayValue   string `json:"displayValue"`
	SearchValue    string `json:"searchValue"`
	ReferencedType bool   `json:"referencedType"`
}
type Attributes struct {
	ID                    int                     `json:"id"`
	ObjectTypeAttribute   ObjectTypeAttribute     `json:"objectTypeAttribute,omitempty"`
	ObjectTypeAttributeID int                     `json:"objectTypeAttributeId"`
	ObjectAttributeValues []ObjectAttributeValues `json:"objectAttributeValues"`
	ObjectID              int                     `json:"objectId"`
}
type Links struct {
	Self string `json:"self"`
}
type ObjectEntry struct {
	ID         int          `json:"id"`
	Label      string       `json:"label"`
	ObjectKey  string       `json:"objectKey"`
	Avatar     Avatar       `json:"avatar"`
	ObjectType ObjectType   `json:"objectType"`
	Created    time.Time    `json:"created"`
	Updated    time.Time    `json:"updated"`
	HasAvatar  bool         `json:"hasAvatar"`
	Timestamp  int64        `json:"timestamp"`
	Attributes []Attributes `json:"attributes"`
	Links      Links        `json:"_links"`
	Name       string       `json:"name"`
}
type DefaultType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type ReferenceType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	URL16       string `json:"url16"`
	Removable   bool   `json:"removable"`
}
type ReferenceObjectType struct {
	ID                        int       `json:"id"`
	Name                      string    `json:"name"`
	Type                      int       `json:"type"`
	Description               string    `json:"description"`
	Icon                      Icon      `json:"icon"`
	Position                  int       `json:"position"`
	Created                   time.Time `json:"created"`
	Updated                   time.Time `json:"updated"`
	ObjectCount               int       `json:"objectCount"`
	ParentObjectTypeID        int       `json:"parentObjectTypeId"`
	ObjectSchemaID            int       `json:"objectSchemaId"`
	Inherited                 bool      `json:"inherited"`
	AbstractObjectType        bool      `json:"abstractObjectType"`
	ParentObjectTypeInherited bool      `json:"parentObjectTypeInherited"`
}
type ObjectTypeAttribute struct {
	ID                      int                 `json:"id"`
	Name                    string              `json:"name"`
	Label                   bool                `json:"label"`
	Type                    int                 `json:"type"`
	DefaultType             DefaultType         `json:"defaultType,omitempty"`
	Editable                bool                `json:"editable"`
	System                  bool                `json:"system"`
	Sortable                bool                `json:"sortable"`
	Summable                bool                `json:"summable"`
	Indexed                 bool                `json:"indexed"`
	MinimumCardinality      int                 `json:"minimumCardinality"`
	MaximumCardinality      int                 `json:"maximumCardinality"`
	Removable               bool                `json:"removable"`
	Hidden                  bool                `json:"hidden"`
	IncludeChildObjectTypes bool                `json:"includeChildObjectTypes"`
	UniqueAttribute         bool                `json:"uniqueAttribute"`
	Options                 string              `json:"options"`
	Position                int                 `json:"position"`
	Description             string              `json:"description,omitempty"`
	Suffix                  string              `json:"suffix,omitempty"`
	RegexValidation         string              `json:"regexValidation,omitempty"`
	QlQuery                 string              `json:"qlQuery,omitempty"`
	Iql                     string              `json:"iql,omitempty"`
	TypeValueMulti          []string            `json:"typeValueMulti,omitempty"`
	AdditionalValue         string              `json:"additionalValue,omitempty"`
	ReferenceType           ReferenceType       `json:"referenceType,omitempty"`
	ReferenceObjectTypeID   int                 `json:"referenceObjectTypeId,omitempty"`
	ReferenceObjectType     ReferenceObjectType `json:"referenceObjectType,omitempty"`
}
