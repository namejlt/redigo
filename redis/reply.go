// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package redis

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// ErrNil indicates that a reply value is nil.
var ErrNil = errors.New("redigo: nil returned")

// Int is a helper that converts a command reply to an integer. If err is not
// equal to nil, then Int returns 0, err. Otherwise, Int converts the
// reply to an int as follows:
//
//  Reply type    Result
//  integer       int(reply), nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Int(reply interface{}, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case int64:
		x := int(reply)
		if int64(x) != reply {
			return 0, strconv.ErrRange
		}
		return x, nil
	case []byte:
		n, err := strconv.ParseInt(string(reply), 10, 0)
		return int(n), err
	case nil:
		return 0, ErrNil
	case Error:
		return 0, reply
	}
	return 0, fmt.Errorf("redigo: unexpected type for Int, got type %T", reply)
}

// Int64 is a helper that converts a command reply to 64 bit integer. If err is
// not equal to nil, then Int64 returns 0, err. Otherwise, Int64 converts the
// reply to an int64 as follows:
//
//  Reply type    Result
//  integer       reply, nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Int64(reply interface{}, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case int64:
		return reply, nil
	case []byte:
		n, err := strconv.ParseInt(string(reply), 10, 64)
		return n, err
	case nil:
		return 0, ErrNil
	case Error:
		return 0, reply
	}
	return 0, fmt.Errorf("redigo: unexpected type for Int64, got type %T", reply)
}

func errNegativeInt(v int64) error {
	return fmt.Errorf("redigo: unexpected negative value %v for Uint64", v)
}

// Uint64 is a helper that converts a command reply to 64 bit unsigned integer.
// If err is not equal to nil, then Uint64 returns 0, err. Otherwise, Uint64 converts the
// reply to an uint64 as follows:
//
//  Reply type    Result
//  +integer      reply, nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Uint64(reply interface{}, err error) (uint64, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case int64:
		if reply < 0 {
			return 0, errNegativeInt(reply)
		}
		return uint64(reply), nil
	case []byte:
		n, err := strconv.ParseUint(string(reply), 10, 64)
		return n, err
	case nil:
		return 0, ErrNil
	case Error:
		return 0, reply
	}
	return 0, fmt.Errorf("redigo: unexpected type for Uint64, got type %T", reply)
}

// Float64 is a helper that converts a command reply to 64 bit float. If err is
// not equal to nil, then Float64 returns 0, err. Otherwise, Float64 converts
// the reply to a float64 as follows:
//
//  Reply type    Result
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Float64(reply interface{}, err error) (float64, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case []byte:
		n, err := strconv.ParseFloat(string(reply), 64)
		return n, err
	case nil:
		return 0, ErrNil
	case Error:
		return 0, reply
	}
	return 0, fmt.Errorf("redigo: unexpected type for Float64, got type %T", reply)
}

// String is a helper that converts a command reply to a string. If err is not
// equal to nil, then String returns "", err. Otherwise String converts the
// reply to a string as follows:
//
//  Reply type      Result
//  bulk string     string(reply), nil
//  simple string   reply, nil
//  nil             "",  ErrNil
//  other           "",  error
func String(reply interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}
	switch reply := reply.(type) {
	case []byte:
		return string(reply), nil
	case string:
		return reply, nil
	case nil:
		return "", ErrNil
	case Error:
		return "", reply
	}
	return "", fmt.Errorf("redigo: unexpected type for String, got type %T", reply)
}

// Bytes is a helper that converts a command reply to a slice of bytes. If err
// is not equal to nil, then Bytes returns nil, err. Otherwise Bytes converts
// the reply to a slice of bytes as follows:
//
//  Reply type      Result
//  bulk string     reply, nil
//  simple string   []byte(reply), nil
//  nil             nil, ErrNil
//  other           nil, error
func Bytes(reply interface{}, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	switch reply := reply.(type) {
	case []byte:
		return reply, nil
	case string:
		return []byte(reply), nil
	case nil:
		return nil, ErrNil
	case Error:
		return nil, reply
	}
	return nil, fmt.Errorf("redigo: unexpected type for Bytes, got type %T", reply)
}

