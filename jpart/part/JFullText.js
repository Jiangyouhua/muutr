/**
 * JFullText, 全文阅读
 * Jiang Youhua 2019.01.10
 */

let JFullText = function() {
    PART.apply(this, arguments);
    (function(self, args) {
        let a = Array.prototype.slice.call(args);
        self.SetArgs(a);
    })(this, arguments);
}

JFullText.prototype = new PART();

JFullText.prototype.forContent = function() {
    if(!this._data.title || !this._data.text){
        console.log("JFullText Data's title and text is null");
        return;
    }

    let div = new HTML();
    let number = parseInt(this._config);
    let a = ["title", "subtitle", "author", "date", "text", "publisher"];
    for(let i = 0; i < a.length; i ++){
        let k = a[i];
        if(!this._data[k]){
            continue;
        }
        let d = new HTML("div", this._data[k])
        if(number > 1 && k == "text"){
            d.AddCss("column-count", number);
            d.AddCss("column-gap", "3rem");
        }
        d.AddClass(k);
        div.AddContent(d)
    }

    for(let i in this._data){
        if(i == "title" || i == "subtitle" || i == "author" || i == "date" || i == "text" || i == "publisher"){
            continue;
        }
        div.AddAttr(i, this._data[i]);
    }
    this._html = div;
}

/**
 * 默认数据
 * title, not null 标题
 * subtitle, 副标题
 * text, not null 内容
 * author, 作者
 * publisher, 出版者
 * date, 时间
 * 其它标签为各项标签属性
 */
JFullText.prototype._data = {
    title:"JFullText, 标题",
    subtitle:"JFullText, 内容",
    text:"JFullText, 内容",
};

/**
 * 默认配置
 * 分栏数
 * @type {number}
 * @private
 */
JFullText.prototype._config = 1;