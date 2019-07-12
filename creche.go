package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

type client struct {
	*http.Client
}

var maxRetries = 10

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("Not enough arguments (eg email and pass): ", os.Args)
	}
	user := os.Args[1]
	pass := os.Args[2]

	c := newClient()
	c.login(user, pass)
	c.selectCentre()
	c.selectActivityCreche()
	// over 2s, 2 hours = 882
	// under 2s, 2 hours = 752
	c.selectCrecheType("882")
	c.addBooking()
	c.applyVoucher()
	c.complete()
}

func newClient() *client {
	// All users of cookiejar should import "golang.org/x/net/publicsuffix"
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	c := &http.Client{
		Jar: jar,
	}

	return &client{c}
}

func (c *client) login(user, pass string) {
	u, err := url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/account/login")
	if err != nil {
		log.Fatalf("Failed to parse url (%s): %s", u.String(), err)
	}
	// get cookies from login page, no need for the resonse
	for attempt := 1; attempt <= maxRetries; attempt++ {
		res, err := c.Get(u.String())
		if err == nil && res.StatusCode == 200 {
			break
		}
		if attempt >= maxRetries {
			log.Fatalf("Too many failed attempts, giving up (%s): %s", u.String(), err)
		}
	}

	// set form values for login POST request
	v := url.Values{}
	v.Set("login.Email", user)
	v.Set("login.Password", pass)
	v.Set("login.RedirectURL", "")
	// this form value is a cookie received on initial request
	for _, cookie := range c.Jar.Cookies(u) {
		// fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
		if cookie.Name == "__RequestVerificationToken" {
			v.Set("__RequestVerificationToken", cookie.Value)
		}
	}

	// login
	for attempt := 1; attempt <= maxRetries; attempt++ {
		res, err := c.PostForm(u.String(), v)
		if err == nil && res.StatusCode == 200 {
			fmt.Println("Login resonse status code: ", res.StatusCode)
			break
		}
		if attempt >= maxRetries {
			log.Fatalf("Too many failed attempts, giving up (%s): %s", u.String(), err)
		}
	}
}

func (c *client) selectCentre() {
	// Select correct club (Greenwich Centre)
	v := url.Values{}
	v.Set("club", "343")
	v.Set("X-Requested-With", "XMLHttpRequest")

	u, err := url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/bookingscentre/behaviours")
	if err != nil {
		log.Fatalf("Failed to parse url (%s): %s", u.String(), err)
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		res, err := c.PostForm(u.String(), v)
		if err == nil && res.StatusCode == 200 {
			fmt.Println("Posted club to", u.String())
			break
		}
		if attempt >= maxRetries {
			log.Fatalf("Too many failed attempts, giving up (%s): %s", u.String(), err)
		}
	}

}

func (c *client) selectActivityCreche() {
	// Select correct activity category (Creche)
	v := url.Values{}
	v.Set("behaviours", "2366")
	v.Set("bookingType", "1")
	v.Set("X-Requested-With", "XMLHttpRequest")

	u, err := url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/bookingscentre/activities")
	if err != nil {
		log.Fatalf("Failed to parse url (%s): %s", u.String(), err)
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		res, err := c.PostForm(u.String(), v)
		if err == nil && res.StatusCode == 200 {
			fmt.Println("Posted activity category to", u.String())
			break
		}
		if attempt >= maxRetries {
			log.Fatalf("Too many failed attempts, giving up (%s): %s", u.String(), err)
		}
	}
}

func (c *client) selectCrecheType(activity string) {
	// Select correct activity (CreCreche Over 2's - Two hours)
	v := url.Values{}
	// v.Set("activity", "882") Over 2s two hours
	v.Set("activity", activity)
	v.Set("X-Requested-With", "XMLHttpRequest")

	u, err := url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/bookingscentre/activitySelect")
	if err != nil {
		log.Fatalf("Failed to parse url (%s): %s", u.String(), err)
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		res, err := c.PostForm(u.String(), v)
		if err == nil && res.StatusCode == 200 {
			fmt.Println("Posted activity to", u.String())
			break
		}
		if attempt >= maxRetries {
			log.Fatalf("Too many failed attempts, giving up (%s): %s", u.String(), err)
		}
	}
}

