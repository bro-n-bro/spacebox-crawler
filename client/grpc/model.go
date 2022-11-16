package grpc

//
//// FindEventByType searches inside the given tx events for the message having the specified index, in order
//// to find the event having the given type, and returns it.
//// If no such event is found, returns an error instead.
//func (tx Tx) FindEventByType(index int, eventType string) (sdk.StringEvent, error) {
//	for _, ev := range tx.Logs[index].Events {
//		if ev.Type == eventType {
//			return ev, nil
//		}
//	}
//
//	return sdk.StringEvent{}, fmt.Errorf("no %s event found inside tx with hash %s", eventType, tx.TxHash)
//}
//
//// FindAttributeByKey searches inside the specified event of the given tx to find the attribute having the given key.
//// If the specified event does not contain a such attribute, returns an error instead.
//func (tx Tx) FindAttributeByKey(event sdk.StringEvent, attrKey string) (string, error) {
//	for _, attr := range event.Attributes {
//		if attr.Key == attrKey {
//			return attr.Value, nil
//		}
//	}
//
//	return "", fmt.Errorf("no event with attribute %s found inside tx with hash %s", attrKey, tx.TxHash)
//}
