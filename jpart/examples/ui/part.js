/**
 * 数据部分
 * @type {string}
 */

let jp = new Object(); //定义一个根对象

// Logo
jp.logo = "<span class='logo'><span>j</span><span>part</span></span>";

jp.banner = {
    title:jp.logo,
    text:"一起<span style='color: red'>组装</span>您的前端程序",
    "class":"big"
};

// AD
jp.ad = [
    {text:"Simplest"},
    {text:"Quickly"},
    {text:"Component"},
    {text:"Lightest"},
    {text:"Efficient"},
];

jp.load = [
    {"text":"<span data-source=$jp.logo></span> 核心库1.0.0版", href:"jpartcore.js"},
];

// Navigation
jp.nav = [
    {text:"首页", href:"./examples/_index.html", onclick:"jp.loadPage(this)", "data-affect":"#body"},
    {text:"文档", href:"./examples/document.html", onclick:"jp.loadPage(this)", "data-affect":"#body"},
    {text:"示例", href:"./examples/editor.html", onclick:"jp.loadPage(this)", "data-affect":"#body"},
    {text:"下载", href:"./examples/download.html", onclick:"jp.loadPage(this)", "data-affect":"#body"}
];

// index.html
jp.banner1 = {
    title:jp.logo,
    text:"本网站完全使用<span data-source='$jp.logo'></span>来实现",
    "class":"big",
}

jp.preview = [
    {
        // image:"./examples/ui/image/web1.png",
        title:"极简的",
        text:"<span data-source='$jp.logo'></span> 的学习应用是如此简单：它不需要您有拥有丰富的前端知识。您可以从零开始，去了解一点点HTML知识，然后您就可以在前端中，去很好的应用它。"
    },
    {
        // image:"./examples/ui/image/web2.png",
        title:"快速的",

        text:"<span data-source='$jp.logo'></span> 构建网页来是如此快速：我们为您准备好了各式各样网页组件。您所需要做的是下载它、搭建它，连通后台传来的数据，您就完成了您的前端工作。"
    },
    {
        // image:"./examples/ui/image/web3.png",
        title:"标准的",
        text:"<span data-source='$jp.logo'></span> 提供了一套标准去指导组件设计与开发：如果您有较强的前端知识，您可以按我们提供的标准设计开发组件，这此组件可以开放给您或其它的朋友使用。"
    },
];

jp.about = {
    title:"站在前端，我想做点什么",
    subtitle:"一个普通程序员对前端长久的思考",
    text:"<span data-source='$jp.logo'></span> 是我多年前编写的，放在电脑硬盘的角落。直到2014年，我又要为APP编写系统后台，寻找了一此前端框架。当时Angular、React正火，我尝试使用一下，特别是React，我发现它跟我编写的这个非常的相似。于我翻出这个框架，把它初步完善并应用于项目中，整个过程自我感觉还不错，我又给它取了个名字，叫：<span data-source='$jp.logo'></span> ，同时做了一个简单的应用说明，然后把它挂到了github上。<br><br>\n" +
        "这样时间又在不知不觉中过去。直到最近，我接手朋友开发的一个小程序，修改后台时发现它是使用Vue实现前端的，这又让我想起了我的 <span data-source='$jp.logo'></span> 。我去github上看它，它孤零零的在那里，Star依然为零。<br>\n" +
        "<span data-source='$jp.logo'></span> 能让我在前端开发中不断的想起它，是因为它远比Angular、React、Vue来得更简单，对JS、HTML、CSS的分离实现得更彻底。并不是说 <span data-source='$jp.logo'></span> 比Angular、React、Vue更优秀，而是它们间的使命不相同。在大多数前端开发中，我们只做一些小事情，我们并不需要一个高大全的前端框架，我们只需要把后台传来的数据格式化后呈现出来，我们甚至不要求界面有多炫，我们只需要它看起来还不错。我们不想，今天做的事，在昨天也做过，在前天也做过，前一个星期也做过，前一个月也做过，前一年也做过……。我们会不停的编写导航条、横幅、列表、图文预览、全文阅读等等。这让我做了这个框架，准确的说是设计一个标准，想用类似于工业标准的方式来解决这些问题。<span data-source='$jp.logo'></span> 是如何做到的呢？。<br><br>\n" +
        "首先，<span data-source='$jp.logo'></span> 将页面划分为不同的区块，每个区块都放置了组件，每个组件连接后台的数据，并负责将数据格式化后呈现出来：<br>\n" +
        "1. 这此组件可以是导航条、横幅、列表、表格、图文预览、全文阅读等。<br>\n" +
        "2. 每个组件定义了数据加载的接口，用来获取后台数据，格式化后并呈现它。<br>\n" +
        "3. 各个区块与组件组合在一起，那就是您所需要的页面。<br><br>\n" +
        "其次，<span data-source='$jp.logo'></span> 彻底将JavaScript、HTML\CSS分开，让零基础的朋友也能快速的实现前端开发。：<br>\n" +
        "1. 零基础的朋友，可以下载核心类库及组件，并在HTML页面上引用它们，只需要基本的HTML知识就可以开发出前端。<br>\n" +
        "2. 有CSS知识的朋友，可以重新设计自己CSS，这样就拥有了自己的个性化页面。<br>\n" +
        "3. 有JavaScript知识的朋友，可以按组件的标准来设计自己所需要的组件，这样不仅可以贴切自己的需求，还能分享给朋友和未来的自己使用，让人生少了一些重复的工作。<br><br>\n" +
        "最后，<span data-source='$jp.logo'></span> 有着自己对未来的想像：<br>\n" +
        "1. 会有一个公共的平台，让一些动手能力好的朋友提交设计的组件，每个组件必需有组件实现代码（js）、样式呈现代码（CSS）、组件需要的最全默认数据，这样即能保证它的呈现效果，有能说明它的数据接口及样式参考。<br>\n" +
        "2. 有需要的朋友，只需要下载他所需求的组件，然后向后端人员提交组件接口说明，获取数据。就可以基于 <span data-source='$jp.logo'></span> 核心库，像搭积木一样，搭出自己所需要的页面。这个朋友完全可以是零基础。<br>\n" +
        "3. 当共享的组件达到一定数量，也许您就不需要前端人员了，自己来做前端开发。<br><br>\n" +
        "<span data-source='$jp.logo'></span> 是能做到这些的，这也是我再一次完善它并把它推出来的原因。遵循 MIT 开源协议，您可以使用它，您会越来越喜欢它的。期待您把它做到更好，或和我一起把它做到更好，谢谢您。\n",
};

