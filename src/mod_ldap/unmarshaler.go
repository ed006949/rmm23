package mod_ldap

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// UnmarshalLDAP adapts a go-ldap Entry to the generic Unmarshal.
func UnmarshalLDAP(e *ldap.Entry, dst any) error {
	if e == nil {
		return nil
	}

	attrs := make(map[string]AttributeValues, len(e.Attributes))
	for _, a := range e.Attributes {
		// go-ldap exposes both Values (as []string) and ByteValues (as [][]byte).
		// Prefer ByteValues if present; else use Values.
		name := strings.ToLower(a.Name)
		if len(a.ByteValues) > 0 {
			// Clone to avoid retaining underlying buffers (optional)
			bins := make([][]byte, len(a.ByteValues))
			for i := range a.ByteValues {
				if a.ByteValues[i] == nil {
					continue
				}

				cp := make([]byte, len(a.ByteValues[i]))
				copy(cp, a.ByteValues[i])
				bins[i] = cp
			}

			attrs[name] = AttributeValues{Binary: bins}

			continue
		}

		if len(a.Values) > 0 {
			// Copy strings to avoid surprises (optional)
			texts := make([]string, len(a.Values))
			copy(texts, a.Values)
			attrs[name] = AttributeValues{Text: texts}
		}
	}

	entry := Entry{
		DN:         e.DN,
		Attributes: attrs,
	}

	return Unmarshal(entry, dst)
}

// Entry is a minimal LDAP entry abstraction.
// Values can be provided as []string (UTF-8) or [][]byte (raw).
type Entry struct {
	DN         string
	Attributes map[string]AttributeValues
}

// AttributeValues can carry either textual or binary values.
// Only one of Text or Binary should be non-nil for a given attribute.
type AttributeValues struct {
	Text   []string
	Binary [][]byte
}

// Unmarshal fills dst (pointer to struct) from Entry.
// - Uses struct tags `ldap:"attr[,omitempty]"`
// - Single-value fields get first attribute value
// - Slice fields get all values
// - Field types supporting encoding.TextUnmarshaler take precedence
// - time.Time supported via common generalizedTime and RFC3339, plus as Unix seconds if numeric.
func Unmarshal(e Entry, dst any) error {
	if dst == nil {
		return errors.New("dst is nil")
	}

	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("dst must be a non-nil pointer")
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("dst must point to a struct")
	}

	// Special-case: allow a field tagged as "dn" to capture DN
	// Example: `ldap:"dn"`
	dnConsumed := false

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		if !sf.IsExported() {
			continue
		}

		tag := sf.Tag.Get("ldap")
		if tag == "-" {
			continue
		}

		name, opts := parseTag(tag)
		if name == "" {
			// default: field name lowercased
			name = strings.ToLower(sf.Name)
		}

		fv := rv.Field(i)
		if name == "dn" {
			// Assign DN if possible
			if err := assignSingleValue(fv, []byte(e.DN)); err != nil {
				return fmt.Errorf("field %s (dn): %w", sf.Name, err)
			}

			dnConsumed = true

			continue
		}

		attr, ok := e.Attributes[strings.ToLower(name)]
		if !ok {
			// handle omitempty: leave zero value
			if opts.Contains("omitempty") {
				continue
			}
			// If not omitempty, still leave zero value (LDAP may not provide attribute)
			continue
		}

		// Decide single vs multi based on field kind
		if isSliceButNotByteSlice(fv) {
			// all values
			if err := assignMultiValue(fv, attr); err != nil {
				return fmt.Errorf("field %s: %w", sf.Name, err)
			}
		} else {
			// first value
			var first []byte
			if len(attr.Binary) > 0 {
				first = attr.Binary[0]
			} else if len(attr.Text) > 0 {
				first = []byte(attr.Text[0])
			}

			if len(first) == 0 {
				// nothing present; leave zero value
				continue
			}

			if err := assignSingleValue(fv, first); err != nil {
				return fmt.Errorf("field %s: %w", sf.Name, err)
			}
		}
	}

	_ = dnConsumed

	return nil
}

// parseTag splits `attr[,opt1,opt2]`.
func parseTag(tag string) (name string, opts tagOptions) {
	if tag == "" {
		return "", nil
	}

	parts := strings.Split(tag, ",")

	name = strings.TrimSpace(parts[0])
	for _, p := range parts[1:] {
		opts = append(opts, strings.TrimSpace(p))
	}

	return
}

type tagOptions []string

func (o tagOptions) Contains(s string) bool {
	for _, v := range o {
		if v == s {
			return true
		}
	}

	return false
}

// Helpers

func isSliceButNotByteSlice(v reflect.Value) bool {
	if v.Kind() != reflect.Slice {
		return false
	}
	// []byte special-case
	return !(v.Type().Elem().Kind() == reflect.Uint8)
}

