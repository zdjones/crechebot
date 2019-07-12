// Author: dg_dgrainger
// Project: LegendOnline.CMS
// Creation date: 3/24/2009 10:27 AM
// Last modified: 11/3/2009 10:47 AM

/// <summary>
/// Display club
/// </summary>
function displayClub(id) {
    parent.tb_show('Facility Details', '../MembershipCentre/Facility?id=' + id + '&KeepThis=true&TB_iframe=true&height=350&width=450&modal=true', '');
} // displayClub()

function subClub(item) {
    // Submit Clubs
    var clubSubmit = $("#clubSubmit");
    clubSubmit.click();
    // Update Courts
    fnHide($('.viewTable'));
}

function subAct(item) {
    // Update Courts
    var actSub = $("#actSubmit");
    actSub.click();
    checkActivitySelection();
}

function updateBeh() {
    var actSub = $("#behSubmit");
    actSub.click();

    if ($('#View_Day_Radio').is(":checked")) {
        $('#View_Activity').hide();
        $("#SubViewDay").click();
    }

    fnHide($('.viewTable'));
}

function subBeh(item) {
    // (re)set the value of the hidden bookingType input before sending AJAX POST
    var bookingType = $(item).attr('data-booking-type'); // .data() doesn't return anything
    $('#hiddenBookingType').val(bookingType);

    // Update Courts
    var actSub = $("#behSubmit");
    actSub.click();
    fnHide($('.viewTable'));
}

/// <summary>
/// Load hover
/// </summary>
function loadHover() {
    $('a.jTip200').cluetip({ width: '200px', showTitle: false });
} // loadHover()

/// <summary>
/// Load only hover
/// </summary>
function loadOnlyHover() {
    $('a.jTip350').cluetip({ width: '350px', showTitle: false });
} // loadOnlyHover()

/// <summary>
/// Load only hover
/// </summary>
function loadPriceHover() {
    $('a.jTip100').cluetip({ width: '120px', showTitle: false });
} // loadOnlyHover()

/// <summary>
/// Test
/// </summary>
function fnHide(item) {
    $(item).fadeOut("slow");
} // test(text)

function fnShow(item) {
    $(item).fadeIn("slow");
} // test(text)

/// <summary>
/// Check activity selection
/// </summary>
function checkActivitySelection() {
    if ($('.activityCheckbox').is(':checked')) {
        
        // Check how many we have. Only show the top if we have more than 10 results.
        var count = 0;
        $('.activityCheckbox').each(function () {
            count++;
        });

        if (count > 10) {
            setTimeout("fnShow($('#topsubmit'));", 500);
        }

        setTimeout("fnShow($('#bottomsubmit'));", 500);
    }
    else {
        fnHide($('.viewTable'));
    }
} // checkActivitySelection()

/// <summary>
/// Apply json
/// </summary>
function addBooking(id, url) {

    url = (typeof url === "undefined") ? "AddBooking?booking=" + id : url;

    var originalContent = $("#slot" + id).replaceWith("<div id='slot" + id + "'><span class='ajaxLoader'>[ PROCESSING ]</span></div>");

    $.getJSON(url + "&ajax=" + Math.random(), null, function (data) {
        // Check success
        if (data.Success == true) {
            $("#slot" + id).replaceWith("[ IN BASKET ]").html();
        }
        else if (data.AllowRetry) {
            $("#slot" + id).replaceWith(originalContent);
        }
        else {
            $("#slot" + id).replaceWith("[ CONFLICT ]").html();
        }

        var params = $.param(formatDatesFromJsonObject(data));        
        tb_show("Bookings", "Message?" + params + "&KeepThis=true&TB_iframe=true&height=260&width=340&modal=true", "");     // Show modal
    });
} // applyJson()

function addToClassWaitingList(id) {
    tb_show("Bookings", "AddToWaitingListConfirmation?activityInstanceId=" + id + "&successCallback=addToWaitingListSuccessCallback&errorCallback=addToWaitingListErrorCallback&KeepThis=true&TB_iframe=true&height=280&width=340&modal=true", "");     // Show modal
}

var onWaitingListText = 'On Waiting List';

function setWaitingListText(newText) {
    onWaitingListText = newText;
}

function addToWaitingListSuccessCallback(id) {
    $("#slot" + id).replaceWith("[ " + onWaitingListText + " ]").html();
}

function addToWaitingListErrorCallback(id) {
    $("#slot" + id).replaceWith("[ CONFLICT ]").html();
}

function addSportsHallBooking(id, selectedCourts) {
    $(".sportsHallSubItem").fadeOut();

    selectedCourts = (typeof selectedCourts === "undefined") ? '' : selectedCourts;

    addBooking(id, "AddSportsHallBooking?slotId=" + id + "&selectedCourts=" + selectedCourts);
}

/// <summary>
/// Takes in the user Court selection and calls AddBooking
/// </summary>
function confirmCourtSelect(slotId, selectedCourts)
{
    selectedCourts = (selectedCourts == "-1") ? '' : selectedCourts;  // Any selected

    addSportsHallBooking(slotId, selectedCourts);
} // confirmCourtSelect(selectedCourts)

function parseJsonDateForURI(jsonDate)
{
    if (jsonDate == null)
    {
        return "";
    }
    else
    {   // Parse date from milliseconds
        var date = new Date(parseInt(jsonDate.replace("/Date(", "").replace(")/", ""), 10));

        return date.toUTCString();
    }
}

function formatDatesFromJsonObject(data)
{
    for (var key in data)
    {
        var value = data[key];

        if (value != null && value.toString().indexOf("/Date(") > -1)
        {   // JSON date found
            data[key] = parseJsonDateForURI(value);
        }
    }

    return data;
}

/// <summary>
/// View timetable
/// </summary>
function viewTimetable() {
    parent.tb_show("Bookings", "../BookingsCentre/Timetable?KeepThis=true&TB_iframe=true&height=500&width=850&modal=true", "");
} // viewTimetable()

function loadSporsHallExplodedView(link, id) {
    $(".sportsHallSubItem").fadeOut;
    
    var p = $(link);
    var position = p.position();

    $(id).fadeIn();
    $(id).css("left", position.left + 50 + "px");
    $(id).css("top", position.top + 10 + "px");
}

function selectResourceLocation(slotId)
{
    $('.sportsHallSubItem').fadeOut();

    var params = "slotId=" + slotId;
    // Show court selection modal
    tb_show("Bookings", "SelectResourceLocation?" + params + "&KeepThis=true&TB_iframe=true&height=240&width=360&modal=true", "");
}

var activitiesLoaded = function() {
    $('#activitiesLoading').hide();
    $('#activities').show();
};

var activitiesLoading = function () {
    $('#activitiesLoading').show();
    $('#activities').hide();
};