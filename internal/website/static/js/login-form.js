function setLoginFormVisible(visible) {
    if (visible) {
        document.getElementById("login-form").classList.remove("modal-closing")
        document.getElementById("login-form").classList.add("modal-opening")
    } else {
        document.getElementById("login-form").classList.add("modal-closing")
        document.getElementById("login-form").classList.remove("modal-opening")
    }
}