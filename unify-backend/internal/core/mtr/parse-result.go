package mtr

import "encoding/json"

func ParseResult(data []byte) (*Result, error) {
	var raw MtrJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	result := &Result{
		Host:      raw.Report.MTR.Dst,
		TotalHops: len(raw.Report.Hubs),
		Reachable: len(raw.Report.Hubs) > 0,
	}

	var totalRTT float64
	var maxLoss float64

	for i, h := range raw.Report.Hubs {
		result.Hops = append(result.Hops, HopResult{
			Hop:    i + 1,
			IP:     h.Host,
			Loss:   h.Loss,
			AvgRTT: h.Avg,
		})

		totalRTT += h.Avg
		if h.Loss > maxLoss {
			maxLoss = h.Loss
		}
	}

	if len(raw.Report.Hubs) > 0 {
		result.AvgRTT = totalRTT / float64(len(raw.Report.Hubs))
	}

	result.MaxLoss = maxLoss
	return result, nil
}