func (c *client) getTimetableHTML() (*goquery.Document, error) {
	// Submit Timetable button
	v := url.Values{}
	v.Set("X-Requested-With", "XMLHttpRequest")

	u, err := url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/bookingscentre/TimeTable")
	if err != nil {
		log.Fatalf("Failed to parse url (%s): %s", u.String(), err)
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		res, err := c.PostForm(u.String(), v)
		if err == nil && res.StatusCode == 200 {
			fmt.Println("Posted Timetable button to", u.String())
			break
		}
		if attempt >= maxRetries {
			log.Fatalf("Too many failed attempts, giving up (%s): %s", u.String(), err)
		}
	}

	// Get Timetable subdocument html - sometimes after previous (u5)
	// this is the html containing the creche slots and their IDs
	u, err = url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/BookingsCentre/Timetable?KeepThis=true&")
	if err != nil {
		log.Fatalf("Failed to parse url (%s): %s", u.String(), err)
	}

	var res *http.Response
	for attempt := 1; attempt <= maxRetries; attempt++ {
		res, err = c.Get(u.String())
		if err == nil && res.StatusCode == 200 {
			fmt.Println("Got Timetable subdocument from", u.String())
			break
		}
		if attempt >= maxRetries {
			log.Fatalf("Too many failed attempts, giving up (%s): %s", u.String(), err)
		}
	}
	defer res.Body.Close()
	// Load the HTML document
	return goquery.NewDocumentFromReader(res.Body)
}

func getSlotID(doc *goquery.Document) string {
	// Extract slot ID from last slot (1130, two weeks from today)
	// Find the slots
	// wrapper := doc.Find(".sportsHallSlotWrapper").Each(func(i int, s *goquery.Selection) {
	// 	if s.Is
	// wrapper.Find(".sporthallSlot")
	// 	// For each slot div, get the time and the slotId
	// 	datetime
	// 	band := s.Find("a").Text()
	// 	title := s.Find("i").Text()
	// 	fmt.Printf("Review %d: %s - %s\n", i, band, title)
	// })
	// Depth-first search, from bottom to top of doc
	var getSlot func(*html.Node) string
	getSlot = func(n *html.Node) string {
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
			}
		}
		for c := n.LastChild; c != nil; c = c.PrevSibling {
			// If the target is found, unravel the recursion
			if r := getSlot(c); r != "" {
				return r
			}
		}
		return ""
	}
	// return getSlot(doc), nil
	return ""
}

func (c *client) addBooking() {
	doc, err := c.getTimetableHTML()
	if err != nil {
		log.Fatalf("Bad timetable HTML: %s", err)
	}

	slotID := getSlotID(doc)
	if len(slotID) < 1 {
		log.Fatalf("Failed to find slotID")
	}
	// for selecting a booking, see BookingAjax.js
	// Click to select booking calls addSportsHallBooking(bookingID): line 150
	// addSportsHallBooking() sets selectedCourts to empty string if selectedCourts
	// is not provided as second parameter, as in this case
	selectedCourts := ""
	ajax := fmt.Sprintf("%.16f", rand.Float64()) // js NUMBER ~= float64
	// slotID := "5641294"                          //should just be able to use final booking
	// addSportsHallBooking() then calls addBooking(id, url)
	// addBooking(id, "AddSportsHallBooking?slotId=" + id + "&selectedCourts=" + selectedCourts);
	// AddSportsHallBooking?slotId=5857026&selectedCourts=&ajax=0.2919207756868989
	// build the url
	bookingURL := url.URL{}
	bookingURL.Scheme = "https"
	bookingURL.Host = "better.legendonlineservices.co.uk"
	bookingURL.Path = "/east_greenwich/BookingsCentre/AddSportsHallBooking"
	q := bookingURL.Query()
	q.Set("slotId", slotID)
	q.Set("selectedCourts", selectedCourts)
	q.Set("ajax", ajax)
	bookingURL.RawQuery = q.Encode()
	// addBooking(id, url) starts at line 109 of BookinAjax.js
	// need to find how to locaate correct link to use, because I dont
	// notice an obvious generation scheme for booking IDs.
	//                 https://better.legendonlineservices.co.uk/east_greenwich/BookingsCentre/AddSportsHallBooking?ajax=0.6046602879796196&selectedCourts=&slotId=235234
	// get request to: https://better.legendonlineservices.co.uk/east_greenwich/BookingsCentre/AddSportsHallBooking?slotId=5641294&selectedCourts=&ajax=0.09610071300889433
	// params: slotId	5641294
	//   selectedCourts
	//   ajax	0.09610071300889433	---> from math.random, wonder what this does?

	var res *http.Response
	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Println("Requesting to add booking to", bookingURL.String())
		res, err = c.Get(bookingURL.String())
		if err == nil && res.StatusCode == 200 {
			fmt.Println("Added booking:", bookingURL.String())
			break
		}
		if attempt >= maxRetries {
			log.Fatalf("Too many failed attempts, giving up (%s): %s", bookingURL.String(), err)
		}
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed to read addBooking response body: %s", err)
	}
	fmt.Println("Recieved json:", string(body))

	// TODO: get more data from json resonse
	//		1. unix nanoseconds describing the slot start/end
	//			can be used to check correctness before finishing
	data := struct{ Success bool }{false} // anonymous struct for one-off use
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("Failed to parse booking response json")
	}

	fmt.Println("Booking added to basket =", data.Success)
}

