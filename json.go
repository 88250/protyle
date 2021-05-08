package protyle

import jsoniter "github.com/json-iterator/go"

func UnmarshalJSON(data []byte, v interface{}) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, v)
}

func MarshalJSON(v interface{}) ([]byte, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(v)
}

func MarshalIndentJSON(v interface{}, prefix, indent string) ([]byte, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.MarshalIndent(v, prefix, indent)
}
