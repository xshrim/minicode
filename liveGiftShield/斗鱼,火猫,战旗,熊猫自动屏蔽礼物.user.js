// ==UserScript==
// @name         斗鱼,火猫,战旗,熊猫自动屏蔽礼物
// @version      0.3
// @description  网页加载完毕自动屏蔽斗鱼，熊猫礼物信息
// @author       douyu
// @include        *://www.douyu.com/*
// @include        *://www.huomao.com/*
// @include     *://www.zhanqi.tv/*
// @include     *://www.panda.tv/*
// @grant        none
// @namespace https://greasyfork.org/users/15780
// ==/UserScript==
window.onload = setTimeout(init,3000);
function  init()
{
               if(location.href.indexOf("douyu.com")!=-1)
            {
               document.getElementById("shie-switch").click()
            }else if(location.href.indexOf("panda.tv")!=-1)
            {
            	document.getElementById("gift-forbid-option-forbid_chat_gift").click();
            	document.getElementById("gift-forbid-option-forbid_flash_gift").click();
            	document.getElementById("gift-forbid-option-forbid_chat_notice").click();
            }else if(location.href.indexOf("huomao.com")!=-1)
            {
                document.getElementById("gift_fider").click();
            }else if(location.href.indexOf("zhanqi.tv")!=-1)
            {
                document.getElementById("js-gift-shield").click();
            }
}