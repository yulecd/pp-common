package app

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/yulecd/pp-common/util"
)

const OmitRequiredFlag = "OmitRequired"

type OmitRequired struct {
	Key   string
	Error error
}

type Validator interface {
	IsSatisfied(obj interface{}) bool
	GetError() error
	HasValidatorFlag(tagStr string, vfunc *string) bool
}

type Validation struct{}

func (v *Validation) Check(obj interface{}, checks ...Validator) (bool, error) {
	for _, check := range checks {
		if !check.IsSatisfied(obj) {
			return false, check.GetError()
		}
	}

	return true, nil
}

func (or *OmitRequired) IsSatisfied(obj interface{}) bool {
	if obj == nil {
		or.Error = errors.New("invalid check obj, parsed nil")
		return false
	}

	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	if util.IsStructPtr(objT) {
		objT = objT.Elem()
		objV = objV.Elem()
	}
	if !util.IsStruct(objT) {
		return true
	}

	group, errMap := or.makeGroup(objT, objV)
	if group == nil && or.Error != nil {
		return false
	}

	if len(group) > 0 {
		var isSatisfied bool
		var errStr string
		for k, v := range group {
			if v > 0 {
				isSatisfied = true
				break
			}

			errStr += fmt.Sprintf("%s at latest need one;", errMap[k])
		}
		if errStr != "" {
			or.Error = fmt.Errorf(errStr)
		}
		if !isSatisfied {
			return false
		}
	}

	return true
}

func mergeMap(dest, src interface{}) error {
	if reflect.TypeOf(dest).Kind() != reflect.TypeOf(src).Kind() {
		return errors.New("un support different type merge")
	}

	switch src.(type) {
	case map[string]int:
		for k, v := range src.(map[string]int) {
			dest.(map[string]int)[k] = v
		}
	case map[string]string:
		for k, v := range src.(map[string]string) {
			dest.(map[string]string)[k] = v
		}
	case map[string]interface{}:
		for k, v := range src.(map[string]interface{}) {
			dest.(map[string]interface{})[k] = v
		}
	default:
		return errors.New("un support map type to merge")
	}

	return nil
}

func (or *OmitRequired) makeGroup(objT reflect.Type, objV reflect.Value) (map[string]int, map[string]string) {
	group := make(map[string]int)
	errMap := make(map[string]string)
	for i := 0; i < objT.NumField(); i++ {
		if util.IsStruct(objT.Field(i).Type) {
			subGroup, subErrMap := or.makeGroup(objT.Field(i).Type, objV.Field(i))
			err := mergeMap(group, subGroup)
			if err != nil {
				or.Error = err
				return nil, nil
			}
			err = mergeMap(errMap, subErrMap)
			if err != nil {
				or.Error = err
				return nil, nil
			}
			continue
		}
		tag := objT.Field(i).Tag
		tv := tag.Get("xvalid")
		var vfunc string
		if tv == "" || !or.HasValidatorFlag(tv, &vfunc) {
			continue
		}

		var found bool
		gname, err := parseGroup(vfunc)
		if err != nil {
			or.Error = err
			return nil, nil
		}
		if _, ok := group[gname]; ok {
			found = true
			if !objV.Field(i).IsZero() {
				group[gname]++
			}
		}

		if !found {
			errMap[gname] = fmt.Sprintf("%s not satisfied required", gname)
			if _, ok := group[gname]; !ok {
				group[gname] = 0
			}
			if !objV.Field(i).IsZero() {
				group[gname]++
			}
		}
	}

	return group, errMap
}

func (or *OmitRequired) GetError() error {
	return or.Error
}

func (or *OmitRequired) HasValidatorFlag(tagVal string, vfunc *string) bool {
	parts := strings.Split(tagVal, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, OmitRequiredFlag) {
			*vfunc = part
			return true
		}
	}

	return false
}

func parseGroup(vfunc string) (string, error) {
	vfunc = strings.TrimSpace(vfunc)
	start := strings.Index(vfunc, "(")

	if start == -1 {
		err := fmt.Errorf("%s require at latest 1 parameters", vfunc)
		return "", err
	}

	end := strings.Index(vfunc, ")")
	if end == -1 {
		err := fmt.Errorf("invalid valid function")
		return "", err
	}

	group := vfunc[start+1 : end]

	return group, nil
}

var (
	xvalidFun = []Validator{
		&OmitRequired{},
	}
)

func RegisterXValidator(v ...Validator) {
	xvalidFun = append(xvalidFun, v...)
}
