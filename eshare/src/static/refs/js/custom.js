;

function getAttributes($node) {
    var attrs = {};
    if ($node === undefined) {
        return attrs;
    }
    $.each($node[0].attributes, function (index, attribute) {
        if (attribute.name != "id" && attribute.name != "style") {
            attrs[attribute.name] = attribute.value;
        }
    });

    return attrs;
}

function getQueryString(name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
    var r = window.location.search.substr(1).match(reg);
    if (r != null) return unescape(r[2]);
    return null;
}

function submit() {
    var form = document.getElementById('loginForm');
    form.action = "/login?surl=" + GetQueryString("surl")
    //进行下一步
    return true;
}

function change() {
    document.getElementById("photoCover").value = document.getElementById("lefile").value;
}

function formatDateTime(timeStamp) {
    var date = new Date();
    date.setTime(timeStamp * 1000);
    var y = date.getFullYear();
    var m = date.getMonth() + 1;
    m = m < 10 ? ('0' + m) : m;
    var d = date.getDate();
    d = d < 10 ? ('0' + d) : d;
    var h = date.getHours();
    h = h < 10 ? ('0' + h) : h;
    var minute = date.getMinutes();
    var second = date.getSeconds();
    minute = minute < 10 ? ('0' + minute) : minute;
    second = second < 10 ? ('0' + second) : second;
    return y + '-' + m + '-' + d + ' ' + h + ':' + minute + ':' + second;
};

function bytesToSize(bytes) {
    if (bytes === 0) return '0 B';

    var k = 1024;

    sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

    i = Math.floor(Math.log(bytes) / Math.log(k));

    return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i];
    //toPrecision(3) 后面保留一位小数，如1.0GB                                                                                                                  //return (bytes / Math.pow(k, i)).toPrecision(3) + ' ' + sizes[i];
}


function getDocument(docid) {
    $.post("/docinfo", {
        //id: "d714dabdd3286329bb1210646d0acf30",
        //id: "7e6801b711c3efd329d7983ff083be99", // arch.ppt
        //id: "4a43505883f8783d6b97b76cba74f8fa",
        //id: "c94b101114c3e63c69a757f630380572",
        //id: "341189b9762a35d95c0be62612655e8c", // xlsx
        //id: "c68e51c212ad5db49bd866602821b84f", // p.jpg
        //id: "3f92937346d59dc3c0145152b0f72223", // ca.pdf
        id: docid,
    }, function (data, status) {
        if (status == "success") {
            doc = data.res
            if (doc) {
                $("#imgDiv").empty()

                console.log(doc);
                // 图标问题
                // doctype = doc.ext.replace("docx", "word");
                var doctype = doc.ext.replace(".", "");

                $(".wenku-doc-type").attr("alt", doctype + "文档");
                $(".wenku-doc-type").attr("src", "/static/refs/img/" + doctype +
                    "_24.png"); // 增加图标

                $(".wenku-doc-title").data("docid", doc.id);
                $(".wenku-doc-title-span").text(doc.title);

                $(".wenku-doc-catalog").text(doc.catalog);
                $(".wenku-doc-class").text(doc.class);
                $(".wenku-doc-subclass").text(doc.subclass);

                $(".wenku-doc-owner").text(doc.owner);
                $(".wenku-doc-pagenum").text(doc.pagenum);
                $(".wenku-doc-dcnt").text(doc.dcnt);
                $(".wenku-doc-vcnt").text(doc.vcnt);
                $(".wenku-doc-ccnt").text(doc.ccnt);
                $(".wenku-doc-date").text(formatDateTime(doc.date));
                $(".wenku-doc-size").text(bytesToSize(doc.size));
                $(".wenku-doc-perm").text(doc.perm[1]); // 此处perm值应依据用户身份判定(是否交由后台?)

                $(".wenku-star").removeClass("star-35")
                $(".wenku-star").addClass("star-" + doc.score)
                $(".wenku-doc-score").text((Number(doc.score) / 10).toFixed(1))
                // $(".wenku-doc-score").text((doc.score / 10.0).toFixed(1)); // 后台记录评分为原始评分*10

                if (doc.desc) {
                    var descs = doc.desc.split("\n");
                    $.each(descs, function (idx, desc) {
                        // console.log(idx, desc);
                        $tmpp = $('<span>' + desc + '</span><br>')
                        $(".wenku-doc-desc").append($tmpp);
                    })
                }


                $(doc.tag).each(function (i, item) {
                    // console.log(item);
                    $tmptag = $('<label class="tag-label">' + item + '</label>');
                    $(".wenku-doc-tag").append($tmptag);
                });

                $(".wenku-chain-docid").text(doc.id);
                $(".wenku-chain-endorse").text(doc.endorse);
                $(".wenku-chain-block").text(doc.block);
                $(".wenku-chain-txhash").text(doc.txhash);

                getPages(doc.id, "1-3");
            }
        }
    }, "json")
}


function getPages(docid, pgrg) {
    $.post("/pageinfo", {
        //id: "d714dabdd3286329bb1210646d0acf30",
        //id: "7e6801b711c3efd329d7983ff083be99", // arch.ppt
        //id: "4a43505883f8783d6b97b76cba74f8fa",
        //id: "c94b101114c3e63c69a757f630380572",
        //id: "341189b9762a35d95c0be62612655e8c", // xlsx
        //id: "c68e51c212ad5db49bd866602821b84f", // p.jpg
        //id: "3f92937346d59dc3c0145152b0f72223", // ca.pdf
        id: docid,
        range: pgrg
    }, function (data, status) {
        if (status == "success") {
            // console.log(data);
            pagenum = Number($(".wenku-doc-pagenum:eq(0)").text());
            $(data.res).each(function (i, item) {
                console.log(item);
                if (item.prenum == -1 || i < item.prenum) {
                    // svg格式前缀为data:image/svg+xml;base54,
                    $image = $(
                        '<img class="wenku-lazy wenku-viewer-img" alt="第' +
                        item.number +
                        '页/共' + pagenum + '页"' + ' data-next="' + (
                            item.number + 1) + '" ' +
                        'src="data:image/png;base64,' +
                        item.content + '" width="100%" />');
                    $("#imgDiv").append($image);
                    // console.log($("#imgDiv").children().length);
                    $(".wenku-unread-pages").text(pagenum - item.number);
                    $(".wenku-current-page").text(item.number);
                    /*
                    $("#pic").attr("src",
                        "data:image/svg+xml;base64," +
                        item.content);
                        */
                }
            });
        }

    }, "json")
}






$(function () {

    $("#viewdoc").click(function () {
        window.location.href = "/view?docid=" + $(this).data("docid")
    });

    $("#uponemore").click(function () {
        window.location.href = "/share"
    });

    /*
    onclick="$('input[id=lefile]').click();"
    $('#lefile').change(function () {
        console.log("aaa");
        $('#photoCover').val($(this).val());
    });
    */

    $('#imgDiv').hover(function () {
        $('#imgDiv').css({
            'overflow-y': 'auto'
        })
    }, function () {
        $('#imgDiv').css({
            'overflow-y': 'hidden'
        })
    })

    $("#echo").click(function () {
        $.post("/echo", {
            id: "712424cd388514767af651419a4e0407"
        }, function (data, status) {
            // console.log(data, status);
        }, "json")
    });

    $(".wenku-viewer-more-btn").click(function () {
        docid = $(".wenku-doc-title").data("docid");
        start = Number($(".wenku-current-page").text()) + 1;
        end = start + 9; // 每次加载10页
        getPages(docid, start + "-" + end);
    });

});