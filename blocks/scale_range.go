package blocks

import (
    "encoding/json"
    "log"
)

// ScaleRange() rescales data in an online fashion
// All data will be mapped to the specified range, but the
// position each point maps to may change as the empirical
// min and max values for the stream are updated.
// Formula: x_scaled = (x - min / (max - min)) * (scaled_max - scaled_min) + scaled_min 
func ScaleRange(b *Block) {

    type scaleRule struct {
        Key string
        Min string
        Max string
    }

    type scaleData struct {
        X_scaled float64
    }

    data := &avgData{X_scaled: 0.0}
    var rule *scaleRule

    min := 0.0
    max := 0.0

    for {
        select {
        case query := <-b.Routes["scale_range"]:
            marshal(query, data)
        case ruleUpdate := <-b.Routes["set_rule"]:
            if rule == nil {
                rule = &scaleRule{}
            }
            unmarshal(ruleUpdate, rule)
        case msg := <-b.Routes["get_rule"]:
            if rule == nil {
                marshal(msg, &scaleRule{})
            } else {
                marshal(msg, rule)
            }
        case <-b.QuitChan:
            quit(b)
            return
        case msg := <-b.InChan:
            if rule == nil {
                break
            }

            x_val := getKeyValues(msg, rule.Key)[0].(json.Number)
            x, err := x_val.Float64()
            if err != nil {
                log.Println(err.Error())
            }
            if x < min {
                min = x
            }
            if x > max {
                max = x
            }

            max_val := getKeyValues(msg, rule.Max)[0].(json.Number)
            scaled_max, err_max := max_val.Float64()
            if err_max != nil {
                log.Println(err.Error())
            }

            min_val := getKeyValues(msg, rule.Min)[0].(json.Number)
            scaled_min, err_min := min_val.Float64()
            if err_min != nil {
                log.Println(err.Error())
            }

            data.X_prime = ((x - min) / (max - min)) * (scaled_max - scaled_min) + scaled_min
        }
    }
}
