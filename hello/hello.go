package hello

import (
    "fmt"
    "http"
    "appengine"
    "appengine/datastore"
    "appengine/urlfetch"
    "log"
    "io/ioutil"
    "bytes"
    "url"
    "strings"
    "template"
    "time"
)

type FreeNumbers struct {
    Date		datastore.Time
    Gynaecologists	int // акушер-гинекол.
    Venereologist	int // дерма-венеролог
    Cardiologist	int // кардиолог
    Neurologist 	int // невропатолог
    Otolaryngologist	int // отоларинголог
    Ophthalmologist	int // офтальмолог
    Proctologist	int // проктолог
    Rheumatologist	int // ревматолог
    Physician		int // терапевт участ.
    Urologist		int // уролог
    Surgeon		int // хирург
    Endocrinologist	int // эндокринолог
}

func init() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/cron", update)
}

func handler(w http.ResponseWriter, r *http.Request) {
    data := get_olddata(w, r)
    if err := IndexTemplate.Execute(w, data); err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
    }
}

func update(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    a := FreeNumbers{Date: datastore.SecondsToTime(time.Seconds())}

    var s = post(r, "http://nikmed.spb.ru/cgi-bin/tcgi1.exe", "COMMAND=2")
    var prefix = "<span style=\"text-align:center;font-size:large;font-family:arial\">"
    var prefix_len = len(prefix)
    var suffix_lef = len("</span>")
    for _, i := range strings.Split(s, "\n") {
        if strings.HasPrefix(i, prefix) {
	    var s = i[prefix_len:len(i)-suffix_lef-1]
	    s = strings.TrimSpace(strings.ToLower(s))
	    fields := strings.Split(s, " ")
	    f_len := len(fields)
	    name := strings.Join(fields[0:f_len-2], " ")
	    count := fields[f_len-1]
	    var count_int int
	    if _, e := fmt.Sscan(count, &count_int); e != nil {
	    	count_int = 0
	    }
	    switch name {
		case "акушер-гинекол.": {
			a.Gynaecologists = count_int
		}
		case "дерма-венеролог": {
			a.Venereologist = count_int
		}
		case "кардиолог": {
			a.Cardiologist = count_int
		}
		case "невропатолог": {
			a.Neurologist = count_int
		}
		case "отоларинголог": {
			a.Otolaryngologist = count_int
		}
		case "офтальмолог": {
			a.Ophthalmologist = count_int
		}
		case "проктолог": {
			a.Proctologist = count_int
		}
		case "ревматолог": {
			a.Rheumatologist = count_int
		}
		case "терапевт участ.": {
			a.Physician = count_int
		}
		case "уролог": {
			a.Urologist = count_int
		}
		case "хирург": {
			a.Surgeon = count_int
		}
		case "эндокринолог": {
			a.Endocrinologist = count_int
		}
	    }
            // fmt.Fprint(w, name, " -> ", count_int, "\n")
	}
    }
    _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "FreeNumbers", nil), &a)
    if err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
        return
    }
}


func get_olddata(w http.ResponseWriter, r *http.Request) *bytes.Buffer {
    c := appengine.NewContext(r)
    q := datastore.NewQuery("FreeNumbers").
        Filter("Date >", 1000).
        Order("-Date").
        Limit(5000)
    nums := make([]FreeNumbers, 0, 10)
    if _, err := q.GetAll(c, &nums); err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
    }
    data1 := bytes.NewBufferString("");
    data2 := bytes.NewBufferString("");
    for _, n := range nums {
         fmt.Fprintf(data1, "\t  [new Date(%d), \t%d, %d, %d, %d, %d, %d],\n",
               n.Date/1000, n.Gynaecologists, n.Venereologist, n.Cardiologist,
                            n.Neurologist, n.Otolaryngologist, n.Ophthalmologist)
         fmt.Fprintf(data2, "\t  [new Date(%d), \t%d, %d, %d, %d, %d, %d],\n",
               n.Date/1000, n.Proctologist, n.Rheumatologist, n.Physician,
                            n.Urologist, n.Surgeon, n.Endocrinologist)
    }

    data := bytes.NewBufferString("");
    fmt.Fprintf(data, "data.addRows([\n%s\n]);\ndata2.addRows([\n%s\n]);", data1, data2)
    return data
}

