package gjm

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// GetProperty returns a property if it exist.
//
//    property, err := GetProperty(document, "one.two.three[0]", ".")
//
// Property type is `interface{}`
func GetProperty(original_data map[string]interface{}, path, separator string) (path_parsed interface{}, err error) {
	if len(separator) == 0 {
		separator = "."
	}

	// Protect the original map :D
	data := make(map[string]interface{})
	for k, v := range original_data {
		data[k] = v
	}
	err = fmt.Errorf("Property %s does not exist", path)

	if len(path) == 0 {
		path = separator
	}

	levels_tmp := strings.Split(path, separator)
	levels := make([]string, 0)
	for _, level_tmp := range levels_tmp {
		if len(level_tmp) > 0 {
			levels = append(levels, level_tmp)
		}
	}

	if len(levels) > 0 && path != separator {
		path_level_one := levels[0]

		// If we have a level in path_level_one

		re := regexp.MustCompile(`\w+\[\d+\]{1}`)
		if matched := re.FindString(path_level_one); len(matched) > 0 {
			property_re := regexp.MustCompile(`\w+`)
			index_re := regexp.MustCompile(`\[\d+\]{1}`)
			// Get a property
			// avatars
			property := property_re.FindString(path_level_one)

			// Get an index
			index_found := index_re.FindString(path_level_one)

			// If index > 0 - check if this property is array
			if len(index_found) > 0 {
				if len(property) > 0 {
					path_level_one = property
				}
				index_found = strings.Trim(index_found, "[]")
				if index, err := strconv.Atoi(index_found); err == nil {
					if v, ok := data[property]; ok {
						if IsKind(v, reflect.Slice) {
							slice := reflect.ValueOf(v)
							if index >= 0 && index < slice.Len() {
								value := slice.Index(index).Interface()

								data[property] = value
							} else {
								err = fmt.Errorf(
									"%s: Min index is 0, Max index is %d", property, slice.Len(),
								)
								return path_parsed, err
							}
						} else {
							err = fmt.Errorf(
								"%s: is not an array", property,
							)
							return path_parsed, err
						}
					} else {
						err = fmt.Errorf(
							"Property %s does not exist", property,
						)
						return path_parsed, err
					}
				} else {
					err = fmt.Errorf(
						"%s must be of type %s",
						fmt.Sprintf("%s[%d]", property, index_found),
						"number",
					)
					return path_parsed, err
				}
			}
		}

		if len(levels[1:]) >= 1 {
			if level_one_value, ok := data[path_level_one]; ok {
				if level_one_value != nil {
					switch reflect.TypeOf(level_one_value).Kind() {
					case reflect.Map:
						if mapped_level_one_value, ok := level_one_value.(map[string]interface{}); ok {
							if path_parsed_local, err_local := GetProperty(mapped_level_one_value, strings.Join(levels[1:], separator), separator); err_local != nil {
								return path_parsed, err_local
							} else {
								path_parsed = path_parsed_local
								err = nil
							}
						}
					default:
						path_parsed = map[string]interface{}{
							path_level_one: level_one_value,
						}
						err = nil
					}
				}
			} else {
				err = fmt.Errorf(
					"Property %s does not exist", path_level_one,
				)
				return path_parsed, err
			}
		} else {
			if v, ok := data[path_level_one]; ok {
				path_parsed = v
				err = nil
			}
		}
	} else if path == separator {
		path_parsed = data
		err = nil
	}
	return
}

