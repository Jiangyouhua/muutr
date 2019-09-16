/**
 * JBanner, 站点横幅
 * Jiang Youhua 2019.01.10
 */

let JBanner = function() {
    PART.apply(this, arguments);
    (function(self, args) {
        let a = Array.prototype.slice.call(args);
        self.SetArgs(a);
    })(this, arguments);
}

JBanner.prototype = new PART();

JBanner.prototype.forContent = function() {
    let div = new HTML();

    if (!this._data.image && !this._data.title && !this._data.text) {
        console.log("JBanner Data's image, title and text is null");
        return;
    }

    if(!!this._data.image) {
        div.AddCss("background-image", "url(" + this._data.image + ")");
        div.AddCss("background-repeat", "no-repeat");
        div.AddCss("background-size", "100% 100%");
        div.AddCss("background-position", "no-center");
    }
    if(!!this._data.title){
        let title = new HTML("div", this._data.title);
        title.AddClass("title")
        div.AddContent(title);
    }
    if(!!this._data.text){
        let text = new HTML("div", this._data.text);
        text.AddClass("text")
        div.AddContent(text);
    }

    for(let i in this._data){
        if(i == "image" || i == "title" || i == "text"){
            continue;
        }
        div.AddAttr(i, this._data[i]);
    }

    this._html = div;
}

/**
 * 默认数据
 * image: 底图
 * title:  显示的标题
 * text: 显示的文字
 * image, title, text 三者不能同时为空
 * 其它标签为各项标签属性
 * @type {*[]}
 * @private
 */
JBanner.prototype._data = {
    image:"",
    title:"Jiang",
    text:"YouHua",
};
