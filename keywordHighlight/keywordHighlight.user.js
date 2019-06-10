// ==UserScript==
// @name 	关键词高亮
// @namespace http://tampermonkey.net/
// @version 0.1
// @description 	关键词高亮
// @author xshrim
// @include      http*://*/*
// @require      http://libs.baidu.com/jquery/2.0.0/jquery.min.js
// @grant        GM_addStyle
// ==/UserScript==

GM_addStyle("qs-highlight {background: yellow; color: black; font-style: normal; box-shadow: rgba(0, 0, 0, 0.298039) 1px 1px 3px; text-decoration: none; border-radius: 3px;}");


const DELAY = 0, MAX_LENGTH = 50;



let timer, lastText;
$('body').on('mouseup dbclick', e => {

	clearTimeout(timer);
	timer = setTimeout(() => {

		let selection = getSelection();
		let text = $.trim(selection.toString());

		if (lastText == text) { // do nothing if selection is the same as lastText
			return;
		}

		undoHighlight();
		lastText = null;

		if (text.length <= MAX_LENGTH && text.length > 0) {
			lastText = text;
			doHighight(text, [selection.anchorNode, selection.baseNode]); //ignore anchorNode and baseNode
		}
	}, DELAY);
});



function doHighight(selectedText, ignores = []) {
	let regex;
	if (/[^a-zA-z0-9_\$]/.test(selectedText)) { // 满足这个条件的不用\b, 针对非变量名的情况
		regex = new RegExp(`${escRegExp(selectedText)}`, `g`);
	} else { // 否则加上\b匹配, 要求匹配整个单词
		regex = new RegExp(`\\b${escRegExp(selectedText)}\\b`, `g`);
	}



	let textList = getVisibleTextNodesIn('body');
	textList.each((index, textNode) => {
		// 递归调用, 每次进处理第一个词, 将剩余的词放入下一次递归中处理
		(function handlerFirstTextInStr(textNode) {
			let curText = textNode.nodeValue;
			let parentNode = textNode.parentNode;
			regex.lastIndex = 0;
			if (regex.test(curText) && ignores.indexOf(textNode) == -1) {
				// curText => `textLeft_selectedText_textRemain`
				// one node => three nodes
				let [textLeft = '', ...textRemain] = curText.split(selectedText);
				textNode.nodeValue = textRemain.join(selectedText);
				parentNode.insertBefore(document.createTextNode(textLeft), textNode);
				parentNode.insertBefore($(`<qs-highlight>${selectedText}</qs-highlight>`)[0], textNode);
				handlerFirstTextInStr(textNode);
			}
		})(textNode);
	});
}



function undoHighlight() {
	$('qs-highlight').each((index, el) => {
		el.outerHTML = el.innerHTML;
	});
}



function getVisibleTextNodesIn(el) {
	return $(el).find(":not(iframe, textarea):visible").addBack().contents().filter(function() {
		return this.nodeType == 3;
	});
}


function escRegExp(str) {
	return String(str).replace(/[.*+?^=!:${}()|[\]\/\\]/g, '\\$&');
}