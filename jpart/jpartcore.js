/**
 * Created by jiangyouhua on 2016/3/31.
 * Update 2017.09.15
 * JPART JS组件化UI
 * *************************************
 * 基本类
 * URL类, 重新的解析URl
 * HTML类, 生成HTML元素字符串的对象
 * PART类, 生成HTML组件的对象
 * WEB类, 实现AJAX加载，
 * *************************************
 * 引用量
 * $string = js变量
 * @string = url
 * string = value
 * *************************************
 * 说明
 * * 本JS需要解决的是两个问题，
 * * 一、获取数据：
 * *** 1. 加载不需要处理的数据（data-file）。 如：html, text, js等等；
 * ****** 加载进来直接写入当前位置，只对js进行激活处理
 * *** 2. 加载需要处理的数据（data-source）。 如：AJAX请求，本地JS对象等；
 * ****** a. 前处理（data-func）：加载数据后通过data-func指定的函数进行处理
 * ****** b. 字符串化写入当前位置：
 * ********** ~ 直接字符串化写入当前位置
 * ********** ~ 作参数传入data-part指定的PART中，并格式化PART写入到当前位置
 * ****** c. 后处理（data-back）: 写入后调用data-back指定的函数作为回调函数处理其它业务
 * ****** d. 施加影响（data-affect）：对data-affect指定的位置，施加响应
 * * 二、处理事件：
 * *** 1. 等价处理：对onload、onclick、onkeydown、onchange等等HTML DOM事件进行等价处理
 * ******* ~如：onLoad 等价如 data-load
 * *** 2. 指定影响：等价处理后通过data-affect指定的位置，施加响应
 */

// 获取数据，加载驱动，事件驱动
const DATA_SOURCE = "data-source"; //data-source = value：为该标签提供数据源，value有$变量、@Ajax请求、[{json字符串}]、字面量
const DATA_FUNC = "data-func"; // data-func = value：data-source数据解析后，在输出前可经过该程序再处理，一般用来为JPart子类统一数格式
const DATA_BACK = "data-back"; // data-back = value：data-source数据解析并输出后，回调该程序
const DATA_AFFECT = "data-affect"; // data-affect = selector：加载数据后、及事件处理后，所需要施加影响的区域

// 组件
const DATA_PART = "data-part"; // data-part = value：使用JPart格式该标签，value为JPart子类
const DATA_CONFIG = "data-config"; // data-config = string:

// 标签
const DATA_SOURCE_TAG = "[" + DATA_SOURCE + "]";
const DATA_AFFECT_TAG = "[" + DATA_AFFECT + "]";

/**
 * 自动启用
 */
window.onload = function () {
    // 加载html页面内的组件
    WEB.Source(DATA_SOURCE_TAG);
    WEB.Event(DATA_AFFECT_TAG);
};

