package bigid

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sony/sonyflake"
	"github.com/spf13/cast"
)

type IDGen interface {
	MustNextID() ID
	NextID() (ID, error)
}

type idGen struct {
	*sonyflake.Sonyflake
}

func (g *idGen) MustNextID() ID {
	id, _ := g.Sonyflake.NextID()
	return ID(id)
}

func (g *idGen) NextID() (ID, error) {
	id, err := g.Sonyflake.NextID()
	return ID(id), err
}

func NewIDGen() IDGen {
	return &idGen{Sonyflake: sonyflake.NewSonyflake(sonyflake.Settings{})}
}

type ID uint64

func (id *ID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", cast.ToString(*id))), nil
}

func (id *ID) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err // 如果解析为字符串失败，返回错误
	}
	v, err := cast.ToUint64E(s)
	if err != nil {
		return err
	}
	*id = ID(v)
	return nil
}

func (id *ID) String() string {
	return strconv.FormatInt(int64(*id), 10)
}

func (id *ID) Int64() int64 {
	return int64(*id)
}

func (id *ID) Uint64() uint64 {
	return uint64(*id)
}

func (id *ID) Int() int {
	return int(*id)
}

type id interface {
	String() string
	Int64() int64
	Uint64() uint64
	Int() int
	json.Marshaler
}

var _ id = new(ID)
