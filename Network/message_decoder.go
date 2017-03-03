package Network

import(
	"encoding/json"
)




func Struct_to_json_MainData(struct_object MainData) []byte {
	json_object, _ := json.Marshal(struct_object)

	return json_object
}



func Json_to_struct_MainData(receivedMessageBytes []byte) MainData {
    struct_object := MainData{}
    json.Unmarshal(receivedMessageBytes, &struct_object)

    return struct_object
}

func Struct_to_json_Button(struct_object Button) []byte {
	json_object, _ := json.Marshal(struct_object)

	return json_object
}



func Json_to_struct_Button(receivedMessageBytes []byte) Button {
    struct_object := Button{}
    json.Unmarshal(receivedMessageBytes, &struct_object)

    return struct_object
}