// 全函数
let WEB = {
    /**
     * 获取页面宽与高
     */
    Width: window.innerWidth || document.documentElement.clientWidth || document.body.clientWidth,
    Height: window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight,

    /**
     * 显示提示信息
     */
    Alert: function (message, duration) {
        if (!message) {
            return;
        }
        if (typeof message != "string") {
            return;
        }
        if(!duration){
            duration = 3000
        }
        let div = new HTML("div", message);
        div.AddAttr("id", "JAlert");
        div.AddCss("position", "fixed");
        div.AddCss("display", "inline-block");
        div.AddCss("width", "26em");
        div.AddCss("background", "#222");
        div.AddCss("color", "#fff");
        div.AddCss("padding", "1em 2em");
        div.AddCss("left", "50%");
        div.AddCss("top", "50%");
        div.AddCss("margin-left", '-15rem');
        div.AddCss("margin-top", '-5rem');
        div.AddCss("text-align", "center");
        div.AddCss("z-index", "999999");
        let node = div.Element();
        document.body.appendChild(node);
        setTimeout(function () {
            document.body.removeChild(node);
        }, duration);
        return;
    },

    Loading:function(dom, end){
        let tag = "JLoading";
        if(!dom || (!WEB._isString(dom) && !WEB._isDOM(dom))){
            dom = document.body;
        }else if(WEB._isString(dom)){
            dom = document.querySelector(dom)
        }

        if(end) {
            let nodes = document.querySelectorAll("."+tag);
            for(let i = nodes.length - 1; i >= 0; i --){
                nodes[i].remove()
            }
            return;
        }
        let path1 = new HTML("path");
        path1.AddAttr('opacity', '0.2');
        path1.AddAttr('fill', '#000');
        path1.AddAttr('d', 'M20.201,5.169c-8.254,0-14.946,6.692-14.946,14.946c0,8.255,6.692,14.946,14.946,14.946 s14.946-6.691,14.946-14.946C35.146,11.861,28.455,5.169,20.201,5.169z M20.201,31.749c-6.425,0-11.634-5.208-11.634-11.634 c0-6.425,5.209-11.634,11.634-11.634c6.425,0,11.633,5.209,11.633,11.634C31.834,26.541,26.626,31.749,20.201,31.749z');
        let path2 = new HTML("path");
        path2.AddAttr('fill', '#000');
        path2.AddAttr('d', 'M26.013,10.047l1.654-2.866c-2.198-1.272-4.743-2.012-7.466-2.012h0v3.312h0 C22.32,8.481,24.301,9.057,26.013,10.047z');
        let animate = new HTML('animateTransform');
        animate.AddAttr('attributeType', 'xml');
        animate.AddAttr('attributeName', 'transform');
        animate.AddAttr('type', 'rotate');
        animate.AddAttr('from', '0 20 20');
        animate.AddAttr('to', '360 20 20');
        animate.AddAttr('dur', '0.5s');
        animate.AddAttr('repeatCount', 'indefinite');
        path2.AddContent(animate);
        let svg = new HTML("svg");
        svg.AddClass("JPartLoad");
        svg.AddContent(path1, path2);
        svg.AddAttr("fill", "#ff6700");
        svg.AddAttr("version", "1.1");
        svg.AddAttr("xmlns", "http://www.w3.org/2000/svg");
        svg.AddAttr("xmlns:xlink", "http://www.w3.org/1999/xlink");
        svg.AddAttr("x", "0px");
        svg.AddAttr("y", "0px");
        svg.AddAttr("width", "40px");
        svg.AddAttr("height", "40px");
        svg.AddAttr("viewBox", "0 0 40 40");
        svg.AddAttr("enable-background", "0 0 40 40");
        svg.AddAttr("xml:space", "preserve");
        let div = new HTML('div', svg);
        div.AddAttr("class", "JLoading");
        div.AddCss("position", "fixed");
        div.AddCss("display", "inline-block");
        div.AddCss("left", "50%");
        div.AddCss("top", "50%");
        div.AddCss("margin-left", '-5rem');
        div.AddCss("margin-top", '-5rem');
        div.AddCss("width", "10rem");
        div.AddCss("height", "10rem");
        div.AddCss("color", "#fff");
        div.AddCss("text-align", "center");
        div.AddCss("z-index", "999999");
        let node = div.Element();

        dom.appendChild(node);
    },

    /**
     * AJAX get请求
     */
    Get: function (url, callBack, classify) {
        WEB.Post(url, null, callBack, classify);
    },

    /**
     * AJAX 请求
     * 1. json；2.js；3.html,
     */
    Post: function (url, data, callBack, classify) {
        if(!WEB._isString(url)){
            return;
        }
        // 创建对象
        let xmlhttp;
        if (window.XMLHttpRequest) { // code for IE7+, Firefox, Chrome, Opera, Safari
            xmlhttp = new XMLHttpRequest();
        } else { // code for IE6, IE5
            xmlhttp = new ActiveXObject("Microsoft.XMLHTTP");
        }

        // 设置传递
        if (!data) {
            xmlhttp.open("GET", url, true);
        } else {
            xmlhttp.open("POST", url, true);
            if (!(data instanceof FormData)) {
                xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
            }
        }

        // 监听结果
        xmlhttp.onreadystatechange = function () {
            if (xmlhttp.readyState == 4){
                if(xmlhttp.status == 200){
                    let text = xmlhttp.responseText;
                    let obj;
                    switch (classify) {
                        case "js":
                            obj = WEB.Eval(text);
                            break;
                        case "json":
                            try {
                                obj = JSON.parse(text);
                            }catch(e){
                                obj = text;
                            }
                            break;
                        default:
                            obj =  text;
                    }

                    if(!!callBack && typeof callBack === 'function'){
                        callBack(obj);
                    }
                    return false;
                }
                if(!!callBack && typeof callBack === 'function'){
                    console.log(xmlhttp.statusText);
                    callBack();
                }
            }
        }

        xmlhttp.send(data);
    },

    Eval:function(s){
        if(!WEB._isString(s)){
            return s;
        }
        // let arr = s.split(".")
        // if(arr.length == 1) {
        //     return (new Function("return " + "" + s + ""))();
        // }
        // let str = arr[0];
        // for(let i = 1; i < arr.length - 1; i ++){
        //     str = str + "['" + arr[i] + "']";
        // }
        // let a = arr[arr.length - 1];
        // arr = a.split('(');
        // str = str + "['" + arr[0] + "']";
        // if(arr.length == 1){
        //     return (new Function("return " + "" + str + ""))();
        // }
        return eval(s);
    },

    Hide:function(selector) {
        WEB._display(selector, 'none');
    },

    Show:function(selector) {
        WEB._display(selector, 'block');
    },

    /**
     * 返回全局标签
     * @returns {string}
     */
    Mark: function () {
        let mark = Math.random().toString()
        return mark.replace("0.", "J")
    },

    _display:function(selector, tag){
        if(!WEB._isString(selector) && !WEB._isDOM(selector)){
            return;
        }
        if(WEB._isDOM(selector)){
            selector.style.display = tag;
            return;
        }
        let nodes = document.querySelectorAll(selector);
        for(let i = 0; i < nodes.length; i ++){
            nodes[i].style.display = tag;
        }
    },

    /**
     * 从字符串解析对象, 默认返回对象，也可以返回[{key:k, value:v},...]
     * 1. key:value;key1:value1
     * 2. key=value&key1=value1
     */
    ObjectByString: function (str, split1, split2, isArray) {
        if (!WEB._isString(str)) {
            return null;
        }
        if (!WEB._isString(split1)) {
            return str;
        }
        let arr = str.split(split1);
        if (!WEB._isString(split2)) {
            return arr;
        }

        if (arr.length == 0) {
            return null;
        }

        let li = []
        let obj = {};
        for (let i = 0; i < arr.length; i++) {
            if (!arr[i]) {
                continue;
            }
            let v = arr[i].trim();
            let a = v.split(split2);
            if (!a || a.length == 0) {
                continue;
            }
            let key = a[0].trim();
            let value = a[1].trim();

            obj[key] = value;
            li[i] = {key:key, value:value}
        }
        return  isArray ? li : obj;
    },

    /**
     * 解析part对象
     * @param selector，基于data_source实现组件化
     * @constructor
     */
    Source: function (selector) {
        // 没有选择器
        if (!selector) {
            return;
        }
        // 如果不是字符串又不是节点对象则返回
        if (!WEB._isString(selector) &&!WEB._isDOM(selector)) {
            return;
        }
        let nodes;
        if(WEB._isString(selector)){
            if (selector.indexOf(DATA_SOURCE_TAG) < 0) {
                selector += DATA_SOURCE_TAG;
            }
            nodes = document.querySelectorAll(selector);
            return WEB._source(nodes);
        }
        nodes = selector.querySelectorAll(DATA_SOURCE_TAG)
        WEB._source(nodes)
    },

    /***
     * 根据选择器绑定事件
     * @param selector
     * @constructor
     */
    Event: function(selector){
        // 没有选择器
        if (!selector) {
            return;
        }
        // 如果不是字符串又不是节点对象则返回
        let nodes;
        if(WEB._isString(selector)){
            if (selector.indexOf(DATA_AFFECT_TAG) < 0) {
                selector += DATA_AFFECT_TAG;
            }
            nodes = document.querySelectorAll(selector);
            return WEB._event(nodes)
        }
        nodes = selector.querySelectorAll(DATA_AFFECT_TAG)
        WEB._event(nodes)
    },

    /**
     * 解析所有数据源内容
     * @param array
     * @private
     */
    _source:function(notes){
        // 没有相关内容
        if (!notes || ( notes.length == 0 && !(notes instanceof NodeList ))) {
            return;
        }

        // 解析当前Source
        for (let x = 0; x < notes.length; x++) {
            let dom = notes[x];
            let source = dom.getAttribute(DATA_SOURCE);
            let func = dom.getAttribute(DATA_FUNC);
            let back = dom.getAttribute(DATA_BACK);
            let part = dom.getAttribute(DATA_PART);
            let affect = dom.getAttribute(DATA_AFFECT);
            WEB._sourceToString(dom, source, func, back, part, affect);
        }
    },

    /**
     * 解析数据
     * @param source
     * @param func
     * @param back
     * @param part
     * @param config
     * @private
     */
    _sourceToString:function(dom, source, func, back, part, affect) {
        if(!WEB._isDOM(dom)){
            return;
        }
        WEB._parseValue(source, function (data) {
            // 没有数据则返回空
            if (!data && !WEB._isString(part)) {
                return;
            }
            // 处理数据
            let re = WEB._parseFunc(func, data);
            // 判断是否为组件
            if (WEB._isString(part)) {
                // 是组件，尝试格式化组件
                try {
                    let p = new (WEB.Eval(part));
                    let config = dom.getAttribute(DATA_CONFIG)
                    // 字符串化PART
                    p.SetIt(dom).SetData(re).SetConfig(config).SetBack(back).Out();
                    if(!!affect){
                        WEB.Source(affect);
                    }
                    WEB.Source(dom);
                    WEB.Event(dom);
                } catch (e) {
                    console.log(e);
                }
                return;
            }

            // 不是组件
            if (typeof re == "object") {
                re = JSON.stringify(re); // 字符串化对象
            }

            switch (dom.tagName.toLowerCase()) {
                case "select":
                case "input":
                    dom.value = re;
                    break;
                case "img":
                case "script":
                    dom.setAttribute("src", re);
                    break;
                case "link":
                    dom.setAttribute("href", re);
                    break;
                default:
                    dom.innerHTML = re;
            }
            WEB._parseFunc(back, re)
            WEB._parseJS(dom);
            if(!!affect){
                WEB.Source(affect);
            }
            WEB.Source(dom);
            WEB.Event(dom);
        }, dom);
    },

    _event:function(notes){
        // 没有相关内容
        if (!notes || ( notes.length == 0 && !(notes instanceof NodeList ))) {
            return;
        }

        // 解析当前affect
        for (let x = 0; x < notes.length; x++) {
            let doc = notes[x];
            if(!WEB._isDOM(doc)){
                return;
            }
            let affect = doc.getAttribute(DATA_AFFECT);
            if(!WEB._isString(affect)){
                continue;
            }

            for(let i = 0; i < doc.attributes.length; i++){
                let a = doc.attributes[i];
                if(a.name.indexOf('on') == 0){
                    WEB._eventAddListener(doc, a.name, a.value, affect)
                }
            }
        }
    },

    /**
     * 解析表单
     * @constructor
     */
    _eventAddListener:function(element, name, func, affect){
        if(!WEB._isDOM(element)){
            return;
        }
        if(!WEB._isString(name)){
            return;
        }
        if(!WEB._isString(func)){
            return;
        }
        if(!WEB._isString(affect)){
            return;
        }
        // 绑定事件
        element.addEventListener(name.substr(2), function(e){
            e.preventDefault();
            WEB.Source(affect);
        });
    },

    /**
     * 判断是否为非空字符串
     * @param s
     * @returns {boolean}
     * @private
     */
    _isString:function(s){
        return (!!s && typeof s == "string" && s.length > 0);
    },

    _isArray:function(s){
        return (!!s && (s instanceof Array || s instanceof NodeList ) && s.length > 0);
    },

    _isObject:function(s){
        return (!!s && JSON.stringify(s) != '{}');
    },

    /**
     * 判断是否为函数
     * @param f
     * @returns {boolean}
     * @private
     */
    _isFunction:function(f){
        return (!!f && typeof f === "function");
    },

    /**
     * 判断URL是否有效
     * @param u
     * @returns {boolean}
     * @private
     */
    _isUrl:function(url){
        if(!WEB._isString(url)){
            return false;
        }
        let re = new RegExp("^\./{1}|^/{1}|\.\w{2,6}$|\.\w{2,6}\?|\.\w{2,6}\#");
        return re.test(url);
    },

    /**
     * 判断DOM
     * @param obj
     * @returns {boolean}
     * @private
     */
    _isDOM:function(obj){
        return ( typeof HTMLElement === 'object' ) ? obj instanceof HTMLElement : obj && typeof obj === 'object' && obj.nodeType === 1 && typeof obj.nodeName === 'string';
    },

    /**
     * 解析变量，1.从js定义中获取；2.从表单中获取；从url.search中获取
     * @param value
     * @param callback
     * @private
     */
    _parseValue:function( _value, callback, dom){
        // 1. 不是有效字符串则返回, 2. 不是变量， 不是URL则返回
        if( !WEB._isString(_value) || (_value.indexOf("$") < 0 && !WEB._isUrl(_value))){
            return WEB._isFunction(callback)? callback(_value) : _value;
        }
        // 变量字符串
        if(_value.indexOf("$") == 0 ){
            try{
                let v = WEB.Eval(_value.substr(1))
                if(!!v) {
                    return WEB._parseValue(v, callback, dom);
                }
                return WEB._parseOther(_value, callback, dom);
            }catch (e) {
                return WEB._parseOther(_value, callback, dom);
            }
        }

        // url字符串
        let url =  WEB._parseUrl(_value);
        if(!dom){
            return WEB._parseAjax(url, callback);
        }

        switch (dom.tagName.toLowerCase()) {
            case "script":
            case "img":
            case "link":
            case "iframe":
            case "frame":
                return WEB._isFunction(callback)? callback(url) : url;
            default:
                return WEB._parseAjax(url, callback);
        }
    },

    _parseAjax:function(url, callback){
        if(!WEB._isUrl(url)){
            return WEB._isFunction(callback)? callback(url) : url;
        }
        // 获取请求的类别
        let arr = url.split(".");
        let classify;
        if(arr.length > 0){
            classify = arr[arr.length - 1].toLowerCase();
        }
        WEB.Get(url, callback, classify);
        return url
    },

    _parseOther:function(_value, callback, dom){
        // 获取表单值
        let word = _value.substr(1);
        let element = document.getElementsByName(word);
        if (element.length > 0) {
            return WEB._parseValue( element[0].value, callback, dom);
        }
        // 从window.location.search获取值
        _value = URL.SearchGet(word);
        if (!!_value) {
            return WEB._parseValue( _value, callback, dom);
        }

        return WEB._isFunction(callback)? callback(_value) : _value;
    },

    _parseFunc:function(name, data){
        if (!WEB._isString(name)){
            return data;
        }
        try {
            if (!data) {
                return WEB.Eval(name);
            }
            let func = name.substr(0, name.indexOf("("));
            let args = "";
            if(func.length > 0) {
                let start = name.indexOf("(") + 1;
                let end = name.indexOf(")");
                args = name.substr(name.indexOf("(") + 1, end - start);
                if(args.length > 0){
                    args = "," + args;
                }
                name = func;
            }
            if(WEB._isString(data)){
                data = "'"+data+"'";
            }
            return WEB.Eval(name + "(" + data + args + ")");
        }catch (e) {
            console.log(e);
            return data;
        }
    },

    /**
     * 处理JS
     * @param doc
     * @private
     */
    _parseJS:function(doc){
        if(!WEB._isDOM(doc)){
            WEB.Source(doc);
            return;
        }
        let nodes = doc.getElementsByTagName("script")
        if(!nodes || nodes.length == 0){
            WEB.Source(doc);
            return;
        }
        // 解析link
        let link = new Array();
        for(let i = 0; i < nodes.length; i++){
            let element = nodes[i];
            let src = element.getAttribute("src");
            if (WEB._isString(src)){
                link.push(src);
                continue;
            }
            let script = document.createElement('script');
            script.type = "text/javascript";
            script.text = element.innerHTML;
            document.getElementsByTagName('head')[0].appendChild(script);
            document.head.removeChild(document.head.lastChild);
        }
        if(link.length == 0){
            WEB.Source(doc);
        }

        for(let i = 0; i < link.length; i++){
            let url = link[i];
            WEB.Source(doc);
            WEB.Get(url, null,"js");
        }
        return;
    },

    /**
     * 解析url中包含的$value值
     * @param url
     * @returns {void | string | never}
     * @private
     */
    _parseUrl: function (url) {
        let reg = new RegExp("(\{\$\w+\})|(\$\w+)");
        return url.replace(reg, function (word) {
            let v = word.match(/\w+/);
            return WEB._parseValue(v);
        })
    },

}