func (c *client) applyVoucher() {
	// Get basket

	uBasket, err := url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/Basket/Index")
	if err != nil {
		log.Fatal(err)
	}

	var res *http.Response
	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Println("Requesting basket", uBasket.String())
		res, err = c.Get(uBasket.String())
		if err == nil && res.StatusCode == 200 {
			fmt.Println("Got Basket from", uBasket.String())
			break
		}
		if attempt >= maxRetries {
			log.Fatalf("Too many failed attempts, giving up (%s): %s", uBasket.String(), err)
		}
	}

	// Find Voucher button/id/href
	defer res.Body.Close()
	doc, err := html.Parse(res.Body)
	if err != nil {
		log.Fatalf("Failed to parse html from basket: %s", err)
	}

	// Depth-first search, from bottom to top of doc
	var getVoucher func(*html.Node) string
	getVoucher = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					// looking for href = "/east_greenwich/Basket/AllocateBookingCredit?reservationId=49600430"
					if strings.Contains(a.Val, "AllocateBookingCredit") {
						fmt.Println("Found Voucher link to", a.Val)
						return a.Val
					}
					break
				}
			}
		}
		for c := n.LastChild; c != nil; c = c.PrevSibling {
			// If the target is found, unravel the recursion
			if r := getVoucher(c); r != "" {
				return r
			}
		}
		return ""
	}
	uVoucherPath := getVoucher(doc)

	// Apply Voucher
	// https://better.legendonlineservices.co.uk/east_greenwich/Basket/AllocateBookingCredit?reservationId=49600430
	uVoucher, err := url.Parse("https://better.legendonlineservices.co.uk" + uVoucherPath)
	if err != nil {
		log.Fatal(err)
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Println("Applying Voucher", uVoucher.String())
		res, err := c.Get(uVoucher.String())
		if err == nil && res.StatusCode == 200 {
			fmt.Println("Applied Voucher at", uVoucher.String())
			break
		}
		if attempt >= maxRetries {
			log.Fatalf("Too many failed attempts, giving up (%s): %s", uVoucher.String(), err)
		}
	}
}

func (c *client) complete() {
	// Complete Purchase/Check Out
	// https://better.legendonlineservices.co.uk/east_greenwich/Basket/Pay
	uPay, err := url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/Basket/Pay")
	if err != nil {
		log.Fatal(err)
	}

	var res *http.Response
	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Println("Completing Transaction (Pay zero balance and confirm booking)", uPay.String())
		res, err = c.Get(uPay.String())
		if err == nil && res.StatusCode == 200 {
			fmt.Println("Completed Transaction!!", uPay.String())
			break
		}
		if attempt >= maxRetries {
			log.Fatalf("Too many failed attempts, giving up (%s): %s", uPay.String(), err)
		}
	}
}
