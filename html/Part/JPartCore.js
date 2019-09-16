/**
 * Created by jiangyouhua on 2016/3/31.
 * Update 2017.09.15
 * 自动启用JPart
 * URL, 重新的解析URl
 * HTML, 生成HTML元素字符串的对象
 * Part, 生成HTML组件的对象
 * Web, 实现AJAX加载，
 */

/**
 * 自动启用JPart
 */
window.addEventListener("load",function(){
    // 加载html页面内的组件
    Web.Parse("[data-source]")
    Web.Form()
});

var Web = {
    /**
     * 获取页面宽与高
     */
    Width: window.innerWidth || document.documentElement.clientWidth || document.body.clientWidth,
    Height: window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight,

    /**
     *
     * */
    Url: function (key) {
        var reg = new RegExp("[\\?&]" + key + "=([^&]+)", "gim")
        reg.test(window.location.href)
        return RegExp.$1
    },

    StringToBase64: function (s) {
        return window.btoa(encodeURIComponent(s))
    },

    Base64ToString: function (b) {
        return decodeURIComponent(window.atob(b))
    },

    // 取代默认表单提交与恢复事件
    Form:function () {
        var elements = document.getElementsByTagName("form")
        for (x in elements) {
            elements[x].onsubmit = function (e) {
                // 有没路经
                var param = this.getAttribute("data-form-action")
                if (!param) {
                    return;
                }
                var url = param;
                var t = param.substr(0, 1)
                // 变量，$name
                if (t == "$") {
                    var v = param.substr(1)
                    url = Web._specialWord(v)
                }
                if(t == "@"){
                    url = param.substr(1)
                }
                e.preventDefault();
                // 判断数据是否有效
                var filter = this.getAttribute("data-form-filter")
                if (!!filter && !Web._formatFunc(filter, this)) {
                    console.log("form data is err")
                    return;
                }

                e.preventDefault();
                var data = new FormData(this)
                var router = this.getAttribute("data-form-router")

                // 发送数据
                var it = this
                Web.Post(url, data, function (re) {
                    Web._formatFunc(router, re, it)
                }, "josn")
            }
            elements[x].onreset = function (e) {
                var source = this.getAttribute("data-source")
                if (source.substr(0, 1) != "$") {
                    return
                }
                e.preventDefault();
                var obj = eval(source.substr(1))
                for(x in obj){
                    obj[x] = ""
                }
                Web.Parse(this)
            }
        }x
    },

    /**
     * 显示提示信息
     */
    Alert: function (message, duration) {
        if (!message) {
            return
        }
        if (!duration) {
            duration = 2000
        }
        if (typeof message != "string") {
            return
        }
        var div = new HTML("div", message)
        div.AddAttr("id", "JAlert")
        div.AddCss("position", "absolute")
        div.AddCss("display", "inline-block")
        div.AddCss("width", "26em")
        div.AddCss("background", "#222")
        div.AddCss("color", "#fff")
        div.AddCss("padding", "1em 2em")
        div.AddCss("left", "50%")
        div.AddCss("top", "50%")
        div.AddCss("margin-left", '-15em')
        div.AddCss("margin-top", '-5em')
        div.AddCss("text-align", "center")
        div.AddCss("z-index", "999999")
        var node = div.Element()
        document.body.appendChild(node)
        setTimeout(function () {
            var a = document.getElementById("JAlert")
            document.body.removeChild(a)
        }, duration)
        return;
    },
    /**
     * 本文档的标准错误输出
     * @c, class name
     * @f, function name
     * @n, arg name
     * @a, arg
     * @t, type name
     */
    ArgErr: function (c, f, n, a, t) {
        console.log(c + "." + f + "is Error, typeof " + n + " must is" + a + ", Not type is " + a.toString())
    },

    /**
     * AJAX get请求
     */
    Get: function (url, callBack, classify) {
        var request = Web._request("GET", url)
        // 监听结果
        Web._requestSend(request, null, callBack, classify)
    },

    /**
     * AJAX 请求
     * 1. json；2.js；3.html,
     */
    Post: function (url, data, callBack, classify) {
        var request = Web._request("POST", url)
        if (typeof data == 'object' && !(data instanceof FormData)) {
            data = JSON.stringify(data)
            request.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        }
        Web._requestSend(request, data, callBack, classify)
    },

    _requestSend: function (request, data, callBack, classify) {
        request.onreadystatechange = function () {
            if (request.readyState == 4 && request.status == 200) {
                var obj = Web._response(classify, request.responseText)
                if (!!callBack && typeof callBack === "function") {
                    callBack(obj)
                }
                return false
            }
        }
        request.send(data)
    },

    /**
     * 处理函数，字符串转函数
     * @param func
     * @param data
     * @private
     */
    _formatFunc: function (func, data) {
        var args = Array.prototype.slice.call(arguments);
        args.shift();

        // 不是字符串返回
        if (!(typeof func === 'string')) {
            return data
        }

        // 解析函数
        var start = func.indexOf("(")
        var end = func.lastIndexOf(")")

        // 函数名直接调用
        if (start < 0) {
            return window[func].apply(null, args)
        }

        // 处理参数
        var p = func.substr(start + 1, end - start - 1)
        if (p.length > 0) {
            var as = p.split(",")
            for (var x = as.length - 1; x > -1; x--) {
                var t = as[x]
                if (!t) {
                    continue;
                }
                t = t.trim()
                t = t.replace(/^["']{1}|["']{1}$/g, "");
                args.unshift(t)
            }
        }

        // 处理函数名
        var f = func.substr(0, start)
        var fs = f.split(".")

        // 方法
        var obj = {}
        for (var i = 0; i < fs.length; i++) {
            var k = fs[i];
            k = k.trim()
            if (i == 0) {
                obj = window[k]
                continue
            }
            obj = obj[k]
            if (!obj) {
                return data
            }
        }
        return obj.apply(null, args)
    },

    _request: function (key, url) {
        var request;
        if (window.XMLHttpRequest) { // code for IE7+, Firefox, Chrome, Opera, Safari
            request = new XMLHttpRequest();
        } else { // code for IE6, IE5
            request = new ActiveXObject("Microsoft.XMLHTTP");
        }

        // 设置传递
        request.open(key, url, true);
        return request
    },

    _response: function (classify, text) {
        switch (classify) {
            case "js":
                return eval(text)
            case "html":
                return text
            default:
                try {
                    return JSON.parse(text)
                } catch (e) {
                    return text
                }
        }
    },

    /**
     * 返回全局标签
     * @returns {string}
     */
    Mark: function () {
        var mark = Math.random().toString()
        return mark.replace("0.", "J")
    }
    ,

    /**
     * 从字符串解析对象
     * 1. key:value;key1:value1
     * 2. key=value&key1=value1
     */
    ObjectByString: function (string, split1, split2) {
        var obj = {}
        if (typeof string != "string") {
            return obj
        }
        if (typeof split1 != "string") {
            return obj
        }
        if (typeof split2 != "string") {
            return obj
        }
        var arr = string.split(split1)
        if (!arr) {
            return obj
        }

        var prev = ""
        for (var i = 0; i < arr.length; i++) {
            if (!arr[i]) {
                continue
            }
            var v = arr[i].trim()
            var a = v.split(split2)
            if (!a || a.length == 0) {
                continue
            }
            var key = a[0].trim()
            if (a.length == 1) {
                if (!prev) {
                    obj[key] = false
                } else {
                    obj[prev] += (" " + key)
                }
                continue
            } else {
                prev = key
            }

            obj[key] = a[1].trim()
        }
        return obj
    }
    ,

    /**
     * 解析part对象
     * @param selector，基于data_source实现组件化
     * @constructor
     */
    Parse: function (selector) {
        // 没有选择器
        if (!selector) {
            return
        }
        var b = (typeof selector == "string") || (selector instanceof Element)
        if (!b) {
            return
        }

        // 只对含有JPart的html标签进行组件应用
        var array = []
        if (typeof selector == "string") {
            if (selector.indexOf("[data-source]") < 0) {
                selector += "[data-source]"
            }
            array = document.querySelectorAll(selector)
            Web._parse(array)
            return
        }
        array = selector.querySelectorAll("[data-source]")
        if (array.length > 0) {
            Web._parse(array)
            return
        }
        Web._parse(new Array(selector))
    },

    _parse: function (array) {
        // 没有相关内容
        if (!array || array.length == 0) {
            return
        }

        // 解析当前JPart
        for (var x = 0; x < array.length; x++) {
            var doc = array[x]
            var source = doc.getAttribute("data-source")
            if (!source || source.length == 0) {
                continue
            }
            var config = doc.getAttribute("data-config")
            var func = doc.getAttribute("data-func")
            var callback = doc.getAttribute("data-back")
            var name = doc.getAttribute("data-part")

            // 插入元素
            if (!name) {
                Web._parseDom(doc, source, func, callback)
                continue
            }
            // 解析组件插入
            Web._parsePart(doc, name, source, config, func, callback)
        }
    }
    ,

    _parseJs: function (selector) {
        // 解析对象内容的Js，有js连接文件的加完成后再解析
        var js = selector.getElementsByTagName("script")
        if (!js) {
            return false
        }

        var link = new Array()
        for (var i = 0; i < js.length; i++) {
            var src = js[i].getAttribute("src")
            if (!!src) {
                link.push(src)
            }
            eval(js[i].innerHTML)
        }

        var num = link.length
        if (num == 0) {
            return false
        }
        for (var i = 0; i < link.length; i++) {
            Web.Get(link[i], function () {
                num--
                if (num > 0) {
                    return;
                }
                var array = selector.querySelectorAll("[data-source]")
                Web._parse(array)
            }, "js")
        }
    },

    /**
     * 为DOM 节点处理绑定数据
     * @param doc, 节点
     * @param source, 绑定的数据源
     * @param place, 指定字符串化后的内容插入位置
     * @param back, 回调函数
     * @private
     */
    _parseDom: function (doc, source, func,  back) {
        Web._dataFromParma(source, function (data) {
            if (!data) {
                return
            }
            if(!!func){
                data = Web._formatFunc(func, data)
            }
            Web.SetElementValue(doc, data)
            // 调用回调函数
            if (!back) {
                return
            }
            back()
        })
    }
    ,

    SetElementValue: function (element, value) {
        if (!element) {
            return
        }
        var re = value
        if (typeof re == "object") {
            re = JSON.stringify(re)
        }
        var tag = element.tagName
        if (!!tag) {
            tag = tag.toLowerCase()
        }
        switch (tag) {
            case "form":
                if (value instanceof Object) {
                    for (x in value) {
                        var v = value[x]
                        var d = element.querySelector("[name='" + x + "']")
                        Web.SetElementValue(d, v)
                    }
                }
                break;
            case "input":
                if(element.getAttribute("type") == "file"){
                    break
                }
                element.value = re
                break
            case "select":
                if(!re){
                    element.querySelector("[selected]").selected = false
                }else{
                    element.querySelector("[value='" + re + "']").selected = true
                }
                break
            case "img":
                element.src = re
                break
            default:
                element.innerHTML = re
        }
    },

    /**
     * 为Part 组件输出
     * @param doc， 节点
     * @param name， 组件名称
     * @param source， 绑定的资源
     * @param config， 配置参数
     * @param func， 数据处理函数
     * @param callback， 输出后的回调函数
     * @private
     */
    _parsePart: function (doc, name, source, config, func, callback) {
        var c = Web._dataFromParma(config)
        Web._dataFromParma(source, function (data, inpage) {
            // 判断是否有数据处理函数
            if (data < 1) {
                console.log(data.info)
                return
            }
            // 处理AJAX请求的数
            if (!inpage) {
                data = data.data
            }
            var re = Web._formatFunc(func, data, doc)
            try {
                var part = new (eval(name))
            } catch (e) {
                console.log(e)
                return
            }
            // 字符串化JPart
            part.SetIt(doc).SetData(re).SetConfig(c).SetBack(callback).Out()
        })
    }
    ,


    /**
     * dom 所绑定的参数转为数据
     * 1. 传入非字符串，直接返回或处理
     * 2. 字符串判断为：路径，JSON，变量，其它
     * 3. 路则直接Ajax请求
     * 4. JSON，则转JSON对象
     * 5. 变量，直接获取变量值
     * @param param， 绑定的参数
     * @param callback，回调函数
     * @returns {*}
     * @private
     */
    _dataFromParma: function (param, callback) {
        var b = typeof callback === "function"
        if (typeof param != "string") {
            return !b ? param : callback(param, true)
        }

        var t = param.substr(0, 1)
        // 变量，$name
        if (t == "$") {
            var v = param.substr(1)
            v = Web._specialWord(v)
            return !b ? v : callback(v, true)
        }

        // AJAX 请求
        if (t == "@") {
            var arr = param.split(".")
            var type, classify
            if (arr.length > 0) {
                type = arr[arr.length - 1].toLowerCase()
            }

            switch (type) {
                case "js":
                    classify = "js"
                    break
                case "html":
                    classify = "html"
                    break
            }
            var v = param.substr(1)
            v = Web._parseUrl(v)
            if (!callback) {
                return v
            }
            Web.Get(v, callback, classify)
            return v
        }

        // json object
        if (t == "{" || t == "[") {
            try {
                var v = JSON.parse(param)
                return !b ? v : callback(v, true)
            } catch (e) {
                console.log(e)
                return !b ? param : callback(param, true)
            }
        }

        // 字面值
        return !b ? param : callback(param, true)
    }
    ,

    /**
     * $name类为特别字符，这种变量可取值于变量与表单值
     */
    _specialWord: function (word) {
        try {
            // 转为变量
            return eval(word)
        } catch (e) {
            // 获取表单值
            element = document.querySelector("[name='" + word + "']")
            if (!!element) {
                return element.value
            }

            // 从window.location.search获取值
            return Web.Url(word)
        }
    }
    ,

    /**
     * @ AJAX 处理
     * 1. server/$value/a{$value1}/page.suf, GET
     * 2. server/?key=value&key1=$value, GET
     */
    _parseUrl: function (url) {
        var reg = new RegExp("(\{\$\w+\})|(\$\w+)")
        return url.replace(reg, function (word) {
            var v = word.match(/\w+/)
            v = Web._searchParse(v)
            return v
        })
    }
    ,
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
var HTML = function () {
    //实例属性
    var _tag = "div"
    var _attr = {}
    var _class = []
    var _css = {}
    var _content = null
    var _element = null

        //构造函数，可去掉该闭包
    ;
    (function (self, args) {
        var a = Array.prototype.slice.call(args)
        if (a.length == 0) {
            return
        }
        self.SetTag(a.shift())
        //内容
        if (a.length == 0) {
            return
        }
        self.AddContent(a.shift())
        //属性
        if (a.length == 0) {
            return
        }
        self._attribute(a)
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
        if (typeof name != "string") {
            Web.ArgErr("HTML", "SetTag", "name", "string", name)
            return
        }
        this._tag = name.trim()
        return this
    },

    /**
     * 初化并设置属性，接收键值对与非键值对
     * @param args
     * @param init
     * @returns {HTML}
     * @Jiang youhua
     */
    _attribute: function (args) {
        var _attr = false
        var _class = false
        var _css = false
        // 无参数
        if (!(args instanceof Array) || args.length == 0) {
            return
        }
        for (var i = 0; i < args.length; i++) {
            var value = args[i]
            var obj = {}
            var b = false

            // 判断是否为字符串 "key=value key1=value1", "class=value value1"
            if (typeof value == 'string') {
                var object = Web.ObjectByString(value, " ", "=")
                for (var x in object) {
                    obj[x] = object[x]
                    b = true
                }
                if (!b) {
                    continue
                }
            } else {
                // 判断是否为对象
                if ((typeof value != 'object') || (value instanceof Array)) {
                    continue
                }
                obj = value
            }

            for (var x in obj) {
                var v = obj[x]

                // 样式类型
                if (x == "class") {
                    if (!_class) {
                        this.SetClass(v)
                        _class = true
                        continue
                    }
                    this.AddClass(v)
                    continue
                }
                // 样式分格
                if (x == "style") {
                    var object = Web.ObjectByString(v, ";", ":")
                    for (var x in object) {
                        if (!object[x]) {
                            continue
                        }
                        if (!_css) {
                            this.SetCss(x, object[x])
                            _css = true
                            continue
                        }
                        this.AddCss(x, object[x])
                        continue
                    }
                }

                // 样式属性
                if (!_attr) {
                    this.SetAttr(x, v)
                    _attr = true
                    continue
                }
                this.AddAttr(x, v)
            }
        }
        return this
    },

    /**
     * 设置属性
     */
    SetAttr: function (key, value) {
        this._attr = {}
        return this.AddAttr(key, value)
    },

    /**
     * 添加属性，接收键值对与非键值对
     * @param key
     * @param value
     * @returns {HTML}
     * @Jiang youhua
     */
    AddAttr: function (key, value) {
        if (typeof key != "string") {
            Web.ArgErr("HTML", "AddAttr", "name", "string", key)
            return
        }

        key = key.trim()
        //如果是类别
        if (key == 'class') {
            this.AddClass(value)
            return this
        }

        //如果是风格
        if (key == 'style') {
            var object = Web.ObjectByString(value, ";", ":")
            for (var x in object) {
                this.AddCss(x, object[x])
            }
            return this
        }

        //初始属性
        if (!this._attr) {
            this._attr = {}
        }
        //value为空
        if (!value && value != 0) {
            this._attr[key] = ''
            return this
        }
        if (value instanceof Object) {
            for (var x in value) {
                this.AddAttr(key + "-" + x, value[x])
            }
            return this
        }
        if (typeof value != "string") {
            this._attr[key] = value.trim()
            return this
        }
        this._attr[key] = value.toString()
        return this
    },

    /**
     * 初始化并设置类型
     * @param args，接收字符串与数组
     * @returns {HTML}
     * @Jiang youhua
     */
    SetClass: function () {
        //空值
        this._class = {key: {}, value: []}
        var a = Array.prototype.slice.call(arguments);
        return this._classByArray(a)
    },

    /**
     * 添加类型
     * @param value，类型名
     * @returns {HTML}
     * @Jiang youhua
     */
    AddClass: function () {
        var a = Array.prototype.slice.call(arguments);
        return this._classByArray(a)
    },

    /**
     * 接收数组，转为式样类型
     */
    _classByArray: function (a) {
        if (!this._class) {
            this._class = {key: {}, value: []}
        }

        if (!(a instanceof Array)) {
            return
        }
        if (a.length == 0) {
            return
        }
        for (var i = 0; i < a.length; i++) {
            var val = a[i]
            if (!val) {
                continue
            }
            // 如果是数组
            if (val instanceof Array) {
                this._classByArray(val)
                continue
            }
            // 非字符串
            if (typeof val != "string") {
                continue
            }
            //字符串
            var arr = val.split(" ")
            if (arr.length > 1) {
                this._classByArray(arr)
                continue
            }

            var v = val.trim()
            if (this._class.key[v]) {
                continue
            }
            this._class.key[v] = true
            this._class.value.push(v)
        }
    },

    // 获取样式类型
    GetClass: function () {
        if (!this._class || !this._class.value) {
            return ""
        }
        return this._class.value.join(" ")
    },

    // 设置样式属性，重置样式为空后再添加样式
    SetCss: function (key, value) {
        this._css = {}
        return this.AddCss(key, value)
    },

    // 添加样式属性
    AddCss: function (key, value) {
        if (typeof key != "string") {
            return this
        }
        if (typeof value != "string") {
            return this
        }
        if (!this._css) {
            this._css = {}
        }
        this._css[key] = value
        return this
    },

    // 获取样式属性
    GetCss: function (key) {
        if (typeof key != "stirng") {
            return
        }
        if (!this._css) {
            return
        }
        return this._css[key]
    },

    /**
     * 初始化并设置内容
     * @param contents，内容，支持所有对象
     * @returns {HTML}
     * @Jiang youhua
     */
    SetContent: function () {
        this._content = []
        var arr = Array.prototype.slice.call(arguments)
        return this._contentByArray(arr)
    },

    /**
     * 添加内容
     * @param content，内容，支持所有对象
     * @returns {HTML}
     * @Jiang youhua
     */
    AddContent: function () {
        var arr = Array.prototype.slice.call(arguments)
        return this._contentByArray(arr)
    },

    _contentByArray: function (arr) {
        //空值
        if (!(arr instanceof Array) || arr.length == 0) {
            return this
        }
        for (var x in arr) {
            var v = arr[x]
            if (this == v) {
                return this
            }
            //初始化
            if (!this._content) {
                this._content = []
            }
            this._content.push(v)
        }
        return this
    },

    /**
     * 格式化内容为字符串
     * @param contents，内容
     * @returns {*}
     */
    _forContent: function (contents) {
        //非对象，转字符串返回
        if (!contents && contents != 0) {
            return ''
        }
        if (typeof contents != 'object') {
            return contents.toString()
        }
        //html对象，转字符串返回
        if (contents instanceof HTML) {
            return contents.String()
        }
        //对象，遍历递归调用
        s = ''
        for (var x in contents) {
            s += this._forContent(contents[x])
        }
        return s
    },

    /**
     * HTML对象输出文档对象节点
     * @return {dom.Element}
     */
    Element: function () {
        var div = document.createElement("div");
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
            this._tag = "div"
        }
        // 处理类别
        var c = ''
        if (!!this._class && !!this._class.value) {
            c = this._class.value.join(" ")
            c = c.trim()
            if (!!c) {
                c = 'class="' + c + '"'
            }
        }
        // 处理属性
        var a = ''
        for (var x in this._attr) {
            if (!this._attr[x]) {
                a = a + " " + x //单词属性
            } else {
                a = a + ' ' + x + '="' + this._attr[x] + '"'; //key=value
            }
        }
        // 处理风格
        var s = ''
        for (var x in this._css) {
            if (!this._css[x]) {
                continue
            }
            if (!s) {
                s += 'style="'
            }
            var v = x + ":" + this._css[x]
            s = s + v + ";"
        }
        if (!!s) {
            s += '"'
        }

        // 非内容部分
        a = a.trim()
        var h = "<" + this._tag + " " + c + " " + a + " " + s
        switch (this._tag) {
            case "img":
            case "input":
            case "br":
            case "hr":
                return h + "/>"
        }
        var str = this._forContent(this._content)
        return h + ">" + str + "</" + this._tag + ">"
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
var Part = function () {
    this._content = [];
    this._isforContent = false;
    //构造函数，可去掉该闭包
    ;
    (function (self, args) {
        var arr = Array.prototype.slice.call(args)
        self.SetArgs(arr)
    })(this, arguments)
}

Part.prototype = {
    SetArgs: function (arr) {
        if (!(arr instanceof Array) || arr.length == 0) {
            return this
        }
        // 0th, it, 当前使用part的对象
        this.SetIt(a.shift())
        if (a.length == 0) {
            return this
        }

        // 1th, data, 设置标签
        this.SetHtml(a.shift())
        if (a.length == 0) {
            return this
        }

        // 2th, data, 数据源解数据后经适匹后的数据
        this.SetData(a.shift())
        if (a.length == 0) {
            return this
        }
        // 3th, config, 配置参数
        this.SetConfg(a.shift())
        if (a.length == 0) {
            return this
        }

        // 4th, back
        this.SetBack(a.shift())
        if (a.length == 0) {
            return this
        }
        // 5th, content
        this.SetAttr(a)
        return this
    },

    /**
     * 需要应用JPart的HMTLElement
     */
    SetIt: function (it) {
        if (!it || !it.nodeType || it.nodeType != 1) {
            return this
        }
        this._it = it
        return this
    },

    /**
     * 设置数据
     * 判断数据是否正确
     * @param data， 数据
     * @returns {Part}
     * @Jiang youhua
     */
    SetData: function (data) {
        if (!data) {
            return this
        }
        this._data = data
        return this
    },

    /**
     * 设置配置参数
     * @param config， 参数
     * @returns {Part}
     * @Jiang youhua
     */
    SetConfig: function (config) {
        if (!config) {
            return this
        }
        this._config = config
        return this
    },

    /**
     * 设置返调孙函数
     */
    SetBack: function (back) {
        if (!back) {
            return this
        }
        if (typeof back != "string") {
            return this
        }
        this._back = back
        return this
    },

    /**
     * 设置组件内容
     */
    SetContent: function () {
        this._content = []
        var arr = Array.prototype.slice.call(arguments)
        return AddContent(arr)
    },

    /**
     * 添加组件内容
     */
    AddContent: function (content) {
        var arr = Array.prototype.slice.call(arguments)
        return AddContent(arr)
    },

    _contentByArray: function (arr) {
        if (!(arr instanceof Array) || arr.length == 0) {
            return this
        }
        if (!this._contnet) {
            this._contnet = []
        }
        for (var x in arr) {
            var v = arr[x]
            if (typeof v == 'string' || v instanceof HTML) {
                this._content.push(content)
                continue
            }
            if (typeof content == 'object') {
                this._contentByArray(v)
            }
        }
        return this
    },

    /**
     * 输出html对象
     * @returns 一个或多个HTML对象
     * @Jiang youhua
     */
    Html: function () {
        // 数据不合格则返回空
        if (!this.checkData()) {
            return new HTML()
        }
        // 参数不合格，则用默认参数
        this.checkConfig()
        this.forContent()
        return this._html
    },

    /**
     * 输出
     */
    Out: function () {
        if (!this._it) {
            return
        }
        // 将内容处理为
        var s = ""
        var obj = this.Html()
        if (obj instanceof HTML) {
            s = obj.String()
        }
        if (obj instanceof Object) {
            for (x in obj) {
                if (obj[x] instanceof HTML) {
                    s += obj[x].String()
                }
            }
        }
        this._it.innerHTML = s
        Web._formatFunc(this._back, this._data, this._it)
    },

    /** 抽象方法区 start */

    /**
     * 检查数据是否合格
     */
    checkData: function () {
        console.log("!!! checkData is not overwrite!")
        return false
    },

    /**
     * 检查配置文件是否合格，不合格用默认值
     */
    checkConfig: function () {
        console.log("!!! checkConfig is not overwrite")
    },

    /**
     * 接收数据进行写成HTML内容
     */
    forContent: function () {
        console.log("!!! forContent is not overwrite!")
    },

    /** 抽象方法区 end */
}
