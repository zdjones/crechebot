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
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

// Shorthand http methods
const (
	GET  = http.MethodGet
	POST = http.MethodPost
)

type step struct {
	url       *url.URL
	rawurl    string
	method    string
	retries   int
	query     map[string]string
	values    url.Values
	rawvalues map[string]string
	before    func(*http.Client, map[string]string)
	after     func(*http.Client, *http.Response, map[string]string)
	next      string
}

var steps = map[string]*step{}
var share = map[string]string{}

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("Not enough arguments (eg email and pass): ", os.Args)
	}

	buildSteps()
	// All users of cookiejar should import "golang.org/x/net/publicsuffix"
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar: jar,
	}

	// For each step
	//	prepare
	//	before
	//	build request
	//	do request
	//		retries every x seconds
	//	after
	// Next
	current := steps["loginCookies"]
	for current != nil {
		current.prepare()
		if current.before != nil {
			current.before(client, share)
		}

		log.Printf("Making %s request to %s with data %v\n", current.method, current.url.String(), current.values)
		var res *http.Response
		var err error
		switch current.method {
		case GET:
			res, err = client.Get(current.url.String())
		case POST:
			res, err = client.PostForm(current.url.String(), current.values)
		}
		if err != nil {
			if current.retries < 5 {
				current.retries++
				log.Printf("Trying again (attempt %d) Error: %s\n", current.retries, err)
				time.Sleep(5 * time.Second)
				continue
			} else {
				log.Fatalf("Too many errors calling %s, giving up: %s\n", current.url.String(), err)
			}
		}

		log.Printf("Completed %s request to %s\n", current.method, current.url.String())
		if current.after != nil {
			current.after(client, res, share)
		}
		current = steps[current.next]
	}
	fmt.Println(share)

}

func (s *step) prepare() {
	var err error
	s.url, err = url.Parse(s.rawurl)
	if err != nil {
		log.Fatalf("Bad rawurl (%s): %s\n", s.rawurl, err)
	}

	s.values = url.Values{}
	for k, v := range s.rawvalues {
		s.values.Set(k, v)
	}

}

