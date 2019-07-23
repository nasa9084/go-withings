package withings

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/xerrors"
)

const (
	DefaultWithingsAPIEndpoint = "https://wbsapi.withings.net"
)

type Logger interface {
}

// ensure nilLogger implements Logger interface
var _ Logger = (*nilLogger)(nil)

type nilLogger struct{}

type httpClient struct {
	client   *http.Client
	endpoint string
	logger   Logger
}

func (c *httpClient) makeURL(path string) (*url.URL, error) {
	u, err := url.Parse(c.endpoint + path)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (c *httpClient) get(ctx context.Context, path string, query url.Values, data interface{}) error {
	u, err := c.makeURL(path)
	if err != nil {
		return xerrors.Errorf("making url: %w", err)
	}
	u.RawQuery = query.Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return xerrors.Errorf("creating new request: %w", err)
	}
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return xerrors.Errorf("doing request: %w", err)
	}
	defer closeResponse(resp)

	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return xerrors.Errorf("decoding response: %w", err)
	}
	return nil
}

func closeResponse(resp *http.Response) {
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
}

type Client struct {
	user    *UserService
	measure *MeasureService
	sleep   *SleepService
	notify  *NotifyService
}

func New(opts ...Option) *Client {
	httpcl := &httpClient{
		client:   http.DefaultClient,
		endpoint: DefaultWithingsAPIEndpoint,
		logger:   nilLogger{},
	}
	for _, opt := range opts {
		opt(httpcl)
	}
	return &Client{
		user:    &UserService{client: httpcl},
		measure: &MeasureService{client: httpcl},
		sleep:   &SleepService{client: httpcl},
		notify:  &NotifyService{client: httpcl},
	}
}

func (c *Client) User() *UserService       { return c.user }
func (c *Client) Measure() *MeasureService { return c.measure }

// Option allows
type Option func(*httpClient)

// WithEndpoint allows to use specific withings api endpoint.
func WithEndpoint(endpoint string) Option {
	return func(c *httpClient) {
		c.endpoint = endpoint
	}
}

// WithHTTPClient allows to use specific *http.Client.
func WithHTTPClient(httpcl *http.Client) Option {
	return func(c *httpClient) {
		c.client = httpcl
	}
}

type UserService struct {
	client *httpClient
}

func (svc *UserService) GetDevice(ctx context.Context) (UserGetDeviceResponse, error) {
	query := url.Values{}
	query.Set("action", "getdevice")
	var resp UserGetDeviceResponse
	if err := svc.client.get(ctx, "/v2/user", query, &resp); err != nil {
		return UserGetDeviceResponse{}, err
	}
	return resp, nil
}

type UserGetDeviceResponse struct {
	Status int
	Body   UserGetDeviceResponseBody
}

type UserGetDeviceResponseBody struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	Type     DeviceType  `json:"type"`
	Model    DeviceModel `json:"model"`
	Battery  string      `json:"battery"`
	DeviceID string      `json:"deviceid"`
	Timezone string      `json:"timezone"`
}

// DeviceType is a type of the device.
type DeviceType string

// DeviceType Constants.
const (
	Scale                     DeviceType = "Scale"
	Babyphone                 DeviceType = "Babyphone"
	BloodPressureMonitor      DeviceType = "Blood Pressure Monitor"
	ActivityTracker           DeviceType = "Activity Tracker"
	SleepMonitor              DeviceType = "Sleep Monitor"
	SmartConnectedThermometer DeviceType = "Smart Connected Thermometer"
)

// DeviceModel is a device model.
type DeviceModel string

