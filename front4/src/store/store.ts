import {CurrentRoute, Feed, Post, Query, User, VkAuthURL} from "../types";

//const store: any = {};
/** @ts-ignore */
window.store = {};

export function writeStore(item: any) {
    return
    if (typeof item === "object") {
        const itemFlat: any = {};
        for (const key in item) {
            if (typeof item[key] !== "object") {
                itemFlat[key] = item[key];
            } else {
                writeStore(item[key]);
            }
        }

        //store[item['id']] = itemFlat;
    }
}

/*function fillData() {
    store["query"] = {
        "type": "Query",
        "id": "query",
        "viewer": "VXNlcklEOnsidXNlcklkIjoxMH0",
        "feed": "feed:fake:id"
    }
    store["feed:fake:id"] = {
        "type": "Feed",
        "id": "feed:fake:id",
        "feed": [
            "UG9zdElEOnsicG9zdElkIjoxNjU2MDAyNjgxfQ",
            "UG9zdElEOnsicG9zdElkIjoxNjUzODE2NzYzfQ",
            "UG9zdElEOnsicG9zdElkIjoxNjUzODEzMjg3fQ",
            "UG9zdElEOnsicG9zdElkIjoxNjUzODExNzAyfQ",
        ]
    }
    store["UG9zdElEOnsicG9zdElkIjoxNjU2MDAyNjgxfQ"] = {
        "type": "Post",
        "id": "UG9zdElEOnsicG9zdElkIjoxNjU2MDAyNjgxfQ",
        "text": "Сегодня я ничего не делал!",
        "user": "VXNlcklEOnsidXNlcklkIjoxMH0",
        "date": 1656002681
    }
    store["UG9zdElEOnsicG9zdElkIjoxNjUzODE2NzYzfQ"] = {
        "type": "Post",
        "id": "UG9zdElEOnsicG9zdElkIjoxNjUzODE2NzYzfQ",
        "text": "sdfsf",
        "user": "VXNlcklEOnsidXNlcklkIjoxMH0",
        "date": 1653816763
    }
    store["UG9zdElEOnsicG9zdElkIjoxNjUzODEzMjg3fQ"] = {
        "type": "Post",
        "id": "UG9zdElEOnsicG9zdElkIjoxNjUzODEzMjg3fQ",
        "text": "sdf",
        "user": "VXNlcklEOnsidXNlcklkIjoxMH0",
        "date": 1653813287
    }
    store["UG9zdElEOnsicG9zdElkIjoxNjUzODExNzAyfQ"] = {
        "type": "Post",
        "id": "UG9zdElEOnsicG9zdElkIjoxNjUzODExNzAyfQ",
        "text": "qqqq",
        "user": "VXNlcklEOnsidXNlcklkIjoxMH0",
        "date": 1653811702
    }
    store["VXNlcklEOnsidXNlcklkIjoxMH0"] = {
        "id": "VXNlcklEOnsidXNlcklkIjoxMH0",
        "type": "User",
        "avatar": "https://689809.selcdn.ru/meme-files/avatars/bc5cc7b1c10dcc46919aa5087d027b658c6e530ec6b17bd81208ec2070117927",
        "name": "User 10"
    }
}*/

//fillData();

export function getByID(id: string): (Feed | Post | User | Query) {
    // @ts-ignore
    return window.store[id] || {};
}

export function getByType(type: string): (Feed | Post | User | Query | CurrentRoute | VkAuthURL | null) {
    // @ts-ignore
    for(const [, item] of Object.entries(window.store)) {
        if (item.type === type) {
            return item;
        }
    }

    return null;
}

const waitingCallbacks: (() => void)[] = [];

// @ts-ignore
window.waitingCallbacks = waitingCallbacks;

export function storeOnChanged() {
    for (let callback of waitingCallbacks) {
        callback();
    }
}

export function storeSubscribe(callback: () => void) {
    waitingCallbacks.push(callback);
}

export function storeUnsubscribe(callback: () => void) {
    waitingCallbacks.forEach((item, idx) => {
        if (item == callback) {
            waitingCallbacks.splice(idx, 1);
        }
    })
}