func post(r *http.Request, uri, postvars string) string {
    c:=appengine.NewContext(r)
    t := urlfetch.Transport{Context:c, DeadlineSeconds: 5.0}
    client := http.Client{Transport: &t}
    // resp, err := client.Get("http://punklan.net/cp1251")
    var params, _ = url.ParseQuery(postvars)
    resp, err := client.PostForm(uri, params)
    if err != nil {
        log.Print("err %s", err.String())
    }
    b,err:=ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Print("err %s", err.String())
    }
    buffer := bytes.NewBufferString("");
    for _, char := range b {
    	 var ch = utf(char)
         fmt.Fprintf(buffer, "%c", ch)
    }
    return string(buffer.Bytes())
}

var IndexTemplate = template.Must(template.New("book").Parse(IndexHTML))

const IndexHTML =
`<html>
  <head>
    <script type='text/javascript' src='http://www.google.com/jsapi'></script>
    <script type='text/javascript'>
      google.load('visualization', '1', {'packages':['annotatedtimeline']});
      google.setOnLoadCallback(drawChart);
      function drawChart() {
        var data = new google.visualization.DataTable();
        var data2 = new google.visualization.DataTable();
        data.addColumn('date', 'Date');
	data.addColumn('number', 'гинекол.');
	data.addColumn('number', 'венеролог');
	data.addColumn('number', 'кардиолог');
	data.addColumn('number', 'невропатолог');
	data.addColumn('number', 'отоларинголог');
	data.addColumn('number', 'офтальмолог');
        data2.addColumn('date', 'Date');
	data2.addColumn('number', 'проктолог');
	data2.addColumn('number', 'ревматолог');
	data2.addColumn('number', 'терапевт');
	data2.addColumn('number', 'уролог');
	data2.addColumn('number', 'хирург');
	data2.addColumn('number', 'эндокринолог');

{{.}}

        var chart = new google.visualization.AnnotatedTimeLine(document.getElementById('chart_div'));
        chart.draw(data, {displayAnnotations: true});
        var chart = new google.visualization.AnnotatedTimeLine(document.getElementById('chart_div2'));
        chart.draw(data2, {displayAnnotations: true});
      }
    </script>
  </head>

  <body>
    <div id='chart_div' style='width: 980px; height: 340px;'></div>
    <div id='chart_div2' style='width: 980px; height: 340px;'></div>
  </body>
</html>`

