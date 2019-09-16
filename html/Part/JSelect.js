/**
 * JSelect, 下拉列表
 * ** data
 * 1. [{text:string, value:url}, ...]
 * 2. text显示的文字, value:对应的值
 * ** config
 * 1. noting
 */

var JSelect = function() {
    Part.apply(this, arguments);
    (function(self, args) {
        var a = Array.prototype.slice.call(args)
        self.SetArgs(a)
    })(this, arguments)
}

JSelect.prototype = new Part()
JSelect.prototype.checkData = function() {
    // 空值
    if (!this._data) {
        console.log("JSelect Data is nil")
        return false
    }

    // 不是数组
    if (!(this._data instanceof Array)) {
        console.log("JSelect Data's type is not array")
        return false
    }

    // 空数组
    if (!this._data[0] || (!this._data[0].text && !this._data[0].value)) {
        console.log("JSelect Data's type is err")
        return false
    }
    return true
}

JSelect.prototype.checkConfig = function() {
    console.log("JSelect Config is noting")
}

JSelect.prototype.forContent = function() {
    this._html = this._recursion(this._data)
}

JSelect.prototype._recursion = function(data) {
    if (!data) {
        return
    }
    var arr = new Array()
    for (var x in data) {
        var obj = data[x]
        var option = new HTML("option")
        option.AddAttr("value", obj.value)
        option.AddContent(obj.text)
        arr.push(option)
    }
    return arr
}