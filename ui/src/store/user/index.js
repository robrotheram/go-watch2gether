import store from "../index"
export * from './user.actions'
export * from './user.reducer'

export const GetUsername = () => {
    let user = store.getState().user
    console.log(user)
    return user.username
}

export const GetWatcher = () => {
    let user = store.getState().user
    return user
}

export const GetID = () => {
    let user = store.getState().user
    console.log(user)
    return user.id
}
