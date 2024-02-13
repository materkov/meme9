import {cookieAuthToken, getCookie} from "../utils/cookie";

export function uploadFile(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
        const body = new FormData();
        body.set("file", file);

        fetch("/upload", {
            method: 'POST',
            body: body,
            headers: {
                'authorization': 'Bearer ' + getCookie(cookieAuthToken),
            }
        }).then(r => r.text()).then(r => {
            resolve(r);
        })
    })
}
