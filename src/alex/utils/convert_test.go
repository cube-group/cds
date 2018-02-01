package utils

import (
    "testing"
    "encoding/xml"
    "fmt"
)

type MyXml struct {
    XMLName xml.Name `xml:"root"`
    Items   []MyXmlItem `xml:"item"`
}

type MyXmlItem struct {
    XMLName   xml.Name `xml:"item"`
    Name      string `xml:"name,attr"`
    Level     string `xml:"level,attr"`
    InnerText string `xml:",innerxml"`
}

type MyObject struct {
    Name     string `json:"name"`
    Value    string `json:"value"`
    Children []MyObject `json:"children"`
}

type StructMap struct {
    Name  string
    Value string
}

func TestJsonDecode(t *testing.T) {
    jsonString := "{\"name\":\"linyang\",\"value\":\"123\",\"children\":[{\"name\":\"linyang\",\"value\":\"123\"}]}"
    jsonObject := MyObject{}
    if _, err := JsonDecode(jsonString, &jsonObject); err != nil {
        t.Errorf("TestJsonDecode Error %v", err)
    } else {
        t.Logf("TestJsonDecode Success %v", jsonObject)
    }
}

func TestJsonEncode(t *testing.T) {
    jsonObject2 := MyObject{Name:"root", Value:"1111"}
    jsonObject2Items := []MyObject{{"linyang", "12312313", []MyObject{}}}
    jsonObject2.Children = append([]MyObject{}, jsonObject2Items...)
    if result, err := JsonEncode(jsonObject2); err != nil {
        t.Errorf("TestJsonDecode Error %v", err)
    } else {
        fmt.Printf("TestJsonEncode Result %v", result)
    }
}

func TestXmlDecode(t *testing.T) {
    xmlString := "<root><item name='linyang' age='3'>哈哈</item><item name='linyang' age='3'>呵呵</item></root>"
    xmlObject := MyXml{}
    if _, err := XmlDecode(xmlString, &xmlObject); err != nil {
        t.Errorf("TestXmlDecode Error %v", err)
    } else {
        t.Logf("TestXmlDecode Success %v", xmlObject)
    }
}

func TestXmlEncode(t *testing.T) {
    jsonObject1 := MyObject{Name:"root", Value:"1111"}
    jsonObject1Items := []MyObject{{"linyang", "12312313", []MyObject{}}}
    jsonObject1.Children = append([]MyObject{}, jsonObject1Items...)
    if result, err := JsonEncode(jsonObject1); err != nil {
        t.Errorf("TestXmlEncode1 Error %v", err)
    } else {
        t.Logf("TestXmlEncode1 Success %v", result)
    }

    jsonObject2 := []interface{}{"Name", "root", "Value", 1111}
    if result2, err2 := JsonEncode(jsonObject2); err2 != nil {
        t.Errorf("TestXmlEncode2 Error %v", err2)
    } else {
        t.Logf("TestXmlEncode2 Success %v", result2)
    }
}

func TestGetStructMapData(t *testing.T) {
    s := &StructMap{Name:"test", Value:"value"}
    m, err := GetStructMapData(*s)
    if err != nil {
        t.Errorf("TestGetStructMapData Error %v", err)
    } else {
        fmt.Println(m)
    }
}

func TestSetStructData(t *testing.T) {
    s := &StructMap{Name:"test", Value:"value"}
    err := SetStructData(s, map[string]string{"aaa":"111", "Name":"aaa", "Value":"bbbb"})
    if err != nil {
        t.Errorf("TestSetStructData Error %v", err)
    } else {
        fmt.Println(s)
    }
}