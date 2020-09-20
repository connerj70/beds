function getUserIDFromCookies(cookies) {
    let cookiesSplit = cookies.split(";")

    for (let i = 0; i < cookiesSplit.length; i++) {
        if (cookiesSplit[i].includes("user_id")) {
            let userID = cookiesSplit[i].split("=")[1]
            return userID
        }
    }

    return ""
}

export default getUserIDFromCookies