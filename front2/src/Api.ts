export function api<TReq, TResp>(method: string, args: TReq): Promise<TResp> {
    return new Promise((resolve, reject) => {
        fetch("/api/" + method, {
            method: 'POST',
            credentials: 'include',
            body: JSON.stringify(args),
        })
            .then(r => {
                if (r.status !== 200) {
                    reject();
                    return
                }

                return r.json()
            })
            .then(r => {
                resolve(r);
            })
            .catch(() => reject())
    })
}