// Bool is a helper that converts a command reply to a boolean. If err is not
// equal to nil, then Bool returns false, err. Otherwise Bool converts the
// reply to boolean as follows:
//
//  Reply type      Result
//  integer         value != 0, nil
//  bulk string     strconv.ParseBool(reply)
//  nil             false, ErrNil
//  other           false, error
func Bool(reply interface{}, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	switch reply := reply.(type) {
	case int64:
		return reply != 0, nil
	case []byte:
		return strconv.ParseBool(string(reply))
	case nil:
		return false, ErrNil
	case Error:
		return false, reply
	}
	return false, fmt.Errorf("redigo: unexpected type for Bool, got type %T", reply)
}

// MultiBulk is a helper that converts an array command reply to a []interface{}.
//
// Deprecated: Use Values instead.
func MultiBulk(reply interface{}, err error) ([]interface{}, error) { return Values(reply, err) }

// Values is a helper that converts an array command reply to a []interface{}.
// If err is not equal to nil, then Values returns nil, err. Otherwise, Values
// converts the reply as follows:
//
//  Reply type      Result
//  array           reply, nil
//  nil             nil, ErrNil
//  other           nil, error
func Values(reply interface{}, err error) ([]interface{}, error) {
	if err != nil {
		return nil, err
	}
	switch reply := reply.(type) {
	case []interface{}:
		return reply, nil
	case nil:
		return nil, ErrNil
	case Error:
		return nil, reply
	}
	return nil, fmt.Errorf("redigo: unexpected type for Values, got type %T", reply)
}

func sliceHelper(reply interface{}, err error, name string, makeSlice func(int), assign func(int, interface{}) error) error {
	if err != nil {
		return err
	}
	switch reply := reply.(type) {
	case []interface{}:
		makeSlice(len(reply))
		for i := range reply {
			if reply[i] == nil {
				continue
			}
			if err := assign(i, reply[i]); err != nil {
				return err
			}
		}
		return nil
	case nil:
		return ErrNil
	case Error:
		return reply
	}
	return fmt.Errorf("redigo: unexpected type for %s, got type %T", name, reply)
}

// Float64s is a helper that converts an array command reply to a []float64. If
// err is not equal to nil, then Float64s returns nil, err. Nil array items are
// converted to 0 in the output slice. Floats64 returns an error if an array
// item is not a bulk string or nil.
func Float64s(reply interface{}, err error) ([]float64, error) {
	var result []float64
	err = sliceHelper(reply, err, "Float64s", func(n int) { result = make([]float64, n) }, func(i int, v interface{}) error {
		switch v := v.(type) {
		case []byte:
			f, err := strconv.ParseFloat(string(v), 64)
			result[i] = f
			return err
		case Error:
			return v
		default:
			return fmt.Errorf("redigo: unexpected element type for Float64s, got type %T", v)
		}
	})
	return result, err
}

// Strings is a helper that converts an array command reply to a []string. If
// err is not equal to nil, then Strings returns nil, err. Nil array items are
// converted to "" in the output slice. Strings returns an error if an array
// item is not a bulk string or nil.
func Strings(reply interface{}, err error) ([]string, error) {
	var result []string
	err = sliceHelper(reply, err, "Strings", func(n int) { result = make([]string, n) }, func(i int, v interface{}) error {
		switch v := v.(type) {
		case string:
			result[i] = v
			return nil
		case []byte:
			result[i] = string(v)
			return nil
		case Error:
			return v
		default:
			return fmt.Errorf("redigo: unexpected element type for Strings, got type %T", v)
		}
	})
	return result, err
}

// ByteSlices is a helper that converts an array command reply to a [][]byte.
// If err is not equal to nil, then ByteSlices returns nil, err. Nil array
// items are stay nil. ByteSlices returns an error if an array item is not a
// bulk string or nil.
func ByteSlices(reply interface{}, err error) ([][]byte, error) {
	var result [][]byte
	err = sliceHelper(reply, err, "ByteSlices", func(n int) { result = make([][]byte, n) }, func(i int, v interface{}) error {
		switch v := v.(type) {
		case []byte:
			result[i] = v
			return nil
		case Error:
			return v
		default:
			return fmt.Errorf("redigo: unexpected element type for ByteSlices, got type %T", v)
		}
	})
	return result, err
}