jp.tail = {
    title:"谢谢您",
    subtitle:"遵循 MIT 开源协议",
    text:"Copyright © 2014-2019 Jiang YouHua"
};

// document
jp.catalog = [{
    text: "概述",
    href: "#document-start",
    onclick:"jp.anchorDocument(this)",
}, {
    text: "页面设计",
    href: "#page-design",
    onclick:"jp.anchorDocument(this)",
    sub: [{
        text: "一、data-source 属性",
        onclick:"jp.anchorDocument(this)",
        href: "#data-source"
    }, {
        text: "二、data-func 属性",
        onclick:"jp.anchorDocument(this)",
        href: "#data-func"
    }, {
        text: "三、data-part 属性",
        onclick:"jp.anchorDocument(this)",
        href: "#data-part"
    }, {
        text: "四、data-back 属性",
        onclick:"jp.anchorDocument(this)",
        href: "#data-back"
    }, {
        text: "五、data-affect 属性",
        onclick:"jp.anchorDocument(this)",
        href: "#data-affect"
    }]
}, {
    text: "组件设计",
    href: "#part-design",
    onclick:"jp.anchorDocument(this)",
    sub: [{
        text: "一、组件类的基本结构与设计要求",
        onclick:"jp.anchorDocument(this)",
        href: "#part-base"
    }, {
        text: "二. 组件类的基本接口",
        onclick:"jp.anchorDocument(this)",
        href: "#part-interface"
    }]
}, {
    text: "内容类说明",
    href: "#default-class",
    onclick:"jp.anchorDocument(this)",
    sub:[{
        text: "一、WEB类",
        onclick:"jp.anchorDocument(this)",
        href: "#web-class"
    }, {
        text: "二、HTML类",
        onclick:"jp.anchorDocument(this)",
        href: "#html-class"
    }]
},{
    text: "流程说明",
    onclick:"jp.anchorDocument(this)",
    href: "#document-end",
    sub:[{
        text: "一、jpart加载数据的流程",
        onclick:"jp.anchorDocument(this)",
        href: "#nodes-text"
    }, {
        text: "二、jpart流程图",
        onclick:"jp.anchorDocument(this)",
        href: "#nodes-image"
    }]
}];

jp.banner2 = {
    image:"",
    title:jp.logo,
    text:"应用文档",
    "class":"title",
};

/**
 * 文档数据
 */

