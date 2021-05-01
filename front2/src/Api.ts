export function api<TReq, TResp>(method: string, args: TReq): Promise<TResp> {
    return new Promise((resolve, reject) => {
        fetch("http://localhost:8000/api/" + method, {
            method: 'POST',
            credentials: 'include',
            body: JSON.stringify(args),
        }).then(r => r.json()).then(r => {
            resolve(r);
        })
    })
}
