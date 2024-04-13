package openapi

// VersionResponse struct to store banner version
type VersionResponse struct {
	Content   *map[string]interface{} `json:"content"`
	TagIds    *[]int                  `json:"tag_ids"`
	FeatureId *int                    `json:"feature_id"`
	Version   *int                    `json:"version"`
}
