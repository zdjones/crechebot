package main

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestHtml(t *testing.T) {
	s := `<!DOCTYPE html
    PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">

<html xmlns="http://www.w3.org/1999/xhtml">

<head id="Head1">
    <title>
        Web Bookings
    </title>
    <link href="/sitecss/gllbetter/thickbox.css?v=63688688113" rel="stylesheet" type="text/css" />
    <link href="/sitecss/gllbetter/jtip.css?v=63688688113" rel="stylesheet" type="text/css" />
    <link href="/sitecss/gllbetter/BookingsCss.css?v=63688688113" rel="stylesheet" type="text/css" />
    <link href="/sitecss/gllbetter/MacroCss.css?v=63688688113" rel="stylesheet" type="text/css" />
    <link href="/sitecss/gllbetter/TextCss.css?v=63688688113" rel="stylesheet" type="text/css" />
    <link href="/sitecss/gllbetter/CSCCss.css?v=63688688113" rel="stylesheet" type="text/css" />
    <script src="/sitescripts/jquery.min.js?v=63688688113" type="text/javascript"></script>
    <script src="/sitescripts/jquery-migrate.min.js?v=63688688113" type="text/javascript"></script>
    <script src="/sitescripts/MicrosoftAjax.js?v=63688688113" type="text/javascript"></script>
    <script src="/sitescripts/MicrosoftMvcAjax.js?v=63688688113" type="text/javascript"></script>
    <script src="/sitescripts/BookingsAJAX.js?v=63688688113" type="text/javascript"></script>
    <script src="/sitescripts/jtip.js?v=63688688113" type="text/javascript"></script>
    <script src="/sitescripts/thickbox.js?v=63688688113" type="text/javascript"></script>
</head>

<body class="timetableBody">
    <div id="resultUpdate">
        <div id="dayLinks">
            <a href="#" onclick='parent.tb_remove()' class='MembersButtonLarge'
                style="float: right; margin-top: -3px;">Close / Change Selection</a>
            <b class='TextMembers'>Showing results up to 04 April 2019 19:36.</b>
        </div>
        <div id="resultContainer">

            <div>
                <br style='clear:both;' />
                <div><br /><br />
                    <div class='activityHeader'>Creche Over 2's - Two hours </div><br />
                    <div class='sportsHallSlotWrapper'>
                        <div style='clear:both; padding: 5px 0;'><b class='TextMembers'>Friday 22 March 2019</b></div>
                        <div class='sporthallSlot'>Greenwich Centre<br />09:30 <a
                                href='GetSportsHallPrice?slotId=5545413' rel='GetSportsHallPrice?slotId=5545413'
                                id='price5545413' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br />Fully Booked</div>
                        <div class='sporthallSlot'>Greenwich Centre<br />11:30 <a
                                href='GetSportsHallPrice?slotId=5528959' rel='GetSportsHallPrice?slotId=5528959'
                                id='price5528959' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br />Fully Booked</div>
                        <div style='clear:both; padding: 5px 0;'><b class='TextMembers'>Monday 25 March 2019</b></div>
                        <div class='sporthallSlot'>Greenwich Centre<br />09:30 <a
                                href='GetSportsHallPrice?slotId=5629048' rel='GetSportsHallPrice?slotId=5629048'
                                id='price5629048' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br />Fully Booked</div>
                        <div class='sporthallSlot'>Greenwich Centre<br />11:30 <a
                                href='GetSportsHallPrice?slotId=5615760' rel='GetSportsHallPrice?slotId=5615760'
                                id='price5615760' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br />Fully Booked</div>
                        <div style='clear:both; padding: 5px 0;'><b class='TextMembers'>Tuesday 26 March 2019</b></div>
                        <div class='sporthallSlot'>Greenwich Centre<br />09:30 <a
                                href='GetSportsHallPrice?slotId=5662163' rel='GetSportsHallPrice?slotId=5662163'
                                id='price5662163' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br />Fully Booked</div>
                        <div class='sporthallSlot'>Greenwich Centre<br />11:30 <a
                                href='GetSportsHallPrice?slotId=5641294' rel='GetSportsHallPrice?slotId=5641294'
                                id='price5641294' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br /><a id='slot5641294'
                                class='sporthallSlotAddLink' href='#'
                                onclick='addSportsHallBooking(5641294); return false;' class='addLink'>2 Available</a>
                        </div>
                        <div style='clear:both; padding: 5px 0;'><b class='TextMembers'>Wednesday 27 March 2019</b>
                        </div>
                        <div class='sporthallSlot'>Greenwich Centre<br />09:30 <a
                                href='GetSportsHallPrice?slotId=5683061' rel='GetSportsHallPrice?slotId=5683061'
                                id='price5683061' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br />Fully Booked</div>
                        <div class='sporthallSlot'>Greenwich Centre<br />11:30 <a
                                href='GetSportsHallPrice?slotId=5696464' rel='GetSportsHallPrice?slotId=5696464'
                                id='price5696464' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br />Fully Booked</div>
                        <div style='clear:both; padding: 5px 0;'><b class='TextMembers'>Thursday 28 March 2019</b></div>
                        <div class='sporthallSlot'>Greenwich Centre<br />09:30 <a
                                href='GetSportsHallPrice?slotId=5733309' rel='GetSportsHallPrice?slotId=5733309'
                                id='price5733309' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br />Fully Booked</div>
                        <div class='sporthallSlot'>Greenwich Centre<br />11:30 <a
                                href='GetSportsHallPrice?slotId=5716521' rel='GetSportsHallPrice?slotId=5716521'
                                id='price5716521' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br />Fully Booked</div>
                        <div style='clear:both; padding: 5px 0;'><b class='TextMembers'>Friday 29 March 2019</b></div>
                        <div class='sporthallSlot'>Greenwich Centre<br />09:30 <a
                                href='GetSportsHallPrice?slotId=5755999' rel='GetSportsHallPrice?slotId=5755999'
                                id='price5755999' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br />Fully Booked</div>
                        <div class='sporthallSlot'>Greenwich Centre<br />11:30 <a
                                href='GetSportsHallPrice?slotId=5739314' rel='GetSportsHallPrice?slotId=5739314'
                                id='price5739314' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br />Fully Booked</div>
                        <div style='clear:both; padding: 5px 0;'><b class='TextMembers'>Monday 01 April 2019</b></div>
                        <div class='sporthallSlot'>Greenwich Centre<br />09:30 <a
                                href='GetSportsHallPrice?slotId=5820288' rel='GetSportsHallPrice?slotId=5820288'
                                id='price5820288' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br />Fully Booked</div>
                        <div class='sporthallSlot'>Greenwich Centre<br />11:30 <a
                                href='GetSportsHallPrice?slotId=5829610' rel='GetSportsHallPrice?slotId=5829610'
                                id='price5829610' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br />Fully Booked</div>
                        <div style='clear:both; padding: 5px 0;'><b class='TextMembers'>Tuesday 02 April 2019</b></div>
                        <div class='sporthallSlot'>Greenwich Centre<br />09:30 <a
                                href='GetSportsHallPrice?slotId=5870375' rel='GetSportsHallPrice?slotId=5870375'
                                id='price5870375' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br /><a id='slot5870375'
                                class='sporthallSlotAddLink' href='#'
                                onclick='addSportsHallBooking(5870375); return false;' class='addLink'>1 Available</a>
                        </div>
                        <div class='sporthallSlot'>Greenwich Centre<br />11:30 <a
                                href='GetSportsHallPrice?slotId=5857026' rel='GetSportsHallPrice?slotId=5857026'
                                id='price5857026' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br /><a id='slot5857026'
                                class='sporthallSlotAddLink' href='#'
                                onclick='addSportsHallBooking(5857026); return false;' class='addLink'>4 Available</a>
                        </div>
                        <div style='clear:both; padding: 5px 0;'><b class='TextMembers'>Wednesday 03 April 2019</b>
                        </div>
                        <div class='sporthallSlot'>Greenwich Centre<br />09:30 <a
                                href='GetSportsHallPrice?slotId=5892202' rel='GetSportsHallPrice?slotId=5892202'
                                id='price5892202' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br />Fully Booked</div>
                        <div class='sporthallSlot'>Greenwich Centre<br />11:30 <a
                                href='GetSportsHallPrice?slotId=5905462' rel='GetSportsHallPrice?slotId=5905462'
                                id='price5905462' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br />Fully Booked</div>
                        <div style='clear:both; padding: 5px 0;'><b class='TextMembers'>Thursday 04 April 2019</b></div>
                        <div class='sporthallSlot'>Greenwich Centre<br />09:30 <a
                                href='GetSportsHallPrice?slotId=5916755' rel='GetSportsHallPrice?slotId=5916755'
                                id='price5916755' class='jTip100'><img src='/media/siteimages/moreInfoBlue.gif'
                                    alt='Click for More info' style='border:0px;' /></a><br /><a id='slot5916755'
                                class='sporthallSlotAddLink' href='#'
                                onclick='addSportsHallBooking(5916755); return false;' class='addLink'>3 Available</a>
                        </div>
                        <div class='sporthallSlot'>Greenwich Centre<br />11:30 
                            <a href='GetSportsHallPrice?slotId=5926337' rel='GetSportsHallPrice?slotId=5926337' id='price5926337'
                                class='jTip100'>
                                <img src='/media/siteimages/moreInfoBlue.gif' alt='Click for More info'
                                    style='border:0px;' />
                            </a>
                            <br />
                            <a id='slot5926337' class='sporthallSlotAddLink' href='#' onclick='addSportsHallBooking(5926337); return false;'
                                class='addLink'>2 Available
                            </a>
                        </div>
                    </div>
                </div>
            </div>

        </div>
    </div>

    <script type="text/javascript">
        loadPriceHover();
    </script>
</body>

</html>`
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node) string
	f = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				// find <a id='slot5926337' class='sporthallSlotAddLink' href='#' onclick='addSportsHallBooking(5926337); return false;' class='addLink'>
				if a.Key == "id" {
					if strings.HasPrefix(a.Val, "slot") {
						fmt.Println("Found booking slot", a.Val)
						return strings.TrimPrefix(a.Val, "slot")
					}
					break
				}
				// looking for href = "/east_greenwich/Basket/AllocateBookingCredit?reservationId=49600430"
				// if a.Key == "href" {
				// 	if strings.Contains(a.Val, "AllocateBookingCredit") {
				// 		fmt.Println("Found Voucher link to", a.Val)
				// 		return a.Val
				// 	}
				// 	break
				// }
			}
		}
		for c := n.LastChild; c != nil; c = c.PrevSibling {
			if r := f(c); r != "" {
				return r
			}
		}
		return ""
	}
	fmt.Println(f(doc))
}
