package common

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"strconv"
	"strings"
)

type UID struct {
	localID    uint32
	objectType int
	shardID    uint32
}

func NewUID(localID uint32, objType int, shardId uint32) UID {
	return UID{
		localID:    localID,
		objectType: objType,
		shardID:    shardId,
	}
}

func (uid UID) String() string {
	val := uint64(uid.localID)<<28 | uint64(uid.objectType)<<18 | uint64(uid.shardID)<<0
	return base58.Encode([]byte(fmt.Sprintf("%v", val)))
}

func (uid UID) GetLocalID() uint32 {
	return uid.localID
}

func (uid UID) GetShardID() uint32 {
	return uid.shardID
}

func (uid UID) GetObjectType() int {
	return uid.objectType
}

func DecomposeUID(s string) (UID, error) {
	uid, err := strconv.ParseUint(s, 10, 64)

	if err != nil {
		return UID{}, err
	}

	if (1 << 28) > uid {
		return UID{}, errors.New("wrong uid.")
	}

	u := UID{
		localID:    uint32(uid >> 28),
		objectType: int(uid >> 18 & 0x3FF),
		shardID:    uint32(uid >> 18 & 0x3FF),
	}

	return u, nil
}

func FromBase58(s string) (UID, error) {
	return DecomposeUID(string(base58.Decode(s)))
}

func (uid UID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", uid.String())), nil
}

func (uid *UID) UnmarshalJSON(data []byte) error {
	decodeUID, err := FromBase58(strings.Replace(string(data), "\"", "", -1))

	if err != nil {
		return err
	}

	uid.localID = decodeUID.localID
	uid.shardID = decodeUID.shardID
	uid.objectType = decodeUID.objectType

	return nil
}

func (uid *UID) Value() (driver.Value, error) {
	if uid == nil {
		return nil, nil
	}
	return int64(uid.localID), nil
}