/**
 * HTML类：简化js输出html元素
 * 1. String(), 输出字符串
 * 2. Element(), 输出DOM.Element
 *
 *  类属性
 * _tag, html元素的标签名
 * _attr, html元素的属性
 * _class, html元素的类别
 * _css, html元素的风格
 * _content, 内html元素的容
 */
let HTML = function () {
    //实例属性
    let _tag = "div";
    let _attr = {};
    let _class = [];
    let _css = {};
    let _content = null;
    let _element = null;

    //构造函数，可去掉该闭包
    ;
    (function (self, args) {
        let a = Array.prototype.slice.call(args);
        if (a.length == 0) {
            return;
        }
        self.SetTag(a.shift());
        //内容
        if (a.length == 0) {
            return;
        }
        self.AddContent(a.shift());
        //属性
        if (a.length == 0) {
            return;
        }
        self._attribute(a);
    })(this, arguments)
}

/**
 * 原型方法
 */
HTML.prototype = {
    /**
     * 设置标签名
     * @param name， 标签名称
     * @returns {HTML}
     * @Jiang youhua
     */
    SetTag: function (name) {
        if (!WEB._isString(name)) {
            return this;
        }
        this._tag = name.trim();
        return this;
    },

    /**
     * 初化并设置属性，接收键值对与非键值对
     * @param args
     * @param init
     * @returns {HTML}
     * @Jiang youhua
     */
    _attribute: function (args) {
        let _attr = false;
        let _class = false;
        let _css = false;
        // 无参数
        if (!(args instanceof Array) || args.length == 0) {
            return;
        }
        for (let i = 0; i < args.length; i++) {
            let value = args[i];
            let obj = {};
            let b = false;

            // 判断是否为字符串 "key=value key1=value1", "class=value value1"
            if (typeof value == 'string') {
                let object = WEB.ObjectByString(value, " ", "=");
                for (let x in object) {
                    obj[x] = object[x];
                    b = true;
                }
                if (!b) {
                    continue;
                }
            } else {
                // 判断是否为对象
                if ((typeof value != 'object') || (value instanceof Array)) {
                    continue;
                }
                obj = value;
            }

            for (let x in obj) {
                let v = obj[x];

                // 样式类型
                if (x == "class") {
                    if (!_class) {
                        this.SetClass(v);
                        _class = true;
                        continue;
                    }
                    this.AddClass(v);
                    continue;
                }
                // 样式分格
                if (x == "style") {
                    let object = WEB.ObjectByString(v, ";", ":");
                    for (let x in object) {
                        if (!object[x]) {
                            continue;
                        }
                        if (!_css) {
                            this.SetCss(x, object[x]);
                            _css = true;
                            continue;
                        }
                        this.AddCss(x, object[x]);
                        continue;
                    }
                }

                // 样式属性
                if (!_attr) {
                    this.SetAttr(x, v);
                    _attr = true;
                    continue;
                }
                this.AddAttr(x, v);
            }
        }
        return this;
    },

    /**
     * 设置属性
     */
    SetAttr: function (key, value) {
        this._attr = {};
        return this.AddAttr(key, value);
    },

    /**
     * 添加属性，接收键值对与非键值对
     * @param key
     * @param value
     * @returns {HTML}
     * @Jiang youhua
     */
    AddAttr: function (key, value) {
        if(!WEB._isString(key)){
            return this;
        }

        key = key.trim();
        //如果是类别
        if (key == 'class') {
            this.AddClass(value);
            return this;
        }

        //如果是风格
        if (key == 'style') {
            let object = WEB.ObjectByString(value, ";", ":");
            for (let x in object) {
                this.AddCss(x, object[x]);
            }
            return this;
        }

        //初始属性
        if (!this._attr) {
            this._attr = {};
        }
        //value为空
        if (!value && value != 0) {
            this._attr[key] = '';
            return this;
        }
        if (value instanceof Object) {
            for (let x in value) {
                this.AddAttr(key + "-" + x, value[x]);
            }
            return this;
        }
        if (typeof value == "string") {
            this._attr[key] = value.trim();
            return this;
        }
        this._attr[key] = value.toString();
        return this;
    },

    /**
     * 初始化并设置类型
     * @param args，接收字符串与数组
     * @returns {HTML}
     * @Jiang youhua
     */
    SetClass: function () {
        //空值
        this._class = {key: {}, value: []};
        let a = Array.prototype.slice.call(arguments);
        return this._classByArray(a);
    },

    /**
     * 添加类型
     * @param value，类型名
     * @returns {HTML}
     * @Jiang youhua
     */
    AddClass: function () {
        let a = Array.prototype.slice.call(arguments);
        return this._classByArray(a);
    },

    /**
     * 接收数组，转为式样类型
     */
    _classByArray: function (a) {
        if (!this._class) {
            this._class = {key: {}, value: []};
        }

        if (!(a instanceof Array)) {
            return this;
        }
        if (a.length == 0) {
            return this;
        }
        for (let i = 0; i < a.length; i++) {
            let val = a[i];
            if (!val) {
                continue;
            }
            // 如果是数组
            if (val instanceof Array) {
                this._classByArray(val);
                continue;
            }
            // 非字符串
            if (typeof val != "string") {
                continue;
            }
            //字符串
            let arr = val.split(" ");
            if (arr.length > 1) {
                this._classByArray(arr);
                continue;
            }

            let v = val.trim();
            if (this._class.key[v]) {
                continue;
            }
            this._class.key[v] = true;
            this._class.value.push(v);
        }
        return this;
    },

    // 获取样式类型
    GetClass: function () {
        if (!this._class || !this._class.value) {
            return "";
        }
        return this._class.value.join(" ");
    },

    // 设置样式属性，重置样式为空后再添加样式
    SetCss: function (key, value) {
        this._css = {};
        return this.AddCss(key, value);
    },

    // 添加样式属性
    AddCss: function (key, value) {
        if (!WEB._isString(key)) {
            return this;
        }
        if (!value || (typeof value != "string" && typeof value != 'number')) {
            return this;
        }
        if (!this._css) {
            this._css = {};
        }
        this._css[key] = value;
        return this;
    },

    // 获取样式属性
    GetCss: function (key) {
        if (!WEB._isString(key)) {
            return ;
        }
        if (!this._css) {
            return ;
        }
        return this._css[key];
    },

    /**
     * 初始化并设置内容
     * @param contents，内容，支持所有对象
     * @returns {HTML}
     * @Jiang youhua
     */
    SetContent: function () {
        this._content = [];
        let arr = Array.prototype.slice.call(arguments);
        return this._contentByArray(arr);
    },

    /**
     * 添加内容
     * @param content，内容，支持所有对象
     * @returns {HTML}
     * @Jiang youhua
     */
    AddContent: function () {
        let arr = Array.prototype.slice.call(arguments);
        return this._contentByArray(arr);
    },

    _contentByArray: function (arr) {
        //空值
        if (!(arr instanceof Array) || arr.length == 0) {
            return this;
        }
        for (let x in arr) {
            let v = arr[x];
            if (this == v) {
                return this;
            }
            //初始化
            if (!this._content) {
                this._content = [];
            }
            this._content.push(v);
        }
        return this;
    },

    /**
     * 格式化内容为字符串
     * @param contents，内容
     * @returns {*}
     */
    _forContent: function (contents) {
        //非对象，转字符串返回
        if (!contents && contents != 0) {
            return '';
        }
        if (typeof contents != 'object') {
            return contents.toString();
        }
        //html对象，转字符串返回
        if (contents instanceof HTML) {
            return contents.String();
        }
        //对象，遍历递归调用
        let s = '';
        for (let x in contents) {
            s += this._forContent(contents[x]);
        }
        return s;
    },

    /**
     * HTML对象输出文档对象节点
     * @return {dom.Element}
     */
    Element: function () {
        let div = document.createElement("div");
        div.innerHTML = this.String();
        return div.childNodes[0];
    },

    /**
     * 本实例输出为HTML字符串
     * @returns {string}
     * @Jiang youhua
     */
    String: function () {
        // 处理标签
        if (!this._tag) {
            this._tag = "div";
        }
        // 处理类别
        let c = '';
        if (!!this._class && !!this._class.value) {
            c = this._class.value.join(" ");
            c = c.trim();
            if (!!c) {
                c = 'class="' + c + '"';
            }
        }
        // 处理属性
        let a = ''
        for (let x in this._attr) {
            if (!this._attr[x]) {
                a = a + " " + x; //单词属性
            } else {
                a = a + ' ' + x + '="' + this._attr[x] + '"'; //key=value
            }
        }
        // 处理风格
        let s = ''
        for (let x in this._css) {
            if (!this._css[x]) {
                continue;
            }
            if (!s) {
                s += 'style="';
            }
            let v = x + ":" + this._css[x];
            s = s + v + ";";
        }
        if (!!s) {
            s += '"';
        }

        // 非内容部分
        a = a.trim();
        let h = "<" + this._tag + " " + c + " " + a + " " + s;
        switch (this._tag) {
            case "br":
            case "hr":
            case "img":
            case "input":
            case "meta":
            case "param":
            case "link":
            case "animateTransform":
                return h + "/>"
        }
        let str = this._forContent(this._content);
        return h + ">" + str + "</" + this._tag + ">";
    }
}


