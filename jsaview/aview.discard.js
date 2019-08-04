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

//var $ = unsafeWindow.jQuery;
var $ = window.jQuery;

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

function mousePosition(ev){
    ev = ev || window.event;
    if(ev.pageX || ev.pageY){
        return {x:ev.pageX, y:ev.pageY};
    }
    return {
        x:ev.clientX + document.body.scrollLeft - document.body.clientLeft,
        y:ev.clientY + document.body.scrollTop - document.body.clientTop
    };
}

function getVideoCode(title){
    /*
    var t = title.match(/[A-Za-z]+\-\d+/g);
    if(!t){
        t = title.match(/heyzo[\-\_]?\d{4}/g);
    }
    if(!t){
        t = title.match(/\d{6}[\-\_]\d{3}/g);
    }
    if(!t){
        t = title.match(/[A-Za-z]+\d+/g);
    }
    return t;
    */
    return title.match(/([A-Za-z0-9]+[\-\_]\d+)|(heyzo[\-\_]?\d{4})|(\d{6}[\-\_]\d{3})|([A-Za-z]+\d+)|(\d{5}[\-\_]\d{4})|(\d{5}[\-\_]\d{3})/g);
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
                if (title !== undefined) {
                    var info_html ="<div class='imgpop_item_border'><h3>" + title + "</h3></div>";
                    if ($(".imgpop_item_border").length > 0) {
                        $(".imgpop_item_border").remove();
                    }
                    $(".imgpop").append(info_html);
                }
                var img_url = xhr_data.find("div.photo-frame img").attr("src");
                if (img_url !== undefined) {
                    if ($(".imgpop img").length > 0) {
                        $(".imgpop img").remove();
                    }
                    $(".imgpop").append("<img src='" + img_url.replace("ps.j","pl.j") + "'>");
                }
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
                if (title !== undefined) {
                    var info_html ="<div class='imgpop_item_border'><h3>" + title + "</h3></div>";
                    $(".imgpop").append(info_html);
                }
                var details_url = xhr_data.find("a.movie-box").attr("href");
                if (details_url !== undefined) {
                    GM_xmlhttpRequest({
                        method: "GET",
                        url: details_url,
                        onload: temp => {

                            var img = $(temp.responseText).find("a.bigImage").attr("href");
                            //$(".imgpop").find(".imgpop_item_border").append("<img src='" + img + "'>");
                            $(".imgpop").append("<img src='" + img + "'>");
                        }
                    });
                }
            } else {
                //$(".imgpop").hide();
                //$(".imgpop")[0].innerHTML = "";
            }
        }
    })
}

