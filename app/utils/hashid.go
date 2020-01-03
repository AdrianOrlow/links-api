package utils

import (
	"errors"
	"github.com/AdrianOrlow/links-api/config"
	"github.com/speps/go-hashids"
)

func InitializeHashId(config *config.Config) error {
	hd := hashids.NewData()
	hd.Salt = config.HashID.Salt
	hd.MinLength = config.HashID.MinLength
	hid, err := hashids.NewWithData(hd)
	utils.hashID = hid
	return err
}

func EncodeId(id int) string {
	e, _ := utils.hashID.Encode([]int{id})
	return e
}

func DecodeId(hashId string) (int, error) {
	d, _ := utils.hashID.DecodeWithError(hashId)
	if len(d) != 1 {
		return 0, errors.New("not valid hashId")
	} else {
		return d[0], nil
	}
}