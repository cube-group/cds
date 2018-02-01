/**
 * Created by linyang on 2017/12/19.
 */
function displayError(value) {
    window.alert(value)
}


/**
 * 判断是否请求成功
 *
 * @param data 接口返回数据
 * @returns {boolean}
 */
function isDataSuccess(data) {
    return data.code == '0'
}

/**
 * 显示错误信息
 *
 * @param msg
 * @param delaySecond 几秒后隐藏
 */
function showError(msg, delaySecond) {
    if (delaySecond == undefined) {
        delaySecond = 3;
    }
    $('#ajaxSuccess').hide();
    $('#ajaxError').hide().stop(false, true).css({
        "width": "500px",
        "position": "fixed",
        "top": 0,
        "display": "block",
        "opacity": 0,
        "left": "50%",
        "margin-left": "-250px",
        "z-index": getMaxZIndex() + 1
    }).find('.message').html(msg);
    //定位到3分之一处
    $('#ajaxError').css("margin-top", Math.max(0, ($(window).height() - $('#ajaxError').height()) / 3));
    $('#ajaxError').animate({"opacity": 0.9}).delay(1000 * delaySecond).fadeOut();
}


/**
 * 显示ajax成功信息
 *
 * @param msg
 * @param delaySecond 几秒后隐藏
 */
function showSuccess(msg, delaySecond) {
    if (delaySecond == undefined) {
        delaySecond = 3;
    }
    $('#ajaxError').hide();
    $('#ajaxSuccess').hide().stop(false, true).css({
        "width": "500px",
        "position": "fixed",
        "top": 0,
        "display": "block",
        "opacity": 0,
        "left": "50%",
        "margin-left": "-250px",
        "z-index": getMaxZIndex() + 1
    }).find('.message').html(msg);
    //定位到3分之一处
    $('#ajaxSuccess').css("margin-top", Math.max(0, ($(window).height() - $('#ajaxSuccess').height()) / 3));
    $('#ajaxSuccess').animate({"opacity": 0.9}).delay(1000 * delaySecond).fadeOut();
}


/**
 * 延时1秒跳转到某个地址(为了看到弹出的提示)
 *
 * @param url 要跳转的地址
 * @param delaySecond 延时时间(默认1秒)
 */
function jumpUrl(url, delaySecond) {
    if (delaySecond == undefined) {
        delaySecond = 1;
    }
    setTimeout("window.location.href='" + url + "'", 1000 * delaySecond);
}

/**
 * 延时1秒跳转到某个地址(为了看到弹出的提示)
 *
 * @param url 要跳转的地址
 * @param delaySecond 延时时间(默认1秒)
 */
function delayReload(url, delaySecond) {
    if (delaySecond == undefined) {
        delaySecond = 1;
    }
    setTimeout("window.location.href='" + window.location.href + "'", 1000 * delaySecond);
}

/**
 * 显示接口的成功信息,传了defaultMsg如果msg不存在展示defaultMsg
 *
 * @param data
 * @param defaultMsg
 */
function showDataSuccess(data, defaultMsg) {
    var msg = '';
    if (data.msg != undefined) {
        msg = data.msg;
    } else if (defaultMsg != undefined) {
        msg = defaultMsg
    } else {
        msg = "ok"
    }
    showSuccess(msg);
}

/**
 * 显示接口的失败信息,传了defaultMsg如果msg不存在展示defaultMsg
 *
 * @param data
 * @param defaultMsg
 */
function showDataError(data, defaultMsg) {
    var msg = '';
    if (data.msg != undefined) {
        msg = data.msg;
    } else if (defaultMsg != undefined) {
        msg = defaultMsg
    } else {
        msg = "ok"
    }
    showError(msg);
}


function getMaxZIndex() {
    var maxZ = Math.max.apply(null,
        $.map($('body *'), function (e, n) {
            if ($(e).css('position') != 'static')
                return parseInt($(e).css('z-index')) || 1;
        }));
    return maxZ;
};


/**
 * 设置模态框的值
 *
 * @param prefix 表单名字前缀
 * @param dataObj 存放字段的对象
 * @param attributes 设置的对象属性数组
 */
function setInputValue(prefix, dataObj, attributes) {
    for (var i = 0; i < attributes.length; i++) {
        var attribute = attributes[i];
        var inputId = prefix + ucfirst(attribute);
        $("#" + inputId).val(dataObj.data(attribute))
    }
}

/**
 * 获取模态框的form参数值
 *
 * @param prefix 表单名字前缀
 * @param attributes 获取的对象属性数组
 * @returns {{}}
 */
function getInputParams(prefix, attributes) {
    var params = {};
    for (var i = 0; i < attributes.length; i++) {
        var attribute = attributes[i];
        var inputId = prefix + ucfirst(attribute);
        params[attribute] = $("#" + inputId).val()
    }
    return params
}


/**
 * 成功展示信息后跳转，失败展示信息
 *
 * @param data
 * @param modalId 成功隐藏模态框的ID
 * @param url 跳转的地址，如果没传，当前页面刷新
 */
function showMsgJump(data, modalId, url) {
    if (url == undefined) {
        url = window.location.href
    }
    if (isDataSuccess(data)) {
        showDataSuccess(data);
        $('#' + modalId).modal("hide");
        jumpUrl(url)
    } else {
        showDataError(data);
    }
}


/**
 * 首字母大小
 * @param str
 */
function ucfirst(str) {
    return str.substring(0, 1).toUpperCase() + str.substring(1)
}