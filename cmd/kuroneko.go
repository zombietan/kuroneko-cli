package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/codegangsta/cli"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func makeSpace(count int) string {
	s := "　"
	return strings.Repeat(s, count)
}

var Kuroneko = func(c *cli.Context) {
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