// Int64s is a helper that converts an array command reply to a []int64.
// If err is not equal to nil, then Int64s returns nil, err. Nil array
// items are stay nil. Int64s returns an error if an array item is not a
// bulk string or nil.
func Int64s(reply interface{}, err error) ([]int64, error) {
	var result []int64
	err = sliceHelper(reply, err, "Int64s", func(n int) { result = make([]int64, n) }, func(i int, v interface{}) error {
		switch v := v.(type) {
		case int64:
			result[i] = v
			return nil
		case []byte:
			n, err := strconv.ParseInt(string(v), 10, 64)
			result[i] = n
			return err
		case Error:
			return v
		default:
			return fmt.Errorf("redigo: unexpected element type for Int64s, got type %T", v)
		}
	})
	return result, err
}

// Ints is a helper that converts an array command reply to a []int.
// If err is not equal to nil, then Ints returns nil, err. Nil array
// items are stay nil. Ints returns an error if an array item is not a
// bulk string or nil.
func Ints(reply interface{}, err error) ([]int, error) {
	var result []int
	err = sliceHelper(reply, err, "Ints", func(n int) { result = make([]int, n) }, func(i int, v interface{}) error {
		switch v := v.(type) {
		case int64:
			n := int(v)
			if int64(n) != v {
				return strconv.ErrRange
			}
			result[i] = n
			return nil
		case []byte:
			n, err := strconv.Atoi(string(v))
			result[i] = n
			return err
		case Error:
			return v
		default:
			return fmt.Errorf("redigo: unexpected element type for Ints, got type %T", v)
		}
	})
	return result, err
}

// mapHelper builds a map from the data in reply.
func mapHelper(reply interface{}, err error, name string, makeMap func(int), assign func(key string, value interface{}) error) error {
	values, err := Values(reply, err)
	if err != nil {
		return err
	}

	if len(values)%2 != 0 {
		return fmt.Errorf("redigo: %s expects even number of values result, got %d", name, len(values))
	}

	makeMap(len(values) / 2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].([]byte)
		if !ok {
			return fmt.Errorf("redigo: %s key[%d] not a bulk string value, got %T", name, i, values[i])
		}

		if err := assign(string(key), values[i+1]); err != nil {
			return err
		}
	}

	return nil
}

// StringMap is a helper that converts an array of strings (alternating key, value)
// into a map[string]string. The HGETALL and CONFIG GET commands return replies in this format.
// Requires an even number of values in result.
func StringMap(reply interface{}, err error) (map[string]string, error) {
	var result map[string]string
	err = mapHelper(reply, err, "StringMap",
		func(n int) {
			result = make(map[string]string, n)
		}, func(key string, v interface{}) error {
			value, ok := v.([]byte)
			if !ok {
				return fmt.Errorf("redigo: StringMap for %q not a bulk string value, got %T", key, v)
			}

			result[key] = string(value)

			return nil
		},
	)

	return result, err
}

// IntMap is a helper that converts an array of strings (alternating key, value)
// into a map[string]int. The HGETALL commands return replies in this format.
// Requires an even number of values in result.
func IntMap(result interface{}, err error) (map[string]int, error) {
	var m map[string]int
	err = mapHelper(result, err, "IntMap",
		func(n int) {
			m = make(map[string]int, n)
		}, func(key string, v interface{}) error {
			value, err := Int(v, nil)
			if err != nil {
				return err
			}

			m[key] = value

			return nil
		},
	)

	return m, err
}

// Int64Map is a helper that converts an array of strings (alternating key, value)
// into a map[string]int64. The HGETALL commands return replies in this format.
// Requires an even number of values in result.
func Int64Map(result interface{}, err error) (map[string]int64, error) {
	var m map[string]int64
	err = mapHelper(result, err, "Int64Map",
		func(n int) {
			m = make(map[string]int64, n)
		}, func(key string, v interface{}) error {
			value, err := Int64(v, nil)
			if err != nil {
				return err
			}

			m[key] = value

			return nil
		},
	)

	return m, err
}

// Float64Map is a helper that converts an array of strings (alternating key, value)
// into a map[string]float64. The HGETALL commands return replies in this format.
// Requires an even number of values in result.
func Float64Map(result interface{}, err error) (map[string]float64, error) {
	var m map[string]float64
	err = mapHelper(result, err, "Float64Map",
		func(n int) {
			m = make(map[string]float64, n)
		}, func(key string, v interface{}) error {
			value, err := Float64(v, nil)
			if err != nil {
				return err
			}

			m[key] = value

			return nil
		},
	)

	return m, err
}

