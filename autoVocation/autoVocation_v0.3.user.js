// ==UserScript==
// @name 	自动报工
// @namespace http://tampermonkey.net/
// @version 0.3
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
var content = null;//报工内容

//判断是否为数组
function isArray(v) {
    return v && typeof v.length == 'number' && typeof v.splice == 'function';
}
//创建元素
function createEle(tagName) {
    return document.createElement(tagName);
}
//在body中添加子元素
function appChild(eleName) {
    return document.body.appendChild(eleName);
}
//从body中移除子元素
function remChild(eleName) {
    return document.body.removeChild(eleName);
}
//弹出窗口，标题（html形式）、html、默认关闭动作、默认关闭动作的参数、宽度、高度、是否为模式对话框(true,false)、按钮（关闭按钮为默认，格式为['按钮1',fun1,'按钮2',fun2]数组形式，前面为按钮值，后面为按钮onclick事件）
function showWindow(title, html, defaultAction, args, width, height, modal, buttons) {
    //避免窗体过小
    if (width < 300) {
        width = 300;
    }
    if (height < 200) {
        height = 200;
    }

    //声明mask的宽度和高度（也即整个屏幕的宽度和高度）
    var w, h;
    //可见区域宽度和高度
    var cw = document.body.clientWidth;
    var ch = document.body.clientHeight;
    //正文的宽度和高度
    var sw = document.body.scrollWidth;
    var sh = document.body.scrollHeight;
    //可见区域顶部距离body顶部和左边距离
    var st = document.body.scrollTop;
    var sl = document.body.scrollLeft;

    w = cw > sw ? cw : sw;
    h = ch > sh ? ch : sh;

    //避免窗体过大
    if (width > w) {
        width = w;
    }
    if (height > h) {
        height = h;
    }

    //如果modal为true，即模式对话框的话，就要创建一透明的掩膜
    if (modal) {
        var mask = createEle('div');
        mask.style.cssText = "position:absolute;left:0;top:0px;background:#fff;filter:Alpha(Opacity=30);opacity:0.5;z-index:100;width:" + w + "px;height:" + h + "px;";
        appChild(mask);
    }

    //这是主窗体
    var win = createEle('div');
    win.style.cssText = "position:absolute;left:" + (sl + cw / 2 - width / 2) + "px;top:" + (st + ch / 2 - height / 2) + "px;background:#f0f0f0;z-index:101;width:" + width + "px;height:" + height + "px;border:solid 2px #afccfe;";
    //标题栏
    var tBar = createEle('div');
    //afccfe,dce8ff,2b2f79
    tBar.style.cssText = "margin:0;width:" + width + "px;height:25px;cursor:move;";
    //标题栏中的文字部分
    var titleCon = createEle('div');
    titleCon.style.cssText = "float:left;margin:3px; font-weight: bold; font-size: 13px";
    titleCon.innerHTML = title; //firefox不支持innerText，所以这里用innerHTML
    tBar.appendChild(titleCon);
    //标题栏中的“关闭按钮”
    var closeCon = createEle('div');
    closeCon.style.cssText = "float:right;width:20px;margin:3px;cursor:pointer;"; //cursor:hand在firefox中不可用
    closeCon.innerHTML = "×";
    tBar.appendChild(closeCon);
    win.appendChild(tBar);
    //窗体的内容部分，CSS中的overflow使得当内容大小超过此范围时，会出现滚动条

    var htmlCon = createEle('div');
    htmlCon.style.cssText = "text-align:center;width:" + width + "px;height:" + (height - 50) + "px;overflow:auto;";
     if(html === ""){
    htmlCon.innerHTML = "<div><textarea id='contentArea' style='height: 95%;width: 95%;overflow:auto; font-weight: bold; color: #0000CD; font-size: 13px'></textarea></div>";
     }else{
         htmlCon.innerHTML = html;
     }
    win.appendChild(htmlCon);

    //窗体底部的按钮部分
    var btnCon = createEle('div');
    btnCon.style.cssText = "width:" + width + "px;height:35px;text-height:30px;;text-align:center;padding-top:2px;";

    //如果参数buttons为数组的话，就会创建自定义按钮
    if (isArray(buttons)) {
        var length = buttons.length;
        if (length > 0) {
            if (length % 2 === 0) {
                for (var i = 0; i < length; i = i + 2) {
                    //按钮数组
                    var btn = createEle('button');
                    btn.innerHTML = buttons[i]; //firefox不支持value属性，所以这里用innerHTML
                    //                        btn.value = buttons[i];
                    btn.onclick = buttons[i + 1];
                    btnCon.appendChild(btn);
                    //用户填充按钮之间的空白
                    var nbsp = createEle('label');
                    nbsp.innerHTML = "&nbsp&nbsp";
                    btnCon.appendChild(nbsp);
                }
            }
        }
    }
     //这是默认的确认按钮
    var obtn = createEle('button');
    //        obtn.setAttribute("value","确定");
    obtn.innerHTML = "确定";
    //        obtn.value = '确定';
    obtn.style.cssText = "margin:5px 30px 30px 5px; font-weight: bold; font-size: 15px; color: #4B4B4B; border: 1px solid #CCC; width: 60px";
    btnCon.appendChild(obtn);
    //这是默认的关闭按钮
    var cbtn = createEle('button');
    //        cbtn.setAttribute("value","关闭");
    cbtn.innerHTML = "关闭";
    //        cbtn.value = '关闭';
    cbtn.style.cssText = "margin:5px 30px 30px 5px; font-weight: bold; font-size: 15px; color: #4B4B4B; border: 1px solid #CCC; width: 60px";
    btnCon.appendChild(cbtn);

    win.appendChild(btnCon);
    appChild(win);

    /******************************************************以下为拖动窗体事件************************************************/
    //鼠标停留的X坐标
    var mouseX = 0;
    //鼠标停留的Y坐标
    var mouseY = 0;
    //窗体到body顶部的距离
    var toTop = 0;
    //窗体到body左边的距离
    var toLeft = 0;
    //判断窗体是否可以移动
    var moveable = false;

    //标题栏的按下鼠标事件
    tBar.onmousedown = function() {
        var eve = getEvent();
        moveable = true;
        mouseX = eve.clientX;
        mouseY = eve.clientY;
        toTop = parseInt(win.style.top);
        toLeft = parseInt(win.style.left);

        //移动鼠标事件

        tBar.onmousemove = function()　　 {
            if (moveable)　　 {
                var eve = getEvent();
                var x = toLeft + eve.clientX - mouseX;
                var y = toTop + eve.clientY - mouseY;
                if (x > 0 && (x + width < w) && y > 0 && (y + height < h))　　 {
                    win.style.left = x + "px";
                    win.style.top = y + "px";
                }
            }
        };
        //放下鼠标事件，注意这里是document而不是tBar

        document.onmouseup = function()　　 {
            moveable = false;
        };
    };

    //获取浏览器事件的方法，兼容ie和firefox
    function getEvent() {
        return window.event || arguments.callee.caller.arguments[0];
    }

    //顶部的标题栏和底部的按钮栏中的“关闭按钮”的关闭事件
    obtn.onclick = cbtn.onclick = closeCon.onclick = function(event) {
        event = event ? event : window.event;
        var obj = event.srcElement ? event.srcElement : event.target;
        var content = document.getElementById('contentArea').value;
        remChild(win);
        if (mask) {
            remChild(mask);
        }
        if(obj.innerHTML == '确定'){
            defaultAction(args, content);
        }
    };
}
/*
//未使用功能
function show() {
    var str = "<div><textarea id='contentArea' cols ='60' rows = '8'></textarea></div>";
    showWindow('我的提示框', "", doit, 450, 60, true, ['确定', fun1]);
}

function getStr(elename) {
    var str = document.getElementById(elename).value;
    return str;
}

function fun1() {
    content = getStr('contentArea');
    //doit(content);
    //alert('你输入的内容为：' + str);
}
*/

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

