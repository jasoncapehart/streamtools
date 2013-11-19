package blocks

import (
	"encoding/json"
	"log"
)

var (
	Library     BlockLibrary
	LibraryBlob string
)

// A block library is a collection of possible block templates
type BlockLibrary map[string]*BlockTemplate

func BuildLibrary() {
	log.Println("building block library")
	Library = make(map[string]*BlockTemplate)

	templates := []*BlockTemplate{
		&BlockTemplate{
			BlockType:  "connection",
			RouteNames: []string{"last_message", "rate"},
			Routine:    Connection,
		},
		////////// TESTING BLOCKS
		&BlockTemplate{
			BlockType:  "blocked",
			RouteNames: []string{"get_rule"},
			Routine:    Blocked,
		},
		////////////////////
		&BlockTemplate{
			BlockType:  "ticker",
			RouteNames: []string{"set_rule", "get_rule"},
			Routine:    Ticker,
		},
		&BlockTemplate{
			BlockType:  "tolog",
			RouteNames: []string{},
			Routine:    ToLog,
		},
		&BlockTemplate{
			BlockType:  "random",
			RouteNames: []string{"set_rule", "get_rule"},
			Routine:    Random,
		},
		&BlockTemplate{
			BlockType:  "count",
			RouteNames: []string{"set_rule", "get_rule", "count"},
			Routine:    Count,
		},
		&BlockTemplate{
			BlockType:  "mask",
			RouteNames: []string{"set_rule", "get_rule"},
			Routine:    Mask,
		},
		&BlockTemplate{
			BlockType:  "sync",
			RouteNames: []string{"set_rule", "get_rule"},
			Routine:    Sync,
		},
		&BlockTemplate{
			BlockType:  "postValue",
			RouteNames: []string{"set_rule", "get_rule"},
			Routine:    PostValue,
		},
		&BlockTemplate{
			BlockType:  "post",
			RouteNames: []string{"set_rule", "get_rule"},
			Routine:    Post,
		},
		&BlockTemplate{
			BlockType:  "date",
			RouteNames: []string{"set_rule", "get_rule"},
			Routine:    Date,
		},
		&BlockTemplate{
			BlockType:  "fromNSQ",
			RouteNames: []string{"set_rule", "get_rule"},
			Routine:    FromNSQ,
		},
		&BlockTemplate{
			BlockType:  "pollS3",
			RouteNames: []string{"set_rule", "poll_now"},
			Routine:    PollS3,
		},
		&BlockTemplate{
			BlockType:  "tofile",
			RouteNames: []string{"set_rule", "get_rule"},
			Routine:    ToFile,
		},
		&BlockTemplate{
			BlockType:  "toNSQ",
			RouteNames: []string{"set_rule", "get_rule"},
			Routine:    ToNSQ,
		},
		&BlockTemplate{
			BlockType:  "filter",
			RouteNames: []string{"set_rule", "get_rule"},
			Routine:    Filter,
		},
		&BlockTemplate{
			BlockType:  "postto",
			RouteNames: []string{"in"},
			Routine:    PostTo,
		},
		&BlockTemplate{
			BlockType:  "timeseries",
			RouteNames: []string{"set_rule", "timeseries", "get_rule"},
			Routine:    Timeseries,
		},
		&BlockTemplate{
			BlockType:  "histogram",
			RouteNames: []string{"set_rule", "histogram", "get_rule"},
			Routine:    Histogram,
		},
		&BlockTemplate{
			BlockType:  "groupHistogram",
			RouteNames: []string{"set_rule", "histogram", "get_rule", "list"},
			Routine:    GroupHistogram,
		},
		&BlockTemplate{
			BlockType:  "bunch",
			RouteNames: []string{"set_rule", "get_rule"},
			Routine:    Bunch,
		},
        &BlockTemplate{
            BlockType:  "avg",
            RouteNames: []string{"set_rule", "get_rule", "avg"},
            Routine:    Avg,
        },
        &BlockTemplate{
            BlockType:  "sd",
            RouteNames: []string{"set_rule", "get_rule", "sd"},
            Routine:    Sd,
        },
        &BlockTemplate{
            BlockType:  "var",
            RouteNames: []string{"set_rule", "get_rule", "var"},
            Routine:    Var,
        },
        &BlockTemplate{
			BlockType:  "longHTTP",
			RouteNames: []string{"set_rule", "get_rule"},
			Routine:    LongHTTP,
		},
        &BlockTemplate{
            BlockType: "scale_range",
            RouteNames: []string{"set_rule", "get_rule", "x_scaled"},
            Routine:    ScaleRange,
        },
	}

	libraryList := []map[string]interface{}{}
	for _, t := range templates {
		blockItem := make(map[string]interface{})
		blockItem["BlockType"] = t.BlockType
		blockItem["RouteNames"] = t.RouteNames
		libraryList = append(libraryList, blockItem)

		Library[t.BlockType] = t
	}

	blob, _ := json.Marshal(libraryList)
	LibraryBlob = string(blob)
}
