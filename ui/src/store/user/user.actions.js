import axios from 'axios'
import {BASE_URL, history} from '../index'
import { AUTH_LOGIN } from "./user.types";

export const checklogin = () => {
    return dispatch => {
        axios.get(BASE_URL+`/auth/user`).then(res => {
            console.log("auth", res.data)
            dispatch({
                type: AUTH_LOGIN,
                auth: true,
                username: res.data.user.username,
                icon: res.data.user.avatar_icon,
                guilds: res.data.guilds,
            })
            history.push('/app');
        }).catch(e => {
            console.log("Unathoiriz")
            dispatch({
                type: AUTH_LOGIN,
                auth: false,
                username: "",
                guilds: []
            })
        });
    }
}