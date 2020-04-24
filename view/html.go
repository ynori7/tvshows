package view

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ynori7/tvshows/streamer"
	"github.com/ynori7/tvshows/tvshow"
	"html/template"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strings"
)

type HtmlTemplate struct {
	NewTvShows []tvshow.TvShow
	ReturningTvShows []tvshow.TvShow
}

func NewHtmlTemplate(newTvShows []tvshow.TvShow, returningTvShows []tvshow.TvShow) HtmlTemplate {
	return HtmlTemplate{
		NewTvShows: newTvShows,
		ReturningTvShows: returningTvShows,
	}
}

func (h HtmlTemplate) ExecuteHtmlTemplate() (string, error) {
	t := template.Must(template.New("html").
		Funcs(template.FuncMap{
			"mod": func(i, j int) bool { return i%j == 0 },
			"getStreamer": func(s streamer.Streamer) string {
				if s == "" {
					return ""
				}
				return fmt.Sprintf("Available on %s", s)
			},
			"genres": func(genres []string) string {
				return strings.Join(genres, ", ")
			},
			"formatNumber": func(num int) string {
				p := message.NewPrinter(language.English)
				return p.Sprintf("%d", num)
			},
		}).
		Parse(htmlTemplate))

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	err := t.Execute(w, h)
	if err != nil {
		return "", err
	}

	w.Flush()
	return b.String(), nil
}

const htmlTemplate = `<html>
<head>
	<style>
		.premieresList {
			border-bottom: 1px;
    		border-bottom-style: solid;
    		border-color: #e3e3e3;
    		border-spacing: 0 15px;
		}
        .ratings_wrapper {
            height:50px;
        }
        .ratings_wrapper .imdbRating {
            background: url(https://m.media-amazon.com/images/G/01/imdb/images/title/title_overview_sprite-1705639977._V_.png) no-repeat;
            background-position: -15px -118px;
            float: left;
            font-size: 11px;
            height: 30px;
            line-height: 13px;
            padding: 5px 0 0 34px;
            width: 100%;
        }
        .ratings_wrapper .imdbRating .ratingValue strong {
            font-size: 18px;
            font-weight: normal;
            font-family: Arial;
            line-height: 18px;
        }
        .grey {
            color: #6b6b6b;
            font-size: 10px;
        }
        .small {
            font-size: 10px;
        }
        .left-part {
            float:left;
            width: 44%;
            margin-right: 1%;
        }
        .left-part img {
            max-width: 100%;
        }
        .right-part {
            float:right;
            width: 54%;
            margin-left: 1%;
        }
        .genres {
            margin-top:10px;
        }
        .title {
            font: 18px Arial,sans-serif;
            font-weight: normal;
            line-height: 110%;
            margin: 0px;
            padding-bottom: 10px;
        }
        .title a,.title a:visited {
            text-decoration: none;
            color: #333;
        }
        .description {
            padding: 10px 10px 10px 0;
            font-size: 12px;
            text-align: justify;
        }
        .score .scoreValue {
            font-size: 14px;
            color: #333;
        }
		tr {
			margin-bottom: 15px;
		}
		.ratings_wrapper, .score, .genres {
			width: 100%;
		}
	</style>
</head>
<body>
<h1>Returning Series</h1>
<table width="640" cellpadding="0" cellspacing="0" border="0" bgcolor="#FFFFFF">
    <tbody><tr>
        <td height="10" style="font-size:10px;line-height:10px">&nbsp;</td>
    </tr>
    <tr>
        <td align="center" valign="top">
            <table width="600" cellpadding="0" cellspacing="0" border="0" class="premieresList">
                <tbody>
				{{range $i, $val := .ReturningTvShows}}
					{{ if eq $i 0 }}<tr>{{ else if mod $i 2 }}</tr><tr>{{ else }}<td width="2%" align="center" valign="top">&nbsp;</td>{{ end }}
					<td width="49%" align="left" valign="top">
    			        <div class="title"><a href="{{ $val.Link }}">{{ $val.Title }}</a></div>
						<div class="left-part">
							<img alt="{{ $val.Title }} Poster" title="{{ $val.Title }} Poster" src="{{ $val.Image }}">
            			</div>
						<div class="right-part">
							<div class="ratings_wrapper">
								<div class="imdbRating">
									<div class="ratingValue">
										<strong title="{{ $val.Rating.AverageRating }} based on {{ formatNumber $val.Rating.RatingCount }} user ratings"><span>{{ $val.Rating.AverageRating }}</span></strong><span class="grey">/</span><span class="grey" itemprop="bestRating">10</span>
									</div>
									<span class="small" itemprop="ratingCount">{{ formatNumber $val.Rating.RatingCount }}</span>
								</div>
							</div>
							<div class="score small grey"><span class="scoreValue">{{ $val.Score }}</span>/100</div>
							<div class="genres small grey">{{ genres $val.Genres }}</div>
						</div>
						<div style="clear:both">
							<div class="description">
								<span>{{ $val.Description }}</span>
					    	</div>
            				<div class="streamer">
                				<span class="small" style="font-weight:bold">{{ getStreamer $val.StreamingOption }}</span>
            				</div>
					    </div>
					</td>
				{{ end }}
              </tr>
            </tbody></table>
        </td>
    </tr>
    <tr>
        <td height="10" style="font-size:10px;line-height:10px">&nbsp;</td>
    </tr>
</tbody></table>

<h1>New Series</h1>
<table width="640" cellpadding="0" cellspacing="0" border="0" bgcolor="#FFFFFF">
    <tbody><tr>
        <td height="10" style="font-size:10px;line-height:10px">&nbsp;</td>
    </tr>
    <tr>
        <td align="center" valign="top">
            <table width="600" cellpadding="0" cellspacing="0" border="0" class="premieresList">
                <tbody>
				{{range $i, $val := .NewTvShows}}
					{{ if eq $i 0 }}<tr>{{ else if mod $i 2 }}</tr><tr>{{ else }}<td width="2%" align="center" valign="top">&nbsp;</td>{{ end }}
					<td width="32%" align="left" valign="top">
    			        <div class="title"><a href="{{ $val.Link }}">{{ $val.Title }}</a></div>
						<div class="left-part">
							<img alt="{{ $val.Title }} Poster" title="{{ $val.Title }} Poster" src="{{ $val.Image }}">
            			</div>
						<div class="right-part">
							<div class="ratings_wrapper">
								<div class="imdbRating">
									<div class="ratingValue">
										<strong title="{{ $val.Rating.AverageRating }} based on {{ $val.Rating.RatingCount }} user ratings"><span>{{ $val.Rating.AverageRating }}</span></strong><span class="grey">/</span><span class="grey" itemprop="bestRating">10</span>
									</div>
									<span class="small" itemprop="ratingCount">{{ $val.Rating.RatingCount }}</span>
								</div>
							</div>
							<div class="score small grey"><span class="scoreValue">{{ $val.Score }}</span>/100</div>
							<div class="genres small grey">{{ genres $val.Genres }}</div>
						</div>
						<div style="clear:both">
							<div class="description">
								<span>{{ $val.Description }}</span>
					    	</div>
            				<div class="streamer">
                				<span class="small" style="font-weight:bold">{{ getStreamer $val.StreamingOption }}</span>
            				</div>
					    </div>
					</td>
				{{ end }}
              </tr>
            </tbody></table>
        </td>
    </tr>
    <tr>
        <td height="10" style="font-size:10px;line-height:10px">&nbsp;</td>
    </tr>
</tbody></table>
</body>
</html>
`
