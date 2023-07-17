import axios from "axios";

const instance = axios.create({
    baseURL: "http://localhost:8080",
    withCredentials: true
})

instance.interceptors.response.use(function (resp) {
    if (resp.status == 401) {
        window.location.href="/users/login"
    }
    return resp
}, (err) => {
    if (err.response.status == 401) {
        window.location.href="/users/login"
    }
    return err
})

export default instance