//获取上周开始日期
function getLastWeekStartDate() {
    var lastWeekStartDate = new Date(nowYear, nowMonth, nowDay - nowDayOfWeek - 7);
    return formatDate(lastWeekStartDate);
}
//获取上周结束日期
function getLastWeekEndDate() {
    var lastWeekEndDate = new Date(nowYear, nowMonth, nowDay + (6 - nowDayOfWeek) - 7);
    return formatDate(lastWeekEndDate);
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
function doit(type, content){
    var viid = '0';
    var vstartDate = '';
    var vendDate = '';
    if(type === '0'){
        vstartDate = getWeekStartDate();
        vendDate = getWeekEndDate();
    }else{
        vstartDate = getLastWeekStartDate();
        vendDate = getLastWeekEndDate();
    }
    var dates = getAll(vstartDate, vendDate);
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
    //var content=prompt("请输入每日工作内容(以|分隔!)","");
    var contents = content.replace(/；/g, ';').split(';');
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
                console.log($.trim(data));
            }
                  );
        }
    });
    window.location.reload();
    /*
    if(type === '1'){
        $("#rightFrame").contents().find("tr.x-toolbar-left-row").find("div#ext-comp-1004").find("img")[0].click();
    }
    */
}
//生成自动报工按钮
$(document).ready(function(){
　'use strict';
    var $octd=$('<td class="x-toolbar-cell" ><table id="oneclickBtn" cellspacing="0" class="x-btn   x-btn-text-icon" style="width: auto;"><tbody class="x-btn-small x-btn-icon-small-left"><tr><td class="x-btn-tl"><i>&nbsp;</i></td><td class="x-btn-tc"></td><td class="x-btn-tr"><i>&nbsp;</i></td></tr><tr><td class="x-btn-ml"><i>&nbsp;</i></td><td class="x-btn-mc"><em class="" unselectable="on"><button type="button" id="ext-octd" class=" x-btn-text btn_icon_005">一键报工</button></em></td><td class="x-btn-mr"><i>&nbsp;</i></td></tr><tr><td class="x-btn-bl"><i>&nbsp;</i></td><td class="x-btn-bc"></td><td class="x-btn-br"><i>&nbsp;</i></td></tr></tbody></table></td>');
    var $poctd=$('<td class="x-toolbar-cell" ><table id="poneclickBtn" cellspacing="0" class="x-btn   x-btn-text-icon" style="width: auto;"><tbody class="x-btn-small x-btn-icon-small-left"><tr><td class="x-btn-tl"><i>&nbsp;</i></td><td class="x-btn-tc"></td><td class="x-btn-tr"><i>&nbsp;</i></td></tr><tr><td class="x-btn-ml"><i>&nbsp;</i></td><td class="x-btn-mc"><em class="" unselectable="on"><button type="button" id="ext-poctd" class=" x-btn-text btn_icon_005">一键补报</button></em></td><td class="x-btn-mr"><i>&nbsp;</i></td></tr><tr><td class="x-btn-bl"><i>&nbsp;</i></td><td class="x-btn-bc"></td><td class="x-btn-br"><i>&nbsp;</i></td></tr></tbody></table></td>');
    document.getElementById("rightFrame").onload = function() {
        $("#rightFrame").contents().find("tr.x-toolbar-left-row").append($octd);
        $("#rightFrame").contents().find("tr.x-toolbar-left-row").append($poctd);
        //$("#rightFrame").contents().find("tr.x-toolbar-left-row").on('mouseover', '#oneclickBtn', function(){
        //$("#rightFrame").contents().find("'#oneclickBtn").attr('class', 'x-btn   x-btn-text-icon  x-btn-over');
        //});
        $("#rightFrame").contents().find("tr.x-toolbar-left-row").on('click', '#ext-octd', function(){
             showWindow('报工 <--> 输入每日工作内容(以分号分隔)', "", doit, '0', 450, 60, true);
        });
        $("#rightFrame").contents().find("tr.x-toolbar-left-row").on('click', '#ext-poctd', function(){
             showWindow('补报 <--> 输入每日工作内容(以分号分隔)', "", doit, '1', 450, 60, true);
        });
    };
});

