package main

import "encoding/xml"

// Objs was generated 2020-02-03 16:37:40 by riccardo on HackBookPro.local.
type Objs struct {
	XMLName xml.Name `xml:"Objs"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"Version,attr"`
	Xmlns   string   `xml:"xmlns,attr"`
	Obj     []struct {
		Text  string `xml:",chardata"`
		RefId string `xml:"RefId,attr"`
		TN    struct {
			Text  string   `xml:",chardata"`
			RefId string   `xml:"RefId,attr"`
			T     []string `xml:"T"`
		} `xml:"TN"`
		MS struct {
			Text string `xml:",chardata"`
			Obj  []struct {
				Text  string `xml:",chardata"`
				N     string `xml:"N,attr"`
				RefId string `xml:"RefId,attr"`
				TNRef struct {
					Text  string `xml:",chardata"`
					RefId string `xml:"RefId,attr"`
				} `xml:"TNRef"`
				MS struct {
					Text string `xml:",chardata"`
					S    []struct {
						Text string `xml:",chardata"`
						N    string `xml:"N,attr"`
					} `xml:"S"`
					Obj []struct {
						Text  string `xml:",chardata"`
						N     string `xml:"N,attr"`
						RefId string `xml:"RefId,attr"`
						TN    struct {
							Text  string   `xml:",chardata"`
							RefId string   `xml:"RefId,attr"`
							T     []string `xml:"T"`
						} `xml:"TN"`
						LST struct {
							Text string   `xml:",chardata"`
							S    []string `xml:"S"`
						} `xml:"LST"`
						TNRef struct {
							Text  string `xml:",chardata"`
							RefId string `xml:"RefId,attr"`
						} `xml:"TNRef"`
					} `xml:"Obj"`
				} `xml:"MS"`
			} `xml:"Obj"`
			S []struct {
				Text string `xml:",chardata"`
				N    string `xml:"N,attr"`
			} `xml:"S"`
		} `xml:"MS"`
		TNRef struct {
			Text  string `xml:",chardata"`
			RefId string `xml:"RefId,attr"`
		} `xml:"TNRef"`
	} `xml:"Obj"`
}
