function loadLocalToken() {
    this.jwtTokenData = localStorage.getItem("luckToken");
    if (!jwtTokenData) {
        window.jwtToken = null;
        return
    }
    try {
        this.jwtToken = JSON.parse(jwtTokenData);
        if (jwtToken.Token.length === 0) {
            window.jwtToken = null;
            return;
        }
        jwtToken.ExpiresAt = new Date(jwtToken.ExpiresAt);
        window.jwtToken = jwtToken
    } catch (error) {
        window.jwtToken = null;
    }
}
function Auth() {
    if (!window.jwtToken) {
        window.loginState = "未登錄";
        return false;
    }
    try {
        if (new Date() > window.jwtToken.ExpiresAt) return window.loginState = "登錄已過期,請重新登錄", false;
        return true;
    } catch (error) {
        window.loginState = "未登錄";
        return false;
    }
}

(function () {
    loadLocalToken();
    if (!Auth()) {
        document.querySelector("html").innerHTML = window.loginState;
    } else {
        var appCode = '<div class="container"><div div class="content" ><div id="notebooks"><input type="text" id="query" /><div class="selectbox"><div class="total">總共<span id="total">0</span>條記錄</div><select id="orderList"><option value="asc">Newest</option><option value="desc">Oldest</option></select></div><ul id="notebook_ul"></ul></div><div class="genbox card card-2"><div class="genboxtop"><label for="lineid">lineID</label><input id="lineid" maxlength="8" minlength="4" type="text"><div class="bubbly-button" id="genbtn">gen</div></div><div id="genmsg"></div></div></div ></div > ';
        document.body.insertAdjacentHTML("beforeend", appCode);
        var cssE = document.createElement("link");
        cssE.rel = "stylesheet";
        cssE.href = "/luck/static/css/admin.min.css";
        document.head.appendChild(cssE);
        var jsE = document.createElement("script");
        jsE.src = "/luck/static/js/admin.min.js";
        document.body.appendChild(jsE);
    }
})();