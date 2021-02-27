$(function() {
    //console.log(csrftoken);

    getLikes();
    //console.log(likeCookies);
    for (let i = 0; i < likeCookies.length; i++) {
        $(`#${likeCookies[i]}`).addClass("liked");
    };

    //　改行
    lineBreak();

    // input change
    $(".genderInp, .kindInp").on('click', (e) => {
        $(`.${$(e.currentTarget).attr("class")}`).prop("checked", false);
        $(e.currentTarget).prop("checked", true);
    });

    // confirm modal open
    $("#openConfirm").on('click', (e) => {
        if (!$(e.currentTarget).hasClass("disabled")) {
            // kind
            if ($('input[name="kinds"]:checked').val() == "iina") {
                $("#confirmKind").text("あったらいいな");
            } else if ($('input[name="kinds"]:checked').val() == "huben") {
                $("#confirmKind").text("不便");
            } else {
                $("#confirmKind").text("困ったこと");
            }
            // content
            $("#confirmCON p").text($("#CON").val());
            // gender
            if ($('input[name="genderSelect"]:checked').val() == "men") {
                $("#confirmGender").text("男性");
            } else {
                $("#confirmGender").text("女性");
            }
            // modal
            $(".modal").removeClass("off");
            $(".modalConConfirm").removeClass("off");
        }
    });

    // modal close
    $(".modal, .closeBtn").on('click', (e) => {
        $(".modal").addClass("off");
        $(".modalCon").addClass("off");
    });

    // ボタンdisabled
    $(".kindInp, .genderInp, #CON").on('input', (ev) => {
        able(ev);
    });

    // good button
    $(".goodBtn").on('click', (e) => {
        let id_ = $(e.currentTarget).attr("id");
        let pk_ = id_.replace("good_", "");
        fetch (`/api/detail/${pk_}`, {
            method: "GET",
        })
        .then(() => {
            if ($(e.currentTarget).hasClass("liked")) {
                $(e.currentTarget).removeClass("liked");
                let after_ = parseInt($(e.currentTarget).prev().text()) - 1;
                $(e.currentTarget).prev().text(after_);
            } else {
                $(e.currentTarget).addClass("liked");
                let after_ = parseInt($(e.currentTarget).prev().text()) + 1;
                $(e.currentTarget).prev().text(after_);
            }
        })
        .catch((error) => {
            console.log(error);
        })
    });

    // post new
    $("#posting").on('click', () => {
        //kind
        if ($('input[name="kinds"]:checked').val() == "iina") {
            kind_ = "あったらいいな";
        } else if ($('input[name="kinds"]:checked').val() == "huben") {
            kind_ = "不便";
        } else {
            kind_ = "困ったこと";
        }
        //content
        content_ = $("#CON").val();
        //gender
        if ($('input[name="genderSelect"]:checked').val() == "men") {
            gend_ = "男性";
        } else {
            gend_ = "女性";
        }
        // data dict
        data = {
            Kind: `${kind_}`,
            Content: `${content_}`,
            Gender: `${gend_}`,
        }
        //fetch post
        fetch ("/create", {
            method: "POST",
            body: JSON.stringify(data),
            headers: {
                "Content-Type": "application/json; charset=utf-8",
            },
        })
        .then((response) => {
            return response.json();
        })
        .then(() => {
            location.href = "/";
        })
        .catch((err) => {
            console.log(err);
        })
    })
});

let able = (e) => {
    //console.log($('input[name="kinds"]:checked').val());
    //console.log($('input[name="genderSelect"]:checked').val());
    //console.log($("#CON").val());
    if ($('input[name="kinds"]:checked').val() != null 
    && $('input[name="genderSelect"]:checked').val() != null 
    && $("#CON").val() != "") {
        $("#openConfirm").removeClass("disabled");
    }
};

let kind_, gend_, content_;
let likeCookies = [];

const getCookie = name => {
    if (document.cookie && document.cookie !== '') {
        for (const cookie of document.cookie.split(';')){
            const [key, value] = cookie.trim().split('=');
            if(key === name) {
                return decodeURIComponent(value);
            }
        }
    }
};

const csrftoken = getCookie('csrftoken');

// get like and push to likeCookies[]
const getLikes = () => {
    if (document.cookie && document.cookie !== '') {
        for (const cookie of document.cookie.split(';')) {
            const [key, value] = cookie.trim().split('=');
            if (key.indexOf('good_') != -1 && value == "1") {
                likeCookies.push(key)
            }
        }
    }
}

// 改行
let lineBreak = () => {
    let contents = $("._con");
    for (let i = 0; i < contents.length; i++) {
        $(contents[i]).html($(contents[i]).text().replace(/\r?\n/g, '<br>'));
    }
}