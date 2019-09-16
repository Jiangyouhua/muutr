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

let JPagination = function() {
    PART.apply(this, arguments);
    (function(self, args) {
        let a = Array.prototype.slice.call(args)
        self.SetArgs(a)
    })(this, arguments)
}

JPagination.prototype = new PART()

JPagination.prototype.forContent = function() {
    // 前一页
    let arr = [this._config.prev]

    // 索引部分
    if (!!this._config.show) {
        // 显示页去前与后
        let show = this._config.show - 2
        let start = this._data.page - Math.ceil(show / 2)

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
        let i = 0
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

JPagination.prototype._item = function(arr) {
    if (!(arr instanceof Array)) {
        return
    }
    let span = new HTML("span")
    span.AddClass("part-paging")
    for (let x in arr) {
        let a = new HTML("a", arr[x])
        if (arr[x] != "...") {
            a.AddAttr("href", "#")
        }
        span.AddContent(a)
    }
    return span
}