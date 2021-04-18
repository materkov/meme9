export function api<TReq, TResp>(service: string, method: string, args: TReq): Promise<TResp> {
    return new Promise((resolve, reject) => {
        fetch("http://127.0.0.1:8000/api/" + service + "/" + method, {
            method: 'POST',
            body: JSON.stringify(args),
        }).then(r => r.json()).then(r => {
            resolve(r);
        })
    })
}
