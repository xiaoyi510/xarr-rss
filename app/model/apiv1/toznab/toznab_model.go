package toznab

import "encoding/xml"

///////////////////////////////////////////////////////////////////////////

type Apiv1ToznabCapsRes struct {
	XMLName    xml.Name                  `xml:"caps"`
	Text       string                    `xml:",chardata"`
	Server     Apiv1ToznabCapsServer     `xml:"server"`
	Limits     Apiv1ToznabCapsLimits     `xml:"limits"`
	Searching  Apiv1ToznabCapsSearching  `xml:"searching"`
	Categories Apiv1ToznabCapsCategories `xml:"categories"`
	Tags       Apiv1ToznabCapsTags       `xml:"tags"`
}

type Apiv1ToznabCapsServer struct {
	Text  string `xml:",chardata"`
	Title string `xml:"title,attr"`
}

type Apiv1ToznabCapsLimits struct {
	Text    string `xml:",chardata"`
	Max     string `xml:"max,attr"`
	Default string `xml:"default,attr"`
}
type Apiv1ToznabCapsSearchingSearch struct {
	Text            string `xml:",chardata"`
	Available       string `xml:"available,attr"`
	SupportedParams string `xml:"supportedParams,attr"`
}
type Apiv1ToznabCapsSearchingTvSearch struct {
	Text            string `xml:",chardata"`
	Available       string `xml:"available,attr"`
	SupportedParams string `xml:"supportedParams,attr"`
}
type Apiv1ToznabCapsSearchingMovieSearch struct {
	Text            string `xml:",chardata"`
	Available       string `xml:"available,attr"`
	SupportedParams string `xml:"supportedParams,attr"`
}
type Apiv1ToznabCapsSearchingMusicSearch struct {
	//Text            string `xml:",chardata"`
	Available       string `xml:"available,attr"`
	SupportedParams string `xml:"supportedParams,attr"`
}
type Apiv1ToznabCapsSearchingAudioSearch struct {
	Text            string `xml:",chardata"`
	Available       string `xml:"available,attr"`
	SupportedParams string `xml:"supportedParams,attr"`
}
type Apiv1ToznabCapsSearchingBookSearch struct {
	Text            string `xml:",chardata"`
	Available       string `xml:"available,attr"`
	SupportedParams string `xml:"supportedParams,attr"`
}
type Apiv1ToznabCapsSearching struct {
	Text        string                              `xml:",chardata"`
	Search      Apiv1ToznabCapsSearchingSearch      `xml:"search"`
	TvSearch    Apiv1ToznabCapsSearchingTvSearch    `xml:"tv-search"`
	MovieSearch Apiv1ToznabCapsSearchingMovieSearch `xml:"movie-search"`
	MusicSearch Apiv1ToznabCapsSearchingMusicSearch `xml:"music-search"`
	AudioSearch Apiv1ToznabCapsSearchingAudioSearch `xml:"audio-search"`
	BookSearch  Apiv1ToznabCapsSearchingBookSearch  `xml:"book-search"`
}

