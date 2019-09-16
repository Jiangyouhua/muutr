/**
 * JPreview, 多图文预览
 * Jiang Youhua 2019.01.10
 */

let JPreview = function() {
    PART.apply(this, arguments);
    (function(self, args) {
        let a = Array.prototype.slice.call(args);
        self.SetArgs(a);
    })(this, arguments);
}

JPreview.prototype = new PART()
JPreview.prototype.checkData = function() {
    // 空值
    if (!this._data) {
        console.log("JPreview Data is nil");
        return false;
    }

    // 不是对象
    if (!(this._data instanceof Object)) {
        console.log("JPreview Data's type is Object");
        return false;
    }
    return true;
}

JPreview.prototype.forContent = function() {
    let number = this._config ? parseInt(this._config) : 4;
    let tr;
    let table = new HTML("table");
    for(let x = 0; x < this._data.length; x ++) {
        let obj = this._data[x];
        if (!obj.image && !obj.text && !obj.title) {
            console.log("JPreview Data's image, text or title is null");
            return;
        }
        // 分行
        if(x % number == 0){
            tr = new HTML("tr");
            table.AddContent(tr)
        }
        let td = new HTML("td");
        td.AddCss("width", (1.0/number*100)+"%");
        tr.AddContent(td);
        // 图片
        if(!!obj.image) {
            let image = new HTML("img");
            image.AddAttr("src", obj.image);
            image.AddClass("image");
            td.AddContent(new HTML("div", image));
        }
        // 标题
        if(!!obj.title) {
            let title = new HTML("div", obj.title);
            title.AddClass("title");
            td.AddContent( title);
        }
        // 文字
        if(!!obj.text) {
            let text = new HTML("div", obj.text);
            text.AddClass("text");
            td.AddContent(text);
        }

        for(let i in obj){
            if(i == "image" || i == "title" || i == "text"){
                continue;
            }
            tr.AddAttr(i, obj[i]);
        }
    }
    this._html = table
}

/**
 * 默认数据
 * title 标题
 * image 图片
 * text 内容
 * title, image, text 三者不能同时为空
 * 其它键值为各项标签属性
 * @type {*[]}
 * @private
 */
JPreview.prototype._data = [
    {
        image:"./examples/ui/image/web1.png",
        title:"Jiang",
        text:"JPreview组件，是用来呈现图文预览， 元素有图片，标题，内容，按从上到下的方式排列，三项元素必须存在一个。"
    },
    {
        image:"./examples/ui/image/web2.png",
        title:"You",

        text:"JPreview组件，是用来呈现图文预览， 元素有图片，标题，内容，按从上到下的方式排列，三项元素必须存在一个。"
    },
    {
        image:"./examples/ui/image/web3.png",
        title:"Hua",
        text:"JPreview组件，是用来呈现图文预览， 元素有图片，标题，内容，按从上到下的方式排列，三项元素必须存在一个。"
    },
];

/**
 * 默认配置
 * 指定分行数量
 * @type {number}
 * @private
 */
JPreview.prototype._config = 3;