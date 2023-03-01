package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/jackc/pgx"
)

var responseCode = "response"

type responseBody struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Data    []eventResponse `json:"data"`
}

type eventResponse struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	DateTimeStart time.Time `json:"dateTimeStart"`
	DateTimeEnd   time.Time `json:"dateTimeEnd"`
	Description   string    `json:"description"`
	CreatedBy     int32     `json:"createdBy"`
	RemindFrom    time.Time `json:"remindFrom"`
}

type response struct {
	code      int
	body      *responseBody
	bodyPlain []byte
}

type apiFeature struct {
	client *http.Client
	host   string
	port   string
	conn   *pgx.Conn
}

func (a *apiFeature) sendRequestToWithData(ctx context.Context, method, endpoint string, data *godog.DocString) (context.Context, error) {
	url := fmt.Sprintf("%s:%s%s", a.host, a.port, endpoint)
	req, errReq := http.NewRequest(method, url, strings.NewReader(data.Content))
	if errReq != nil {
		return ctx, errReq
	}

	resHTTP, errDo := a.client.Do(req)
	if errDo != nil {
		return ctx, errDo
	}
	res, errP := parseHTTPResponse(resHTTP)
	if errP != nil {
		return ctx, errP
	}
	ctx = context.WithValue(ctx, responseCode, res)
	return ctx, nil
}

func (a *apiFeature) responseCodeShouldBe(ctx context.Context, code int) (context.Context, error) {
	res, err := getResponse(ctx)
	if err != nil {
		return ctx, err
	}
	if code != res.code {
		return ctx, fmt.Errorf("expected code=%d , but code=%d", code, res.code)
	}
	return ctx, nil
}

func (a *apiFeature) responseDataHasCode(ctx context.Context, code string) (context.Context, error) {
	res, err := getResponse(ctx)
	if err != nil {
		return ctx, err
	}
	if res.body.Code != code {
		return ctx, fmt.Errorf("expected response data code=%s , but code=%s", code, res.body.Code)
	}

	return ctx, nil
}

func (a *apiFeature) sendRequestTo(ctx context.Context, method, endpoint string) (context.Context, error) {
	url := fmt.Sprintf("%s:%s%s", a.host, a.port, endpoint)
	req, errReq := http.NewRequest(method, url, nil)
	if errReq != nil {
		return ctx, errReq
	}
	resHTTP, err := a.client.Do(req)
	if err != nil {
		return ctx, err
	}
	res, errP := parseHTTPResponse(resHTTP)
	if errP != nil {
		return ctx, errP
	}
	ctx = context.WithValue(ctx, responseCode, res)

	return ctx, nil
}

func (a *apiFeature) hasEventWithTitle(ctx context.Context, title string) (context.Context, error) {
	res, err := getResponse(ctx)
	if err != nil {
		return ctx, err
	}
	if res.body.Data == nil {
		return ctx, errors.New("not found events")
	}
	for _, e := range res.body.Data {
		if e.Title == title {
			return ctx, nil
		}
	}
	return ctx, fmt.Errorf("not fount events with title %s", title)
}

func (a *apiFeature) hasNotEventWithTitle(ctx context.Context, title string) (context.Context, error) {
	res, err := getResponse(ctx)
	if err != nil {
		return ctx, err
	}
	if res.body.Data == nil {
		return ctx, nil
	}
	for _, e := range res.body.Data {
		if e.Title == title {
			return ctx, fmt.Errorf("fount events with title %s", title)
		}
	}
	return ctx, nil
}

func (a *apiFeature) findInLogEventWithTitle(ctx context.Context, title string) (context.Context, error) {
	rows, err := a.conn.Query("SELECT body FROM public.logs")
	if err != nil {
		return ctx, err
	}
	defer rows.Close()
	for rows.Next() {
		var body string
		var event eventResponse
		if errS := rows.Scan(&body); errS != nil {
			return ctx, errS
		}
		if errU := json.Unmarshal([]byte(body), &event); errU != nil {
			return ctx, errU
		}
		if event.Title == title {
			return ctx, nil
		}

	}
	if errR := rows.Err(); errR != nil {
		return nil, errR
	}

	return ctx, fmt.Errorf("not fount notification logs with event title %s", title)
}