// DeleteProperty removes a property from map
//
//    err := GetProperty(document, "one.two.three[0]", ".")
//
func DeleteProperty(original_data map[string]interface{}, path, separator string) (err error) {
	if len(separator) == 0 {
		separator = "."
	}

	// If we have a property
	if _, err = GetProperty(original_data, path, separator); err != nil {
		return
	}

	if len(path) == 0 {
		path = separator
	}

	levels_tmp := strings.Split(path, separator)
	levels := make([]string, 0)
	for _, level_tmp := range levels_tmp {
		if len(level_tmp) > 0 {
			levels = append(levels, level_tmp)
		}
	}

	if len(levels) > 0 && path != separator {
		path_level_one := levels[0]

		// If we have a level in path_level_one

		re := regexp.MustCompile(`\w+\[\d+\]{1}`)
		if matched := re.FindString(path_level_one); len(matched) > 0 {
			property_re := regexp.MustCompile(`\w+`)
			index_re := regexp.MustCompile(`\[\d+\]{1}`)
			// Get a property
			// avatars
			property := property_re.FindString(path_level_one)

			// Get an index
			index_found := index_re.FindString(path_level_one)

			// If index > 0 - check if this property is array
			if len(index_found) > 0 {
				if len(property) > 0 {
					path_level_one = property
				}
				index_found = strings.Trim(index_found, "[]")
				if index, err := strconv.Atoi(index_found); err == nil {
					if v, ok := original_data[property]; ok {
						if IsKind(v, reflect.Slice) {
							slice := reflect.ValueOf(v)
							if index >= 0 && index < slice.Len() {
								value := slice.Index(index).Interface()
								// If len of other levels greater than 0
								if len(levels[1:]) >= 1 {
									if IsKind(value, reflect.Map) {
										mapped_value := value.(map[string]interface{})
										err = DeleteProperty(mapped_value, strings.Join(levels[1:], separator), separator)
										if err == nil {
											// If we have an empty value inside of a slice - remove it
											if len(mapped_value) == 0 {
												slices := make([]interface{}, 0)
												for i := 0; i < slice.Len(); i++ {
													if i != index {
														slices = append(slices, slice.Index(i).Interface())
													}
												}
												original_data[property] = slices
											}
										}
										return err
									}
								} else {
									// if this is a `property[1]` in a path like `path.to.property[1]`
									slices := make([]interface{}, 0)
									for i := 0; i < slice.Len(); i++ {
										if i != index {
											slices = append(slices, slice.Index(i).Interface())
										}
									}
									original_data[property] = slices
									return err
								}
							} else {
								err = fmt.Errorf(
									"%s: Min index is 0, Max index is %d", property, slice.Len(),
								)
							}
						} else {
							err = fmt.Errorf(
								"%s: is not an array", property,
							)
						}
					} else {
						err = fmt.Errorf(
							"Property %s does not exist", property,
						)
					}
				} else {
					err = fmt.Errorf(
						"%s must be of type %s",
						fmt.Sprintf("%s[%d]", property, index_found),
						"number",
					)
				}
			}
		}
		if err != nil {
			return err
		}

		if len(levels[1:]) >= 1 {
			if level_one_value, ok := original_data[path_level_one]; ok {
				if level_one_value != nil {
					switch reflect.TypeOf(level_one_value).Kind() {
					case reflect.Map:
						if mapped_level_one_value, ok := level_one_value.(map[string]interface{}); ok {
							err = DeleteProperty(mapped_level_one_value, strings.Join(levels[1:], separator), separator)
							if err != nil {
								return
							}
						}
					default:
						delete(original_data, path)
					}

				}
			} else {
				err = fmt.Errorf(
					"Property %s does not exist", path_level_one,
				)
				return
			}
		} else {
			delete(original_data, path_level_one)
		}
	} else if path == separator {
		for k := range original_data {
			delete(original_data, k)
		}
	}

	return
}

// AddProperty adds a property to map. Returns an error if property already exists
//
//    err := AddProperty(document, "one.two.three[0]", ".", "string value")
//
func AddProperty(original_data map[string]interface{}, path, separator string, value interface{}) (err error) {
	if len(separator) == 0 {
		separator = "."
	}

	// If we have a property - raise an error
	if _, err = GetProperty(original_data, path, separator); err == nil {
		err = fmt.Errorf(
			"Property %s already exists", path,
		)
		return
	} else {
		err = nil
	}

	if len(path) == 0 {
		path = separator
	}

	levels_tmp := strings.Split(path, separator)
	levels := make([]string, 0)
	for _, level_tmp := range levels_tmp {
		if len(level_tmp) > 0 {
			levels = append(levels, level_tmp)
		}
	}

	if len(levels) > 0 && path != separator {
		path_level_one := levels[0]

		// If we have a level in path_level_one

		re := regexp.MustCompile(`\w+\[\d+\]{1}`)
		if matched := re.FindString(path_level_one); len(matched) > 0 {
			property_re := regexp.MustCompile(`\w+`)
			index_re := regexp.MustCompile(`\[\d+\]{1}`)
			// Get a property
			// avatars
			property := property_re.FindString(path_level_one)

			// Get an index
			index_found := index_re.FindString(path_level_one)

			// If index > 0 - check if this property is array
			if len(index_found) > 0 {
				if len(property) > 0 {
					path_level_one = property
				}
				index_found = strings.Trim(index_found, "[]")
				if index, err := strconv.Atoi(index_found); err == nil {
					if v, ok := original_data[property]; ok {
						if IsKind(v, reflect.Slice) {
							slice := reflect.ValueOf(v)
							var dest_value interface{}
							if index >= 0 && index < slice.Len() {
								dest_value = slice.Index(index).Interface()
							}
							// If len of other levels greater than 0
							if len(levels[1:]) >= 1 {
								if IsKind(dest_value, reflect.Map) {
									mapped_value := dest_value.(map[string]interface{})
									err = AddProperty(mapped_value, strings.Join(levels[1:], separator), separator, value)
									return err
								}
							} else {
								// if this is a `property[1]` in a path like `path.to.property[1]`
								slices := make([]interface{}, 0)
								for i := 0; i < slice.Len(); i++ {
									slices = append(slices, slice.Index(i).Interface())
								}
								if index >= slice.Len() {
									slices = append(slices, value)
								}

								original_data[path_level_one] = slices
								return err
							}
						} else {
							err = fmt.Errorf(
								"%s: is not an array", property,
							)
						}
					} else {
						err = fmt.Errorf(
							"Property %s does not exist", property,
						)
					}
				} else {
					err = fmt.Errorf(
						"%s must be of type %s",
						fmt.Sprintf("%s[%d]", property, index_found),
						"number",
					)
				}
			}
		}

		if err != nil {
			return err
		}

		if len(levels[1:]) >= 1 {
			if level_one_value, ok := original_data[path_level_one]; ok {
				if level_one_value != nil {
					switch reflect.TypeOf(level_one_value).Kind() {
					case reflect.Map:
						if mapped_level_one_value, ok := level_one_value.(map[string]interface{}); ok {
							err = AddProperty(mapped_level_one_value, strings.Join(levels[1:], separator), separator, value)
							if err != nil {
								return
							}
						}
					default:
						original_data[path] = value
					}

				}
			} else {
				err = fmt.Errorf(
					"Property %s does not exist", path_level_one,
				)
				return
			}
		} else {
			// If a map does not contain a last node property
			if _, ok := original_data[path_level_one]; !ok {
				original_data[path_level_one] = value
			}
		}
	} else if path == separator {
		original_data[path] = value
	}
	return
}

