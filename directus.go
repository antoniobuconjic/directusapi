package directusapi

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type Version int

const (
	V8 Version = iota
	V9
)

type PrimaryKey interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~string
}

// API is a generic API client for any Directus collection
// R is a read model
// W is a write model
// PK is a type of primary key
type API[R, W any, PK PrimaryKey] struct {
	Scheme         string
	Host           string
	Namespace      string
	CollectionName string
	BearerToken    string
	HTTPClient     *http.Client
	queryFields    []string
	debug          bool
	Version        Version
}

// CreateToken uses provided credentials to generate server token
//
// Related Directus reference:
// https://v8.docs.directus.io/api/authentication.html#retrieve-a-temporary-access-token
func (d API[R, W, PK]) CreateToken(ctx context.Context, email, password string) (string, error) {
	u := fmt.Sprintf("%s://%s/%s/auth/authenticate", d.Scheme, d.Host, d.Namespace)

	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		email,
		password,
	}

	req := request{
		ctx,
		http.MethodPost,
		u,
		nil,
		body,
	}

	var respBody struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}

	err := d.executeRequest(req, http.StatusOK, &respBody)
	if err != nil {
		return "", fmt.Errorf("execute create token request: %w", err)
	}
	return respBody.Data.Token, nil
}

// Insert attempts to insert new item
//
// Related Directus reference:
// https://v8.docs.directus.io/api/items.html#create-an-item
func (d API[R, W, PK]) Insert(ctx context.Context, item W) (R, error) {
	var empty R
	u := fmt.Sprintf("%s://%s/%s/items/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName)

	req := request{
		ctx,
		http.MethodPost,
		u,
		map[string]string{
			"fields": strings.Join(d.jsonFieldsR(), ","),
		},
		item,
	}
	var respBody struct {
		Data R `json:"data"`
	}
	err := d.executeRequest(req, http.StatusOK, &respBody)
	if err != nil {
		return empty, fmt.Errorf("execute insert request: %w", err)
	}
	return respBody.Data, nil
}

// Create attempts to create new item with partials
//
// Related Directus reference:
// https://v8.docs.directus.io/api/items.html#create-an-item
func (d API[R, W, PK]) Create(ctx context.Context, partials map[string]any) (R, error) {
	var empty R
	u := fmt.Sprintf("%s://%s/%s/items/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName)

	req := request{
		ctx,
		http.MethodPost,
		u,
		map[string]string{
			"fields": strings.Join(d.jsonFieldsR(), ","),
		},
		partials,
	}

	var respBody struct {
		Data R `json:"data"`
	}
	err := d.executeRequest(req, http.StatusOK, &respBody)
	if err != nil {
		return empty, fmt.Errorf("execute create request: %w", err)
	}
	return respBody.Data, nil

}

// GetByID reads a single item by given ID
//
// Related Directus reference:
// https://v8.docs.directus.io/api/items.html#retrieve-an-item
func (d API[R, W, PK]) GetByID(ctx context.Context, id PK) (R, error) {
	u := fmt.Sprintf("%s://%s/%s/items/%s/%v", d.Scheme, d.Host, d.Namespace, d.CollectionName, id)

	req := request{
		ctx,
		http.MethodGet,
		u,
		map[string]string{
			"fields": strings.Join(d.jsonFieldsR(), ","),
		},
		nil,
	}

	var respBody struct {
		Data R `json:"data"`
	}
	var empty R
	err := d.executeRequest(req, http.StatusOK, &respBody)
	if err != nil {
		return empty, fmt.Errorf("execute get by id request: %w", err)
	}
	return respBody.Data, nil
}

// Update performs partial update of an item with given id
//
// Related Directus reference:
// https://v8.docs.directus.io/api/items.html#update-an-item
func (d API[R, W, PK]) Update(ctx context.Context, id PK, partials map[string]any) (R, error) {
	var empty R
	u := fmt.Sprintf("%s://%s/%s/items/%s/%v", d.Scheme, d.Host, d.Namespace, d.CollectionName, id)

	req := request{
		ctx,
		http.MethodPatch,
		u,
		map[string]string{
			"fields": strings.Join(d.jsonFieldsR(), ","),
		},
		partials,
	}

	var respBody struct {
		Data R `json:"data"`
	}
	err := d.executeRequest(req, http.StatusOK, &respBody)
	if err != nil {
		return empty, fmt.Errorf("execute update request: %w", err)
	}
	return respBody.Data, nil
}

// Set performs an update of an item with given id
//
// Related Directus reference:
// https://v8.docs.directus.io/api/items.html#update-an-item
func (d API[R, W, PK]) Set(ctx context.Context, id PK, item W) (R, error) {
	var empty R
	u := fmt.Sprintf("%s://%s/%s/items/%s/%v", d.Scheme, d.Host, d.Namespace, d.CollectionName, id)

	req := request{
		ctx,
		http.MethodPatch,
		u,
		map[string]string{
			"fields": strings.Join(d.jsonFieldsR(), ","),
		},
		item,
	}

	var respBody struct {
		Data R `json:"data"`
	}
	err := d.executeRequest(req, http.StatusOK, &respBody)
	if err != nil {
		return empty, fmt.Errorf("execute set request: %w", err)
	}
	return respBody.Data, nil
}

// Delete removes item with a given id
//
// Related Directus reference:
// https://v8.docs.directus.io/api/items.html#update-an-item
func (d API[R, W, PK]) Delete(ctx context.Context, id PK) error {
	u := fmt.Sprintf("%s://%s/%s/items/%s/%v", d.Scheme, d.Host, d.Namespace, d.CollectionName, id)
	req := request{
		ctx,
		http.MethodDelete,
		u,
		nil,
		nil,
	}

	err := d.executeRequest(req, http.StatusNoContent, nil)
	if err != nil {
		return fmt.Errorf("execute delete request: %w", err)
	}
	return nil
}

// Items retrieves a collection of items
//
// Related Directus reference:
// https://v8.docs.directus.io/api/items.html#update-an-item
func (d API[R, W, PK]) Items(ctx context.Context, q query) ([]R, error) {
	u := fmt.Sprintf("%s://%s/%s/items/%s", d.Scheme, d.Host, d.Namespace, d.CollectionName)
	qv := q.asKeyValue(d.Version)
	qv["fields"] = strings.Join(d.jsonFieldsR(), ",")

	req := request{
		ctx,
		http.MethodGet,
		u,
		qv,
		nil,
	}
	var respBody struct {
		Data []R `json:"data"`
	}
	err := d.executeRequest(req, http.StatusOK, &respBody)
	if err != nil {
		return nil, fmt.Errorf("execute items request: %w", err)
	}
	return respBody.Data, nil
}

func (d *API[R, W, PK]) jsonFieldsR() []string {
	if d.queryFields == nil {
		var x R
		t := reflect.TypeOf(x)
		d.queryFields = iterateFields(t, "")
	}
	return d.queryFields
}

// iterateFields returns fields for all struct's fields
func iterateFields(t reflect.Type, prefix string) []string {
	fields := []string{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fields = append(fields, structFields(f, prefix)...)
	}
	return fields
}

// structFields returns fields for a signle struct field
func structFields(f reflect.StructField, prefix string) []string {
	tagVal := ""
	if v, ok := f.Tag.Lookup(tagName); ok {
		tagVal = v
	} else {
		tagVal = f.Name
	}
	switch f.Type.Kind() {
	case reflect.Struct:
		var t Time
		isTime := f.Type.ConvertibleTo(reflect.TypeOf(t))
		isOptional := f.Type.Implements(reflect.TypeOf(new(isOpt)).Elem())
		switch {
		case isOptional:
			val := reflect.New(f.Type).Interface().(isOpt)
			if prefix == "" {
				return val.fields(tagVal)
			} else {
				return val.fields(prefix + "." + tagVal)
			}
		case isTime:
			return []string{prefix}
		default:
			p := prefix
			if p != "" {
				p = p + "." + tagVal
			} else {
				p = tagVal
			}
			return iterateFields(f.Type, p)
		}
	case reflect.Slice:
		p := prefix
		if p != "" {
			p = p + "." + tagVal
		} else {
			p = tagVal
		}
		if f.Type.Elem().Kind() == reflect.Struct {
			return iterateFields(f.Type.Elem(), p)
		}
		return []string{p}
	case
		reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64,
		reflect.String, reflect.Map:
		// field is not nested
		v := tagVal
		if prefix != "" {
			v = prefix + "." + tagVal
		}
		return []string{v}
	case reflect.Pointer:
		t := f.Type.Elem().String()
		panic(f.Name + "(" + t + "," + prefix + "): pointer is not supported, use directus.Optional instead")
	default:
		panic(f.Type.Kind().String() + " not implemented")
	}
}

type isOpt interface {
	getOp() operation
	fields(prefix string) []string
}
