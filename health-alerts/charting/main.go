package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type CaseResponse struct {
	Payload    []Payload `json:"payload"`
	StatusCode int       `json:"status_code"`
}
type Payload struct {
	ID       int    `json:"id"`
	Date     string `json:"date"`
	Reported int    `json:"reported"`
	Total    int    `json:"total"`
}

type PlotlyConfig struct {
	Data []struct {
		Type string   `json:"type"`
		X    []string `json:"x"`
		Y    []int    `json:"y"`
	} `json:"data"`
	Layout struct {
		Height   int `json:"height"`
		Template struct {
			Data struct {
				Bar []struct {
					ErrorX struct {
						Color string `json:"color"`
					} `json:"error_x"`
					ErrorY struct {
						Color string `json:"color"`
					} `json:"error_y"`
					Marker struct {
						Line struct {
							Color string  `json:"color"`
							Width float64 `json:"width"`
						} `json:"line"`
					} `json:"marker"`
					Type string `json:"type"`
				} `json:"bar"`
				Barpolar []struct {
					Marker struct {
						Line struct {
							Color string  `json:"color"`
							Width float64 `json:"width"`
						} `json:"line"`
					} `json:"marker"`
					Type string `json:"type"`
				} `json:"barpolar"`
				Carpet []struct {
					Aaxis struct {
						Endlinecolor   string `json:"endlinecolor"`
						Gridcolor      string `json:"gridcolor"`
						Linecolor      string `json:"linecolor"`
						Minorgridcolor string `json:"minorgridcolor"`
						Startlinecolor string `json:"startlinecolor"`
					} `json:"aaxis"`
					Baxis struct {
						Endlinecolor   string `json:"endlinecolor"`
						Gridcolor      string `json:"gridcolor"`
						Linecolor      string `json:"linecolor"`
						Minorgridcolor string `json:"minorgridcolor"`
						Startlinecolor string `json:"startlinecolor"`
					} `json:"baxis"`
					Type string `json:"type"`
				} `json:"carpet"`
				Choropleth []struct {
					Colorbar struct {
						Outlinewidth int    `json:"outlinewidth"`
						Ticks        string `json:"ticks"`
					} `json:"colorbar"`
					Type string `json:"type"`
				} `json:"choropleth"`
				Contour []struct {
					Colorbar struct {
						Outlinewidth int    `json:"outlinewidth"`
						Ticks        string `json:"ticks"`
					} `json:"colorbar"`
					Colorscale [][]interface{} `json:"colorscale"`
					Type       string          `json:"type"`
				} `json:"contour"`
				Contourcarpet []struct {
					Colorbar struct {
						Outlinewidth int    `json:"outlinewidth"`
						Ticks        string `json:"ticks"`
					} `json:"colorbar"`
					Type string `json:"type"`
				} `json:"contourcarpet"`
				Heatmap []struct {
					Colorbar struct {
						Outlinewidth int    `json:"outlinewidth"`
						Ticks        string `json:"ticks"`
					} `json:"colorbar"`
					Colorscale [][]interface{} `json:"colorscale"`
					Type       string          `json:"type"`
				} `json:"heatmap"`
				Heatmapgl []struct {
					Colorbar struct {
						Outlinewidth int    `json:"outlinewidth"`
						Ticks        string `json:"ticks"`
					} `json:"colorbar"`
					Colorscale [][]interface{} `json:"colorscale"`
					Type       string          `json:"type"`
				} `json:"heatmapgl"`
				Histogram []struct {
					Marker struct {
						Colorbar struct {
							Outlinewidth int    `json:"outlinewidth"`
							Ticks        string `json:"ticks"`
						} `json:"colorbar"`
					} `json:"marker"`
					Type string `json:"type"`
				} `json:"histogram"`
				Histogram2D []struct {
					Colorbar struct {
						Outlinewidth int    `json:"outlinewidth"`
						Ticks        string `json:"ticks"`
					} `json:"colorbar"`
					Colorscale [][]interface{} `json:"colorscale"`
					Type       string          `json:"type"`
				} `json:"histogram2d"`
				Histogram2Dcontour []struct {
					Colorbar struct {
						Outlinewidth int    `json:"outlinewidth"`
						Ticks        string `json:"ticks"`
					} `json:"colorbar"`
					Colorscale [][]interface{} `json:"colorscale"`
					Type       string          `json:"type"`
				} `json:"histogram2dcontour"`
				Mesh3D []struct {
					Colorbar struct {
						Outlinewidth int    `json:"outlinewidth"`
						Ticks        string `json:"ticks"`
					} `json:"colorbar"`
					Type string `json:"type"`
				} `json:"mesh3d"`
				Parcoords []struct {
					Line struct {
						Colorbar struct {
							Outlinewidth int    `json:"outlinewidth"`
							Ticks        string `json:"ticks"`
						} `json:"colorbar"`
					} `json:"line"`
					Type string `json:"type"`
				} `json:"parcoords"`
				Pie []struct {
					Automargin bool   `json:"automargin"`
					Type       string `json:"type"`
				} `json:"pie"`
				Scatter []struct {
					Marker struct {
						Line struct {
							Color string `json:"color"`
						} `json:"line"`
					} `json:"marker"`
					Type string `json:"type"`
				} `json:"scatter"`
				Scatter3D []struct {
					Line struct {
						Colorbar struct {
							Outlinewidth int    `json:"outlinewidth"`
							Ticks        string `json:"ticks"`
						} `json:"colorbar"`
					} `json:"line"`
					Marker struct {
						Colorbar struct {
							Outlinewidth int    `json:"outlinewidth"`
							Ticks        string `json:"ticks"`
						} `json:"colorbar"`
					} `json:"marker"`
					Type string `json:"type"`
				} `json:"scatter3d"`
				Scattercarpet []struct {
					Marker struct {
						Colorbar struct {
							Outlinewidth int    `json:"outlinewidth"`
							Ticks        string `json:"ticks"`
						} `json:"colorbar"`
					} `json:"marker"`
					Type string `json:"type"`
				} `json:"scattercarpet"`
				Scattergeo []struct {
					Marker struct {
						Colorbar struct {
							Outlinewidth int    `json:"outlinewidth"`
							Ticks        string `json:"ticks"`
						} `json:"colorbar"`
					} `json:"marker"`
					Type string `json:"type"`
				} `json:"scattergeo"`
				Scattergl []struct {
					Marker struct {
						Line struct {
							Color string `json:"color"`
						} `json:"line"`
					} `json:"marker"`
					Type string `json:"type"`
				} `json:"scattergl"`
				Scattermapbox []struct {
					Marker struct {
						Colorbar struct {
							Outlinewidth int    `json:"outlinewidth"`
							Ticks        string `json:"ticks"`
						} `json:"colorbar"`
					} `json:"marker"`
					Type string `json:"type"`
				} `json:"scattermapbox"`
				Scatterpolar []struct {
					Marker struct {
						Colorbar struct {
							Outlinewidth int    `json:"outlinewidth"`
							Ticks        string `json:"ticks"`
						} `json:"colorbar"`
					} `json:"marker"`
					Type string `json:"type"`
				} `json:"scatterpolar"`
				Scatterpolargl []struct {
					Marker struct {
						Colorbar struct {
							Outlinewidth int    `json:"outlinewidth"`
							Ticks        string `json:"ticks"`
						} `json:"colorbar"`
					} `json:"marker"`
					Type string `json:"type"`
				} `json:"scatterpolargl"`
				Scatterternary []struct {
					Marker struct {
						Colorbar struct {
							Outlinewidth int    `json:"outlinewidth"`
							Ticks        string `json:"ticks"`
						} `json:"colorbar"`
					} `json:"marker"`
					Type string `json:"type"`
				} `json:"scatterternary"`
				Surface []struct {
					Colorbar struct {
						Outlinewidth int    `json:"outlinewidth"`
						Ticks        string `json:"ticks"`
					} `json:"colorbar"`
					Colorscale [][]interface{} `json:"colorscale"`
					Type       string          `json:"type"`
				} `json:"surface"`
				Table []struct {
					Cells struct {
						Fill struct {
							Color string `json:"color"`
						} `json:"fill"`
						Line struct {
							Color string `json:"color"`
						} `json:"line"`
					} `json:"cells"`
					Header struct {
						Fill struct {
							Color string `json:"color"`
						} `json:"fill"`
						Line struct {
							Color string `json:"color"`
						} `json:"line"`
					} `json:"header"`
					Type string `json:"type"`
				} `json:"table"`
			} `json:"data"`
			Layout struct {
				Annotationdefaults struct {
					Arrowcolor string `json:"arrowcolor"`
					Arrowhead  int    `json:"arrowhead"`
					Arrowwidth int    `json:"arrowwidth"`
				} `json:"annotationdefaults"`
				Coloraxis struct {
					Colorbar struct {
						Outlinewidth int    `json:"outlinewidth"`
						Ticks        string `json:"ticks"`
					} `json:"colorbar"`
				} `json:"coloraxis"`
				Colorscale struct {
					Diverging       [][]interface{} `json:"diverging"`
					Sequential      [][]interface{} `json:"sequential"`
					Sequentialminus [][]interface{} `json:"sequentialminus"`
				} `json:"colorscale"`
				Colorway []string `json:"colorway"`
				Font     struct {
					Color string `json:"color"`
				} `json:"font"`
				Geo struct {
					Bgcolor      string `json:"bgcolor"`
					Lakecolor    string `json:"lakecolor"`
					Landcolor    string `json:"landcolor"`
					Showlakes    bool   `json:"showlakes"`
					Showland     bool   `json:"showland"`
					Subunitcolor string `json:"subunitcolor"`
				} `json:"geo"`
				Hoverlabel struct {
					Align string `json:"align"`
				} `json:"hoverlabel"`
				Hovermode string `json:"hovermode"`
				Mapbox    struct {
					Style string `json:"style"`
				} `json:"mapbox"`
				PaperBgcolor string `json:"paper_bgcolor"`
				PlotBgcolor  string `json:"plot_bgcolor"`
				Polar        struct {
					Angularaxis struct {
						Gridcolor string `json:"gridcolor"`
						Linecolor string `json:"linecolor"`
						Ticks     string `json:"ticks"`
					} `json:"angularaxis"`
					Bgcolor    string `json:"bgcolor"`
					Radialaxis struct {
						Gridcolor string `json:"gridcolor"`
						Linecolor string `json:"linecolor"`
						Ticks     string `json:"ticks"`
					} `json:"radialaxis"`
				} `json:"polar"`
				Scene struct {
					Xaxis struct {
						Backgroundcolor string `json:"backgroundcolor"`
						Gridcolor       string `json:"gridcolor"`
						Gridwidth       int    `json:"gridwidth"`
						Linecolor       string `json:"linecolor"`
						Showbackground  bool   `json:"showbackground"`
						Ticks           string `json:"ticks"`
						Zerolinecolor   string `json:"zerolinecolor"`
					} `json:"xaxis"`
					Yaxis struct {
						Backgroundcolor string `json:"backgroundcolor"`
						Gridcolor       string `json:"gridcolor"`
						Gridwidth       int    `json:"gridwidth"`
						Linecolor       string `json:"linecolor"`
						Showbackground  bool   `json:"showbackground"`
						Ticks           string `json:"ticks"`
						Zerolinecolor   string `json:"zerolinecolor"`
					} `json:"yaxis"`
					Zaxis struct {
						Backgroundcolor string `json:"backgroundcolor"`
						Gridcolor       string `json:"gridcolor"`
						Gridwidth       int    `json:"gridwidth"`
						Linecolor       string `json:"linecolor"`
						Showbackground  bool   `json:"showbackground"`
						Ticks           string `json:"ticks"`
						Zerolinecolor   string `json:"zerolinecolor"`
					} `json:"zaxis"`
				} `json:"scene"`
				Shapedefaults struct {
					Line struct {
						Color string `json:"color"`
					} `json:"line"`
				} `json:"shapedefaults"`
				Sliderdefaults struct {
					Bgcolor     string `json:"bgcolor"`
					Bordercolor string `json:"bordercolor"`
					Borderwidth int    `json:"borderwidth"`
					Tickwidth   int    `json:"tickwidth"`
				} `json:"sliderdefaults"`
				Ternary struct {
					Aaxis struct {
						Gridcolor string `json:"gridcolor"`
						Linecolor string `json:"linecolor"`
						Ticks     string `json:"ticks"`
					} `json:"aaxis"`
					Baxis struct {
						Gridcolor string `json:"gridcolor"`
						Linecolor string `json:"linecolor"`
						Ticks     string `json:"ticks"`
					} `json:"baxis"`
					Bgcolor string `json:"bgcolor"`
					Caxis   struct {
						Gridcolor string `json:"gridcolor"`
						Linecolor string `json:"linecolor"`
						Ticks     string `json:"ticks"`
					} `json:"caxis"`
				} `json:"ternary"`
				Title struct {
					X float64 `json:"x"`
				} `json:"title"`
				Updatemenudefaults struct {
					Bgcolor     string `json:"bgcolor"`
					Borderwidth int    `json:"borderwidth"`
				} `json:"updatemenudefaults"`
				Xaxis struct {
					Automargin bool   `json:"automargin"`
					Gridcolor  string `json:"gridcolor"`
					Linecolor  string `json:"linecolor"`
					Ticks      string `json:"ticks"`
					Title      struct {
						Standoff int `json:"standoff"`
					} `json:"title"`
					Zerolinecolor string `json:"zerolinecolor"`
					Zerolinewidth int    `json:"zerolinewidth"`
				} `json:"xaxis"`
				Yaxis struct {
					Automargin bool   `json:"automargin"`
					Gridcolor  string `json:"gridcolor"`
					Linecolor  string `json:"linecolor"`
					Ticks      string `json:"ticks"`
					Title      struct {
						Standoff int `json:"standoff"`
					} `json:"title"`
					Zerolinecolor string `json:"zerolinecolor"`
					Zerolinewidth int    `json:"zerolinewidth"`
				} `json:"yaxis"`
			} `json:"layout"`
		} `json:"template"`
		Title struct {
			Text string `json:"text"`
		} `json:"title"`
		Width int `json:"width"`
	} `json:"layout"`
}

func main() {
	res, err := http.Get("https://api.aditya.diwakar.io/gt-jpj/cases")
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	casesWrapped := CaseResponse{}
	json.NewDecoder(res.Body).Decode(&casesWrapped)

	json.NewEncoder(os.Stdout).Encode(casesWrapped.Payload)
	dates := []string{}
	reported := []int{}

	for _, data := range casesWrapped.Payload {
		reported = append(reported, data.Reported)

		t, _ := time.Parse("January 2, 2006", data.Date)
		dates = append(dates, t.Format("2006-01-02"))
	}

	config := PlotlyConfig{}
	json.Unmarshal([]byte(jsonTemplate), &config)

	config.Data[0].X = dates
	config.Data[0].Y = reported

	json.NewEncoder(os.Stdout).Encode(config)
}
