package content

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bosssauce/ponzu/management/editor"
	"github.com/bosssauce/ponzu/system/api"
)

// Referenceable enures there is a way to reference the implenting type from
// within another type's editor and from type-scoped API calls
type Referenceable interface {
	Referenced() []byte
}

// Select returns the []byte of a <select> HTML element plus internal <options> with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func Select(fieldName string, p interface{}, attrs map[string]string, contentType string) []byte {
	ct, ok := Types[contentType]
	if !ok {
		log.Println("Cannot reference an invalid content type:", contentType)
		return nil
	}

	// get a handle to the underlying interface type for decoding
	t := ct()

	// decode all content type from db into options map
	// map["?type=<contentType>&id=<id>"]t.String()
	options := make(map[string]string)
	// jj := db.ContentAll(contentType + "__sorted") // make this an API call
	jj := api.ContentAll(contentType)

	for i := range jj {
		err := json.Unmarshal(jj[i], t)
		if err != nil {
			log.Println("Error decoding into reference handle:", contentType, err)
		}

		// make sure it is an Identifiable
		item, ok := t.(Identifiable)
		if !ok {
			log.Println("Cannot use type", contentType, "as a reference since it does not implement Identifiable")
			return nil
		}

		k := fmt.Sprintf("?type=%s&id=%d", contentType, item.ItemID())
		v := item.String()
		options[k] = v
	}

	return editor.Select(fieldName, p, attrs, options)
}
