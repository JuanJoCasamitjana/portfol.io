package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/microcosm-cc/bluemonday"
	xhtml "golang.org/x/net/html"
)

var IMGBB_API_KEY string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	IMGBB_API_KEY = os.Getenv("IMGBB_API_KEY")
}

func uploadImageToImgbb(img []byte) (map[string]string, error) {
	res := make(map[string]string, 3)
	query_params := url.Values{}
	query_params.Add("key", IMGBB_API_KEY)
	req_url := "https://api.imgbb.com/1/upload" + "?" + query_params.Encode()
	req_body := &bytes.Buffer{}
	writer := multipart.NewWriter(req_body)
	part, err := writer.CreateFormFile("image", "image.jpg")
	if err != nil {
		return res, err
	}
	_, err = part.Write(img)
	if err != nil {
		return res, err
	}
	err = writer.Close()
	if err != nil {
		return res, err
	}
	resp, err := http.Post(req_url, writer.FormDataContentType(), req_body)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close() //skipcq GO-S2307
	var imgbbResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&imgbbResponse)
	if err != nil {
		return res, err
	}
	image_url := imgbbResponse["data"].(map[string]interface{})["url"].(string)
	delete_url := imgbbResponse["data"].(map[string]interface{})["delete_url"].(string)
	thumb_url := imgbbResponse["data"].(map[string]interface{})["thumb"].(map[string]interface{})["url"].(string)
	res["image_url"] = image_url
	res["delete_url"] = delete_url
	res["thumb_url"] = thumb_url
	return res, nil
}

func convertFileToBytes(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}

/* func deleteImageFromImgbb(delete_url string) error {
	req, err := http.NewRequest("DELETE", delete_url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
} */

func sanitizeHTML(htmlstr string) string {
	p := bluemonday.NewPolicy()
	p.AllowStandardURLs()
	p.AllowAttrs("href").OnElements("a")
	p.AllowElements("p", "h1", "h2", "h3", "h4", "h5", "h6", "strong", "i", "b",
		"em", "u", "s", "a", "img", "ul", "ol", "li", "blockquote", "code", "pre")
	p.AllowAttrs("src").OnElements("img")
	p.AllowAttrs("alt").OnElements("img")
	p.AllowAttrs("data-filename").OnElements("img")
	p.AllowAttrs("style").OnElements("img")
	p.RequireParseableURLs(true)

	return p.Sanitize(htmlstr)
}

func findAndUploadArticlesImages(htmlstr string) (string, error) {
	doc, err := xhtml.Parse(strings.NewReader(htmlstr))
	if err != nil {
		return "", err
	}
	new_html, err := navigateAndUploadImages(doc)
	if err != nil {
		return "", err
	}
	return new_html, nil
}

func htmlNodeToString(n *xhtml.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	xhtml.Render(w, n)
	return buf.String()
}

func processHTML(htmlstr string) (string, error) {
	htmlstr, err := findAndUploadArticlesImages(htmlstr)
	htmlstr = sanitizeHTML(htmlstr)
	if err != nil {
		return "", err
	}
	return htmlstr, nil
}

func navigateAndUploadImages(n *xhtml.Node) (string, error) {
	visited := make([]*xhtml.Node, 0)
	stack := make([]*xhtml.Node, 0)
	stack = append(stack, n)
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if node.Type == xhtml.ElementNode && node.Data == "img" {
			for i, attr := range node.Attr {
				if attr.Key == "src" && strings.HasPrefix(attr.Val, "data:image") {
					encoded_img := strings.Split(attr.Val, ",")[1]
					decoded_img, err := base64.StdEncoding.DecodeString(encoded_img)
					if err != nil {
						return "", err
					}
					imgbbResponse, err := uploadImageToImgbb(decoded_img)
					if err != nil {
						return "", err
					}
					new_url := imgbbResponse["image_url"]
					attr.Val = new_url
					node.Attr[i] = attr
				}
				if attr.Key == "style" {
					attr.Val = "max-width: 100%;"
					node.Attr[i] = attr
				}
			}
		}
		child := node.LastChild
		for child != nil {
			if !isIn(visited, child) {
				stack = append(stack, child)
			}
			child = child.PrevSibling
		}
		visited = append(visited, node)
	}

	new_html := htmlNodeToString(n)
	return new_html, nil
}

func isIn(slice []*xhtml.Node, node *xhtml.Node) bool {
	for i := range slice {
		if slice[i] == node {
			return true
		}
	}
	return false
}