jp.code = {
    string:"&lt;div <b>data-source='我是字面量'</b>&gt;&lt;/div&gt;",
    jsString:"&lt;script&gt;<br>\n" +
        "<b>jp.docValue</b> = '我是JS变量值';<br>\n" +
        "&lt;/script&gt;<br>\n" +
        "&lt;div <b>data-source='$jp.docValue'</b>&gt;&lt;/div&gt;",
    file:"&lt;script&gt;<br>\n" +
        "<b>jp.docFile</b> = '/jpart/examples/file/info.txt';<br>\n" +
        "&lt;/script&gt;<br>\n" +
        "&lt;div <b>data-source='$jp.docFile'</b>&gt;&lt;/div&gt;",
    ajax:"&lt;script&gt;<br>\n" +
        "<b>jp.docAjax</b> = '/jpart/examples/file/value.json';<br>\n" +
        "&lt;/script&gt;<br>\n" +
        "&lt;div <b>data-source='$jp.docAjax'</b>&gt;&lt;/div&gt;",
    input:"&lt;input autocomplete='off' type='text' <b>data-source='jpart'</b> /&gt;",
    select:"&lt;select <b>data-source='3'</b> &gt; <br>\n" +
        "&nbsp;&nbsp;&nbsp;&lt;option value='0'&gt; Angular &lt;/option&gt;<br>\n" +
        "&nbsp;&nbsp;&nbsp;&lt;option value='1'&gt; React &lt;/option&gt;<br>\n" +
        "&nbsp;&nbsp;&nbsp;&lt;option value='2'&gt; Vue &lt;/option&gt;<br>\n" +
        "&nbsp;&nbsp;&nbsp;&lt;option value='3'&gt; jpart &lt;/option&gt;<br>\n" +
        "&lt;/select&gt;",
    img: "&lt;img <b>data-source='/jpart/examples/ui/image/logo.png'</b> /&gt;",
    func:"&lt;script&gt;<br>\n" +
        "<b>jp.docFunc = function(data)</b>{<br>\n" +
        "&nbsp;&nbsp;&nbsp;return data + ': 我是JiangYouhua'; <br>\n" +
        "}<br>\n" +
        "&lt;/script&gt;<br>\n" +
        "&lt;div <b>data-source='jpart'</b> <b>data-func='jp.docFunc'</b>&gt;&lt;/div&gt;",
    funcArgs:"&lt;script&gt;<br>\n" +
        "<b>jp.docFuncArgs = function(data, a, b, c)</b>{<br>\n" +
        "&nbsp;&nbsp;&nbsp;return data = {<br>\n" +
        "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;data:data,<br>\n" +
        "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;number:a + b + c<br>\n" +
        "&nbsp;&nbsp;&nbsp;}; <br>\n" +
        "}<br>\n" +
        "&lt;/script&gt;<br>\n" +
        "&lt;div <b>data-source='jpart'</b> <b>data-func='jp.docFuncArgs(1, 2, 3)'</b>&gt;&lt;/div&gt;",
    part:"&lt;script&gt;<br>\n" +
        "<b>jp.docNavigation</b> = [{<br>\n" +
        "&nbsp;&nbsp;&nbsp;text:'Angular',<br>\n" +
        "&nbsp;&nbsp;&nbsp;onclick:\"alert('Angular')\"<br>\n" +
        "},{<br>\n" +
        "&nbsp;&nbsp;&nbsp;text:'React',<br>\n" +
        "&nbsp;&nbsp;&nbsp;onclick:\"alert('React')\"<br>\n" +
        "},{<br>\n" +
        "&nbsp;&nbsp;&nbsp;text:'Vue',<br>\n" +
        "&nbsp;&nbsp;&nbsp;onclick:\"alert('Vue')\"<br>\n" +
        "},{<br>\n" +
        "&nbsp;&nbsp;&nbsp;text:'jpart',<br>\n" +
        "&nbsp;&nbsp;&nbsp;onclick:\"alert('jpart')\"<br>\n" +
        "}];<br>\n" +
        "&lt;/script&gt;<br>\n" +
        "&lt;div <b>data-source='$jp.docNavigation'</b> <b>data-part='JNavigation'</b>&gt;&lt;/div&gt;",
    partConfig:"&lt;script&gt;<br>\n" +
        "<b>jp.docNavigation</b> = [{<br>\n" +
        "&nbsp;&nbsp;&nbsp;text:'Angular',<br>\n" +
        "&nbsp;&nbsp;&nbsp;onclick:\"alert('Angular')\"<br>\n" +
        "},{<br>\n" +
        "&nbsp;&nbsp;&nbsp;text:'React',<br>\n" +
        "&nbsp;&nbsp;&nbsp;onclick:\"alert('React')\"<br>\n" +
        "},{<br>\n" +
        "&nbsp;&nbsp;&nbsp;text:'Vue',<br>\n" +
        "&nbsp;&nbsp;&nbsp;onclick:\"alert('Vue')\"<br>\n" +
        "},{<br>\n" +
        "&nbsp;&nbsp;&nbsp;text:'jpart',<br>\n" +
        "&nbsp;&nbsp;&nbsp;onclick:\"alert('jpart')\"<br>\n" +
        "}];<br>\n" +
        "&lt;/script&gt;<br>\n" +
        "&lt;div <b>data-source='$jp.docNavigation'</b> <b>data-part='JNavigation'</b> <b>data-config=2</b> &gt;&lt;/div&gt;",
    back:"&lt;script&gt;<br>\n" +
        "<b>jp.docIndex</b> = 0;<br>\n" +
        "<b>jp.docNavigation</b> = [{<br>\n" +
        "&nbsp;&nbsp;&nbsp;text:'Angular',<br>\n" +
        "&nbsp;&nbsp;&nbsp;onclick:\"alert('Angular')\"<br>\n" +
        "},{<br>\n" +
        "&nbsp;&nbsp;&nbsp;text:'React',<br>\n" +
        "&nbsp;&nbsp;&nbsp;onclick:\"alert('React')\"<br>\n" +
        "},{<br>\n" +
        "&nbsp;&nbsp;&nbsp;text:'Vue',<br>\n" +
        "&nbsp;&nbsp;&nbsp;onclick:\"alert('Vue')\"<br>\n" +
        "},{<br>\n" +
        "&nbsp;&nbsp;&nbsp;text:'jpart',<br>\n" +
        "&nbsp;&nbsp;&nbsp;onclick:\"alert('jpart')\"<br>\n" +
        "}];<br>\n" +
        "<b>jp.docBack = function(id)</b>{<br>\n" +
        "&nbsp;&nbsp;&nbsp;dom =document.getElementById(id).querySelectorAll(\"a\");<br>\n" +
        "&nbsp;&nbsp;&nbsp;dom[jp.docIndex].style.background = \"#ff8888\";<br>\n" +
        "}\n" +
        "&lt;/script&gt;<br>\n" +
        "&lt;div <b>data-source='$jp.docNavigation'</b> <b>data-part='JNavigation'</b><br> <b>data-back='jp.docBack'</b></b>&gt;&lt;/div&gt;",
    backArgs:"&lt;script&gt;<br>\n" +
        "<b>jp.docNavigation</b> = [{<br>\n" +
        "&nbsp;&nbsp;&nbsp;text:'Angular',<br>\n" +
        "&nbsp;&nbsp;&nbsp;onclick:\"alert('Angular')\"<br>\n" +
        "},{<br>\n" +
        "&nbsp;&nbsp;&nbsp;text:'React',<br>\n" +
        "&nbsp;&nbsp;&nbsp;onclick:\"alert('React')\"<br>\n" +
        "},{<br>\n" +
        "&nbsp;&nbsp;&nbsp;text:'Vue',<br>\n" +
        "&nbsp;&nbsp;&nbsp;onclick:\"alert('Vue')\"<br>\n" +
        "},{<br>\n" +
        "&nbsp;&nbsp;&nbsp;text:'jpart',<br>\n" +
        "&nbsp;&nbsp;&nbsp;onclick:\"alert('jpart')\"<br>\n" +
        "}];<br>\n" +
        "<b>jp.docBackArgs = function(id, num)</b>{<br>\n" +
        "&nbsp;&nbsp;&nbsp;dom =document.getElementById(id).querySelectorAll(\"a\");<br>\n" +
        "&nbsp;&nbsp;&nbsp;dom[num].style.background = \"#ff8888\";<br>\n" +
        "}\n" +
        "&lt;/script&gt;<br>\n" +
        "&lt;div <b>data-source='$jp.docNavigation'</b> <b>data-part='JNavigation'</b><br> <b>data-back='jp.docBackArgs(2)'</b></b>&gt;&lt;/div&gt;",
    affect:"&lt;script&gt;<br>\n" +
        "<b>jp.docInput</b> = 'Hello word';<br>\n" +
        "<b>jp.docAffect = function(data)</b>{<br>\n" +
        "&nbsp;&nbsp;&nbsp;docInput = data + ': Hello Jiang Youhua'; <br>\n" +
        "}<br>\n" +
        "&lt;/script&gt;<br>\n" +
        "&lt;div data-source='jpart' <b>data-back='jp.docAffect(\"A\")'</b><br> <b>data-affect='#theInput'</b>&gt;&lt;/div&gt;<br>\n" +
        "&lt;input autocomplete='off' <b>id='theInput'</b> type='text' <b>data-source='$jp.docInput'</b> /&gt;\n" +
        "                                ",
    affectEvent:"&lt;script&gt;<br>\n" +
        "<b>jp.jsEvent = function(it)</b>{<br>\n" +
        "&nbsp;&nbsp;&nbsp;// TODO<br>\n" +
        "}<br>\n" +
        "&lt;/script&gt;<br>\n" +
        "1. &lt;input autocomplete='off' type='text' <b>name='jp.docTogether'</b> <b>onkeyup='jp.jsEvent(this)'</b><br> <b>data-affect='#theTogether'</b> /&gt;<br>\n" +
        "2. &lt;input autocomplete='off' type='text' <b>id='theTogether'</b> <b>data-source='$jp.docTogether'</b> /&gt;",
};

