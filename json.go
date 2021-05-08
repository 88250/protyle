// Protyle - 解析渲染思源笔记文档数据
// Copyright (c) 2021-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

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