// DeviceModel Constants.
const (
	WithingsWBS01                  DeviceModel = "Withings WBS01"
	WS30                           DeviceModel = "WS30"
	KidScale                       DeviceModel = "Kid Scale"
	SmartBodyAnalyzer              DeviceModel = "Smart Body Analyzer"
	BodyPlus                       DeviceModel = "Body+"
	BodyCardio                     DeviceModel = "Body Cardio"
	Boady                          DeviceModel = "Boady"
	SmartBabyMonitor               DeviceModel = "Smart Baby Monitor"
	WhithingsHome                  DeviceModel = "Whithings Home"
	WithingsBloodPressureMonitorV1 DeviceModel = "Withings Blood Pressure Monitor V1"
	WithingsBloodPressureMonitorV2 DeviceModel = "Withings Blood Pressure Monitor V2"
	WithingsBloodPressureMonitorV3 DeviceModel = "Withings Blood Pressure Monitor V3"
	Pulse                          DeviceModel = "Pulse"
	Activite                       DeviceModel = "Activite"
	ActivitePopSteel               DeviceModel = "Activite (Pop, Steel)"
	WithingsGo                     DeviceModel = "Withings Go"
	ActiviteSteelHR                DeviceModel = "Activite Steel HR"
	ActiviteSteelHRSportEdition    DeviceModel = "Activite Steel HR Sport Edition"
	PulseHR                        DeviceModel = "Pulse HR"
	AuraDock                       DeviceModel = "Aura Dock"
	AuraSensor                     DeviceModel = "Aura Sensor"
	AuraSensorV2                   DeviceModel = "Aura Sensor V2"
	Thermo                         DeviceModel = "Thermo"
)

type MeasureService struct {
	client *httpClient
}

func (svc *MeasureService) GetMeas(ctx context.Context, meastype MeasureType, category MeasureCategory, startdate, enddate time.Time, offset int, lastupdate time.Time) (MeasureGetMeasResponse, error) {
	query := url.Values{}
	query.Set("action", "getmeas")
	query.Set("meastype", strconv.Itoa(int(meastype)))
	query.Set("category", strconv.Itoa(int(category)))
	query.Set("startdate", strconv.FormatInt(startdate.Unix(), 10))
	query.Set("enddate", strconv.FormatInt(enddate.Unix(), 10))
	query.Set("offset", strconv.Itoa(offset))
	query.Set("lastupdate", strconv.FormatInt(lastupdate.Unix(), 10))
	var resp MeasureGetMeasResponse
	if err := svc.client.get(ctx, "/measure", query, &resp); err != nil {
		return MeasureGetMeasResponse{}, err
	}
	return resp, nil
}

type MeasureType int

// MeasureType Constants.
const (
	UnknownMeasureType     MeasureType = 0
	Weight                 MeasureType = 1
	Height                 MeasureType = 4
	FatFreeMass            MeasureType = 5
	FatRatio               MeasureType = 6
	FatMassWeight          MeasureType = 8
	DiastolicBloodPressure MeasureType = 9
	SystolicBloodPressure  MeasureType = 10
	HeartPulse             MeasureType = 11
	Temperature            MeasureType = 12
	SP02                   MeasureType = 54
	BodyTemperature        MeasureType = 71
	SkinTemperature        MeasureType = 73
	MuscleMass             MeasureType = 76
	Hydration              MeasureType = 77
	BoneMass               MeasureType = 88
	PulsWaveVelocity       MeasureType = 91
)

type MeasureCategory int

// MeasureCategory Constants.
const (
	UnknownMeasureCategory MeasureCategory = 0
	RealMeasure            MeasureCategory = 1
	UserObjective          MeasureCategory = 2
)

type MeasureGetMeasResponse struct {
	Status int
	Body   MeasureGetMeasResponseBody
}

type MeasureGetMeasResponseBody struct {
	UpdateTime    string         `json:"updatetime"`
	Timezone      string         `json:"timezone"`
	MeasureGroups []MeasureGroup `json:"measuregrps"`
	More          bool           `json:"more"`
	Offset        int            `json:"offset"`
}

type MeasureGroup struct {
	GroupID   int             `json:"grpid"`
	Attribute int             `json:"attrib"` // this is constant, but they do not have name
	Date      time.Time       `json:"date"`
	Created   time.Time       `json:"created"`
	Category  MeasureCategory `json:"category"`
	DeviceID  string          `json:"deviceid"`
	Measures  []Measure       `json:"measures"`
	Comment   string          `json:"comment"` // deprecated
}

func (m *MeasureGroup) UnmarshalJSON(data []byte) error {
	proxy := struct {
		GroupID   int             `json:"grpid"`
		Attribute int             `json:"attrib"`
		Date      int64           `json:"date"`
		Created   int64           `json:"created"`
		Category  MeasureCategory `json:"category"`
		DeviceID  string          `json:"deviceid"`
		Measures  []Measure       `json:"measures"`
		Comment   string          `json:"comment"`
	}{}
	if err := json.Unmarshal(data, &proxy); err != nil {
		return err
	}
	m.GroupID = proxy.GroupID
	m.Attribute = proxy.Attribute
	m.Date = time.Unix(proxy.Date, 0)
	m.Created = time.Unix(proxy.Created, 0)
	m.Category = proxy.Category
	m.DeviceID = proxy.DeviceID
	m.Measures = proxy.Measures
	m.Comment = proxy.Comment
	return nil
}

