/**
 * JTable, 表格
 * Jiangyouhua 32019.01.17
 */

let JTable = function() {
    PART.apply(this, arguments);
    (function(self, args) {
        let a = Array.prototype.slice.call(args)
        self.SetArgs(a)
    })(this, arguments)
}

JTable.prototype = new PART()

JTable.prototype.forContent = function() {
    if(!this._data || this._data.length == 0){
        console.log("JTable Data is null");
        return;
    }
    this._html = new HTML("table");

    let b = WEB._isString(this._config);
    let config = [];
    if(b) {
        config = WEB.ObjectByString(this._config, ",", ":", true);
        b = !!config && config.length > 0;
    }

    let tr = new HTML("tr");
    let head = new HTML("thead", tr);
    this._html.AddContent(head)
    // 表头, 未设置没有表头
    if(b){
        for (let i = 0; i < config.length; i ++) {
            let v = config[i];
            let th = new HTML("th")
            if(v.key == 'checkbox'){
                if(v.value.length > 0) {
                    let check = new HTML("input", "", "type=checkbox");
                    check.AddAttr("name", v.value);
                    th.AddContent(check);
                    tr.AddContent(th);
                }
                continue;
            }
            th.AddContent(v.value);
            tr.AddContent(th);
        }
    }else {
        for (let x in this._data[0]) {
            let th = new HTML("th", x);
            tr.AddContent(th);
        }
    }

    // 表身
    let body = new HTML("tbody");
    this._html.AddContent(body);
    for (let i = 0;  i < this._data.length; i ++) {
        let tr = new HTML("tr");
        body.AddContent (tr);
        let obj = this._data[i];
        // 有参数
        if(b){
            for (let j = 0; j < config.length; j ++) {
                let v = config[j];
                let td = new HTML("td");
                if(v.key == 'checkbox'){
                    if( v.value.length > 0) {
                        let check = new HTML("input", "", "type=checkbox");
                        check.AddAttr("name", v.value);
                        check.AddAttr("value", i);
                        td.AddContent(check);
                        tr.AddContent(td);
                    }
                    continue;
                }
                td.AddContent(obj[v.key]);
                tr.AddContent(td);
            }
        }else{
            // 无参数
            for (let x in obj) {
                let v = obj[x];
                let td = new HTML("td", v);
                tr.AddContent(td);
            }
        }
    }
};

JTable.prototype.addEvent = function(){
    let dom = document.getElementById(this._id);
    if(!dom){
        return;
    }
    let check = dom.querySelector('th input');
    if(!check){
        return;
    }
    check.addEventListener('change', function (e) {
        let b = this.checked;
        let nodes = dom.querySelectorAll("input[type='checkbox']");
        for(let i = 0; i < nodes.length; i++ ){
            nodes[i].checked = b;
        }
    })

};

/**
 * 数据输出为表，没有Config, key为表头
 * @type {{explain: string, name: string, id: number, status: string}[]}
 * @private
 */
JTable.prototype._data = [{
    id:1, name:"Angular", explain:"Google开发的前端框架", status:"0"
},{
    id:1, name:"React", explain:"Facebook开发的前端框架", status:"0"
},{
    id:1, name:"Vue", explain:"国产优秀前端框架", status:"0"
},{
    id:1, name:"jpart", explain:"Jiang Youhua开发的前端框架", status:"1"
}];

/**
 * 确定显示的列，与列的显示名称
 * @type {string}
 * @private
 */
JTable.prototype._config = "checkbox:,id:序号,name:名称,explain:说明,status:状态";