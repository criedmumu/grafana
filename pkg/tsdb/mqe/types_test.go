package mqe

import (
	"testing"

	"time"

	"fmt"

	"github.com/grafana/grafana/pkg/tsdb"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWildcardExpansion(t *testing.T) {
	availableMetrics := []string{
		"os.cpu.all.idle",
		"os.cpu.1.idle",
		"os.cpu.2.idle",
		"os.cpu.3.idle",
	}

	now := time.Now()
	from := now.Add((time.Minute*5)*-1).UnixNano() / int64(time.Millisecond)
	to := now.UnixNano() / int64(time.Millisecond)

	Convey("Can expanding query", t, func() {

		Convey("Without wildcard series", func() {
			query := &MQEQuery{
				Metrics: []MQEMetric{
					MQEMetric{
						Metric: "os.cpu.3.idle",
						Alias:  "cpu on core 3",
					},
					MQEMetric{
						Metric: "os.cpu.2.idle",
						Alias:  "cpu on core 2",
					},
				},
				Hosts:          []string{"staples-lab-1", "staples-lab-2"},
				Apps:           []string{"demoapp-1", "demoapp-2"},
				AddAppToAlias:  false,
				AddHostToAlias: false,
				TimeRange:      &tsdb.TimeRange{Now: now, From: "5m", To: "now"},
			}

			expandeQueries, err := query.Build(availableMetrics)
			So(err, ShouldBeNil)
			So(len(expandeQueries), ShouldEqual, 2)
			So(expandeQueries[0], ShouldEqual, fmt.Sprintf("`os.cpu.3.idle` where app in ('demoapp-1', 'demoapp-2') and host in ('staples-lab-1', 'staples-lab-2') from %v to %v", from, to))
			So(expandeQueries[1], ShouldEqual, fmt.Sprintf("`os.cpu.2.idle` where app in ('demoapp-1', 'demoapp-2') and host in ('staples-lab-1', 'staples-lab-2') from %v to %v", from, to))
		})

		Convey("Containg wildcard series", func() {
			query := &MQEQuery{
				Metrics: []MQEMetric{
					MQEMetric{
						Metric: "os.cpu*",
						Alias:  "cpu on core *",
					},
				},
				Hosts:          []string{"staples-lab-1"},
				AddAppToAlias:  false,
				AddHostToAlias: false,
				TimeRange:      &tsdb.TimeRange{Now: now, From: "5m", To: "now"},
			}

			expandeQueries, err := query.Build(availableMetrics)
			So(err, ShouldBeNil)
			So(len(expandeQueries), ShouldEqual, 4)

			So(expandeQueries[0], ShouldEqual, fmt.Sprintf("`os.cpu.all.idle` where host in ('staples-lab-1') from %v to %v", from, to))
			So(expandeQueries[1], ShouldEqual, fmt.Sprintf("`os.cpu.1.idle` where host in ('staples-lab-1') from %v to %v", from, to))
			So(expandeQueries[2], ShouldEqual, fmt.Sprintf("`os.cpu.2.idle` where host in ('staples-lab-1') from %v to %v", from, to))
			So(expandeQueries[3], ShouldEqual, fmt.Sprintf("`os.cpu.3.idle` where host in ('staples-lab-1') from %v to %v", from, to))

		})
	})
}
