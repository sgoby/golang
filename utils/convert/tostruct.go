package convert

import (
	"fmt"
	"reflect"
)

var specialType []reflect.Kind = []reflect.Kind{
	reflect.Array, reflect.Func, reflect.Map, reflect.Slice,
	reflect.Struct, reflect.UnsafePointer, reflect.Ptr, reflect.Chan,
	reflect.Complex64, reflect.Complex128, reflect.Interface,
}
//
func InterfaceToStruct(vmap interface{}, mstuct interface{}) (err error) {
	defer func() {
		if panidErr := recover(); panidErr != nil {
			err = fmt.Errorf("%v", panidErr)
		}
	}()
	if vmap == nil || reflect.TypeOf(vmap) == nil {
		return fmt.Errorf("The nil value can not InterfaceToStruct....")
	}
	//
	vmapValue := reflect.ValueOf(vmap)
	if vmapValue.Kind() != reflect.Map {
		return fmt.Errorf("the source interface type is not Map")
	}
	keys := vmapValue.MapKeys()
	if keys[0].Kind() != reflect.String {
		return fmt.Errorf("the map key type must string")
	}
	newValue := reflect.ValueOf(mstuct)
	if newValue.Kind() != reflect.Ptr || newValue.Pointer() == 0 {
		return fmt.Errorf("Target stuct must be a pointer")
	}
	return valueToValue(vmapValue, newValue)
}

/*
value 复制,带基础类型自动转化功能
*/
func valueToValue(vmapValue reflect.Value, mstuct reflect.Value) (err error) {
	vmapValue = getElemValue(vmapValue)
	mstuct = getElemValue(mstuct)
	//
	if vmapValue.Kind() == mstuct.Kind() && isBaseType(mstuct.Kind()) {
		mstuct.Set(vmapValue)
	} else if vmapValue.Kind() == mstuct.Kind() && !isBaseType(mstuct.Kind()) {
		//此处可以跟据需求，自行扩展
		if mstuct.Kind() == reflect.Slice || mstuct.Kind() == reflect.Array {
			//
			elemType := mstuct.Type().Elem()
			mstuctSliceV := reflect.MakeSlice(mstuct.Type(), 0, 0)
			//
			for i := 0; i < vmapValue.Len(); i++ {
				t := reflect.New(elemType).Elem()
				m := vmapValue.Index(i)
				if err := valueToValue(m, t); err != nil {
					return err
				}
				mstuctSliceV = reflect.Append(mstuctSliceV, t)
			}
			mstuct.Set(mstuctSliceV)
			return nil
		}
		//=========================以下部分是类型不相等的===================================
	} else if isBaseType(vmapValue.Kind()) && isBaseType(mstuct.Kind()) {
		nVal, err := valueConvert(vmapValue, mstuct.Type())
		if err != nil {
			return err
		}
		mstuct.Set(nVal)
		//
	} else if vmapValue.Kind() == reflect.Interface {
		i := vmapValue.Interface()
		v := reflect.ValueOf(i)
		//
		if isBaseType(v.Kind()) {
			nVal, err := valueConvert(v, mstuct.Type())
			if err != nil {
				return err
			}
			mstuct.Set(nVal)
		} else {
			if err := valueToValue(v, mstuct); err != nil {
				return err
			}
		}
	} else if vmapValue.Kind() == reflect.Map && mstuct.Kind() == reflect.Struct {
		slicSF := getStructFields(mstuct)
		for _, sf := range slicSF {
			mapKeyValue := vmapValue.MapIndex(reflect.ValueOf(sf.Name)) //.Elem();
			mapKeyValue = getElemValue(mapKeyValue)
			if !mapKeyValue.IsValid() {
				continue
			}
			fieldValue := mstuct.FieldByName(sf.Name)
			if fieldValue.IsValid() && fieldValue.CanSet() {
				if err := valueToValue(mapKeyValue, fieldValue); err != nil {
					return err
				}
			}
		}
	} else if mstuct.IsNil() {
		mstuct.Set(reflect.New(mstuct.Type().Elem()))
		if err := valueToValue(vmapValue, mstuct); err != nil {
			return err
		}
	} else {
		//此处可以跟据需求，自行扩展
		return fmt.Errorf("Not support %s  Convert to %s ", vmapValue.Kind(), mstuct.Kind())
	}
	return nil
}

//
func valueConvert(v reflect.Value, t reflect.Type) (newVal reflect.Value, err error) {
	defer func() {
		if panidErr := recover(); panidErr != nil {
			err = fmt.Errorf("'%v' %v", v, panidErr)
		}
	}()
	return v.Convert(t), nil
}

/*
判断是否为基础类型,跟据需要可自行扩展
*/
func isBaseType(k reflect.Kind) bool {
	for _, vk := range specialType {
		if k == vk {
			return false
		}
	}
	return true
}

/*
找到有效value
*/
func getElemValue(v reflect.Value) reflect.Value {
	for { //找到内容
		if v.Kind() != reflect.Ptr || v.IsNil() {
			break
		}
		v = v.Elem()
	}
	return v
}

/*
获取Struct 的所有字段
*/
func getStructFields(val reflect.Value) []reflect.StructField {
	pType := val.Type()
	if pType == nil {
		return nil
	}
	for { //找到内容
		if pType.Kind() != reflect.Ptr {
			break
		}
		pType = pType.Elem()
	}
	//
	if pType.NumField() <= 0 {
		return nil
	}
	var slicStructField []reflect.StructField
	for i := 0; i < pType.NumField(); i++ {
		slicStructField = append(slicStructField, pType.Field(i))
	}
	return slicStructField
}