func assignMultiValue(fv reflect.Value, av AttributeValues) error {
	// If destination is []byte, that's not "multi" path.
	if fv.Kind() == reflect.Slice && fv.Type().Elem().Kind() == reflect.Uint8 {
		// ambiguous; but by convention multi-path only for non-[]byte slices
		// Choose first value
		var b []byte
		if len(av.Binary) > 0 {
			b = av.Binary[0]
		} else if len(av.Text) > 0 {
			b = []byte(av.Text[0])
		}

		return assignSingleValue(fv, b)
	}

	// Ensure it's a slice
	if fv.Kind() != reflect.Slice {
		return fmt.Errorf("expected slice for multi-value attribute, got %s", fv.Kind())
	}

	elemT := fv.Type().Elem()
	// Prepare a new slice
	var count int
	if len(av.Binary) > 0 {
		count = len(av.Binary)
	} else {
		count = len(av.Text)
	}

	slice := reflect.MakeSlice(fv.Type(), 0, count)

	for i := 0; i < count; i++ {
		var b []byte
		if len(av.Binary) > 0 {
			b = av.Binary[i]
		} else {
			b = []byte(av.Text[i])
		}

		ev := reflect.New(elemT).Elem()
		if err := assignSingleValue(ev, b); err != nil {
			return err
		}

		slice = reflect.Append(slice, ev)
	}

	fv.Set(slice)

	return nil
}

func assignSingleValue(fv reflect.Value, b []byte) error {
	// Handle pointers
	if fv.Kind() == reflect.Pointer {
		if fv.IsNil() {
			fv.Set(reflect.New(fv.Type().Elem()))
		}

		return assignSingleValue(fv.Elem(), b)
	}

	// If destination is time.Time, handle first (textual or numeric)
	if fv.Type() == reflect.TypeOf(time.Time{}) {
		// 1) Try LDAP Generalized Time and RFC3339
		if t, err := parseLDAPTime(b); err == nil {
			fv.Set(reflect.ValueOf(t))

			return nil
		}
		// 2) Try Unix seconds (numeric)
		if secs, err := strconv.ParseInt(strings.TrimSpace(string(b)), 10, 64); err == nil {
			fv.Set(reflect.ValueOf(time.Unix(secs, 0).UTC()))

			return nil
		}

		return fmt.Errorf("unsupported time format: %q", string(b))
	}

	// TextUnmarshaler check AFTER time.Time special-case to avoid custom types masking time handling
	if tu, ok := asTextUnmarshaler(fv); ok {
		return tu.UnmarshalText(b)
	}

	switch fv.Kind() {
	case reflect.String:
		fv.SetString(string(b))

		return nil
	case reflect.Bool:
		x, err := strconv.ParseBool(strings.TrimSpace(string(b)))
		if err != nil {
			return err
		}

		fv.SetBool(x)

		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err := strconv.ParseInt(strings.TrimSpace(string(b)), 10, fv.Type().Bits())
		if err != nil {
			return err
		}

		fv.SetInt(x)

		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		x, err := strconv.ParseUint(strings.TrimSpace(string(b)), 10, fv.Type().Bits())
		if err != nil {
			return err
		}

		fv.SetUint(x)

		return nil
	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(strings.TrimSpace(string(b)), fv.Type().Bits())
		if err != nil {
			return err
		}

		fv.SetFloat(x)

		return nil
	case reflect.Slice:
		if fv.Type().Elem().Kind() == reflect.Uint8 {
			dst := make([]byte, len(b))
			copy(dst, b)
			fv.SetBytes(dst)

			return nil
		}
	case reflect.Struct:
		// Other structs only via TextUnmarshaler
	}

	return fmt.Errorf("cannot assign to kind %s (type %s)", fv.Kind(), fv.Type())
}

func asTextUnmarshaler(v reflect.Value) (encoding.TextUnmarshaler, bool) {
	// If addressable, check pointer receiver too
	if v.CanAddr() {
		if tu, ok := v.Addr().Interface().(encoding.TextUnmarshaler); ok {
			return tu, true
		}
	}

	tu, ok := v.Interface().(encoding.TextUnmarshaler)

	return tu, ok
}

func parseLDAPTime(b []byte) (time.Time, error) {
	s := strings.TrimSpace(string(b))

	// Fast paths for common LDAP forms (Zulu, optional fraction, or numeric offset)
	layouts := []string{
		"20060102150405Z",
		"20060102150405.999999999Z", // allow up to 9 fractional digits
		"20060102150405-0700",
		"20060102150405.999999999-0700",
		"200601021504Z", // no seconds
		"2006010215Z",   // no minutes/seconds
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC(), nil
		}
	}

	// If AD-style requires a fixed ".0Z", handle equivalent normalization
	if strings.HasSuffix(s, "Z") && strings.Contains(s, ".") {
		// normalize fractional part length to up to 9 for Go parsing
		// (optional: implement robust rewriter)
		// If needed, pad or trim fraction then parse with .999999999Z layout
		parts := strings.SplitN(s, ".", 2)
		if len(parts) == 2 {
			fracZ := parts[1]
			if strings.HasSuffix(fracZ, "Z") {
				frac := strings.TrimSuffix(fracZ, "Z")
				if len(frac) > 0 && len(frac) <= 9 {
					fracPadded := frac + strings.Repeat("0", 9-len(frac))

					candidate := parts[0] + "." + fracPadded + "Z"
					if t, err := time.Parse("20060102150405.999999999Z", candidate); err == nil {
						return t.UTC(), nil
					}
				}
			}
		}
	}

	// RFC3339 fallbacks for non-LDAP servers
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t.UTC(), nil
	}

	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UTC(), nil
	}

	return time.Time{}, fmt.Errorf("unsupported LDAP time %q", s)
}
