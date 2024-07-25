function GetQueryString(name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)");
    var r = window.location.search.substr(1).match(reg);
    if(r != null){
        //解决中文乱码
        return decodeURI(r[2]);
    }
    return null;
}


var reRegExpChar = /[\\^$.*+?()[\]{}|]/g,
    reHasRegExpChar = RegExp(reRegExpChar.source);

function escapeRegExp(string) {
    return (string && reHasRegExpChar.test(string))
        ? string.replace(reRegExpChar, '\\$&')
        : string;
}


function getSelectPosition(oTxt) {
    var nullvalue = -1;
    var selectStart; //选中开始位置
    var selectEnd; //选中结束位置
    var position; //焦点位置
    var selectText; //选中内容
    if(oTxt.setSelectionRange) { //非IE浏览器
        selectStart = oTxt.selectionStart;
        selectEnd = oTxt.selectionEnd;
        if(selectStart == selectEnd) {
            position = oTxt.selectionStart;
            selectStart = nullvalue;
            selectEnd = nullvalue;
        } else {
            position = nullvalue;
        }
        selectText = oTxt.value.substring(selectStart, selectEnd);
    } else { //IE
        var range = document.selection.createRange();
        selectText = range.text;
        range.moveStart("character", -oTxt.value.length);
        position = range.text.length;
        selectStart = position - (selectText.length);
        selectEnd = selectStart + (selectText.length);
        if(selectStart != selectEnd) {
            position = nullvalue;
        } else {
            selectStart = nullvalue;
            selectEnd = nullvalue;
        }
    }
    return [selectStart,selectEnd,selectText]
}