    // get_post测试
    function get_post() {
        axios({
                url: 'http://127.0.0.1:8090/v1/user',
                data: {
                    "user_id": 10000,
                },
                method: 'post',
                responseType: 'json',
            })
            .then(function (response) {
                console.log(response.data);
            })
            .catch(function (error) {
                console.log(error);
            });

        axios({
                method: 'get',
                url: 'http://127.0.0.1:8090/v1/user/1',
                responseType: 'json'
            })
            .then(function (response) {
                console.log(response.data);
            })
            .catch(function (error) {
                console.log(error);
            });

    }
    // 显示子菜单 0
    function show_sub_tab() {
        var lis = document.getElementById("sub")
        var lt = document.getElementById("left")
        var nv = document.getElementById("nav")
        lt.className = "col-1 border-end border-primary border-3 "
        lis.className = "col-1 border-end border-warning border-3"
        nv.className = "nav flex-column nav-pills"
    }
    // 获取账户管理列表子菜单 1
    function get_account_list() {
        show_sub_tab()
        document.getElementById("show").innerHTML =        
        `<table id="zhgl_list" class="table table-bordered table-hover d-none">
            <thead>
                <tr>
                    <th scope="col">序号</th>
                    <th scope="col">账户</th>
                    <th scope="col">身份ID</th>
                    <th scope="col">名称</th>
                    <th scope="col">激活时间</th>
                    <th scope="col">到期时间</th>
                    <th scope="col">绑定</th>
                    <th scope="col">解除</th>
                </tr>
            </thead>
            <tbody id="zhgllist">
                <!-- 账户列表 -->
            </tbody>
        </table>
        <nav aria-label="Page navigation example">
            <ul class="pagination justify-content-end" id="barcon"></ul>
        </nav>`
    }

    // 账户管理子菜单：获取管理员、座席、终端的账户数量
    function get_accounts(type) {
        document.getElementById("zhgl_list").className = "table table-bordered table-hover text-center"
        goPage(1)
        var topic = "enterprise1." + type.id
        axios({
                method: 'get',
                url: 'http://127.0.0.1:8002/n_node/v1.0/kv/seq',
                params: {
                    "topic": topic,
                    "keys": ""
                },
                responseType: 'json'
            })
            .then(function (response) {
                console.log(response.data);
                var value = response.data.data
                if (value) {
                    value = value
                } else {
                    value = 0
                }
                document.getElementById("zhgllist").innerHTML = ''
                for (let index = 0; index < value.length; index++) {
                    var obj = value[index]["value"]
                    var key = value[index]["key"]
                    document.getElementById("zhgllist").innerHTML +=
                        `<tr>
                    <th scope="row">${index+1}</th>
                    <td>${key}</td>
                    <td>${obj["finger"]}</td>
                    <td>${obj["name"]}</td>
                    <td>${obj["active_time"]}</td>
                    <td>${obj["end_time"]}</td>
                    <td><button type="button" class="btn btn-primary btn-sm" data-bs-toggle="modal"
                    data-bs-target="#info2account" id="${key}" name="managers" onclick="get_account(this)">绑定</button></td>
                    <td>
                    <button type="button" class="btn btn-danger btn-sm" id="${key}" onclick="unbinding(this)">解除</button>
                    </td>
                    </tr>`
                }
            })
            .catch(function (error) {
                console.log(error);
            });
    }

    // 账户配置子菜单：正文绑定按钮：获取绑定账户信息
    function get_account(topic) {

        axios({
                method: 'get',
                url: 'http://www.51tianyue.cn:20004/n_account/v1.0/device',
                // params: {
                //     "user": topic.name,
                //     "finger": topic.name,
                // },
                responseType: 'json'
            })
            .then(function (response) {
                console.log(response.data);
                var value = response.data.data
                if (value) {
                    value = value
                } else {
                    value = 0
                }
                document.getElementById("model-info").innerHTML = ''
                for (let index = 0; index < value.length; index++) {
                    var obj = value[index]
                    document.getElementById("model-info").innerHTML +=
                        `<tr>
                    <th scope="row">${index+1}</th>
                    <td>${obj["nick"]}</td>
                    <td>${obj["user"]}</td>
                    <td>${obj["finger"]}</td>
                    <td><input class="form-check-input" type="radio"  id="${topic.id}" name="info" data="${topic.name}" value="${obj["finger"]}"></td>
                    </tr>`
                }
                var info2account = document.getElementById('info2account')

                info2account.addEventListener('shown.bs.modal', function (event) {
                
                    console.log("account modal opend")
                })
                info2account.addEventListener('hidden.bs.modal', function (event) {
                    console.log("account modal closed")
                    // get_accounts()
                })
            })
            .catch(function (error) {
                console.log(error);
            });
    }

    // 账户管理子菜单 INFO Modal框确定按钮
    function binding() {
        var infos = document.getElementsByName("info")
        for (v in infos) {
            if (infos[v].checked) {
                var name = infos[v].parentNode.previousElementSibling.previousElementSibling.innerText
                var users = infos[v].parentNode.previousElementSibling.previousElementSibling.previousElementSibling.innerText
                var finger = infos[v].parentNode.previousElementSibling.innerText
            }
        }

        var info_name = infos[0].attributes.data.value

        var type
        if (info_name == "managers") {
            type = 1
        } else if (info_name == "seat") {
            type = 2
        } else {
            type = 3
        }

        var data = {
            "data": {
                "enterprise_id": "enterprise1",
                "finger": trim(finger), //指纹（对于终端是硬件指纹，对于坐席是权限系统分配） 
                "name": trim(name), //显示的名称
                "user": trim(users), //登录的用户名
                "accountID": trim(infos[0].id), //绑定的账户ID
                "type": type //1管理员 2坐席 3终端
            }
        }
        console.log("绑定数据：", data)

        axios({
                method: 'post',
                url: 'http://127.0.0.1:8002/d_auth/v1.0/bindAccount',
                data: data,
                responseType: 'json'
            })
            .then(function (response) {
                console.log(response.data);
                var infomodal = document.getElementById('info2account')
                infomodal.addEventListener('shown.bs.modal', function (event) {
                    console.log("account modal opend")
                })
                infomodal.addEventListener('hidden.bs.modal', function (event) {
                    console.log("account modal closed")
                    // get_accounts()
                })

            })
            .catch(function (error) {
                console.log("error：" + error);
            });
    }

    // 账户管理子菜单：正文表格：解绑
    function unbinding(id) {
        console.log("正在解除绑定")
        var data = {
            "data": {
                "enterprise_id": "enterprise1",
                "accountID": id.id,
                "type": 1
            }
        }
        axios({
                method: 'delete',
                url: 'http://127.0.0.1:8002/d_auth/v1.0/unbindAccount',
                data: data,
                responseType: 'json'
            })
            .then(function (response) {
                console.log(response.data);
            })
            .catch(function (error) {
                console.log(error);
            });
    }

    // 座席分配菜单：获取子菜单座席数量列表
    function get_seat_list() {
        show_sub_tab()
        document.getElementById("show").innerHTML = ""
        axios({
                method: 'get',
                url: 'http://127.0.0.1:8002/n_node/v1.0/kv/seq',
                params: {
                    "topic": "enterprise1.seats",
                    "keys": ""

                },
                responseType: 'json'
            })
            .then(function (response) {
                console.log(response.data);
                var value = response.data.data
                if (value) {
                    value = value
                } else {
                    value = 0
                }
                var zxlist = document.getElementById("list-tab zx")
                zxlist.innerHTML = ""
                for (let index = 0; index < value.length; index++) {
                    var obj = value[index]["value"]
                    var key = value[index]["key"]
                    zxlist.innerHTML +=
                        `
                    <a class="list-group-item list-group-item-action" id="${obj["finger"]}" 
                    data-bs-toggle="list" href="#" role="tab" onclick="get_seat_devs(this)">
                    <i class="bi bi-headset"></i> 座席${index+1}
                    </a>`
                }

            })
            .catch(function (error) {
                console.log(error);
            });


    }
    // 座席分配子菜单：获取指定座席下终端数量
    function get_seat_devs(id) {
        var topic = ""
        if (id.id) {
            var topic = id.id + ".devs"
        }
        console.log("seat_id:" + id)
        axios({
                method: 'get',
                url: 'http://127.0.0.1:8002/n_node/v1.0/kv/seq',
                params: {
                    "topic": topic,
                    "keys": ""
                },
                responseType: 'json'
            })
            .then(function (response) {
                console.log(response.data);
                var value = response.data.data
                if (value) {
                    value = value
                } else {
                    value = 0
                }
                document.getElementById("show").innerHTML =
                    `                    
                <div id="zxfp" class="d-block">
                    <ul class="nav justify-content-between align-middle text-center">
                        <li class="nav-item align-middle text-center">
                            <p id="count_zxzd" class="h3">终端个数：</p>
                        </li>
                        <li class="nav-item">
            
                            <button type="button" class="btn btn-primary" name="fp" data-bs-toggle="modal"
                                data-bs-target="#dev2seat" id="${id.id}"  value = "${id.id}" onclick="get_dev_list(this)">分配</button>
                        </li>
                    </ul>
                    <hr class="border-bottom border-danger" />
                    <table class="table table-light table-bordered table-hover text-center">
                        <thead>
                            <tr>
                                <th scope="col">序号</th>
                                <th scope="col">终端识别</th>
                                <th scope="col">终端名称</th>
                            </tr>
                        </thead>
                        <tbody id="zxfplist">
                            <!-- 终端列表 -->
                        </tbody>
                    </table>

                </div>
                `
                var count_zxzd = document.getElementById("count_zxzd")
                count_zxzd.innerText = "终端个数：" + value.length
                for (let index = 0; index < value.length; index++) {
                    var obj = value[index]["value"]
                    var key = value[index]["key"]
                    document.getElementById("zxfplist").innerHTML +=
                        `<tr>
                        <th scope="row">${index+1}</th>
                        <td>${key}</td>
                        <td>${obj["name"]}</td>
                    </tr>`
                }
            })
            .catch(function (error) {
                console.log(error);
            });

    }
    // 座席分配子菜单:正文分配按钮获取所有终端列表
    function get_dev_list(id) {
        axios({
                method: 'get',
                url: 'http://127.0.0.1:8002/n_node/v1.0/kv/seq',
                params: {
                    "topic": "enterprise1.devs",
                    "keys": ""
                },
                responseType: 'json'
            })
            .then(function (response) {
                console.log(response.data);
                var value = response.data.data
                if (value) {
                    value = value
                } else {
                    value = 0
                }
                document.getElementById("model-devs").innerHTML = ''
                for (let index = 0; index < value.length; index++) {
                    var obj = value[index]["value"]
                    var key = value[index]["key"]
                    document.getElementById("model-devs").innerHTML +=
                        `<tr>
                        <th scope="row">${index+1}</th>
                        <td>${obj["finger"]}</td>
                        <td>${obj["name"]}</td>
                        <td><input class="form-check-input" type="checkbox"  id="seatfinger${index+1}" name="seat-finger" data="${id.id}" value="123"></td>
                    </tr>`
                }
                var myModalEl = document.getElementById('dev2seat')

                myModalEl.addEventListener('shown.bs.modal', function (event) {
                    checkall()
               
                    console.log("de modal open")
                })
                myModalEl.addEventListener('hidden.bs.modal', function (event) {
                    get_seat_devs(id)
                    console.log("dev modal close")
                })
            })
            .catch(function (error) {
                console.log(error);
            });
    }

    // 座席分配子菜单 DEV Modal框确定按钮
    function devstoseat() {
        var seatfinger = document.getElementsByName("seat-finger")
        var dev_fingers = Array()
        var dev = ""
        for (v in seatfinger) {
            if (seatfinger[v].checked) {
                dev = seatfinger[v].parentNode.previousElementSibling.previousElementSibling.innerText,
                    dev_fingers.push(dev)
            }
        }
        var seat_finger = seatfinger[0].attributes.data.value

        var data = {
            "data": {
                "enterprise_id": "enterprise1",
                "seat_finger": seat_finger,
                "dev_fingers": dev_fingers
            }
        }
        console.log(data)
        axios({
                method: 'post',
                url: 'http://127.0.0.1:8002/d_auth/v1.0/devsToSeat',
                data: data,
                responseType: 'json'
            })
            .then(function (response) {
                console.log(response.data);
                var devmodal = document.getElementById('dev2seat')
                devmodal.addEventListener('hidden.bs.modal', function (event) {
                    console.log("modal closed")
                })
            })
            .catch(function (error) {
                console.log(error);
            });
    }




    // -----------以下未使用------------------
    function get_addauth_list() {
        show_sub_tab()
    }

    function info2account() {
        axios({
                method: 'post',
                url: 'http://127.0.0.1:8002/d_auth/v1.0/authInfo',
                data: {
                    "enterprise_id": "enterprise1",
                    "enterprise": "xxx",
                    "admin_num": 20,
                    "dev_num": 2,
                    "seat_num": 2,
                    "active_time": 1616941977,
                    "end_time": 1617084030
                },
                responseType: 'json'
            })
            .then(function (response) {
                console.log(response.data);
            })
            .catch(function (error) {
                console.log(error);
            });
    }

    function checkall() {
        var cekall = document.getElementById("allcheck")
        var inp = document.querySelector('#model-devs').getElementsByTagName('input');

        // 注册事件
        // cekall.addEventListener("click",function(){
        //     console.log(this.checked);
        //     for (var i = 0; i < inp.length; i++) {
        //         inp[i].checked = this.checked;
        //     }
        // })
        cekall.onchange = function () {
            // this.checked  可以得到当前复选框的选中状态，如果是 true 就是选中，如果是 false 就是未选中
            for (var i = 0; i < inp.length; i++) {
                inp[i].checked = this.checked;
            }
        }
        for (var i = 0; i < inp.length; i++) {
            inp[i].onclick = function () {
                // 设置一个变量来控制按钮是否全部选中
                var flag = true;
                // 每次点击下面的复选框都要检查下面的按钮是否被全部选中。
                for (var i = 0; i < inp.length; i++) {
                    if (!inp[i].checked) {
                        flag = false;
                    }
                }
                cekall.checked = flag;
            }
        }
    }

    function trim(s) {
        return s.replace(/^\s+|\s+$/g, '');
    }


    function goPage(pno) {
        var itable = document.getElementById("zhgl_list");
        var num = itable.rows.length; //表格所有行数(所有记录数)
        var totalPage = 0; //总页数
        var pageSize = 10; //每页显示行数
        //总共分几页
        if (num / pageSize > parseInt(num / pageSize)) {
            totalPage = parseInt(num / pageSize) + 1;
        } else {
            totalPage = parseInt(num / pageSize);
        }
        var currentPage = pno; //当前页数
        var startRow = (currentPage - 1) * pageSize + 1; //开始显示的行  31
        var endRow = currentPage * pageSize; //结束显示的行   40
        endRow = (endRow > num) ? num : endRow; //40
        // console.log(endRow);
        //遍历显示数据实现分页
        for (var i = 1; i < (num + 1); i++) {
            var irow = itable.rows[i - 1];
            if (i >= startRow && i <= endRow) {
                irow.style.display = "table-row";
            } else {
                irow.style.display = "none";
            }
        }
        var pageEnd = document.getElementById("pageEnd");
        var tempStr =
            `<li class="page-item"><a class="page-link" href='javascript:void(0);'>共 ${totalPage} 页</a></li>`
        if (currentPage > 1) {
            tempStr += `<li class="page-item"><a class="page-link" href='javascript:goPage(` + (1) +  `);'>首页</a></li>`;
            tempStr += `<li class="page-item"><a class="page-link" href='javascript:goPage(` + (currentPage - 1) +  `);'>上一页</a></li>`;
        } else {}
        for (var pageIndex = 1; pageIndex < totalPage + 1; pageIndex++) {
            if (currentPage == pageIndex) {
                tempStr += `<li class="page-item active"><a class="page-link" href='javascript:goPage(` + pageIndex +  `);' >${pageIndex}</a></li>`;
            } else {
                tempStr += `<li class="page-item "><a class="page-link" href='javascript:goPage(` + (pageIndex) +  `);' >${pageIndex}</a></li>`;
            }
        }
        if (currentPage < totalPage) {
            tempStr += `<li class="page-item"><a class="page-link" href='javascript:goPage(` + (currentPage + 1) +  `);'>下一页</a></li>`;
            tempStr += `<li class="page-item"><a class="page-link" href='javascript:goPage(` + (totalPage) +  `);'>尾页</a></li>`;
        } else {}
        document.getElementById("barcon").innerHTML = tempStr;
    }