package Network

import(
	"encoding/json"
)




func Struct_to_json(struct_object MainData) []byte {
	json_object, _ := json.Marshal(struct_object)

	return json_object
}



func Json_to_struct(receivedMessageBytes []byte) MainData {
    struct_object := MainData{}
    json.Unmarshal(receivedMessageBytes, &struct_object)

    return struct_object
}