type Apiv1ToznabCapsCategories struct {
	Text     string                              `xml:",chardata"`
	Category []Apiv1ToznabCapsCategoriesCategory `xml:"category"`
}
type Apiv1ToznabCapsCategoriesCategory struct {
	Text   string              `xml:",chardata"`
	ID     string              `xml:"id,attr"`
	Name   string              `xml:"name,attr"`
	Subcat []Apiv1ToznabSubCat `xml:"subcat,omitempty"`
}
type Apiv1ToznabCapsTags struct {
	Text string `xml:",chardata"`
}
type Apiv1ToznabSubCat struct {
	Text string `xml:",chardata"`
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

///////////////////////////////////////////////////////////////////////////

type Apiv1ToznabTvSearchRes struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Atom    string   `xml:"atom,attr"`
	Torznab string   `xml:"torznab,attr"`
	Channel struct {
		Text string `xml:",chardata"`
		Link struct {
			Text string `xml:",chardata"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Title struct {
			Text string `xml:",chardata"`
		} `xml:"title"`
		Item []struct {
			Text  string `xml:",chardata"`
			Title struct {
				Text string `xml:",chardata"`
			} `xml:"title"`
			Guid struct {
				Text string `xml:",chardata"`
			} `xml:"guid"`
			XArrRssindexer struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
			} `xml:"XArrRssIndexer"`
			Comments struct {
				Text string `xml:",chardata"`
			} `xml:"comments"`
			PubDate struct {
				Text string `xml:",chardata"`
			} `xml:"pubDate"`
			Size struct {
				Text string `xml:",chardata"`
			} `xml:"size"`
			Link struct {
				Text string `xml:",chardata"`
			} `xml:"link"`
			Category struct {
				Text string `xml:",chardata"`
			} `xml:"category"`
			Enclosure struct {
				Text   string `xml:",chardata"`
				URL    string `xml:"url,attr"`
				Length string `xml:"length,attr"`
				Type   string `xml:"type,attr"`
			} `xml:"enclosure"`
			Attr []struct {
				Text  string `xml:",chardata"`
				Name  string `xml:"name,attr"`
				Value string `xml:"value,attr"`
			} `xml:"attr"`
		} `xml:"item"`
	} `xml:"channel"`
}

type Apiv1ToznabSearchChannel struct {
	Text  string                         `xml:",chardata"`
	Link  Apiv1ToznabSearchChannelLink   `xml:"atom:link"`
	Title Apiv1ToznabSearchChannelTitle  `xml:"title"`
	Item  []Apiv1ToznabSearchChannelItem `xml:"item"`
}
type Apiv1ToznabSearchChannelLink struct {
	Text string `xml:",chardata"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}
type Apiv1ToznabSearchChannelTitle struct {
	Text string `xml:",chardata"`
}
type Apiv1ToznabSearchChannelItem struct {
	Text           string                                     `xml:",chardata"`
	Title          Apiv1ToznabSearchChannelItemTitle          `xml:"title"`
	OriginalTitle  Apiv1ToznabSearchChannelItemTitle          `xml:"originalTitle"`
	OtherTitle     Apiv1ToznabSearchChannelItemTitle          `xml:"OtherTitle"`
	Guid           Apiv1ToznabSearchChannelItemGuid           `xml:"guid"`
	XArrRssindexer Apiv1ToznabSearchChannelItemXArrRssindexer `xml:"XArrRssIndexer"`
	Comments       Apiv1ToznabSearchChannelItemComments       `xml:"comments"`
	PubDate        Apiv1ToznabSearchChannelItemPubDate        `xml:"pubDate"`
	Size           Apiv1ToznabSearchChannelItemSize           `xml:"size"`
	Link           Apiv1ToznabSearchChannelItemLink           `xml:"link"`
	Category       Apiv1ToznabSearchChannelItemCategory       `xml:"category"`
	Enclosure      Apiv1ToznabSearchChannelItemEnclosure      `xml:"enclosure"`
	Attr           []Apiv1ToznabSearchChannelItemAttr         `xml:"torznab:attr"`
}

type Apiv1ToznabSearchChannelItemTitle struct {
	Text string `xml:",chardata"`
}
type Apiv1ToznabSearchChannelItemGuid struct {
	Text string `xml:",chardata"`
}
type Apiv1ToznabSearchChannelItemXArrRssindexer struct {
	Text string `xml:",chardata"`
	ID   string `xml:"id,attr"`
}
type Apiv1ToznabSearchChannelItemComments struct {
	Text string `xml:",chardata"`
}
type Apiv1ToznabSearchChannelItemPubDate struct {
	Text string `xml:",chardata"`
}
type Apiv1ToznabSearchChannelItemSize struct {
	Text string `xml:",chardata"`
}
type Apiv1ToznabSearchChannelItemLink struct {
	Text string `xml:",chardata"`
}
type Apiv1ToznabSearchChannelItemCategory struct {
	Text string `xml:",chardata"`
}
type Apiv1ToznabSearchChannelItemEnclosure struct {
	Text   string `xml:",chardata"`
	URL    string `xml:"url,attr"`
	Length string `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}
type Apiv1ToznabSearchChannelItemAttr struct {
	Text  string `xml:",chardata"`
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Apiv1ToznabSearchRes struct {
	XMLName xml.Name                 `xml:"rss"`
	Text    string                   `xml:",chardata"`
	Version string                   `xml:"version,attr"`
	Atom    string                   `xml:"xmlns:atom,attr"`
	Torznab string                   `xml:"xmlns:torznab,attr"`
	Channel Apiv1ToznabSearchChannel `xml:"channel"`
}