type Measure struct {
	Value int         `json:"value"`
	Type  MeasureType `json:"type"`
	Unit  int         `json:"unit"`
	Algo  int         `json:"algo"` // deprecated
	Fm    int         `json:"fm"`   // deprecated
	Fw    int         `json:"fw"`   // deprecated
}

type SleepService struct {
	client *httpClient
}

type NotifyService struct {
	client *httpClient
}

func (svc *MeasureService) GetActivity(ctx context.Context, startdate, enddate time.Time, offset int, dataFields []ActivityDataField, lastupdate time.Time) (MeasureGetActivityResponse, error) {
	query := url.Values{}
	query.Set("action", "getactivity")
	query.Set("startdateymd", startdate.Format("2006-01-02"))
	query.Set("enddateymd", enddate.Format("2006-01-02"))
	query.Set("offset", strconv.Itoa(offset))
	var dfStr string
	switch len(dataFields) {
	case 0:
		return MeasureGetActivityResponse{}, errors.New("dataFields must contain at least 1")
	case 1:
		dfStr = string(dataFields[0])
	default:
		n := len(dataFields) - 1
		for i := 0; i < len(dataFields); i++ {
			n += len(dataFields[i])
		}

		var b strings.Builder
		b.Grow(n)
		b.WriteString(string(dataFields[0]))
		for _, s := range dataFields[1:] {
			b.WriteString(",")
			b.WriteString(string(s))
		}
		dfStr = b.String()
	}
	query.Set("data_fields", dfStr)
	query.Set("lastupdate", strconv.FormatInt(lastupdate.Unix(), 10))

	var resp MeasureGetActivityResponse
	if err := svc.client.get(ctx, "/v2/measure", query, &resp); err != nil {
		return MeasureGetActivityResponse{}, err
	}
	return resp, nil
}

type MeasureGetActivityResponse struct {
	Status int                            `json:"status"`
	Body   MeasureGetActivityResponseBody `json:"body"`
}

type MeasureGetActivityResponseBody struct {
	Activities []Activity `json:"activities"`
	More       bool       `json:"more"`
	Offset     int        `json:"offset"`
}

type Activity struct {
	Date          string  `json:"date"`
	Timezone      string  `json:"timezone"`
	DeviceID      string  `json:"deviceid"`
	Brand         int     `json:"brand"`
	IsTracker     bool    `json:"is_tracker"`
	Steps         int     `json:"steps"`
	Distance      float64 `json:"distance"`
	Elevation     float64 `json:"elevation"`
	Soft          int     `json:"soft"`
	Moderate      int     `json:"moderate"`
	Intense       int     `json:"intense"`
	Active        int     `json:"active"`
	Calories      float64 `json:"calories"`
	TotalCalories float64 `json:"totalcalories"`
	HRAverage     int     `json:"hr_average"`
	HRMin         int     `json:"hr_min"`
	HRMax         int     `json:"hr_max"`
	HRZone0       int     `json:"hr_zone_0"`
	HRZone1       int     `json:"hr_zone_1"`
	HRZone2       int     `json:"hr_zone_2"`
	HRZone3       int     `json:"hr_zone_3"`
}

type ActivityDataField string

const (
	Steps         ActivityDataField = "steps"
	Distance      ActivityDataField = "distance"
	Elevation     ActivityDataField = "elevation"
	Soft          ActivityDataField = "soft"
	Moderate      ActivityDataField = "moderate"
	Intense       ActivityDataField = "intense"
	Active        ActivityDataField = "active"
	Calories      ActivityDataField = "calories"
	TotalCalories ActivityDataField = "totalcalories"
	HRAverage     ActivityDataField = "hr_average"
	HRMin         ActivityDataField = "hr_min"
	HRMax         ActivityDataField = "hr_max"
	HRZone0       ActivityDataField = "hr_zone_0"
	HRZone1       ActivityDataField = "hr_zone_1"
	HRZone2       ActivityDataField = "hr_zone_2"
	HRZone3       ActivityDataField = "hr_zone_3"
)
