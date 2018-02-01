var init = false;
var webSocket;
var scpu, smem, myChart;
var serviceFailures = {};

function connectStatus() {

    setTimeout(function(){

        if(!webSocket.isConnected()){

            $("#notice").html("Lost connection to websocket, <a href='/'>reconnect?</a>");

        } else {

            $("#notice").html();

        }

        connectStatus();

    }, 3000);

}

var restart = function(s) {

    $.getJSON("/restart/"+s+"/", {}, function (message) {

        $("#notice").html(message);

        setTimeout(function(){
            $("#notice").html("");
        }, 5000);
    })

};

var redrawFailures = function () {

    for (var service in serviceFailures) {

        if (serviceFailures.hasOwnProperty(service)) {

            $("#messages-"+service).html(serviceFailures[service].join(" | "));

        }
    }

};

var initWS = function(host, wsport){

    webSocket = $.simpleWebSocket({ url: 'ws://'+host+':'+wsport+'/' });

    webSocket.listen(function(message) {

        if(message.pids) {

            var c = $("#services-container");

            newLayout = "<table><tr><th>..</th><th>Service</th><th>PID</th><th>TCP conns</th><th>Online errs</th></tr>";

            var cats = [];

            for (var i=0; i<message.pids.length; i++) {

                var msgp = "";

                if(message.pids[i].PID.length > 1) {

                    msgp = message.pids[i].PID.length + " processes";

                } else {

                    if (message.pids[i].PID.length == 1) {

                        msgp = message.pids[i].PID[0];

                    }

                }


                newLayout += "<tr>";

                newLayout += "<td><a href='/log/"+message.pids[i].Name+"/'><i class='fa fa-bars'></i></a> <a href='#' onclick='restart(\""+message.pids[i].Name+"\");'><i class='fa fa-refresh'></i></a></td>";
                newLayout += "<td>"+message.pids[i].Name+"</td>";
                newLayout += "<td><span class='trim'>"+msgp+"</span><span class='alert'>"+message.pids[i].Error+"</span></td>";
                newLayout += "<td>"+message.pids[i].Conns+"</td>";
                newLayout += "<td id='messages-"+message.pids[i].Name+"'> - </td>";

                newLayout += "</tr>";
                cats.push(message.pids[i].Name);

            }

            newLayout += "</table>";

            c.html(newLayout);

            redrawFailures();

            if (!init) {
                scpu.xAxis[0].setCategories(cats, true);
                smem.xAxis[0].setCategories(cats, true);
                init = true;
            }

        }

        if(message.openTCPSockets) {
            $("#tcp-sockets-open").html("Open TCP sockets: <b>"+message.openTCPSockets+"</b>");
        }

        if(message.memStat) {
            $("#av-mem-c").html(Math.round(message.memStat[0]/1024/1024, 0) + "MB");
            $("#used-mem-c").html(Math.round(message.memStat[1]/1024/1024, 0) + "MB");
            $("#sc-mem-c").html(Math.round(message.memStat[2]/1024/1024, 0) + "MB");
        }

        if(message.serviceFailure){

            if (!serviceFailures.hasOwnProperty(message.serviceFailure.service)) {

                serviceFailures[message.serviceFailure.service] = new(Array);

            }

            serviceFailures[message.serviceFailure.service].push(
                '<a href="/message/'+message.serviceFailure.id+'/">'+message.serviceFailure.name+'</a>'
            );

            redrawFailures();
        }

    });

    connectStatus();

    scpu = Highcharts.chart('container-cpu-s', {
        chart: {
            type: 'bar',
            animation: Highcharts.svg,
            events: {
                load: function () {

                    webSocket.listen(function(message) {

                        if(!init){
                            console.log("Too early");
                            return;
                        }

                        if(message.ps) {

                            var dataI = [];

                            for(var i=0; i<scpu.xAxis[0].categories.length; i++) {

                                dataI.push(Math.round((message.ps[scpu.xAxis[0].categories[i]].TimeP * 100), 4))

                            }

                            scpu.series[0].setData(dataI);
                        }

                    });

                }
            }
        },
        title: {
            text: 'CPU by service'
        },
        subtitle: {
            text: 'Source: /proc/pid'
        },
        xAxis: {
            categories: ['-','-','-','-','-','-','-','-','-','-'],
            title: {
                text: null
            },
            labels: {
                step: 1,
                padding: 0
            }
        },
        yAxis: {
            min: 0,
            title: {
                text: 'CPU by service',
                align: 'high'
            },
            labels: {
                overflow: 'justify'
            }
        },
        tooltip: {
            valueSuffix: ' %'
        },
        plotOptions: {
            bar: {
                dataLabels: {
                    enabled: false
                }
            }
        },
        legend: {
            layout: 'vertical',
            align: 'right',
            verticalAlign: 'top',
            x: -40,
            y: 80,
            floating: true,
            borderWidth: 1,
            backgroundColor: ((Highcharts.theme && Highcharts.theme.legendBackgroundColor) || '#FFFFFF'),
            shadow: true
        },
        credits: {
            enabled: false
        },
        series: [{
            name: "cpu %",
            data: [0,0,0,0,0,0,0,0,0,0]
        }]
    });

    smem = Highcharts.chart('container-mem-s', {
        chart: {
            type: 'bar',
            animation: Highcharts.svg,
            events: {
                load: function () {

                    webSocket.listen(function(message) {

                        if(!init){
                            console.log("Too early");
                            return;
                        }

                        if(message.ps) {

                            var dataI = [[],[],[]];

                            for(var i=0; i<smem.xAxis[0].categories.length; i++) {
                                dataI[0].push(Math.round(((message.ps[smem.xAxis[0].categories[i]].Vsize) / 1024), 0));
                                dataI[1].push(Math.round(((message.ps[smem.xAxis[0].categories[i]].Rss) / 1024), 0));
                                dataI[2].push(Math.round(((message.ps[smem.xAxis[0].categories[i]].Swap) / 1024), 0));

                            }

                            smem.series[0].setData(dataI[0]);
                            smem.series[1].setData(dataI[1]);
                            smem.series[2].setData(dataI[2]);
                        }

                    });

                }
            }
        },
        title: {
            text: 'Memory consumption by service'
        },
        subtitle: {
            text: 'Source: /proc/pid'
        },
        xAxis: {
            categories: ['-','-','-','-','-','-','-','-','-','-'],
            title: {
                text: null
            },
            labels: {
                step: 1,
                padding: 0
            }
        },
        yAxis: {
            min: 0,
            title: {
                text: 'VmData, VmStk, VmExe, VmRSS, VmSwap by service',
                align: 'high'
            },
            labels: {
                overflow: 'justify'
            }
        },
        tooltip: {
            valueSuffix: ' MB'
        },
        plotOptions: {
            bar: {
                dataLabels: {
                    enabled: false
                }
            }
        },
        legend: {
            layout: 'vertical',
            align: 'right',
            verticalAlign: 'top',
            x: -40,
            y: 80,
            floating: true,
            borderWidth: 1,
            backgroundColor: ((Highcharts.theme && Highcharts.theme.legendBackgroundColor) || '#FFFFFF'),
            shadow: true
        },
        credits: {
            enabled: false
        },
        series: [{
            name: "VmData, VmStk, VmExe",
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
        },
        {
            name: "VmRSS",
            data: [0,0,0,0,0,0,0,0,0,0]
        },
        {
            name: "VmSwap",
            data: [0,0,0,0,0,0,0,0,0,0]
        }]
    });

    myChart = Highcharts.chart('container', {
        chart: {
            type: 'line',
            animation: Highcharts.svg,
            events: {
                load: function () {

                    webSocket.listen(function(message) {

                        if(message.cpuPercent) {

                            var x = (new Date()).getTime();

                            var t = 0;

                            for(var i=0; i<message.cpuPercent.length; i++) {

                                var y = Math.round(message.cpuPercent[i], 2);

                                var series = myChart.series[i];

                                shift = series.data.length > 60;

                                series.addPoint([x, y], true, shift);
                            }

                        }

                    });

                }
            }
        },
        title: {
            text: 'CPU/RAM system'
        },
        xAxis: {
            type: 'datetime',
            title: {
                enabled: false
            },
            tickPixelInterval: 150
        },
        yAxis: {
            title: {
                text: 'Percent'
            },
            labels: {
                formatter: function () {
                    return this.value + '%';
                }
            }
        },
        tooltip: {
            split: true,
            valueSuffix: ' %'
        },
        plotOptions: {
            line: {
                dataLabels: {
                    enabled: true
                },
                enableMouseTracking: false
            }
        },
        series: [
            {
                name: 'Total CPU',
                data: []
            },
            {
                name: "Total used mem",
                data: []
            }
        ]
    });

};