function bold_text() {

    document.querySelectorAll(".form-control_content").execCommand('bold');
}
function set_input_bold_typing() {

    var inp = document.querySelectorAll(".form-control_content");
    var control = document.getElementById("control_content");

    // var start = control.selectionStart;

    // var end = control.selectionEnd;
    
    if (start != end) {
        for (var i = start; i < 1000000; i++) {
            if (inp[i].value) inp[i].style.fontWeight = "bold";

            inp[i].addEventListener("input", function () {
                this.style.fontWeight = this.value ? "bold" : "300";
            });
        }
    }
    // var inp = document.getElementById("control_content");


    // var start = control.selectionStart;

    // var end = control.selectionEnd;
    // alert(`hello ${inp.value[start].style}, vvv ${inp.value[start]}`);


    // if (start != end) {
    //     alert(`hello ${start}, ${end}`);

    //     for (var i = start; i < end; i++) {
    //         alert(`VTOROI  ${inp.value[i]}`);
    //         inp.value[i].addEventListener("input", function () {
    //             if (inp.value[i]) {
    //                 alert("dikii");
    //                 inp.value[i].style.fontWeight = "bold";
    //             }
    //             this.style.fontWeight = this.value ? "bold" : "300";
    //         });
    //     }
    // }
}


function addTag(open, close) {

    var control = document.getElementById("control_content");

    var start = control.selectionStart;

    var end = control.selectionEnd;

    if (start != end) {
        var text = document.getElementById('control_content').value;
        alert(text.substring(start, end));
        var text1 = document.getElementById('control_content')
        text1.execCommand('bold');
        // control.innerHTML = (text.substring(0, start) + open + text.substring(start, end).execCommand("bold") + close + text.substring(end));

        control.focus();

        var sel = end + (open + close).length;

        control.setSelectionRange(sel, sel);
    }

    return false;
}
// Жирный
function bold_button() {
    return addTag('<b>', '</b>');
}

// Курсив
$('#button-i').click(function () {
    return addTag('<i>', '</i>');
});

// Подчеркнутый
$('#button-u').click(function () {
    return addTag('<u>', '</u>');
});

// Зачеркнутый
$('#button-s').click(function () {
    return addTag('<strike>', '</strike>');
});

// Ссылка
$('#button-a').click(function () {
    return addTag('<a href="' + prompt('Введите адрес', '') + '">', '</a>');
});

// При клике на кнопки не снимаем фокус с textarea.
$('a').on('mousedown', function () {
    return false;
});