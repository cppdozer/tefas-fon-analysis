package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	TelegramBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nleeper/goment"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type FundData struct {
	Tarih            string  `json:"TARIH"`
	FonKodu          string  `json:"FONKODU"`
	FonUnvan         string  `json:"FONUNVAN"`
	Fiyat            float64 `json:"FIYAT"`
	TedPaySayisi     float64 `json:"TEDPAYSAYISI"`
	KisiSayisi       float64 `json:"KISISAYISI"`
	PortfoyBuyukluk  float64 `json:"PORTFOYBUYUKLUK"`
	BorsaBultenFiyat string  `json:"BORSABULTENFIYAT"`
}

type Response struct {
	Draw            int        `json:"draw"`
	RecordsTotal    int        `json:"recordsTotal"`
	RecordsFiltered int        `json:"recordsFiltered"`
	Data            []FundData `json:"data"`
}

func fonPrices(name string) string {

	now := time.Now()

	// Calculate the duration of one week
	oneWeek := 6 * 24 * time.Hour

	// Subtract one week from the current time
	oneWeekAgo := now.Add(-oneWeek)

	today := time.Now().Format("02/01/2006")

	var API_URL string = "https://www.tefas.gov.tr/api/DB/BindHistoryInfo"

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,                                    // Use TLS 1.2 or higher
			CipherSuites:       []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384}, // Choose appropriate cipher suites
			InsecureSkipVerify: false,                                               // Set to true to skip certificate verification, but it is not recommended for production use.
		},
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
	}

	values := map[string]string{
		"fontip":      "YAT",
		"sfontur":     "",
		"fonkod":      string(name),
		"fongrup":     "",
		"bastarih":    string(oneWeekAgo.Format("02/01/2006")),
		"bittarih":    string(today),
		"fonturkod":   "",
		"fonunvantip": "",
	}

	jsonData, err := json.Marshal(values)

	response, err := client.Post(API_URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making GET request:", err)
	}
	fmt.Println("\n\n")
	fmt.Println(response.Body)
	fmt.Println("\n\n")
	defer response.Body.Close()

	x, err := ioutil.ReadAll(response.Body)
	var d Response
	err = json.Unmarshal([]byte(x), &d)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}
	var prices []float64
	var dates2 []*goment.Goment
	var isoDayOfWeek int
	for i := 0; i < 7; i++ {
		// Get the current time in the local timezone
		now2 := time.Now()

		// Calculate the duration to subtract from the current time to get the desired date
		temp := time.Duration(i) * 24 * time.Hour

		// Subtract the duration from the current time to get the desired date
		g, err := goment.New(now2.Add(-temp))
		if err != nil {
			fmt.Println(err)
		}

		// Get the ISO day of the week (Monday: 1, Tuesday: 2, ..., Sunday: 7)
		isoDayOfWeek = g.Weekday()

		if isoDayOfWeek != 0 && isoDayOfWeek != 6 {
			dates2 = append(dates2, g)
		}
	}
	var formattedString string = "\n<b>Son 1 Haftalık veriler</b>\n\n"
	for index, fund := range d.Data {

		prices = append(prices, fund.Fiyat)
		date := dates2[index].Format("DD/MM/YYYY")
		price := strconv.FormatFloat(fund.Fiyat, 'f', -1, 64)
		formattedString += "<b>" + date + "</b> : " + price + "\n"
	}

	if isoDayOfWeek == 0 || isoDayOfWeek == 6 {
		formattedString += "<b>Haftasonu borsa kapalı olduğundan dolayı veriler en erken pazartesi günü saat 09:30'da güncellenecektir.</b>"
	}

	return formattedString

}

func fetchFon(name string) string {

	url := "https://www.tefas.gov.tr/FonAnaliz.aspx?FonKod=" + name

	// Create a custom HTTP Transport with the specified options
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,                                    // Use TLS 1.2 or higher
			CipherSuites:       []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384}, // Choose appropriate cipher suites
			InsecureSkipVerify: false,                                               // Set to true to skip certificate verification, but it is not recommended for production use.
		},
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
	}

	response, err := client.Get(url)

	if err != nil {
		fmt.Println("Error making GET request:", err)
		return "<b>" + name + "</b> adlı fona ait veri bulunamamıştır."
	}
	defer response.Body.Close()

	// Check if the request was successful (status code 200)
	if response.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status code:", response.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	header := doc.Find("#MainContent_FormViewMainIndicators_LabelFund").Text()
	var title string = "<b>" + header + "</b>\n\n"

	if header == "Fon" {
		return "<b>" + name + "</b> adlı fona ait veri bulunamamıştır."
	}

	var prices [4]string
	doc.Find(".price-indicators > ul > li > span").Each(func(index int, item *goquery.Selection) {
		prices[index] = item.Text()
	})

	var periodicPrices string = "<b>Son 1 aylık getirisi</b> : " + prices[0] + "\n" + "<b>Son 3 aylık getirisi</b> : " + prices[1] + "\n" + "<b>Son 6 aylık getirisi</b> : " + prices[2] + "\n" + "<b>Son 1 yıllık getirisi</b> : " + prices[3] + "\n"
	mm := fonPrices(name)
	return title + periodicPrices + mm
}

func main() {
	Bot, err := TelegramBot.NewBotAPI("BOT_TOKEN")
	if err != nil {
		log.Panic(err)
	}
	Bot.Debug = true

	log.Printf("Authorized on account %s", Bot.Self.UserName)

	Update := TelegramBot.NewUpdate(0)
	Update.Timeout = 60

	Updates := Bot.GetUpdatesChan(Update)

	for update := range Updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if update.Message.NewChatMembers != nil {
			message := TelegramBot.NewMessage(update.Message.Chat.ID, "")
			message.ParseMode = "HTML"
			message.Text = "İyi günler, </b>/fon fonkodu</b> şeklinde verilere erişebilirsiniz.\\n\\n <b>Örnek kullanım : /fon GSP</b>"
			if _, err := Bot.Send(message); err != nil {
				log.Panic(err)
			}
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := TelegramBot.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = "HTML"
		btn1 := TelegramBot.NewInlineKeyboardButtonURL("Fon Takip Kanal", "https://t.me/fontakipbotu")
		btn2 := TelegramBot.NewInlineKeyboardButtonURL("Fon Takip Sohbet", "https://t.me/fontakipsohbet")

		// Create a new row for the buttons
		row := []TelegramBot.InlineKeyboardButton{btn1, btn2}

		// Create an inline keyboard markup
		inlineKeyboard := TelegramBot.NewInlineKeyboardMarkup(row)

		// Set the inline keyboard markup to the message
		msg.ReplyMarkup = inlineKeyboard
		var args string = strings.ToUpper(update.Message.CommandArguments())
		switch update.Message.Command() {
		case "start":
			msg.Text = "<b>Örnek kullanım : /fonkodu -> /fon GSP</b>"
		case "fon":
			if args == "" {
				msg.Text = "<b>Örnek kullanım : /fonkodu -> /fon GSP</b>"
			} else {
				msg.Text = fetchFon(args)
			}
		default:
			msg.Text = "<b>Örnek kullanım : /fonkodu -> /fon GSP</b>"
		}

		if _, err := Bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
