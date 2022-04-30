package handlers

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"fakebook/internal/backend"
)

type HomePage struct {
	backend      *backend.Backend
	basicURL     string
	pageTemplate string
}

func NewHomePage(backend *backend.Backend, basicURL string) *HomePage {
	return &HomePage{
		backend:  backend,
		basicURL: basicURL,
	}
}

func (h *HomePage) Handle(ctx *gin.Context) {
	username := ctx.Param("username")

	userProfile, err := h.backend.GetProfileByUsername(username)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if userProfile == nil {
		ctx.String(http.StatusNotFound, "not found :(")
		return
	}

	htmlPage, err := h.renderPage(userProfile)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Data(http.StatusOK, "text/html", []byte(htmlPage))
}

func (h *HomePage) renderPage(userProfile *backend.UserProfile) (string, error) {
	htmlPage, err := h.newTemplate()
	if err != nil {
		return "", fmt.Errorf("cannot create new page template: %w", err)
	}

	name := fmt.Sprint(userProfile.FirstName, " ", userProfile.LastName)
	htmlPage = strings.Replace(htmlPage, "${name}", name, 1)

	htmlPage = renderDateOfBirth(htmlPage, userProfile.DateOfBirth)
	htmlPage = renderCity(htmlPage, userProfile.City)
	htmlPage = renderAboutMe(htmlPage, userProfile.Info)

	return htmlPage, nil
}

func (h *HomePage) newTemplate() (string, error) {
	if h.pageTemplate != "" {
		return h.pageTemplate, nil
	}

	data, err := os.ReadFile("site/home.html")
	if err != nil {
		return "", fmt.Errorf("cannot read page template from file %w", err)
	}

	h.pageTemplate = strings.Replace(string(data), "${base_url}", h.basicURL, 1)

	return h.pageTemplate, nil
}

func renderDateOfBirth(htmlTemplate, dateOfBirth string) (newTemplate string) {
	if dateOfBirth != "" {
		return renderProfileRow(htmlTemplate, "date_of_birth", dateOfBirth)
	}
	return hideProfileRow(htmlTemplate, "date_of_birth")
}

func renderCity(htmlTemplate, city string) (newTemplate string) {
	if city != "" {
		return renderProfileRow(htmlTemplate, "city", city)
	}
	return hideProfileRow(htmlTemplate, "city")
}

func renderAboutMe(htmlPage, text string) (newTemplate string) {
	if text == "" {
		return hideProfileRow(htmlPage, "about_me")
	}

	paragraphs := splitIntoParagraphs(text)

	htmlText := strings.Builder{}
	for _, paragraph := range paragraphs {
		escaped := html.EscapeString(paragraph)
		escaped = strings.ReplaceAll(escaped, "\n", "<br>")
		p := fmt.Sprintf("<p>%s</p>\n", escaped)
		htmlText.WriteString(p)
	}

	return renderProfileRow(htmlPage, "about_me", htmlText.String())
}

func splitIntoParagraphs(text string) []string {
	lines := strings.Split(text, "\n")

	paragraphs := make([]string, 0)

	currentParagraph := strings.Builder{}
	readingParagraph := false
	readingDelimiter := true

	for _, line := range lines {
		switch {
		case readingParagraph && line != "":
			currentParagraph.WriteString(line)
			currentParagraph.WriteByte('\n')

		case readingParagraph && line == "":
			paragraphs = append(paragraphs, currentParagraph.String())
			readingParagraph = false
			readingDelimiter = true
			currentParagraph.Reset()

		case readingDelimiter && line != "":
			currentParagraph.WriteString(line)
			currentParagraph.WriteByte('\n')
			readingDelimiter = false
			readingParagraph = true

		case readingDelimiter && line == "":
			break
		}
	}

	if currentParagraph.Len() != 0 || len(paragraphs) == 0 {
		paragraphs = append(paragraphs, currentParagraph.String())
	}

	return paragraphs
}

func renderProfileRow(htmlTemplate, variable, value string) (newTemplate string) {
	visibility := fmt.Sprintf("${hide_%s}", variable)
	htmlTemplate = strings.Replace(htmlTemplate, visibility, "", 1)

	variable = fmt.Sprintf("${%s}", variable)
	htmlTemplate = strings.Replace(htmlTemplate, variable, value, 1)

	return htmlTemplate
}

func hideProfileRow(htmlTemplate, variable string) (newTemplate string) {
	visibility := fmt.Sprintf("${hide_%s}", variable)
	htmlTemplate = strings.Replace(htmlTemplate, visibility, "hidden", 1)

	variable = fmt.Sprintf("${%s}", variable)
	htmlTemplate = strings.Replace(htmlTemplate, variable, "", 1)

	return htmlTemplate
}