func utf(a byte) (b uint16) { // http://unicode.org/Public/MAPPINGS/VENDORS/MICSFT/WINDOWS/CP1251.TXT
    switch a {
	case 0x00: b = 0x0000
	case 0x01: b = 0x0001
	case 0x02: b = 0x0002
	case 0x03: b = 0x0003
	case 0x04: b = 0x0004
	case 0x05: b = 0x0005
	case 0x06: b = 0x0006
	case 0x07: b = 0x0007
	case 0x08: b = 0x0008
	case 0x09: b = 0x0009
	case 0x0A: b = 0x000A
	case 0x0B: b = 0x000B
	case 0x0C: b = 0x000C
	case 0x0D: b = 0x000D
	case 0x0E: b = 0x000E
	case 0x0F: b = 0x000F
	case 0x10: b = 0x0010
	case 0x11: b = 0x0011
	case 0x12: b = 0x0012
	case 0x13: b = 0x0013
	case 0x14: b = 0x0014
	case 0x15: b = 0x0015
	case 0x16: b = 0x0016
	case 0x17: b = 0x0017
	case 0x18: b = 0x0018
	case 0x19: b = 0x0019
	case 0x1A: b = 0x001A
	case 0x1B: b = 0x001B
	case 0x1C: b = 0x001C
	case 0x1D: b = 0x001D
	case 0x1E: b = 0x001E
	case 0x1F: b = 0x001F
	case 0x20: b = 0x0020
	case 0x21: b = 0x0021
	case 0x22: b = 0x0022
	case 0x23: b = 0x0023
	case 0x24: b = 0x0024
	case 0x25: b = 0x0025
	case 0x26: b = 0x0026
	case 0x27: b = 0x0027
	case 0x28: b = 0x0028
	case 0x29: b = 0x0029
	case 0x2A: b = 0x002A
	case 0x2B: b = 0x002B
	case 0x2C: b = 0x002C
	case 0x2D: b = 0x002D
	case 0x2E: b = 0x002E
	case 0x2F: b = 0x002F
	case 0x30: b = 0x0030
	case 0x31: b = 0x0031
	case 0x32: b = 0x0032
	case 0x33: b = 0x0033
	case 0x34: b = 0x0034
	case 0x35: b = 0x0035
	case 0x36: b = 0x0036
	case 0x37: b = 0x0037
	case 0x38: b = 0x0038
	case 0x39: b = 0x0039
	case 0x3A: b = 0x003A
	case 0x3B: b = 0x003B
	case 0x3C: b = 0x003C
	case 0x3D: b = 0x003D
	case 0x3E: b = 0x003E
	case 0x3F: b = 0x003F
	case 0x40: b = 0x0040
	case 0x41: b = 0x0041
	case 0x42: b = 0x0042
	case 0x43: b = 0x0043
	case 0x44: b = 0x0044
	case 0x45: b = 0x0045
	case 0x46: b = 0x0046
	case 0x47: b = 0x0047
	case 0x48: b = 0x0048
	case 0x49: b = 0x0049
	case 0x4A: b = 0x004A
	case 0x4B: b = 0x004B
	case 0x4C: b = 0x004C
	case 0x4D: b = 0x004D
	case 0x4E: b = 0x004E
	case 0x4F: b = 0x004F
	case 0x50: b = 0x0050
	case 0x51: b = 0x0051
	case 0x52: b = 0x0052
	case 0x53: b = 0x0053
	case 0x54: b = 0x0054
	case 0x55: b = 0x0055
	case 0x56: b = 0x0056
	case 0x57: b = 0x0057
	case 0x58: b = 0x0058
	case 0x59: b = 0x0059
	case 0x5A: b = 0x005A
	case 0x5B: b = 0x005B
	case 0x5C: b = 0x005C
	case 0x5D: b = 0x005D
	case 0x5E: b = 0x005E
	case 0x5F: b = 0x005F
	case 0x60: b = 0x0060
	case 0x61: b = 0x0061
	case 0x62: b = 0x0062
	case 0x63: b = 0x0063
	case 0x64: b = 0x0064
	case 0x65: b = 0x0065
	case 0x66: b = 0x0066
	case 0x67: b = 0x0067
	case 0x68: b = 0x0068
	case 0x69: b = 0x0069
	case 0x6A: b = 0x006A
	case 0x6B: b = 0x006B
	case 0x6C: b = 0x006C
	case 0x6D: b = 0x006D
	case 0x6E: b = 0x006E
	case 0x6F: b = 0x006F
	case 0x70: b = 0x0070
	case 0x71: b = 0x0071
	case 0x72: b = 0x0072
	case 0x73: b = 0x0073
	case 0x74: b = 0x0074
	case 0x75: b = 0x0075
	case 0x76: b = 0x0076
	case 0x77: b = 0x0077
	case 0x78: b = 0x0078
	case 0x79: b = 0x0079
	case 0x7A: b = 0x007A
	case 0x7B: b = 0x007B
	case 0x7C: b = 0x007C
	case 0x7D: b = 0x007D
	case 0x7E: b = 0x007E
	case 0x7F: b = 0x007F
	case 0x80: b = 0x0402
	case 0x81: b = 0x0403
	case 0x82: b = 0x201A
	case 0x83: b = 0x0453
	case 0x84: b = 0x201E
	case 0x85: b = 0x2026
	case 0x86: b = 0x2020
	case 0x87: b = 0x2021
	case 0x88: b = 0x20AC
	case 0x89: b = 0x2030
	case 0x8A: b = 0x0409
	case 0x8B: b = 0x2039
	case 0x8C: b = 0x040A
	case 0x8D: b = 0x040C
	case 0x8E: b = 0x040B
	case 0x8F: b = 0x040F
	case 0x90: b = 0x0452
	case 0x91: b = 0x2018
	case 0x92: b = 0x2019
	case 0x93: b = 0x201C
	case 0x94: b = 0x201D
	case 0x95: b = 0x2022
	case 0x96: b = 0x2013
	case 0x97: b = 0x2014
	case 0x98: b = 0x98 // #UNDEFINED
	case 0x99: b = 0x2122
	case 0x9A: b = 0x0459
	case 0x9B: b = 0x203A
	case 0x9C: b = 0x045A
	case 0x9D: b = 0x045C
	case 0x9E: b = 0x045B
	case 0x9F: b = 0x045F
	case 0xA0: b = 0x00A0
	case 0xA1: b = 0x040E
	case 0xA2: b = 0x045E
	case 0xA3: b = 0x0408
	case 0xA4: b = 0x00A4
	case 0xA5: b = 0x0490
	case 0xA6: b = 0x00A6
	case 0xA7: b = 0x00A7
	case 0xA8: b = 0x0401
	case 0xA9: b = 0x00A9
	case 0xAA: b = 0x0404
	case 0xAB: b = 0x00AB
	case 0xAC: b = 0x00AC
	case 0xAD: b = 0x00AD
	case 0xAE: b = 0x00AE
	case 0xAF: b = 0x0407
	case 0xB0: b = 0x00B0
	case 0xB1: b = 0x00B1
	case 0xB2: b = 0x0406
	case 0xB3: b = 0x0456
	case 0xB4: b = 0x0491
	case 0xB5: b = 0x00B5
	case 0xB6: b = 0x00B6
	case 0xB7: b = 0x00B7
	case 0xB8: b = 0x0451
	case 0xB9: b = 0x2116
	case 0xBA: b = 0x0454
	case 0xBB: b = 0x00BB
	case 0xBC: b = 0x0458
	case 0xBD: b = 0x0405
	case 0xBE: b = 0x0455
	case 0xBF: b = 0x0457
	case 0xC0: b = 0x0410
	case 0xC1: b = 0x0411
	case 0xC2: b = 0x0412
	case 0xC3: b = 0x0413
	case 0xC4: b = 0x0414
	case 0xC5: b = 0x0415
	case 0xC6: b = 0x0416
	case 0xC7: b = 0x0417
	case 0xC8: b = 0x0418
	case 0xC9: b = 0x0419
	case 0xCA: b = 0x041A
	case 0xCB: b = 0x041B
	case 0xCC: b = 0x041C
	case 0xCD: b = 0x041D
	case 0xCE: b = 0x041E
	case 0xCF: b = 0x041F
	case 0xD0: b = 0x0420
	case 0xD1: b = 0x0421
	case 0xD2: b = 0x0422
	case 0xD3: b = 0x0423
	case 0xD4: b = 0x0424
	case 0xD5: b = 0x0425
	case 0xD6: b = 0x0426
	case 0xD7: b = 0x0427
	case 0xD8: b = 0x0428
	case 0xD9: b = 0x0429
	case 0xDA: b = 0x042A
	case 0xDB: b = 0x042B
	case 0xDC: b = 0x042C
	case 0xDD: b = 0x042D
	case 0xDE: b = 0x042E
	case 0xDF: b = 0x042F
	case 0xE0: b = 0x0430
	case 0xE1: b = 0x0431
	case 0xE2: b = 0x0432
	case 0xE3: b = 0x0433
	case 0xE4: b = 0x0434
	case 0xE5: b = 0x0435
	case 0xE6: b = 0x0436
	case 0xE7: b = 0x0437
	case 0xE8: b = 0x0438
	case 0xE9: b = 0x0439
	case 0xEA: b = 0x043A
	case 0xEB: b = 0x043B
	case 0xEC: b = 0x043C
	case 0xED: b = 0x043D
	case 0xEE: b = 0x043E
	case 0xEF: b = 0x043F
	case 0xF0: b = 0x0440
	case 0xF1: b = 0x0441
	case 0xF2: b = 0x0442
	case 0xF3: b = 0x0443
	case 0xF4: b = 0x0444
	case 0xF5: b = 0x0445
	case 0xF6: b = 0x0446
	case 0xF7: b = 0x0447
	case 0xF8: b = 0x0448
	case 0xF9: b = 0x0449
	case 0xFA: b = 0x044A
	case 0xFB: b = 0x044B
	case 0xFC: b = 0x044C
	case 0xFD: b = 0x044D
	case 0xFE: b = 0x044E
	case 0xFF: b = 0x044F
	}
    return
}
