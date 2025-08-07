package mod_ldap

// // Helper: Get LDAP attr values by name.
// func getAttributeValues(e *ldap.Entry, name string) (values []string) {
// 	for _, attr := range e.Attributes {
// 		switch {
// 		case attr.Name == name:
// 			return attr.Values
// 		}
// 	}
//
// 	return
// }
//
// // UnmarshalEntry maps ldap.Entry attributes to struct fields (by `ldap` tag), supporting TextUnmarshaler and slices.
// //
// // Pointers to scalars (*T)
// //
// // Pointers to slices (*[]T)
// //
// // Slices of pointers ([]*T)
// //
// // Slices of values ([]T)
// //
// // Scalars (T)
// //
// // All of the above with support for encoding.TextUnmarshaler (for both value and pointer receivers)
// //
// // DN shortcut and field skipping logic.
// func UnmarshalEntry(e *ldap.Entry, out interface{}) error {
// 	var (
// 		vo = reflect.ValueOf(out)
// 	)
// 	switch {
// 	case vo.Kind() != reflect.Ptr || vo.IsNil() || vo.Elem().Kind() != reflect.Struct:
// 		return errors.New("UnmarshalEntry: expected pointer to struct")
// 	}
//
// 	var (
// 		val = vo.Elem()
// 		typ = val.Type()
// 	)
// 	for i := 0; i < typ.NumField(); i++ {
// 		var (
// 			field = typ.Field(i)
// 			fv    = val.Field(i)
// 		)
// 		switch {
// 		case len(field.PkgPath) != 0:
// 			continue
// 		}
//
// 		var (
// 			tag = field.Tag.Get("ldap")
// 		)
// 		switch {
// 		case len(tag) == 0:
// 			continue
// 		}
//
// 		// Handle "dn" shortcut
// 		switch {
// 		case tag == "dn":
// 			switch {
// 			case fv.Kind() == reflect.String && fv.CanSet():
// 				fv.SetString(e.DN)
// 			}
//
// 			continue
// 		}
//
// 		var (
// 			values = getAttributeValues(e, tag)
// 		)
// 		switch {
// 		case len(values) == 0:
// 			continue
// 		}
//
// 		// Pointer to scalar (e.g. *string, *int64, *CustomType)
// 		switch {
// 		case fv.Kind() == reflect.Ptr && fv.CanSet():
// 			var (
// 				elemType = fv.Type().Elem()
// 				elemPtr  = fv
// 			)
// 			switch {
// 			case fv.IsNil():
// 				elemPtr = reflect.New(elemType)
// 				fv.Set(elemPtr)
// 			}
//
// 			var (
// 				elem = fv.Elem() // TextUnmarshaler support for pointer
// 			)
//
// 			switch u, ok := elemPtr.Interface().(encoding.TextUnmarshaler); {
// 			case ok:
// 				switch err := u.UnmarshalText([]byte(values[0])); {
// 				case err != nil:
// 					return err
// 				}
//
// 				continue
// 			}
//
// 			switch elem.Kind() {
// 			case reflect.String:
// 				elem.SetString(values[0])
// 			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 				var (
// 					n, err = strconv.ParseInt(values[0], 10, 64)
// 				)
// 				switch {
// 				case err != nil:
// 					return err
// 				}
//
// 				elem.SetInt(n)
// 			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
// 				var (
// 					n, err = strconv.ParseUint(values[0], 10, 64)
// 				)
// 				switch {
// 				case err != nil:
// 					return err
// 				}
//
// 				elem.SetUint(n)
// 			case reflect.Bool:
// 				var (
// 					b, err = strconv.ParseBool(values[0])
// 				)
// 				switch {
// 				case err != nil:
// 					return err
// 				}
//
// 				elem.SetBool(b)
// 			default:
// 				return fmt.Errorf("unsupported pointer type: %v", elem.Type())
// 			}
//
// 			continue
// 		}
//
// 		// Pointer to slice (e.g. *[]string, *[]int64, *[]CustomType)
// 		switch {
// 		case fv.Kind() == reflect.Ptr && fv.CanSet() && fv.Type().Elem().Kind() == reflect.Slice:
// 			var (
// 				sliceType = fv.Type().Elem()
// 				elemType  = sliceType.Elem()
// 			)
// 			switch {
// 			case fv.IsNil():
// 				fv.Set(reflect.New(sliceType))
// 			}
//
// 			var (
// 				slice = reflect.MakeSlice(sliceType, len(values), len(values))
// 			)
// 			for j, v := range values {
// 				var (
// 					elemPtr = reflect.New(elemType) // Check for TextUnmarshaler for each element
// 				)
// 				switch u, ok := elemPtr.Interface().(encoding.TextUnmarshaler); {
// 				case ok:
// 					switch err := u.UnmarshalText([]byte(v)); {
// 					case err != nil:
// 						return fmt.Errorf("unmarshal '%s': %w", v, err)
// 					}
//
// 					slice.Index(j).Set(elemPtr.Elem())
//
// 					continue
// 				}
// 				// Standard types
// 				switch elemType.Kind() {
// 				case reflect.String:
// 					slice.Index(j).SetString(v)
// 				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 					var (
// 						n, err = strconv.ParseInt(v, 10, 64)
// 					)
// 					switch {
// 					case err != nil:
// 						return err
// 					}
//
// 					slice.Index(j).SetInt(n)
// 				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
// 					var (
// 						n, err = strconv.ParseUint(v, 10, 64)
// 					)
// 					switch {
// 					case err != nil:
// 						return err
// 					}
//
// 					slice.Index(j).SetUint(n)
// 				case reflect.Bool:
// 					var (
// 						b, err = strconv.ParseBool(v)
// 					)
// 					switch {
// 					case err != nil:
// 						return err
// 					}
//
// 					slice.Index(j).SetBool(b)
// 				default:
// 					return fmt.Errorf("unsupported pointer-to-slice element type: %v", elemType)
// 				}
// 			}
//
// 			fv.Elem().Set(slice)
//
// 			continue
// 		}
//
// 		// Slices of pointers (e.g. []*string, []*int, []*CustomType)
// 		switch {
// 		case fv.Kind() == reflect.Slice && fv.Type().Elem().Kind() == reflect.Ptr && fv.CanSet():
// 			var (
// 				elemPtrType = fv.Type().Elem()   // *T
// 				elemType    = elemPtrType.Elem() // T
// 				slice       = reflect.MakeSlice(fv.Type(), len(values), len(values))
// 			)
// 			for j, v := range values {
// 				var (
// 					elemPtr = reflect.New(elemType) // *T
// 				)
// 				// Check for TextUnmarshaler
// 				switch u, ok := elemPtr.Interface().(encoding.TextUnmarshaler); {
// 				case ok:
// 					switch err := u.UnmarshalText([]byte(v)); {
// 					case err != nil:
// 						return fmt.Errorf("unmarshal '%s': %w", v, err)
// 					}
//
// 					slice.Index(j).Set(elemPtr)
//
// 					continue
// 				}
//
// 				switch elemType.Kind() {
// 				case reflect.String:
// 					elemPtr.Elem().SetString(v)
// 				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 					var (
// 						n, err = strconv.ParseInt(v, 10, 64)
// 					)
// 					switch {
// 					case err != nil:
// 						return err
// 					}
//
// 					elemPtr.Elem().SetInt(n)
// 				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
// 					var (
// 						n, err = strconv.ParseUint(v, 10, 64)
// 					)
// 					switch {
// 					case err != nil:
// 						return err
// 					}
//
// 					elemPtr.Elem().SetUint(n)
// 				case reflect.Bool:
// 					var (
// 						b, err = strconv.ParseBool(v)
// 					)
// 					switch {
// 					case err != nil:
// 						return err
// 					}
//
// 					elemPtr.Elem().SetBool(b)
// 				default:
// 					return fmt.Errorf("unsupported slice-of-pointers elem type: %v", elemType)
// 				}
//
// 				slice.Index(j).Set(elemPtr)
// 			}
//
// 			fv.Set(slice)
//
// 			continue
// 		}
//
// 		// Slices: check UnmarshalText first
// 		switch {
// 		case fv.Kind() == reflect.Slice && fv.CanSet():
// 			var (
// 				elemType = fv.Type().Elem()
// 				slice    = reflect.MakeSlice(fv.Type(), len(values), len(values))
// 			)
// 			for j, v := range values {
// 				var (
// 					elemPtr = reflect.New(elemType)
// 				)
// 				switch u, ok := elemPtr.Interface().(encoding.TextUnmarshaler); {
// 				case ok:
// 					switch err := u.UnmarshalText([]byte(v)); {
// 					case err != nil:
// 						return fmt.Errorf("unmarshal '%s': %w", v, err)
// 					}
//
// 					slice.Index(j).Set(elemPtr.Elem())
//
// 					continue
// 				}
// 				// Standard slice element types
// 				switch elemType.Kind() {
// 				case reflect.String:
// 					slice.Index(j).SetString(v)
// 				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 					var (
// 						n, err = strconv.ParseInt(v, 10, 64)
// 					)
// 					switch {
// 					case err != nil:
// 						return err
// 					}
//
// 					slice.Index(j).SetInt(n)
// 				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
// 					var (
// 						n, err = strconv.ParseUint(v, 10, 64)
// 					)
// 					switch {
// 					case err != nil:
// 						return err
// 					}
//
// 					slice.Index(j).SetUint(n)
// 				case reflect.Bool:
// 					var (
// 						b, err = strconv.ParseBool(v)
// 					)
// 					switch {
// 					case err != nil:
// 						return err
// 					}
//
// 					slice.Index(j).SetBool(b)
// 				default:
// 					return fmt.Errorf("unsupported slice element type: %v", elemType)
// 				}
// 			}
//
// 			fv.Set(slice)
//
// 			continue
// 		}
//
// 		var (
// 			// Scalar: check UnmarshalText first
// 			ptr = fv
// 		)
// 		switch {
// 		case fv.CanAddr():
// 			ptr = fv.Addr()
// 		}
//
// 		switch u, ok := ptr.Interface().(encoding.TextUnmarshaler); {
// 		case ok:
// 			switch err := u.UnmarshalText([]byte(values[0])); {
// 			case err != nil:
// 				return err
// 			}
//
// 			continue
// 		}
//
// 		// Standard types
// 		switch fv.Kind() {
// 		case reflect.String:
// 			fv.SetString(values[0])
// 		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
// 			var (
// 				n, err = strconv.ParseUint(values[0], 10, 64)
// 			)
// 			switch {
// 			case err != nil:
// 				return err
// 			}
//
// 			fv.SetUint(n)
// 		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 			var (
// 				n, err = strconv.ParseInt(values[0], 10, 64)
// 			)
// 			switch {
// 			case err != nil:
// 				return err
// 			}
//
// 			fv.SetInt(n)
// 		case reflect.Bool:
// 			var (
// 				b, err = strconv.ParseBool(values[0])
// 			)
// 			switch {
// 			case err != nil:
// 				return err
// 			}
//
// 			fv.SetBool(b)
// 		default:
// 			return fmt.Errorf("unsupported element type: %v", fv.Type().String())
// 		}
// 	}
//
// 	return nil
// }