func buildSteps() {

	// Login "/east_greenwich/account/login"
	// GET for cookies
	steps["loginCookies"] = &step{
		rawurl: "https://better.legendonlineservices.co.uk/east_greenwich/account/login",
		method: GET,
		next:   "login",
	}

	// Need cookie __RequestVerificationToken added as param for login POST
	// POST login.Email, login.Password, login.RedirectURL, __RequestVerificationToken
	steps["login"] = &step{
		rawurl: "https://better.legendonlineservices.co.uk/east_greenwich/account/login",
		method: POST,
		rawvalues: map[string]string{
			"login.Email":                os.Args[1],
			"login.Password":             os.Args[2],
			"login.RedirectURL":          "",
			"__RequestVerificationToken": "",
		},
		next: "behaviours",
	}
	steps["login"].before = func(c *http.Client, data map[string]string) {
		for _, cookie := range c.Jar.Cookies(steps["login"].url) {
			if cookie.Name == "__RequestVerificationToken" {
				steps["login"].rawvalues[cookie.Name] = cookie.Value
			}
		}
	}

	// Behaviours "/east_greenwich/bookingscentre/behaviours"
	// POST ("club", "343"), ("X-Requested-With", "XMLHttpRequest")
	steps["behaviours"] = &step{
		rawurl: "https://better.legendonlineservices.co.uk/east_greenwich/bookingscentre/behaviours",
		method: POST,
		rawvalues: map[string]string{
			"club":             "343",
			"X-Requested-With": "XMLHttpRequest",
		},
		next: "activities",
	}

	// Activities "/east_greenwich/bookingscentre/activities"
	// POST ("behaviours", "2366"), ("bookingType", "1"), ("X-Requested-With", "XMLHttpRequest")
	steps["activities"] = &step{
		rawurl: "https://better.legendonlineservices.co.uk/east_greenwich/bookingscentre/activities",
		method: POST,
		rawvalues: map[string]string{
			"behaviours":       "2366",
			"bookingType":      "1",
			"X-Requested-With": "XMLHttpRequest",
		},
		next: "activitySelect",
	}

	// Activity Select	"/east_greenwich/bookingscentre/activitySelect"
	// POST ("activity", "882"), ("X-Requested-With", "XMLHttpRequest")
	// creche for Under 2 2 hours is activity 752
	steps["activitySelect"] = &step{
		rawurl: "https://better.legendonlineservices.co.uk/east_greenwich/bookingscentre/activitySelect",
		method: POST,
		rawvalues: map[string]string{
			"activity":         "882",
			"X-Requested-With": "XMLHttpRequest",
		},
		next: "timetableSubmit",
	}

	// Timetable "/east_greenwich/bookingscentre/TimeTable"
	// POST ("X-Requested-With", "XMLHttpRequest")
	steps["timetableSubmit"] = &step{
		rawurl: "https://better.legendonlineservices.co.uk/east_greenwich/bookingscentre/TimeTable",
		method: POST,
		rawvalues: map[string]string{
			"X-Requested-With": "XMLHttpRequest",
		},
		next: "timetableRead",
	}

	// "/east_greenwich/BookingsCentre/Timetable?KeepThis=true&"
	// 	GET ("KeepThis", "true")
	// 	Search response Body for slotID (curently just grabs last one on page = 1130)
	steps["timetableRead"] = &step{
		rawurl: "https://better.legendonlineservices.co.uk/east_greenwich/BookingsCentre/Timetable",
		query: map[string]string{
			"KeepThis": "true",
		},
		method: GET,
		after: func(c *http.Client, r *http.Response, share map[string]string) {
			// Extract slot ID from last slot (1130, two weeks from today)
			doc, err := html.Parse(r.Body)
			// html.Parse consumes the entire Body (to EOF), so I can close Body here.
			r.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			// Depth-first search, from bottom to top of doc
			var getSlot func(*html.Node) string
			getSlot = func(n *html.Node) string {
				if n.Type == html.ElementNode && n.Data == "a" {
					for _, a := range n.Attr {
						fmt.Println("Checking for slot id in: ", a)
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
			share["slotID"] = getSlot(doc)
		},
		// next: "addBooking",
	}

	// Add Booking "/east_greenwich/BookingsCentre/AddSportsHallBooking?ajax=0.6046602879796196&selectedCourts=&slotId=235234"
	// GET ("slotId", slotID), ("selectedCourts", selectedCourts), ("ajax", ajax)
	steps["addBooking"] = &step{
		rawurl: "https://better.legendonlineservices.co.uk/east_greenwich/BookingsCentre/AddSportsHallBooking",
		query: map[string]string{
			"selectedCourts": "",
			"ajax":           "",
			"slotID":         "",
		},
		method: GET,
		before: func(c *http.Client, share map[string]string) {
			steps["addBooking"].query["ajax"] = fmt.Sprintf("%.16f", rand.Float64()) // js NUMBER ~= float64
			steps["addBooking"].query["slotID"] = share["slotID"]
		},
		//  Read JSON from response Body to get succes (and slot unix times?)
		//	??Check Json message to verify correct slot before continuing?
		after: func(c *http.Client, r *http.Response, share map[string]string) {
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Recieved json:", string(body))
			// TODO: get more data from json response
			//		1. unix nanoseconds describing the slot start/end
			//			can be used to check correctness before finishing
			data := struct{ Success bool }{} // anonymous struct for one-off use
			err = json.Unmarshal(body, &data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Booking added to basket =", data.Success)
		},
		next: "basket",
	}

	// Message???		"/east_greenwich/BookingsCentre/Message?Success=true&AllowRetry=false&Message=Booking+added+to+basket&StartTime=Tue,+26+Mar+2019+11:30:00+GMT&EndTime=Tue,+26+Mar+2019+13:30:00+GMT&FacilityName=Greenwich+Centre&ActivityName=Creche+Over+2"
	// GET -not required? params from JSON above
	// NOT IMPLEMENTED

	// Basket "/east_greenwich/Basket/Index"
	// GET
	// Search for reservation ID (curently just grabs last one on page, I think this is the last added booking)
	steps["basket"] = &step{
		rawurl: "https://better.legendonlineservices.co.uk/east_greenwich/Basket/Index",
		method: GET,
		after: func(c *http.Client, r *http.Response, share map[string]string) {
			// Find Voucher button/id/href
			doc, err := html.Parse(r.Body)
			// html.Parse consumes the entire Body (to EOF), so I can close Body here.
			r.Body.Close()
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
			u, err := url.Parse(uVoucherPath)
			if err != nil {
				log.Printf("basket step, error parsing url %s: %s\n", uVoucherPath, err)
			}
			q, err := url.ParseQuery(u.RawQuery)
			if err != nil {
				log.Printf("basket step, error parsing query string %s: %s\n", u.RawQuery, err)
			}
			share["reservationId"] = q.Get("reservationId")
		},
		next: "applyVoucher",
	}

	// Apply Voucher "/east_greenwich/Basket/AllocateBookingCredit?reservationId=49600430"
	// GET ("reservationID", "XXXXXXXXX")
	steps["applyVoucher"] = &step{
		rawurl: "https://better.legendonlineservices.co.uk/east_greenwich/Basket/AllocateBookingCredit",
		method: GET,
		rawvalues: map[string]string{
			"reservationId": "",
		},
		before: func(c *http.Client, share map[string]string) {
			steps["applyVoucher"].query["reservationId"] = share["reservationId"]
		},
		next: "complete",
	}

	// Pay "/east_greenwich/Basket/Pay"
	// GET
	steps["complete"] = &step{
		rawurl: "https://better.legendonlineservices.co.uk/east_greenwich/Basket/Pay",
		method: GET,
	}

	// Check response for successful completion?????
	// go again if failed (x5?)

}
