var fk_url = "http://www.8219333.com/img2014/swt/img/left_fk.png";
var zx_url = "http://www.8219333.com/img2014/swt/img/left_zx.png";

var left_style = '\
    <style>\
    .nz_left_float{position:fixed; top: 23%; left: 0; }\
    </style>';
var imgsrc = fk_url;
var left_div = '<div class="nz_left_float"><a href="javascript:;" onclick="openZoosUrl()"><img src="'+imgsrc+'" /></a></div>';

document.write(left_style);
document.write(left_div);