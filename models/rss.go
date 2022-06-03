package models

var Address = []string{"https://www.technologyreview.com/c/business/rss/", "https://www.technologyreview.com/c/biomedicine/rss/", "https://www.technologyreview.com/c/computing/rss/", "https://www.technologyreview.com/c/energy/rss/", "https://www.technologyreview.com/c/mobile/rss/", "https://www.technologyreview.com/c/robotics/rss/", "https://www.technologyreview.com/views/rss/", "https://www.technologyreview.com/stories.rss", "https://www.technologyreview.com/topnews.rss", "http://oasis.col.org/feed/rss_2.0/site", "http://oasis.col.org/feed/rss_1.0/site", "https://www.business-and-management.org/rss.php", "http://oasis.col.org/feed/rss_1.0/site", "http://oasis.col.org/feed/rss_2.0/site", "http://oasis.col.org/feed/atom_1.0/site", "https://www.tandfonline.com/feed/rss/uajb20", "https://www.tandfonline.com/feed/rss/uacn20", "https://www.tandfonline.com/feed/rss/hjsr20", "https://www.tandfonline.com/feed/rss/rfdj20", "https://www.tandfonline.com/feed/rss/rwrd20", "https://www.tandfonline.com/feed/rss/rffc20", "https://www.tandfonline.com/feed/rss/ujoa20", "https://www.tandfonline.com/feed/rss/cesr20", "https://www.tandfonline.com/feed/rss/rrse20", "https://www.tandfonline.com/feed/rss/iafd20", "https://www.tandfonline.com/feed/rss/iamy20", "https://www.tandfonline.com/feed/rss/iirp20", "https://www.tandfonline.com/feed/rss/kgmi20", "https://www.tandfonline.com/feed/rss/iwbp20", "https://www.tandfonline.com/feed/rss/ilab20", "https://www.tandfonline.com/feed/rss/wccq20", "https://www.tandfonline.com/feed/rss/mmis20", "https://www.tandfonline.com/feed/rss/ugit20", "https://www.tandfonline.com/feed/rss/rica20", "https://www.tandfonline.com/feed/rss/rics20", "https://www.tandfonline.com/feed/rss/hmep20", "https://www.tandfonline.com/feed/rss/idmr20", "https://www.tandfonline.com/feed/rss/itxc20", "https://www.tandfonline.com/feed/rss/ictx20", "https://www.tandfonline.com/feed/rss/upcp20", "https://www.tandfonline.com/feed/rss/fjps20", "https://www.tandfonline.com/feed/rss/ctrt20", "https://www.tandfonline.com/feed/rss/cjms20", "https://www.tandfonline.com/feed/rss/tsrm20", "https://www.tandfonline.com/feed/rss/reus20", "https://www.tandfonline.com/feed/rss/rwle20", "https://www.tandfonline.com/feed/rss/rsus20", "https://www.tandfonline.com/feed/rss/rcit20", "https://www.mdpi.com/rss/journal/admsci", "http://rss.sciencedirect.com/publication/science/24448834"}

type Enclosure struct {
	Url    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

type Item struct {
	Title     string    `xml:"title"`
	Link      string    `xml:"link"`
	Desc      string    `xml:"description"`
	City      string    `xml:"city"`
	Company   string    `xml:"company"`
	Logo      string    `xml:"logo"`
	JobType   string    `xml:"jobtype"`
	Category  string    `xml:"category"`
	PubDate   string    `xml:"date"`
	Enclosure Enclosure `xml:"enclosure"`
}

type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Desc  string `xml:"description"`
	Items []Item `xml:"item"`
}

type Rss struct {
	Channel Channel `xml:"channel"`
}
