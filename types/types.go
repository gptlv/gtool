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

type InsightUser struct {
	ObjectEntries []struct {
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
		} `json:"objectType,omitempty"`
		Created    time.Time `json:"created"`
		Updated    time.Time `json:"updated"`
		HasAvatar  bool      `json:"hasAvatar"`
		Timestamp  int64     `json:"timestamp"`
		Attributes []struct {
			ID                    int `json:"id"`
			ObjectTypeAttributeID int `json:"objectTypeAttributeId"`
			ObjectAttributeValues []struct {
				Value          string `json:"value"`
				DisplayValue   string `json:"displayValue"`
				SearchValue    string `json:"searchValue"`
				ReferencedType bool   `json:"referencedType"`
			} `json:"objectAttributeValues"`
			ObjectID int `json:"objectId"`
		} `json:"attributes"`
		Links struct {
			Self string `json:"self"`
		} `json:"_links"`
		Name        string `json:"name"`
		ObjectType0 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				URL16 string `json:"url16"`
				URL48 string `json:"url48"`
			} `json:"icon"`
			Position                  int       `json:"position"`
			Created                   time.Time `json:"created"`
			Updated                   time.Time `json:"updated"`
			ObjectCount               int       `json:"objectCount"`
			ObjectSchemaID            int       `json:"objectSchemaId"`
			Inherited                 bool      `json:"inherited"`
			AbstractObjectType        bool      `json:"abstractObjectType"`
			ParentObjectTypeInherited bool      `json:"parentObjectTypeInherited"`
		} `json:"objectType,omitempty"`
	} `json:"objectEntries"`
	ObjectTypeAttributes []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Label       bool   `json:"label"`
		Type        int    `json:"type"`
		DefaultType struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"defaultType,omitempty"`
		Editable                bool     `json:"editable"`
		System                  bool     `json:"system"`
		Sortable                bool     `json:"sortable"`
		Summable                bool     `json:"summable"`
		Indexed                 bool     `json:"indexed"`
		MinimumCardinality      int      `json:"minimumCardinality"`
		MaximumCardinality      int      `json:"maximumCardinality"`
		Removable               bool     `json:"removable"`
		Hidden                  bool     `json:"hidden"`
		IncludeChildObjectTypes bool     `json:"includeChildObjectTypes"`
		UniqueAttribute         bool     `json:"uniqueAttribute"`
		Options                 string   `json:"options"`
		Position                int      `json:"position"`
		Description             string   `json:"description,omitempty"`
		Suffix                  string   `json:"suffix,omitempty"`
		RegexValidation         string   `json:"regexValidation,omitempty"`
		QlQuery                 string   `json:"qlQuery,omitempty"`
		Iql                     string   `json:"iql,omitempty"`
		TypeValueMulti          []string `json:"typeValueMulti,omitempty"`
		AdditionalValue         string   `json:"additionalValue,omitempty"`
		ReferenceType           struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Color       string `json:"color"`
			URL16       string `json:"url16"`
			Removable   bool   `json:"removable"`
		} `json:"referenceType,omitempty"`
		ReferenceObjectTypeID int `json:"referenceObjectTypeId,omitempty"`
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
		} `json:"referenceObjectType,omitempty"`
		ReferenceType0 struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			Color          string `json:"color"`
			URL16          string `json:"url16"`
			Removable      bool   `json:"removable"`
			ObjectSchemaID int    `json:"objectSchemaId"`
		} `json:"referenceType,omitempty"`
		ReferenceObjectType0 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				URL16 string `json:"url16"`
				URL48 string `json:"url48"`
			} `json:"icon"`
			Position                  int       `json:"position"`
			Created                   time.Time `json:"created"`
			Updated                   time.Time `json:"updated"`
			ObjectCount               int       `json:"objectCount"`
			ObjectSchemaID            int       `json:"objectSchemaId"`
			Inherited                 bool      `json:"inherited"`
			AbstractObjectType        bool      `json:"abstractObjectType"`
			ParentObjectTypeInherited bool      `json:"parentObjectTypeInherited"`
		} `json:"referenceObjectType,omitempty"`
		ReferenceType1 struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			Color          string `json:"color"`
			URL16          string `json:"url16"`
			Removable      bool   `json:"removable"`
			ObjectSchemaID int    `json:"objectSchemaId"`
		} `json:"referenceType,omitempty"`
		ReferenceObjectType1 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				URL16 string `json:"url16"`
				URL48 string `json:"url48"`
			} `json:"icon"`
			Position                  int       `json:"position"`
			Created                   time.Time `json:"created"`
			Updated                   time.Time `json:"updated"`
			ObjectCount               int       `json:"objectCount"`
			ObjectSchemaID            int       `json:"objectSchemaId"`
			Inherited                 bool      `json:"inherited"`
			AbstractObjectType        bool      `json:"abstractObjectType"`
			ParentObjectTypeInherited bool      `json:"parentObjectTypeInherited"`
		} `json:"referenceObjectType,omitempty"`
		ReferenceType2 struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			Color          string `json:"color"`
			URL16          string `json:"url16"`
			Removable      bool   `json:"removable"`
			ObjectSchemaID int    `json:"objectSchemaId"`
		} `json:"referenceType,omitempty"`
		ReferenceObjectType2 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				URL16 string `json:"url16"`
				URL48 string `json:"url48"`
			} `json:"icon"`
			Position                  int       `json:"position"`
			Created                   time.Time `json:"created"`
			Updated                   time.Time `json:"updated"`
			ObjectCount               int       `json:"objectCount"`
			ObjectSchemaID            int       `json:"objectSchemaId"`
			Inherited                 bool      `json:"inherited"`
			AbstractObjectType        bool      `json:"abstractObjectType"`
			ParentObjectTypeInherited bool      `json:"parentObjectTypeInherited"`
		} `json:"referenceObjectType,omitempty"`
		ReferenceType3 struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			Color          string `json:"color"`
			URL16          string `json:"url16"`
			Removable      bool   `json:"removable"`
			ObjectSchemaID int    `json:"objectSchemaId"`
		} `json:"referenceType,omitempty"`
		ReferenceObjectType3 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				URL16 string `json:"url16"`
				URL48 string `json:"url48"`
			} `json:"icon"`
			Position                  int       `json:"position"`
			Created                   time.Time `json:"created"`
			Updated                   time.Time `json:"updated"`
			ObjectCount               int       `json:"objectCount"`
			ObjectSchemaID            int       `json:"objectSchemaId"`
			Inherited                 bool      `json:"inherited"`
			AbstractObjectType        bool      `json:"abstractObjectType"`
			ParentObjectTypeInherited bool      `json:"parentObjectTypeInherited"`
		} `json:"referenceObjectType,omitempty"`
		ReferenceType4 struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			Color          string `json:"color"`
			URL16          string `json:"url16"`
			Removable      bool   `json:"removable"`
			ObjectSchemaID int    `json:"objectSchemaId"`
		} `json:"referenceType,omitempty"`
		ReferenceObjectType4 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				URL16 string `json:"url16"`
				URL48 string `json:"url48"`
			} `json:"icon"`
			Position                  int       `json:"position"`
			Created                   time.Time `json:"created"`
			Updated                   time.Time `json:"updated"`
			ObjectCount               int       `json:"objectCount"`
			ObjectSchemaID            int       `json:"objectSchemaId"`
			Inherited                 bool      `json:"inherited"`
			AbstractObjectType        bool      `json:"abstractObjectType"`
			ParentObjectTypeInherited bool      `json:"parentObjectTypeInherited"`
		} `json:"referenceObjectType,omitempty"`
		ReferenceObjectType5 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
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
		} `json:"referenceObjectType,omitempty"`
		ReferenceObjectType6 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
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
		} `json:"referenceObjectType,omitempty"`
		ReferenceType5 struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			Color          string `json:"color"`
			URL16          string `json:"url16"`
			Removable      bool   `json:"removable"`
			ObjectSchemaID int    `json:"objectSchemaId"`
		} `json:"referenceType,omitempty"`
		ReferenceObjectType7 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				URL16 string `json:"url16"`
				URL48 string `json:"url48"`
			} `json:"icon"`
			Position                  int       `json:"position"`
			Created                   time.Time `json:"created"`
			Updated                   time.Time `json:"updated"`
			ObjectCount               int       `json:"objectCount"`
			ObjectSchemaID            int       `json:"objectSchemaId"`
			Inherited                 bool      `json:"inherited"`
			AbstractObjectType        bool      `json:"abstractObjectType"`
			ParentObjectTypeInherited bool      `json:"parentObjectTypeInherited"`
		} `json:"referenceObjectType,omitempty"`
		ReferenceObjectType8 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
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
		} `json:"referenceObjectType,omitempty"`
		ReferenceObjectType9 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
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
		} `json:"referenceObjectType,omitempty"`
		ReferenceObjectType10 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
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
		} `json:"referenceObjectType,omitempty"`
		ReferenceObjectType11 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
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
		} `json:"referenceObjectType,omitempty"`
		ReferenceObjectType12 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
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
		} `json:"referenceObjectType,omitempty"`
		ReferenceObjectType13 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
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
		} `json:"referenceObjectType,omitempty"`
		ReferenceObjectType14 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
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
		} `json:"referenceObjectType,omitempty"`
		ReferenceType6 struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			Color          string `json:"color"`
			URL16          string `json:"url16"`
			Removable      bool   `json:"removable"`
			ObjectSchemaID int    `json:"objectSchemaId"`
		} `json:"referenceType,omitempty"`
		ReferenceObjectType15 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				URL16 string `json:"url16"`
				URL48 string `json:"url48"`
			} `json:"icon"`
			Position                  int       `json:"position"`
			Created                   time.Time `json:"created"`
			Updated                   time.Time `json:"updated"`
			ObjectCount               int       `json:"objectCount"`
			ObjectSchemaID            int       `json:"objectSchemaId"`
			Inherited                 bool      `json:"inherited"`
			AbstractObjectType        bool      `json:"abstractObjectType"`
			ParentObjectTypeInherited bool      `json:"parentObjectTypeInherited"`
		} `json:"referenceObjectType,omitempty"`
		ReferenceType7 struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			Color          string `json:"color"`
			URL16          string `json:"url16"`
			Removable      bool   `json:"removable"`
			ObjectSchemaID int    `json:"objectSchemaId"`
		} `json:"referenceType,omitempty"`
		ReferenceObjectType16 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				URL16 string `json:"url16"`
				URL48 string `json:"url48"`
			} `json:"icon"`
			Position                  int       `json:"position"`
			Created                   time.Time `json:"created"`
			Updated                   time.Time `json:"updated"`
			ObjectCount               int       `json:"objectCount"`
			ObjectSchemaID            int       `json:"objectSchemaId"`
			Inherited                 bool      `json:"inherited"`
			AbstractObjectType        bool      `json:"abstractObjectType"`
			ParentObjectTypeInherited bool      `json:"parentObjectTypeInherited"`
		} `json:"referenceObjectType,omitempty"`
		ReferenceObjectType17 struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type int    `json:"type"`
			Icon struct {
				ID    int    `json:"id"`
				Name  string `json:"name"`
				URL16 string `json:"url16"`
				URL48 string `json:"url48"`
			} `json:"icon"`
			Position                  int       `json:"position"`
			Created                   time.Time `json:"created"`
			Updated                   time.Time `json:"updated"`
			ObjectCount               int       `json:"objectCount"`
			ObjectSchemaID            int       `json:"objectSchemaId"`
			Inherited                 bool      `json:"inherited"`
			AbstractObjectType        bool      `json:"abstractObjectType"`
			ParentObjectTypeInherited bool      `json:"parentObjectTypeInherited"`
		} `json:"referenceObjectType,omitempty"`
	} `json:"objectTypeAttributes"`
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
