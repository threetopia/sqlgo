package sqlgo

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
)

type SQLGoParameter interface {
	GetSQLGoParameter() SQLGoParameter
	SetSQLParameter(value interface{}) SQLGoParameter
	GetSQLParameterCount(value interface{}) int
	GetSQLParameter() []interface{}
}

type sqlGoParameter struct {
	parameterList []interface{}
	parameterMap  map[interface{}]int
}

func NewSQLGoParameter() SQLGoParameter {
	return &sqlGoParameter{}
}

func (s *sqlGoParameter) GetSQLGoParameter() SQLGoParameter {
	return s
}

func (s *sqlGoParameter) SetSQLParameter(value interface{}) SQLGoParameter {
	if value == nil {
		return s
	}

	hashVal := hash(value)
	if s.parameterMap == nil {
		s.parameterMap = make(map[interface{}]int)
	} else if _, ok := s.parameterMap[hashVal]; ok {
		return s
	}

	s.parameterList = append(s.parameterList, value)
	s.parameterMap[hashVal] = len(s.parameterList)
	return s
}

func (s *sqlGoParameter) GetSQLParameterCount(value interface{}) int {
	hashVal := hash(value)
	if count, ok := s.parameterMap[hashVal]; ok {
		return count
	}
	return -1
}

func (s *sqlGoParameter) GetSQLParameter() []interface{} {
	return s.parameterList
}

func hash(i interface{}) string {
	jVal, _ := json.Marshal(i)
	h := sha1.New()
	h.Write(jVal)
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}
