package cache

import (
	"encoding/json"
	"fmt"
)

type RedisStorageAdapter struct {
	UserBanner map[string]any
	Bytes      []byte
	FeatureId  int
	TagId      int
}

func (adapter RedisStorageAdapter) Key() string {
	return fmt.Sprintf("feature: %d, tag: %d", adapter.FeatureId, adapter.TagId)
}

func (adapter RedisStorageAdapter) Value() ([]byte, error) {
	b, err := json.Marshal(&adapter.UserBanner)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (adapter RedisStorageAdapter) Content() (map[string]any, error) {
	err := json.Unmarshal(adapter.Bytes, &adapter.UserBanner)
	if err != nil {
		return nil, err
	}
	return adapter.UserBanner, nil
}
