// ==UserScript==
// @name        steam卡牌利润最大化
// @namespace   https://github.com/lzghzr/GreasemonkeyJS
// @version     0.2.22
// @author      lzghzr
// @description 按照美元区出价, 最大化steam卡牌卖出的利润
// @supportURL  https://github.com/lzghzr/GreasemonkeyJS/issues
// @include     /^https?:\/\/steamcommunity\.com\/.*\/inventory/
// @license     MIT
// @grant       GM_xmlhttpRequest
// @run-at      document-end
// ==/UserScript==
"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = y[op[0] & 2 ? "return" : op[0] ? "throw" : "next"]) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [0, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
var SteamCardMaximumProfit = (function () {
    function SteamCardMaximumProfit() {
        this._D = document;
        this._W = unsafeWindow || window;
        this._divItems = [];
        this.quickSells = [];
    }
    SteamCardMaximumProfit.prototype.Start = function () {
        var _this = this;
        this._AddUI();
        this._DoLoop();
        var elmDivActiveInventoryPage = this._D.querySelector('#inventories');
        var observer = new MutationObserver(function (rec) {
            if (location.hash.match(/^#753|^$/)) {
                for (var _i = 0, rec_1 = rec; _i < rec_1.length; _i++) {
                    var r = rec_1[_i];
                    var rt = r.target;
                    if (rt.classList.contains('inventory_page')) {
                        var itemHolders = rt.querySelectorAll('.itemHolder');
                        for (var i = 0; i < itemHolders.length; i++) {
                            var rgItem = _this._GetRgItem(itemHolders[i]);
                            if (rgItem != null && _this._divItems.indexOf(rgItem.element) === -1 && rgItem.description.appid === 753 && rgItem.description.marketable === 1) {
                                _this._divItems.push(rgItem.element);
                                var elmDiv = _this._D.createElement('div');
                                elmDiv.classList.add('scmpItemCheckbox');
                                rgItem.element.appendChild(elmDiv);
                            }
                        }
                    }
                }
            }
        });
        observer.observe(elmDivActiveInventoryPage, { childList: true, subtree: true, attributes: true, attributeFilter: ['style'] });
    };
    SteamCardMaximumProfit.prototype._AddUI = function () {
        return __awaiter(this, void 0, void 0, function () {
            var _this = this;
            var elmStyle, elmDivInventoryPageRight, elmDiv, elmSpanQuickSellItems, elmDivQuickSellItem, i, baiduExch;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        elmStyle = this._D.createElement('style');
                        elmStyle.innerHTML = "\n.scmpItemSelect {\n  background: yellow;\n}\n.scmpItemRun {\n  background: blue;\n}\n.scmpItemSuccess {\n  background: green;\n}\n.scmpItemError {\n  background: red;\n}\n.scmpQuickSell {\n  margin: 0 0 1em;\n}\n.scmpItemCheckbox {\n  position: absolute;\n  z-index: 100;\n  top: 0;\n  left: 0;\n  width: 20px;\n  height: 20px;\n  border: 2px solid yellow;\n  opacity: 0.7;\n  cursor: default;\n}\n.scmpItemCheckbox:hover {\n  opacity: 1;\n}\n#scmpExch {\n  width: 5em;\n}";
                        this._D.body.appendChild(elmStyle);
                        elmDivInventoryPageRight = this._D.querySelector('.inventory_page_right'), elmDiv = this._D.createElement('div');
                        elmDiv.innerHTML = "\n<div class=\"scmpQuickSell\">\u5EFA\u8BAE\u6700\u4F4E\u552E\u4EF7:\n  <span class=\"btn_green_white_innerfade scmpQuickSellItem\">null</span>\n  <span class=\"btn_green_white_innerfade scmpQuickSellItem\">null</span>\n</div>\n<div>\n  \u6C47\u7387:\n  <input class=\"filter_search_box\" id=\"scmpExch\" type=\"text\">\n  <span class=\"btn_green_white_innerfade\" id=\"scmpQuickAllItem\">\u5FEB\u901F\u51FA\u552E</span>\n  \u5269\u4F59:\n  <span id=\"scmpQuickSurplus\">0</span>\n  \u5931\u8D25:\n  <span id=\"scmpQuickError\">0</span>\n</div>";
                        elmDivInventoryPageRight.appendChild(elmDiv);
                        elmSpanQuickSellItems = elmDiv.querySelectorAll('.scmpQuickSellItem');
                        this.spanQuickSurplus = elmDiv.querySelector('#scmpQuickSurplus');
                        this.spanQuickError = elmDiv.querySelector('#scmpQuickError');
                        this._D.addEventListener('click', function (ev) { return __awaiter(_this, void 0, void 0, function () {
                            var evt, rgItem, itemInfo, priceOverview, rgItem, select_1, ChangeClass, start, end, someDivItems, _i, someDivItems_1, y;
                            return __generator(this, function (_a) {
                                switch (_a.label) {
                                    case 0:
                                        evt = ev.target;
                                        if (!(evt.className === 'inventory_item_link')) return [3, 2];
                                        elmSpanQuickSellItems[0].innerText = 'null';
                                        elmSpanQuickSellItems[1].innerText = 'null';
                                        rgItem = this._GetRgItem(evt.parentNode), itemInfo = new ItemInfo(rgItem);
                                        return [4, this._GetPriceOverview(itemInfo)];
                                    case 1:
                                        priceOverview = _a.sent();
                                        if (priceOverview != null) {
                                            elmSpanQuickSellItems[0].innerText = priceOverview.firstFormatPrice;
                                            elmSpanQuickSellItems[1].innerText = priceOverview.secondFormatPrice;
                                        }
                                        return [3, 3];
                                    case 2:
                                        if (evt.classList.contains('scmpItemCheckbox')) {
                                            rgItem = this._GetRgItem(evt.parentNode), select_1 = evt.classList.contains('scmpItemSelect'), ChangeClass = function (elmDiv) {
                                                var elmCheckbox = elmDiv.querySelector('.scmpItemCheckbox');
                                                if (elmDiv.parentNode.style.display !== 'none' && !elmCheckbox.classList.contains('scmpItemSuccess')) {
                                                    elmCheckbox.classList.remove('scmpItemError');
                                                    elmCheckbox.classList.toggle('scmpItemSelect', !select_1);
                                                }
                                            };
                                            if (this._divLastChecked != null && ev.shiftKey) {
                                                start = this._divItems.indexOf(this._divLastChecked), end = this._divItems.indexOf(rgItem.element), someDivItems = this._divItems.slice(Math.min(start, end), Math.max(start, end) + 1);
                                                for (_i = 0, someDivItems_1 = someDivItems; _i < someDivItems_1.length; _i++) {
                                                    y = someDivItems_1[_i];
                                                    ChangeClass(y);
                                                }
                                            }
                                            else
                                                ChangeClass(rgItem.element);
                                            this._divLastChecked = rgItem.element;
                                        }
                                        _a.label = 3;
                                    case 3: return [2];
                                }
                            });
                        }); });
                        elmDivQuickSellItem = this._D.querySelectorAll('.scmpQuickSellItem');
                        for (i = 0; i < elmDivQuickSellItem.length; i++) {
                            elmDivQuickSellItem[i].addEventListener('click', function (ev) {
                                var evt = ev.target, rgItem = _this._GetRgItem(_this._D.querySelector('.activeInfo'));
                                if (!rgItem.element.querySelector('.scmpItemCheckbox').classList.contains('scmpItemSuccess') && evt.innerText != 'null') {
                                    var price = _this._W.GetPriceValueAsInt(evt.innerText), itemInfo = new ItemInfo(rgItem, price);
                                    _this._QuickSellItem(itemInfo);
                                }
                            });
                        }
                        this._D.querySelector('#scmpQuickAllItem').addEventListener('click', function () {
                            var itemInfos = _this._D.querySelectorAll('.scmpItemSelect');
                            for (var i = 0; i < itemInfos.length; i++) {
                                var rgItem = _this._GetRgItem(itemInfos[i].parentNode), itemInfo = new ItemInfo(rgItem);
                                if (rgItem.description.marketable === 1)
                                    _this.quickSells.push(itemInfo);
                            }
                        });
                        this._inputUSDCNY = elmDiv.querySelector('#scmpExch');
                        this._inputUSDCNY.value = '6.50';
                        return [4, tools.XHR({
                                method: 'GET',
                                url: "https://sp0.baidu.com/8aQDcjqpAAV3otqbppnN2DJv/api.php?query=1%E7%BE%8E%E5%85%83%E7%AD%89%E4%BA%8E%E5%A4%9A%E5%B0%91%E4%BA%BA%E6%B0%91%E5%B8%81&resource_id=6017&t=" + Date.now() + "&ie=utf8&oe=utf8&format=json&tn=baidu",
                                responseType: 'json',
                                GM_xmlhttpRequest: true
                            }).catch(console.log)];
                    case 1:
                        baiduExch = _a.sent();
                        if (baiduExch != null)
                            this._inputUSDCNY.value = baiduExch.data[0].number2;
                        return [2];
                }
            });
        });
    };
    SteamCardMaximumProfit.prototype._GetPriceOverview = function (itemInfo) {
        return __awaiter(this, void 0, void 0, function () {
            var priceoverview, stop, marketListings, marketLoadOrderSpread, itemordershistogram;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4, tools.XHR({
                            method: 'GET',
                            url: "/market/priceoverview/?country=US&currency=1&appid=" + itemInfo.rgItem.description.appid + "&market_hash_name=" + encodeURIComponent(this._W.GetMarketHashName(itemInfo.rgItem.description)),
                            responseType: 'json'
                        }).catch(console.log)];
                    case 1:
                        priceoverview = _a.sent(), stop = function () {
                            itemInfo.status = 'error';
                            return;
                        };
                        if (!(priceoverview != null && priceoverview.success && priceoverview.lowest_price !== '')) return [3, 2];
                        itemInfo.lowestPrice = priceoverview.lowest_price.replace('$', '');
                        return [2, this._CalculatePrice(itemInfo)];
                    case 2: return [4, tools.XHR({
                            method: 'GET',
                            url: "/market/listings/" + itemInfo.rgItem.description.appid + "/" + encodeURIComponent(this._W.GetMarketHashName(itemInfo.rgItem.description)),
                            responseType: 'text'
                        }).catch(console.log)];
                    case 3:
                        marketListings = _a.sent();
                        if (marketListings == null)
                            return [2, stop()];
                        marketLoadOrderSpread = marketListings.toString().match(/Market_LoadOrderSpread\( (\d+)/);
                        if (marketLoadOrderSpread == null)
                            return [2, stop()];
                        return [4, tools.XHR({
                                method: 'GET',
                                url: "/market/itemordershistogram/?country=US&language=english&currency=1&item_nameid=" + marketLoadOrderSpread[1] + "&two_factor=0",
                                responseType: 'json'
                            }).catch(console.log)];
                    case 4:
                        itemordershistogram = _a.sent();
                        if (itemordershistogram == null)
                            return [2, stop()];
                        if (itemordershistogram.success) {
                            itemInfo.lowestPrice = ' ' + itemordershistogram.sell_order_graph[0][0];
                            return [2, this._CalculatePrice(itemInfo)];
                        }
                        else
                            return [2, stop()];
                        _a.label = 5;
                    case 5: return [2];
                }
            });
        });
    };
    SteamCardMaximumProfit.prototype._CalculatePrice = function (itemInfo) {
        var firstPrice = this._W.GetPriceValueAsInt(itemInfo.lowestPrice), publisherFee = itemInfo.rgItem.description.market_fee || this._W.g_rgWalletInfo.wallet_publisher_fee_percent_default, feeInfo = this._W.CalculateFeeAmount(firstPrice, publisherFee);
        firstPrice = firstPrice - feeInfo.fees;
        itemInfo.firstPrice = Math.floor(firstPrice * parseFloat(this._inputUSDCNY.value));
        itemInfo.secondPrice = Math.floor((firstPrice + 1) * parseFloat(this._inputUSDCNY.value));
        itemInfo.firstFormatPrice = this._W.v_currencyformat(itemInfo.firstPrice, this._W.GetCurrencyCode(this._W.g_rgWalletInfo.wallet_currency));
        itemInfo.secondFormatPrice = this._W.v_currencyformat(itemInfo.secondPrice, this._W.GetCurrencyCode(this._W.g_rgWalletInfo.wallet_currency));
        return itemInfo;
    };
    SteamCardMaximumProfit.prototype._QuickSellItem = function (itemInfo) {
        return __awaiter(this, void 0, void 0, function () {
            var price, sellitem;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        itemInfo.status = 'run';
                        price = itemInfo.price || itemInfo.firstPrice;
                        return [4, tools.XHR({
                                method: 'POST',
                                url: 'https://steamcommunity.com/market/sellitem/',
                                data: "sessionid=" + this._W.g_sessionID + "&appid=" + itemInfo.rgItem.description.appid + "&contextid=" + itemInfo.rgItem.contextid + "&assetid=" + itemInfo.rgItem.assetid + "&amount=1&price=" + price,
                                responseType: 'json',
                                cookie: true
                            }).catch(console.log)];
                    case 1:
                        sellitem = _a.sent();
                        if (sellitem == null || !sellitem.success)
                            itemInfo.status = 'error';
                        else
                            itemInfo.status = 'success';
                        return [2];
                }
            });
        });
    };
    SteamCardMaximumProfit.prototype._DoLoop = function () {
        return __awaiter(this, void 0, void 0, function () {
            var _this = this;
            var itemInfo, loop, priceOverview;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        itemInfo = this.quickSells.shift(), loop = function () {
                            setTimeout(function () {
                                _this._DoLoop();
                            }, 500);
                        };
                        if (!(itemInfo != null)) return [3, 5];
                        return [4, this._GetPriceOverview(itemInfo)];
                    case 1:
                        priceOverview = _a.sent();
                        if (!(priceOverview != null)) return [3, 3];
                        return [4, this._QuickSellItem(priceOverview)];
                    case 2:
                        _a.sent();
                        this._DoLoop();
                        return [3, 4];
                    case 3:
                        loop();
                        _a.label = 4;
                    case 4: return [3, 6];
                    case 5:
                        loop();
                        _a.label = 6;
                    case 6: return [2];
                }
            });
        });
    };
    SteamCardMaximumProfit.prototype._GetRgItem = function (elmDiv) {
        return ('wrappedJSObject' in elmDiv) ? elmDiv.wrappedJSObject.rgItem : elmDiv.rgItem;
    };
    return SteamCardMaximumProfit;
}());
var ItemInfo = (function () {
    function ItemInfo(rgItem, price) {
        this.rgItem = rgItem;
        if (price != null)
            this.price = price;
    }
    Object.defineProperty(ItemInfo.prototype, "status", {
        get: function () {
            return this._status || '';
        },
        set: function (valve) {
            this._status = valve;
            var elmCheckbox = this.rgItem.element.querySelector('.scmpItemCheckbox');
            if (elmCheckbox == null)
                return;
            switch (valve) {
                case 'run':
                    elmCheckbox.classList.remove('scmpItemError');
                    elmCheckbox.classList.remove('scmpItemSelect');
                    elmCheckbox.classList.add('scmpItemRun');
                    break;
                case 'success':
                    scmp.spanQuickSurplus.innerText = scmp.quickSells.length.toString();
                    elmCheckbox.classList.remove('scmpItemError');
                    elmCheckbox.classList.remove('scmpItemRun');
                    elmCheckbox.classList.add('scmpItemSuccess');
                    break;
                case 'error':
                    scmp.spanQuickSurplus.innerText = scmp.quickSells.length.toString();
                    scmp.spanQuickError.innerText = (parseInt(scmp.spanQuickError.innerText) + 1).toString();
                    elmCheckbox.classList.remove('scmpItemRun');
                    elmCheckbox.classList.add('scmpItemError');
                    elmCheckbox.classList.add('scmpItemSelect');
                    break;
                default:
                    break;
            }
        },
        enumerable: true,
        configurable: true
    });
    return ItemInfo;
}());
var tools = (function () {
    function tools() {
    }
    tools.XHR = function (XHROptions) {
        return new Promise(function (resolve, reject) {
            if (XHROptions.GM_xmlhttpRequest) {
                GM_xmlhttpRequest({
                    method: XHROptions.method,
                    url: XHROptions.url,
                    user: XHROptions.user,
                    password: XHROptions.password,
                    responseType: XHROptions.responseType || '',
                    timeout: 3e4,
                    onload: function (res) {
                        if (res.status === 200)
                            resolve(res.response);
                        else
                            reject(res);
                    },
                    onerror: reject,
                    ontimeout: reject
                });
            }
            else {
                var xhr = new XMLHttpRequest();
                xhr.open(XHROptions.method, XHROptions.url, true, XHROptions.user, XHROptions.password);
                if (XHROptions.method === 'POST')
                    xhr.setRequestHeader('content-type', 'application/x-www-form-urlencoded; charset=utf-8');
                if (XHROptions.cookie)
                    xhr.withCredentials = true;
                xhr.responseType = XHROptions.responseType || '';
                xhr.timeout = 3e4;
                xhr.onload = function (ev) {
                    var evt = ev.target;
                    if (evt.status === 200)
                        resolve(evt.response);
                    else
                        reject(evt);
                };
                xhr.onerror = reject;
                xhr.ontimeout = reject;
                xhr.send(XHROptions.data);
            }
        });
    };
    return tools;
}());
var scmp = new SteamCardMaximumProfit();
scmp.Start();
