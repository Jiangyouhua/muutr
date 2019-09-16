/**
 * JTable, 表格
 * ** data 
 * 1. [{key:value, sub:it, other...}, ...]
 * 2. key任意项，至少有一项
 * 3. sub子项，结构与父项相同
 * ** config
 * 1. {keys:[string,], alias:[string,],sub:key,char:string }
 * 2. keys, 确定数据行中哪些字段显示在表格上
 * 3. alias, 该字段在表头上显示的名称
 * 4. sub, 确定在哪个字段显示级别
 * 5. char, 通过该字符以缩进的方式显示
 */

var JTable = function() {
    Part.apply(this, arguments);
    (function(self, args) {
        var a = Array.prototype.slice.call(args)
        self.SetArgs(a)
    })(this, arguments)
}

JTable.prototype = new Part()
JTable.prototype.checkData = function() {
    // 空值
    if (!this._data) {
        console.log("JTable Data is nil")
        return false
    }

    // 不是数组
    if (!(this._data instanceof Array)) {
        console.log("JTable Data's type is array")
        return false
    }
    return true
}

JTable.prototype.checkConfig = function() {
    if (!this._config) {
        this._config = { keys: [], alias: [], sub: "", char: "--" }
        console.log("JTable Config is nil")
    }

    // 判断keys, 是否是数组
    if (!(this._config.keys instanceof Array)) {
        this._config.keys = []
    }
    // 判断alias,
    if(!(this._config.alias instanceof Array)){
        this._config.alias = []
    }

    // 判断sub,
    if (typeof this._config.sub != "string") {
        this._config.sub = ""
    }
    // 判断char
    if(typeof this._config.char != "string"){
        this._config.char = "--"
    }

    var b = false
    if (this._config.keys.length == 0) {
        for (x in this._data[0]) {
            this._config.keys.push(x)
            if (x  == this._config.sub){
                b = true
            }
        }
    }

    if (!b) {
        this._config.sub = ""
    }
}

JTable.prototype.forContent = function() {
    this._html = new HTML("table")
    this._header()
    this._body = new HTML("tbody")
    this._recursion(this._data, 0)
    this._html.AddContent(this._body)
}

// 表头
JTable.prototype._header = function() {
    var tr = new HTML("tr")
    var header = new HTML("thead", tr)
    for (var x in this._config.keys) {
        var v = !!this._config.alias && !!this._config.alias[x] ? this._config.alias[x] : this._config.keys[x]
        var th = new HTML("th", v)
        th.SetClass("part-" + this._config.keys[x])
        tr.AddContent(th)
    }
    this._html.AddContent(header)
}

// 表身
JTable.prototype._recursion = function(data, level) {
    for (var x in data) {
        var obj = data[x]
        var tr = new HTML("tr")
        this._body.AddContent(tr)
        for (var i in this._config.keys) {
            var v = obj[this._config.keys[i]]
                // 添加子类缩进
            if (!!this._config && this._config.sub && i == this._config.sub) {
                var prefix = ""
                for (var j = 0; j < level; j++) {
                    prefix += this._config.char
                }
                v = prefix + v
            }
            var td = new HTML("td", v)
            td.SetClass("part-" + this._config.keys[i])
            tr.AddContent(td)
        }
        this._recursion(obj.sub, level + 1)
    }
}