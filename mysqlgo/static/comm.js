// 原生js代替jquery的方法: https://www.cnblogs.com/xiaopen/p/5540884.html
function parseParams(data) {
    // from:
    // {name: 'zhangsan', age: 100}
    // to:
    // "name=zhangsan&age=100"
    try {
        var tempArr = [];
        for (var i in data) {
            var key = encodeURIComponent(i);
            var value = encodeURIComponent(data[i]);
            tempArr.push(key + '=' + value);
        }
        var urlParamsStr = tempArr.join('&');
        return urlParamsStr;
    } catch (err) {
        return '';
    }
}

function extractParams(url) {
    // from:
    // var urlStr = 'http://www.xxx.com/test?name=zhangshan&age=100#hello';
    // to:
    // {name: "zhangshan", age: "100"}
    try {
        var index = url.indexOf('?');
        url = url.match(/\?([^#]+)/)[1];
        var obj = {}, arr = url.split('&');
        for (var i = 0; i < arr.length; i++) {
            var subArr = arr[i].split('=');
            var key = decodeURIComponent(subArr[0]);
            var value = decodeURIComponent(subArr[1]);
            obj[key] = value;
        }
        return obj;

    } catch (err) {
        return null;
    }
}

function getParams() {
    // var isql = document.querySelector("#code").value;
    var isql = editor.getValue();  // textarea已被codemirror接管
    var ikind = document.querySelector("#kind").value;
    var ihost = document.querySelector("#host").value;
    var iport = document.querySelector("#port").value;
    var icharset = document.querySelector("#charset").value;
    var iuser = document.querySelector("#user").value;
    var ipasswd = document.querySelector("#passwd").value;

    var idbname;
    var dbselect = document.querySelector("#dbname");
    if (dbselect.selectedIndex < 0) {
        idbname = "";
    } else {
        idbname = dbselect.options[dbselect.selectedIndex].text;
    }

    var itbname;
    var tbselect = document.querySelector("#tbname");
    if (tbselect.selectedIndex < 0) {
        itbname = "";
    } else {
        itbname = tbselect.options[tbselect.selectedIndex].text;
    }

    var params = {  // 构建json格式传入数据, 然后使用parseParams函数将数据转换为js原生ajax支持的格式
        kind: ikind,
        host: ihost,
        port: iport,
        charset: icharset,
        user: iuser,
        passwd: ipasswd,
        dbname: idbname,
        tbname: itbname,
        sql: isql
    }

    return params
}

function selectTheme() {
    var input = document.getElementById("select");
    var theme = input.options[input.selectedIndex].textContent;
    console.log(theme);
    editor.setOption("theme", theme);
    //location.hash = "#" + theme;
}

function removeDiv(el) {
    //var resbed = document.getElementById("resdiv");

    while (el.firstChild) el.removeChild(el.firstChild);  // 删除所有子元素table
}

function createDiv(el, title, headers, datas) {
    //var resbed = document.getElementById("resdiv");

    var div = document.createElement("div");
    el.appendChild(div);

    // 创建表
    var table = document.createElement("table");
    //将表格追加到父容器中
    div.appendChild(table);
    //设置table的样式
    table.id = "table-3";
    table.cellSpacing = 0;
    table.cellPadding = 0;
    table.border = "1px";

    // 创建标题
    if (title !== "") {
        var caption = document.createElement("caption")
        table.appendChild(caption);
        caption.innerHTML = "<strong>" + title + "</strong>";
    }

    // 创建表头
    var thead = document.createElement("thead");
    //将标题追加到table中
    table.appendChild(thead);
    // 创建tr
    var tr = document.createElement("tr");
    // 将tr追加到thead中
    // 设置tr的样式属性
    tr.style.height = "30px";
    tr.style.backgroundColor = "lightgray";
    thead.appendChild(tr);
    // 遍历headers中的数据
    for (var i = 0; i < headers.length; i++) {
        // 创建th
        var th = document.createElement("th");
        th.innerHTML = headers[i];
        // 将th追加到thead中的tr中
        tr.appendChild(th);
    }

    // 创建表行
    for (var idx in datas) {
        var tr = document.createElement("tr");
        // 将tr追加到thead中
        // 设置tr的样式属性
        tr.style.height = "30px";
        //tr.style.backgroundColor = "lightgray";
        thead.appendChild(tr);
        for (var i in headers) {
            var td = document.createElement("td");
            td.innerHTML = datas[idx][headers[i]];
            // 将th追加到thead中的tr中
            tr.appendChild(td);
        }
    }

    // 创建分隔线
    var sep = document.createElement("hr");
    // sep.style.width = "100%";
    // sep.style.color = "#987cb9";
    sep.className = "style-two";
    el.appendChild(sep);
}

function tab(el) {
    var lis = document.querySelectorAll("#head ul li");
    var divs = document.querySelectorAll("#content > div")
    //editor.focus();
    //console.log(el.textContent)
    for (var idx in lis) {
        if (lis[idx] == el) {
            lis[idx].className = "working";
            divs[idx].className = "";
            doTabDiv(lis[idx].textContent, divs[idx]);
        } else {
            lis[idx].className = "";
            divs[idx].className = "hide";
        }
    }
}

function doTabDiv(title, el) {
    console.log(title, el);
    if (title == "概览") {
        var params = getParams();

        params.sql = "select * from performance_schema.global_status";
        var el = document.querySelector("#viewdiv");
        var loading = document.querySelector("#loading");
        loading.style.display = "";
        ajax_method("./info", parseParams(params), "post", function (data) {
            loading.style.display = "none";
            if (data !== undefined) {
                console.log(data);
                removeDiv(el);
                // var stat = params.sql;
                jobj = JSON.parse(data);

                var stat = "";
                var rows;
                var fields = new Array();

                if (Array.isArray(jobj)) {  // 结果转为JSONArray
                    rows = jobj;
                } else {
                    rows = new Array();
                    rows.push(jobj);
                }

                for (var key in rows[0]) {  // 从第一条记录中获取key列表
                    fields.push(key);
                }

                createDiv(el, stat, fields, rows);

            }
        });
    } else if (title == "结构") {
        var params = getParams();
        if (params.tbname == "") {
            return
        }
        params.sql = "desc " + params.tbname;
        var el = document.querySelector("#structdiv");
        var loading = document.querySelector("#loading");
        loading.style.display = "";
        ajax_method("./info", parseParams(params), "post", function (data) {
            loading.style.display = "none";
            if (data !== undefined) {
                console.log(data);
                removeDiv(el);
                // var stat = params.sql;
                jobj = JSON.parse(data);

                var stat = "";
                var rows;
                var fields = new Array();

                if (Array.isArray(jobj)) {  // 结果转为JSONArray
                    rows = jobj;
                } else {
                    rows = new Array();
                    rows.push(jobj);
                }

                for (var key in rows[0]) {  // 从第一条记录中获取key列表
                    fields.push(key);
                }

                createDiv(el, stat, fields, rows);

            }
        });

    } else if (title == "数据") {
        var params = getParams();
        if (params.tbname == "") {
            return
        }
        params.sql = "select * from " + params.tbname;
        var el = document.querySelector("#datadiv");
        var loading = document.querySelector("#loading");
        loading.style.display = "";
        ajax_method("./info", parseParams(params), "post", function (data) {
            loading.style.display = "none";
            if (data !== undefined) {
                console.log(data);
                removeDiv(el);
                // var stat = params.sql;
                jobj = JSON.parse(data);

                var stat = "";
                var rows;
                var fields = new Array();

                if (Array.isArray(jobj)) {  // 结果转为JSONArray
                    rows = jobj;
                } else {
                    rows = new Array();
                    rows.push(jobj);
                }

                for (var key in rows[0]) {  // 从第一条记录中获取key列表
                    fields.push(key);
                }

                createDiv(el, stat, fields, rows);

            }
        });

    } else if (title == "执行") {
        editor.focus();
        editor.refresh();
    }
}

function ajax_method(url, data, method, success) {
    // 异步对象
    var ajax = new XMLHttpRequest();

    // get 跟post  需要分别写不同的代码
    if (method == 'get') {
        // get请求
        if (data) {
            // 如果有值
            url += '?';
            url += data;
        } else {

        }
        // 设置 方法 以及 url
        ajax.open(method, url);

        // send即可
        ajax.send();
    } else {
        // post请求
        // post请求 url 是不需要改变
        ajax.open(method, url);

        // 需要设置请求报文
        ajax.setRequestHeader("Content-type", "application/x-www-form-urlencoded");

        // 判断data send发送数据
        if (data) {
            // 如果有值 从send发送
            ajax.send(data);
        } else {
            // 木有值 直接发送即可
            ajax.send();
        }
    }

    // 注册事件
    ajax.onreadystatechange = function () {
        // 在事件中 获取数据 并修改界面显示
        if (ajax.readyState == 4 && ajax.status == 200) {
            // console.log(ajax.responseText);

            // 将 数据 让 外面可以使用
            // return ajax.responseText;

            // 当 onreadystatechange 调用时 说明 数据回来了
            // ajax.responseText;

            // 如果说 外面可以传入一个 function 作为参数 success
            success(ajax.responseText);
        }
    }
}

window.onload = function () {
    // 初始化代码编辑器
    var mime = 'text/x-mariadb';
    // get mime type
    if (window.location.href.indexOf('mime=') > -1) {
        mime = window.location.href.substr(window.location.href.indexOf('mime=') + 5);
    }
    window.editor = CodeMirror.fromTextArea(document.getElementById('code'), {
        mode: mime,
        //theme: "liquibyte",
        indentWithTabs: true,
        smartIndent: true,
        //lineNumbers: true,
        matchBrackets: true,
        autofocus: true,
        extraKeys: { "Ctrl-Space": "autocomplete" },
        hintOptions: {
            tables: {
                users: ["name", "score", "birthDate"],
                countries: ["name", "population", "size"]
            }
        }
    });
    editor.setSize(null, 200);
    //editor.setOption("theme", "liquibyte");

    var connbtn = document.querySelector("#conn");
    connbtn.onclick = function () {
        var loading = document.querySelector("#loading");
        loading.style.display = "";
        console.log(loading);
        var params = getParams();
        params.sql = "show databases";
        ajax_method("./info", parseParams(params), "post", function (data) {
            loading.style.display = "none";
            if (data !== undefined) {
                console.log(data);
                var dbname = document.querySelector("#dbname");
                while (dbname.firstChild) dbname.removeChild(dbname.firstChild);
                var dopt = document.createElement("option");
                dopt.innerHTML = "";
                dbname.appendChild(dopt);
                if (data.indexOf("ERROR:") != 0) {
                    document.querySelector("#connflag").style.color = "green";
                    jdata = JSON.parse(data);
                    for (var idx in jdata) {
                        var opt = document.createElement("option");
                        for (var key in jdata[idx]) {
                            opt.innerHTML = jdata[idx][key];
                        }
                        dbname.appendChild(opt);
                    }
                } else {
                    document.querySelector("#connflag").style.color = "red";
                    alert(data);
                }
            }
        });

        tab(document.querySelector("#head ul li"));
    }

    var dbselect = document.querySelector("#dbname");
    dbselect.onchange = function () {
        var params = getParams();
        params.sql = "show tables";
        var loading = document.querySelector("#loading");
        loading.style.display = "";
        ajax_method("./info", parseParams(params), "post", function (data) {
            loading.style.display = "none";
            if (data !== undefined) {
                console.log(data);
                var tbname = document.querySelector("#tbname");
                while (tbname.firstChild) tbname.removeChild(tbname.firstChild);
                // var dopt = document.createElement("option");
                // dopt.innerHTML = "";
                // dbname.appendChild(dopt);
                if (data.indexOf("ERROR:") != 0) {
                    //document.querySelector("#connflag").style.color = "green";
                    jdata = JSON.parse(data);
                    for (var idx in jdata) {
                        var opt = document.createElement("option");
                        for (var key in jdata[idx]) {
                            opt.innerHTML = jdata[idx][key];
                        }
                        tbname.appendChild(opt);
                    }
                } else {
                    //document.querySelector("#connflag").style.color = "red";
                    //alert(data);
                }
            }
        });
    }

    var rstbtn = document.querySelector("#reset");
    rstbtn.onclick = function () {
        //console.log(editor.getValue());  //经过转义的数据
        //console.log(editor.getTextArea().value); //未经转义的数据
        //editor.setValue('select * from user');
        editor.setValue('');
    }


    //do something
    var smtbtn = document.querySelector("#submit");
    smtbtn.onclick = function () {
        var params = getParams();
        var el = document.querySelector("#resdiv");
        var loading = document.querySelector("#loading");
        loading.style.display = "";
        ajax_method("./exe", parseParams(params), "post", function (data) {
            loading.style.display = "none";
            if (data !== undefined) {
                console.log(data);
                data = data.replace(/\(MISSING\)/, "");
                removeDiv(el);
                jdata = JSON.parse(data);
                for (var idx in jdata) {
                    var stat;
                    var res;
                    var item = jdata[idx];
                    if (Array.isArray(item)) {
                        stat = item[0];
                        res = item[1];
                    } else {
                        stat = "ERROR";
                        res = JSON.stringify(item);
                    }
                    // data = '{ "Database": "blockchain" }'
                    var jobj = JSON.parse(res);

                    var rows;
                    var fields = new Array();

                    if (Array.isArray(jobj)) {  // 结果转为JSONArray
                        rows = jobj;
                    } else {
                        rows = new Array();
                        rows.push(jobj);
                    }

                    for (var key in rows[0]) {  // 从第一条记录中获取key列表
                        fields.push(key);
                    }

                    createDiv(el, stat, fields, rows);

                }
            }
        });
    }
}