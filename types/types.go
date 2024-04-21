package types

import "time"

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

type InsightObjectEntries struct {
	ObjectEntries         []InsightObject
	ObjectTypeID          int    `json:"objectTypeId"`
	ObjectTypeIsInherited bool   `json:"objectTypeIsInherited"`
	AbstractObjectType    bool   `json:"abstractObjectType"`
	TotalFilterCount      int    `json:"totalFilterCount"`
	StartIndex            int    `json:"startIndex"`
	ToIndex               int    `json:"toIndex"`
	PageObjectSize        int    `json:"pageObjectSize"`
	PageNumber            int    `json:"pageNumber"`
	OrderWay              string `json:"orderWay"`
	QlQuery               string `json:"qlQuery"`
	QlQuerySearchResult   bool   `json:"qlQuerySearchResult"`
	ConversionPossible    bool   `json:"conversionPossible"`
	Iql                   string `json:"iql"`
	IqlSearchResult       bool   `json:"iqlSearchResult"`
	PageSize              int    `json:"pageSize"`
}

type InsightObject struct {
	ID        int    `json:"id"`
	Label     string `json:"label"`
	ObjectKey string `json:"objectKey"`
	Avatar    struct {
		URL16    string `json:"url16"`
		URL48    string `json:"url48"`
		URL72    string `json:"url72"`
		URL144   string `json:"url144"`
		URL288   string `json:"url288"`
		ObjectID int    `json:"objectId"`
	} `json:"avatar"`
	ObjectType struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Type        int    `json:"type"`
		Description string `json:"description"`
		Icon        struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			URL16 string `json:"url16"`
			URL48 string `json:"url48"`
		} `json:"icon"`
		Position                  int       `json:"position"`
		Created                   time.Time `json:"created"`
		Updated                   time.Time `json:"updated"`
		ObjectCount               int       `json:"objectCount"`
		ParentObjectTypeID        int       `json:"parentObjectTypeId"`
		ObjectSchemaID            int       `json:"objectSchemaId"`
		Inherited                 bool      `json:"inherited"`
		AbstractObjectType        bool      `json:"abstractObjectType"`
		ParentObjectTypeInherited bool      `json:"parentObjectTypeInherited"`
	} `json:"objectType"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
	HasAvatar  bool      `json:"hasAvatar"`
	Timestamp  int64     `json:"timestamp"`
	Attributes []struct {
		ID                  int `json:"id"`
		ObjectTypeAttribute struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Label       bool   `json:"label"`
			Type        int    `json:"type"`
			DefaultType struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"defaultType"`
			Editable                bool   `json:"editable"`
			System                  bool   `json:"system"`
			Sortable                bool   `json:"sortable"`
			Summable                bool   `json:"summable"`
			Indexed                 bool   `json:"indexed"`
			MinimumCardinality      int    `json:"minimumCardinality"`
			MaximumCardinality      int    `json:"maximumCardinality"`
			Removable               bool   `json:"removable"`
			Hidden                  bool   `json:"hidden"`
			IncludeChildObjectTypes bool   `json:"includeChildObjectTypes"`
			UniqueAttribute         bool   `json:"uniqueAttribute"`
			Options                 string `json:"options"`
			Position                int    `json:"position"`
		} `json:"objectTypeAttribute,omitempty"`
		ObjectTypeAttributeID int `json:"objectTypeAttributeId"`
		ObjectAttributeValues []struct {
			Value          string `json:"value"`
			DisplayValue   string `json:"displayValue"`
			SearchValue    string `json:"searchValue"`
			ReferencedType bool   `json:"referencedType"`
		} `json:"objectAttributeValues"`
		ObjectID             int `json:"objectId"`
		ObjectTypeAttribute0 struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Label       bool   `json:"label"`
			Type        int    `json:"type"`
			Description string `json:"description"`
			DefaultType struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"defaultType"`
			Editable                bool   `json:"editable"`
			System                  bool   `json:"system"`
			Sortable                bool   `json:"sortable"`
			Summable                bool   `json:"summable"`
			Indexed                 bool   `json:"indexed"`
			MinimumCardinality      int    `json:"minimumCardinality"`
			MaximumCardinality      int    `json:"maximumCardinality"`
			Suffix                  string `json:"suffix"`
			Removable               bool   `json:"removable"`
			Hidden                  bool   `json:"hidden"`
			IncludeChildObjectTypes bool   `json:"includeChildObjectTypes"`
			UniqueAttribute         bool   `json:"uniqueAttribute"`
			RegexValidation         string `json:"regexValidation"`
			QlQuery                 string `json:"qlQuery"`
			Options                 string `json:"options"`
			Position                int    `json:"position"`
			Iql                     string `json:"iql"`
		} `json:"objectTypeAttribute,omitempty"`
		ObjectTypeAttribute1 struct {
			ID                      int      `json:"id"`
			Name                    string   `json:"name"`
			Label                   bool     `json:"label"`
			Type                    int      `json:"type"`
			TypeValueMulti          []string `json:"typeValueMulti"`
			Editable                bool     `json:"editable"`
			System                  bool     `json:"system"`
			Sortable                bool     `json:"sortable"`
			Summable                bool     `json:"summable"`
			Indexed                 bool     `json:"indexed"`
			MinimumCardinality      int      `json:"minimumCardinality"`
			MaximumCardinality      int      `json:"maximumCardinality"`
			Suffix                  string   `json:"suffix"`
			Removable               bool     `json:"removable"`
			Hidden                  bool     `json:"hidden"`
			IncludeChildObjectTypes bool     `json:"includeChildObjectTypes"`
			UniqueAttribute         bool     `json:"uniqueAttribute"`
			RegexValidation         string   `json:"regexValidation"`
			QlQuery                 string   `json:"qlQuery"`
			Options                 string   `json:"options"`
			Position                int      `json:"position"`
			Iql                     string   `json:"iql"`
		} `json:"objectTypeAttribute,omitempty"`
		ObjectTypeAttribute2 struct {
			ID                      int    `json:"id"`
			Name                    string `json:"name"`
			Label                   bool   `json:"label"`
			Type                    int    `json:"type"`
			Description             string `json:"description"`
			AdditionalValue         string `json:"additionalValue"`
			Editable                bool   `json:"editable"`
			System                  bool   `json:"system"`
			Sortable                bool   `json:"sortable"`
			Summable                bool   `json:"summable"`
			Indexed                 bool   `json:"indexed"`
			MinimumCardinality      int    `json:"minimumCardinality"`
			MaximumCardinality      int    `json:"maximumCardinality"`
			Removable               bool   `json:"removable"`
			Hidden                  bool   `json:"hidden"`
			IncludeChildObjectTypes bool   `json:"includeChildObjectTypes"`
			UniqueAttribute         bool   `json:"uniqueAttribute"`
			Options                 string `json:"options"`
			Position                int    `json:"position"`
		} `json:"objectTypeAttribute,omitempty"`
		ObjectTypeAttribute3 struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Label       bool   `json:"label"`
			Type        int    `json:"type"`
			Description string `json:"description"`
			DefaultType struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"defaultType"`
			Editable                bool   `json:"editable"`
			System                  bool   `json:"system"`
			Sortable                bool   `json:"sortable"`
			Summable                bool   `json:"summable"`
			Indexed                 bool   `json:"indexed"`
			MinimumCardinality      int    `json:"minimumCardinality"`
			MaximumCardinality      int    `json:"maximumCardinality"`
			Suffix                  string `json:"suffix"`
			Removable               bool   `json:"removable"`
			Hidden                  bool   `json:"hidden"`
			IncludeChildObjectTypes bool   `json:"includeChildObjectTypes"`
			UniqueAttribute         bool   `json:"uniqueAttribute"`
			RegexValidation         string `json:"regexValidation"`
			QlQuery                 string `json:"qlQuery"`
			Options                 string `json:"options"`
			Position                int    `json:"position"`
			Iql                     string `json:"iql"`
		} `json:"objectTypeAttribute,omitempty"`
		ObjectTypeAttribute4 struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Label       bool   `json:"label"`
			Type        int    `json:"type"`
			Description string `json:"description"`
			DefaultType struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"defaultType"`
			Editable                bool   `json:"editable"`
			System                  bool   `json:"system"`
			Sortable                bool   `json:"sortable"`
			Summable                bool   `json:"summable"`
			Indexed                 bool   `json:"indexed"`
			MinimumCardinality      int    `json:"minimumCardinality"`
			MaximumCardinality      int    `json:"maximumCardinality"`
			Suffix                  string `json:"suffix"`
			Removable               bool   `json:"removable"`
			Hidden                  bool   `json:"hidden"`
			IncludeChildObjectTypes bool   `json:"includeChildObjectTypes"`
			UniqueAttribute         bool   `json:"uniqueAttribute"`
			RegexValidation         string `json:"regexValidation"`
			QlQuery                 string `json:"qlQuery"`
			Options                 string `json:"options"`
			Position                int    `json:"position"`
			Iql                     string `json:"iql"`
		} `json:"objectTypeAttribute,omitempty"`
		ObjectTypeAttribute5 struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Label       bool   `json:"label"`
			Type        int    `json:"type"`
			Description string `json:"description"`
			DefaultType struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"defaultType"`
			Editable                bool   `json:"editable"`
			System                  bool   `json:"system"`
			Sortable                bool   `json:"sortable"`
			Summable                bool   `json:"summable"`
			Indexed                 bool   `json:"indexed"`
			MinimumCardinality      int    `json:"minimumCardinality"`
			MaximumCardinality      int    `json:"maximumCardinality"`
			Suffix                  string `json:"suffix"`
			Removable               bool   `json:"removable"`
			Hidden                  bool   `json:"hidden"`
			IncludeChildObjectTypes bool   `json:"includeChildObjectTypes"`
			UniqueAttribute         bool   `json:"uniqueAttribute"`
			RegexValidation         string `json:"regexValidation"`
			QlQuery                 string `json:"qlQuery"`
			Options                 string `json:"options"`
			Position                int    `json:"position"`
			Iql                     string `json:"iql"`
		} `json:"objectTypeAttribute,omitempty"`
		ObjectTypeAttribute6 struct {
			ID            int    `json:"id"`
			Name          string `json:"name"`
			Label         bool   `json:"label"`
			Type          int    `json:"type"`
			ReferenceType struct {
				ID          int    `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Color       string `json:"color"`
				URL16       string `json:"url16"`
				Removable   bool   `json:"removable"`
			} `json:"referenceType"`
			ReferenceObjectTypeID int `json:"referenceObjectTypeId"`
			ReferenceObjectType   struct {
				ID          int    `json:"id"`
				Name        string `json:"name"`
				Type        int    `json:"type"`
				Description string `json:"description"`
				Icon        struct {
					ID    int    `json:"id"`
					Name  string `json:"name"`
					URL16 string `json:"url16"`
					URL48 string `json:"url48"`
				} `json:"icon"`
				Position                  int       `json:"position"`
				Created                   time.Time `json:"created"`
				Updated                   time.Time `json:"updated"`
				ObjectCount               int       `json:"objectCount"`
				ParentObjectTypeID        int       `json:"parentObjectTypeId"`
				ObjectSchemaID            int       `json:"objectSchemaId"`
				Inherited                 bool      `json:"inherited"`
				AbstractObjectType        bool      `json:"abstractObjectType"`
				ParentObjectTypeInherited bool      `json:"parentObjectTypeInherited"`
			} `json:"referenceObjectType"`
			Editable                bool   `json:"editable"`
			System                  bool   `json:"system"`
			Sortable                bool   `json:"sortable"`
			Summable                bool   `json:"summable"`
			Indexed                 bool   `json:"indexed"`
			MinimumCardinality      int    `json:"minimumCardinality"`
			MaximumCardinality      int    `json:"maximumCardinality"`
			Removable               bool   `json:"removable"`
			Hidden                  bool   `json:"hidden"`
			IncludeChildObjectTypes bool   `json:"includeChildObjectTypes"`
			UniqueAttribute         bool   `json:"uniqueAttribute"`
			Options                 string `json:"options"`
			Position                int    `json:"position"`
		} `json:"objectTypeAttribute,omitempty"`
		ObjectTypeAttribute7 struct {
			ID                      int    `json:"id"`
			Name                    string `json:"name"`
			Label                   bool   `json:"label"`
			Type                    int    `json:"type"`
			Description             string `json:"description"`
			AdditionalValue         string `json:"additionalValue"`
			Editable                bool   `json:"editable"`
			System                  bool   `json:"system"`
			Sortable                bool   `json:"sortable"`
			Summable                bool   `json:"summable"`
			Indexed                 bool   `json:"indexed"`
			MinimumCardinality      int    `json:"minimumCardinality"`
			MaximumCardinality      int    `json:"maximumCardinality"`
			Suffix                  string `json:"suffix"`
			Removable               bool   `json:"removable"`
			Hidden                  bool   `json:"hidden"`
			IncludeChildObjectTypes bool   `json:"includeChildObjectTypes"`
			UniqueAttribute         bool   `json:"uniqueAttribute"`
			RegexValidation         string `json:"regexValidation"`
			QlQuery                 string `json:"qlQuery"`
			Options                 string `json:"options"`
			Position                int    `json:"position"`
			Iql                     string `json:"iql"`
		} `json:"objectTypeAttribute,omitempty"`
	} `json:"attributes"`
	ExtendedInfo struct {
		OpenIssuesExists  bool `json:"openIssuesExists"`
		AttachmentsExists bool `json:"attachmentsExists"`
	} `json:"extendedInfo"`
	Links struct {
		Self string `json:"self"`
	} `json:"_links"`
	Name string `json:"name"`
}