// func readTag(f reflect.StructField) (options string, flag bool) {
// 	var (
// 		val, ok = f.Tag.Lookup(ldapTagName)
// 	)
// 	switch {
// 	case !ok:
// 		return f.Name, false
// 	}
//
// 	var (
// 		// opts = mod_slices.SplitString(val, mod_strings.TagSeparator, mod_slices.FlagNormalize)
// 		opts = strings.Split(val, ",")
// 	)
// 	switch {
// 	case len(opts) == mod_slices.KVElements:
// 		flag = opts[1] == ldapTagOptionOmitEmpty
// 	}
//
// 	return opts[0], flag
// }
//
// func UnmarshalEntry(e *ldap.Entry, i interface{}) (err error) {
// 	var (
// 		sv reflect.Value
// 		st reflect.Type
// 	)
// 	switch sv, st, err = mod_reflect.GetStructSVST(i); {
// 	case err != nil:
// 		return
// 	}
//
// 	for n := 0; n < st.NumField(); n++ {
// 		var (
// 			fv = sv.Field(n) // Holds struct field value
// 			ft = st.Field(n) // Holds struct field type
// 		)
// 		switch {
// 		case len(ft.PkgPath) != 0: // skip unexported fields
// 			continue
// 		}
//
// 		// omitempty can be safely discarded, as it's not needed when unmarshalling
// 		var (
// 			fieldTag, _ = readTag(ft)
// 		)
//
// 		// Fill the field with the distinguishedName if the tag key is `dn`
// 		switch fieldTag {
// 		case "dn":
// 			switch {
// 			case fv.CanSet():
// 				switch unmarshaler, ok := fv.Addr().Interface().(LDAPAttributeUnmarshaler); {
// 				case ok:
// 					switch err = unmarshaler.UnmarshalLDAPAttr([]string{e.DN}); {
// 					case err != nil:
// 						return
// 					}
//
// 					continue
// 				case fv.Kind() == reflect.String:
// 					fv.SetString(e.DN)
//
// 					continue
// 				}
// 			}
//
// 			return mod_errors.EParse
// 		}
//
// 		var (
// 			values = e.GetAttributeValues(fieldTag)
// 		)
// 		switch {
// 		case len(values) == 0:
// 			continue
// 		}
//
// 		// var (
// 		// 	// Handle types implementing encoding.TextUnmarshaler
// 		// 	ptr = fv
// 		// )
// 		// switch {
// 		// case fv.Kind() != reflect.Ptr && fv.CanAddr():
// 		// 	ptr = fv.Addr()
// 		// }
//
// 		// switch textUnmarshaler, ok := ptr.Interface().(encoding.TextUnmarshaler); {
// 		// case ok:
// 		// 	switch err = textUnmarshaler.UnmarshalText([]byte(values[0])); {
// 		// 	case err != nil:
// 		// 		return
// 		// 	}
// 		// 	continue
// 		// }
//
// 		switch fieldType := fv.Interface().(type) {
// 		default:
// 			switch rt := reflect.TypeOf(fieldType); rt.Kind() {
// 			case reflect.Map:
// 				var (
// 					mapVal = reflect.MakeMap(rt)
// 					ptrVal = reflect.New(rt)
// 				)
//
// 				ptrVal.Elem().Set(mapVal)
//
// 				switch unmarshaler, ok := ptrVal.Interface().(LDAPAttributeUnmarshaler); {
// 				case ok:
// 					switch err = unmarshaler.UnmarshalLDAPAttr(values); {
// 					case err != nil:
// 						return
// 					}
//
// 					fv.Set(ptrVal.Elem())
// 				default:
// 					return
// 				}
// 			default:
// 				var (
// 					ptrVal = reflect.New(rt)
// 				)
// 				switch unmarshaler, ok := ptrVal.Interface().(LDAPAttributeUnmarshaler); {
// 				case ok:
// 					switch err = unmarshaler.UnmarshalLDAPAttr(values); {
// 					case err != nil:
// 						return
// 					}
//
// 					fv.Set(ptrVal.Elem())
// 				}
// 			}
// 		}
// 	}
//
// 	return
// }
