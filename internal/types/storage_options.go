package types

import (
	"github.com/matzefriedrich/containerssh-authserver/internal/utils"
)

type StorageOptions struct {
	values map[string]string
}

func NewStorageOptions() *StorageOptions {
	return &StorageOptions{
		values: make(map[string]string),
	}
}

func (opt *StorageOptions) AddOrUpdate(other map[string]string) *StorageOptions {
	merged := utils.MergeMaps(opt.values, other, utils.AddMissing|utils.Override)
	opt.values = merged
	return opt
}

func (opt *StorageOptions) Size(value Quota) *StorageOptions {
	opt.values[("size")] = value.String()
	return opt
}

func (opt *StorageOptions) AsMap() map[string]string {
	return opt.values
}
