$(document).ready(function() {

    $('.tooltipped').tooltip({delay: 10});

    $('.parallax').parallax();

    $("#settings-menu").sideNav();

    $(".my-blueprints-control").click(function () {

        toggleMyBluePrints($(this).attr("data-typeid"));

    });

    $(document).ready(function(){
        $('.collapsible').collapsible();
    });

    ajaxAutoComplete({inputId: 'region', ajaxUrl: '/list/regions/'});

    $('select').material_select();

    $('#search').on("keypress", function (e) {
        if (e.which === 13) {

            window.location.href = "/?q=" + $(this).val();

        }
    });

    $('#search-btn').on("click", function (e) {
        window.location.href = "/?q=" + $(this).val();
    });

    $("#settings-form").submit(function (e) {

        e.preventDefault();

        e.stopPropagation();

        $.ajax({
            type: 'POST',
            url: "/settings/append/",
            data: $(this).serializeArray(),
            success: function (data) {
                if (data.Error != null) {
                    Materialize.toast('Error: <span class="red-text">' + data.Error + '</span>', 4000)
                } else {
                    Materialize.toast('Settings saved', 4000)
                }
            }
        });

        return false;

    });
});

function toggleMyBluePrints(typeID) {
    $.ajax({
        type: 'POST',
        url: '/toggle-my-blueprint/',
        data: {typeid : typeID},
        success: function (data) {
            if(data.Error != null) {

                Materialize.toast('Error: <span class="red-text">'+data.Error+'</span>', 4000)

            } else {

                var $toastContent = "";

                var cc = parseInt($("#pinned-count").html());

                if (data.Result) {

                    $toastContent = $('<span>Blueprint pinned</span>').add($('<a class="btn-flat toast-action" href="/my-blueprints/">View pinned blueprints</a>'));

//                    $("#bpo-menu").html(cnt);

                } else {

                    $toastContent = $('<span>Blueprint unpinned</span>').add($('<a class="btn-flat toast-action" href="/my-blueprints/">View pinned blueprints</a>'));

//                    $("#bpo-menu").html(cnt);

                }

                Materialize.toast($toastContent, 5000);

            }
        }
    });
}

function ajaxAutoComplete(options)
{

    var defaults = {
        inputId:null,
        ajaxUrl:false,
        data: {},
        minLength: 2
    };

    options = $.extend(defaults, options);
    var $input = $("#" + options.inputId);


    if (options.ajaxUrl){


        var $autocomplete = $('<ul id="ac" class="autocomplete-content dropdown-content"'
            + 'style="position:absolute"></ul>'),
            $inputDiv = $input.closest('.input-field'),
            request,
            runningRequest = false,
            timeout,
            liSelected;

        if ($inputDiv.length) {
            $inputDiv.append($autocomplete); // Set ul in body
        } else {
            $input.after($autocomplete);
        }

        var highlight = function (string, match) {
            var matchStart = string.toLowerCase().indexOf("" + match.toLowerCase() + ""),
                matchEnd = matchStart + match.length - 1,
                beforeMatch = string.slice(0, matchStart),
                matchText = string.slice(matchStart, matchEnd + 1),
                afterMatch = string.slice(matchEnd + 1);
            string = "<span>" + beforeMatch + "<span class='highlight'>" +
                matchText + "</span>" + afterMatch + "</span>";
            return string;

        };

        $autocomplete.on('click', 'li', function () {

            $("#regionID").val(($(this).attr("data-id")));

            $input.val($(this).text().trim());
            $autocomplete.empty();
        });

        $input.on('keyup', function (e) {

            if (timeout) { // comment to remove timeout
                clearTimeout(timeout);
            }

            if (runningRequest) {
                request.abort();
            }

            if (e.which === 13) { // select element with enter key
                liSelected[0].click();
                return;
            }

            // scroll ul with arrow keys
            if (e.which === 40) {   // down arrow
                if (liSelected) {
                    liSelected.removeClass('selected');
                    next = liSelected.next();
                    if (next.length > 0) {
                        liSelected = next.addClass('selected');
                    } else {
                        liSelected = $autocomplete.find('li').eq(0).addClass('selected');
                    }
                } else {
                    liSelected = $autocomplete.find('li').eq(0).addClass('selected');
                }
                return; // stop new AJAX call
            } else if (e.which === 38) { // up arrow
                if (liSelected) {
                    liSelected.removeClass('selected');
                    next = liSelected.prev();
                    if (next.length > 0) {
                        liSelected = next.addClass('selected');
                    } else {
                        liSelected = $autocomplete.find('li').last().addClass('selected');
                    }
                } else {
                    liSelected = $autocomplete.find('li').last().addClass('selected');
                }
                return;
            }

            // escape these keys
            if (e.which === 9 ||        // tab
                e.which === 16 ||       // shift
                e.which === 17 ||       // ctrl
                e.which === 18 ||       // alt
                e.which === 20 ||       // caps lock
                e.which === 35 ||       // end
                e.which === 36 ||       // home
                e.which === 37 ||       // left arrow
                e.which === 39) {       // right arrow
                return;
            } else if (e.which === 27) { // Esc. Close ul
                $autocomplete.empty();
                return;
            }

            var val = $input.val().toLowerCase();
            $autocomplete.empty();

            if (val.length > options.minLength) {

                timeout = setTimeout(function () { // comment this line to remove timeout
                    runningRequest = true;

                    request = $.ajax({
                        type: 'GET',
                        url: options.ajaxUrl + '?term=' + val,
                        success: function (data) {
                            if (!$.isEmptyObject(data)) { // (or other) check for empty result
                                var appendList = '';
                                for (var key in data) {

                                    var li = '';
                                        li += '<li data-id="'+data[key].value+'">' + highlight(data[key].label, val) + '</li>';
                                    appendList += li;
                                }
                                $autocomplete.append(appendList);
                            }else{
                                $autocomplete.append($('<li>No matches</li>'));
                            }
                        },
                        complete: function (d) {
                            runningRequest = false;
                        }
                    });
                }, 250);        // comment this line to remove timeout
            }
        });

        $(document).click(function (event) { // close ul if clicked outside
            if (!$(event.target).closest($autocomplete).length) {
                $autocomplete.empty();
            }
        });
    }
}


