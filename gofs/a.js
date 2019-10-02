function getVideoInfo(id){
	console.log("KKKKKKK")
$.ajax({
             type: "get",
             async: false,
             url: "https://avmoo.asia/cn/search/" + id,
             dataType: "jsonp",
             jsonp: "callback",//传递给请求处理程序或页面的，用以获得jsonp回调函数名的参数名(一般默认为:callback)
             jsonpCallback:"flightHandler",//自定义的jsonp回调函数名称，默认为jQuery自动生成的随机函数名，也可以写"?"，jQuery会自动为你处理数据
             success: function(json){
                 console.log(json);
             },
             error: function(){
                 alert('fail');
             }
         });
}

document.addEventListener("DOMContentLoaded", getVideoInfo("ipz-371"));
