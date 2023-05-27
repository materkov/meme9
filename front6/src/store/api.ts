export function api<T>(method: string, args: any): Promise<T> {
    return new Promise((resolve, reject) => {
        fetch("/api/" + method, {
            method: 'POST',
            body: JSON.stringify(args),
            headers: {
                'Authorization': 'Bearer ' + localStorage.getItem('authToken'),
            }
        })
            .then(r => r.json())
            .then(r => {
                if (r.error) {
                    reject();
                } else {
                    resolve(r)
                }
            })
            .catch(reject)
    })
}