func (a *apiFeature) notFindInLogEventWithTitle(ctx context.Context, title string) (context.Context, error) {
	rows, err := a.conn.Query("SELECT body FROM public.logs")
	if err != nil {
		return ctx, err
	}
	defer rows.Close()
	for rows.Next() {
		var body string
		var event eventResponse
		if errS := rows.Scan(&body); errS != nil {
			return ctx, errS
		}
		if errU := json.Unmarshal([]byte(body), &event); errU != nil {
			return ctx, errU
		}
		if event.Title == title {
			return ctx, fmt.Errorf("fount notification logs with event title %s", title)
		}

	}
	if errR := rows.Err(); errR != nil {
		return nil, errR
	}

	return ctx, nil
}

func (a *apiFeature) waitWhenSchedulerSendAllNotification(ctx context.Context, duration string) (context.Context, error) {
	dur, err := time.ParseDuration(duration)
	if err != nil {
		return ctx, err
	}
	time.Sleep(dur)
	return ctx, nil
}

func (a *apiFeature) connectToDB(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	if a.conn != nil {
		return ctx, nil
	}
	port, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		return ctx, err
	}
	conf := pgx.ConnConfig{
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     uint16(port),
		Database: os.Getenv("DATABASE_DB"),
	}
	conn, err := pgx.Connect(conf)
	if err != nil {
		return ctx, err
	}
	a.conn = conn

	return ctx, err
}

func (a *apiFeature) Close(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
	if a.conn == nil {
		return ctx, nil
	}
	if errC := a.conn.Close(); errC != nil {
		return ctx, errC
	}
	a.conn = nil
	return ctx, nil
}

func parseHTTPResponse(resHTTP *http.Response) (res response, err error) {
	res.code = resHTTP.StatusCode
	res.bodyPlain, err = io.ReadAll(resHTTP.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(res.bodyPlain, &res.body); err != nil {
		return
	}
	return res, nil
}

func getResponse(ctx context.Context) (*response, error) {
	res, ok := ctx.Value(responseCode).(response)
	if !ok {
		return nil, errors.New("not found response in context")
	}
	return &res, nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	api := &apiFeature{
		host:   os.Getenv("CALENDAR_HOST"),
		port:   os.Getenv("CALENDAR_PORT"),
		client: &http.Client{},
	}
	ctx.Before(api.connectToDB)
	ctx.Step(`^I send "([^"]*)" request to "([^"]*)" with data$`, api.sendRequestToWithData)
	ctx.Step(`^I response code should be (\d+)$`, api.responseCodeShouldBe)
	ctx.Step(`^Response data has code "([^"]*)"$`, api.responseDataHasCode)
	ctx.Step(`^I send "([^"]*)" request to "([^"]*)"$`, api.sendRequestTo)
	ctx.Step(`^Has event with title "([^"]*)"$`, api.hasEventWithTitle)
	ctx.Step(`^Has not event with title "([^"]*)"$`, api.hasNotEventWithTitle)
	ctx.Step(`^Find in log event with title "([^"]*)"$`, api.findInLogEventWithTitle)
	ctx.Step(`^Wait "([^"]*)" when scheduler send all notification$`, api.waitWhenSchedulerSendAllNotification)
	ctx.Step(`^Not find in log event with title "([^"]*)"$`, api.notFindInLogEventWithTitle)
	ctx.After(api.Close)
}

func TestMain(m *testing.M) {
	status := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:    "progress", // Замените на "pretty" для лучшего вывода
			Paths:     []string{"features"},
			Randomize: 0, // Последовательный порядок исполнения
		},
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