/**
 * 组件类，为所有组件的父类
 * 1. 抽象类
 * 2. 所有组件均要继承该类
 * _it, 当前对象
 * _data, 数据
 * _config, 配置
 * _func, 数据处理函数
 * _back, 实现化完成后的回调函数
 * @Jiang youhua
 */
let PART = function () {
    this._id = WEB.Mark();
    this._result;
    //构造函数，可去掉该闭包
    ;
    (function (self, args) {
        let arr = Array.prototype.slice.call(args)
        self.SetArgs(arr)
    })(this, arguments)
}

PART.prototype = {
    SetArgs: function (arr) {
        if (!(arr instanceof Array) || arr.length == 0) {
            return this;
        }
        // 0th, it, 当前使用part的对象
        this.SetIt(a.shift());
        if (a.length == 0) {
            return this;
        }

        // 1th, data, 设置标签
        this.SetHtml(a.shift());
        if (a.length == 0) {
            return this;
        }

        // 2th, data, 数据源解数据后经适匹后的数据
        this.SetData(a.shift());
        if (a.length == 0) {
            return this;
        }

        // 3th, back
        this.SetBack(a.shift());
        if (a.length == 0) {
            return this;
        }
        // 4th, content
        this.SetAttr(a);
        return this;
    },

    /**
     * 需要应用JPART的HMTLElement
     */
    SetIt: function (it) {
        if (!WEB._isDOM(it)) {
            return this;
        }
        this._it = it;
        return this;
    },

    /**
     * 设置数据
     * 判断数据是否正确
     * @param data， 数据
     * @returns {PART}
     * @Jiang youhua
     */
    SetData: function (data) {
        if (!data) {
            return this;
        }
        if(WEB._isDOM(data)){
            return this;
        }
        this._data = data;
        return this;
    },

    /**
     * 设置返回函数
     * @param back
     * @returns {PART}
     * @constructor
     */
    SetBack: function (back) {
        if (!WEB._isString(back)){
            return this;
        }
        this._back = back;
        return this;
    },

    /**
     * 设置设置项
     * @param config
     * @returns {PART}
     * @constructor
     */
    SetConfig: function(config) {
        if (!WEB._isString(config)){
            return this;
        }
        this._config = config;
        return this;
    },

    /**
     * 输出html对象
     * @returns 一个或多个HTML对象
     * @Jiang youhua
     */
    Html: function () {
        // 参数不合格，则用默认参数
        this.forContent();
        this._html.AddAttr("id", this._id)
        return this._html;
    },

    /**
     * 输出
     * @constructor
     */
    Out: function () {
        if (!this._it) {
            return;
        }
        // 将内容处理为
        let s = "";
        let obj = this.Html();
        if (obj instanceof HTML) {
            s = obj.String();
        }
        if (obj instanceof Object) {
            for (let x in obj) {
                if (obj[x] instanceof HTML) {
                    s += obj[x].String();
                }
            }
        }
        this._it.innerHTML = s;
        this.addEvent();
        if (!WEB._isString(this._back)) {
            return;
        }
        try {
            WEB._parseFunc(this._back, this._result);
        } catch (e) {
            console.log(e);
        }
    },

    /** 抽象方法区 start */

    /**
     * 接收数据进行写成HTML内容
     */
    forContent: function () {
        console.log("!!! forContent is not overwrite!");
    },

    addEvent:function () {}

    /** 为组件设置 _data，_config 默认值

     /** 抽象方法区 end */
}

/**
 * URL, 重新的解析URl
 */
let URL = (function () {
    let obj = {};

    // 解析URL search
    let _searchParse = function () {
        return WEB.ObjectByString(window.location.search, "&", "=");
    }

    if (window.location.search != "") {
        obj = _searchParse();
    }

    // 重构URL
    return {
        SearchSet: function (key, value) {
            obj[key] = value;
        },
        SearchGet: function (key) {
            return obj[key];
        },
        Host: function () {
            return window.location.host;
        },
        Path: function () {
            return window.location.pathname;
        },
        Hash: function () {
            return window.location.hash;
        },
        Url: function () {
            return window.location.href;
        }
    }
}());