package session

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type GachaSessionOptions struct {
	UserAgent string
	Timeout   time.Duration
}

type GachaSession struct {
	client           *http.Client
	userAgent        string
	currentCharacter string
	character1       string
	character2       string
	serial           string
	rollCount        int
}

func NewGachaSession() (*GachaSession, error) {
	return NewGachaSessionWithOptions(GachaSessionOptions{
		Timeout: 5 * time.Second,
	})
}

func NewGachaSessionWithOptions(options GachaSessionOptions) (*GachaSession, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	s := &GachaSession{
		client: &http.Client{
			Jar:     jar,
			Timeout: options.Timeout,
		},
		userAgent: options.UserAgent,
		rollCount: 0,
	}

	if err := s.init(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *GachaSession) newRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if s.userAgent != "" {
		req.Header.Add("User-Agent", s.userAgent)
	}

	return req, nil
}

func (s *GachaSession) init() error {
	req, err := s.newRequest(http.MethodGet, gachaInitialURL, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrNotOK
	}

	return nil
}

func (s *GachaSession) Roll() error {
	if s.serial != "" {
		return ErrAlreadyGotSerial
	}

	if s.rollCount >= RollLimit {
		return ErrRollLimitExceeded
	}

	req, err := s.newRequest(http.MethodGet, gachaRollURL, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrNotOK
	}

	dom, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	resultsImg := dom.Find(`.result`)
	resultsCount := len(resultsImg.Nodes)

	if resultsCount == 1 {
		// first round
		src, _ := resultsImg.First().Attr("src")
		s.currentCharacter = imgSrcToCharacter(src)
		s.character1 = s.currentCharacter
	} else {
		// other round
		src1, _ := resultsImg.First().Attr("src")
		src2, _ := resultsImg.Last().Attr("src")
		s.character1 = imgSrcToCharacter(src1)
		s.character2 = imgSrcToCharacter(src2)
		s.currentCharacter = s.character1
	}

	s.rollCount++

	return nil
}

func (s *GachaSession) GetCurrentCharacter() string {
	return s.currentCharacter
}

func (s *GachaSession) GetAlternateCharacter() string {
	if s.character2 == "" {
		return ""
	}
	switch s.currentCharacter {
	case s.character1:
		return s.character2
	case s.character2:
		return s.character1
	}
	return ""
}

func (s *GachaSession) GetCharacter1() string {
	return s.character1
}

func (s *GachaSession) GetCharacter2() string {
	return s.character2
}

func (s *GachaSession) SelectCharacter1() error {
	if s.currentCharacter == "" {
		return ErrNeverRolled
	}

	if s.currentCharacter == s.character1 {
		return nil
	}

	if s.serial != "" {
		return ErrAlreadyGotSerial
	}

	req, err := s.newRequest(http.MethodPost, gachaSelectCharacter1URL, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrNotOK
	}

	s.currentCharacter = s.character1

	return nil
}

func (s *GachaSession) SelectCharacter2() error {
	if s.currentCharacter == "" {
		return ErrNeverRolled
	}

	if s.currentCharacter == s.character2 {
		return nil
	}

	if s.serial != "" {
		return ErrAlreadyGotSerial
	}

	req, err := s.newRequest(http.MethodPost, gachaSelectCharacter2URL, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrNotOK
	}

	s.currentCharacter = s.character2

	return nil
}

func (s *GachaSession) SwitchCharacter() error {
	if s.currentCharacter == "" {
		return ErrNeverRolled
	}

	switch s.currentCharacter {
	case s.character1:
		return s.SelectCharacter2()
	case s.character2:
		return s.SelectCharacter1()
	}

	return ErrUnknown
}

func (s *GachaSession) ObtainSerial() error {
	if s.currentCharacter == "" {
		return ErrNeverRolled
	}

	if s.serial != "" {
		return nil
	}

	req, err := s.newRequest(http.MethodGet, gachaObtainSerialURL, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	dom, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	serial, _ := dom.Find(`#code`).First().Attr("value")

	s.serial = serial

	return nil
}

func (s *GachaSession) GetSerial() string {
	return s.serial
}

func (s *GachaSession) GetRollCount() int {
	return s.rollCount
}
