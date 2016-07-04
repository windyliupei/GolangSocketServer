package xmlConfig

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

//Get Host address and prot frm configuration file.
func GetConfig() map[string]string {

	var token xml.Token
	var err error

	configMap := make(map[string]string)

	currentPath, err := os.Getwd()
	//CheckError(err)

	configByte, err := ioutil.ReadFile(currentPath + "/config.xml")

	decoder := xml.NewDecoder(bytes.NewBuffer(configByte))

	var elementCount int

	var keyList []string
	var valueList []string

	for token, err = decoder.Token(); err == nil; token, err = decoder.Token() {

		switch token := token.(type) {
		// 处理元素
		case xml.StartElement:
			name := token.Name.Local
			keyList = append(keyList, name)
			fmt.Println(keyList)
			elementCount++
		case xml.CharData:
			content := string([]byte(token))
			if content != "\n" {
				valueList = append(valueList, content)
				fmt.Println(valueList)
			}

		}
	}
	for index := 0; index < elementCount; index++ {
		configMap[keyList[index]] = valueList[index]
	}
	fmt.Println(elementCount)

	return configMap
}

func GetElement(key string, themap map[string]string) string {
	if value, ok := themap[key]; ok {
		return value
	}

	return ""
}
