package lbp_xml

import (
	"encoding/xml"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/lbp_xml/npdata"
)

//<news>
//    <subcategory>
//        <id></id>
//        <title></title>
//        <item>
//            <id></id>
//            <subject></subject>
//            <content>
//                <!--There can be multiple frame tags-->
//                <frame width="">
//                    <title></title>
//                    <item width="">
//                        <slot type="">
//                            <id></id>
//                        </slot>
//                        <npHandle icon=""></npHandle>
//                        <content></content>
//                    </item>
//                </frame>
//            </content>
//            <timestamp></timestamp>
//        </item>
//    </subcategory>
//</news>

type News struct {
	XMLName    xml.Name      `xml:"news"`
	Categories []SubCategory `xml:",innerxml"`
}

type SubCategory struct {
	XMLName xml.Name   `xml:"subcategory"`
	ID      string     `xml:"id"`
	Title   string     `xml:"title"`
	Items   []NewsItem `xml:",innerxml"`
}
type NewsItem struct {
	XMLName   xml.Name        `xml:"item"`
	ID        string          `xml:"id"`
	Subject   string          `xml:"subject"`
	Content   NewsItemContent `xml:"content"`
	Timestamp int64           `xml:"timestamp"`
}

type NewsItemContent struct {
	Content []NewsFrame
}

type NewsFrame struct {
	XMLName xml.Name        `xml:"frame"`
	Width   string          `xml:"width,attr"`
	Title   string          `xml:"title"`
	Item    []NewsFrameItem `xml:"item"`
}
type NewsFrameItem struct {
	Width    string             `xml:"width,attr"`
	NpHandle *npdata.NpHandle   `xml:"npHandle,omitempty"`
	Slot     *NewsFrameItemSlot `xml:"slot,omitempty"`
	Content  string             `xml:"content"`
}
type NewsFrameItemSlot struct {
	Type string `xml:"type,attr"`
	ID   uint64 `xml:"id"`
}
