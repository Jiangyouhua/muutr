/**
 * JNavigation, 站点导航
 * Jiang Youhua 2019.01.10
 */

let JNavigation = function() {
    PART.apply(this, arguments);
    (function(self, args) {
        let a = Array.prototype.slice.call(args);
        self.SetArgs(a);
    })(this, arguments);
}

JNavigation.prototype = new PART()

JNavigation.prototype.forContent = function() {
    let div = new HTML();
    for(let x = 0; x < this._data.length; x ++) {
        let obj = this._data[x];
        if (!obj.text) {
            console.log("JNavigation Data's text is null");
            return;
        }

        let a = new HTML("a", obj.text);
        if(x == parseInt(this._config)){
            a.AddClass("action");
        }
        for(let i in obj){
            if( i == "text"){
                continue;
            }
            a.AddAttr(i, obj[i]);
        }
        a.AddCss("width", 100 / this._data.length + "%");
        div.AddContent(a);
    }
    this._html = div;
    this._result = this._id;
}

/**
 * 默认数据
 * text: not null, 菜单文字
 * href: not null, 跳转连接
 * 其它键值为各项标签属性
 * data-*，需要使用引号
 * @type {*[]}
 * @private
 */
JNavigation.prototype._data = [
    {text:"Jiang", href:"#", onclick:"alert(1)", "data-affect":"#id"},
    {text:"You", href:"#", onclick:"alert(2)", "data-affect":"#id"},
    {text:"Hua", href:"#", onclick:"alert(3)", "data-affect":"#id"},
];

/**
 * 默认配置
 * 该项为选择状
 * @type {number}
 * @private
 */
JNavigation.prototype._config = 0;