function ajaxTimeLine(opts){

    $.getJSON('/timeline-data/'+opts.regionId+'/?batch='+opts.batch, function (data_i) {

        if (data_i.Error != null) {
            $("#container").html("Error: " + data_i.Error);
            return
        }

        data = data_i.Result;

        // split the data set into ohlc and volume
        var categories = [],
            volume = [],
            avg = [],
            dataLength = data.C.length,
            // set the allowed units for data grouping
            groupingUnits = [[
                'week',                         // unit name
                [1]                             // allowed multiples
            ], [
                'month',
                [1, 2, 3, 4, 6]
            ]],

            i = 0;

        for (i; i < dataLength; i += 1) {

            for (j=0; j < data.C[i].length; j += 1) {

                categories.push(data.C[i][j].TypeName + ' x'+data.C[i][j].Quantity);

                volume.push({
                    x: Date.parse(data.C[i][j].Start),
                    x2: Date.parse(data.C[i][j].End),
                    y: categories.length - 1
                })
            }

        }

        // create the chart
        Highcharts.chart('container', {
            chart: {
                type: 'xrange'
            },
            title: {
                text: 'Research and production timeline'
            },
            xAxis: {
                type: 'datetime'
            },
            yAxis: {
                title: {
                    text: ''
                },
                categories: categories,
                reversed: true
            },
            series: [{
                name: "Manufacturing",
                pointPadding: 0,
                groupPadding: 0,
                borderColor: 'gray',
                pointWidth: 20,
                colorByPoint: false,
                color: '#78909c',
                data: volume,
                dataLabels: {
                    enabled: true,
                    formatter: function(){return formatDuration(this.x2 - this.x);}
                }
            }
            ]

        });
    });
}


function formatDuration (seconds) {

    function numberEnding (number) {
        return '';
    }

    seconds = seconds / 1000;

    if (seconds > 0){
        //var years = Math.floor(seconds / 31536000);
        var days = Math.floor(seconds / 86400);
        var hours = Math.floor((seconds % 86400)  / 3600);
        var minutes = Math.floor((((seconds % 86400) % 3600) %  60));
        var second = (((seconds % 31536000) % 86400) % 3600) % 0;
        //var r = (years > 0 ) ? years + " year" + numberEnding(years) : "";
        var x = (days > 0) ? days + "d" + numberEnding(days) : "";
        var y = (hours > 0) ? hours + "h" + numberEnding(hours) : "";
        var z = (minutes > 0) ? minutes + "m" + numberEnding(minutes) : "";
        var u = (second > 0) ? second + "s" + numberEnding(second) : "";
        return x + y + z + u;

    }
    else {
        return "now"
    }
}
