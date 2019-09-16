/**
 * JPaging, 分页组件
 * ** data 
 * 1. {total:int, page:int, other...}
 * 2. total总页数, 
 * 3. page,当前页数
 * 4. other...为其他信息
 * ** config
 * 1. {prev:string, next:string, show:int}
 * 2. prev, 前一页的显示文字
 * 3. next, 后一页的显示文字
 * 2. show, 显示页数, 取值0, 5〜20
 */

var JPaging = function() {
    Part.apply(this, arguments);
    (function(self, args) {
        var a = Array.prototype.slice.call(args)
        self.SetArgs(a)
    })(this, arguments)
}

JPaging.prototype = new Part()
JPaging.prototype.checkData = function() {
    // 空值
    if (!this._data) {
        console.log("JPaging Data is nil")
        return false
    }

    // 不是对象
    if (!(this._data instanceof Object)) {
        console.log("JPaging Data's type is Object")
        return false
    }

    // 空数组
    if (!this._data || !this._data.total || !this._data.page) {
        console.log("JPaging Data's total or page is nil")
        return false
    }
    return true
}

JPaging.prototype.checkConfig = function() {
    if (!this._config) {
        this._config = { prev: "<", next: ">", show: 0 }
    }
    if (typeof this._config.prev != "string" || this._config.prev == "") {
        this._config.prev = "<"
    }
    if (typeof this._config.next != "string" || this._config.next == "") {
        this._config.next = ">"
    }
    // 设置有效值范围
    if (this._config.show > 0 && this._config.show < 2) {
        this._config.show = 5
    }
    if (this._config.show > 20) {
        this._config = 20
    }
}

JPaging.prototype.forContent = function() {
    // 前一页
    var arr = [this._config.prev]

    // 索引部分
    if (!!this._config.show) {
        // 显示页去前与后
        var show = this._config.show - 2
        var start = this._data.page - Math.ceil(show / 2)

        //至少从索2开始
        if (start > this._data.total - show - 1) {
            start = this._data.total - show - 1
        }
        if (start < 2) {
            start = 2
        }

        // 显示最小索引，及后项的省略
        arr.push(1)
        if (start > 2) {
            arr.push("...")
        }
        var i = 0
        for (; i <= show && i < this._data.total; i++) {
            arr.push(start + i)
        }
        // 显示最大索引，及前基的省略
        if (i + start < this._data.total - 1) {
            arr.push("...")
        }
        arr.push(this._data.total)
    }
    // 后一页
    arr.push(this._config.next)
    this._html = this._item(arr)
}

JPaging.prototype._item = function(arr) {
    if (!(arr instanceof Array)) {
        return
    }
    var span = new HTML("span")
    span.AddClass("part-paging")
    for (var x in arr) {
        var a = new HTML("a", arr[x])
        if (arr[x] != "...") {
            a.AddAttr("href", "#")
        }
        span.AddContent(a)
    }
    return span
}