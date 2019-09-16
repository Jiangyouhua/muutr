/**
 * JList, 导航、菜单、列表
 * ** data 
 * 1. [{text:string, image:url, sub:it, other...}, ...]
 * 2. text显示的文字, image:显示的图片, 图片在文字前
 * 3. text,image两项必须有一项，不能两项都为空
 * 4. sub子项，结构与父项相同
 * 5. other...当前项A标签的属性
 * ** config
 * 1. 无
 */

var JList = function() {
    Part.apply(this, arguments);
    (function(self, args) {
        var a = Array.prototype.slice.call(args)
        self.SetArgs(a)
    })(this, arguments)
}

JList.prototype = new Part()
JList.prototype.checkData = function() {
    // 空值
    if (!this._data) {
        console.log("JList Data is nil")
        return false
    }

    // 不是数组
    if (!(this._data instanceof Array)) {
        console.log("JList Data's type is not array")
        return false
    }

    // 空数组
    if (!this._data[0] || (!this._data[0].text && !this._data[0].image)) {
        console.log("JList Data's type is err")
        return false
    }
    return true
}

JList.prototype.checkConfig = function() {
    console.log("JList have not config")
}

JList.prototype.forContent = function() {
    this._html = this._recursion(this._data)
}

JList.prototype._recursion = function(data) {
    if (!data) {
        return
    }
    var ul = new HTML('ul')
    for (var x in data) {
        var obj = data[x]
        var a = new HTML("a")
        if (!!obj.image) {
            var img = new HTML("img")
            img.SetAttr("src", obj.image)
            img.AddAttr("align", "absmiddle")
            a.AddContent(img)
        }
        if (!!obj.text) {
            var span = new HTML("span", obj.text)
            span.SetCss("display", "inline")
            span.AddCss("Padding-left", ".5em")
            a.AddContent(span)
        }
        for (var x in obj) {
            if (x == "text" || x == "image" || x == obj.sub) {
                continue
            }
            a.AddAttr(x, obj[x])
        }
        var li = new HTML('li', a)
        li.AddContent(this._recursion(obj.sub))
        ul.AddContent(li)
    }
    return ul
}