/**
 * JList, 导航、菜单、列表
 *  * Jiang Youhua 2019.01.10
 */

let JList = function() {
    PART.apply(this, arguments);
    (function(self, args) {
        let a = Array.prototype.slice.call(args)
        self.SetArgs(a)
    })(this, arguments)
}

JList.prototype = new PART()

JList.prototype.forContent = function() {
    this._html = this._recursion(this._data)
}

JList.prototype._recursion = function(data) {
    if (!data) {
        return
    }
    let ul = new HTML('ul')
    for (let x in data) {
        let obj = data[x]
        if(!obj.text){
            console.log("JList Data's text is null");
            return;
        }
        let a = new HTML("a")
        if (!!obj.icon) {
            let arr = obj.icon.split(".")
            let suffix = arr.length > 1 ? arr[length - 1] : null;
            if(!suffix){
                a.AddContent(obj.icon);
            } else {
                let span = new HTML("span", "&nbsp;&nbsp;&nbsp;&nbsp;");
                span.AddCss("background-image", "url("+obj.icon+")");
                span.AddCss("background-repeat", "no-repeat");
                span.AddCss("background-size", "100% 100%");
                span.AddCss("background-position", "no-center");
                a.AddContent(span);
            }
        }
        if (!!obj.text) {
            let span = new HTML("span", obj.text)
            a.AddContent(span)
        }
        if (!obj.href){
            obj.href = "#";
        }
        for (let i in obj) {
            if (i == "text" || i == "icon" || i == "sub") {
                continue;
            }
            a.AddAttr(i, obj[i])
        }
        let li = new HTML('li', a)
        li.AddContent(this._recursion(obj.sub))
        ul.AddContent(li)
    }
    return ul
}

/**
 * 默认数据
 * text, not null 项目文字
 * icon, 项目图标，1.image, 2. icon class : <span class="book"><span>
 * href, 跳转链接
 * sub, 子数据
 * 其它标签为各项标签属性
 */
JList.prototype._data = [{
    text: "一级列表一",
    sub:[{
        text: "二一级列表一",
    }, {
        text: "二一级列表二",
    }, {
        text: "二一级列表三",
    }],
}, {
    text: "一级列表二",
    sub:[{
        text: "二二级列表一",
    }, {
        text: "二二级列表二",
    }, {
        text: "二二级列表三",
    }],
}, {
    text: "一级列表三",
    sub:[{
        text: "二三级列表一",
    }, {
        text: "二三级列表二",
    }, {
        text: "二三级列表三",
    }],
}]