/**
 * 下载
 */

jp.parts = [
    {id:1, name:"jpartcore", explain:"jpart核心库", example: "jpart核心库1.0.0"},
    {id:2, name:"JAd", explain:"AD、幻灯片组件", example: "<div data-source='' data-part='JAd'"},
    {id:3, name:"JBanner", explain:"横幅组件", example: "<div data-source='' data-part='JBanner'"},
    {id:4, name:"JFullText", explain:"全文阅读组件", example: "<div data-source='' data-part='JFullText'"},
    {id:5, name:"JList", explain:"列表组件", example: "<div data-source='' data-part='JList'"},
    {id:6, name:"JNavigation", explain:"导航组件", example: "<div data-source='' data-part='JNavigation'"},
    {id:7, name:"JTable", explain:"列表组件", example: "<div data-source='' data-part='JTable'"},
]

/**
 * 函数部分
 * @type {string}
 */

jp.pageUrl = "./examples/_index.html";

// 导航页面跳转
jp.loadPage = function (it) {
    jp.pageUrl = it.getAttribute("href");
    let nodes = it.parentNode.childNodes;
    for(let x = 0; x < nodes.length; x ++){
        let node = nodes[x];
        node.classList.remove("action")
    }
    it.classList.add("action");
    document.documentElement.scrollTop = 0;
    document.getElementById("body").innerHTML = "";
    document.getElementById("head").style.position = "inherit";
    if(it.innerHTML == "首页"){
        WEB.Show(".ad");
        WEB.Hide("#top");
        return;
    }
    WEB.Hide(".ad")
    WEB.Show("#top");
    document.getElementById("head").style.position = "fixed";
};

