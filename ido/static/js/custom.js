function savestorage() {
  localStorage.setItem("user", document.querySelector("#user").value);
  localStorage.setItem("app", document.querySelector("#app").selectedIndex);
  localStorage.setItem("client", document.querySelector("#client").value);
  localStorage.setItem("scope", document.querySelector("#scope").value);
  localStorage.setItem("secret", document.querySelector("#secret").value);
  localStorage.setItem("reqapp", document.querySelector("#reqapp").selectedIndex);
  localStorage.setItem("reqmethod", document.querySelector("#reqmethod").selectedIndex);
  localStorage.setItem("requrl", document.querySelector("#requrl").value);
  localStorage.setItem("reqbody", document.querySelector("#reqbody").value);
}

function loadstorage() {
  var iuser = localStorage.getItem("user");
  var iapp = localStorage.getItem("app");
  var iclient = localStorage.getItem("client");
  var iscope = localStorage.getItem("scope");
  var isecret = localStorage.getItem("secret");
  var ireqapp = localStorage.getItem("reqapp");
  var ireqmethod = localStorage.getItem("reqmethod");
  var irequrl = localStorage.getItem("requrl");
  var ireqbody = localStorage.getItem("reqbody");

  document.querySelector("#user").value = iuser;
  document.querySelector("#app").selectedIndex = iapp;
  document.querySelector("#client").value = iclient;
  document.querySelector("#scope").value = iscope;
  document.querySelector("#secret").value = isecret;
  document.querySelector("#reqapp").selectedIndex = ireqapp;
  document.querySelector("#reqmethod").selectedIndex = ireqmethod;
  document.querySelector("#requrl").value = irequrl;
  document.querySelector("#reqbody").value = ireqbody;

  localStorage.clear();
}

function isJsonString(str) {
  try {
    if (typeof JSON.parse(str) == "object") {
      return true;
    }
  } catch (e) {
  }
  return false;
}

function syntaxHighlight(data) {
  var json
  if (typeof data == 'string') {
    json = JSON.parse(data)
  }

  json = JSON.stringify(json, undefined, 2);
  json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
  return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
    var cls = 'number';
    if (/^"/.test(match)) {
      if (/:$/.test(match)) {
        cls = 'key';
      } else {
        cls = 'string';
      }
    } else if (/true|false/.test(match)) {
      cls = 'boolean';
    } else if (/null/.test(match)) {
      cls = 'null';
    }
    return '<span class="' + cls + '">' + match + '</span>';
  });
}

function callback(state, code, resp) {
  if (state == 4 && code == 200) {
    console.log(resp);
    data = JSON.parse(resp);
    func = data["msg"];
    if (func == "authorize") {
      savestorage(); // 保证页面刷新后各元素值仍保留

      win = window.open(data["data"].toString());  // 刷新页面

      var checkConnect = setInterval(function () {  // 弹出窗口关闭后刷新此窗口
        if (!win || !win.closed) return;
        clearInterval(checkConnect);
        window.location.reload();
      }, 100);
    } else if (func == "authcheck") {
      if (data["data"].toString() == "ok") {
        auth.style.background = "#9DC45F";
      } else {
        auth.style.background = "orange";
      }
    } else if (func == "request") {
      loading.style.display = "none";
      result.innerHTML = "";
      var data = data["data"].toString().replace(/\u0000|\u0001|\u0002|\u0003|\u0004|\u0005|\u0006|\u0007|\u0008|\u0009|\u000a|\u000b|\u000c|\u000d|\u000e|\u000f|\u0010|\u0011|\u0012|\u0013|\u0014|\u0015|\u0016|\u0017|\u0018|\u0019|\u001a|\u001b|\u001c|\u001d|\u001e|\u001f/g, "");
      if (isJsonString(data)) {
        json = syntaxHighlight(data);
        result.appendChild(document.createElement('pre')).innerHTML = json;
        // reqbody.value = json;
      } else {
        result.appendChild(document.createElement('pre')).innerHTML = '<span class="error">Error: ' + data + '</span>';
      }
    } else if (func == "link") {
      alert(data["data"].toString());
    }
  }
}

function ajax(url, method, data, mode, handler) {
  //1.创建AJAX对象
  var ajax = new XMLHttpRequest();

  //4.给AJAX设置事件(这里最多感知4[1-4]个状态)
  ajax.onreadystatechange = function () {
    handler(ajax.readyState, ajax.status, ajax.responseText);
    //5.获取响应
    //responseText		以字符串的形式接收服务器返回的信息
    // if (ajax.readyState == 4 && ajax.status == 200) {
    // }
  }

  var info = null;

  //2.创建http请求,并设置请求地址
  ajax.open(method, url);

  if (mode == "json") {
    ajax.setRequestHeader("content-type", "application/json");
    info = JSON.stringify(data);
  } else {
    //post方式传递数据是模仿form表单传递给服务器的,要设置header头协议
    ajax.setRequestHeader("content-type", "application/x-www-form-urlencoded");

    //3.发送请求(get--null    post--数据)
    if (method === "post") {
      info = "";
      for (var val in data) {
        info += val + "=" + encodeURIComponent(data[val]) + "&";  //将请求信息组成请求字符串
      }
      info = info.substr(0, info.length - 1);
    }
  }

  ajax.send(info);
}

