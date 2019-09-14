package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"golang.org/x/net/publicsuffix"
)

type client struct {
	*http.Client
}

var maxRetries = 10

const (
	OVER2_1HOUR   = "752" // over 2s, 1 hour = 752
	OVER2_2HOURS  = "882" // over 2s, 2 hours = 882
	UNDER2_1HOUR  = "751" // under 2s, 1 hour = 751
	UNDER2_2HOURS = "883" // under 2s, 2 hours = 883
)

func main() {
	// if len(os.Args) < 3 {
	// 	log.Fatalln("Not enough arguments (eg email and pass): ", os.Args)
	// }
	var under2s bool
	var early bool
	flag.BoolVar(&under2s, "under2", false, "Book for Under 2s (default is Over 2s)")
	flag.BoolVar(&early, "early", false, "Book for 930 (default is 1130)")
	// long := flag.String("l", flag.Args[1], "long url")
	// short := flag.String("s", flag.Args[1], "short url")
	// addr := flag.String("api", flag.Args[2], "API endpoint")
	flag.Parse()

	if len(flag.Args()) < 2 {
		log.Fatalln("Not enough arguments (eg email and pass): ", flag.Args())
	}

	log.Println("BEGIN creche bot at", time.Now().Format(time.RFC822))

	user := flag.Arg(0)
	pass := flag.Arg(1)

	// over 2s, 1 hour = 752
	// over 2s, 2 hours = 882
	// under 2s, 1 hour = 751
	// under 2s, 2 hours = 883
	var crecheType string
	if under2s {
		crecheType = UNDER2_2HOURS
	} else {
		crecheType = OVER2_2HOURS
	}

	c := newClient()
	c.login(user, pass)
	c.selectCentre()
	c.selectActivityCreche()
	c.selectCrecheType(crecheType)
	c.addBooking(early)
	c.applyVoucher()
	c.complete()

	log.Println("END creche bot without failures")
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
		log.Printf("Failed attempt at %s. StatusCode %d, error: %s\n", u.Path, res.StatusCode, err)
		time.Sleep(5 * time.Second)
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
		log.Printf("Failed attempt at %s. StatusCode %d, error: %s", u.Path, res.StatusCode, err)
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
		log.Printf("Failed attempt at %s. StatusCode %d, error: %s", u.Path, res.StatusCode, err)
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
		log.Printf("Failed attempt at %s. StatusCode %d, error: %s", u.Path, res.StatusCode, err)
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
		log.Printf("Failed attempt at %s. StatusCode %d, error: %s", u.Path, res.StatusCode, err)
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
		log.Printf("Failed attempt at %s. StatusCode %d, error: %s", u.Path, res.StatusCode, err)
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
		log.Printf("Failed attempt at %s. StatusCode %d, error: %s", u.Path, res.StatusCode, err)
		if attempt >= maxRetries {
			log.Fatalf("Too many failed attempts, giving up (%s): %s", u.String(), err)
		}
	}
	defer res.Body.Close()
	// Load the HTML document
	return goquery.NewDocumentFromReader(res.Body)
}

func getSlotID(doc *goquery.Document, early bool) string {
	// Extract slot ID from last slot (1130, two weeks from today)
	// Find the slots
	targetDate := time.Now().UTC().AddDate(0, 0, 14)
	// For targeting an arbirtrary date
	// TODO: Paramaterize target date and time
	// targetDate, _ := time.Parse("Monday 2 January 2006", "Monday 15 July 2019")
	fmt.Println("Looking for slots on", targetDate.Format("Monday 2 January 2006"))

	var targetTime string
	switch early {
	case true:
		targetTime = "9:30"
	case false:
		targetTime = "11:30"
	}
	slotDate := time.Time{}
	slotID := ""
	doc.Find(".sportsHallSlotWrapper").Children().EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each slot div, get the time and the slotId
		// Mon Jan 2 15:04:05 -0700 MST 2006
		// Monday 15 July 2019
		if s.Is(".sporthallSlot") {
			if slotDate.YearDay() == targetDate.YearDay() && slotDate.Year() == targetDate.Year() {
				text := strings.TrimSpace(s.Text())
				if strings.Contains(text, targetTime) {
					link := s.Find(".sporthallSlotAddLink")
					if id, ok := link.Attr("id"); ok {
						slotID = strings.TrimPrefix(id, "slot")
						log.Println("Found slot id for date:", slotDate.Format("Mon 2 Jan 2006"), text, slotID)
						return false
					}
					fmt.Println("Unable to find slot id for date:", slotDate.Format("Mon 2 Jan 2006"), s.Text())
					return true
				}
				fmt.Println("Wrong time, skip this slot:", slotDate.Format("Mon 2 Jan 2006"), s.Text())
				return true
			}
			fmt.Println("Wrong date, skip this slot:", slotDate.Format("Mon 2 Jan 2006"), s.Text())
			return true
		} else {
			text := strings.TrimSpace(s.Text())
			var err error
			slotDate, err = time.Parse("Monday 2 January 2006", text)
			if err != nil {
				log.Println("Can't parse date string from timetable: ", text)
				slotDate = time.Time{}
			} else {
				fmt.Println("Next slots are for date:", slotDate.Format("Mon 2 Jan 2006"))
			}
		}
		return true
	})
	return slotID
}

