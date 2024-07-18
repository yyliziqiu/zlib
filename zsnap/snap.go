package zsnap

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yyliziqiu/zlib/zfile"
)

type Snap struct {
	Path string
	Data interface{}
}

func New(path string) *Snap {
	return &Snap{
		Path: path,
	}
}

func (s *Snap) SaveData(data interface{}) error {
	err := zfile.MakeDirIfNotExist(filepath.Dir(s.Path))
	if err != nil {
		return fmt.Errorf("mkdir snapshot dir [%s] failed [%v]", filepath.Dir(s.Path), err)
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal snapshot data [%s] failed [%v]", s.Path, err)
	}

	err = os.WriteFile(s.Path+snapTempExt, bytes, 0644)
	if err != nil {
		return fmt.Errorf("save snapshot file [%s] failed [%v]", s.Path, err)
	}

	err = os.Rename(s.Path+snapTempExt, s.Path)
	if err != nil {
		return fmt.Errorf("rename snapshot file [%s] failed [%v]", s.Path, err)
	}

	return nil
}

func (s *Snap) LoadData(data interface{}) error {
	ok, err := zfile.Exist(s.Path)
	if err != nil {
		return fmt.Errorf("check snapshot file [%s] failed [%v]", s.Path, err)
	}
	if !ok {
		return nil
	}

	bytes, err := os.ReadFile(s.Path)
	if err != nil {
		return fmt.Errorf("load snapshot file [%s] failed [%v]", s.Path, err)
	}

	return json.Unmarshal(bytes, data)
}

func NewWithData(path string, data interface{}) *Snap {
	return &Snap{
		Path: path,
		Data: data,
	}
}

func (s *Snap) Save() error {
	return s.SaveData(s.Data)
}

func (s *Snap) Load() error {
	return s.LoadData(s.Data)
}