function authcheck() {
  var selopt = app.options[app.selectedIndex];
  var data = {
    "user": user.value,
    "app": selopt.value
  }

  ajax("/ido/authcheck", "post", data, "form", callback);
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
  if (title == "API调试") {
    var el = document.querySelector("#apidiv");

  } else if (title == "应用授权") {
    var el = document.querySelector("#authdiv");

  } else if (title == "关联同步") {
    var el = document.querySelector("#flowdiv");

  }
}

window.onload = function () {

  var umodal = document.querySelector("#umodal");

  var user = document.querySelector("#user");

  var auth = document.querySelector("#auth");

  var send = document.querySelector("#send");

  var link = document.querySelector("#link");

  var flow = document.querySelector("#flow");

  var app = document.querySelector("#app");

  var reqapp = document.querySelector("#reqapp");

  var reqmethod = document.querySelector("#reqmethod");

  var requrl = document.querySelector("#requrl");

  var loading = document.querySelector("#loading");

  var rmodal = document.querySelector("#rmodal");

  var reqbody = document.querySelector("#reqbody");

  var extend = document.querySelector("#extend");

  var result = document.querySelector("#result");

  // close.onclick = function () {
  //   modal.style.display = "none";
  // }

  // window.onload = function () {
  //   init()
  // }

  window.onclick = function (event) {
    if (event.target != extend && event.target != rmodal && !extend.contains(event.target) && !rmodal.contains(event.target)) {
      rmodal.style.display = "none";
    }
  }

  app.onchange = function () {
    authcheck();
  }

  auth.onclick = function () {
    var selopt = app.options[app.selectedIndex];

    var data = {
      "app": selopt.value,
      "url": selopt.dataset.url,
      "user": user.value,
      "client": document.getElementById('client').value,
      "scope": document.getElementById('scope').value,
      "secret": document.getElementById('secret').value
    }

    // var data = {
    //   "app": selopt.value,
    //   "url": selopt.dataset.url,
    //   "user": user.value,
    //   "client": selopt.dataset.client,
    //   "scope": selopt.dataset.scope,
    //   "secret": selopt.dataset.secret
    // }

    // var data = {
    //   "app": selopt.value,
    //   "url": "https://gitlab.com/oauth",
    //   "user": user.value,
    //   "client": "2adf26541ad465a7d2e36d1c7ca6370455b288be4472e15125cd5d69dea5b750",
    //   "scope": "api",
    //   "secret": "b3b3e6b3639449b0c5f2e906054b37eae276d695376633304e446db1f1d186ef"
    // }

    ajax("/ido/authorize", "post", data, "form", callback);
  }

  document.querySelector('#serv').onclick = function (event) {
    var selopt = app.options[app.selectedIndex];
    window.open(selopt.dataset.regapp);
  }

  // document.querySelector('#req').onclick = function (event) {
  //   var selopt = app.options[app.selectedIndex];
  //   window.open(selopt.dataset.apidoc);
  // }

  user.onblur = function () {
    authcheck();
  }

  extend.onclick = function () {
    if (rmodal.style.display === "" || rmodal.style.display === "none") {
      rmodal.style.display = "block";
    } else {
      rmodal.style.display = "none";
    }
  }

  send.onclick = function () {

    // if (auth.style.background == "orange") {
    //   result.innerHTML = "";
    //   result.appendChild(document.createElement('pre')).innerHTML = '<span class="error">Error: need authorization</span>';
    //   return;
    // }

    if (requrl.value == "") {
      result.innerHTML = "";
      result.appendChild(document.createElement('pre')).innerHTML = '<span class="error">Error: request url is empty</span>';
      return;
    }

    loading.style.display = "block";

    var data = {
      "user": user.value,
      "app": reqapp.options[reqapp.selectedIndex].value,
      "method": reqmethod.options[reqmethod.selectedIndex].value,
      "url": requrl.value,
      "body": reqbody.value
    }

    ajax("/ido/request", "post", data, "form", callback);
  }

  link.onclick = function () {
    var data = {
      "id": "1",
      "name": "gitlab到mstodo同步",
      "owner": "me",
      "event": {
        "app": "gitlab",
        "class": "issue/comment/repo",
        "scope": "group/project/item",
        "item": ["*", "ebcpaas", "baas"],
        "op": ["create", "close", "update", "reopen", "destroy"]
      },
      "actions": [
        {
          "app": "mstodo",
          "class": "task",
          "scope": "group/task",
          "mapping": {
            "group": "{group}",
            "task": "[{projectName}]{issue}",
            "content": "{content}"
          },
          "op": {
            "create": "create",
            "close": "done",
            "update": "modify",
            "reopen": "start",
            "destroy": "delete"
          },
          "status": "run"
        },
        {
          "app": "github",
          "class": "issue",
          "scope": "group",
          "mapping": {
            "group": "{group}",
            "project": "{projectName}",
            "issue": "{issue}",
            "content": "{content}"
          },
          "op": {
            "create": "create",
            "close": "close",
            "update": "update",
            "reopen": "reopen",
            "destroy": "destroy"
          },
          "status": "stop"
        }
      ],
      "timestamp": 156677654,
      "status": "run"
    }

    flow.value = JSON.stringify(data, null, 2);

    ajax("/ido/link", "post", JSON.parse(flow.value), "json", callback);
  }

  loadstorage();

  authcheck();
}