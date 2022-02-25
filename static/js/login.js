var loginbtnE = document.getElementById("loginbtn");
var userE = document.getElementById("user");
userE.onchange = verify;
var passE = document.getElementById("pass");
passE.onchange = verify;
var autoE = document.getElementById("check");
var alertE = document.getElementById("alert");
loginbtnE.addEventListener('click', Login);
function verify() {
    if (this.value.trim().length === 0) {
        this.classList.add("error");
        return
    } else {
        this.classList.remove("error");
    }
}
function Login() {
    if (loginbtnE.dataset.state) {
        return;
    }
    loginbtnE.dataset.state = true;

    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");

    var raw = JSON.stringify({
        "Username": userE.value,
        "Password": passE.value,
        "AutoLogin": autoE.checked
    });
    var requestOptions = {
        method: 'POST',
        headers: myHeaders,
        body: raw,
        redirect: 'follow'
    };
    fetch("/luck/login/verify", requestOptions)
        .then(response => response.json())
        .then(result => {
            try {
                alertE.innerText = result.message;
                if (result.code !== 200) {
                    loginbtnE.dataset.state = false;
                    alertE.classList.add("errortxt");
                    return
                }
                alertE.classList.remove("errortxt");
                var tokenData = JSON.stringify(result.data);
                localStorage.setItem("luckToken", tokenData);
                location.href = "/luck/admin";
            } catch (error) {
                loginbtnE.dataset.state = false;
                alertE.innerText = error;
                alertE.classList.add("errortxt");
                console.log(error);
            }
        })
        .catch(error => { console.log('error', error); loginbtnE.dataset.loginState = 0; });
}