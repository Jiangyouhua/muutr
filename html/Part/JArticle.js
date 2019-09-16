/**
 * JArticle, 文章预览
 * ** data 
 * 1. {title:string, text:string, href:string, other...}
 * 2. title文章标题, 文章内容
 * 3. title,text不能为空
 * 4. other...为其他信息
 * ** config
 * 1. {keys:[string,], alias:[string,],href:key,limit:int }
 * 2. keys, 确定数据行中哪些字段显示在表格上
 * 3. alias, 该字段在表头上显示的名称
 * 4. href, 确定在哪个字段显示级别, 默认在title
 * 5. limit, 可见的文字数，默认全文显示
 */

var JArticle = function() {
    Part.apply(this, arguments);
    (function(self, args) {
        var a = Array.prototype.slice.call(args)
        self.SetArgs(a)
    })(this, arguments)
}

JArticle.prototype = new Part()
JArticle.prototype.checkData = function() {
    // 空值
    if (!this._data) {
        console.log("JArticle Data is nil")
        return false
    }

    // 不是对象
    if (!(this._data instanceof Array)) {
        console.log("JArticle Data's type is Array")
        return false
    }

    // 空数组
    if (!this._data[0] || !this._data[0].title || !this._data[0].text) {
        console.log("JArticle Data's text and title is nil")
        return false
    }
    return true
}

JArticle.prototype.checkConfig = function() {
    var keys = ["title", "text"]
    var href = "title"
    var limit = 0
    if (!this._config) {
        this._config = { keys: keys, alses: null, href: title, limit: limit }
    }

    // 判断keys, 是否是数组
    if (!(this._config.keys instanceof Array)) {
        this._config.keys = keys
    }
    if (this._config.keys.lenght == 0) {
        var keys = keys
    }

    // 判断alses, 是否是数组并设置
    if (!(this._config.alias instanceof Array) || this._config.alias.lenght == 0) {
        this._config.alias = null
    }

    // 判断sub, 不存在设置为sub, 数据字段没有该项, 则设置为空
    if (typeof this._config.href != "string") {
        this._config.href = href
    }
    var b = false
    for (x in this._data) {
        if (x == this._config.href) {
            b = true
            break
        }
    }
    if (!b) {
        this._config.href = href
    }

    // 判断char, 不合格则设置为默认值
    if (isNaN(this._config.limit)) {
        this._config.limit = limit
    }
    if (this._config.limit < limit) {
        this.cofig.limit = limit
    }
}

JArticle.prototype.forContent = function() {
    var arr = []
    for (var x in this._data) {
        arr.push(this._item(this._data[x]))
    }
    if (arr.length == 0) {
        return
    }
    if (arr.length == 1) {
        this._html = arr[0]
        return
    }
    this._html.AddContent(arr)
}

JArticle.prototype._item = function(data) {
    var html = new HTML()
    for (var x in this._config.keys) {
        var k = this._config.keys[x]
        if (k == "href") {
            continue
        }
        if (!!this._config.alias) {
            var a = this._config.alias[x]
            var p = new HTML("p", a, "class=part-alias")
            html.AddContent(p)
        }
        var v = data[k]
        if (k == "text" && !!this._config.limit) {
            v = this._text(v)
        }
        var d = new HTML("div", v)
        d.AddClass("part-" + k)
        if (!!this._config.href && this._config.href == k) {
            d.AddAttr("href", data.href)
        }
        html.AddContent(d)
    }
    return html
}

JArticle.prototype._text = function(text) {
    text = text.replace(/<.*?(?:>|\/>)/gi, function(word) {
        return ""
    })
    text = text.substring(0, this._config.limit)
    return text + "..."
}