func (c *client) addBooking(early bool) {
	doc, err := c.getTimetableHTML()
	if err != nil {
		log.Fatalf("Bad timetable HTML: %s", err)
	}

	slotID := getSlotID(doc, early)
	if len(slotID) < 1 {
		log.Fatalf("Failed to find slotID")
	}
	// for selecting a booking, see BookingAjax.js
	// Click to select booking calls addSportsHallBooking(bookingID): line 150
	// addSportsHallBooking() sets selectedCourts to empty string if selectedCourts
	// is not provided as second parameter, as in this case
	// selectedCourts := ""
	// slotID := "5641294"                          //should just be able to use final booking
	// addSportsHallBooking() then calls addBooking(id, url)
	// addBooking(id, "AddSportsHallBooking?slotId=" + id + "&selectedCourts=" + selectedCourts);
	// AddSportsHallBooking?slotId=5857026&selectedCourts=&ajax=0.2919207756868989
	// build the url
	// bookingURL := url.URL{}
	// bookingURL.Scheme = "https"
	// bookingURL.Host = "better.legendonlineservices.co.uk"
	// bookingURL.Path = "/east_greenwich/BookingsCentre/AddSportsHallBooking"
	// addBooking(id, url) starts at line 109 of BookinAjax.js
	// need to find how to locaate correct link to use, because I dont
	// notice an obvious generation scheme for booking IDs.
	//                 https://better.legendonlineservices.co.uk/east_greenwich/BookingsCentre/AddSportsHallBooking?ajax=0.6046602879796196&selectedCourts=&slotId=235234
	// get request to: https://better.legendonlineservices.co.uk/east_greenwich/BookingsCentre/AddSportsHallBooking?slotId=5641294&selectedCourts=&ajax=0.09610071300889433
	// params: slotId	5641294
	//   selectedCourts
	//   ajax	0.09610071300889433	---> from math.random, wonder what this does?
	u, err := url.Parse("https://better.legendonlineservices.co.uk/east_greenwich/BookingsCentre/AddSportsHallBooking")
	if err != nil {
		log.Fatalf("Failed to parse url (%s): %s", u.String(), err)
	}

	var res *http.Response
	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Println("Requesting to add booking to", u.String())
		ajax := fmt.Sprintf("%.16f", rand.Float64()) // js NUMBER ~= float64
		q := u.Query()
		q.Set("slotId", slotID)
		q.Set("selectedCourts", "")
		q.Set("ajax", ajax)
		u.RawQuery = q.Encode()
		res, err = c.Get(u.String())
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Failed to read addBooking response body: %s", err)
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
		if err == nil && res.StatusCode == 200 && data.Success {
			fmt.Println("Added booking:", u.String())
			break
		}
		log.Printf("Failed attempt at %s. StatusCode %d, error: %s", u.Path, res.StatusCode, err)
		if attempt >= maxRetries {
			log.Fatalf("Too many failed attempts, giving up (%s): %s", u.String(), err)
		}
		// sleep longer after each attempt
		time.Sleep(time.Duration(attempt) * 5 * time.Second)
	}
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
		log.Printf("Failed attempt at %s. StatusCode %d, error: %s", uBasket.Path, res.StatusCode, err)
		if attempt >= maxRetries {
			log.Fatalf("Too many failed attempts, giving up (%s): %s", uBasket.String(), err)
		}
	}

	// Find Voucher button/id/href
	voucherPath := ""
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Failed to parse html from basket: %s", err)
	}
	doc.Find(".basketItem").Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		text := s.Text()
		if strings.Contains(text, "Use Voucher") {
			var ok bool
			voucherPath, ok = s.Attr("href")
			if ok {
				fmt.Println("Found Voucher link:", voucherPath)
				return false
			}
			log.Println("Can't find href in Voucher link")
		}
		return true
	})

	if len(voucherPath) < 1 {
		log.Fatalln("Failed to find voucher path, abandonind booking")
	}

	// Apply Voucher
	// https://better.legendonlineservices.co.uk/east_greenwich/Basket/AllocateBookingCredit?reservationId=49600430
	uVoucher, err := url.Parse("https://better.legendonlineservices.co.uk" + voucherPath)
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
		log.Printf("Failed attempt at %s. StatusCode %d, error: %s", uVoucher.Path, res.StatusCode, err)
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
		log.Printf("Failed attempt at %s. StatusCode %d, error: %s", uPay.Path, res.StatusCode, err)
		if attempt >= maxRetries {
			log.Fatalf("Too many failed attempts, giving up (%s): %s", uPay.String(), err)
		}
	}
}
