package withings_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	withings "github.com/nasa9084/go-withings"
	"golang.org/x/oauth2"
)

var tokenSrc = oauth2.StaticTokenSource(
	&oauth2.Token{
		AccessToken: "this_is_oauth2_token",
	},
)

var httpClient = oauth2.NewClient(oauth2.NoContext, tokenSrc)

func assertEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected result:\n  got:  %+v\n  want: %+v", got, want)
	}
}

func handler(t *testing.T, response string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token := r.Header.Get("Authorization"); token == "" {
			t.Fatal("Authorization header is empty or undefined")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}

const UserGetDeviceSuccessResponse = `{
  "status": 0,
  "body": {
    "devices": [
      {
        "type": "string",
        "model": "string",
        "battery": "string",
        "deviceid": "string",
        "timezone": "string"
      }
    ]
  }
}`

func TestUserGetDevice(t *testing.T) {
	srv := httptest.NewServer(handler(t, UserGetDeviceSuccessResponse))
	defer srv.Close()

	c := withings.New(withings.WithEndpoint(srv.URL), withings.WithHTTPClient(httpClient))
	got, err := c.User().GetDevice(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := withings.UserGetDeviceResponse{
		Status: 0,
		Body: withings.UserGetDeviceResponseBody{
			Devices: []withings.Device{
				withings.Device{
					Type:     "string",
					Model:    "string",
					Battery:  "string",
					DeviceID: "string",
					Timezone: "string",
				},
			},
		},
	}
	assertEqual(t, got, want)
}

const MeasureGetMeasSuccessResponse = `{
  "status": 0,
  "body": {
    "updatetime": "string",
    "timezone": "string",
    "measuregrps": [
      {
        "grpid": 0,
        "attrib": 0,
        "date": 0,
        "created": 0,
        "category": 0,
        "deviceid": "string",
        "measures": [
          {
            "value": 0,
            "type": 0,
            "unit": 0,
            "algo": 0,
            "fm": 0,
            "fw": 0
          }
        ],
        "comment": "string"
      }
    ],
    "more": true,
    "offset": 0
  }
}`

func TestMeasureGetMeas(t *testing.T) {
	srv := httptest.NewServer(handler(t, MeasureGetMeasSuccessResponse))
	defer srv.Close()

	c := withings.New(withings.WithEndpoint(srv.URL), withings.WithHTTPClient(httpClient))
	got, err := c.Measure().GetMeas(context.Background(), withings.Weight, withings.RealMeasure, time.Now(), time.Now(), 0, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	want := withings.MeasureGetMeasResponse{
		Status: 0,
		Body: withings.MeasureGetMeasResponseBody{
			UpdateTime: "string",
			Timezone:   "string",
			MeasureGroups: []withings.MeasureGroup{
				withings.MeasureGroup{
					GroupID:   0,
					Attribute: 0,
					Date:      time.Unix(0, 0),
					Created:   time.Unix(0, 0),
					Category:  withings.UnknownMeasureCategory,
					DeviceID:  "string",
					Measures: []withings.Measure{
						withings.Measure{
							Value: 0,
							Type:  withings.UnknownMeasureType,
							Unit:  0,
							Algo:  0,
							Fm:    0,
							Fw:    0,
						},
					},
					Comment: "string",
				},
			},
			More:   true,
			Offset: 0,
		},
	}
	assertEqual(t, got, want)
}

const MeasureGetActivitySuccessResponse = `{
  "status": 0,
  "body": {
    "activities": [
      {
        "date": "string",
        "timezone": "string",
        "deviceid": "string",
        "brand": 0,
        "is_tracker": true,
        "steps": 0,
        "distance": 0,
        "elevation": 0,
        "soft": 0,
        "moderate": 0,
        "intense": 0,
        "active": 0,
        "calories": 0,
        "totalcalories": 0,
        "hr_average": 0,
        "hr_min": 0,
        "hr_max": 0,
        "hr_zone_0": 0,
        "hr_zone_1": 0,
        "hr_zone_2": 0,
        "hr_zone_3": 0
      }
    ],
    "more": true,
    "offset": 0
  }
}`

func TestMeasureGetActivity(t *testing.T) {
	srv := httptest.NewServer(handler(t, MeasureGetActivitySuccessResponse))
	defer srv.Close()

	c := withings.New(withings.WithEndpoint(srv.URL), withings.WithHTTPClient(httpClient))
	got, err := c.Measure().GetActivity(context.Background(), time.Now(), time.Now(), 0, []withings.ActivityDataField{withings.Steps}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	want := withings.MeasureGetActivityResponse{
		Status: 0,
		Body: withings.MeasureGetActivityResponseBody{
			Activities: []withings.Activity{
				withings.Activity{
					Date:          "string",
					Timezone:      "string",
					DeviceID:      "string",
					Brand:         0,
					IsTracker:     true,
					Steps:         0,
					Distance:      0,
					Elevation:     0,
					Soft:          0,
					Moderate:      0,
					Intense:       0,
					Active:        0,
					Calories:      0,
					TotalCalories: 0,
					HRAverage:     0,
					HRMin:         0,
					HRMax:         0,
					HRZone0:       0,
					HRZone1:       0,
					HRZone2:       0,
					HRZone3:       0,
				},
			},
			More:   true,
			Offset: 0,
		},
	}
	assertEqual(t, got, want)
}
