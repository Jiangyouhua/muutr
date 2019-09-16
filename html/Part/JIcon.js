/**
 * JIcon, 站点ICON
 * ** data 
 * 1. {text:string, image:url, other...}
 * 2. text显示的文字, image:显示的图片, 图片在文字前
 * 3. text,image两项必须有一项，不能两项都为空
 * 4. other...将为顶级HTML节点下A标签的属性
 * ** config
 * 1. 无
 */

var JIcon = function() {
    Part.apply(this, arguments);
    (function(self, args) {
        var a = Array.prototype.slice.call(args)
        self.SetArgs(a)
    })(this, arguments)
}

JIcon.prototype = new Part()
JIcon.prototype.checkData = function() {
    // 空值
    if (!this._data) {
        console.log("JIcon Data is nil")
        return false
    }

    // 不是对象
    if (!(this._data instanceof Object)) {
        console.log("JIcon Data's type is Object")
        return false
    }

    // 空数组
    if (!this._data.text || !this._data.image) {
        console.log("JIcon Data's text and image is nil")
        return false
    }
    return true
}

JIcon.prototype.checkConfig = function() {
    console.log("JIcon have not config")
}

JIcon.prototype.forContent = function() {
    var a = new HTML("a")
    if (!!this._data.image) {
        var img = new HTML("img")
        img.SetAttr("src", this._data.image)
        img.AddAttr("align", "absmiddle")
        a.AddContent(img)
    }
    if (!!this._data.text) {
        var span = new HTML("h1", this._data.text)
        span.SetCss("display", "inline")
        span.AddCss("Padding-left", ".5em")
        a.AddContent(span)
    }
    for (var x in this._data) {
        if (x == "text" || x == "image") {
            continue
        }
        a.AddAttr(x, this._data[x])
    }
    this._html = a
}