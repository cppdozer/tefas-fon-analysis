package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	TelegramBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nleeper/goment"
	"io"
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

	Today := time.Now()                                                  // Anlık tarih verisi alındı.
	twoWeek := 13 * 24 * time.Hour                                       // 14 günlük zaman dilimi atandı.
	twoWeekAgo := Today.Add(-twoWeek)                                    // 14 gün öncesinin tarihi anlık tarihten çıkarılarak elde edildi.
	MinerUrl := "https://www.tefas.gov.tr/FonAnaliz.aspx?FonKod=" + name // Verilerin çekileceği adres
	ApiUrl := "https://www.tefas.gov.tr/api/DB/BindHistoryInfo"          // Sabit API adresi
	var prices []float64                                                 // API'dan dönen verilerden FIYAT içeriğini bu array'e ekleyecek
	var validDates []*goment.Goment                                      // Borsanın açık olduğu gün ve saatler dahilindeki tarihler indexlenecek.
	var isoDayOfWeek int                                                 // Tarih kontrol parametresi
	var PostJson Response
	var periodicPrices [4]string

	values := map[string]string{ // Gönderilecek verinin JSON formatına dönüştürülmeden önce map ile oluşturulmuş hali.
		"fontip":      "YAT",
		"sfontur":     "",
		"fonkod":      name,
		"fongrup":     "",
		"bastarih":    twoWeekAgo.Format("02/01/2006"),
		"bittarih":    Today.Format("02/01/2006"),
		"fonturkod":   "",
		"fonunvantip": "",
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,                                    // Use TLS 1.2 or higher
			CipherSuites:       []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384}, // Choose appropriate cipher suites
			InsecureSkipVerify: false,
		},
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableKeepAlives:   true, // İşi biten bağlantının düşürülmesi için
	}

	client := &http.Client{
		Transport: transport,
	}

	formData, err := json.Marshal(values) // Gönderilecek veriler json formatına dönüştürülüyor.

	// HTTP Request işlemi yapılıyor.
	PostResponse, err := client.Post(ApiUrl, "application/json", bytes.NewBuffer(formData))
	if err != nil {
		fmt.Println("Error making POST request:", err)
	}

	GetResponse, err := client.Get(MinerUrl)
	if err != nil {
		fmt.Println("Error making GET request:", err)
	}

	Document, err := goquery.NewDocumentFromReader(GetResponse.Body)

	if err != nil {
		log.Fatal(err)
	}

	Title := Document.Find("#MainContent_FormViewMainIndicators_LabelFund").Text()

	if Title == "Fon" {
		return "<b>" + name + "</b> adlı fona ait veri bulunamamıştır."
	}

	Header := "<b>" + Title + "</b>"
	DailyProfit := Document.Find("ul.top-list > li:nth-child(2) > span").Text()
	Header += "\n\n<b>Günlük getirisi</b> : " + DailyProfit + "\n\n"

	Document.Find(".price-indicators > ul > li > span").Each(func(index int, item *goquery.Selection) {
		periodicPrices[index] = item.Text()
	})

	Body := "<b>Son 1 aylık getirisi</b> : " + periodicPrices[0] + "\n" + "<b>Son 3 aylık getirisi</b> : " + periodicPrices[1] + "\n" + "<b>Son 6 aylık getirisi</b> : " + periodicPrices[2] + "\n" + "<b>Son 1 yıllık getirisi</b> : " + periodicPrices[3] + "\n"

	// Fonksiyon return etmeden önce yapılacaklar
	defer func(PostBody io.ReadCloser, GetBody io.ReadCloser) {
		err := PostBody.Close()
		if err != nil {
			fmt.Println(err)
		}

		err = GetBody.Close()
		if err != nil {
			fmt.Println(err)
		}

	}(PostResponse.Body, GetResponse.Body)

	// Byte olarak dönen saf veri burada Response veri tipine dönüştürülüyor.
	apiRawData, err := io.ReadAll(PostResponse.Body)
	err = json.Unmarshal([]byte(apiRawData), &PostJson)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}

	beforeOpening := time.Date(Today.Year(), Today.Month(), Today.Day(), 9, 30, 0, 0, Today.Location())
	afterClosing := time.Date(Today.Year(), Today.Month(), Today.Day(), 18, 00, 0, 0, Today.Location())

	for i := 0; i < 14; i++ {

		temp := time.Duration(i) * 24 * time.Hour

		subtractedDate, err := goment.New(Today.Add(-temp))
		if err != nil {
			fmt.Println(err)
		}

		isoDayOfWeek = subtractedDate.Weekday()

		if isoDayOfWeek != 0 && isoDayOfWeek != 6 {
			if i == 0 && Today.Before(beforeOpening) {
				// pass
			} else {
				validDates = append(validDates, subtractedDate)
			}
		}
	}

	Footer := "\n<b>Son 14 Günlük veriler</b>\n\n"

	for index, fund := range PostJson.Data {

		prices = append(prices, fund.Fiyat)
		date := validDates[index].Format("DD/MM/YYYY")
		price := strconv.FormatFloat(fund.Fiyat, 'f', -1, 64)
		Footer += "<b>" + date + "</b> : " + price + "\n"
	}

	if isoDayOfWeek == 0 || isoDayOfWeek == 6 {
		Footer += "<b>\nHaftasonu borsa kapalı olduğundan dolayı veriler en erken pazartesi günü saat 09:30'da güncellenecektir.</b>"
	} else {
		if Today.Before(beforeOpening) || Today.After(afterClosing) {
			Footer += "\n<b>Borsanın durumu</b>: KAPALI\n\n"
			Footer += "\n<b>Hatırlatma</b> : Borsa haftaiçi saat 09:30 ile 18:00 saatleri arasında açıktır."
		} else {
			Footer += "\n<b>Borsanın durumu</b>: AÇIK"
		}
	}

	return Header + Body + Footer
}

func main() {
	Bot, err := TelegramBot.NewBotAPI("6404844010:AAFdHSrMgHvX0NUjcSb0Lx9uzMoqEcLsOUM")
	Bot.Debug = true

	if err != nil {
		log.Panic(err)
	}

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

		Message := TelegramBot.NewMessage(update.Message.Chat.ID, "")
		Message.ParseMode = "HTML"
		ChannelButton := TelegramBot.NewInlineKeyboardButtonURL("Fon Takip Kanal", "https://t.me/fontakipbotu")
		GroupButton := TelegramBot.NewInlineKeyboardButtonURL("Fon Takip Sohbet", "https://t.me/fontakipsohbet")

		ButtonField := []TelegramBot.InlineKeyboardButton{ChannelButton, GroupButton}

		// Create an inline keyboard markup
		inlineKeyboard := TelegramBot.NewInlineKeyboardMarkup(ButtonField)

		// Set the inline keyboard markup to the message
		Message.ReplyMarkup = inlineKeyboard
		args := strings.ToUpper(update.Message.CommandArguments())
		switch update.Message.Command() {
		case "start":
			Message.Text = "<b>Örnek kullanım : /fonkodu -> /fon GSP</b>"
		case "fon":
			if args == "" {
				Message.Text = "<b>Örnek kullanım : /fonkodu -> /fon GSP</b>"
			} else {
				Message.Text = fonPrices(args)
			}
		default:
			Message.Text = "<b>Örnek kullanım : /fonkodu -> /fon GSP</b>"
		}

		if _, err := Bot.Send(Message); err != nil {
			log.Panic(err)
		}
	}
}