function showdiv(e) {
    //var e = event || window.event;

    /* 相对于页面的坐标
    var scrollX = document.documentElement.scrollLeft || document.body.scrollLeft;
    var scrollY = document.documentElement.scrollTop || document.body.scrollTop;
    var x = e.pageX || e.clientX + scrollX;
    var y = e.pageY || e.clientY + scrollY;
    */

    /*
    var w =  document.body.clientWidth;
    var h = document.body.clientHeight;
    var x = e.clientX ;
    var y = e.clientY;
    */

    //console.log({'x': x, 'y': y}, {'sx': w, 'sy': h});

    var mousePos = mousePosition(e)


    if ( $(".imgpop").length <= 0 ) {
        //pop = document.createElement("iframe");
        var pop = document.createElement("div");
        //pop.attr("id", "imgpop");
        $(pop).attr("class", "imgpop");

        pop.style.cssText = "position:absolute;left:" + mousePos.x + "px;top:" + mousePos.y + "px;background:#f0f0f0;z-index:101;border:solid 2px #afccfe;cursor:move;";

        var ismousedown = false;
		var popleft,poptop;
		var downX,downY;

        popleft = parseInt(pop.style.left);
	    poptop = parseInt(pop.style.top);

        pop.onmousedown = function(e){
            if (!e.ctrlKey){
			    ismousedown = true;
                pop.style.cursor="move";
            } else {
                ismousedown = false;
                pop.style.cursor="text";
            }
			downX = e.clientX;
			downY = e.clientY;
            //console.log(downX, downY);
		}
        pop.onmousemove = function(e){
				if(ismousedown){
				pop.style.top = e.clientY - downY + poptop + "px";
				pop.style.left = e.clientX - downX + popleft + "px";
				}
		}
        document.onmouseup = function(){
				popleft = parseInt(pop.style.left);
				poptop = parseInt(pop.style.top);
				ismousedown = false;
			}
        document.body.appendChild(pop);
    } else {
        console.log($(".imgpop").is(":visible"));
    }
    //$(".imgpop")[0].style.cssText = "position:absolute;left:" + x + "px;top:" + y + "px;background:#f0f0f0;z-index:101;width:" + 800 + "px;height:" + 600 + "px;border:solid 2px #afccfe;";
    //$(".imgpop")[0].style.cssText = "position:absolute;left:" + mousePos.x + "px;top:" + mousePos.y + "px;background:#f0f0f0;z-index:101;border:solid 2px #afccfe;cursor:move;";
/*
    if (y > h / 2) {
        if (x > w / 2) {
            $(".imgpop")[0].style.cssText = "position:absolute;right:" + (w - x) + "px;bottom:" + (h - y) + "px;background:#f0f0f0;z-index:101;px;border:solid 2px #afccfe;";
        }else{
            $(".imgpop")[0].style.cssText = "position:absolute;left:" + x + "px;bottom:" + (h - y) + "px;background:#f0f0f0;z-index:101;px;border:solid 2px #afccfe;";
        }
    } else {
        if (x > w / 2) {
            $(".imgpop")[0].style.cssText = "position:absolute;right:" + (w - x) + "px;top:" + y + "px;background:#f0f0f0;z-index:101;px;border:solid 2px #afccfe;";
        } else {
            $(".imgpop")[0].style.cssText = "position:absolute;left:" + x + "px;top:" + y + "px;background:#f0f0f0;z-index:101;px;border:solid 2px #afccfe;";
        }
    }
*/

    $(".imgpop").hide();
    //$(".imgpop")[0].innerHTML = "";
    $(".imgpop").empty();
}

(function() {
    'use strict';
    // show();
    // Your code here...
    var pretarget;
    //var lasty=0;

    document.body.onclick = function(event){
        if (event.target.className !== "imgpop" && event.target.className !== "imgpop_item_border" && event.target.parentNode.className !== "imgpop" && event.target.parentNode.className !== "imgpop_item_border") {
            if ($(".imgpop").length > 0 && $(".imgpop").is(":visible")) {
                $(".imgpop").hide();
                //$(".imgpop")[0].innerHTML = "";
                $(".imgpop").empty();
            }
        }
    }

    document.body.onmouseover = function(event){
        if (!event.altKey) {
            return
        }

       if (event.target === pretarget) { // 鼠标移动但元素不变时无动作
           return
       }
        //console.log('当前鼠标在', el, '元素上');//在控制台中打印该变量

       //if (Math.abs(event.clientY-lasty) < 10) {
       //    return
       //}
       //lasty = event.clientY;

       showdiv(event);

       //console.log(el.children.length);
       if (event.target.children.length <= 1) {
           //console.log(el.innerText);
           var code = getVideoCode(event.target.innerText);
           if (code !== null && code !== undefined) {
               $.each(code ,function(index,value){
                   console.log("=====  ", value, "  =====");
                   getVideoInfo(value);
               });
               $(".imgpop").show();
               //$(".imgpop").focus();
           } else {
               $(".imgpop").hide();
               $(".imgpop").empty();
               //$(".imgpop")[0].innerHTML = "";
           }
       }
    }
/*
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
*/
})(); 
