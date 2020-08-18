---
marp: true
size: 16:9
#theme: gaia
#theme: uncover
#class:
#  - invert
#  - lead
paginate: true
header: sample
footer: page
#backgroundColor: white
#backgroundImage: url(a1.jpg)
#color: orange

---

<style>
img[alt~="center"] {
  display: block;
  margin: 0 auto;
}

section {
  padding: 100px;
  //background: gray;
}

p { color: orange; }

</style>

<!--_backgroundColor: aqua -->
<!--_color: red -->
<!--_theme: default -->

<!-- _class: invert -->

html directive设置
- 默认: 对当前与其后的页面有效
- $: 对所有页面有效
- _: 仅仅对当前页面有效


<u>hogehoge</u>

| Markdown元素        | HTML元素   |
|---------------------|------------|
| #                   | h1         |
| 普通文字            | p          |
| 换号                | br         |
| 分隔线              | hr         |
| 空格                | &nbsp;     |
| >                   | blockquote |
| \[]()   < > http:// | a          |
| \!\[]()             | img        |
| ** **               | strong     |
| * *                 | i em       |
| ~~ ~~               | s          |
| - + *               | ul -> li   |
| 1.    2.            | ol -> li   |
| `     ```           | code       |
| :----:              | table      |
| $ $     $$ $$       | p -> svg   |

---

> marp -p --server --html true . 

```flow
st=>start: 开始
op=>operation: My Operation
cond=>condition: Yes or No?
e=>end
st->op->cond
cond(yes)->e
cond(no)->op
```

1. 我经常去的几个网站[Google][1], [Leanote][2]和 [Marp][3]

[1]:http://www.google.com "Google"
[2]:http://www.leanote.com "Leanote"
[3]:https://marpit.marp.app/theme-css "Marp"

$E=mc^2$
$$f(x_1,x_x,\ldots,x_n) = x_1^2 + x_2^2 + \cdots + x_n^2 $$

---

<style scoped>span { background-color: #abc; }</style>
<style scoped>h4 { color: green; }</style>
<style scoped>li { color: blue; }</style>
#### 这是一张 &nbsp;&nbsp;&nbsp;&nbsp; 图片

<span style="border-bottom:2px dashed purple;">所添加的需要加下划线的行内文字</span>

![w:300 h:400 bg right:30% auto](bg1.jpg)

![w:300 center brightness contrast drop-show invert sepia opacity grayscale saturate](bg2.jpg)

- <img src="https://assets.zeit.co/image/upload/front/assets/design/now-black.svg" width="24" height="24" valign="center" /> **[Now](https://marp-cli-example.yhatt.now.sh/)**  `code`

###### <span style="color:red">背景前景</span> <my@example.com>

> 并列  http://g.com <img src="progit.png" width="80" height="80" align=right/>

![w:200](a2.jpg) ![w:200](bg3.jpg)
