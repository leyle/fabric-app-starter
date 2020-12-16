package main

import (
	"fmt"
)

func getId(app, dataId string) string {
	return fmt.Sprintf("%s|%s", app, dataId)
}

func copyToStorageOut(in *StorageIn) *StorageOut {
	if in == nil {
		return nil
	}
	out := &StorageOut{
		Id:        in.Id,
		App:       in.App,
		DataId:    in.DataId,
		Data:      string(in.Data),
		CreatedAt: in.CreatedAt,
	}

	return out
}
