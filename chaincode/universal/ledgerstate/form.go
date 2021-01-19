package ledgerstate

// private data input form, we use transient data
type TransientForm struct {
	Id           string `json:"id"`
	PublicDataId string `json:"publicDataId"`

	CollectionName string `json:"collectionName"`
	Data           string `json:"data"`
}

// public state create form
type PublicStateForm struct {
	Id      string `json:"id"`
	AppName string `json:"appName"`
	Data    string `json:"data"`

	// private meta info
	PrivateMetaInfo *PrivateMetaInfo `json:"privateMetaInfo"`
}