// Positions is a helper that converts an array of positions (lat, long)
// into a [][2]float64. The GEOPOS command returns replies in this format.
func Positions(result interface{}, err error) ([]*[2]float64, error) {
	values, err := Values(result, err)
	if err != nil {
		return nil, err
	}
	positions := make([]*[2]float64, len(values))
	for i := range values {
		if values[i] == nil {
			continue
		}

		p, ok := values[i].([]interface{})
		if !ok {
			return nil, fmt.Errorf("redigo: unexpected element type for interface slice, got type %T", values[i])
		}

		if len(p) != 2 {
			return nil, fmt.Errorf("redigo: unexpected number of values for a member position, got %d", len(p))
		}

		lat, err := Float64(p[0], nil)
		if err != nil {
			return nil, err
		}

		long, err := Float64(p[1], nil)
		if err != nil {
			return nil, err
		}

		positions[i] = &[2]float64{lat, long}
	}
	return positions, nil
}

// Uint64s is a helper that converts an array command reply to a []uint64.
// If err is not equal to nil, then Uint64s returns nil, err. Nil array
// items are stay nil. Uint64s returns an error if an array item is not a
// bulk string or nil.
func Uint64s(reply interface{}, err error) ([]uint64, error) {
	var result []uint64
	err = sliceHelper(reply, err, "Uint64s", func(n int) { result = make([]uint64, n) }, func(i int, v interface{}) error {
		switch v := v.(type) {
		case uint64:
			result[i] = v
			return nil
		case []byte:
			n, err := strconv.ParseUint(string(v), 10, 64)
			result[i] = n
			return err
		case Error:
			return v
		default:
			return fmt.Errorf("redigo: unexpected element type for Uint64s, got type %T", v)
		}
	})
	return result, err
}

// Uint64Map is a helper that converts an array of strings (alternating key, value)
// into a map[string]uint64. The HGETALL commands return replies in this format.
// Requires an even number of values in result.
func Uint64Map(result interface{}, err error) (map[string]uint64, error) {
	var m map[string]uint64
	err = mapHelper(result, err, "Uint64Map",
		func(n int) {
			m = make(map[string]uint64, n)
		}, func(key string, v interface{}) error {
			value, err := Uint64(v, nil)
			if err != nil {
				return err
			}

			m[key] = value

			return nil
		},
	)

	return m, err
}

// SlowLogs is a helper that parse the SLOWLOG GET command output and
// return the array of SlowLog
func SlowLogs(result interface{}, err error) ([]SlowLog, error) {
	rawLogs, err := Values(result, err)
	if err != nil {
		return nil, err
	}
	logs := make([]SlowLog, len(rawLogs))
	for i, e := range rawLogs {
		rawLog, ok := e.([]interface{})
		if !ok {
			return nil, fmt.Errorf("redigo: slowlog element is not an array, got %T", e)
		}

		var log SlowLog
		if len(rawLog) < 4 {
			return nil, fmt.Errorf("redigo: slowlog element has %d elements, expected at least 4", len(rawLog))
		}

		log.ID, ok = rawLog[0].(int64)
		if !ok {
			return nil, fmt.Errorf("redigo: slowlog element[0] not an int64, got %T", rawLog[0])
		}

		timestamp, ok := rawLog[1].(int64)
		if !ok {
			return nil, fmt.Errorf("redigo: slowlog element[1] not an int64, got %T", rawLog[1])
		}

		log.Time = time.Unix(timestamp, 0)
		duration, ok := rawLog[2].(int64)
		if !ok {
			return nil, fmt.Errorf("redigo: slowlog element[2] not an int64, got %T", rawLog[2])
		}

		log.ExecutionTime = time.Duration(duration) * time.Microsecond

		log.Args, err = Strings(rawLog[3], nil)
		if err != nil {
			return nil, fmt.Errorf("redigo: slowlog element[3] is not array of strings: %w", err)
		}

		if len(rawLog) >= 6 {
			log.ClientAddr, err = String(rawLog[4], nil)
			if err != nil {
				return nil, fmt.Errorf("redigo: slowlog element[4] is not a string: %w", err)
			}

			log.ClientName, err = String(rawLog[5], nil)
			if err != nil {
				return nil, fmt.Errorf("redigo: slowlog element[5] is not a string: %w", err)
			}
		}
		logs[i] = log
	}
	return logs, nil
}