// UpdateProperty updates a property in a map. It will create or update existing property
//
//    err := UpdateProperty(document, "one.two.three[0]", ".", "string value")
//
func UpdateProperty(original_data map[string]interface{}, path, separator string, value interface{}) (err error) {
	// If we have a property - update it, otherwise add it
	if _, err = GetProperty(original_data, path, separator); err != nil {
		err = AddProperty(original_data, path, separator, value)
	} else {
		if len(path) == 0 {
			path = separator
		}

		levels_tmp := strings.Split(path, separator)
		levels := make([]string, 0)
		for _, level_tmp := range levels_tmp {
			if len(level_tmp) > 0 {
				levels = append(levels, level_tmp)
			}
		}

		if len(levels) > 0 && path != separator {
			path_level_one := levels[0]

			// If we have a level in path_level_one

			re := regexp.MustCompile(`\w+\[\d+\]{1}`)
			if matched := re.FindString(path_level_one); len(matched) > 0 {
				property_re := regexp.MustCompile(`\w+`)
				index_re := regexp.MustCompile(`\[\d+\]{1}`)
				property := property_re.FindString(path_level_one)

				// Get an index
				index_found := index_re.FindString(path_level_one)

				// If index > 0 - check if this property is array
				if len(index_found) > 0 {
					if len(property) > 0 {
						path_level_one = property
					}
					index_found = strings.Trim(index_found, "[]")
					if index, err := strconv.Atoi(index_found); err == nil {
						if v, ok := original_data[property]; ok {
							if IsKind(v, reflect.Slice) {
								slice := reflect.ValueOf(v)
								var dest_value interface{}
								if index >= 0 && index < slice.Len() {
									dest_value = slice.Index(index).Interface()
								}
								// If len of other levels greater than 0
								if len(levels[1:]) >= 1 {
									if IsKind(dest_value, reflect.Map) {
										mapped_value := dest_value.(map[string]interface{})
										err = UpdateProperty(mapped_value, strings.Join(levels[1:], separator), separator, value)
										return err
									}
								} else {
									// if this is a `property[1]` in a path like `path.to.property[1]`
									slices := make([]interface{}, 0)
									for i := 0; i < slice.Len(); i++ {
										if i != index {
											slices = append(slices, slice.Index(i).Interface())
										} else {
											slices = append(slices, value)
										}
									}
									if index >= slice.Len() {
										slices = append(slices, value)
									}

									original_data[path_level_one] = slices
									return err
								}
							}
						} else {
							err = fmt.Errorf(
								"%s: is not an array", property,
							)
						}
					} else {
						err = fmt.Errorf(
							"Property %s does not exist", property,
						)
					}
				} else {
					err = fmt.Errorf(
						"%s must be of type %s",
						fmt.Sprintf("%s[%d]", property, index_found),
						"number",
					)
				}
			}

			if len(levels[1:]) >= 1 {
				if level_one_value, ok := original_data[path_level_one]; ok {
					if level_one_value != nil {
						switch reflect.TypeOf(level_one_value).Kind() {
						case reflect.Map:
							if mapped_level_one_value, ok := level_one_value.(map[string]interface{}); ok {
								err = UpdateProperty(mapped_level_one_value, strings.Join(levels[1:], separator), separator, value)
								if err != nil {
									return
								}
							}
						default:
							original_data[path] = value
						}
					}
				} else {
					err = fmt.Errorf(
						"Property %s does not exist", path_level_one,
					)
					return
				}
			} else {
				original_data[path_level_one] = value
			}
		} else if path == separator {
			original_data[path] = value
		}
	}

	return
}
