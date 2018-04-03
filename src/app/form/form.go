package form

import (
	"encoding/json"
	"app/crypto"
	"net/http"
	"sort"
	"fmt"
)

type Form struct {
	Action   string
	Name     string
	Title    string
	Elements []*Element `json:"e"`
	salt     string
	Method   string
	skipFields map[string]interface {}
}

func New(t, n, s string) *Form {

	return &Form {
		Action  :  "",
		Name    :  n,
		Title   :  t,
		Elements:  []*Element{},
		Method  :  "POST",
		skipFields : map[string]interface {}{},
		salt    : s,
	}
}

func (f *Form) GetElement(name string) *Element {

	for _, el := range f.Elements {
		if el.Name == name {
			return el
		}
	}

	return nil
}

func (f *Form) Fields(els ...*Element) {

	for _, el := range els {

		el.HashName = crypto.Encrypt(el.Name, f.salt)

		f.Elements = append(f.Elements, el)

	}

}

func (f *Form) ParseForm(reciever interface {}, req *http.Request) error {

	req.ParseForm()

	for _, e := range f.Elements {

		switch e.Type {

		case Checkbox:

			e.Value = fmt.Sprintf("%v",req.Form[e.HashName] != nil)

		default:

			e.Value = req.PostFormValue(e.HashName)

		}

	}

	if reciever == nil {

		return nil
	}

	return f.unmarshal(reciever)
}

func (f *Form) Values() map[string]string {

	res := map[string]string{}

	for _, el := range f.Elements {

		res[el.Name] = el.GetValue()

	}

	return res
}

func (f *Form) UnmarshalValues(src string) error {

	vals := map[string]interface {}{}

	if err := json.Unmarshal([]byte(src), &vals); err != nil {
		return err
	}

	for _, e := range f.Elements {

		e.Value = vals[e.Name]

	}

	return nil
}

func (f *Form) Context() *Form {

	sort.Sort(f)

	return f
}

func (f *Form) Skip(fields ...string) {

	for _, nm := range fields {

		f.skipFields[nm] = "skip"
	}

}

func (f *Form) Len() int {
	return len(f.Elements)
}

func (f *Form) Swap(i, j int) {
	f.Elements[i], f.Elements[j] = f.Elements[j], f.Elements[i]
}

func (f *Form) Less(i, j int) bool {
	return f.Elements[i].Order < f.Elements[j].Order
}
