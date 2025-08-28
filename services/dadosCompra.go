package services

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/filipegms5/MoneyFlow-Backend/models"
	"golang.org/x/net/html"
)

var transacao models.Transacao
var estabelecimento models.Estabelecimento
var formaPagamento models.FormaPagamento

func FetchTransacao(url string) (models.Transacao, error) {
	doc, err := fetch(url)
	if err != nil {
		return models.Transacao{}, err
	}
	transacao = models.Transacao{}
	scrapeAll(doc)
	transacao.Estabelecimento = &estabelecimento
	transacao.FormaPagamento = &formaPagamento

	return transacao, err
}

// Fetches and parses the HTML document
func fetch(url string) (*html.Node, error) {
	// Create a custom HTTP client with TLS configuration to skip SSL verification
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: customTransport}

	// Send a GET request
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func scrapeAll(n *html.Node) {

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {

		storeName(n)

		// Store address
		storeAdress(n)

		// Date
		date(n)

		// Products
		//products(n)

		// Sale info
		saleInfo(n)

		// Traverse children
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
}

func storeName(n *html.Node) {
	// Store name
	if n.Type == html.ElementNode && n.Data == "h4" {
		for b := n.FirstChild; b != nil; b = b.NextSibling {
			if b.Type == html.ElementNode && b.Data == "b" && b.FirstChild != nil {
				estabelecimento.RazaoSocial = b.FirstChild.Data
			}
		}
	}
	// Store CNPJ
	if n.Type == html.ElementNode && n.Data == "td" {
		for _, attr := range n.Attr {
			if attr.Key == "style" && attr.Val == "border-top: 0px;" && n.FirstChild != nil {
				if strings.Contains(n.FirstChild.Data, "CNPJ:") {
					cnpjRegex := regexp.MustCompile(`\d{2}\.\d{3}\.\d{3}/\d{4}-\d{2}`)
					cnpj := cnpjRegex.FindString(n.FirstChild.Data)
					estabelecimento.CNPJ = cnpj
				}
			}
		}
	}
}

func date(n *html.Node) {
	dateRegex := regexp.MustCompile(`\d{2}/\d{2}/\d{4} \d{2}:\d{2}:\d{2}`)
	if n.Type == html.TextNode {
		if date := strings.TrimSpace(n.Data); dateRegex.MatchString(date) {

			// Parse the input date string
			layoutIn := "02/01/2006 15:04:05"
			t, err := time.Parse(layoutIn, date)
			if err != nil {
				panic(err)
			}

			// Format to ISO 8601 (UTC with Z)
			layoutOut := "2006-01-02T15:04:05Z"
			formatted := t.UTC().Format(layoutOut)

			transacao.Data = formatted

		}
	}
}

func storeAdress(n *html.Node) {
	// Store address
	if n.Type == html.ElementNode && n.Data == "td" {
		for _, attr := range n.Attr {
			if attr.Key == "style" && attr.Val == "border-top: 0px; display: block; font-style: italic;" && n.FirstChild != nil {
				estabelecimento.Endereco = n.FirstChild.Data
			}
		}
	}
}

func saleInfo(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "strong" && n.FirstChild != nil {
		text := n.FirstChild.Data
		if strings.Contains(text, ".") {
			var err error
			transacao.Valor, err = strconv.ParseFloat(text, 64)
			if err != nil {
				fmt.Printf("Error converting string to float64: %v\n", err)
				return
			}

		} else if n.FirstChild.Type == html.ElementNode && n.FirstChild.Data == "div" {
			formaPagamento.Nome = n.FirstChild.FirstChild.Data
		}
	}
}

// func products(n *html.Node) {
// 	if n.Type == html.ElementNode && n.Data == "tr" {
// 		produto := models.Produto{}
// 		count := 0
// 		for td := n.FirstChild; td != nil; td = td.NextSibling {
// 			if td.Type == html.ElementNode && td.Data == "td" {
// 				text := ""
// 				if td.FirstChild != nil {
// 					if td.FirstChild.Type == html.ElementNode && td.FirstChild.Data == "h7" {
// 						text = td.FirstChild.FirstChild.Data
// 						produto.Nome = strings.Split(text, "\n")[0]
// 					} else {
// 						text = td.FirstChild.Data
// 						switch count {
// 						case 1:
// 							parts := strings.Split(text, ":")
// 							if len(parts) > 1 {
// 								produto.Quantidade = strings.TrimSpace(parts[1])
// 							}
// 						case 2:
// 							parts := strings.Split(text, ":")
// 							if len(parts) > 1 {
// 								produto.Unidade = strings.TrimSpace(parts[1])
// 							}
// 						case 3:
// 							parts := strings.Split(text, ": R$ ")
// 							if len(parts) > 1 {
// 								produto.Valor = strings.TrimSpace(parts[1])
// 							}
// 						}
// 					}
// 				}
// 				count++
// 			}
// 		}
// 		if produto.Nome != "" {
// 			transacao.Produtos = append(transacao.Produtos, produto)
// 		}
// 	}
// }