jp.anchorDocument = function(e){
    if ( e && e.preventDefault )
        e.preventDefault();
    else
        window.event.returnValue = false;
    document.documentElement.scrollTop
    let x = document.getElementById("document-top").offsetTop;
    let id = e.getAttribute("href")
    let top = document.querySelector(id).offsetTop;
    let t =  document.documentElement.scrollTop ;
    document.documentElement.scrollTop = top - x;
};

jp.loadLib = function () {
    let nodes = document.querySelectorAll("td [name='lib']");
    if(!nodes || nodes.length == 0){
        alert("您没有选择需要下载的库与组件")
        return ;
    }
    let dir = "./part/"
    let zip = new JSZip();
    let n = 0;
    for(let i = 0; i < nodes.length; i++){
        let node = nodes[i];
        if(!node.checked){
            continue;
        }
        n += 2;
        // 核心库
        if(!i){
            WEB.Get("./jpartcore.js", function (re) {
                zip.file("jpartcore.js", re);
                n = n - 2;
            })
            continue;
        }
        // js组件
        let name = jp.parts[i].name;
        WEB.Get(dir+name+".js", function (re) {
            zip.file(name+".js", re);
            n --;
        })
        //对应的CSS
        WEB.Get(dir+name+".css", function (re) {
            zip.file(name+".css", re);
            n --;
        })
    }


    let t = setInterval(function () {
        if (n <=0){
            clearInterval(t);
            if(!WEB._isObject(zip.files)){
                return;
            }
            zip.generateAsync({type:"blob"})
                .then(function(content) {
                    // see FileSaver.js
                    saveAs(content, "example.zip");
                });
        }
    }, 100);

}
