package app

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
    "app/cp1251_utf8"
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
		case "акушер-гинекол.":	a.Gynaecologists = count_int
		case "дерма-венеролог":	a.Venereologist = count_int
		case "кардиолог":	a.Cardiologist = count_int
		case "невропатолог":	a.Neurologist = count_int
		case "отоларинголог":	a.Otolaryngologist = count_int
		case "офтальмолог":	a.Ophthalmologist = count_int
		case "проктолог":	a.Proctologist = count_int
		case "ревматолог":	a.Rheumatologist = count_int
		case "терапевт участ.":	a.Physician = count_int
		case "уролог":		a.Urologist = count_int
		case "хирург":		a.Surgeon = count_int
		case "эндокринолог":	a.Endocrinologist = count_int
	    }
            fmt.Fprint(w, name, " -> ", count_int, "\n")
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
        Order("-Date").
        Limit(4032)
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
         var ch = cp1251_utf8.Utf(char)
         fmt.Fprintf(buffer, "%c", ch)
    }
    return string(buffer.Bytes())
}

var IndexTemplate = template.Must(template.New("book").Parse(IndexHTML))

const IndexHTML =
`<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>Статистика по свободным номеркам Николаевской больницы Петродворцового района</title>
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
    <style type="text/css">
        body {margin:0;padding:0;}
        .wrapper {width: 985px; margin: 10px auto;}
        .intro {margin: 10px; text-align: center;}
        .footer {margin: 15px auto; text-align: center;}
    </style>
  </head>

  <body>
    <div class="wrapper">
      <div class="intro">
        <p>На этой страничке вы можете посмотреть статистику занятости номерков на прием в поликлинике
        <a target="_blank" href="http://nikmed.spb.ru/">Николаевской больницы</a>.</p>
        <p>Информация о доступности номерков проверяется каждые 5 минут.</p>
      </div>
      <div id='chart_div' style='width: 980px; height: 340px;'></div>
      <div id='chart_div2' style='width: 980px; height: 340px;'></div>
      <div class="footer">
        <a rel="me" href="https://plus.google.com/104701071096191312666/about">JLarky</a>
          &copy; 2012 <a href="https://github.com/JLarky/nikmed-stats">
          <img src="images/github_logo.png" style="vertical-align: -6px;" alt="github"></a>
      </div>
    </div>
  </body>
 <script type="text/javascript">
  var _gaq = _gaq || [];
  _gaq.push(['_setAccount', 'UA-5041583-9']);
  _gaq.push(['_trackPageview']);

  (function() {
    var ga = document.createElement('script'); ga.type = 'text/javascript'; ga.async = true;
    ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
    var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(ga, s);
  })();
 </script>
</html>`
