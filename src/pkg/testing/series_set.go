package testing

import "github.com/prometheus/prometheus/storage"

type Series struct {
	Labels map[string]string
	Points []Point
}

type Point struct {
	Time  int64
	Value float64
}

func ExplodeSeriesSet(seriesSet storage.SeriesSet) []Series {
	series := []Series{}

	for {
		hasSeries := seriesSet.Next()
		if !hasSeries {
			break
		}

		currentSeries := seriesSet.At()

		newSeries := Series{Labels: make(map[string]string)}
		for _, label := range currentSeries.Labels() {
			newSeries.Labels[label.Name] = label.Value
		}

		iterator := currentSeries.Iterator()
		for {
			hasPoint := iterator.Next()
			if !hasPoint {
				break
			}

			timestamp, value := iterator.At()
			newSeries.Points = append(newSeries.Points, Point{Time: timestamp, Value: value})
		}

		series = append(series, newSeries)

	}

	return series

}
