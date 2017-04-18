package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func makeSpace(count int) string {
	s := "　"
	return strings.Repeat(s, count)
}

var TrackNumber = func(c *cli.Context) {
	if c.NArg() < 1 {
		fmt.Println("伝票番号を入力してください")
		return
	}
	slipNumber := c.Args()[0]
	values := url.Values{}
	values.Add("number00", "1")
	values.Add("number01", slipNumber)

	url := "http://toi.kuronekoyamato.co.jp/cgi-bin/tneko"
	resp, err := http.PostForm(url, values)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	utfBody := transform.NewReader(bufio.NewReader(resp.Body), japanese.ShiftJIS.NewDecoder())

	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		fmt.Println(err)
		return
	}

	doc.Find(".saisin td").Each(func(_ int, args *goquery.Selection) {
		if args.HasClass("bold") || args.HasClass("font14") {
			text := args.Text()
			fmt.Printf(" %s\n", text)
		}
	})

	fmt.Print("\n")

	doc.Find(".meisai tr").Each(func(i int, args *goquery.Selection) {
		if i != 0 {
			textArray := args.Find("td").Map(func(_ int, s *goquery.Selection) string {
				text := s.Text()
				return text
			})
			// fmt.Printf("%#v\n", textArray)
			detailArray := textArray[1:6]
			statusLength := utf8.RuneCountInString(detailArray[0])
			width := 15 - statusLength
			statusSpace := makeSpace(width)
			status := detailArray[0] + statusSpace
			storeLength := utf8.RuneCountInString(detailArray[3])
			width = 20 - storeLength
			storeSpace := makeSpace(width)
			store := detailArray[3] + storeSpace
			date, times, code := detailArray[1], detailArray[2], detailArray[4]
			fmt.Printf(" %s| %s | %s | %s| %s |\n", status, date, times, store, code)
		}
	})

	underLine := strings.Repeat("-", 99)
	fmt.Println(underLine)
}

var TrackSerialNumbers = func(c *cli.Context) {
	if c.NArg() < 1 {
		fmt.Println("伝票番号を入力してください")
		return
	}
	rowNumber := c.Args()[0]

	slipNumber := hyphenRemove(rowNumber)

	if !isInt(slipNumber) {
		fmt.Println("不正な数値です")
		return
	}

	if !is12Digits(slipNumber) {
		fmt.Println("12桁の伝票番号を入力してください")
		return
	}

	ch := checkdigit(slipNumber[:11])
	values := url.Values{}
	values.Add("number00", "1")
	var i uint
	for i = 0; i < c.Uint("serial"); i++ {
		s := fmt.Sprintf("number%02d", i+1)
		values.Add(s, <-ch)
	}

	url := "http://toi.kuronekoyamato.co.jp/cgi-bin/tneko"
	resp, err := http.PostForm(url, values)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	utfBody := transform.NewReader(bufio.NewReader(resp.Body), japanese.ShiftJIS.NewDecoder())

	docs, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		fmt.Println(err)
		return
	}

	d := color.New(color.FgYellow, color.Bold)
	docs.Find("center").Each(func(_ int, doc *goquery.Selection) {

		doc.Find(".saisin td").Each(func(_ int, args *goquery.Selection) {
			if args.HasClass("number") {
				subject := args.Text()
				// fmt.Printf(" %s\n", subject)
				d.Printf(" %s\n", subject)
			}
			if args.HasClass("bold") || args.HasClass("font14") {
				text := args.Text()
				fmt.Printf(" %s\n", text)
			}
		})

		fmt.Print("\n")

		doc.Find(".meisai tr").Each(func(i int, args *goquery.Selection) {
			if i != 0 {
				textArray := args.Find("td").Map(func(_ int, s *goquery.Selection) string {
					text := s.Text()
					return text
				})
				// fmt.Printf("%#v\n", textArray)
				detailArray := textArray[1:6]
				statusLength := utf8.RuneCountInString(detailArray[0])
				width := 15 - statusLength
				statusSpace := makeSpace(width)
				status := detailArray[0] + statusSpace
				storeLength := utf8.RuneCountInString(detailArray[3])
				width = 20 - storeLength
				storeSpace := makeSpace(width)
				store := detailArray[3] + storeSpace
				date, times, code := detailArray[1], detailArray[2], detailArray[4]
				fmt.Printf(" %s| %s | %s | %s| %s |\n", status, date, times, store, code)
			}
		})

		doc.Find("hr").Each(func(_ int, args *goquery.Selection) {
			if args.HasClass("middle") {
				underLine := strings.Repeat("-", 99)
				fmt.Println(underLine)
			}
		})
	})

}

func checkdigit(n string) <-chan string {
	ch := make(chan string)
	var conf int64 = 7
	go func() {
		num, _ := strconv.ParseInt(n, 10, 64)
		for {
			digitString := num % conf
			digit := strconv.FormatInt(digitString, 10)
			slipNum := strconv.FormatInt(num, 10) + digit
			ch <- slipNum
			num++
		}
	}()
	return ch
}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func is12Digits(s string) bool {
	if len(s) == 12 {
		return true
	}
	return false
}

func hyphenRemove(s string) string {
	if strings.Contains(s, "-") {
		removed := strings.Replace(s, "-", "", -1)
		return removed
	}

	return s
}
