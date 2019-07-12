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
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

func old() {
	loginURL, err := url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/account/login")
	if err != nil {
		log.Fatal(err)
	}

	// All users of cookiejar should import "golang.org/x/net/publicsuffix"
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar: jar,
	}

	// betterURL, err := url.Parse("https://better.legendonlineservices.co.uk/")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// TODO: add user agent?
	// TODO: add request timeout?

	// get cookies from login page
	_, err = client.Get(loginURL.String())
	if err != nil {
		log.Fatal(err)
	}

	// set form values for login POST request
	v := url.Values{}
	v.Set("login.Email", "zachj1@gmail.com")
	v.Set("login.Password", "london11")
	v.Set("login.RedirectURL", "")
	// this form value is a cookie received on initial request
	for _, cookie := range jar.Cookies(loginURL) {
		// fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
		if cookie.Name == "__RequestVerificationToken" {
			v.Set("__RequestVerificationToken", cookie.Value)
		}
	}

	resp, err := client.PostForm(loginURL.String(), v)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Login response status code: ", resp.StatusCode)

	// Select correct club (Greenwich Centre)
	v2 := url.Values{}
	v2.Set("club", "343")
	v2.Set("X-Requested-With", "XMLHttpRequest")

	u2, err := url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/bookingscentre/behaviours")
	if err != nil {
		log.Fatal(err)
	}

	resp, err = client.PostForm(u2.String(), v2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Posted club to", u2.String())

	// Select correct activity category (Creche)
	v3 := url.Values{}
	v3.Set("behaviours", "2366")
	v3.Set("bookingType", "1")
	v3.Set("X-Requested-With", "XMLHttpRequest")

	u3, err := url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/bookingscentre/activities")
	if err != nil {
		log.Fatal(err)
	}

	resp, err = client.PostForm(u3.String(), v3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Posted activity category to", u3.String())

	// Select correct activity (CreCreche Over 2's - Two hours)
	v4 := url.Values{}
	v4.Set("activity", "882")
	v4.Set("X-Requested-With", "XMLHttpRequest")

	u4, err := url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/bookingscentre/activitySelect")
	if err != nil {
		log.Fatal(err)
	}

	resp, err = client.PostForm(u4.String(), v4)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Posted activity to", u4.String())

	// Submit Timetable button
	v5 := url.Values{}
	v5.Set("X-Requested-With", "XMLHttpRequest")

	u5, err := url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/bookingscentre/TimeTable")
	if err != nil {
		log.Fatal(err)
	}

	resp, err = client.PostForm(u5.String(), v5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Posted Timetable button to", u5.String())

	// Get Timetable subdocument html - sometimes after previous (u5)
	// this is the html containing the creche slots and their IDs
	u6, err := url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/BookingsCentre/Timetable?KeepThis=true&")
	if err != nil {
		log.Fatal(err)
	}

	resp, err = client.Get(u6.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Got Timetable subdocument from", u6.String())

	// Extract slot ID from last slot (1130, two weeks from today)
	doc, err := html.Parse(resp.Body)
	// html.Parse consumes the entire Body (to EOF), so I can close Body here.
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

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
	slotID := getSlot(doc)

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

	fmt.Println("Requesting to add booking to", bookingURL.String())
	resp, err = client.Get(bookingURL.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Recieved json:", string(body))
	// TODO: get more data from json response
	//		1. unix nanoseconds describing the slot start/end
	//			can be used to check correctness before finishing
	data := struct{ Success bool }{} // anonymous struct for one-of use
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Booking added to basket =", data.Success)

	// click to select also requests:
	// https://better.legendonlineservices.co.uk/east_greenwich/BookingsCentre/Message?Success=true&AllowRetry=false&Message=Booking+added+to+basket&StartTime=Tue,+26+Mar+2019+11:30:00+GMT&EndTime=Tue,+26+Mar+2019+13:30:00+GMT&FacilityName=Greenwich+Centre&ActivityName=Creche+Over+2
	// which responds with little box to choose another or go to checkout

	// Get basket

	uBasket, err := url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/Basket/Index")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Requesting basket", uBasket.String())
	resp, err = client.Get(uBasket.String())
	if err != nil {
		log.Fatal(err)
	}

	// Find Voucher button/id/href
	doc, err = html.Parse(resp.Body)
	// html.Parse consumes the entire Body (to EOF), so I can close Body here.
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
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

	fmt.Println("Applying Voucher", uVoucher.String())
	resp, err = client.Get(uVoucher.String())
	if err != nil {
		log.Fatal(err)
	}

	// Complete Purchase/Check Out
	// https://better.legendonlineservices.co.uk/east_greenwich/Basket/Pay
	uPay, err := url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/Basket/Pay")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Completing Transaction (Pay zero balance and confirm booking)", uPay.String())
	resp, err = client.Get(uPay.String())
	if err != nil {
		log.Fatal(err)
	}
}
