package mod_xml

import (
	"encoding/xml"
)

type GenericXMLElement struct {
	XMLName    xml.Name
	Attrs      []xml.Attr          `xml:",any,attr"`
	InnerText  string              `xml:",chardata"`
	ChildNodes []GenericXMLElement `xml:",any"`
}
