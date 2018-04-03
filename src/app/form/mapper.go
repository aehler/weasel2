package form

import (
	"app/crypto"
	"reflect"
	"strings"
	"strconv"
	"errors"
	"fmt"
)

//Takes any linear struct, searches for tag `weaselform:"elementType"` and tries to create appropriate elements in Form.Elements
func (f *Form) MapStruct(s interface {}) error {

	st := reflect.TypeOf(s)

	if st.Kind() == reflect.Ptr {

		st = st.Elem()

	}

	if st.Kind() != reflect.Struct {

		return errors.New(fmt.Sprintf("Form MapStruct recieved %s, but needs Struct", st.Kind().String()))
	}

	valSt := reflect.ValueOf(s)

	for i := 0; i < st.NumField(); i++ {

		field := st.Field(i)

		valField := valSt.Field(i)

		if field.Tag.Get("weaselform") == "" {

			continue
		}

		if f.skipFields[field.Name] != nil {

			continue
		}

		if elementType[field.Tag.Get("weaselform")] == 0 {

			return errors.New(fmt.Sprintf("Inappropriate element type %s", field.Tag.Get("weaselform")))

		}

		val := valField.Interface()
		if valField.Kind() == reflect.Ptr {
			val = valField.Addr().Interface()
		}

		e := Element{
			Name : field.Name,
			HashName : crypto.Encrypt(field.Name, f.salt),
			Label : field.Tag.Get("formLabel"),
			Order : uint(i),
			Type : elementType[field.Tag.Get("weaselform")],
			TypeName : field.Tag.Get("weaselform"),
			Value: val,
		}

		//meth, ok := field.Type.MethodByName("Options")
		//
		//if ok {
		//
		//	opts := meth.Func.Interface().(func() Options)
		//
		//	e.Options = opts
		//}

		f.Elements = append(f.Elements, &e)

	}

	return nil
}

//Set values to struct
func (f *Form) unmarshal(s interface {}) error {

	v := reflect.ValueOf(s)

	st := reflect.TypeOf(s)

	if st.Kind() == reflect.Ptr {

		st = st.Elem()

	}

	if st.Kind() != reflect.Struct {

		return errors.New(fmt.Sprintf("Form unmarshal recieved %s, but needs *Struct", st.Kind().String()))
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	vals := map[string]string{}

	for _, e := range f.Elements {

		vals[e.Name] = e.GetValue()

		switch v.FieldByName(e.Name).Kind() {

		case reflect.Bool:

			if e.GetValue() == "" || e.GetValue() == "0" {
				v.FieldByName(e.Name).Set(reflect.ValueOf(false))
			} else {
				v.FieldByName(e.Name).Set(reflect.ValueOf(true))
			}

		case reflect.String:
			v.FieldByName(e.Name).SetString(e.GetValue())

		case reflect.Uint :

			val, err := strconv.ParseUint(e.GetValue(), 10, 64)

			if err != nil {
				val = 0
			}

			v.FieldByName(e.Name).Set(reflect.ValueOf(uint(val)))

		case reflect.Float64,
			reflect.Float32:

			val, err := strconv.ParseFloat(e.GetValue(), 64)

			if err != nil {
				val = 0
			}

			v.FieldByName(e.Name).Set(reflect.ValueOf(float64(val)))

		case reflect.Slice :
			v.FieldByName(e.Name).Set(reflect.ValueOf(strings.Split(strings.Trim(e.GetValue(), " ")," ")))

		default:

			continue;
			//return errors.New(fmt.Sprintf("Cannot unmarshal type %s, element %s", v..FieldByName(e.Name).Kind(), e.Name))

		}

	}

	//@todo Try to map dimensions

//	for n := 0; n < v..NumField(); n++ {
//
//		sd := v..Field(n)
//
//		if v..Field(n).Kind() == reflect.Ptr {
//			sd = sd.Elem()
//		}
//
//		if sd.Type().Name() == "Dimensions" {
//
//			fmt.Println(sd.NumMethod())
//
//			//v..Field(n).MethodByName("MapValues").Call([]reflect.Value{reflect.ValueOf(vals)})
//
//		}
//
//	}


	return nil
}

//Caller for unmarshal
func (f *Form) MapValues(r interface {}) error {

	return f.unmarshal(r)
}

//Set values from struct to elements
func (f *Form) SetValues(s interface {}) error {

	st := reflect.TypeOf(s)

	sv := reflect.ValueOf(s)

	if st.Kind() == reflect.Ptr {

		st = st.Elem()
		sv = sv.Elem()

	}

	if st.Kind() != reflect.Struct {

		return errors.New(fmt.Sprintf("Form SetValues recieved %s, but needs Struct", st.Kind().String()))
	}

	for i := 0; i < st.NumField(); i++ {

		field := st.Field(i)

		if field.Tag.Get("weaselform") == "" {

			continue
		}

		if f.skipFields[field.Name] != nil {

			continue
		}

		for _, e := range f.Elements {

			if e.Name == field.Name {

				e.Value = sv.FieldByName(field.Name).Interface()

			}

		}

	}

	return nil
}
