// ==UserScript==
// @name           预览
// @namespace      http://tampermonkey.net/
// @version        0.1
// @description    detect code and view the cover
// @author         xshrim
// @include        *
// @grant          GM_setValue
// @grant          GM_getValue
// @grant          GM_setClipboard
// @grant          unsafeWindow
// @grant          window.close
// @grant          window.focus
// @grant          GM_log
// @grant          GM_addStyle
// @grant          GM_xmlhttpRequest
// @grant          GM_getResourceText
// @require        http://code.jquery.com/jquery-2.1.1.min.js

// ==/UserScript==

var $ = unsafeWindow.jQuery;

// cookie
function getCookie(name) {
  var nameEQ = name + "=";
  var ca = document.cookie.split(';');
  for(var i=0;i < ca.length;i++) {
    var c = ca[i];
    while (c.charAt(0)==' ') c = c.substring(1,c.length);
    if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
  }
  return null;
}

function setCookie(c_name, value, expiredays) {
    var exdate = new Date();
    exdate.setDate(exdate.getDate()+expiredays);
    document.cookie = c_name + "=" + escape(value) + ((expiredays==null) ?
        "" :
        ";expires="+exdate.toUTCString() + ";path=/");
}

function delete_cookie( name ) {
      document.cookie = name + '=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}

function getVideoCode(title){
    var t = title.match(/[A-Za-z]+\-\d+/);
    if(!t){
        t = title.match(/heyzo[\-\_]?\d{4}/);
    }
    if(!t){
        t = title.match(/\d{6}[\-\_]\d{3}/);
    }
    if(!t){
        t = title.match(/[A-Za-z]+\d+/);
    }
    return t;
}

function getVideoInfo(id){
    //pop.innerHTML = "";
    var info = "<div id='" + id + "' class='item_info'></div>";
    $(".imgpop")[0].innerHTML = info;
    GM_xmlhttpRequest({
        method: "GET",
        url: "https://avmoo.asia/cn/search/" + id,
        onload: xhr => {
            var xhr_data = $(xhr.responseText);
            if(!(xhr_data.find("div.alert").length)){
                var title = xhr_data.find("div.photo-info span").html();
                var info_html ="<div class='item_border'><h4>" + title + "</h4></div>";
                $(".imgpop").append(info_html);
                var img = xhr_data.find("div.photo-frame img").attr("src").replace("ps.j","pl.j");;
                $(".imgpop").append("<img src='" + img + "'>");
            }else{
                getUncensored(id);
            }
        }
    })
}

function getUncensored(id){
    GM_xmlhttpRequest({
        method: "GET",
        url: "https://avsox.asia/cn/search/" + id,
        onload: xhr => {
            var xhr_data = $(xhr.responseText);
            if(!(xhr_data.find("div.alert").length)){
                var title = xhr_data.find("div.photo-info span").html();
                var info_html ="<div class='item_border'><h4>" + title + "</h4></div>";
                $(".imgpop").append(info_html);
                var details_url = xhr_data.find("a.movie-box").attr("href");
                GM_xmlhttpRequest({
                    method: "GET",
                    url: details_url,
                    onload: temp => {
                        var img = $(temp.responseText).find("a.bigImage").attr("href");
                        $(".imgpop").find(".item_border").append("<img src='" + img + "'>");
                    }
                });
            } else {
                $(".imgpop").hide();
                $(".imgpop")[0].innerHTML = "";
            }
        }
    })
}

function showdiv(e) {
    //var e = event || window.event;
    var scrollX = document.documentElement.scrollLeft || document.body.scrollLeft;
    var scrollY = document.documentElement.scrollTop || document.body.scrollTop;
    var x = e.pageX || e.clientX + scrollX;
    var y = e.pageY || e.clientY + scrollY;
    //alert('x: ' + x + '\ny: ' + y);
    console.log({'x': x, 'y': y});
    if ( $(".imgpop").length <= 0 ) {
        //pop = document.createElement("iframe");
        var pop = document.createElement("div");
        //pop.attr("id", "imgpop");
        $(pop).attr("class", "imgpop");
        document.body.appendChild(pop);
    }
    //$(".imgpop")[0].style.cssText = "position:absolute;left:" + x + "px;top:" + y + "px;background:#f0f0f0;z-index:101;width:" + 800 + "px;height:" + 600 + "px;border:solid 2px #afccfe;";
    $(".imgpop")[0].style.cssText = "position:absolute;left:" + x + "px;top:" + y + "px;background:#f0f0f0;z-index:101;px;border:solid 2px #afccfe;";
    $(".imgpop")[0].innerHTML = "";
    $(".imgpop").hide();
}

(function() {
    'use strict';
    // show();
    // Your code here...
    var pretarget;

    document.body.onmouseover = function(event){
       var el = event.target;//鼠标每经过一个元素，就把该元素赋值给变量el
       console.log('当前鼠标在', el, '元素上');//在控制台中打印该变量
       if (el === pretarget) { // 鼠标移动但元素不变时无动作
           return
       }
       showdiv(event);

       //console.log(el.children.length);
       if (el.children.length <= 2) {
           //console.log(el.innerText);
           var code = getVideoCode(el.innerText);
           if (code !== null && code !== undefined) {
               console.log("================================code:", code[0]);
               getVideoInfo(code[0]);
               $(".imgpop").show();
               //$(".imgpop").focus();
           } else {
               $(".imgpop").hide();
               $(".imgpop")[0].innerHTML = "";
           }
       }
    }

    item.mouseenter(
            function(f){
                item_name = $(this).attr("title");
                var id = getVideoCode(item_name);
                if(id){
                    if(!(item_list.find("div#"+id).length))
                        getVideoInfo(id);
                    showVideoInfo(item_list.find("div#"+id),f.clientX,f.clientY);
                }
            }
        );
        item.mouseleave(
            function(){
                item_name = $(this).attr("title");
                hiddenVideoInfo(getVideoCode(item_name));
            }
        );
})();
