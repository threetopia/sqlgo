package sqlgo

import (
	"fmt"
)

type SQLGoParameter interface {
	SetSQLParameter(value interface{}) SQLGoParameter
	GetSQLParameterCount(value interface{}) int
	GetSQLParameterSign(value interface{}) string
	GetSQLParameter() []interface{}
}

type sqlGoParameter struct {
	parameterList []interface{}
	parameterMap  map[string]int
}

func NewSQLGoParameter() SQLGoParameter {
	return new(sqlGoParameter)
}

func (s *sqlGoParameter) GetSQLGoParameter() SQLGoParameter {
	return s
}

func (s *sqlGoParameter) SetSQLParameter(value interface{}) SQLGoParameter {
	if value == nil {
		return s
	}

	hashVal := hash(value, fmt.Sprintf("%T", value))
	if s.parameterMap == nil {
		s.parameterMap = make(map[string]int)
	} else if _, ok := s.parameterMap[hashVal]; ok {
		return s
	}
	s.parameterList = append(s.parameterList, value)
	s.parameterMap[hashVal] = len(s.parameterList)
	return s
}

func (s *sqlGoParameter) GetSQLParameterCount(value interface{}) int {
	hashVal := hash(value, fmt.Sprintf("%T", value))
	if count, ok := s.parameterMap[hashVal]; ok {
		return count
	}
	return -1
}

func (s *sqlGoParameter) GetSQLParameterSign(value interface{}) string {
	return fmt.Sprintf("$%d", s.GetSQLParameterCount(value))
}

func (s *sqlGoParameter) GetSQLParameter() []interface{} {
	return s.parameterList
}
