import {createStore, applyMiddleware} from 'redux'
import { combineReducers } from 'redux';
import {roomReducer} from './room/room.reducer';

import thunk from 'redux-thunk'
import {composeWithDevTools} from 'redux-devtools-extension'
import { connectRouter } from 'connected-react-router'
import { routerMiddleware } from 'connected-react-router'
import reduxWebsocket from '@giantmachines/redux-websocket';
import {createBrowserHistory} from "history"
export const history  = createBrowserHistory();

const URL = "/api/v1/"

const getWSURL = () => {
    let loc = window.location, new_uri;
    if (loc.protocol === "https:") {
        new_uri = "wss:";
    } else {
        new_uri = "ws:";
    }
    new_uri += "//" + loc.host + URL;
    return new_uri
}

const rootReducer = (history) => combineReducers({
    room: roomReducer,
    router: connectRouter(history),
});

const middleware = [thunk,  routerMiddleware(history), reduxWebsocket()]

export const getStoreFromLocalStore = () => {
    let store =  JSON.parse( sessionStorage.getItem('watch2gether') );
    if (store === null ){
        return {}
    }
    if (store.version !== window.w2g_version) {
        return {}
    }
    if (store.room !== undefined){
        store.room.error = "";
        store.room.active = false;
    }
    return store;
}
const store = createStore(rootReducer(history), getStoreFromLocalStore(), composeWithDevTools(applyMiddleware(...middleware)))
store.subscribe(()=>{
    let store = store.getState()
    store.version = window.w2g_version
    sessionStorage.setItem('watch2gether', JSON.stringify(store.getState()))
})



let apiUrl = URL
let wsUrl = getWSURL();

if ( process.env.NODE_ENV === "development"){
    apiUrl = "http://localhost:8080"+URL
    wsUrl = "ws://localhost:8080"+URL
} 

export const API_URL = apiUrl
export const WS_URL = wsUrl

export default store;
