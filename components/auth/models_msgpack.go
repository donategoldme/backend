package auth


//func init() {
//    msgpack.Register(reflect.TypeOf(Provider{}), encodeProvider, decodeProvider)
//}
//
//func encodeProvider(e *msgpack.Encoder, v reflect.Value) error {
//	p := v.Interface().(Provider)
//	if err := e.EncodeArrayLen(4); err != nil {
//		return err
//	}
//    if err := e.EncodeString(p.UID); err != nil {
//		return err
//	}
//	if err := e.EncodeUint(p.UserID); err != nil {
//		return err
//	}
//	if err := e.EncodeString(p.TypeProvider); err != nil {
//		return err
//	}
//    if err := e.EncodeString(p.AccessToken); err != nil {
//		return err
//	}
//	return nil
//}
//
//func decodeProvider(d *msgpack.Decoder, v reflect.Value) error {
//	var err error
//	var l int
//	p := v.Addr().Interface().(*Provider)
//	if l, err = d.DecodeArrayLen(); err != nil {
//		return err
//	}
//	if l != 4 {
//		return fmt.Errorf("array len doesn't match: %d", l)
//	}
//	if p.UID, err = d.DecodeString(); err != nil {
//		return err
//	}
//	return nil
//}
