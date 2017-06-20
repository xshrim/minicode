// ==UserScript==
// @name 	自动报工
// @namespace http://tampermonkey.net/
// @version 0.1
// @description 	一键生成报工事项！
// @author xshrim
// @match *://130.1.12.49:9080/project/index.jsp
// @Grant none
// ==/UserScript==

var now = new Date(); //当前日期
var nowDayOfWeek = now.getDay() - 1; //今天本周的第几天
var nowDay = now.getDate(); //当前日
var nowMonth = now.getMonth(); //当前月
var nowYear = now.getFullYear(); //当前年

//数字补齐(未使用)
function pad(num, n) {
  return Array(n>num?(n-(''+num).length+1):0).join(0)+num;
}

//格式化日期：yyyy-MM-dd
function formatDate(date) {
    var myyear = date.getFullYear();
    var mymonth = date.getMonth() + 1;
    var myweekday = date.getDate();
    if(mymonth < 10){
        mymonth = "0" + mymonth;
    }
    if(myweekday < 10){
        myweekday = "0" + myweekday;
    }
    return (myyear+"-"+mymonth + "-" + myweekday);
}

//获取当前日期
function getWeekCurrentDate() {
    var weekCurrentDate = new Date(nowYear, nowMonth, nowDay);
    return formatDate(weekCurrentDate);
}

//获取本周开始日期
function getWeekStartDate() {
    var weekStartDate = new Date(nowYear, nowMonth, nowDay - nowDayOfWeek);
    return formatDate(weekStartDate);
}

//获取本周结束日期
function getWeekEndDate() {
    var weekEndDate = new Date(nowYear, nowMonth, nowDay + (6 - nowDayOfWeek));
    return formatDate(weekEndDate);
}

//获取两个日期间的所有日期
function getAll(begin,end){
    var dates = [];
    var ab = begin.split("-");
    var ae = end.split("-");
    var db = new Date();
    db.setUTCFullYear(ab[0], ab[1]-1, ab[2]);
    var de = new Date();
    de.setUTCFullYear(ae[0], ae[1]-1, ae[2]);
    var unixDb=db.getTime();
    var unixDe=de.getTime();
    for(var k=unixDb;k<=unixDe;){
        dates.push(formatDate(new Date(parseInt(k))));
        k=k+24*60*60*1000;
    }
    return dates;
}

//根据年和周数获取该周所有日期(未使用)
function yugi(year, index) {
    var d = new Date(year, 0, 1);
    while (d.getDay() != 1) {
        d.setDate(d.getDate() + 1);
    }
    var to = new Date(year + 1, 0, 1);
    var i = 1;
    var arr = [];
    for (var from = d; from < to;) {
        if (i == index) {
            arr.push(from.getFullYear() + "-" + pad((from.getMonth() + 1), 2) + "-" + pad(from.getDate(), 2));
        }
        var j = 6;
        while (j > 0) {
            from.setDate(from.getDate() + 1);
            if (i == index) {
                arr.push(from.getFullYear() + "-" + pad((from.getMonth() + 1) , 2)+ "-" + pad(from.getDate(), 2));
            }
            j--;
        }
        if (i == index) {
            return arr;
        }
        from.setDate(from.getDate() + 1);
        i++;
    }
}

//自动报工实现
window.onload = function() {
    'use strict';
    var dates = getAll(getWeekStartDate(), getWeekEndDate());
    var $octd=$('<td class="x-toolbar-cell" id="ext-gen12"><table id="oneclickBtn" cellspacing="0" class="x-btn   x-btn-text-icon" style="width: auto;"><tbody class="x-btn-small x-btn-icon-small-left"><tr><td class="x-btn-tl"><i>&nbsp;</i></td><td class="x-btn-tc"></td><td class="x-btn-tr"><i>&nbsp;</i></td></tr><tr><td class="x-btn-ml"><i>&nbsp;</i></td><td class="x-btn-mc"><em class="" unselectable="on"><button type="button" id="ext-octd" class=" x-btn-text btn_icon_005">一键报工</button></em></td><td class="x-btn-mr"><i>&nbsp;</i></td></tr><tr><td class="x-btn-bl"><i>&nbsp;</i></td><td class="x-btn-bc"></td><td class="x-btn-br"><i>&nbsp;</i></td></tr></tbody></table></td>');
    document.getElementById("rightFrame").onload = function() {
        $("#rightFrame").contents().find("tr.x-toolbar-left-row").append($octd);
        //$("#rightFrame").contents().find("tr.x-toolbar-left-row").on('mouseover', '#oneclickBtn', function(){
        //$("#rightFrame").contents().find("'#oneclickBtn").attr('class', 'x-btn   x-btn-text-icon  x-btn-over');
        //});
        $("#rightFrame").contents().find("tr.x-toolbar-left-row").on('click', '#ext-octd', function(){
            var viid = '0';
            var vstartDate =getWeekStartDate();
            var vendDate = getWeekEndDate();
            var vtaskid = '';
            //var vddate = getWeekCurrentDate();
            var vstype = '3';
            var vprojectid = '';
            var vtasktype = '0';
            var vtaskidIn = '';
            var vtasknameIn = '';
            var vtimetypeid = '990008007';
            var vtimetypename = '技术支持';
            var vworkname = '';
            var vworknameyw2 = '';
            var vproduct = '';
            var vchanpinname = '';
            var vmanhour = '8.0';
            var vmanovertime = '0.0';
            var vremark = '^_^';
            var content=prompt("请输入每日工作内容(以|分隔!)","");
            var contents = content.split('|');
            $.each(dates,function(index,vddate){
                if(index < 5){
                    if(contents === null || contents.length < index + 1)
                        vremark = "^_^";
                    else
                        vremark = contents[index];
                    if($.trim(vremark) === "")
                        vremark = "^_^";
                    $.post("http://130.1.12.49:9080/project/vocation/timesheet/handle/outTimesheet.jsp?method=save&submitstate=0", { iid: viid, startDate: vstartDate, endDate: vendDate, taskid: vtaskid, ddate: vddate, stype: vstype, projectid: vprojectid, tasktype: vtasktype, taskidIn: vtaskidIn, tasknameIn: vtasknameIn, timetypeid: vtimetypeid, workname: vworkname, worknameyw2: vworknameyw2, product: vproduct, chanpinname: vchanpinname, manhour: vmanhour, manovertime: vmanovertime, remark: vremark },
                           function(data){
                             consloe.log($.trim(data));
                          }
                    );
                }
            });
            window.location.reload();
        });
    };
};