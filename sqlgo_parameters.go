package sqlgo

type SQLGoParameter interface {
	GetSQLGoParameter() SQLGoParameter
	SetSQLParameter(value interface{}) SQLGoParameter
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
	if s.parameterMap == nil {
		s.parameterMap = make(map[interface{}]int)
	} else if _, ok := s.parameterMap[value]; ok {
		return s
	}

	s.parameterList = append(s.parameterList, value)
	s.parameterMap[value] = len(s.parameterList)
	return s
}
