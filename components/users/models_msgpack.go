package users

import (
	"fmt"
	"gopkg.in/vmihailenco/msgpack.v2"
	"reflect"
)

func init() {
	msgpack.Register(reflect.TypeOf(User{}), encodeUser, decodeUser)
	msgpack.Register(reflect.TypeOf(AccessToken{}), encodeAccessToken, decodeAccessToken)
}

func encodeUser(e *msgpack.Encoder, v reflect.Value) error {
    u := v.Interface().(User)
    if err := e.EncodeArrayLen(3); err != nil {
        return err
    }
    if err := e.EncodeUint(u.ID); err != nil {
        return err
    }
    if err := e.EncodeString(u.Username); err != nil {
        return err
    }
    if err := e.EncodeString(u.Email); err != nil {
        return err
    }
    return nil
}

func decodeUser(d *msgpack.Decoder, v reflect.Value) error {
    var err error
    var l int
    u := v.Addr().Interface().(*User)
    if l, err = d.DecodeArrayLen(); err != nil {
        return err
    }
    if l != 3 {
        return fmt.Errorf("array len doesn't match: %d", l)
    }
    if u.ID, err = d.DecodeUint(); err != nil {
        return err
    }
    if u.Username, err = d.DecodeString(); err != nil {
        return err
    }
    if u.Email, err = d.DecodeString(); err != nil {
        return err
    }
    return nil
}

func encodeAccessToken(e *msgpack.Encoder, v reflect.Value) error {
	at := v.Interface().(AccessToken)
	if err := e.EncodeArrayLen(4); err != nil {
		return err
	}
	if err := e.EncodeString(at.Token); err != nil {
		return err
	}
    if err := e.EncodeUint(at.User.ID); err != nil {
        return err
    }
    if err := e.EncodeString(at.User.Username); err != nil {
        return err
    }
    if err := e.EncodeString(at.User.Email); err != nil {
        return err
    }
	return nil
}

func decodeAccessToken(d *msgpack.Decoder, v reflect.Value) error {
	var err error
	var l int
	at := v.Addr().Interface().(*AccessToken)
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	if l != 4 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}
	if at.Token, err = d.DecodeString(); err != nil {
		return err
	}
    if at.User.ID, err = d.DecodeUint(); err != nil {
        return err
    }
    if at.User.Username, err = d.DecodeString(); err != nil {
        return err
    }
    if at.User.Email, err = d.DecodeString(); err != nil {
        return err
    }
	return nil
}
