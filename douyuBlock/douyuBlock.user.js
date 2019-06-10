// ==UserScript==
// @name 	斗鱼屏蔽
// @namespace http://tampermonkey.net/
// @version 0.2
// @description 	斗鱼屏蔽指定主播！
// @author jswh
// @match *://*.douyu.com/*
// @Grant none
// ==/UserScript==

var blocklist = ['xiao8', 'yyf', '820'];

function children(curEle,tagName){
    var nodeList = curEle.children;
    var ary = [];
    if(/MSIE(6|7|8)/.test(navigator.userAgent)){
        for(var i=0;i<nodeList.length;i++){
            var curNode = nodeList[i];
            if(curNode.nodeType ===1){
                ary[ary.length] = curNode;
            }
        }
    }else{
        ary = Array.prototype.slice.call(curEle.children);
    }
    // 获取指定子元素
    if(typeof tagName === "string"){
        for(var k=0;k<ary.length;k++){
            curTag = ary[k];
            if(curTag.nodeName.toLowerCase() !== tagName.toLowerCase()){
                ary.splice(k,1);
                k--;
            }
        }
    }

    return ary;
}

(function() {
    'use strict';
    var vlist = document.getElementById('live-list-contentbox');
    var vitem = children(vlist, 'li');
    for(var i = 0; i< vitem.length; i++){
        for(var k = 0; k< blocklist.length; k++){
            if(vitem[i].children[0].children[1].children[1].children[0].innerHTML.indexOf(blocklist[k]) >= 0){
                vlist.removeChild(vitem[i]);
            }
        }
    }

})();